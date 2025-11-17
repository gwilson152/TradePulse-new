package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"` // Never send password hash to client
	HasPassword    bool       `json:"has_password"` // Indicates if user has set a password
	PlanType       string     `json:"plan_type"`
	PlanStatus     string     `json:"plan_status"`
	PlanSelectedAt *time.Time `json:"plan_selected_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	LastLogin      *time.Time `json:"last_login,omitempty"`
	Preferences    string     `json:"preferences,omitempty"` // JSONB stored as string
}

type MagicLink struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
