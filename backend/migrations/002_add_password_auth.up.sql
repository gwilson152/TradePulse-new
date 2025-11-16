-- Add password authentication to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT;

-- Add index for faster password lookups
CREATE INDEX IF NOT EXISTS idx_users_email_password ON users(email) WHERE password_hash IS NOT NULL;
