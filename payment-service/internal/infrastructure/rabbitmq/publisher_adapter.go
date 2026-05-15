package rabbitmq

import "payment-service/internal/usecase"

type PublisherAdapter struct {
	producer *Producer
}

func NewPublisherAdapter(producer *Producer) *PublisherAdapter {
	return &PublisherAdapter{producer: producer}
}

func (a *PublisherAdapter) PublishPaymentCompleted(event usecase.PaymentCompletedEvent) error {
	return a.producer.PublishPaymentCompleted(PaymentEvent{
		EventID:       event.EventID,
		OrderID:       event.OrderID,
		Amount:        event.Amount,
		CustomerEmail: event.CustomerEmail,
		Status:        event.Status,
	})
}
