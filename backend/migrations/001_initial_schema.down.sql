-- Drop triggers
DROP TRIGGER IF EXISTS calculate_trade_pnl ON trades;
DROP TRIGGER IF EXISTS update_journal_entries_updated_at ON journal_entries;
DROP TRIGGER IF EXISTS update_trades_updated_at ON trades;

-- Drop functions
DROP FUNCTION IF EXISTS calculate_pnl();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS trade_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS journal_entries;
DROP TABLE IF EXISTS trades;
DROP TABLE IF EXISTS magic_links;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "pgcrypto";
