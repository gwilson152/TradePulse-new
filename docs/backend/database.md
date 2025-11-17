# Database

## Connection

**Host**: postgres1.drivenw.local:5432
**Database**: tradepulse
**User**: tradepulse
**Password**: drv13llc!

## Schema

### users
```sql
id          UUID PRIMARY KEY
email       VARCHAR(255) UNIQUE NOT NULL
created_at  TIMESTAMP WITH TIME ZONE
last_login  TIMESTAMP WITH TIME ZONE
preferences JSONB
```

### magic_links
```sql
id          UUID PRIMARY KEY
user_id     UUID → users(id)
token       VARCHAR(255) UNIQUE NOT NULL
expires_at  TIMESTAMP WITH TIME ZONE NOT NULL
used_at     TIMESTAMP WITH TIME ZONE
created_at  TIMESTAMP WITH TIME ZONE
```

### trades
```sql
id          UUID PRIMARY KEY
user_id     UUID → users(id)
symbol      VARCHAR(20) NOT NULL
trade_type  VARCHAR(10) CHECK (LONG|SHORT)
quantity    DECIMAL(18,8) NOT NULL
entry_price DECIMAL(18,8) NOT NULL
exit_price  DECIMAL(18,8)
fees        DECIMAL(18,8) DEFAULT 0
pnl         DECIMAL(18,8)  -- Auto-calculated
opened_at   TIMESTAMP WITH TIME ZONE NOT NULL
closed_at   TIMESTAMP WITH TIME ZONE
created_at  TIMESTAMP WITH TIME ZONE
updated_at  TIMESTAMP WITH TIME ZONE
```

**Trigger**: `calculate_pnl()` - Auto-calculates P&L on insert/update

### journal_entries
```sql
id              UUID PRIMARY KEY
trade_id        UUID → trades(id)
user_id         UUID → users(id)
content         TEXT
emotional_state JSONB
created_at      TIMESTAMP WITH TIME ZONE
updated_at      TIMESTAMP WITH TIME ZONE
```

### attachments
```sql
id              UUID PRIMARY KEY
entry_id        UUID → journal_entries(id)
attachment_type VARCHAR(20) CHECK (screenshot|voice)
storage_path    VARCHAR(500) NOT NULL
filename        VARCHAR(255) NOT NULL
file_size       INTEGER
mime_type       VARCHAR(100)
uploaded_at     TIMESTAMP WITH TIME ZONE
```

### tags
```sql
id         UUID PRIMARY KEY
user_id    UUID → users(id)
name       VARCHAR(50) NOT NULL
color      VARCHAR(7)
created_at TIMESTAMP WITH TIME ZONE
UNIQUE(user_id, name)
```

### trade_tags
```sql
trade_id UUID → trades(id)
tag_id   UUID → tags(id)
PRIMARY KEY (trade_id, tag_id)
```

## Migrations

See [migrations.md](./migrations.md) for complete migration documentation.

**Quick Start:**
- Migrations run automatically on server start
- VSCode Tasks: `Ctrl+Shift+P` → "Migrate Up" / "Migration Status"
- Uses [golang-migrate](https://github.com/golang-migrate/migrate)

## Query Patterns

### Always Use Parameterized Queries

**Good**:
```go
row := db.QueryRow("SELECT id FROM users WHERE email = $1", email)
```

**Bad** (SQL injection risk):
```go
query := fmt.Sprintf("SELECT id FROM users WHERE email = '%s'", email)
row := db.QueryRow(query)
```

### Select Specific Columns

**Good**:
```go
db.Query("SELECT id, email, created_at FROM users WHERE ...")
```

**Bad**:
```go
db.Query("SELECT * FROM users WHERE ...")
```

### Use Transactions for Multi-Step Operations

```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback() // Rollback if not committed

// Execute multiple queries
_, err = tx.Exec("INSERT INTO trades ...")
_, err = tx.Exec("INSERT INTO journal_entries ...")

if err != nil {
    return err
}

return tx.Commit()
```

### Filter by User ID

**Always** filter by user ID to prevent data leaks:

```go
row := db.QueryRow(`
    SELECT id, symbol, pnl
    FROM trades
    WHERE id = $1 AND user_id = $2
`, tradeID, userID)
```

## Indexes

Key indexes for performance:

```sql
-- Users
CREATE INDEX idx_users_email ON users(email);

-- Trades
CREATE INDEX idx_trades_user_id ON trades(user_id);
CREATE INDEX idx_trades_symbol ON trades(symbol);
CREATE INDEX idx_trades_opened_at ON trades(opened_at);

-- Journal Entries
CREATE INDEX idx_journal_entries_trade_id ON journal_entries(trade_id);
CREATE INDEX idx_journal_entries_user_id ON journal_entries(user_id);

-- Magic Links
CREATE INDEX idx_magic_links_token ON magic_links(token);
CREATE INDEX idx_magic_links_user_id ON magic_links(user_id);
```

## Triggers

### P&L Auto-Calculation

Automatically calculates P&L when trade is closed:

```sql
CREATE OR REPLACE FUNCTION calculate_pnl()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.exit_price IS NOT NULL THEN
        IF NEW.trade_type = 'LONG' THEN
            NEW.pnl = (NEW.exit_price - NEW.entry_price) * NEW.quantity - NEW.fees;
        ELSE
            NEW.pnl = (NEW.entry_price - NEW.exit_price) * NEW.quantity - NEW.fees;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER calculate_trade_pnl
BEFORE INSERT OR UPDATE ON trades
FOR EACH ROW EXECUTE FUNCTION calculate_pnl();
```

### Updated At Timestamp

Automatically updates `updated_at` on row change:

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_trades_updated_at
BEFORE UPDATE ON trades
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## Backup & Restore

### Backup
```bash
pg_dump -h postgres1.drivenw.local -U tradepulse -d tradepulse > backup.sql
```

### Restore
```bash
psql -h postgres1.drivenw.local -U tradepulse -d tradepulse < backup.sql
```

## Best Practices

1. **Use UUIDs for IDs** - Prevents enumeration attacks
2. **Use TIMESTAMP WITH TIME ZONE** - Consistent timezone handling
3. **Add indexes for frequently queried columns** - Improves performance
4. **Use foreign key constraints** - Maintains data integrity
5. **Use CHECK constraints** - Validates data at database level
6. **Use triggers for computed fields** - Keeps logic in database
7. **Always filter by user_id** - Prevents data leaks
8. **Use transactions for related operations** - Maintains consistency
