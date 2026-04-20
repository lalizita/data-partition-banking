package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DBShard0DSN  string `env:"DB_SHARD_0_DSN"`
	DBShard1DSN  string `env:"DB_SHARD_1_DSN"`
	DBAccountDSN string `env:"DB_ACCOUNT_DSN"`

	TransactionShards string `env:"TRANSACTION_SHARDS_UNITS"`
}

func Load(ctx context.Context) (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	return &cfg, nil
}
