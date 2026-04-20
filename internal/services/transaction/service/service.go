package service

import (
	"context"

	"github.com/lalizita/shard-banking/internal/services/transaction/model"
	"github.com/lalizita/shard-banking/internal/services/transaction/repository"
	"github.com/lalizita/shard-banking/pkg/sharding"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, tx model.Transaction) error
}

type TransactionServiceImpl struct {
	repo        repository.ITransactionRepository
	shardRouter sharding.IShardRouter
}

func NewTransactionService(repo repository.ITransactionRepository) *TransactionServiceImpl {
	return &TransactionServiceImpl{repo: repo}
}

func (s *TransactionServiceImpl) CreateTransaction(ctx context.Context, tx model.Transaction) error {
	tx.Status = model.TransactionStatusInitialized
	return s.repo.CreateTransaction(ctx, tx)
}
