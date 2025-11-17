package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/tradepulse/api/internal/integrations"
)

type FetchPropReportsInput struct {
	Site     string `json:"site"`
	Username string `json:"username"`
	Password string `json:"password"`
	FromDate string `json:"from_date,omitempty"` // Optional: YYYY-MM-DD format
	ToDate   string `json:"to_date,omitempty"`   // Optional: YYYY-MM-DD format
}

// FetchPropReportsTrades fetches trades from PropReports API
func FetchPropReportsTrades(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input FetchPropReportsInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
			return
		}

		if input.Site == "" || input.Username == "" || input.Password == "" {
			writeError(w, http.StatusBadRequest, "MISSING_CREDENTIALS", "Site, username, and password are required")
			return
		}

		// Fetch trades from PropReports
		client := integrations.NewPropReportsClient(input.Site, input.Username, input.Password)
		trades, err := client.FetchTrades(input.FromDate, input.ToDate)
		if err != nil {
			logger.Error("Failed to fetch PropReports trades", "error", err, "site", input.Site, "username", input.Username)
			writeError(w, http.StatusInternalServerError, "FETCH_ERROR", "Failed to fetch trades from PropReports: "+err.Error())
			return
		}

		logger.Info("Successfully fetched trades from PropReports", "site", input.Site, "username", input.Username, "count", len(trades))

		// Log sample of trades for debugging
		if len(trades) > 0 {
			logger.Info("=== Sample of Imported Trades ===")
			sampleSize := 5
			if len(trades) < sampleSize {
				sampleSize = len(trades)
			}
			for i := 0; i < sampleSize; i++ {
				t := trades[i]
				logger.Info("Trade Sample",
					"index", i+1,
					"symbol", t.Symbol,
					"type", t.TradeType,
					"quantity", t.Quantity,
					"entry_price", t.EntryPrice,
					"exit_price", t.ExitPrice,
					"opened_at", t.OpenedAt.Format("2006-01-02 15:04:05"),
					"closed_at", func() string {
						if t.ClosedAt != nil {
							return t.ClosedAt.Format("2006-01-02 15:04:05")
						}
						return "OPEN"
					}(),
					"pnl", t.PnL,
					"fees", t.Fees,
				)
			}
			if len(trades) > sampleSize {
				logger.Info("... and more", "remaining", len(trades)-sampleSize)
			}
			logger.Info("=================================")
		}

		writeSuccess(w, http.StatusOK, trades)
	}
}
