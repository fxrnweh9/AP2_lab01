package domain

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "Pending"
	StatusPaid      OrderStatus = "Paid"
	StatusFailed    OrderStatus = "Failed"
	StatusCancelled OrderStatus = "Cancelled"
)

type Order struct {
	ID         string
	CustomerID string
	ItemName   string
	Amount     int64
	Status     OrderStatus
	CreatedAt  time.Time
}

func (o *Order) Validate() error {
	if o.Amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}
