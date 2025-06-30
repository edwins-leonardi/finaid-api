-- Drop indexes first
DROP INDEX IF EXISTS idx_expense_subcategories_created_at;
DROP INDEX IF EXISTS idx_expense_subcategories_name;
DROP INDEX IF EXISTS idx_expense_subcategories_category_id;

-- Drop the table (foreign key constraints will be dropped automatically)
DROP TABLE IF EXISTS expense_subcategories; 