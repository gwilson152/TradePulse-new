package models

import (
	"time"

	"github.com/google/uuid"
)

type TradeType string

const (
	TradeLong  TradeType = "LONG"
	TradeShort TradeType = "SHORT"
)

type Trade struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Symbol     string     `json:"symbol"`
	TradeType  TradeType  `json:"trade_type"`
	Quantity   float64    `json:"quantity"`
	EntryPrice float64    `json:"entry_price"`
	ExitPrice  *float64   `json:"exit_price,omitempty"`
	Fees       float64    `json:"fees"`
	PnL        *float64   `json:"pnl,omitempty"`
	OpenedAt   time.Time  `json:"opened_at"`
	ClosedAt   *time.Time `json:"closed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	HasJournal bool       `json:"has_journal,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
}

type Tag struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	Color      string    `json:"color,omitempty"`
	UsageCount int       `json:"usage_count,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
