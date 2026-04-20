package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "ACTIVE"
	AccountStatusSuspended AccountStatus = "SUSPENDED"
	AccountStatusClosed    AccountStatus = "CLOSED"

	InitialDailyLimit float64 = 0
)

type Account struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Email      string        `json:"email"`
	Status     AccountStatus `json:"status"`
	Balance    float64       `json:"balance"`
	DailyLimit float64       `json:"daily_limit"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type CreateAccountRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateAccountResponse struct {
	ID string `json:"client_id"`
}

type ClientShardRouting struct {
	ClientID uuid.UUID
	ShardID  int
}
