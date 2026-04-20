package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/lalizita/shard-banking/internal/services/account/model"
	"github.com/lalizita/shard-banking/internal/services/account/repository"
)

type AccountServiceImpl struct {
	repo repository.IAccountRepository
}

type IAccountService interface {
	CreateAccount(ctx context.Context, account model.Account) (model.Account, error)
}

func NewAccountService(repo repository.IAccountRepository) *AccountServiceImpl {
	return &AccountServiceImpl{repo: repo}
}

func (s *AccountServiceImpl) CreateAccount(ctx context.Context, account model.Account) (model.Account, error) {
	account.Balance = 0
	account.Status = model.AccountStatusActive
	account.DailyLimit = model.InitialDailyLimit

	created, err := s.repo.CreateAccount(ctx, account)
	if err != nil {
		return model.Account{}, err
	}

	accID, err := uuid.Parse(created.ID)
	if err != nil {
		return model.Account{}, err
	}

	clientShard := model.ClientShardRouting{
		ClientID: accID,
	}
	err = s.repo.CreateClientShard(ctx, clientShard)
	if err != nil {
		return model.Account{}, err
	}

	return created, nil
}
