package consumer

import (
	"context"
	"encoding/json"
	"log"
	"notification-service/internal/handler"
	"notification-service/internal/worker"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

const QueueName = "payment.completed"

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	handler *handler.NotificationHandler
	redis   *redis.Client
}

func NewRabbitMQConsumer(
	amqpURL string,
	h *handler.NotificationHandler,
	redisClient *redis.Client,
) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(QueueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
		handler: h,
		redis:   redisClient,
	}, nil
}

func (c *RabbitMQConsumer) Start() error {
	msgs, err := c.channel.Consume(QueueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Println("[Consumer] Waiting for messages... (parallel workers, prefetch=10)")

	sem := make(chan struct{}, 10)

	for msg := range msgs {
		msg := msg // capture for goroutine

		var event handler.PaymentEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("[Consumer] Bad message: %v", err)
			msg.Nack(false, false)
			continue
		}

		sem <- struct{}{}

		go func() {
			defer func() { <-sem }()
			ctx := context.Background()
			key := "notified:" + event.EventID

			// Redis idempotency check
			exists, _ := c.redis.Exists(ctx, key).Result()
			if exists > 0 {
				log.Printf("[Worker] Duplicate event_id=%s, skipping", event.EventID)
				msg.Ack(false)
				return
			}

			err := worker.WithRetry(5, func() error {
				return c.handler.Handle(event)
			})

			if err != nil {
				log.Printf("[Worker] All 5 retries failed for event_id=%s: %v", event.EventID, err)
				msg.Nack(false, false)
				return
			}

			c.redis.Set(ctx, key, "done", 24*time.Hour)
			msg.Ack(false)
			log.Printf("[Worker] Successfully processed event_id=%s", event.EventID)
		}()
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
