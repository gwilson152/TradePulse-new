package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// New creates a new database connection
func New(cfg Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// RunMigrations runs all pending database migrations
func (db *DB) RunMigrations() error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Check if migration 001 has been applied
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = 1").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migrations: %w", err)
	}

	// If migration hasn't been applied, run it
	if count == 0 {
		// Read and execute migration file
		// For simplicity, we'll execute the SQL directly here
		// In production, you'd read from the migrations directory
		err = db.executeMigration001()
		if err != nil {
			return fmt.Errorf("failed to run migration 001: %w", err)
		}

		// Mark migration as applied
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (1)")
		if err != nil {
			return fmt.Errorf("failed to mark migration as applied: %w", err)
		}
	}

	// Check if migration 002 has been applied
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = 2").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migrations: %w", err)
	}

	// If migration hasn't been applied, run it
	if count == 0 {
		err = db.executeMigration002()
		if err != nil {
			return fmt.Errorf("failed to run migration 002: %w", err)
		}

		// Mark migration as applied
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (2)")
		if err != nil {
			return fmt.Errorf("failed to mark migration as applied: %w", err)
		}
	}

	return nil
}

func (db *DB) executeMigration001() error {
	migrationSQL := `
		-- Enable UUID extension
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";

		-- Users table
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			last_login TIMESTAMP WITH TIME ZONE,
			preferences JSONB DEFAULT '{}'::jsonb
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

		-- Magic links table
		CREATE TABLE IF NOT EXISTS magic_links (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(255) UNIQUE NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			used_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_magic_links_token ON magic_links(token);
		CREATE INDEX IF NOT EXISTS idx_magic_links_user_id ON magic_links(user_id);
		CREATE INDEX IF NOT EXISTS idx_magic_links_expires_at ON magic_links(expires_at);

		-- Trades table
		CREATE TABLE IF NOT EXISTS trades (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			symbol VARCHAR(20) NOT NULL,
			trade_type VARCHAR(10) NOT NULL CHECK (trade_type IN ('LONG', 'SHORT')),
			quantity DECIMAL(18, 8) NOT NULL,
			entry_price DECIMAL(18, 8) NOT NULL,
			exit_price DECIMAL(18, 8),
			fees DECIMAL(18, 8) DEFAULT 0,
			pnl DECIMAL(18, 8),
			opened_at TIMESTAMP WITH TIME ZONE NOT NULL,
			closed_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_trades_user_id ON trades(user_id);
		CREATE INDEX IF NOT EXISTS idx_trades_symbol ON trades(symbol);
		CREATE INDEX IF NOT EXISTS idx_trades_opened_at ON trades(opened_at);
		CREATE INDEX IF NOT EXISTS idx_trades_closed_at ON trades(closed_at);
	`

	_, err := db.Exec(migrationSQL)
	return err
}

func (db *DB) executeMigration002() error {
	migrationSQL := `
		-- Add password authentication to users table
		ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT;

		-- Add index for faster password lookups
		CREATE INDEX IF NOT EXISTS idx_users_email_password ON users(email) WHERE password_hash IS NOT NULL;
	`

	_, err := db.Exec(migrationSQL)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
