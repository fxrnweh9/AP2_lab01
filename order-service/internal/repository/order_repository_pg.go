package repository

import (
	"context"
	"order-service/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	query := `
		INSERT INTO orders (id, customer_id, item_name, amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, query,
		order.ID,
		order.CustomerID,
		order.ItemName,
		order.Amount,
		order.Status,
		order.CreatedAt,
	)

	return err
}

func (r *OrderRepository) GetByID(id string) (*domain.Order, error) {
	query := `
		SELECT id, customer_id, item_name, amount, status, created_at
		FROM orders
		WHERE id=$1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRow(ctx, query, id)

	var order domain.Order

	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.ItemName,
		&order.Amount,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) Update(order *domain.Order) error {
	query := `
		UPDATE orders
		SET status=$1
		WHERE id=$2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, query, order.Status, order.ID)
	return err
}

func (r *OrderRepository) GetRecent(limit int) ([]*domain.Order, error) {
	query := `
		SELECT id, customer_id, item_name, amount, status, created_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {
		var o domain.Order
		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.ItemName,
			&o.Amount,
			&o.Status,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}
