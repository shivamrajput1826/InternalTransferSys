# Internal Transfer System

Internal transfer system for managing account transfers.

## Setup

1. Copy `.env.example` to `.env` and configure your environment variables
2. Run `docker-compose up -d` to start PostgreSQL
3. Run migrations
4. Start the application: `go run cmd/api/main.go`

## Project Structure

- `cmd/api/` - Application entry point
- `internal/` - Internal application code
  - `models/` - Data models
  - `handlers/` - HTTP handlers
  - `repository/` - Database operations
  - `service/` - Business logic
  - `database/` - Database connection setup
- `migrations/` - Database migrations
- `config/` - Configuration management
- `tests/` - Test files

## API Testing

### 1. Create Account 1
```bash
curl -X POST http://localhost:8080/api/accounts \
  -H "Content-Type: application/json" \
  -d '{"account_id": 123, "initial_balance": "1000.00"}'
```

### 2. Create Account 2
```bash
curl -X POST http://localhost:8080/api/accounts \
  -H "Content-Type: application/json" \
  -d '{"account_id": 456, "initial_balance": "500.00"}'
```

### 3. Check Account 1 Balance
```bash
curl -X GET http://localhost:8080/api/accounts/123
```

### 4. Transfer 100 from Account 1 to Account 2
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"source_account_id": 123, "destination_account_id": 456, "amount": "100.00"}'
```

### 5. Verify Balances Changed - Account 1
```bash
curl -X GET http://localhost:8080/api/accounts/123
```

### 6. Verify Balances Changed - Account 2
```bash
curl -X GET http://localhost:8080/api/accounts/456
```

**Expected Results:**
- Account 1 balance: `900.00` (1000 - 100)
- Account 2 balance: `600.00` (500 + 100)