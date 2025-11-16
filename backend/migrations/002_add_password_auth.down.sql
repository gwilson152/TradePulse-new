-- Rollback password authentication
DROP INDEX IF EXISTS idx_users_email_password;
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;
