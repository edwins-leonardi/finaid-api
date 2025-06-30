-- Drop indexes first
DROP INDEX IF EXISTS idx_expense_categories_created_at;
DROP INDEX IF EXISTS idx_expense_categories_name;

-- Drop the table
DROP TABLE IF EXISTS expense_categories; 