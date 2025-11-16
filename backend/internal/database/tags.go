package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/models"
)

// ListTags retrieves all tags for a user
func (db *DB) ListTags(ctx context.Context, userID uuid.UUID) ([]models.Tag, error) {
	query := `
		SELECT
			t.id, t.user_id, t.name, t.color, t.created_at,
			COALESCE(
				(SELECT COUNT(*) FROM trade_tags tt WHERE tt.tag_id = t.id),
				0
			) as usage_count
		FROM tags t
		WHERE t.user_id = $1
		ORDER BY t.name ASC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(
			&tag.ID, &tag.UserID, &tag.Name, &tag.Color, &tag.CreatedAt, &tag.UsageCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tags: %w", err)
	}

	return tags, nil
}

// GetTag retrieves a single tag by ID
func (db *DB) GetTag(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Tag, error) {
	query := `
		SELECT
			t.id, t.user_id, t.name, t.color, t.created_at,
			COALESCE(
				(SELECT COUNT(*) FROM trade_tags tt WHERE tt.tag_id = t.id),
				0
			) as usage_count
		FROM tags t
		WHERE t.id = $1 AND t.user_id = $2`

	var tag models.Tag
	err := db.QueryRow(query, id, userID).Scan(
		&tag.ID, &tag.UserID, &tag.Name, &tag.Color, &tag.CreatedAt, &tag.UsageCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return &tag, nil
}

// GetTagByName retrieves a tag by name for a user
func (db *DB) GetTagByName(ctx context.Context, name string, userID uuid.UUID) (*models.Tag, error) {
	query := `
		SELECT
			t.id, t.user_id, t.name, t.color, t.created_at,
			COALESCE(
				(SELECT COUNT(*) FROM trade_tags tt WHERE tt.tag_id = t.id),
				0
			) as usage_count
		FROM tags t
		WHERE LOWER(t.name) = LOWER($1) AND t.user_id = $2`

	var tag models.Tag
	err := db.QueryRow(query, name, userID).Scan(
		&tag.ID, &tag.UserID, &tag.Name, &tag.Color, &tag.CreatedAt, &tag.UsageCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by name: %w", err)
	}

	return &tag, nil
}

// CreateTag inserts a new tag
func (db *DB) CreateTag(ctx context.Context, tag *models.Tag) error {
	// Check if tag with same name already exists for this user
	existing, err := db.GetTagByName(ctx, tag.Name, tag.UserID)
	if err != nil {
		return fmt.Errorf("failed to check for existing tag: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("tag with name '%s' already exists", tag.Name)
	}

	query := `
		INSERT INTO tags (user_id, name, color)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err = db.QueryRow(
		query,
		tag.UserID, tag.Name, tag.Color,
	).Scan(&tag.ID, &tag.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	tag.UsageCount = 0

	return nil
}

// UpdateTag updates an existing tag
func (db *DB) UpdateTag(ctx context.Context, id uuid.UUID, userID uuid.UUID, tag *models.Tag) error {
	// Check if another tag with the same name exists
	existing, err := db.GetTagByName(ctx, tag.Name, userID)
	if err != nil {
		return fmt.Errorf("failed to check for existing tag: %w", err)
	}
	if existing != nil && existing.ID != id {
		return fmt.Errorf("tag with name '%s' already exists", tag.Name)
	}

	query := `
		UPDATE tags
		SET name = $3, color = $4
		WHERE id = $1 AND user_id = $2
		RETURNING created_at`

	err = db.QueryRow(query, id, userID, tag.Name, tag.Color).Scan(&tag.CreatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("tag not found or unauthorized")
	}
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	tag.ID = id
	tag.UserID = userID

	// Get usage count
	var usageCount int
	err = db.QueryRow("SELECT COUNT(*) FROM trade_tags WHERE tag_id = $1", id).Scan(&usageCount)
	if err == nil {
		tag.UsageCount = usageCount
	}

	return nil
}

// DeleteTag deletes a tag
func (db *DB) DeleteTag(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM tags WHERE id = $1 AND user_id = $2`

	result, err := db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tag not found or unauthorized")
	}

	return nil
}
