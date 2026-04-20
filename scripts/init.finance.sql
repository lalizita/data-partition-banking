-- ENUM for operation types
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE transaction_entry_type AS ENUM ('CREDIT', 'DEBIT');
CREATE TYPE transaction_type AS ENUM ('PIX', 'BANK_TRANSFER');
CREATE TYPE transaction_status AS ENUM ('INITIALIZED', 'PENDING', 'COMPLETED', 'FAILED');
CREATE TYPE fraud_analysis_result AS ENUM ('APPROVED', 'BLOCKED', 'MANUAL_REVIEW');

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status transaction_status NOT NULL DEFAULT 'INITIALIZED',
    shard_id SMALLINT NOT NULL,
    entry_type transaction_entry_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tx_account_sequence
ON transactions (client_id, sequence);

CREATE TABLE IF NOT EXISTS fraud_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    fraud_score SMALLINT NOT NULL,
    rules_triggred JSONB NOT NULL DEFAULT '[]',
    result fraud_analysis_result NOT NULL DEFAULT 'MANUAL_REVIEW',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

