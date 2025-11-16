package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/database"
	"github.com/tradepulse/api/internal/middleware"
	"github.com/tradepulse/api/internal/models"
	"github.com/tradepulse/api/internal/notifications"
)

// Helper functions for JSON responses
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		slog.Error(message, "error", err)
	}
	sendJSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]string{
			"message": message,
		},
	})
}

// getUserID extracts user ID from context (set by auth middleware)
func getUserID(r *http.Request) uuid.UUID {
	userID, _ := middleware.GetUserID(r)
	return userID
}

type TradesHandler struct {
	db  *database.DB
	bus *notifications.Bus
}

func NewTradesHandler(db *database.DB, bus *notifications.Bus) *TradesHandler {
	return &TradesHandler{db: db, bus: bus}
}

// ListTrades handles GET /api/trades
func (h *TradesHandler) ListTrades(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	// Parse query parameters
	filters := database.TradeFilters{
		Symbol:    r.URL.Query().Get("symbol"),
		TradeType: r.URL.Query().Get("trade_type"),
		Status:    r.URL.Query().Get("status"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	// Parse pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	trades, err := h.db.ListTrades(r.Context(), userID, filters)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch trades", err)
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    trades,
	})
}

// GetTrade handles GET /api/trades/{id}
func (h *TradesHandler) GetTrade(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	tradeID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid trade ID", err)
		return
	}

	trade, err := h.db.GetTrade(r.Context(), tradeID, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch trade", err)
		return
	}
	if trade == nil {
		sendError(w, http.StatusNotFound, "Trade not found", nil)
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    trade,
	})
}

// CreateTrade handles POST /api/trades
func (h *TradesHandler) CreateTrade(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	var trade models.Trade
	if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Set user ID
	trade.UserID = userID

	// Validate required fields
	if trade.Symbol == "" || trade.Quantity <= 0 || trade.EntryPrice <= 0 {
		sendError(w, http.StatusBadRequest, "Missing required fields", nil)
		return
	}

	// Validate trade type
	if trade.TradeType != models.TradeLong && trade.TradeType != models.TradeShort {
		sendError(w, http.StatusBadRequest, "Invalid trade type", nil)
		return
	}

	// Create trade
	if err := h.db.CreateTrade(r.Context(), &trade); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create trade", err)
		return
	}

	// Send notification
	h.bus.Publish(
		notifications.NotificationTypeTradeCreated,
		userID,
		"Trade Created",
		"New trade added successfully",
		map[string]interface{}{
			"id":     trade.ID,
			"symbol": trade.Symbol,
		},
	)

	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    trade,
	})
}

// UpdateTrade handles PUT /api/trades/{id}
func (h *TradesHandler) UpdateTrade(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	tradeID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid trade ID", err)
		return
	}

	var trade models.Trade
	if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate required fields
	if trade.Symbol == "" || trade.Quantity <= 0 || trade.EntryPrice <= 0 {
		sendError(w, http.StatusBadRequest, "Missing required fields", nil)
		return
	}

	// Validate trade type
	if trade.TradeType != models.TradeLong && trade.TradeType != models.TradeShort {
		sendError(w, http.StatusBadRequest, "Invalid trade type", nil)
		return
	}

	// Update trade
	if err := h.db.UpdateTrade(r.Context(), tradeID, userID, &trade); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update trade", err)
		return
	}

	// Send notification
	h.bus.Publish(
		notifications.NotificationTypeTradeUpdated,
		userID,
		"Trade Updated",
		"Trade updated successfully",
		map[string]interface{}{
			"id":     trade.ID,
			"symbol": trade.Symbol,
		},
	)

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    trade,
	})
}

// DeleteTrade handles DELETE /api/trades/{id}
func (h *TradesHandler) DeleteTrade(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	tradeID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid trade ID", err)
		return
	}

	if err := h.db.DeleteTrade(r.Context(), tradeID, userID); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to delete trade", err)
		return
	}

	// Send notification
	h.bus.Publish(
		notifications.NotificationTypeTradeDeleted,
		userID,
		"Trade Deleted",
		"Trade deleted successfully",
		map[string]interface{}{
			"id": tradeID,
		},
	)

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Trade deleted successfully",
	})
}

// AddTagToTrade handles POST /api/trades/{id}/tags
func (h *TradesHandler) AddTagToTrade(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	tradeID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid trade ID", err)
		return
	}

	var req struct {
		TagID string `json:"tag_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	tagID, err := uuid.Parse(req.TagID)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid tag ID", err)
		return
	}

	if err := h.db.AddTagToTrade(r.Context(), tradeID, tagID, userID); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to add tag to trade", err)
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Tag added to trade successfully",
	})
}

// RemoveTagFromTrade handles DELETE /api/trades/{tradeId}/tags/{tagId}
func (h *TradesHandler) RemoveTagFromTrade(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	tradeID, err := uuid.Parse(chi.URLParam(r, "tradeId"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid trade ID", err)
		return
	}

	tagID, err := uuid.Parse(chi.URLParam(r, "tagId"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid tag ID", err)
		return
	}

	if err := h.db.RemoveTagFromTrade(r.Context(), tradeID, tagID, userID); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to remove tag from trade", err)
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Tag removed from trade successfully",
	})
}
