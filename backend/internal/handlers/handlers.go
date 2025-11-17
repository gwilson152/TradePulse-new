package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/database"
	"github.com/tradepulse/api/internal/middleware"
	"github.com/tradepulse/api/internal/models"
)

// APIResponse is the standard response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeSuccess writes a successful JSON response
func writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, APIResponse{
		Success: true,
		Data:    data,
	})
}

// writeError writes an error JSON response
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// Placeholder handlers - these will be implemented in separate files

func Logout(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeSuccess(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
	}
}

func RefreshToken(logger *slog.Logger, jwtSecret, jwtExpiry string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func ListTrades(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func CreateTrade(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func GetTrade(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func UpdateTrade(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func DeleteTrade(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func ImportCSV(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func AddTradeTag(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func RemoveTradeTag(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func ListJournalEntries(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		// Parse pagination parameters
		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}
		if limit <= 0 || limit > 100 {
			limit = 20
		}

		offset := 0
		if o := r.URL.Query().Get("offset"); o != "" {
			fmt.Sscanf(o, "%d", &offset)
		}
		if offset < 0 {
			offset = 0
		}

		// Get journal entries
		entries, total, err := db.ListJournalEntries(r.Context(), userID, limit, offset)
		if err != nil {
			logger.Error("Failed to list journal entries", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to list journal entries")
			return
		}

		writeSuccess(w, http.StatusOK, map[string]interface{}{
			"entries": entries,
			"total":   total,
		})
	}
}

func CreateJournalEntry(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		var input struct {
			TradeID        *string                `json:"trade_id"`
			Content        string                 `json:"content"`
			EmotionalState map[string]interface{} `json:"emotional_state"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Failed to decode request body", "error", err)
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		if input.Content == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Content is required")
			return
		}

		entry := &models.JournalEntry{
			UserID:  userID,
			Content: input.Content,
		}

		// Parse trade ID if provided
		if input.TradeID != nil && *input.TradeID != "" {
			tradeID, err := uuid.Parse(*input.TradeID)
			if err != nil {
				writeError(w, http.StatusBadRequest, "INVALID_TRADE_ID", "Invalid trade ID format")
				return
			}
			entry.TradeID = tradeID
		} else {
			entry.TradeID = uuid.Nil
		}

		// Marshal emotional state to JSON string
		if input.EmotionalState != nil {
			emotionalStateJSON, err := json.Marshal(input.EmotionalState)
			if err != nil {
				logger.Error("Failed to marshal emotional state", "error", err)
				writeError(w, http.StatusInternalServerError, "PROCESSING_ERROR", "Failed to process emotional state")
				return
			}
			entry.EmotionalState = string(emotionalStateJSON)
		}

		if err := db.CreateJournalEntry(r.Context(), entry); err != nil {
			logger.Error("Failed to create journal entry", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to create journal entry")
			return
		}

		logger.Info("Journal entry created", "id", entry.ID, "user_id", userID)
		writeSuccess(w, http.StatusCreated, entry)
	}
}

func GetJournalEntry(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid journal entry ID")
			return
		}

		entry, err := db.GetJournalEntry(r.Context(), id, userID)
		if err != nil {
			logger.Error("Failed to get journal entry", "error", err, "id", id)
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Journal entry not found")
			return
		}

		writeSuccess(w, http.StatusOK, entry)
	}
}

func UpdateJournalEntry(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func DeleteJournalEntry(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func GetJournalEntriesByTradeID(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		// Note: chi router will need the import
		// For now, get from query param as fallback
		tradeIDStr := r.URL.Query().Get("tradeId")

		tradeID, err := uuid.Parse(tradeIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_TRADE_ID", "Invalid trade ID")
			return
		}

		entries, err := db.GetJournalEntriesByTradeID(r.Context(), tradeID, userID)
		if err != nil {
			logger.Error("Failed to get journal entries for trade", "error", err, "trade_id", tradeID)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to get journal entries")
			return
		}

		writeSuccess(w, http.StatusOK, entries)
	}
}

func UploadAttachment(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func GetAttachment(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func DeleteAttachment(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func ListTags(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func CreateTag(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func GetSummaryMetrics(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func GetMetricsBySymbol(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

func GetDailyPerformance(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not yet implemented")
	}
}

// RuleSet handlers
func ListRuleSets(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		ruleSets, err := db.ListRuleSets(r.Context(), userID)
		if err != nil {
			logger.Error("Failed to list rule sets", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to list rule sets")
			return
		}

		writeSuccess(w, http.StatusOK, ruleSets)
	}
}

func CreateRuleSet(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		var input struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			IsActive    bool   `json:"is_active"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		if input.Name == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
			return
		}

		ruleSet := &models.RuleSet{
			UserID:      userID,
			Name:        input.Name,
			Description: input.Description,
			IsActive:    input.IsActive,
		}

		if err := db.CreateRuleSet(r.Context(), ruleSet); err != nil {
			logger.Error("Failed to create rule set", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to create rule set")
			return
		}

		writeSuccess(w, http.StatusCreated, ruleSet)
	}
}

func GetRuleSet(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rule set ID")
			return
		}

		ruleSet, err := db.GetRuleSet(r.Context(), id, userID)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Rule set not found")
			return
		}

		writeSuccess(w, http.StatusOK, ruleSet)
	}
}

func UpdateRuleSet(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rule set ID")
			return
		}

		var input struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			IsActive    bool   `json:"is_active"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		ruleSet := &models.RuleSet{
			ID:          id,
			UserID:      userID,
			Name:        input.Name,
			Description: input.Description,
			IsActive:    input.IsActive,
		}

		if err := db.UpdateRuleSet(r.Context(), ruleSet); err != nil {
			logger.Error("Failed to update rule set", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update rule set")
			return
		}

		writeSuccess(w, http.StatusOK, ruleSet)
	}
}

func DeleteRuleSet(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rule set ID")
			return
		}

		if err := db.DeleteRuleSet(r.Context(), id, userID); err != nil {
			logger.Error("Failed to delete rule set", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to delete rule set")
			return
		}

		writeSuccess(w, http.StatusOK, map[string]string{"message": "rule set deleted successfully"})
	}
}

func CreateRule(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		ruleSetIDStr := r.URL.Query().Get("ruleSetId")
		ruleSetID, err := uuid.Parse(ruleSetIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rule set ID")
			return
		}

		// Verify user owns this rule set
		_, err = db.GetRuleSet(r.Context(), ruleSetID, userID)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Rule set not found")
			return
		}

		var input struct {
			Title       string              `json:"title"`
			Description string              `json:"description"`
			Weight      int                 `json:"weight"`
			Phase       models.RulePhase    `json:"phase"`
			Category    models.RuleCategory `json:"category"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		rule := &models.Rule{
			RuleSetID:   ruleSetID,
			Title:       input.Title,
			Description: input.Description,
			Weight:      input.Weight,
			Phase:       input.Phase,
			Category:    input.Category,
		}

		if err := db.CreateRule(r.Context(), rule); err != nil {
			logger.Error("Failed to create rule", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to create rule")
			return
		}

		writeSuccess(w, http.StatusCreated, rule)
	}
}

func UpdateRule(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		ruleSetIDStr := r.URL.Query().Get("ruleSetId")
		ruleSetID, err := uuid.Parse(ruleSetIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rule set ID")
			return
		}

		ruleIDStr := r.URL.Query().Get("ruleId")
		ruleID, err := uuid.Parse(ruleIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rule ID")
			return
		}

		// Verify user owns this rule set
		_, err = db.GetRuleSet(r.Context(), ruleSetID, userID)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Rule set not found")
			return
		}

		var input struct {
			Title       string              `json:"title"`
			Description string              `json:"description"`
			Weight      int                 `json:"weight"`
			Phase       models.RulePhase    `json:"phase"`
			Category    models.RuleCategory `json:"category"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		rule := &models.Rule{
			ID:          ruleID,
			RuleSetID:   ruleSetID,
			Title:       input.Title,
			Description: input.Description,
			Weight:      input.Weight,
			Phase:       input.Phase,
			Category:    input.Category,
		}

		if err := db.UpdateRule(r.Context(), rule); err != nil {
			logger.Error("Failed to update rule", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update rule")
			return
		}

		writeSuccess(w, http.StatusOK, rule)
	}
}

func DeleteRule(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		ruleSetIDStr := r.URL.Query().Get("ruleSetId")
		ruleSetID, err := uuid.Parse(ruleSetIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rule set ID")
			return
		}

		ruleIDStr := r.URL.Query().Get("ruleId")
		ruleID, err := uuid.Parse(ruleIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rule ID")
			return
		}

		// Verify user owns this rule set
		_, err = db.GetRuleSet(r.Context(), ruleSetID, userID)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Rule set not found")
			return
		}

		if err := db.DeleteRule(r.Context(), ruleID); err != nil {
			logger.Error("Failed to delete rule", "error", err)
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to delete rule")
			return
		}

		writeSuccess(w, http.StatusOK, map[string]string{"message": "rule deleted successfully"})
	}
}
