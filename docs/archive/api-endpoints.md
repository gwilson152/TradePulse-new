# Backend API Endpoints Documentation

**Last Updated:** January 2025

This document provides comprehensive documentation for all API endpoints in the TradePulse backend.

---

## Authentication Endpoints

### POST /api/auth/request-magic-link
Request a magic link for passwordless authentication.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Magic link sent to your email"
  }
}
```

---

### POST /api/auth/signup
Create a new user account with a plan selection.

**Request Body:**
```json
{
  "email": "user@example.com",
  "plan_type": "starter"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": { ... },
    "message": "Account created. Check your email for verification link."
  }
}
```

---

### GET /api/auth/verify
Verify a magic link token and return JWT.

**Query Parameters:**
- `token` - The magic link token

**Response:**
```json
{
  "success": true,
  "data": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": { ... }
  }
}
```

---

### POST /api/auth/login
Login with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": { ... }
  }
}
```

---

### GET /api/auth/me
Get current authenticated user.

**Headers:**
- `Authorization: Bearer <jwt>`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "company": "Acme Corp",
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "United States",
    "timezone": "America/New_York",
    "profile_completed": true,
    "plan_type": "starter",
    "plan_status": "beta_free",
    "created_at": "2025-01-01T00:00:00Z",
    "last_login": "2025-01-15T12:00:00Z"
  }
}
```

---

### POST /api/auth/set-password
Set or update user password.

**Headers:**
- `Authorization: Bearer <jwt>`

**Request Body:**
```json
{
  "password": "newsecurepassword"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Password set successfully"
  }
}
```

---

### POST /api/auth/logout
Logout current user.

**Headers:**
- `Authorization: Bearer <jwt>`

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Logged out successfully"
  }
}
```

---

### POST /api/auth/refresh
Refresh JWT token.

**Headers:**
- `Authorization: Bearer <jwt>`

**Response:**
```json
{
  "success": true,
  "data": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

## User Profile Endpoints

### PUT /api/users/profile
Update user profile information.

**Headers:**
- `Authorization: Bearer <jwt>`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "company": "Acme Corp",
  "address_line1": "123 Main St",
  "address_line2": "Apt 4B",
  "city": "New York",
  "state": "NY",
  "postal_code": "10001",
  "country": "United States",
  "timezone": "America/New_York",
  "profile_completed": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    ...
  }
}
```

**Notes:**
- All fields are optional (use `null` to clear a field)
- Profile fields are stored as nullable strings in database
- `profile_completed` is used to track onboarding status

---

## Account Management Endpoints

### POST /api/account/reset-data
Delete all user data (trades, journal entries, rule sets, etc.) while keeping the account.

**Headers:**
- `Authorization: Bearer <jwt>`

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "All account data has been reset successfully"
  }
}
```

**Notes:**
- Deletes data in transaction-safe manner
- Respects foreign key constraints
- Does NOT delete user account or authentication data
- Deletion order:
  1. Journal attachments
  2. Journal entries
  3. Trade exits
  4. Trade entries
  5. Trade tags associations
  6. Trades
  7. Rules
  8. Rule sets
  9. Tags

---

## Trade Endpoints

### GET /api/trades
List all trades for authenticated user.

**Headers:**
- `Authorization: Bearer <jwt>`

**Query Parameters:**
- `limit` - Number of results (default: 50)
- `offset` - Offset for pagination (default: 0)
- `symbol` - Filter by symbol
- `trade_type` - Filter by LONG or SHORT
- `status` - Filter by open or closed
- `start_date` - Filter by date range
- `end_date` - Filter by date range

**Response:**
```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "total": 100,
    "limit": 50,
    "offset": 0
  }
}
```

---

### POST /api/trades
Create a new trade.

**Headers:**
- `Authorization: Bearer <jwt>`

**Request Body:**
```json
{
  "symbol": "AAPL",
  "trade_type": "LONG",
  "opened_at": "2025-01-15T09:30:00Z"
}
```

---

### GET /api/trades/{id}
Get a specific trade by ID.

---

### PUT /api/trades/{id}
Update a trade.

---

### DELETE /api/trades/{id}
Delete a trade.

---

## Trade Entry/Exit Endpoints

### POST /api/trades/{id}/entries
Add an entry to a trade.

**Request Body:**
```json
{
  "price": 150.25,
  "quantity": 100,
  "timestamp": "2025-01-15T09:30:00Z",
  "fees": 1.50,
  "notes": "Entry at support level"
}
```

---

### GET /api/trades/{id}/entries
Get all entries for a trade.

---

### DELETE /api/trades/{id}/entries/{entryId}
Delete an entry from a trade.

---

### POST /api/trades/{id}/exits
Add an exit to a trade.

**Request Body:**
```json
{
  "price": 155.75,
  "quantity": 50,
  "timestamp": "2025-01-15T11:30:00Z",
  "fees": 1.25,
  "notes": "Partial exit at resistance"
}
```

**Notes:**
- Backend automatically calculates P&L for the exit
- Position metrics are recalculated
- Trade auto-closes when position size reaches 0

---

### GET /api/trades/{id}/exits
Get all exits for a trade.

---

### DELETE /api/trades/{id}/exits/{exitId}
Delete an exit from a trade.

---

## Journal Endpoints

### GET /api/journal
List all journal entries.

### POST /api/journal
Create a journal entry.

### GET /api/journal/{id}
Get a specific journal entry.

### PUT /api/journal/{id}
Update a journal entry.

### DELETE /api/journal/{id}
Delete a journal entry.

### GET /api/trades/{tradeId}/journal
Get all journal entries for a specific trade.

---

## Rule Set Endpoints

### GET /api/rulesets
List all rule sets.

### POST /api/rulesets
Create a rule set.

### GET /api/rulesets/{id}
Get a specific rule set.

### PUT /api/rulesets/{id}
Update a rule set.

### DELETE /api/rulesets/{id}
Delete a rule set.

### POST /api/rulesets/{ruleSetId}/rules
Add a rule to a rule set.

### PUT /api/rulesets/{ruleSetId}/rules/{ruleId}
Update a rule.

### DELETE /api/rulesets/{ruleSetId}/rules/{ruleId}
Delete a rule.

---

## Metrics Endpoints

### GET /api/metrics/summary
Get summary metrics for all trades.

### GET /api/metrics/by-symbol
Get metrics grouped by symbol.

### GET /api/metrics/daily
Get daily performance metrics.

---

## Integration Endpoints

### POST /api/integrations/propreports/fetch
Fetch trades from PropReports API.

---

## WebSocket Endpoints

### GET /api/ws
WebSocket connection for real-time notifications.

### GET /api/notifications/stats
Get notification bus statistics.

---

## Database Schema Changes

### Migration 006: User Profile Fields

**Added Columns to `users` table:**
- `first_name` VARCHAR(100)
- `last_name` VARCHAR(100)
- `phone` VARCHAR(50)
- `company` VARCHAR(255)
- `address_line1` VARCHAR(255)
- `address_line2` VARCHAR(255)
- `city` VARCHAR(100)
- `state` VARCHAR(100)
- `postal_code` VARCHAR(20)
- `country` VARCHAR(100)
- `timezone` VARCHAR(100) DEFAULT 'America/New_York'
- `profile_completed` BOOLEAN DEFAULT FALSE

**Indexes:**
- `idx_users_profile_completed` on `profile_completed`

---

## Error Responses

All endpoints return errors in the following format:

```json
{
  "success": false,
  "error": {
    "message": "Error description",
    "code": "ERROR_CODE"
  }
}
```

**Common HTTP Status Codes:**
- `200` - Success
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid JWT)
- `404` - Not Found
- `500` - Internal Server Error

---

## Authentication

Most endpoints require authentication via JWT token in the Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Exceptions (Public Endpoints):**
- POST /api/auth/request-magic-link
- POST /api/auth/signup
- GET /api/auth/verify
- POST /api/auth/login
- GET /health

---

## CORS Configuration

**Allowed Origins:**
- https://tradepulse.drivenw.com

**Allowed Methods:**
- GET, POST, PUT, DELETE, OPTIONS

**Allowed Headers:**
- Accept, Authorization, Content-Type

---

## Rate Limiting

Currently no rate limiting is implemented. Future versions may add:
- Per-user rate limits
- Per-IP rate limits
- WebSocket connection limits
