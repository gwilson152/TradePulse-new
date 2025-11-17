-- Create rule_sets table
CREATE TABLE IF NOT EXISTS rule_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create rules table
CREATE TABLE IF NOT EXISTS rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_set_id UUID NOT NULL REFERENCES rule_sets(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    weight INTEGER NOT NULL DEFAULT 3 CHECK (weight BETWEEN 1 AND 5),
    phase VARCHAR(50) NOT NULL CHECK (phase IN ('PRE_TRADE', 'DURING_TRADE', 'POST_TRADE')),
    category VARCHAR(50) NOT NULL CHECK (category IN ('RISK_MANAGEMENT', 'ENTRY', 'EXIT', 'POSITION_SIZING', 'TIMING', 'PSYCHOLOGY', 'GENERAL')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_rule_sets_user_id ON rule_sets(user_id);
CREATE INDEX IF NOT EXISTS idx_rule_sets_is_active ON rule_sets(is_active);
CREATE INDEX IF NOT EXISTS idx_rules_rule_set_id ON rules(rule_set_id);

-- Create updated_at trigger for rule_sets
DROP TRIGGER IF EXISTS update_rule_sets_updated_at ON rule_sets;
CREATE TRIGGER update_rule_sets_updated_at BEFORE UPDATE ON rule_sets
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
