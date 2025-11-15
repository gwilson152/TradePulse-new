package notifications

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client represents a WebSocket client connection
type Client struct {
	UserID uuid.UUID
	conn   *websocket.Conn
	send   chan *Notification
	bus    *Bus
	logger *slog.Logger
}

// NewClient creates a new WebSocket client
func NewClient(userID uuid.UUID, conn *websocket.Conn, bus *Bus, logger *slog.Logger) *Client {
	return &Client{
		UserID: userID,
		conn:   conn,
		send:   make(chan *Notification, 256),
		bus:    bus,
		logger: logger,
	}
}

// ReadPump pumps messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.bus.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("WebSocket read error", "error", err, "user_id", c.UserID)
			}
			break
		}
		// We don't expect any messages from the client, just ignore them
	}
}

// WritePump pumps messages from the bus to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case notification, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The bus closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Marshal notification to JSON
			data, err := MarshalNotification(notification)
			if err != nil {
				c.logger.Error("Failed to marshal notification", "error", err)
				continue
			}

			// Send notification
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.logger.Error("Failed to write notification", "error", err, "user_id", c.UserID)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
