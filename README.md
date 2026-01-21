# Digital Wallet API

REST API for digital wallet operations using Go and PostgreSQL.

## Features

- Balance inquiry
- Withdraw with transaction logging
- Pessimistic locking to prevent race conditions

## Tech Stack

- Go 1.21+
- PostgreSQL 16
- gorilla/mux
- Docker

## Project Structure

```
.
├── cmd/api/           # Entry point
├── internal/
│   ├── config/        # Configuration
│   ├── domain/        # Entities
│   ├── handler/       # HTTP handlers
│   ├── middleware/    # Middleware
│   ├── repository/    # Data access
│   └── usecase/       # Business logic
├── pkg/
│   ├── database/      # DB connection
│   └── response/      # Response helpers
└── migrations/        # SQL migrations
```

## Quick Start

```bash
# Start PostgreSQL
docker compose up -d

# Install dependencies
go mod tidy

# Run
go run cmd/api/main.go
```

Server runs at `http://localhost:8080`

## API

### Get Balance
```
GET /api/v1/users/{user_id}/balance
```

### Withdraw
```
POST /api/v1/users/{user_id}/withdraw
Content-Type: application/json

{"amount": 50000, "description": "optional"}
```

## Test Users

```
550e8400-e29b-41d4-a716-446655440001  John Doe     Rp 1,000,000
550e8400-e29b-41d4-a716-446655440002  Jane Smith   Rp 500,000
550e8400-e29b-41d4-a716-446655440003  Bob Wilson   Rp 250,000
```

## Configuration

Copy `.env.example` to `.env` and adjust as needed.

## License

MIT
