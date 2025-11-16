# Trades API Documentation

Complete documentation for the TradePulse Trades API endpoints.

**Base URL:** `https://api.tradepulse.drivenw.com:9000/api`

**Authentication:** All endpoints require JWT token in `Authorization: Bearer <token>` header

---

## Table of Contents

- [List Trades](#list-trades)
- [Get Single Trade](#get-single-trade)
- [Create Trade](#create-trade)
- [Update Trade](#update-trade)
- [Delete Trade](#delete-trade)
- [Add Tag to Trade](#add-tag-to-trade)
- [Remove Tag from Trade](#remove-tag-from-trade)
- [Bulk Import Trades](#bulk-import-trades)

---

## List Trades

Retrieve all trades for the authenticated user with optional filtering and pagination.

### Request

```http
GET /api/trades?symbol=AAPL&trade_type=LONG&status=closed&start_date=2024-01-01&end_date=2024-12-31&limit=50&offset=0
```

### Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `symbol` | string | No | Filter by stock symbol (case-insensitive) |
| `trade_type` | string | No | Filter by trade type: `LONG` or `SHORT` |
| `status` | string | No | Filter by status: `all`, `open`, `closed` |
| `start_date` | string | No | Filter trades opened on or after this date (ISO 8601) |
| `end_date` | string | No | Filter trades opened on or before this date (ISO 8601) |
| `limit` | integer | No | Maximum number of results to return |
| `offset` | integer | No | Number of results to skip (for pagination) |

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "symbol": "AAPL",
      "trade_type": "LONG",
      "quantity": 100,
      "entry_price": 150.25,
      "exit_price": 155.50,
      "fees": 2.50,
      "pnl": 522.50,
      "opened_at": "2024-01-15T09:30:00Z",
      "closed_at": "2024-01-15T15:45:00Z",
      "created_at": "2024-01-15T09:30:05Z",
      "updated_at": "2024-01-15T15:45:10Z",
      "has_journal": true,
      "tags": ["breakout", "morning-trade"]
    }
  ]
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Unique trade identifier |
| `user_id` | UUID | Owner's user ID |
| `symbol` | string | Stock symbol |
| `trade_type` | string | `LONG` or `SHORT` |
| `quantity` | float | Number of shares |
| `entry_price` | float | Average entry price |
| `exit_price` | float | Average exit price (null if still open) |
| `fees` | float | Total fees/commissions |
| `pnl` | float | Profit/loss (calculated automatically, null if open) |
| `opened_at` | timestamp | When position was opened |
| `closed_at` | timestamp | When position was closed (null if still open) |
| `created_at` | timestamp | When record was created |
| `updated_at` | timestamp | When record was last updated |
| `has_journal` | boolean | Whether a journal entry exists for this trade |
| `tags` | array | Array of associated tag names |

### Status Codes

- `200 OK` - Success
- `401 Unauthorized` - Missing or invalid JWT token
- `500 Internal Server Error` - Server error

### Example

```bash
curl -X GET "https://api.tradepulse.drivenw.com:9000/api/trades?status=closed&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Get Single Trade

Retrieve a specific trade by ID.

### Request

```http
GET /api/trades/{id}
```

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | UUID | Yes | Trade ID |

### Response

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "symbol": "TSLA",
    "trade_type": "SHORT",
    "quantity": 50,
    "entry_price": 250.75,
    "exit_price": 245.20,
    "fees": 1.50,
    "pnl": 276.00,
    "opened_at": "2024-01-20T10:15:00Z",
    "closed_at": "2024-01-20T14:30:00Z",
    "created_at": "2024-01-20T10:15:05Z",
    "updated_at": "2024-01-20T14:30:10Z",
    "has_journal": false,
    "tags": ["gap-fade", "afternoon"]
  }
}
```

### Status Codes

- `200 OK` - Success
- `400 Bad Request` - Invalid trade ID format
- `401 Unauthorized` - Missing or invalid JWT token
- `404 Not Found` - Trade not found or not owned by user
- `500 Internal Server Error` - Server error

### Example

```bash
curl -X GET "https://api.tradepulse.drivenw.com:9000/api/trades/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Create Trade

Create a new trade.

### Request

```http
POST /api/trades
Content-Type: application/json
```

### Request Body

```json
{
  "symbol": "NVDA",
  "trade_type": "LONG",
  "quantity": 25,
  "entry_price": 475.50,
  "exit_price": 482.25,
  "fees": 1.00,
  "opened_at": "2024-02-01T09:45:00Z",
  "closed_at": "2024-02-01T11:30:00Z"
}
```

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `symbol` | string | Stock symbol (will be uppercased) |
| `trade_type` | string | Must be `LONG` or `SHORT` |
| `quantity` | float | Number of shares (must be > 0) |
| `entry_price` | float | Average entry price (must be > 0) |
| `opened_at` | timestamp | When position was opened (ISO 8601) |

### Optional Fields

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `exit_price` | float | Average exit price (null for open position) | null |
| `fees` | float | Total fees/commissions | 0.0 |
| `closed_at` | timestamp | When position was closed | null |

### Response

```json
{
  "success": true,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "symbol": "NVDA",
    "trade_type": "LONG",
    "quantity": 25,
    "entry_price": 475.50,
    "exit_price": 482.25,
    "fees": 1.00,
    "pnl": 167.75,
    "opened_at": "2024-02-01T09:45:00Z",
    "closed_at": "2024-02-01T11:30:00Z",
    "created_at": "2024-02-01T11:30:15Z",
    "updated_at": "2024-02-01T11:30:15Z",
    "has_journal": false,
    "tags": []
  }
}
```

### P&L Calculation

P&L is calculated automatically by the database:

**For LONG trades:**
```
PnL = (exit_price - entry_price) × quantity - fees
```

**For SHORT trades:**
```
PnL = (entry_price - exit_price) × quantity - fees
```

If `exit_price` is null (open position), `pnl` will be null.

### Status Codes

- `201 Created` - Trade created successfully
- `400 Bad Request` - Invalid request body or missing required fields
- `401 Unauthorized` - Missing or invalid JWT token
- `500 Internal Server Error` - Server error

### Example

```bash
curl -X POST "https://api.tradepulse.drivenw.com:9000/api/trades" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "NVDA",
    "trade_type": "LONG",
    "quantity": 25,
    "entry_price": 475.50,
    "exit_price": 482.25,
    "fees": 1.00,
    "opened_at": "2024-02-01T09:45:00Z",
    "closed_at": "2024-02-01T11:30:00Z"
  }'
```

---

## Update Trade

Update an existing trade.

### Request

```http
PUT /api/trades/{id}
Content-Type: application/json
```

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | UUID | Yes | Trade ID to update |

### Request Body

Same as [Create Trade](#create-trade), all fields are required.

```json
{
  "symbol": "NVDA",
  "trade_type": "LONG",
  "quantity": 25,
  "entry_price": 475.50,
  "exit_price": 485.00,
  "fees": 1.00,
  "opened_at": "2024-02-01T09:45:00Z",
  "closed_at": "2024-02-01T13:00:00Z"
}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "symbol": "NVDA",
    "trade_type": "LONG",
    "quantity": 25,
    "entry_price": 475.50,
    "exit_price": 485.00,
    "fees": 1.00,
    "pnl": 236.50,
    "opened_at": "2024-02-01T09:45:00Z",
    "closed_at": "2024-02-01T13:00:00Z",
    "created_at": "2024-02-01T11:30:15Z",
    "updated_at": "2024-02-01T13:00:20Z",
    "has_journal": false,
    "tags": []
  }
}
```

### Notes

- P&L is recalculated automatically when updating
- `updated_at` timestamp is automatically updated
- Cannot update trades belonging to other users

### Status Codes

- `200 OK` - Trade updated successfully
- `400 Bad Request` - Invalid request body or trade ID
- `401 Unauthorized` - Missing or invalid JWT token
- `404 Not Found` - Trade not found or not owned by user
- `500 Internal Server Error` - Server error

### Example

```bash
curl -X PUT "https://api.tradepulse.drivenw.com:9000/api/trades/660e8400-e29b-41d4-a716-446655440001" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "NVDA",
    "trade_type": "LONG",
    "quantity": 25,
    "entry_price": 475.50,
    "exit_price": 485.00,
    "fees": 1.00,
    "opened_at": "2024-02-01T09:45:00Z",
    "closed_at": "2024-02-01T13:00:00Z"
  }'
```

---

## Delete Trade

Delete a trade permanently.

### Request

```http
DELETE /api/trades/{id}
```

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | UUID | Yes | Trade ID to delete |

### Response

```json
{
  "success": true,
  "message": "Trade deleted successfully"
}
```

### Notes

- Deletion is permanent and cannot be undone
- Associated journal entries are NOT automatically deleted (separate deletion required)
- Tag associations are automatically removed (CASCADE delete)
- Cannot delete trades belonging to other users

### Status Codes

- `200 OK` - Trade deleted successfully
- `400 Bad Request` - Invalid trade ID format
- `401 Unauthorized` - Missing or invalid JWT token
- `404 Not Found` - Trade not found or not owned by user
- `500 Internal Server Error` - Server error

### Example

```bash
curl -X DELETE "https://api.tradepulse.drivenw.com:9000/api/trades/660e8400-e29b-41d4-a716-446655440001" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Add Tag to Trade

Associate an existing tag with a trade.

### Request

```http
POST /api/trades/{id}/tags
Content-Type: application/json
```

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | UUID | Yes | Trade ID |

### Request Body

```json
{
  "tag_id": "770e8400-e29b-41d4-a716-446655440002"
}
```

### Response

```json
{
  "success": true,
  "message": "Tag added to trade successfully"
}
```

### Notes

- Tag must already exist (create tags via `/api/tags` endpoint)
- Tag must belong to the same user
- Duplicate tag associations are ignored (idempotent)
- Trade ownership is verified

### Status Codes

- `200 OK` - Tag added successfully (or already exists)
- `400 Bad Request` - Invalid trade ID or tag ID
- `401 Unauthorized` - Missing or invalid JWT token
- `404 Not Found` - Trade not found or not owned by user
- `500 Internal Server Error` - Server error

### Example

```bash
curl -X POST "https://api.tradepulse.drivenw.com:9000/api/trades/660e8400-e29b-41d4-a716-446655440001/tags" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "tag_id": "770e8400-e29b-41d4-a716-446655440002"
  }'
```

---

## Remove Tag from Trade

Remove a tag association from a trade.

### Request

```http
DELETE /api/trades/{tradeId}/tags/{tagId}
```

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tradeId` | UUID | Yes | Trade ID |
| `tagId` | UUID | Yes | Tag ID to remove |

### Response

```json
{
  "success": true,
  "message": "Tag removed from trade successfully"
}
```

### Notes

- Trade ownership is verified
- Returns error if tag association doesn't exist

### Status Codes

- `200 OK` - Tag removed successfully
- `400 Bad Request` - Invalid trade ID or tag ID
- `401 Unauthorized` - Missing or invalid JWT token
- `404 Not Found` - Trade not found, tag association not found, or not owned by user
- `500 Internal Server Error` - Server error

### Example

```bash
curl -X DELETE "https://api.tradepulse.drivenw.com:9000/api/trades/660e8400-e29b-41d4-a716-446655440001/tags/770e8400-e29b-41d4-a716-446655440002" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Bulk Import Trades

Import multiple trades at once (typically from CSV).

### Status

⚠️ **NOT YET IMPLEMENTED** - Handler is stubbed, database function ready

### Planned Request

```http
POST /api/trades/import-csv
Content-Type: application/json
```

### Planned Request Body

```json
{
  "trades": [
    {
      "symbol": "AAPL",
      "trade_type": "LONG",
      "quantity": 100,
      "entry_price": 150.00,
      "exit_price": 155.00,
      "fees": 2.00,
      "opened_at": "2024-01-15T09:30:00Z",
      "closed_at": "2024-01-15T15:00:00Z"
    },
    {
      "symbol": "TSLA",
      "trade_type": "SHORT",
      "quantity": 50,
      "entry_price": 250.00,
      "exit_price": 245.00,
      "fees": 1.50,
      "opened_at": "2024-01-15T10:00:00Z",
      "closed_at": "2024-01-15T14:00:00Z"
    }
  ]
}
```

### Planned Response

```json
{
  "success": true,
  "data": {
    "imported_count": 2,
    "trade_ids": [
      "880e8400-e29b-41d4-a716-446655440003",
      "990e8400-e29b-41d4-a716-446655440004"
    ]
  }
}
```

### Implementation Notes

- Backend database function `BulkCreateTrades()` is implemented
- Uses database transaction for atomicity (all or nothing)
- Frontend CSV parsing is complete (DAS Trader, PropReports)
- HTTP handler needs to be created
- Will send WebSocket notification on completion

---

## WebSocket Notifications

All trade operations trigger real-time WebSocket notifications.

### Notification Types

| Operation | Notification Type | Data |
|-----------|------------------|------|
| Create Trade | `trade.created` | `{ id, symbol }` |
| Update Trade | `trade.updated` | `{ id, symbol }` |
| Delete Trade | `trade.deleted` | `{ id }` |

### Example Notification

```json
{
  "id": "notif-123",
  "type": "trade.created",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Trade Created",
  "message": "New trade added successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "symbol": "NVDA"
  },
  "timestamp": "2024-02-01T11:30:15Z"
}
```

### Connecting to WebSocket

```javascript
const ws = new WebSocket('wss://api.tradepulse.drivenw.com:9000/api/ws');

ws.onmessage = (event) => {
  const notification = JSON.parse(event.data);
  console.log('Notification:', notification);
};
```

See [WebSocket Notifications Documentation](../websocket-notifications.md) for details.

---

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": {
    "message": "Descriptive error message"
  }
}
```

### Common Error Codes

| Status Code | Description |
|-------------|-------------|
| `400 Bad Request` | Invalid input, missing required fields, or malformed data |
| `401 Unauthorized` | Missing, invalid, or expired JWT token |
| `404 Not Found` | Resource not found or not owned by user |
| `500 Internal Server Error` | Server-side error occurred |

---

## Rate Limiting

⚠️ **NOT YET IMPLEMENTED**

Planned rate limiting:
- 100 requests per minute per user
- 1000 requests per hour per user
- Bulk import has separate limits

---

## Best Practices

### Creating Trades

1. **Always validate before submitting**
   - Ensure `quantity > 0`
   - Ensure `entry_price > 0`
   - Ensure `trade_type` is `LONG` or `SHORT`

2. **Use proper timestamps**
   - Use ISO 8601 format
   - Include timezone information
   - Ensure `opened_at <= closed_at` for closed positions

3. **Set exit price only for closed positions**
   - Leave `exit_price` and `closed_at` null for open positions
   - Update later when closing the position

### Filtering & Pagination

1. **Use pagination for large datasets**
   - Default to `limit=50` for reasonable performance
   - Use `offset` for subsequent pages

2. **Combine filters for targeted queries**
   - Example: `?symbol=AAPL&status=closed&start_date=2024-01-01`

3. **Cache results when appropriate**
   - Closed trades don't change frequently
   - Consider client-side caching

### Error Handling

1. **Always check the `success` field**
2. **Handle `401 Unauthorized` by refreshing JWT**
3. **Display error messages from API to users**
4. **Retry on `500` errors with exponential backoff**

---

## Code Examples

### JavaScript/TypeScript

```typescript
const API_BASE = 'https://api.tradepulse.drivenw.com:9000/api';

async function createTrade(token: string, trade: Trade) {
  const response = await fetch(`${API_BASE}/trades`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(trade)
  });

  const result = await response.json();

  if (!result.success) {
    throw new Error(result.error.message);
  }

  return result.data;
}

async function listTrades(token: string, filters: TradeFilters = {}) {
  const params = new URLSearchParams(filters);
  const response = await fetch(`${API_BASE}/trades?${params}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });

  const result = await response.json();

  if (!result.success) {
    throw new Error(result.error.message);
  }

  return result.data;
}
```

### Python

```python
import requests

API_BASE = 'https://api.tradepulse.drivenw.com:9000/api'

def create_trade(token: str, trade: dict):
    response = requests.post(
        f'{API_BASE}/trades',
        headers={
            'Authorization': f'Bearer {token}',
            'Content-Type': 'application/json'
        },
        json=trade
    )

    result = response.json()

    if not result['success']:
        raise Exception(result['error']['message'])

    return result['data']

def list_trades(token: str, filters: dict = {}):
    response = requests.get(
        f'{API_BASE}/trades',
        headers={'Authorization': f'Bearer {token}'},
        params=filters
    )

    result = response.json()

    if not result['success']:
        raise Exception(result['error']['message'])

    return result['data']
```

---

## Changelog

### 2024-11-15
- Initial documentation of implemented trades API
- All CRUD endpoints functional
- Tag management implemented
- Bulk import database function ready (handler pending)
