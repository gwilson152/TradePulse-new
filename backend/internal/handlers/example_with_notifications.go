package handlers

// Example of how to use the notification bus in handlers
// This file demonstrates the pattern for publishing notifications

/*
import (
	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/notifications"
)

// Example: Publishing a notification when a trade is created
func CreateTradeWithNotification(db *database.DB, logger *slog.Logger, bus *notifications.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := middleware.GetUserID(r)

		// ... create trade logic ...

		// Publish notification
		bus.Publish(
			notifications.NotificationTypeTradeCreated,
			userID,
			"Trade Created",
			fmt.Sprintf("Successfully created trade for %s", trade.Symbol),
			map[string]interface{}{
				"trade_id": trade.ID,
				"symbol":   trade.Symbol,
				"pnl":      trade.PnL,
			},
		)

		writeSuccess(w, http.StatusCreated, trade)
	}
}

// Example: Publishing an error notification
func handleError(bus *notifications.Bus, userID uuid.UUID, title, message string) {
	bus.Publish(
		notifications.NotificationTypeError,
		userID,
		title,
		message,
		nil,
	)
}

// Example: Publishing CSV import progress
func ImportCSVWithNotifications(db *database.DB, logger *slog.Logger, bus *notifications.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := middleware.GetUserID(r)

		// Start import
		bus.Publish(
			notifications.NotificationTypeInfo,
			userID,
			"CSV Import Started",
			"Processing your CSV file...",
			nil,
		)

		// ... import logic ...

		// Completion notification
		bus.Publish(
			notifications.NotificationTypeCSVImport,
			userID,
			"CSV Import Complete",
			fmt.Sprintf("Successfully imported %d trades", importedCount),
			map[string]interface{}{
				"imported": importedCount,
				"failed":   failedCount,
			},
		)

		writeSuccess(w, http.StatusOK, result)
	}
}
*/
