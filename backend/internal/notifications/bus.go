package notifications

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeTradeCreated   NotificationType = "trade.created"
	NotificationTypeTradeUpdated   NotificationType = "trade.updated"
	NotificationTypeTradeDeleted   NotificationType = "trade.deleted"
	NotificationTypeJournalCreated NotificationType = "journal.created"
	NotificationTypeJournalUpdated NotificationType = "journal.updated"
	NotificationTypeCSVImport      NotificationType = "csv.import"
	NotificationTypeError          NotificationType = "error"
	NotificationTypeInfo           NotificationType = "info"
	NotificationTypeSuccess        NotificationType = "success"
)

// Notification represents a notification message
type Notification struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	UserID    uuid.UUID        `json:"user_id"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Data      interface{}      `json:"data,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
}

// Bus is the central notification bus
type Bus struct {
	clients  map[uuid.UUID]map[*Client]bool
	mu       sync.RWMutex
	logger   *slog.Logger
	register chan *Client
	unregister chan *Client
	broadcast chan *Notification
}

// NewBus creates a new notification bus
func NewBus(logger *slog.Logger) *Bus {
	return &Bus{
		clients:    make(map[uuid.UUID]map[*Client]bool),
		logger:     logger,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Notification, 256),
	}
}

// Run starts the notification bus
func (b *Bus) Run() {
	for {
		select {
		case client := <-b.register:
			b.mu.Lock()
			if b.clients[client.UserID] == nil {
				b.clients[client.UserID] = make(map[*Client]bool)
			}
			b.clients[client.UserID][client] = true
			b.mu.Unlock()
			b.logger.Info("Client registered", "user_id", client.UserID, "total_clients", b.countClients())

		case client := <-b.unregister:
			b.mu.Lock()
			if clients, ok := b.clients[client.UserID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(b.clients, client.UserID)
					}
				}
			}
			b.mu.Unlock()
			b.logger.Info("Client unregistered", "user_id", client.UserID, "total_clients", b.countClients())

		case notification := <-b.broadcast:
			b.mu.RLock()
			clients := b.clients[notification.UserID]
			b.mu.RUnlock()

			for client := range clients {
				select {
				case client.send <- notification:
				default:
					// Client's send channel is full, close the connection
					b.mu.Lock()
					delete(b.clients[notification.UserID], client)
					close(client.send)
					if len(b.clients[notification.UserID]) == 0 {
						delete(b.clients, notification.UserID)
					}
					b.mu.Unlock()
				}
			}
		}
	}
}

// Publish sends a notification to all connected clients for a user
func (b *Bus) Publish(notificationType NotificationType, userID uuid.UUID, title, message string, data interface{}) {
	notification := &Notification{
		ID:        uuid.New().String(),
		Type:      notificationType,
		UserID:    userID,
		Title:     title,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}

	select {
	case b.broadcast <- notification:
		b.logger.Debug("Notification published", "type", notificationType, "user_id", userID)
	default:
		b.logger.Warn("Notification broadcast channel full, dropping notification", "type", notificationType, "user_id", userID)
	}
}

// Register adds a client to the bus
func (b *Bus) Register(client *Client) {
	b.register <- client
}

// Unregister removes a client from the bus
func (b *Bus) Unregister(client *Client) {
	b.unregister <- client
}

// countClients returns the total number of connected clients
func (b *Bus) countClients() int {
	count := 0
	for _, clients := range b.clients {
		count += len(clients)
	}
	return count
}

// GetStats returns statistics about the notification bus
func (b *Bus) GetStats() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return map[string]interface{}{
		"total_users":   len(b.clients),
		"total_clients": b.countClients(),
		"timestamp":     time.Now(),
	}
}

// MarshalNotification converts a notification to JSON bytes
func MarshalNotification(n *Notification) ([]byte, error) {
	return json.Marshal(n)
}
