package repository

import (
	"context"

	"github.com/lalizita/shard-banking/internal/infraestructure/db"
	"github.com/lalizita/shard-banking/internal/services/transaction/model"
)

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, tx model.Transaction) error
}

type TransactionRepository struct {
	db *db.FinanceDB
}

func NewTransactionRepository(db *db.FinanceDB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx model.Transaction) error {
	shardID := tx.ShardID
	pool := r.db.Shards[shardID]
	// Ainda falta alguns campos para serem preenchidos mas dependem da construção de mais funcionalidades
	sqlStatement := `
		INSERT INTO transactions (client_id, amount, entry_type, status, shard_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := pool.Exec(ctx, sqlStatement, tx.ClientID, tx.Amount, tx.EntryType, tx.Status, tx.ShardID)
	if err != nil {
		return err
	}

	return nil
}
