-- Initialize finaid database
-- This file runs automatically when the PostgreSQL container starts for the first time

-- Create a sample table (you can modify or remove this)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data (optional)
INSERT INTO users (email, name) VALUES 
    ('admin@finaid.com', 'Admin User'),
    ('user@finaid.com', 'Test User')
ON CONFLICT (email) DO NOTHING; 