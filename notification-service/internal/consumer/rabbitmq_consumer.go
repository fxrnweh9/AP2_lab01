package consumer

import (
	"encoding/json"
	"log"
	"notification-service/internal/handler"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QueueName = "payment.completed"

type RabbitMQConsumer struct {
	conn              *amqp.Connection
	channel           *amqp.Channel
	handler           *handler.NotificationHandler
	processedEventsMu sync.Mutex
	processedEvents   map[string]bool // idempotency store (in-memory)
}

func NewRabbitMQConsumer(amqpURL string, h *handler.NotificationHandler) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		QueueName,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:            conn,
		channel:         ch,
		handler:         h,
		processedEvents: make(map[string]bool),
	}, nil
}

func (c *RabbitMQConsumer) Start() error {
	msgs, err := c.channel.Consume(
		QueueName,
		"",    // consumer tag
		false, // auto-ack = FALSE (manual ACK!)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("[Consumer] Waiting for messages...")

	for msg := range msgs {
		var event handler.PaymentEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("[Consumer] Failed to parse message: %v", err)
			msg.Nack(false, false)
			continue
		}

		// Idempotency check
		c.processedEventsMu.Lock()
		if c.processedEvents[event.EventID] {
			c.processedEventsMu.Unlock()
			log.Printf("[Consumer] Duplicate event_id=%s, skipping", event.EventID)
			msg.Ack(false)
			continue
		}
		c.processedEvents[event.EventID] = true
		c.processedEventsMu.Unlock()

		// Обрабатываем
		if err := c.handler.Handle(event); err != nil {
			log.Printf("[Consumer] Error handling event: %v", err)
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
	}

	return nil
}

func (c *RabbitMQConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
