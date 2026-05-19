package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/consumer"
	"notification-service/internal/handler"
	"notification-service/internal/provider"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	_ = godotenv.Load()

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"), // "redis:6379"
	})

	var sender provider.EmailSender
	switch os.Getenv("PROVIDER_MODE") {
	case "REAL":
		sender = provider.NewSMTPEmailSender(
			os.Getenv("SMTP_HOST"),
			os.Getenv("SMTP_PORT"),
			os.Getenv("SMTP_USER"),
			os.Getenv("SMTP_PASS"),
			os.Getenv("SMTP_FROM"),
		)
		log.Println("[Notification] Using REAL SMTP provider")
	default:
		sender = provider.NewSimulatedEmailSender()
		log.Println("[Notification] Using SIMULATED provider")
	}

	h := handler.NewNotificationHandler(sender)

	c, err := consumer.NewRabbitMQConsumer(
		os.Getenv("RABBITMQ_URL"),
		h,
		redisClient,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := c.Start(); err != nil {
			log.Printf("Consumer stopped: %v", err)
		}
	}()

	log.Println("[Notification Service] Running. Press Ctrl+C to stop.")
	<-quit
	log.Println("[Notification Service] Shutting down gracefully...")
}
