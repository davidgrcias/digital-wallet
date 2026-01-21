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

| Method | Endpoint | Description |
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

## Technical Notes

- **Concurrency**: I used `SELECT ... FOR UPDATE` (Pessimistic Locking) because in a wallet system, precision beats speed. It prevents race conditions during simultaneous withdrawals.
- **Idempotency**: Requests include an `Idempotency-Key` header. This ensures that if a user retries a request due to a network timeout, they won't be charged twice.
- **Data Types**: I used `float64` for current balance calculations to keep the implementation simple. For a real production system, I'd switch to `int64` (storing cents) or a proper Decimal library to avoid any floating-point precision issues. The database already uses `DECIMAL(18,2)` for safety.

## Author

**David Garcia Saragih**
- GitHub: [@davidgrcias](https://github.com/davidgrcias)

## License

MIT
