package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/models"
)

// TradeFilters represents filters for listing trades
type TradeFilters struct {
	Symbol     string
	TradeType  string // "LONG" or "SHORT"
	Status     string // "open", "closed", "all"
	StartDate  string // ISO 8601 format
	EndDate    string // ISO 8601 format
	Limit      int
	Offset     int
}

// ListTrades retrieves all trades for a user with optional filters
func (db *DB) ListTrades(ctx context.Context, userID uuid.UUID, filters TradeFilters) ([]models.Trade, error) {
	query := `
		SELECT
			t.id, t.user_id, t.symbol, t.trade_type, t.quantity,
			t.entry_price, t.exit_price, t.fees, t.pnl,
			t.opened_at, t.closed_at, t.created_at, t.updated_at,
			EXISTS(SELECT 1 FROM journal_entries je WHERE je.trade_id = t.id) as has_journal,
			COALESCE(
				(SELECT json_agg(tag.name)
				FROM trade_tags tt
				JOIN tags tag ON tt.tag_id = tag.id
				WHERE tt.trade_id = t.id),
				'[]'::json
			) as tags
		FROM trades t
		WHERE t.user_id = $1`

	args := []interface{}{userID}
	argCount := 1

	// Apply filters
	if filters.Symbol != "" {
		argCount++
		query += fmt.Sprintf(" AND UPPER(t.symbol) = UPPER($%d)", argCount)
		args = append(args, filters.Symbol)
	}

	if filters.TradeType != "" {
		argCount++
		query += fmt.Sprintf(" AND t.trade_type = $%d", argCount)
		args = append(args, filters.TradeType)
	}

	if filters.Status == "open" {
		query += " AND t.exit_price IS NULL"
	} else if filters.Status == "closed" {
		query += " AND t.exit_price IS NOT NULL"
	}

	if filters.StartDate != "" {
		argCount++
		query += fmt.Sprintf(" AND t.opened_at >= $%d", argCount)
		args = append(args, filters.StartDate)
	}

	if filters.EndDate != "" {
		argCount++
		query += fmt.Sprintf(" AND t.opened_at <= $%d", argCount)
		args = append(args, filters.EndDate)
	}

	// Order by most recent first
	query += " ORDER BY t.opened_at DESC"

	// Apply pagination
	if filters.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
	}

	if filters.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filters.Offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list trades: %w", err)
	}
	defer rows.Close()

	var trades []models.Trade
	for rows.Next() {
		var trade models.Trade
		var tagsJSON []byte

		err := rows.Scan(
			&trade.ID, &trade.UserID, &trade.Symbol, &trade.TradeType, &trade.Quantity,
			&trade.EntryPrice, &trade.ExitPrice, &trade.Fees, &trade.PnL,
			&trade.OpenedAt, &trade.ClosedAt, &trade.CreatedAt, &trade.UpdatedAt,
			&trade.HasJournal, &tagsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}

		// Parse tags JSON
		if len(tagsJSON) > 0 && string(tagsJSON) != "[]" {
			// Simple JSON array parsing
			tagsStr := string(tagsJSON)
			tagsStr = strings.Trim(tagsStr, "[]")
			if tagsStr != "" {
				tags := strings.Split(tagsStr, ",")
				for i, tag := range tags {
					tags[i] = strings.Trim(tag, `"`)
				}
				trade.Tags = tags
			}
		}

		trades = append(trades, trade)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trades: %w", err)
	}

	return trades, nil
}

// GetTrade retrieves a single trade by ID
func (db *DB) GetTrade(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Trade, error) {
	query := `
		SELECT
			t.id, t.user_id, t.symbol, t.trade_type, t.quantity,
			t.entry_price, t.exit_price, t.fees, t.pnl,
			t.opened_at, t.closed_at, t.created_at, t.updated_at,
			EXISTS(SELECT 1 FROM journal_entries je WHERE je.trade_id = t.id) as has_journal,
			COALESCE(
				(SELECT json_agg(tag.name)
				FROM trade_tags tt
				JOIN tags tag ON tt.tag_id = tag.id
				WHERE tt.trade_id = t.id),
				'[]'::json
			) as tags
		FROM trades t
		WHERE t.id = $1 AND t.user_id = $2`

	var trade models.Trade
	var tagsJSON []byte

	err := db.QueryRow(query, id, userID).Scan(
		&trade.ID, &trade.UserID, &trade.Symbol, &trade.TradeType, &trade.Quantity,
		&trade.EntryPrice, &trade.ExitPrice, &trade.Fees, &trade.PnL,
		&trade.OpenedAt, &trade.ClosedAt, &trade.CreatedAt, &trade.UpdatedAt,
		&trade.HasJournal, &tagsJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	// Parse tags JSON
	if len(tagsJSON) > 0 && string(tagsJSON) != "[]" {
		tagsStr := string(tagsJSON)
		tagsStr = strings.Trim(tagsStr, "[]")
		if tagsStr != "" {
			tags := strings.Split(tagsStr, ",")
			for i, tag := range tags {
				tags[i] = strings.Trim(tag, `"`)
			}
			trade.Tags = tags
		}
	}

	return &trade, nil
}

// CreateTrade inserts a new trade
func (db *DB) CreateTrade(ctx context.Context, trade *models.Trade) error {
	query := `
		INSERT INTO trades (
			user_id, symbol, trade_type, quantity, entry_price, exit_price,
			fees, opened_at, closed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, pnl, created_at, updated_at`

	err := db.QueryRow(
		query,
		trade.UserID, trade.Symbol, trade.TradeType, trade.Quantity,
		trade.EntryPrice, trade.ExitPrice, trade.Fees, trade.OpenedAt, trade.ClosedAt,
	).Scan(&trade.ID, &trade.PnL, &trade.CreatedAt, &trade.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create trade: %w", err)
	}

	return nil
}

// UpdateTrade updates an existing trade
func (db *DB) UpdateTrade(ctx context.Context, id uuid.UUID, userID uuid.UUID, trade *models.Trade) error {
	query := `
		UPDATE trades
		SET symbol = $3, trade_type = $4, quantity = $5, entry_price = $6,
		    exit_price = $7, fees = $8, opened_at = $9, closed_at = $10
		WHERE id = $1 AND user_id = $2
		RETURNING pnl, updated_at`

	err := db.QueryRow(
		query,
		id, userID, trade.Symbol, trade.TradeType, trade.Quantity,
		trade.EntryPrice, trade.ExitPrice, trade.Fees, trade.OpenedAt, trade.ClosedAt,
	).Scan(&trade.PnL, &trade.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("trade not found or unauthorized")
	}
	if err != nil {
		return fmt.Errorf("failed to update trade: %w", err)
	}

	trade.ID = id
	trade.UserID = userID

	return nil
}

// DeleteTrade deletes a trade
func (db *DB) DeleteTrade(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM trades WHERE id = $1 AND user_id = $2`

	result, err := db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete trade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("trade not found or unauthorized")
	}

	return nil
}

// BulkCreateTrades inserts multiple trades (for CSV import)
func (db *DB) BulkCreateTrades(ctx context.Context, trades []models.Trade) ([]uuid.UUID, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt := `
		INSERT INTO trades (
			user_id, symbol, trade_type, quantity, entry_price, exit_price,
			fees, opened_at, closed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	ids := make([]uuid.UUID, 0, len(trades))

	for _, trade := range trades {
		var id uuid.UUID
		err := tx.QueryRow(
			stmt,
			trade.UserID, trade.Symbol, trade.TradeType, trade.Quantity,
			trade.EntryPrice, trade.ExitPrice, trade.Fees, trade.OpenedAt, trade.ClosedAt,
		).Scan(&id)

		if err != nil {
			return nil, fmt.Errorf("failed to insert trade: %w", err)
		}

		ids = append(ids, id)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ids, nil
}

// AddTagToTrade associates a tag with a trade
func (db *DB) AddTagToTrade(ctx context.Context, tradeID uuid.UUID, tagID uuid.UUID, userID uuid.UUID) error {
	// Verify trade ownership
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM trades WHERE id = $1 AND user_id = $2)",
		tradeID, userID,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to verify trade ownership: %w", err)
	}
	if !exists {
		return fmt.Errorf("trade not found or unauthorized")
	}

	// Insert tag association
	query := `
		INSERT INTO trade_tags (trade_id, tag_id)
		VALUES ($1, $2)
		ON CONFLICT (trade_id, tag_id) DO NOTHING`

	_, err = db.Exec(query, tradeID, tagID)
	if err != nil {
		return fmt.Errorf("failed to add tag to trade: %w", err)
	}

	return nil
}

// RemoveTagFromTrade removes a tag association from a trade
func (db *DB) RemoveTagFromTrade(ctx context.Context, tradeID uuid.UUID, tagID uuid.UUID, userID uuid.UUID) error {
	// Verify trade ownership
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM trades WHERE id = $1 AND user_id = $2)",
		tradeID, userID,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to verify trade ownership: %w", err)
	}
	if !exists {
		return fmt.Errorf("trade not found or unauthorized")
	}

	query := `DELETE FROM trade_tags WHERE trade_id = $1 AND tag_id = $2`

	result, err := db.Exec(query, tradeID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from trade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tag association not found")
	}

	return nil
}
