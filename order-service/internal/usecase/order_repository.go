package usecase

import "order-service/internal/domain"

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	Update(order *domain.Order) error

	GetRecent(limit int) ([]*domain.Order, error)
}
