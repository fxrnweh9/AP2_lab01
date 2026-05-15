package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QueueName = "payment.completed"

type PaymentEvent struct {
	EventID       string `json:"event_id"`
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	CustomerEmail string `json:"customer_email"`
	Status        string `json:"status"`
}

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(amqpURL string) (*Producer, error) {
	conn, err := amqp.Dial(amqpURL) // RabbitMQ
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Durable queue
	_, err = ch.QueueDeclare(
		QueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Producer{conn: conn, channel: ch}, nil
}

func (p *Producer) PublishPaymentCompleted(event PaymentEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(ctx,
		"",        // default exchange
		QueueName, // routing key = queue name
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("[RabbitMQ] Published event for order_id=%s status=%s", event.OrderID, event.Status)
	return nil
}

func (p *Producer) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
