package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tradepulse/api/internal/middleware"
	"github.com/tradepulse/api/internal/notifications"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now - this should be restricted in production
		return true
	},
}

// HandleWebSocket upgrades HTTP connection to WebSocket and registers client
func HandleWebSocket(bus *notifications.Bus, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("Failed to upgrade to WebSocket", "error", err, "user_id", userID)
			return
		}

		// Create new client
		client := notifications.NewClient(userID, conn, bus, logger)

		// Register client with the bus
		bus.Register(client)

		logger.Info("WebSocket connection established", "user_id", userID)

		// Start client goroutines
		go client.WritePump()
		go client.ReadPump()
	}
}

// HandleNotificationStats returns statistics about the notification system
func HandleNotificationStats(bus *notifications.Bus, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := bus.GetStats()
		writeSuccess(w, http.StatusOK, stats)
	}
}
