CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    category_id INTEGER NOT NULL,
    subcategory_id INTEGER,
    date DATE NOT NULL,
    payee_id INTEGER NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_expenses_category 
        FOREIGN KEY (category_id) 
        REFERENCES expense_categories(id) 
        ON DELETE RESTRICT,
    
    CONSTRAINT fk_expenses_subcategory 
        FOREIGN KEY (subcategory_id) 
        REFERENCES expense_subcategories(id) 
        ON DELETE SET NULL,
    
    CONSTRAINT fk_expenses_payee 
        FOREIGN KEY (payee_id) 
        REFERENCES person(id) 
        ON DELETE RESTRICT
);

-- Create indexes for better query performance
CREATE INDEX idx_expenses_category_id ON expenses(category_id);
CREATE INDEX idx_expenses_subcategory_id ON expenses(subcategory_id);
CREATE INDEX idx_expenses_payee_id ON expenses(payee_id);
CREATE INDEX idx_expenses_date ON expenses(date);
CREATE INDEX idx_expenses_created_at ON expenses(created_at);

-- Create composite indexes for common query patterns
CREATE INDEX idx_expenses_category_date ON expenses(category_id, date);
CREATE INDEX idx_expenses_payee_date ON expenses(payee_id, date); 