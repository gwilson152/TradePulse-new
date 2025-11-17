package models

import (
	"time"

	"github.com/google/uuid"
)

type RulePhase string
type RuleCategory string

const (
	RulePhasePreTrade    RulePhase = "PRE_TRADE"
	RulePhaseDuringTrade RulePhase = "DURING_TRADE"
	RulePhasePostTrade   RulePhase = "POST_TRADE"
)

const (
	RuleCategoryRiskManagement RuleCategory = "RISK_MANAGEMENT"
	RuleCategoryEntry          RuleCategory = "ENTRY"
	RuleCategoryExit           RuleCategory = "EXIT"
	RuleCategoryPositionSizing RuleCategory = "POSITION_SIZING"
	RuleCategoryTiming         RuleCategory = "TIMING"
	RuleCategoryPsychology     RuleCategory = "PSYCHOLOGY"
	RuleCategoryGeneral        RuleCategory = "GENERAL"
)

type RuleSet struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	Rules       []Rule    `json:"rules,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Rule struct {
	ID          uuid.UUID    `json:"id"`
	RuleSetID   uuid.UUID    `json:"rule_set_id,omitempty"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Weight      int          `json:"weight"` // 1-5 (importance)
	Phase       RulePhase    `json:"phase"`
	Category    RuleCategory `json:"category"`
	CreatedAt   time.Time    `json:"created_at"`
}

type RuleAdherence struct {
	RuleID    uuid.UUID `json:"rule_id"`
	RuleTitle string    `json:"rule_title"`
	Score     int       `json:"score"` // 0, 25, 50, 75, 100
	Notes     string    `json:"notes,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
