package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/models"
)

// CreateJournalEntry creates a new journal entry
func (db *DB) CreateJournalEntry(ctx context.Context, entry *models.JournalEntry) error {
	query := `
		INSERT INTO journal_entries (id, trade_id, user_id, content, emotional_state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	entry.ID = uuid.New()

	err := db.QueryRowContext(
		ctx,
		query,
		entry.ID,
		entry.TradeID,
		entry.UserID,
		entry.Content,
		entry.EmotionalState,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)

	return err
}

// GetJournalEntry retrieves a journal entry by ID
func (db *DB) GetJournalEntry(ctx context.Context, id, userID uuid.UUID) (*models.JournalEntry, error) {
	query := `
		SELECT id, trade_id, user_id, content, emotional_state, created_at, updated_at
		FROM journal_entries
		WHERE id = $1 AND user_id = $2
	`

	entry := &models.JournalEntry{}
	err := db.QueryRowContext(ctx, query, id, userID).Scan(
		&entry.ID,
		&entry.TradeID,
		&entry.UserID,
		&entry.Content,
		&entry.EmotionalState,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, err
	}

	// Load attachments
	attachments, err := db.GetAttachmentsByEntryID(ctx, entry.ID)
	if err != nil {
		return nil, err
	}
	entry.Attachments = attachments

	return entry, nil
}

// ListJournalEntries retrieves all journal entries for a user with pagination
func (db *DB) ListJournalEntries(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.JournalEntry, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM journal_entries WHERE user_id = $1`
	err := db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get entries
	query := `
		SELECT id, trade_id, user_id, content, emotional_state, created_at, updated_at
		FROM journal_entries
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []models.JournalEntry
	for rows.Next() {
		var entry models.JournalEntry
		err := rows.Scan(
			&entry.ID,
			&entry.TradeID,
			&entry.UserID,
			&entry.Content,
			&entry.EmotionalState,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Load attachments for each entry
		attachments, err := db.GetAttachmentsByEntryID(ctx, entry.ID)
		if err != nil {
			return nil, 0, err
		}
		entry.Attachments = attachments

		entries = append(entries, entry)
	}

	return entries, total, nil
}

// UpdateJournalEntry updates an existing journal entry
func (db *DB) UpdateJournalEntry(ctx context.Context, entry *models.JournalEntry) error {
	query := `
		UPDATE journal_entries
		SET content = $1, emotional_state = $2
		WHERE id = $3 AND user_id = $4
		RETURNING updated_at
	`

	err := db.QueryRowContext(
		ctx,
		query,
		entry.Content,
		entry.EmotionalState,
		entry.ID,
		entry.UserID,
	).Scan(&entry.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("journal entry not found")
		}
		return err
	}

	return nil
}

// DeleteJournalEntry deletes a journal entry
func (db *DB) DeleteJournalEntry(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM journal_entries WHERE id = $1 AND user_id = $2`

	result, err := db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("journal entry not found")
	}

	return nil
}

// GetAttachmentsByEntryID retrieves all attachments for a journal entry
func (db *DB) GetAttachmentsByEntryID(ctx context.Context, entryID uuid.UUID) ([]models.Attachment, error) {
	query := `
		SELECT id, entry_id, attachment_type, storage_path, filename, file_size, mime_type, uploaded_at
		FROM attachments
		WHERE entry_id = $1
		ORDER BY uploaded_at ASC
	`

	rows, err := db.QueryContext(ctx, query, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.Attachment
	for rows.Next() {
		var attachment models.Attachment
		err := rows.Scan(
			&attachment.ID,
			&attachment.EntryID,
			&attachment.Type,
			&attachment.StoragePath,
			&attachment.Filename,
			&attachment.FileSize,
			&attachment.MimeType,
			&attachment.UploadedAt,
		)
		if err != nil {
			return nil, err
		}

		// Generate URL for the attachment
		attachment.URL = fmt.Sprintf("/api/attachments/%s", attachment.ID)

		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

// CreateAttachment creates a new attachment record
func (db *DB) CreateAttachment(ctx context.Context, attachment *models.Attachment) error {
	query := `
		INSERT INTO attachments (id, entry_id, attachment_type, storage_path, filename, file_size, mime_type, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, uploaded_at
	`

	attachment.ID = uuid.New()

	err := db.QueryRowContext(
		ctx,
		query,
		attachment.ID,
		attachment.EntryID,
		attachment.Type,
		attachment.StoragePath,
		attachment.Filename,
		attachment.FileSize,
		attachment.MimeType,
	).Scan(&attachment.ID, &attachment.UploadedAt)

	if err != nil {
		return err
	}

	attachment.URL = fmt.Sprintf("/api/attachments/%s", attachment.ID)
	return nil
}

// GetAttachment retrieves an attachment by ID
func (db *DB) GetAttachment(ctx context.Context, id uuid.UUID) (*models.Attachment, error) {
	query := `
		SELECT id, entry_id, attachment_type, storage_path, filename, file_size, mime_type, uploaded_at
		FROM attachments
		WHERE id = $1
	`

	attachment := &models.Attachment{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&attachment.ID,
		&attachment.EntryID,
		&attachment.Type,
		&attachment.StoragePath,
		&attachment.Filename,
		&attachment.FileSize,
		&attachment.MimeType,
		&attachment.UploadedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("attachment not found")
		}
		return nil, err
	}

	attachment.URL = fmt.Sprintf("/api/attachments/%s", attachment.ID)
	return attachment, nil
}

// DeleteAttachment deletes an attachment
func (db *DB) DeleteAttachment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// First verify the user owns the journal entry associated with this attachment
	query := `
		DELETE FROM attachments
		WHERE id = $1
		AND entry_id IN (
			SELECT id FROM journal_entries WHERE user_id = $2
		)
	`

	result, err := db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("attachment not found or access denied")
	}

	return nil
}

// GetJournalEntriesByTradeID retrieves all journal entries for a specific trade
func (db *DB) GetJournalEntriesByTradeID(ctx context.Context, tradeID, userID uuid.UUID) ([]models.JournalEntry, error) {
	query := `
		SELECT id, trade_id, user_id, content, emotional_state, created_at, updated_at
		FROM journal_entries
		WHERE trade_id = $1 AND user_id = $2
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query, tradeID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.JournalEntry
	for rows.Next() {
		var entry models.JournalEntry
		err := rows.Scan(
			&entry.ID,
			&entry.TradeID,
			&entry.UserID,
			&entry.Content,
			&entry.EmotionalState,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Load attachments for each entry
		attachments, err := db.GetAttachmentsByEntryID(ctx, entry.ID)
		if err != nil {
			return nil, err
		}
		entry.Attachments = attachments

		entries = append(entries, entry)
	}

	return entries, nil
}

// Helper function to marshal/unmarshal JSONB data
func marshalJSON(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
