-- Create accounts table migration
CREATE TABLE IF NOT EXISTS accounts (
    account_id BIGINT PRIMARY KEY,
    balance NUMERIC(20,5) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on account_id (already primary key, but explicit for clarity)
CREATE INDEX IF NOT EXISTS idx_accounts_account_id ON accounts(account_id);

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
