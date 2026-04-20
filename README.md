# Partition Study Case — Shard Banking

Projeto de estudo sobre **sharding de banco de dados** aplicado a um sistema bancário simplificado. A ideia central é demonstrar como distribuir dados de transações financeiras em múltiplos bancos de dados (shards) usando o ID do cliente como chave de particionamento.

---

## Visão Geral da Arquitetura

O sistema é composto por dois microsserviços independentes que compartilham uma estratégia de roteamento baseada em hash consistente.

```mermaid
graph TD
    Client["Cliente HTTP"]

    subgraph Services["Microsserviços"]
        AccountSvc["Account Service\n:8080"]
        FinanceSvc["Finance Service\n:8081"]
    end

    subgraph Databases["Bancos de Dados"]
        AccountDB[("account_db\n:5434")]
        Shard0[("finance_db_0\nShard 0\n:5432")]
        Shard1[("finance_db_1\nShard 1\n:5433")]
    end

    Client -->|"POST /accounts"| AccountSvc
    Client -->|"POST /transactions/credit"| FinanceSvc

    AccountSvc -->|"Cria conta"| AccountDB
    AccountSvc -->|"Armazena roteamento\nclient_id → shard_id"| AccountDB

    FinanceSvc -->|"hash(client_id) % 2 == 0"| Shard0
    FinanceSvc -->|"hash(client_id) % 2 == 1"| Shard1
```

### Como o sharding funciona

Quando um cliente é criado, o sistema calcula em qual shard suas transações serão armazenadas:

```
shard_id = FNV-32a(client_uuid) % TRANSACTION_SHARDS_UNITS
```

O mapeamento `client_id → shard_id` é persistido na tabela `clients_shard_routing` do banco de contas. Assim, ao processar uma transação, o serviço financeiro roteia o INSERT direto para o banco correto sem broadcast.

---

## Estrutura do Projeto

```
.
├── cmd/
│   ├── account/main.go       # Entrypoint do Account Service
│   └── finance/main.go       # Entrypoint do Finance Service
├── internal/
│   ├── config/               # Leitura de variáveis de ambiente
│   ├── http/handlers/        # Registro de rotas HTTP (Echo)
│   ├── infraestructure/db/   # Pools de conexão PostgreSQL
│   └── services/
│       ├── account/          # Domain: conta (model, repo, service, handler)
│       └── transaction/      # Domain: transação (model, repo, service, handler)
├── pkg/
│   └── sharding/router.go    # Algoritmo de hash para roteamento de shard
├── scripts/
│   ├── init.account.sql      # Schema do banco de contas
│   └── init.finance.sql      # Schema replicado em cada shard
└── docker-compose.yaml       # 3 instâncias PostgreSQL
```

---

## Serviço: Account

Responsável pelo cadastro de clientes e pela decisão de qual shard cada cliente pertence.

### Endpoints

| Método | Rota        | Descrição                        |
|--------|-------------|----------------------------------|
| GET    | `/accounts` | Health check                     |
| POST   | `/accounts` | Cria uma nova conta              |

**POST /accounts — Request**
```json
{
  "name": "Maria Silva",
  "email": "maria@exemplo.com"
}
```

**POST /accounts — Response 201**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Maria Silva",
  "email": "maria@exemplo.com",
  "status": "ACTIVE",
  "balance": 0,
  "daily_limit": 0,
  "created_at": "2026-04-20T10:00:00Z"
}
```

### Fluxo de criação de conta

```mermaid
sequenceDiagram
    participant C as Cliente
    participant AS as Account Service
    participant SR as ShardRouter
    participant DB as account_db

    C->>AS: POST /accounts {name, email}
    AS->>DB: INSERT INTO accounts → retorna id (UUID)
    AS->>SR: RouteForClientID(client_id)
    SR-->>AS: shard_id = FNV32a(uuid) % 2
    AS->>DB: INSERT INTO clients_shard_routing (client_id, shard_id)
    AS-->>C: 201 Created {account}
```

### Armazenamento de dados — account_db

```mermaid
erDiagram
    accounts {
        UUID id PK
        VARCHAR name
        VARCHAR email
        account_status status
        NUMERIC balance
        NUMERIC daily_limit
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    clients_shard_routing {
        UUID id PK
        UUID client_id FK
        SMALLINT transaction_shard_id
        TIMESTAMP created_at
    }

    accounts ||--o{ clients_shard_routing : "tem roteamento"
```

- **`accounts`** — dados do cliente: saldo, limite diário, status (`ACTIVE`, `SUSPENDED`, `CLOSED`)
- **`clients_shard_routing`** — tabela de lookup: persiste qual shard armazena as transações de cada cliente

---

## Serviço: Finance

Responsável por registrar transações financeiras, roteando cada operação diretamente ao shard correto do cliente.

### Endpoints

| Método | Rota                    | Descrição                       |
|--------|-------------------------|---------------------------------|
| POST   | `/transactions/credit`  | Registra uma transação de crédito |

**POST /transactions/credit — Request**
```json
{
  "client_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 250.00
}
```

**POST /transactions/credit — Response 201**
```
201 Created
```

### Fluxo de criação de transação

```mermaid
sequenceDiagram
    participant C as Cliente
    participant FS as Finance Service
    participant SR as ShardRouter
    participant S0 as finance_db_0 (Shard 0)
    participant S1 as finance_db_1 (Shard 1)

    C->>FS: POST /transactions/credit {client_id, amount}
    FS->>SR: RouteForClientID(client_id)
    SR-->>FS: shard_id = FNV32a(uuid) % 2

    alt shard_id == 0
        FS->>S0: INSERT INTO transactions
    else shard_id == 1
        FS->>S1: INSERT INTO transactions
    end

    FS-->>C: 201 Created
```

### Armazenamento de dados — finance_db_0 / finance_db_1

O mesmo schema é replicado em cada shard. Cada banco armazena apenas as transações dos clientes que foram roteados para ele.

```mermaid
erDiagram
    transactions {
        UUID id PK
        UUID client_id
        DECIMAL amount
        transaction_status status
        SMALLINT shard_id
        transaction_entry_type entry_type
        TIMESTAMP created_at
    }

    fraud_analysis {
        UUID id PK
        UUID transaction_id FK
        SMALLINT fraud_score
        JSONB rules_triggered
        fraud_analysis_result result
        TIMESTAMP created_at
    }

    transactions ||--o| fraud_analysis : "analisada por"
```

- **`transactions`** — registro financeiro: valor, tipo (`CREDIT`, `DEBIT`), status (`INITIALIZED`, `PENDING`, `COMPLETED`, `FAILED`), shard de origem
- **`fraud_analysis`** — resultado de análise antifraude: score, regras disparadas, resultado (`APPROVED`, `BLOCKED`, `MANUAL_REVIEW`)

---

## Algoritmo de Sharding

Implementado em [pkg/sharding/router.go](pkg/sharding/router.go).

```go
func (s *ShardRouter) RouteForClientID(clientID uuid.UUID) int {
    h := fnv.New32a()
    h.Write([]byte(clientID.String()))
    return int(h.Sum32()) % s.ShardUnits
}
```

**Por que FNV-32a?**
- Distribuição uniforme entre shards
- Determinístico: o mesmo `client_id` sempre vai para o mesmo shard
- Baixo custo computacional

---

## Como executar

**Pré-requisitos:** Docker, Go 1.22+

```bash
# Subir os 3 bancos PostgreSQL
make docker-up

# Iniciar o Account Service (porta 8080)
make run.account

# Iniciar o Finance Service (porta 8081)
make run.finance
```

### Variáveis de ambiente (`.env`)

```env
PORT=8080
PORT_FINANCE=8081

DB_ACCOUNT_DSN=postgres://postgres:postgres@localhost:5434/account_db
DB_SHARD_0_DSN=postgres://postgres:postgres@localhost:5432/finance_db_0
DB_SHARD_1_DSN=postgres://postgres:postgres@localhost:5433/finance_db_1

TRANSACTION_SHARDS_UNITS=2
```

---

## Exemplo de uso

```bash
# 1. Criar conta
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"name": "Maria Silva", "email": "maria@exemplo.com"}'

# Resposta: {"id": "<uuid>", "status": "ACTIVE", ...}

# 2. Registrar transação de crédito (usando o id retornado)
curl -X POST http://localhost:8081/transactions/credit \
  -H "Content-Type: application/json" \
  -d '{"client_id": "<uuid>", "amount": 250.00}'
```
