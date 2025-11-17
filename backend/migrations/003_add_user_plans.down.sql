-- Remove constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_plan_type;
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_plan_status;

-- Remove index
DROP INDEX IF EXISTS idx_users_plan_type;

-- Remove plan columns
ALTER TABLE users DROP COLUMN IF EXISTS plan_type;
ALTER TABLE users DROP COLUMN IF EXISTS plan_status;
ALTER TABLE users DROP COLUMN IF EXISTS plan_selected_at;
