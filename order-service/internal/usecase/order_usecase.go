package usecase

import (
	"context"
	"time"

	"order-service/internal/domain"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	repo          OrderRepository
	paymentClient PaymentClient
}

func NewOrderUseCase(r OrderRepository, p PaymentClient) *OrderUseCase {
	return &OrderUseCase{
		repo:          r,
		paymentClient: p,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, customerID, itemName string, amount int64) (*domain.Order, error) {
	order := &domain.Order{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(order); err != nil {
		return nil, err
	}

	status, err := uc.paymentClient.ProcessPayment(order.ID, order.Amount)
	if err != nil {
		order.Status = domain.StatusFailed
		_ = uc.repo.Update(order)
		return nil, err
	}

	if status == "Authorized" {
		order.Status = domain.StatusPaid
	} else {
		order.Status = domain.StatusFailed
	}

	if err := uc.repo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *OrderUseCase) GetOrder(id string) (*domain.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *OrderUseCase) CancelOrder(id string) error {
	order, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status != domain.StatusPending {
		return domain.ErrCannotCancelOrder
	}

	order.Status = domain.StatusCancelled
	return uc.repo.Update(order)
}
