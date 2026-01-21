# Digital Wallet API

[![Go CI](https://github.com/davidgarcia/digital-wallet/actions/workflows/ci.yml/badge.svg)](https://github.com/davidgarcia/digital-wallet/actions/workflows/ci.yml)

RESTful API backend for digital wallet services, built with Go and PostgreSQL.

## Requirements

- Go 1.25+
- Docker & Docker Compose
- Postman (for testing)
- Make (Optional, for shortcuts)

## Quick Start

### Using Makefile (Recommended for Windows/Linux)
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
   - **Get Balance**: Check balance for user John Doe (ID: ...0001).
   - **Withdraw 50k**: Simulate a successful withdrawal.
   - **Withdraw Fail**: Simulate insufficient funds error.

## API Endpoints

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

## Architecture Notes

- **Language:** Go (Golang)
- **Database:** PostgreSQL 16 (via Docker)
- **Architecture:** Clean Architecture (Handler -> Usecase -> Repository)
- **Transactions:** All withdrawals are wrapped in database transactions (ACID).

## License

MIT
