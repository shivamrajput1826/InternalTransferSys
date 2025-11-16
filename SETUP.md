# Database Setup Instructions

## Issue: Password Authentication Failed

Your local PostgreSQL is running but the password doesn't match. Here are solutions:

## Solution 1: Reset PostgreSQL Password (Recommended)

### On macOS (using Homebrew):
```bash
# Stop PostgreSQL
brew services stop postgresql@14  # or your version

# Start PostgreSQL in single-user mode
postgres --single -D /usr/local/var/postgres

# In the PostgreSQL prompt, run:
ALTER USER postgres WITH PASSWORD 'postgres';
\q

# Restart PostgreSQL
brew services start postgresql@14
```

### Alternative: Using psql (if you know current password):
```bash
psql -U postgres
ALTER USER postgres WITH PASSWORD 'postgres';
\q
```

## Solution 2: Create the Database

After fixing the password, create the database:
```bash
psql -U postgres -h localhost
CREATE DATABASE transfers;
\q
```

## Solution 3: Use Docker Instead

If you prefer to use Docker:

1. Stop local PostgreSQL (if possible)
2. Update docker-compose.yml to use port 5433:
   ```yaml
   ports:
     - "5433:5432"
   ```
3. Update config/dev.env:
   ```env
   DB_PORT=5433
   ```
4. Run: `docker compose up -d`

## Solution 4: Update Config with Correct Password

If you know your PostgreSQL password, update `config/dev.env`:
```env
DB_PASSWORD=your_actual_password
```

## Verify Connection

Test the connection:
```bash
psql -h localhost -U postgres -d transfers
```

If it works, you can run the application:
```bash
go run cmd/api/main.go
```

