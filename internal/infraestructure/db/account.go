package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lalizita/shard-banking/internal/config"
)

type AccountDB struct {
	Pool *pgxpool.Pool
}

func NewAccountDB(ctx context.Context, cfg *config.Config) (*AccountDB, error) {
	pool, err := pgxpool.New(ctx, cfg.DBAccountDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create account db pool: %w", err)
	}
	return &AccountDB{Pool: pool}, nil
}
