package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/consumer"
	"notification-service/internal/handler"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	h := handler.NewNotificationHandler()

	c, err := consumer.NewRabbitMQConsumer(rabbitURL, h)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	// Graceful Shutdown
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
