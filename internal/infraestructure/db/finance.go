package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lalizita/shard-banking/internal/config"
)

type FinanceDB struct {
	Shards map[int]*pgxpool.Pool
}

func NewFinanceDB(ctx context.Context, cfg *config.Config) (*FinanceDB, error) {
	shard0, err := pgxpool.New(ctx, cfg.DBShard0DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}
	shard1, err := pgxpool.New(ctx, cfg.DBShard1DSN)
	if err != nil {
		shard0.Close()
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}
	if err := shard0.Ping(ctx); err != nil {
		shard0.Close()
		shard1.Close()
		return nil, fmt.Errorf("finance shard 0 unreachable: %w (start databases with: docker compose up -d)", err)
	}
	if err := shard1.Ping(ctx); err != nil {
		shard0.Close()
		shard1.Close()
		return nil, fmt.Errorf("finance shard 1 unreachable: %w (start databases with: docker compose up -d)", err)
	}
	return &FinanceDB{
		Shards: map[int]*pgxpool.Pool{
			0: shard0,
			1: shard1,
		},
	}, nil
}
