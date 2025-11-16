package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/models"
)

// GetUserByEmail retrieves a user by email
func (db *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	var lastLogin sql.NullTime
	var passwordHash sql.NullString

	err := db.QueryRowContext(ctx, `
		SELECT id, email, COALESCE(password_hash, '') as password_hash, created_at, last_login
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &passwordHash, &user.CreatedAt, &lastLogin)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	if passwordHash.Valid && passwordHash.String != "" {
		user.PasswordHash = passwordHash.String
		user.HasPassword = true
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	var lastLogin sql.NullTime
	var passwordHash sql.NullString

	err := db.QueryRowContext(ctx, `
		SELECT id, email, COALESCE(password_hash, '') as password_hash, created_at, last_login
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Email, &passwordHash, &user.CreatedAt, &lastLogin)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	if passwordHash.Valid && passwordHash.String != "" {
		user.PasswordHash = passwordHash.String
		user.HasPassword = true
	}

	return &user, nil
}

// CreateUser creates a new user
func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (id, email, created_at, last_login)
		VALUES ($1, $2, $3, $4)
	`, user.ID, user.Email, user.CreatedAt, user.LastLogin)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateUserLastLogin updates the user's last login timestamp
func (db *DB) UpdateUserLastLogin(ctx context.Context, id uuid.UUID) error {
	_, err := db.ExecContext(ctx, `
		UPDATE users
		SET last_login = $1
		WHERE id = $2
	`, time.Now(), id)

	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// StoreMagicLinkToken stores a magic link token for a user
func (db *DB) StoreMagicLinkToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	// Delete any existing tokens for this user
	_, err := db.ExecContext(ctx, `
		DELETE FROM magic_links
		WHERE user_id = $1 AND used_at IS NULL
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to delete old tokens: %w", err)
	}

	// Insert new token
	_, err = db.ExecContext(ctx, `
		INSERT INTO magic_links (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`, userID, token, expiresAt, time.Now())

	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	return nil
}

// VerifyMagicLinkToken verifies a magic link token and returns the user ID
// The token is marked as used after verification
func (db *DB) VerifyMagicLinkToken(ctx context.Context, token string) (uuid.UUID, error) {
	var userID uuid.UUID
	var expiresAt time.Time
	var usedAt sql.NullTime

	err := db.QueryRowContext(ctx, `
		SELECT user_id, expires_at, used_at
		FROM magic_links
		WHERE token = $1
	`, token).Scan(&userID, &expiresAt, &usedAt)

	if err == sql.ErrNoRows {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to verify token: %w", err)
	}

	// Check if token has already been used
	if usedAt.Valid {
		return uuid.Nil, fmt.Errorf("token has already been used")
	}

	// Check if token has expired
	if time.Now().After(expiresAt) {
		return uuid.Nil, fmt.Errorf("token has expired")
	}

	// Mark token as used (one-time use)
	_, err = db.ExecContext(ctx, `
		UPDATE magic_links
		SET used_at = $1
		WHERE token = $2
	`, time.Now(), token)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to consume token: %w", err)
	}

	return userID, nil
}

// SetUserPassword sets or updates a user's password hash
func (db *DB) SetUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2
	`, passwordHash, userID)

	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	return nil
}

// VerifyUserPassword retrieves the password hash for a user to verify against
func (db *DB) VerifyUserPassword(ctx context.Context, email string, passwordHash string) (*models.User, error) {
	var user models.User
	var lastLogin sql.NullTime
	var storedHash sql.NullString

	err := db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, created_at, last_login
		FROM users
		WHERE email = $1 AND password_hash IS NOT NULL
	`, email).Scan(&user.ID, &user.Email, &storedHash, &user.CreatedAt, &lastLogin)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found or no password set")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	if storedHash.Valid {
		user.PasswordHash = storedHash.String
		user.HasPassword = true
	}

	return &user, nil
}
