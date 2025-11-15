# WebSocket Real-time Notifications

## Overview

TradePulse implements a centralized WebSocket-based notification system that provides real-time updates to users for all platform actions.

## Architecture

### Backend Components

#### 1. Notification Bus (`internal/notifications/bus.go`)
Central hub that manages all WebSocket connections and broadcasts notifications to users.

**Features:**
- Concurrent client management per user
- Message broadcasting with 256 message buffer
- Automatic client cleanup on disconnect
- Connection statistics tracking

**Key Methods:**
```go
bus.Publish(type, userID, title, message, data)  // Broadcast notification
bus.Register(client)                              // Add new WebSocket connection
bus.Unregister(client)                            // Remove connection
bus.GetStats()                                    // Get connection statistics
```

#### 2. Client Manager (`internal/notifications/client.go`)
Handles individual WebSocket connections with ping/pong keepalive.

**Configuration:**
- Write timeout: 10 seconds
- Pong timeout: 60 seconds
- Ping interval: 54 seconds
- Max message size: 512 bytes

#### 3. WebSocket Handler (`internal/handlers/websocket.go`)
HTTP handler that upgrades connections to WebSocket.

**Endpoint:** `GET /api/ws`
**Authentication:** JWT token required (via middleware)

### Frontend Components

#### 1. Notification Store (`lib/stores/notifications.ts`)
Svelte store managing notification state.

**Features:**
- Stores last 50 notifications
- Tracks unread count
- Connection status tracking
- Mark as read/unread functionality

**Methods:**
```typescript
notificationStore.add(notification)      // Add new notification
notificationStore.markAsRead(id)         // Mark notification as read
notificationStore.markAllAsRead()        // Mark all as read
notificationStore.remove(id)             // Remove notification
notificationStore.clear()                // Clear all notifications
notificationStore.setConnected(bool)     // Update connection status
```

#### 2. WebSocket Client (`lib/api/websocket.ts`)
Manages WebSocket connection lifecycle.

**Features:**
- Auto-connect on authentication
- Auto-reconnect with 5-second interval
- Browser notification support (with permission)
- Converts HTTPS → WSS automatically

**Methods:**
```typescript
wsClient.connect(token)      // Establish connection
wsClient.disconnect()        // Close connection
wsClient.isConnected()       // Check connection status
```

#### 3. UI Components

**NotificationBell** (`lib/components/notifications/NotificationBell.svelte`)
- Displays unread count badge
- Shows connection status indicator
- Toggles notification panel

**NotificationPanel** (`lib/components/notifications/NotificationPanel.svelte`)
- Dropdown panel with notification list
- Mark all as read button
- Clear all button
- Connection status footer

**NotificationItem** (`lib/components/notifications/NotificationItem.svelte`)
- Individual notification display
- Color-coded by type
- Dismiss button
- Time formatting (e.g., "5m ago")

## Notification Types

| Type | Description | Icon | Color |
|------|-------------|------|-------|
| `trade.created` | New trade created | ✓ | Green |
| `trade.updated` | Trade modified | ℹ | Blue |
| `trade.deleted` | Trade removed | 🗑 | Gray |
| `journal.created` | New journal entry | ✓ | Green |
| `journal.updated` | Journal entry modified | ℹ | Blue |
| `csv.import` | CSV import completed | 📁 | Purple |
| `error` | Error occurred | ✕ | Red |
| `info` | Informational message | ℹ | Blue |
| `success` | Success message | ✓ | Green |

## Message Format

```json
{
  "id": "uuid",
  "type": "trade.created",
  "user_id": "user-uuid",
  "title": "Trade Created",
  "message": "Successfully created trade for AAPL",
  "data": {
    "trade_id": "trade-uuid",
    "symbol": "AAPL",
    "pnl": 125.50
  },
  "timestamp": "2024-01-20T15:30:00Z"
}
```

## Usage Examples

### Backend: Publishing Notifications

```go
import "github.com/tradepulse/api/internal/notifications"

// In any handler with access to the notification bus
func CreateTrade(db *database.DB, logger *slog.Logger, bus *notifications.Bus) http.HandlerFunc {
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
```

### Frontend: Using Notifications

```svelte
<script>
  import { wsClient } from '$lib/api/websocket';
  import { notificationStore } from '$lib/stores/notifications';
  import NotificationBell from '$lib/components/notifications/NotificationBell.svelte';

  // Connect WebSocket after authentication
  function handleLogin(token) {
    wsClient.connect(token);
  }

  // Disconnect on logout
  function handleLogout() {
    wsClient.disconnect();
    notificationStore.clear();
  }
</script>

<!-- Add notification bell to header -->
<header>
  <NotificationBell />
</header>

<!-- Access notification state -->
{#if $notificationStore.connected}
  <span>Connected</span>
{/if}
<span>Unread: {$notificationStore.unreadCount}</span>
```

## Connection Lifecycle

1. **Authentication** - User logs in and receives JWT token
2. **Connect** - Frontend establishes WebSocket connection with token
3. **Register** - Backend registers client in notification bus
4. **Active** - Bidirectional ping/pong for keepalive
5. **Receive** - Client receives notifications in real-time
6. **Disconnect** - Connection closes (manual or timeout)
7. **Reconnect** - Auto-reconnect after 5 seconds if unintentional

## Browser Notifications

The system supports native browser notifications when permitted by the user.

**Permission Flow:**
1. First notification triggers permission request
2. If granted, subsequent notifications show in OS
3. Notifications include title, message, and app icon
4. Clicking notification focuses the app

## Performance Considerations

- **Message Buffer**: 256 messages buffered per broadcast
- **Client Limit**: No hard limit, scales with server resources
- **Message Size**: Limited to 512 bytes from client (ignored)
- **History**: Frontend stores last 50 notifications only
- **Cleanup**: Automatic client cleanup on disconnect

## Security

- **Authentication**: JWT token required for WebSocket upgrade
- **Authorization**: Users only receive their own notifications
- **Origin Check**: CORS validation on WebSocket upgrade
- **Token Validation**: JWT verified before connection establishment

## Monitoring

### Get Connection Stats

**Endpoint:** `GET /api/notifications/stats`

**Response:**
```json
{
  "success": true,
  "data": {
    "total_users": 15,
    "total_clients": 23,
    "timestamp": "2024-01-20T15:30:00Z"
  }
}
```

### Logs

Backend logs WebSocket events:
- Client registration
- Client unregistration
- Connection errors
- Notification broadcasts (debug level)

## Testing

### Manual Testing

**Backend:**
```bash
# Start server
go run cmd/api/main.go

# Check connection stats
curl -H "Authorization: Bearer <token>" \
  https://api.tradepulse.drivenw.com/api/notifications/stats
```

**Frontend:**
```bash
# Start dev server
npm run dev

# Open browser console to see WebSocket messages
# Look for "WebSocket connected" message
```

### Integration Testing

Use `websocat` or similar tools for WebSocket testing:

```bash
# Connect to WebSocket
websocat "wss://api.tradepulse.drivenw.com/api/ws?token=<jwt_token>"

# Should receive ping/pong messages
# Trigger actions in another session to receive notifications
```

## Troubleshooting

**Connection Issues:**
- Verify JWT token is valid and not expired
- Check CORS configuration allows WebSocket upgrade
- Ensure `allowedHosts` includes domain in Vite config
- Verify external proxy routes WSS correctly

**No Notifications Received:**
- Check notification bus is running (`go notificationBus.Run()`)
- Verify handlers are calling `bus.Publish()`
- Check user ID matches between notification and connected client
- Inspect browser console for WebSocket errors

**Auto-reconnect Not Working:**
- Ensure `shouldReconnect` flag is true
- Check network connectivity
- Verify server is running and accepting connections

## Future Enhancements

- **Notification Persistence**: Store notifications in database
- **Notification Preferences**: User settings for notification types
- **Read Receipts**: Track when notifications are read
- **Push Notifications**: Mobile push notification support
- **Notification Grouping**: Group related notifications
- **Custom Sounds**: Different sounds per notification type
- **Notification Actions**: Interactive buttons in notifications
