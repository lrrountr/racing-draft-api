package model

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBClient struct {
	pool *pgxpool.Pool
}

func (c DBClient) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	tx, err = c.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start TX: %w", err)
	}
	return tx, nil
}
