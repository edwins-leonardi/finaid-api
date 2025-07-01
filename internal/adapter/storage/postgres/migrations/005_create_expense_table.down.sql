-- Drop indexes first
DROP INDEX IF EXISTS idx_expenses_payee_date;
DROP INDEX IF EXISTS idx_expenses_category_date;
DROP INDEX IF EXISTS idx_expenses_created_at;
DROP INDEX IF EXISTS idx_expenses_date;
DROP INDEX IF EXISTS idx_expenses_payee_id;
DROP INDEX IF EXISTS idx_expenses_subcategory_id;
DROP INDEX IF EXISTS idx_expenses_category_id;

-- Drop the table
DROP TABLE IF EXISTS expenses; 