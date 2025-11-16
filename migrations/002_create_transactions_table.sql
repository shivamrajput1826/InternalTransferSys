-- Create transactions table migration
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    source_account_id BIGINT NOT NULL,
    destination_account_id BIGINT NOT NULL,
    amount NUMERIC(20,5) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_transactions_source_account_id ON transactions(source_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_destination_account_id ON transactions(destination_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);

-- Add foreign key constraints (optional, but recommended for data integrity)
-- Uncomment these if you want referential integrity
-- ALTER TABLE transactions 
--     ADD CONSTRAINT fk_transactions_source_account 
--     FOREIGN KEY (source_account_id) REFERENCES accounts(account_id);
-- 
-- ALTER TABLE transactions 
--     ADD CONSTRAINT fk_transactions_destination_account 
--     FOREIGN KEY (destination_account_id) REFERENCES accounts(account_id);
