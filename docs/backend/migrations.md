# Migrations

Uses [golang-migrate](https://github.com/golang-migrate/migrate) - the standard migration tool for Go.

## Quick Start

**Reset & Run Migrations:**
1. `Ctrl+Shift+P` → "Reset Database"
2. `Ctrl+Shift+P` → "Start Backend" (auto-runs migrations)

## VSCode Tasks

- **Migration Status** - Show current version
- **Migrate Up** - Apply pending migrations
- **Migrate Down** - Rollback last migration
- **Migrate Force Version** - Fix dirty state
- **Reset Database** - Drop all tables

## CLI Commands

```bash
cd backend

# Status
migrate -path migrations -database "postgres://tradepulse:drv13llc!@postgres1.drivenw.local:5432/tradepulse?sslmode=disable" version

# Apply all
migrate -path migrations -database "postgres://..." up

# Rollback one
migrate -path migrations -database "postgres://..." down 1

# Force version (if dirty)
migrate -path migrations -database "postgres://..." force 1
```

## Create New Migration

```bash
cd backend
migrate create -ext sql -dir migrations -seq your_feature_name
```

Creates `NNN_your_feature_name.up.sql` and `NNN_your_feature_name.down.sql`

## Current Migrations

- **001_initial_schema** - All tables, indexes, functions, triggers
- **002_add_password_auth** - Password authentication

## How It Works

- Migrations run automatically on server startup
- `schema_migrations` table tracks applied migrations
- Safe to run `up` multiple times (skips applied migrations)
- Use `IF NOT EXISTS` / `IF EXISTS` for idempotency

## Troubleshooting

**Dirty Migration Error:**
1. Fix database manually
2. `Ctrl+Shift+P` → "Migrate Force Version"

## Migration System Upgrade

Previously used custom hardcoded migrations. Now uses industry-standard golang-migrate.

**What Changed:**
- ✅ SQL files are source of truth (not Go code)
- ✅ CLI tools for management
- ✅ Auto-runs on startup
- ✅ Up/down migration support
- ✅ Standard library approach

**Files Modified:**
- `backend/internal/database/db.go` - Rewritten (263→81 lines)
- `backend/migrations/001_initial_schema.up.sql` - Complete schema
- `.vscode/tasks.json` - Added migration tasks

**Fixed Issues:**
- Missing `journal_entries` table
- Missing `attachments`, `tags`, `trade_tags` tables
- Database functions and triggers
