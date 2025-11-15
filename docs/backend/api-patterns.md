# API Patterns

## Handler Pattern

All handlers follow this structure:

```go
func HandlerName(db *database.DB, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Get user ID from context
        userID, ok := middleware.GetUserID(r)
        if !ok {
            writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
            return
        }

        // 2. Parse and validate request
        var req RequestType
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
            return
        }

        // 3. Perform business logic
        result, err := doSomething(db, userID, req)
        if err != nil {
            logger.Error("Failed to do something", "error", err, "user_id", userID)
            writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Operation failed")
            return
        }

        // 4. Return success response
        writeSuccess(w, http.StatusOK, result)
    }
}
```

## Response Helpers

### Success Response

```go
writeSuccess(w, http.StatusOK, map[string]interface{}{
    "trade": trade,
    "message": "Trade created successfully",
})
```

Response:
```json
{
  "success": true,
  "data": {
    "trade": {...},
    "message": "Trade created successfully"
  }
}
```

### Error Response

```go
writeError(w, http.StatusBadRequest, "INVALID_EMAIL", "Email format is invalid")
```

Response:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_EMAIL",
    "message": "Email format is invalid"
  }
}
```

## Common Patterns

### GET - List Resources

```go
func ListTrades(db *database.DB, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, _ := middleware.GetUserID(r)

        // Parse query parameters
        limit := 50
        if l := r.URL.Query().Get("limit"); l != "" {
            limit, _ = strconv.Atoi(l)
        }

        // Query database
        rows, err := db.Query(`
            SELECT id, symbol, pnl, opened_at
            FROM trades
            WHERE user_id = $1
            ORDER BY opened_at DESC
            LIMIT $2
        `, userID, limit)
        if err != nil {
            logger.Error("Failed to list trades", "error", err)
            writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch trades")
            return
        }
        defer rows.Close()

        // Scan results
        var trades []models.Trade
        for rows.Next() {
            var trade models.Trade
            if err := rows.Scan(&trade.ID, &trade.Symbol, &trade.PnL, &trade.OpenedAt); err != nil {
                continue
            }
            trades = append(trades, trade)
        }

        writeSuccess(w, http.StatusOK, map[string]interface{}{
            "trades": trades,
            "total":  len(trades),
        })
    }
}
```

### POST - Create Resource

```go
func CreateTrade(db *database.DB, logger *slog.Logger, bus *notifications.Bus) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, _ := middleware.GetUserID(r)

        // Parse request
        var req struct {
            Symbol     string  `json:"symbol"`
            TradeType  string  `json:"trade_type"`
            Quantity   float64 `json:"quantity"`
            EntryPrice float64 `json:"entry_price"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
            return
        }

        // Validate
        if req.Symbol == "" {
            writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Symbol is required")
            return
        }

        // Insert to database
        var tradeID uuid.UUID
        err := db.QueryRow(`
            INSERT INTO trades (user_id, symbol, trade_type, quantity, entry_price, opened_at)
            VALUES ($1, $2, $3, $4, $5, NOW())
            RETURNING id
        `, userID, req.Symbol, req.TradeType, req.Quantity, req.EntryPrice).Scan(&tradeID)

        if err != nil {
            logger.Error("Failed to create trade", "error", err)
            writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to create trade")
            return
        }

        // Publish notification
        bus.Publish(
            notifications.NotificationTypeTradeCreated,
            userID,
            "Trade Created",
            fmt.Sprintf("Created trade for %s", req.Symbol),
            map[string]interface{}{"trade_id": tradeID, "symbol": req.Symbol},
        )

        writeSuccess(w, http.StatusCreated, map[string]interface{}{
            "id": tradeID,
            "message": "Trade created successfully",
        })
    }
}
```

### GET - Single Resource

```go
func GetTrade(db *database.DB, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, _ := middleware.GetUserID(r)

        // Get ID from URL
        tradeID := chi.URLParam(r, "id")

        // Query database
        var trade models.Trade
        err := db.QueryRow(`
            SELECT id, symbol, trade_type, quantity, entry_price, pnl, opened_at
            FROM trades
            WHERE id = $1 AND user_id = $2
        `, tradeID, userID).Scan(
            &trade.ID, &trade.Symbol, &trade.TradeType,
            &trade.Quantity, &trade.EntryPrice, &trade.PnL, &trade.OpenedAt,
        )

        if err == sql.ErrNoRows {
            writeError(w, http.StatusNotFound, "NOT_FOUND", "Trade not found")
            return
        }
        if err != nil {
            logger.Error("Failed to get trade", "error", err)
            writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch trade")
            return
        }

        writeSuccess(w, http.StatusOK, trade)
    }
}
```

### PUT - Update Resource

```go
func UpdateTrade(db *database.DB, logger *slog.Logger, bus *notifications.Bus) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, _ := middleware.GetUserID(r)
        tradeID := chi.URLParam(r, "id")

        // Parse request
        var req struct {
            ExitPrice *float64 `json:"exit_price"`
            ClosedAt  *string  `json:"closed_at"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
            return
        }

        // Update database
        result, err := db.Exec(`
            UPDATE trades
            SET exit_price = COALESCE($1, exit_price),
                closed_at = COALESCE($2, closed_at)
            WHERE id = $3 AND user_id = $4
        `, req.ExitPrice, req.ClosedAt, tradeID, userID)

        if err != nil {
            logger.Error("Failed to update trade", "error", err)
            writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update trade")
            return
        }

        rows, _ := result.RowsAffected()
        if rows == 0 {
            writeError(w, http.StatusNotFound, "NOT_FOUND", "Trade not found")
            return
        }

        // Publish notification
        bus.Publish(
            notifications.NotificationTypeTradeUpdated,
            userID,
            "Trade Updated",
            "Trade updated successfully",
            map[string]interface{}{"trade_id": tradeID},
        )

        writeSuccess(w, http.StatusOK, map[string]interface{}{
            "message": "Trade updated successfully",
        })
    }
}
```

### DELETE - Remove Resource

```go
func DeleteTrade(db *database.DB, logger *slog.Logger, bus *notifications.Bus) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, _ := middleware.GetUserID(r)
        tradeID := chi.URLParam(r, "id")

        // Delete from database
        result, err := db.Exec(`
            DELETE FROM trades
            WHERE id = $1 AND user_id = $2
        `, tradeID, userID)

        if err != nil {
            logger.Error("Failed to delete trade", "error", err)
            writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to delete trade")
            return
        }

        rows, _ := result.RowsAffected()
        if rows == 0 {
            writeError(w, http.StatusNotFound, "NOT_FOUND", "Trade not found")
            return
        }

        // Publish notification
        bus.Publish(
            notifications.NotificationTypeTradeDeleted,
            userID,
            "Trade Deleted",
            "Trade deleted successfully",
            nil,
        )

        writeSuccess(w, http.StatusOK, map[string]interface{}{
            "message": "Trade deleted successfully",
        })
    }
}
```

## File Upload Pattern

```go
func UploadFile(db *database.DB, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, _ := middleware.GetUserID(r)

        // Limit upload size
        r.ParseMultipartForm(10 << 20) // 10MB

        // Get file from form
        file, header, err := r.FormFile("file")
        if err != nil {
            writeError(w, http.StatusBadRequest, "NO_FILE", "No file provided")
            return
        }
        defer file.Close()

        // Validate file type
        if !isValidFileType(header.Filename) {
            writeError(w, http.StatusBadRequest, "INVALID_FILE_TYPE", "File type not allowed")
            return
        }

        // Save file
        filepath := fmt.Sprintf("./uploads/%s-%s", uuid.New(), header.Filename)
        dst, err := os.Create(filepath)
        if err != nil {
            writeError(w, http.StatusInternalServerError, "FILE_ERROR", "Failed to save file")
            return
        }
        defer dst.Close()

        if _, err = io.Copy(dst, file); err != nil {
            writeError(w, http.StatusInternalServerError, "FILE_ERROR", "Failed to save file")
            return
        }

        writeSuccess(w, http.StatusCreated, map[string]interface{}{
            "filename": header.Filename,
            "path":     filepath,
        })
    }
}
```

## Error Codes

Standard error codes:

| Code | HTTP Status | Usage |
|------|-------------|-------|
| `INVALID_JSON` | 400 | Request body is not valid JSON |
| `VALIDATION_ERROR` | 400 | Input validation failed |
| `UNAUTHORIZED` | 401 | Missing or invalid auth token |
| `FORBIDDEN` | 403 | User lacks permission |
| `NOT_FOUND` | 404 | Resource doesn't exist |
| `DUPLICATE_ENTRY` | 409 | Resource already exists |
| `FILE_TOO_LARGE` | 413 | File exceeds size limit |
| `INVALID_FILE_TYPE` | 415 | File type not allowed |
| `DATABASE_ERROR` | 500 | Database operation failed |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

## Best Practices

1. **Always filter by user_id** - Prevent unauthorized access
2. **Validate all input** - Never trust user data
3. **Use specific error codes** - Help frontend handle errors
4. **Log errors with context** - Include user_id, resource_id
5. **Publish notifications** - Keep users informed
6. **Use transactions** - For multi-step operations
7. **Return appropriate status codes** - 200, 201, 400, 404, 500
8. **Handle edge cases** - Empty results, missing parameters
