#!/bin/bash

# Test database connection script

echo "Testing PostgreSQL connection..."

# Try to connect with password from config
export PGPASSWORD=postgres
psql -h localhost -U postgres -d postgres -c "SELECT version();" 2>&1

if [ $? -eq 0 ]; then
    echo "✓ Connection successful!"
    echo "Creating database if it doesn't exist..."
    psql -h localhost -U postgres -d postgres -c "CREATE DATABASE transfers;" 2>&1 | grep -v "already exists"
    echo "✓ Database ready!"
else
    echo "✗ Connection failed. Please check:"
    echo "  1. PostgreSQL is running"
    echo "  2. Password is correct (current: postgres)"
    echo "  3. User 'postgres' exists"
    echo ""
    echo "To reset password, run:"
    echo "  psql -U postgres"
    echo "  ALTER USER postgres WITH PASSWORD 'postgres';"
fi

