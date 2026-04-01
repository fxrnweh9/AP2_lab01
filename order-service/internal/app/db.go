package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pgxpool.New(ctx, connString)
}
