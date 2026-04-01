package usecase

import (
	"payment-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	repo PaymentRepository
}

func NewPaymentUseCase(repo PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: repo}
}

func (uc *PaymentUseCase) ProcessPayment(orderID string, amount int64) (*domain.Payment, error) {
	if amount > 100000 {
		p := &domain.Payment{
			ID:        uuid.New().String(),
			OrderID:   orderID,
			Amount:    amount,
			Status:    domain.StatusDeclined,
			CreatedAt: time.Now(),
		}
		_ = uc.repo.Create(p)
		return p, domain.ErrPaymentLimitExceeded
	}

	p := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		Amount:        amount,
		Status:        domain.StatusAuthorized,
		TransactionID: uuid.New().String(),
		CreatedAt:     time.Now(),
	}

	err := uc.repo.Create(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (uc *PaymentUseCase) GetPayment(orderID string) (*domain.Payment, error) {
	return uc.repo.GetByOrderID(orderID)
}
