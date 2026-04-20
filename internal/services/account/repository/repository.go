package repository

import (
	"context"

	"github.com/lalizita/shard-banking/internal/infraestructure/db"
	"github.com/lalizita/shard-banking/internal/services/account/model"
	"github.com/lalizita/shard-banking/pkg/sharding"
)

type IAccountRepository interface {
	CreateAccount(ctx context.Context, account model.Account) (model.Account, error)
	CreateClientShard(ctx context.Context, client model.ClientShardRouting) error
}

type AccountRepository struct {
	db             *db.AccountDB
	shardingRouter sharding.IShardRouter
}

func NewAccountRepository(db *db.AccountDB,
	shardRouter sharding.IShardRouter) *AccountRepository {
	return &AccountRepository{
		db:             db,
		shardingRouter: shardRouter,
	}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account model.Account) (model.Account, error) {
	sqlStatement := `
		INSERT INTO accounts (name, email, status, balance, daily_limit)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var created model.Account
	err := r.db.Pool.QueryRow(ctx, sqlStatement,
		account.Name, account.Email, account.Status, account.Balance, account.DailyLimit,
	).Scan(
		&created.ID,
	)
	if err != nil {
		return model.Account{}, err
	}

	return created, nil
}

func (r *AccountRepository) CreateClientShard(ctx context.Context, client model.ClientShardRouting) error {
	client.ShardID = r.shardingRouter.RouteForClientID(client.ClientID)

	sqlStatement := `
		INSERT INTO clients_shard_routing (client_id, transaction_shard_id)
		VALUES ($1, $2)
	`
	_, err := r.db.Pool.Exec(ctx, sqlStatement, client.ClientID, client.ShardID)
	return err
}
