package usecase

import (
	"context"
	"encoding/json"
	"time"

	"order-service/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type OrderUseCase struct {
	repo          OrderRepository
	paymentClient PaymentClient
	cache         *redis.Client
}

func NewOrderUseCase(
	r OrderRepository,
	p PaymentClient,
	c *redis.Client,
) *OrderUseCase {
	return &OrderUseCase{
		repo:          r,
		paymentClient: p,
		cache:         c,
	}
}

func (uc *OrderUseCase) CreateOrder(
	ctx context.Context,
	customerID,
	itemName string,
	amount int64,
) (*domain.Order, error) {

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

		uc.cache.Del(ctx, "order:"+order.ID)

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

	uc.cache.Del(ctx, "order:"+order.ID)

	return order, nil
}

func (uc *OrderUseCase) GetOrder(
	ctx context.Context,
	id string,
) (*domain.Order, error) {

	val, err := uc.cache.Get(ctx, "order:"+id).Result()
	if err == nil {
		var order domain.Order
		if json.Unmarshal([]byte(val), &order) == nil {
			return &order, nil
		}
	}

	order, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(order)
	// TTL
	uc.cache.Set(
		ctx,
		"order:"+id,
		data,
		5*time.Minute,
	)

	return order, nil
}

func (uc *OrderUseCase) CancelOrder(
	ctx context.Context,
	id string,
) error {

	order, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status != domain.StatusPending {
		return domain.ErrCannotCancelOrder
	}

	order.Status = domain.StatusCancelled

	if err := uc.repo.Update(order); err != nil {
		return err
	}

	uc.cache.Del(ctx, "order:"+order.ID)

	return nil
}

func (uc *OrderUseCase) GetRecentOrders(
	limit int,
) ([]*domain.Order, error) {
	return uc.repo.GetRecent(limit)
}

func (uc *OrderUseCase) WatchOrder(
	orderID string,
) <-chan domain.Order {

	ch := make(chan domain.Order)

	go func() {
		for {
			order, _ := uc.repo.GetByID(orderID)

			ch <- *order

			time.Sleep(1 * time.Second)
		}
	}()

	return ch
}
