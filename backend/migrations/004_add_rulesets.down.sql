-- Drop triggers
DROP TRIGGER IF EXISTS update_rule_sets_updated_at ON rule_sets;

-- Drop indexes
DROP INDEX IF EXISTS idx_rules_rule_set_id;
DROP INDEX IF EXISTS idx_rule_sets_is_active;
DROP INDEX IF EXISTS idx_rule_sets_user_id;

-- Drop tables
DROP TABLE IF EXISTS rules;
DROP TABLE IF EXISTS rule_sets;
