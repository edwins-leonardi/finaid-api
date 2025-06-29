-- Drop indexes first
DROP INDEX IF EXISTS idx_account_second_owner;
DROP INDEX IF EXISTS idx_account_primary_owner;
DROP INDEX IF EXISTS idx_account_currency;
DROP INDEX IF EXISTS idx_account_type;
DROP INDEX IF EXISTS idx_account_name;

-- Drop the table
DROP TABLE IF EXISTS account; 