-- Migration: Drop person table (DOWN)
-- Created: 2024-01-01

-- Drop index
DROP INDEX IF EXISTS idx_person_email;

-- Drop table
DROP TABLE IF EXISTS person; 