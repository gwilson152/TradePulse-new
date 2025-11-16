package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tradepulse/api/internal/database"
	"github.com/tradepulse/api/internal/middleware"
	"github.com/tradepulse/api/internal/models"
	"github.com/tradepulse/api/internal/notifications"
)

type CSVImportHandler struct {
	db  *database.DB
	bus *notifications.Bus
}

func NewCSVImportHandler(db *database.DB, bus *notifications.Bus) *CSVImportHandler {
	return &CSVImportHandler{db: db, bus: bus}
}

// ImportCSV handles POST /api/trades/import-csv
func (h *CSVImportHandler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req struct {
		Trades []models.Trade `json:"trades"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if len(req.Trades) == 0 {
		sendError(w, http.StatusBadRequest, "No trades to import", nil)
		return
	}

	// Set user ID for all trades
	for i := range req.Trades {
		req.Trades[i].UserID = userID
	}

	// Bulk insert trades
	tradeIDs, err := h.db.BulkCreateTrades(r.Context(), req.Trades)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to import trades", err)
		return
	}

	// Send notification
	h.bus.Publish(
		notifications.NotificationTypeCSVImport,
		userID,
		"CSV Import Complete",
		"Successfully imported "+string(rune(len(tradeIDs)))+" trades",
		map[string]interface{}{
			"count":     len(tradeIDs),
			"trade_ids": tradeIDs,
		},
	)

	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"imported_count": len(tradeIDs),
			"trade_ids":      tradeIDs,
		},
	})
}
