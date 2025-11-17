-- Reset Database Script
-- This will drop all tables and recreate them fresh

-- Drop all tables in the correct order (reverse of dependencies)
DROP TABLE IF EXISTS trade_tags CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS attachments CASCADE;
DROP TABLE IF EXISTS journal_entries CASCADE;
DROP TABLE IF EXISTS trades CASCADE;
DROP TABLE IF EXISTS magic_links CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop golang-migrate schema_migrations table
DROP TABLE IF EXISTS schema_migrations CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS calculate_pnl() CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Drop extension (optional - you may want to keep it)
-- DROP EXTENSION IF EXISTS "pgcrypto";

-- The migrations will be re-run automatically by the application on next start
