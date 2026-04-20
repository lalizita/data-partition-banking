package model

import "time"

type TransactionStatus string
type TransactionType string
type TransactionEntryType string

const (
	TransactionStatusInitialized TransactionStatus = "INITIALIZED"
	TransactionStatusPending     TransactionStatus = "PENDING"
	TransactionStatusCompleted   TransactionStatus = "COMPLETED"
	TransactionStatusFailed      TransactionStatus = "FAILED"

	TransactionTypePix          TransactionType = "PIX"
	TransactionTypeBankTransfer TransactionType = "BANK_TRANSFER"

	EntryTypeCredit TransactionEntryType = "CREDIT"
	EntryTypeDebit  TransactionEntryType = "DEBIT"
)

type Transaction struct {
	ID            string               `json:"id"`
	ClientID      string               `json:"client_id"`
	Amount        float64              `json:"amount"`
	BalanceBefore float64              `json:"balance_before"`
	BalanceAfter  float64              `json:"balance_after"`
	Sequence      int64                `json:"sequence"`
	Status        TransactionStatus    `json:"status"`
	ShardID       int                  `json:"shard_id"`
	CreatedAt     time.Time            `json:"created_at"`
	EntryType     TransactionEntryType `json:"entry_type"`
}

type CreateTransactionRequest struct {
	ClientID string  `json:"client_id"`
	Amount   float64 `json:"amount"`
}
