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

