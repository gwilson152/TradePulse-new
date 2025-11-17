package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/models"
)

// CreateRuleSet creates a new rule set
func (db *DB) CreateRuleSet(ctx context.Context, ruleSet *models.RuleSet) error {
	query := `
		INSERT INTO rule_sets (id, user_id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	ruleSet.ID = uuid.New()

	err := db.QueryRowContext(
		ctx,
		query,
		ruleSet.ID,
		ruleSet.UserID,
		ruleSet.Name,
		ruleSet.Description,
		ruleSet.IsActive,
	).Scan(&ruleSet.ID, &ruleSet.CreatedAt, &ruleSet.UpdatedAt)

	return err
}

// GetRuleSet retrieves a rule set by ID
func (db *DB) GetRuleSet(ctx context.Context, id, userID uuid.UUID) (*models.RuleSet, error) {
	query := `
		SELECT id, user_id, name, description, is_active, created_at, updated_at
		FROM rule_sets
		WHERE id = $1 AND user_id = $2
	`

	ruleSet := &models.RuleSet{}
	err := db.QueryRowContext(ctx, query, id, userID).Scan(
		&ruleSet.ID,
		&ruleSet.UserID,
		&ruleSet.Name,
		&ruleSet.Description,
		&ruleSet.IsActive,
		&ruleSet.CreatedAt,
		&ruleSet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("rule set not found")
		}
		return nil, err
	}

	// Load rules for this rule set
	rules, err := db.GetRulesByRuleSetID(ctx, ruleSet.ID)
	if err != nil {
		return nil, err
	}
	ruleSet.Rules = rules

	return ruleSet, nil
}

// ListRuleSets retrieves all rule sets for a user
func (db *DB) ListRuleSets(ctx context.Context, userID uuid.UUID) ([]models.RuleSet, error) {
	query := `
		SELECT id, user_id, name, description, is_active, created_at, updated_at
		FROM rule_sets
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize to empty slice to avoid null JSON serialization
	ruleSets := make([]models.RuleSet, 0)
	for rows.Next() {
		var ruleSet models.RuleSet
		err := rows.Scan(
			&ruleSet.ID,
			&ruleSet.UserID,
			&ruleSet.Name,
			&ruleSet.Description,
			&ruleSet.IsActive,
			&ruleSet.CreatedAt,
			&ruleSet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Load rules for each rule set
		rules, err := db.GetRulesByRuleSetID(ctx, ruleSet.ID)
		if err != nil {
			return nil, err
		}
		ruleSet.Rules = rules

		ruleSets = append(ruleSets, ruleSet)
	}

	return ruleSets, nil
}

// UpdateRuleSet updates an existing rule set
func (db *DB) UpdateRuleSet(ctx context.Context, ruleSet *models.RuleSet) error {
	query := `
		UPDATE rule_sets
		SET name = $1, description = $2, is_active = $3
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at
	`

	err := db.QueryRowContext(
		ctx,
		query,
		ruleSet.Name,
		ruleSet.Description,
		ruleSet.IsActive,
		ruleSet.ID,
		ruleSet.UserID,
	).Scan(&ruleSet.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("rule set not found")
		}
		return err
	}

	return nil
}

// DeleteRuleSet deletes a rule set
func (db *DB) DeleteRuleSet(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM rule_sets WHERE id = $1 AND user_id = $2`

	result, err := db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("rule set not found")
	}

	return nil
}

// CreateRule creates a new rule in a rule set
func (db *DB) CreateRule(ctx context.Context, rule *models.Rule) error {
	query := `
		INSERT INTO rules (id, rule_set_id, title, description, weight, phase, category, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, created_at
	`

	rule.ID = uuid.New()

	err := db.QueryRowContext(
		ctx,
		query,
		rule.ID,
		rule.RuleSetID,
		rule.Title,
		rule.Description,
		rule.Weight,
		rule.Phase,
		rule.Category,
	).Scan(&rule.ID, &rule.CreatedAt)

	return err
}

// GetRule retrieves a rule by ID
func (db *DB) GetRule(ctx context.Context, id uuid.UUID) (*models.Rule, error) {
	query := `
		SELECT id, rule_set_id, title, description, weight, phase, category, created_at
		FROM rules
		WHERE id = $1
	`

	rule := &models.Rule{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.RuleSetID,
		&rule.Title,
		&rule.Description,
		&rule.Weight,
		&rule.Phase,
		&rule.Category,
		&rule.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("rule not found")
		}
		return nil, err
	}

	return rule, nil
}

// GetRulesByRuleSetID retrieves all rules for a rule set
func (db *DB) GetRulesByRuleSetID(ctx context.Context, ruleSetID uuid.UUID) ([]models.Rule, error) {
	query := `
		SELECT id, rule_set_id, title, description, weight, phase, category, created_at
		FROM rules
		WHERE rule_set_id = $1
		ORDER BY phase, category, created_at
	`

	rows, err := db.QueryContext(ctx, query, ruleSetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize to empty slice to avoid null JSON serialization
	rules := make([]models.Rule, 0)
	for rows.Next() {
		var rule models.Rule
		err := rows.Scan(
			&rule.ID,
			&rule.RuleSetID,
			&rule.Title,
			&rule.Description,
			&rule.Weight,
			&rule.Phase,
			&rule.Category,
			&rule.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// UpdateRule updates an existing rule
func (db *DB) UpdateRule(ctx context.Context, rule *models.Rule) error {
	query := `
		UPDATE rules
		SET title = $1, description = $2, weight = $3, phase = $4, category = $5
		WHERE id = $6
	`

	result, err := db.ExecContext(
		ctx,
		query,
		rule.Title,
		rule.Description,
		rule.Weight,
		rule.Phase,
		rule.Category,
		rule.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("rule not found")
	}

	return nil
}

// DeleteRule deletes a rule
func (db *DB) DeleteRule(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM rules WHERE id = $1`

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("rule not found")
	}

	return nil
}
