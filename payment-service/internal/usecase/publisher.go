package usecase

type EventPublisher interface {
	PublishPaymentCompleted(event PaymentCompletedEvent) error
}

type PaymentCompletedEvent struct {
	EventID       string
	OrderID       string
	Amount        int64
	CustomerEmail string
	Status        string
}
