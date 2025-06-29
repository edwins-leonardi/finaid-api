CREATE TABLE IF NOT EXISTS account (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    initial_balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    primary_owner_id BIGINT NOT NULL,
    second_owner_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Foreign key constraints
    CONSTRAINT fk_account_primary_owner FOREIGN KEY (primary_owner_id) REFERENCES person(id) ON DELETE RESTRICT,
    CONSTRAINT fk_account_second_owner FOREIGN KEY (second_owner_id) REFERENCES person(id) ON DELETE SET NULL
);

-- Create indexes for better performance
CREATE INDEX idx_account_name ON account(name);
CREATE INDEX idx_account_type ON account(account_type);
CREATE INDEX idx_account_currency ON account(currency);
CREATE INDEX idx_account_primary_owner ON account(primary_owner_id);
CREATE INDEX idx_account_second_owner ON account(second_owner_id); 