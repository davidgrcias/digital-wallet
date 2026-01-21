# Digital Wallet API

[![Go CI](https://github.com/davidgrcias/digital-wallet/actions/workflows/ci.yml/badge.svg)](https://github.com/davidgrcias/digital-wallet/actions/workflows/ci.yml)

RESTful API backend for digital wallet services, built with Go and PostgreSQL.

## Requirements

- Go 1.21+
- Docker & Docker Compose
- Postman (for testing)
- Make (Optional, for shortcuts)

## Quick Start

### Using Makefile (Recommended)
```bash
make run    # Start Database & Server
make test   # Run Unit Tests
make stop   # Stop application
```

### Manual Setup
1. **Start Database**
   ```bash
   docker compose up -d
   ```
2. **Run Application**
   ```bash
   go mod tidy
   go run cmd/api/main.go
   ```
   Server will start at `http://localhost:8081`.

## How to Test

I've included a Postman Collection to make testing easier without using CLI commands.

1. Open **Postman**.
2. Click **Import** -> Upload file `digital-wallet.postman_collection.json` from this repository.
3. You will see 4 ready-to-use requests:
   - **Health Check**: Verify server is running.
   - **Get Balance**: Check balance for user John Doe.
   - **Withdraw 50k**: Simulate a successful withdrawal.
   - **Withdraw Fail**: Simulate insufficient funds error.

## Configuration

The application is configured via environment variables. Copy `.env.example` to `.env`.

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Port for the API server | `8081` |
| `DB_HOST` | Database hostname | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `wallet_db` |

## API Endpoints

### cURL Examples

If you prefer terminal commands over Postman:

**1. Check Health**
```bash
curl -X GET http://localhost:8081/health
```

**2. Get Balance**
```bash
curl -X GET http://localhost:8081/api/v1/users/550e8400-e29b-41d4-a716-446655440001/balance
```

**3. Withdraw Funds**
```bash
curl -X POST http://localhost:8081/api/v1/users/550e8400-e29b-41d4-a716-446655440001/withdraw \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: unique-key-123" \
  -d '{"amount": 50000, "description": "Coffee money"}'
```

### Route List


|:---|:---|:---|
| GET | `/health` | Server status check |
| GET | `/api/v1/users/{id}/balance` | Get wallet balance |
| POST | `/api/v1/users/{id}/withdraw` | Withdraw funds |

### Security Features
- **Idempotency Key:** Support `Idempotency-Key` header on Withdraw endpoint to prevent double-spending on retries.
- **Race Condition Prevention:** Uses `SELECT ... FOR UPDATE` (Pessimistic Locking).

### Test Data (Pre-seeded Users)

| Name | ID | Initial Balance |
|:---|:---|:---|
| John Doe | `550e8400-e29b-41d4-a716-446655440001` | 1,000,000 |
| Jane Smith | `550e8400-e29b-41d4-a716-446655440002` | 500,000 |
| Bob Wilson | `550e8400-e29b-41d4-a716-446655440003` | 250,000 |

## Architecture

- **Language:** Go (Golang)
- **Database:** PostgreSQL 16 (via Docker)
- **Pattern:** Clean Architecture (Handler -> Usecase -> Repository)
- **Transactions:** All withdrawals are atomic (ACID compliant).

## Engineering Decisions & Trade-offs

Here is a breakdown of the technical decisions made for this project:

### 1. Concurrency Safety
I chose **Pessimistic Locking** (`SELECT ... FOR UPDATE`) over Optimistic Locking.
- **Why?** In a financial system, data correctness is paramount. We want to prevent race conditions at the database level to ensure a user's balance never drops below zero, even if thousands of requests hit the API simultaneously.
- **Trade-off:** Slightly lower throughput compared to Optimistic Locking, but significantly safer for this use case.

### 2. Idempotency Key
Failed network requests shouldn't drain a wallet twice. I implemented a custom middleware compliant with the `Idempotency-Key` header standard.
- **Behavior:** If a client retries a request with the same key, they receive the *cached* response (from the database) instead of triggering a new transaction.

### 3. Floating Point Math
**Note for Reviewers:** I used `float64` for simplicity given the time constraints.
- **Production View:** In a real-world banking ledger, I would strictly use a decimal library (like `shopspring/decimal`) or store values as `int64` (minor units) to avoid floating-point precision errors.
- **Current Mitigation:** Database column uses `DECIMAL(18,2)` to ensure storage precision.

## Author

**David Garcia Saragih**
- GitHub: [@davidgrcias](https://github.com/davidgrcias)

## License

MIT
