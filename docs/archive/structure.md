# Backend Structure

## Package Organization

TradePulse follows Go's standard project layout with an `internal` directory for private code.

```
backend/
├── cmd/api/              # Application entry point
├── internal/             # Private packages
└── migrations/           # Database migrations
```

## cmd/api

**Purpose**: Application entry point

**Files**:
- `main.go` - Initializes dependencies, starts server

**Responsibilities**:
- Load configuration
- Connect to database
- Start notification bus
- Set up router and middleware
- Start HTTP server
- Handle graceful shutdown

## internal/handlers

**Purpose**: HTTP request handlers

**Pattern**:
```go
func HandlerName(db *database.DB, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Handle request
        writeSuccess(w, http.StatusOK, data)
    }
}
```

**Files**:
- `handlers.go` - Helper functions, placeholder handlers
- `websocket.go` - WebSocket upgrade handler
- `auth.go` - Authentication endpoints (to be implemented)
- `trades.go` - Trade CRUD (to be implemented)
- `journal.go` - Journal CRUD (to be implemented)

**Key Functions**:
- `writeSuccess()` - Send success JSON response
- `writeError()` - Send error JSON response
- `writeJSON()` - Generic JSON writer

## internal/middleware

**Purpose**: HTTP middleware functions

**Files**:
- `auth.go` - JWT authentication middleware

**Usage**:
```go
// Protect routes
r.Group(func(r chi.Router) {
    r.Use(middleware.Authenticate(jwtSecret))
    r.Get("/protected", handler)
})

// Get user ID in handler
userID, ok := middleware.GetUserID(r)
```

**Pattern**:
```go
func Middleware(config string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Middleware logic
            next.ServeHTTP(w, r)
        })
    }
}
```

## internal/models

**Purpose**: Data structures

**Files**:
- `user.go` - User, MagicLink
- `trade.go` - Trade, Tag
- `journal.go` - JournalEntry, Attachment, EmotionalState

**Pattern**:
```go
type Model struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    // ... fields
}
```

## internal/database

**Purpose**: Database connection and queries

**Files**:
- `db.go` - Connection management, migration runner

**Usage**:
```go
// Initialize
db, err := database.New(database.Config{
    Host: "postgres1.drivenw.local",
    // ...
})

// Query
row := db.QueryRow("SELECT id FROM users WHERE email = $1", email)
```

**Pattern**:
- Always use parameterized queries ($1, $2, etc.)
- Never concatenate user input into SQL
- Use transactions for multi-step operations

## internal/email

**Purpose**: Email provider abstraction

**Files** (to be created):
- `provider.go` - Email interface
- `smtp.go` - SMTP implementation
- `graph.go` - M365 Graph API
- `gmail.go` - Google API

**Interface**:
```go
type EmailProvider interface {
    Send(to, subject, body string) error
    SendHTML(to, subject, htmlBody string) error
}
```

## internal/notifications

**Purpose**: WebSocket notification system

**Files**:
- `bus.go` - Central notification hub
- `client.go` - WebSocket client manager

**Usage**:
```go
// Publish notification
bus.Publish(
    notifications.NotificationTypeTradeCreated,
    userID,
    "Trade Created",
    "Successfully created trade for AAPL",
    data,
)
```

**Key Components**:
- **Bus**: Manages all connections, broadcasts messages
- **Client**: Individual WebSocket connection with ping/pong

## migrations/

**Purpose**: Database schema changes

**Files**:
- `001_initial_schema.up.sql` - Create tables
- `001_initial_schema.down.sql` - Drop tables

**Naming**: `XXX_description.up.sql` / `XXX_description.down.sql`

**Rules**:
- Never modify existing migrations
- Always create new migrations for changes
- Test both up and down migrations

## Dependency Flow

```
main.go
  ↓
┌─────────────────────────────────────┐
│ Database, Logger, Config, Bus       │
└─────────────────────────────────────┘
  ↓
┌─────────────────────────────────────┐
│ Router + Middleware                 │
└─────────────────────────────────────┘
  ↓
┌─────────────────────────────────────┐
│ Handlers (use DB, Logger, Bus)      │
└─────────────────────────────────────┘
  ↓
┌─────────────────────────────────────┐
│ Models (data structures)            │
└─────────────────────────────────────┘
```

## Adding New Features

### 1. Create Model

```go
// internal/models/feature.go
type Feature struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 2. Create Migration

```sql
-- migrations/002_add_features.up.sql
CREATE TABLE features (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 3. Create Handler

```go
// internal/handlers/features.go
func CreateFeature(db *database.DB, logger *slog.Logger, bus *notifications.Bus) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, _ := middleware.GetUserID(r)

        // Parse request, validate, insert to DB

        // Publish notification
        bus.Publish(
            notifications.NotificationTypeSuccess,
            userID,
            "Feature Created",
            "Successfully created feature",
            nil,
        )

        writeSuccess(w, http.StatusCreated, feature)
    }
}
```

### 4. Register Route

```go
// cmd/api/main.go
r.Group(func(r chi.Router) {
    r.Use(appMiddleware.Authenticate(cfg.jwtSecret))
    r.Post("/features", handlers.CreateFeature(app.db, app.logger, app.notificationBus))
})
```

## Code Style

### Naming
- Packages: lowercase, single word
- Files: snake_case (e.g., `magic_link.go`)
- Exported: PascalCase (e.g., `CreateTrade`)
- Unexported: camelCase (e.g., `validateEmail`)

### Error Handling
```go
if err != nil {
    return fmt.Errorf("failed to do thing: %w", err)
}
```

### Logging
```go
logger.Info("Action completed", "user_id", userID, "trade_id", tradeID)
logger.Error("Action failed", "error", err)
```

### HTTP Responses
```go
// Success
writeSuccess(w, http.StatusOK, data)

// Error
writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Email is required")
```

## Testing

### Unit Tests
```go
func TestHandlerName(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    handler := HandlerName(db, logger)

    // Act
    req := httptest.NewRequest("POST", "/endpoint", body)
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)

    // Assert
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

### Integration Tests
- Test with real database (test instance)
- Use transactions for test isolation
- Clean up after tests

## Best Practices

1. **Always use parameterized queries** - Prevent SQL injection
2. **Log important events** - Help debugging in production
3. **Return specific errors** - Help frontend handle errors
4. **Validate input** - Never trust user data
5. **Use middleware for cross-cutting concerns** - Auth, logging, CORS
6. **Keep handlers thin** - Move logic to separate functions/packages
7. **Use context for request-scoped values** - User ID, request ID
8. **Close resources with defer** - Database connections, files
