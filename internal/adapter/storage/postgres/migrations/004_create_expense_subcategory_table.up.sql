CREATE TABLE IF NOT EXISTS expense_subcategories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    expense_category_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Foreign key constraint
    CONSTRAINT fk_expense_subcategories_category 
        FOREIGN KEY (expense_category_id) 
        REFERENCES expense_categories(id) 
        ON DELETE RESTRICT,
    
    -- Unique constraint for name within the same category
    CONSTRAINT uk_expense_subcategories_name_category 
        UNIQUE (name, expense_category_id)
);

-- Create index on expense_category_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_expense_subcategories_category_id ON expense_subcategories(expense_category_id);

-- Create index on name for faster lookups
CREATE INDEX IF NOT EXISTS idx_expense_subcategories_name ON expense_subcategories(name);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_expense_subcategories_created_at ON expense_subcategories(created_at); 