package handler

import (
	"fmt"
	"notification-service/internal/provider"
)

type PaymentEvent struct {
	EventID       string `json:"event_id"`
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	CustomerEmail string `json:"customer_email"`
	Status        string `json:"status"`
}

type NotificationHandler struct {
	sender provider.EmailSender
}

func NewNotificationHandler(sender provider.EmailSender) *NotificationHandler {
	return &NotificationHandler{sender: sender}
}

func (h *NotificationHandler) Handle(event PaymentEvent) error {
	email := event.CustomerEmail
	if email == "" {
		email = "user@example.com"
	}

	amount := fmt.Sprintf("$%.2f", float64(event.Amount)/100)
	subject := fmt.Sprintf("Payment %s for Order #%s", event.Status, event.OrderID)
	body := fmt.Sprintf(
		"Your payment of %s for Order #%s has been %s.",
		amount, event.OrderID, event.Status,
	)

	return h.sender.Send(email, subject, body)
}
