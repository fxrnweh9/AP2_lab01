package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"order-service/internal/domain"
)

type PostgresOrderRepo struct {
	db *sql.DB
}

func NewPostgresOrderRepo(db *sql.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}

func (r *PostgresOrderRepo) Create(o *domain.Order) error {
	query := `INSERT INTO orders (customer_id, item_name, amount, status, created_at) 
              VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`
	return r.db.QueryRow(query, o.CustomerID, o.ItemName, o.Amount, o.Status).Scan(&o.ID, &o.CreatedAt)
}

func (r *PostgresOrderRepo) GetByID(id string) (*domain.Order, error) {
	o := &domain.Order{}
	query := `SELECT id, customer_id, item_name, amount, status, created_at FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt)
	return o, err
}

func (r *PostgresOrderRepo) UpdateStatus(id string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}
