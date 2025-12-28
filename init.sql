-- PostgreSQL initialization script for gonewproject
-- This file will be executed when the PostgreSQL container starts for the first time
-- Note: GORM will handle table migrations automatically, so manual table creation is not needed

-- Create additional schemas if needed
-- CREATE SCHEMA IF NOT EXISTS app_schema;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- GORM will create the users table with the following structure:
-- - id (uuid, primary key, default: gen_random_uuid())
-- - name (varchar, not null)
-- - email (varchar, unique, not null)
-- - created_at (timestamp)
-- - updated_at (timestamp)
-- - deleted_at (timestamp, soft delete)

-- Grant permissions to the application user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO gonewproject;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO gonewproject;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO PUBLIC;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO PUBLIC;
