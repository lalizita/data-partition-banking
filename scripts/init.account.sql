CREATE TYPE account_status AS ENUM ('ACTIVE', 'SUSPENDED', 'CLOSED');

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    status account_status NOT NULL DEFAULT 'ACTIVE',
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0,
    daily_limit NUMERIC(15, 2) NOT NULL DEFAULT 1000,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Table for routing clients to shards
CREATE TABLE IF NOT EXISTS clients_shard_routing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES accounts(id),
    transaction_shard_id SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
