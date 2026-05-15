package handler

import (
	"fmt"
	"log"
)

type PaymentEvent struct {
	EventID       string `json:"event_id"`
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	CustomerEmail string `json:"customer_email"`
	Status        string `json:"status"`
}

type NotificationHandler struct {
	// в реальной жизни — email client, но здесь просто лог
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) Handle(event PaymentEvent) error {
	email := event.CustomerEmail
	if email == "" {
		email = "user@example.com"
	}

	amountDollars := fmt.Sprintf("$%.2f", float64(event.Amount)/100)

	log.Printf("[Notification] Sent email to %s for Order #%s. Amount: %s. Status: %s",
		email,
		event.OrderID,
		amountDollars,
		event.Status,
	)

	return nil
}
