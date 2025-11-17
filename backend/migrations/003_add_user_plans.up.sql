-- Add plan fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS plan_type VARCHAR(50) DEFAULT 'starter';
ALTER TABLE users ADD COLUMN IF NOT EXISTS plan_status VARCHAR(50) DEFAULT 'beta_free';
ALTER TABLE users ADD COLUMN IF NOT EXISTS plan_selected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Create index for faster plan-based queries
CREATE INDEX IF NOT EXISTS idx_users_plan_type ON users(plan_type);

-- Add check constraint for valid plan types
ALTER TABLE users ADD CONSTRAINT check_plan_type
    CHECK (plan_type IN ('starter', 'pro', 'premium'));

-- Add check constraint for valid plan statuses
ALTER TABLE users ADD CONSTRAINT check_plan_status
    CHECK (plan_status IN ('beta_free', 'active', 'cancelled', 'trial', 'expired'));
