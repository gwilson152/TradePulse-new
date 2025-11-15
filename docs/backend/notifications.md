# Notifications

## Overview

The notification bus broadcasts real-time messages to connected WebSocket clients.

## Publishing Notifications

### Basic Usage

```go
bus.Publish(
    notifications.NotificationTypeTradeCreated,
    userID,
    "Trade Created",
    "Successfully created trade for AAPL",
    map[string]interface{}{
        "trade_id": tradeID,
        "symbol": "AAPL",
    },
)
```

### Notification Types

| Type | Usage |
|------|-------|
| `NotificationTypeTradeCreated` | New trade created |
| `NotificationTypeTradeUpdated` | Trade modified |
| `NotificationTypeTradeDeleted` | Trade deleted |
| `NotificationTypeJournalCreated` | New journal entry |
| `NotificationTypeJournalUpdated` | Journal modified |
| `NotificationTypeCSVImport` | CSV import completed |
| `NotificationTypeError` | Error occurred |
| `NotificationTypeInfo` | Informational message |
| `NotificationTypeSuccess` | Success message |

## Examples

### Trade Created
```go
bus.Publish(
    notifications.NotificationTypeTradeCreated,
    userID,
    "Trade Created",
    fmt.Sprintf("Created %s trade for %s", tradeType, symbol),
    map[string]interface{}{
        "trade_id": tradeID,
        "symbol": symbol,
        "type": tradeType,
    },
)
```

### CSV Import Complete
```go
bus.Publish(
    notifications.NotificationTypeCSVImport,
    userID,
    "Import Complete",
    fmt.Sprintf("Successfully imported %d trades", count),
    map[string]interface{}{
        "imported": count,
        "failed": failedCount,
    },
)
```

### Error
```go
bus.Publish(
    notifications.NotificationTypeError,
    userID,
    "Upload Failed",
    "File size exceeds 10MB limit",
    nil,
)
```

### Success
```go
bus.Publish(
    notifications.NotificationTypeSuccess,
    userID,
    "Settings Saved",
    "Your preferences have been updated",
    nil,
)
```

## In Handlers

```go
func CreateTrade(db *database.DB, logger *slog.Logger, bus *notifications.Bus) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, _ := middleware.GetUserID(r)

        // ... create trade logic ...

        // Send notification
        bus.Publish(
            notifications.NotificationTypeTradeCreated,
            userID,
            "Trade Created",
            fmt.Sprintf("Created trade for %s", trade.Symbol),
            map[string]interface{}{
                "trade_id": trade.ID,
                "symbol": trade.Symbol,
            },
        )

        writeSuccess(w, http.StatusCreated, trade)
    }
}
```

## Best Practices

1. **Always include userID** - Notifications are user-specific
2. **Keep titles short** - 3-5 words max
3. **Messages should be actionable** - Tell user what happened
4. **Include relevant data** - IDs, names, counts
5. **Use appropriate types** - Match notification to action
6. **Don't overuse** - Only for significant events
7. **Handle failures gracefully** - If bus is full, log but don't fail request

## Message Structure

```json
{
  "id": "generated-uuid",
  "type": "trade.created",
  "user_id": "user-uuid",
  "title": "Trade Created",
  "message": "Successfully created trade for AAPL",
  "data": {
    "trade_id": "trade-uuid",
    "symbol": "AAPL"
  },
  "timestamp": "2024-01-20T15:30:00Z"
}
```

## Connection Stats

Check how many users are connected:

```go
stats := bus.GetStats()
// Returns:
// {
//   "total_users": 15,
//   "total_clients": 23,
//   "timestamp": "2024-01-20T15:30:00Z"
// }
```

## See Also

- [WebSocket Documentation](../websocket-notifications.md) - Full WebSocket guide
- [API Patterns](api-patterns.md) - Handler examples
