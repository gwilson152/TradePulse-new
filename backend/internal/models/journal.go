package models

import (
	"time"

	"github.com/google/uuid"
)

type AttachmentType string

const (
	AttachmentScreenshot AttachmentType = "screenshot"
	AttachmentVoice      AttachmentType = "voice"
)

type JournalEntry struct {
	ID             uuid.UUID       `json:"id"`
	TradeID        uuid.UUID       `json:"trade_id"`
	UserID         uuid.UUID       `json:"user_id"`
	Content        string          `json:"content,omitempty"`
	EmotionalState string          `json:"emotional_state,omitempty"` // JSONB stored as string
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Attachments    []Attachment    `json:"attachments,omitempty"`
}

type Attachment struct {
	ID         uuid.UUID      `json:"id"`
	EntryID    uuid.UUID      `json:"entry_id"`
	Type       AttachmentType `json:"attachment_type"`
	StoragePath string        `json:"-"` // Don't expose internal path
	Filename   string         `json:"filename"`
	FileSize   int64          `json:"file_size,omitempty"`
	MimeType   string         `json:"mime_type,omitempty"`
	URL        string         `json:"url"`
	UploadedAt time.Time      `json:"uploaded_at"`
}

type EmotionalState struct {
	PreTradeConfidence   int    `json:"pre_trade_confidence,omitempty"`   // 1-10
	PreTradeClarity      int    `json:"pre_trade_clarity,omitempty"`      // 1-10
	PostTradeDiscipline  int    `json:"post_trade_discipline,omitempty"`  // 1-10
	PostTradeEmotion     string `json:"post_trade_emotion,omitempty"`     // free text or enum
}
