# TradePulse API Specification

**Base URL (Development)**: `https://api.tradepulse.drivenw.com`

**Note**: External proxy routes HTTPS (443) to internal port 9000. No port specification needed in API calls.

All authenticated endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

## Authentication

### Request Magic Link

**Endpoint:** `POST /api/auth/request-magic-link`

**Description:** Send a magic link to the user's email address for passwordless authentication.

**Request:**
```json
{
  "email": "trader@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Magic link sent to trader@example.com",
    "expires_in": 900
  }
}
```

**Status Codes:**
- `200`: Success
- `400`: Invalid email format
- `429`: Too many requests

---

### Verify Magic Link

**Endpoint:** `GET /api/auth/verify`

**Description:** Verify magic link token and receive JWT for authenticated requests.

**Query Parameters:**
- `token` (required): Magic link token from email

**Request:**
```
GET /api/auth/verify?token=abc123...
```

**Response:**
```json
{
  "success": true,
  "data": {
    "jwt": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "trader@example.com",
      "created_at": "2024-01-15T10:30:00Z",
      "last_login": "2024-01-20T14:22:00Z"
    }
  }
}
```

**Status Codes:**
- `200`: Success
- `400`: Invalid or missing token
- `401`: Token expired or already used

---

### Logout

**Endpoint:** `POST /api/auth/logout`

**Authentication:** Required

**Description:** Invalidate current JWT token (client should also clear local storage).

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

### Refresh Token

**Endpoint:** `POST /api/auth/refresh`

**Authentication:** Required (with existing JWT)

**Description:** Get a new JWT token before the current one expires.

**Response:**
```json
{
  "success": true,
  "data": {
    "jwt": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

---

## Trades

### List Trades

**Endpoint:** `GET /api/trades`

**Authentication:** Required

**Description:** Get paginated list of user's trades with optional filtering.

**Query Parameters:**
- `limit` (optional, default: 50): Number of trades to return
- `offset` (optional, default: 0): Number of trades to skip
- `from` (optional): Start date (ISO 8601)
- `to` (optional): End date (ISO 8601)
- `symbol` (optional): Filter by symbol
- `type` (optional): Filter by LONG or SHORT
- `sort` (optional, default: opened_at): Sort field
- `order` (optional, default: desc): asc or desc

**Request:**
```
GET /api/trades?limit=20&from=2024-01-01&symbol=AAPL&sort=pnl&order=desc
```

**Response:**
```json
{
  "success": true,
  "data": {
    "trades": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "symbol": "AAPL",
        "trade_type": "LONG",
        "quantity": 100,
        "entry_price": 150.25,
        "exit_price": 155.80,
        "fees": 2.50,
        "pnl": 552.50,
        "opened_at": "2024-01-15T09:30:00Z",
        "closed_at": "2024-01-15T15:45:00Z",
        "has_journal": true,
        "tags": ["breakout", "tech"]
      }
    ],
    "total": 145,
    "limit": 20,
    "offset": 0
  }
}
```

---

### Get Single Trade

**Endpoint:** `GET /api/trades/:id`

**Authentication:** Required

**Description:** Get detailed information about a specific trade.

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "660e8400-e29b-41d4-a716-446655440001",
    "symbol": "AAPL",
    "trade_type": "LONG",
    "entries": [
      {
        "id": "entry-1",
        "price": 150.25,
        "quantity": 100,
        "timestamp": "2024-01-15T09:30:00Z",
        "notes": "Initial entry on breakout",
        "fees": 1.25
      }
    ],
    "exits": [
      {
        "id": "exit-1",
        "price": 155.80,
        "quantity": 100,
        "timestamp": "2024-01-15T15:45:00Z",
        "notes": "Exiting at target",
        "fees": 1.25,
        "pnl": 552.50
      }
    ],
    "current_position_size": 0,
    "average_entry_price": 150.25,
    "total_fees": 2.50,
    "realized_pnl": 552.50,
    "unrealized_pnl": null,
    "cost_basis_method": "FIFO",
    "opened_at": "2024-01-15T09:30:00Z",
    "closed_at": "2024-01-15T15:45:00Z",
    "created_at": "2024-01-15T16:00:00Z",
    "updated_at": "2024-01-15T16:00:00Z",
    "notes": "Strong breakout pattern",
    "tags": ["breakout", "tech"]
  }
}
```

**Status Codes:**
- `200`: Success
- `404`: Trade not found

---

### Create Trade

**Endpoint:** `POST /api/trades`

**Authentication:** Required

**Description:** Manually create a new trade entry.

**Request:**
```json
{
  "symbol": "TSLA",
  "trade_type": "SHORT",
  "quantity": 50,
  "entry_price": 245.50,
  "exit_price": 240.20,
  "fees": 1.75,
  "opened_at": "2024-01-20T10:15:00Z",
  "closed_at": "2024-01-20T14:30:00Z"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "symbol": "TSLA",
    "trade_type": "SHORT",
    "quantity": 50,
    "entry_price": 245.50,
    "exit_price": 240.20,
    "fees": 1.75,
    "pnl": 263.25,
    "opened_at": "2024-01-20T10:15:00Z",
    "closed_at": "2024-01-20T14:30:00Z",
    "created_at": "2024-01-20T15:00:00Z",
    "updated_at": "2024-01-20T15:00:00Z"
  }
}
```

**Status Codes:**
- `201`: Created
- `400`: Validation error

---

### Update Trade

**Endpoint:** `PUT /api/trades/:id`

**Authentication:** Required

**Description:** Update an existing trade.

**Request:**
```json
{
  "exit_price": 241.00,
  "closed_at": "2024-01-20T15:00:00Z"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "pnl": 223.25,
    "updated_at": "2024-01-20T15:05:00Z"
  }
}
```

---

### Delete Trade

**Endpoint:** `DELETE /api/trades/:id`

**Authentication:** Required

**Description:** Delete a trade and associated journal entries.

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Trade deleted successfully"
  }
}
```

**Status Codes:**
- `200`: Success
- `404`: Trade not found

---

### Add Entry to Position

**Endpoint:** `POST /api/trades/:id/entries`

**Authentication:** Required

**Description:** Add a new entry (scale in) to an existing position.

**Request:**
```json
{
  "price": 152.50,
  "quantity": 50,
  "timestamp": "2024-01-16T10:30:00Z",
  "notes": "Adding to position on pullback",
  "fees": 1.25
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "entry-550e8400-e29b-41d4-a716-446655440000",
    "trade_id": "550e8400-e29b-41d4-a716-446655440000",
    "price": 152.50,
    "quantity": 50,
    "timestamp": "2024-01-16T10:30:00Z",
    "notes": "Adding to position on pullback",
    "fees": 1.25,
    "current_position_size": 150,
    "average_entry_price": 151.17
  }
}
```

**Status Codes:**
- `201`: Created
- `400`: Validation error
- `404`: Trade not found

---

### Add Exit to Position

**Endpoint:** `POST /api/trades/:id/exits`

**Authentication:** Required

**Description:** Add a new exit (scale out) to an existing position.

**Request:**
```json
{
  "price": 158.75,
  "quantity": 75,
  "timestamp": "2024-01-17T14:00:00Z",
  "notes": "Taking partial profit at resistance",
  "fees": 1.50
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "exit-660e8400-e29b-41d4-a716-446655440001",
    "trade_id": "550e8400-e29b-41d4-a716-446655440000",
    "price": 158.75,
    "quantity": 75,
    "timestamp": "2024-01-17T14:00:00Z",
    "notes": "Taking partial profit at resistance",
    "fees": 1.50,
    "pnl": 565.50,
    "current_position_size": 75,
    "realized_pnl": 565.50
  }
}
```

**Status Codes:**
- `201`: Created
- `400`: Validation error (e.g., exiting more than position size)
- `404`: Trade not found

---

### Import Trades from CSV

**Endpoint:** `POST /api/trades/import-csv`

**Authentication:** Required

**Content-Type:** `multipart/form-data`

**Description:** Bulk import trades from a CSV file.

**Request:**
FormData with:
- `file`: CSV file
- `mapping` (optional): JSON string with column mappings

**CSV Expected Columns:**
- date, symbol, type (LONG/SHORT), quantity, entry_price, exit_price, fees

**Response:**
```json
{
  "success": true,
  "data": {
    "imported": 45,
    "failed": 2,
    "errors": [
      {
        "row": 12,
        "error": "Invalid date format"
      },
      {
        "row": 28,
        "error": "Missing required field: symbol"
      }
    ]
  }
}
```

**Status Codes:**
- `200`: Success (even with partial failures)
- `400`: Invalid file or format

---

## Journal Entries

### List Journal Entries

**Endpoint:** `GET /api/journal`

**Authentication:** Required

**Description:** Get all journal entries for the authenticated user.

**Query Parameters:**
- `limit` (optional, default: 50)
- `offset` (optional, default: 0)
- `trade_id` (optional): Filter by specific trade
- `from` (optional): Start date
- `to` (optional): End date
- `search` (optional): Full-text search in content

**Response:**
```json
{
  "success": true,
  "data": {
    "entries": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "trade_id": "550e8400-e29b-41d4-a716-446655440000",
        "entry_date": "2024-01-15T16:00:00Z",
        "content": "Perfect breakout setup on AAPL...",
        "emotional_state": {
          "confidence": 8,
          "stress": 3,
          "discipline": 9,
          "notes": "Felt very confident about this setup"
        },
        "rule_adherence": [
          {
            "rule_id": "rule-770e8400-e29b-41d4-a716-446655440002",
            "rule_title": "Risk no more than 2% per trade",
            "score": 100,
            "notes": "Position sized correctly at 1.8% risk",
            "timestamp": "2024-01-15T16:00:00Z"
          }
        ],
        "adherence_score": 87.5,
        "screenshots": [
          "https://api.tradepulse.drivenw.com/uploads/screenshots/abc123.png"
        ],
        "voice_notes": [
          "https://api.tradepulse.drivenw.com/uploads/voice/def456.mp3"
        ],
        "created_at": "2024-01-15T16:00:00Z",
        "updated_at": "2024-01-15T16:00:00Z"
      }
    ],
    "total": 89,
    "limit": 50,
    "offset": 0
  }
}
```

---

### Get Single Journal Entry

**Endpoint:** `GET /api/journal/:id`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "trade_id": "550e8400-e29b-41d4-a716-446655440000",
    "content": "Perfect breakout setup...",
    "emotional_state": {
      "pre_trade_confidence": 8,
      "pre_trade_clarity": 7,
      "post_trade_discipline": 9,
      "post_trade_emotion": "satisfied"
    },
    "attachments": [],
    "created_at": "2024-01-15T16:00:00Z",
    "updated_at": "2024-01-15T16:00:00Z"
  }
}
```

---

### Create Journal Entry

**Endpoint:** `POST /api/journal`

**Authentication:** Required

**Description:** Create a new journal entry with emotional state, rule adherence tracking, and media attachments.

**Content-Type:** `multipart/form-data` (when uploading files) or `application/json`

**Request (JSON):**
```json
{
  "trade_id": "550e8400-e29b-41d4-a716-446655440000",
  "entry_date": "2024-01-20T16:30:00Z",
  "content": "Excellent execution on this trade. Waited for confirmation...",
  "emotional_state": {
    "confidence": 8,
    "stress": 4,
    "discipline": 9,
    "notes": "Felt very disciplined today"
  },
  "rule_adherence": [
    {
      "rule_id": "rule-770e8400-e29b-41d4-a716-446655440002",
      "rule_title": "Risk no more than 2% per trade",
      "score": 100,
      "notes": "Position sized correctly",
      "timestamp": "2024-01-20T16:30:00Z"
    },
    {
      "rule_id": "rule-880e8400-e29b-41d4-a716-446655440003",
      "rule_title": "Wait for confirmation candle",
      "score": 75,
      "notes": "Entered slightly early but still had confirmation",
      "timestamp": "2024-01-20T16:30:00Z"
    }
  ]
}
```

**Request (with Files - multipart/form-data):**
FormData with:
- `data`: JSON string with entry data (as above)
- `screenshots`: File[] (optional)
- `voice_notes`: File[] (optional)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "trade_id": "550e8400-e29b-41d4-a716-446655440000",
    "entry_date": "2024-01-20T16:30:00Z",
    "content": "Excellent execution on this trade...",
    "emotional_state": {
      "confidence": 8,
      "stress": 4,
      "discipline": 9,
      "notes": "Felt very disciplined today"
    },
    "rule_adherence": [
      {
        "rule_id": "rule-770e8400-e29b-41d4-a716-446655440002",
        "rule_title": "Risk no more than 2% per trade",
        "score": 100,
        "notes": "Position sized correctly",
        "timestamp": "2024-01-20T16:30:00Z"
      }
    ],
    "adherence_score": 87.5,
    "screenshots": [
      "https://api.tradepulse.drivenw.com/uploads/screenshots/abc123.png"
    ],
    "voice_notes": [
      "https://api.tradepulse.drivenw.com/uploads/voice/def456.mp3"
    ],
    "created_at": "2024-01-20T16:30:00Z",
    "updated_at": "2024-01-20T16:30:00Z"
  }
}
```

**Adherence Score Calculation:**
The `adherence_score` is automatically calculated as a weighted average:
- Each rule's score (0, 25, 50, 75, or 100) is multiplied by its weight (1-5)
- Sum of weighted scores divided by sum of weights
- Example: (100×5 + 75×4) ÷ (5+4) = 87.5

---

### Update Journal Entry

**Endpoint:** `PUT /api/journal/:id`

**Authentication:** Required

**Request:**
```json
{
  "content": "Updated reflection after reviewing...",
  "emotional_state": {
    "pre_trade_confidence": 7,
    "pre_trade_clarity": 9,
    "post_trade_discipline": 8,
    "post_trade_emotion": "reflective"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "updated_at": "2024-01-20T17:00:00Z"
  }
}
```

---

### Delete Journal Entry

**Endpoint:** `DELETE /api/journal/:id`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Journal entry deleted successfully"
  }
}
```

---

## Attachments

### Upload Attachment

**Endpoint:** `POST /api/journal/:id/attachments`

**Authentication:** Required

**Content-Type:** `multipart/form-data`

**Description:** Upload a screenshot or voice note to a journal entry.

**Request:**
FormData with:
- `file`: File to upload
- `type`: "screenshot" or "voice"

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "aa0e8400-e29b-41d4-a716-446655440005",
    "entry_id": "990e8400-e29b-41d4-a716-446655440004",
    "type": "screenshot",
    "filename": "chart_analysis.png",
    "file_size": 524288,
    "mime_type": "image/png",
    "url": "/api/attachments/aa0e8400-e29b-41d4-a716-446655440005",
    "uploaded_at": "2024-01-20T16:45:00Z"
  }
}
```

**File Constraints:**
- Max size: 10MB
- Allowed image types: PNG, JPG, JPEG, GIF, WebP
- Allowed audio types: MP3, WAV, M4A, OGG

---

### Get Attachment

**Endpoint:** `GET /api/attachments/:id`

**Authentication:** Required

**Description:** Download/view an attachment file.

**Response:** File stream with appropriate Content-Type header

---

### Delete Attachment

**Endpoint:** `DELETE /api/attachments/:id`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Attachment deleted successfully"
  }
}
```

---

## Tags

### List Tags

**Endpoint:** `GET /api/tags`

**Authentication:** Required

**Description:** Get all tags for the authenticated user.

**Response:**
```json
{
  "success": true,
  "data": {
    "tags": [
      {
        "id": "bb0e8400-e29b-41d4-a716-446655440006",
        "name": "breakout",
        "color": "#3b82f6",
        "usage_count": 23,
        "created_at": "2024-01-10T08:00:00Z"
      },
      {
        "id": "cc0e8400-e29b-41d4-a716-446655440007",
        "name": "mistake",
        "color": "#ef4444",
        "usage_count": 8,
        "created_at": "2024-01-12T10:30:00Z"
      }
    ]
  }
}
```

---

### Create Tag

**Endpoint:** `POST /api/tags`

**Authentication:** Required

**Request:**
```json
{
  "name": "earnings-play",
  "color": "#10b981"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "dd0e8400-e29b-41d4-a716-446655440008",
    "name": "earnings-play",
    "color": "#10b981",
    "created_at": "2024-01-20T18:00:00Z"
  }
}
```

---

### Add Tag to Trade

**Endpoint:** `POST /api/trades/:id/tags`

**Authentication:** Required

**Request:**
```json
{
  "tag_id": "bb0e8400-e29b-41d4-a716-446655440006"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Tag added to trade"
  }
}
```

---

### Remove Tag from Trade

**Endpoint:** `DELETE /api/trades/:trade_id/tags/:tag_id`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Tag removed from trade"
  }
}
```

---

## Rule Sets

### List Rule Sets

**Endpoint:** `GET /api/rulesets`

**Authentication:** Required

**Description:** Get all rule sets for the authenticated user.

**Response:**
```json
{
  "success": true,
  "data": {
    "rulesets": [
      {
        "id": "ruleset-550e8400-e29b-41d4-a716-446655440000",
        "user_id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Day Trading Rules",
        "description": "Core rules for intraday trading",
        "is_active": true,
        "rules": [
          {
            "id": "rule-770e8400-e29b-41d4-a716-446655440002",
            "title": "Risk no more than 2% per trade",
            "description": "Position size should be calculated to risk maximum 2% of account",
            "weight": 5,
            "phase": "PRE_TRADE",
            "category": "RISK_MANAGEMENT",
            "created_at": "2024-01-10T08:00:00Z"
          }
        ],
        "created_at": "2024-01-10T08:00:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

---

### Get Single Rule Set

**Endpoint:** `GET /api/rulesets/:id`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "ruleset-550e8400-e29b-41d4-a716-446655440000",
    "user_id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Day Trading Rules",
    "description": "Core rules for intraday trading",
    "is_active": true,
    "rules": [],
    "created_at": "2024-01-10T08:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Status Codes:**
- `200`: Success
- `404`: Rule set not found

---

### Create Rule Set

**Endpoint:** `POST /api/rulesets`

**Authentication:** Required

**Request:**
```json
{
  "name": "Swing Trading Rules",
  "description": "Rules for multi-day positions",
  "is_active": false
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "ruleset-880e8400-e29b-41d4-a716-446655440003",
    "user_id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Swing Trading Rules",
    "description": "Rules for multi-day positions",
    "is_active": false,
    "rules": [],
    "created_at": "2024-01-20T16:00:00Z",
    "updated_at": "2024-01-20T16:00:00Z"
  }
}
```

**Status Codes:**
- `201`: Created
- `400`: Validation error

---

### Update Rule Set

**Endpoint:** `PUT /api/rulesets/:id`

**Authentication:** Required

**Request:**
```json
{
  "name": "Updated Swing Trading Rules",
  "is_active": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "ruleset-880e8400-e29b-41d4-a716-446655440003",
    "updated_at": "2024-01-20T16:30:00Z"
  }
}
```

---

### Delete Rule Set

**Endpoint:** `DELETE /api/rulesets/:id`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Rule set deleted successfully"
  }
}
```

**Status Codes:**
- `200`: Success
- `404`: Rule set not found

---

### Add Rule to Rule Set

**Endpoint:** `POST /api/rulesets/:id/rules`

**Authentication:** Required

**Request:**
```json
{
  "title": "Wait for confirmation candle",
  "description": "Always wait for a confirmation candle before entering",
  "weight": 4,
  "phase": "PRE_TRADE",
  "category": "ENTRY"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "rule-990e8400-e29b-41d4-a716-446655440004",
    "ruleset_id": "ruleset-880e8400-e29b-41d4-a716-446655440003",
    "title": "Wait for confirmation candle",
    "description": "Always wait for a confirmation candle before entering",
    "weight": 4,
    "phase": "PRE_TRADE",
    "category": "ENTRY",
    "created_at": "2024-01-20T16:45:00Z"
  }
}
```

**Valid Phases:**
- `PRE_TRADE`: Rules to check before entering position
- `DURING_TRADE`: Rules to follow while position is open
- `POST_TRADE`: Rules for reflection after closing position

**Valid Categories:**
- `RISK_MANAGEMENT`
- `ENTRY`
- `EXIT`
- `POSITION_SIZING`
- `TIMING`
- `PSYCHOLOGY`
- `GENERAL`

**Status Codes:**
- `201`: Created
- `400`: Validation error
- `404`: Rule set not found

---

### Update Rule

**Endpoint:** `PUT /api/rulesets/:ruleset_id/rules/:rule_id`

**Authentication:** Required

**Request:**
```json
{
  "title": "Wait for strong confirmation candle",
  "weight": 5
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "rule-990e8400-e29b-41d4-a716-446655440004",
    "updated_at": "2024-01-20T17:00:00Z"
  }
}
```

---

### Delete Rule

**Endpoint:** `DELETE /api/rulesets/:ruleset_id/rules/:rule_id`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Rule deleted successfully"
  }
}
```

**Status Codes:**
- `200`: Success
- `404`: Rule or rule set not found

---

## Metrics

### Get Summary Metrics

**Endpoint:** `GET /api/metrics/summary`

**Authentication:** Required

**Description:** Get overall trading performance metrics.

**Query Parameters:**
- `from` (optional): Start date
- `to` (optional): End date

**Response:**
```json
{
  "success": true,
  "data": {
    "total_trades": 145,
    "winning_trades": 89,
    "losing_trades": 56,
    "win_rate": 61.38,
    "total_pnl": 12450.75,
    "average_win": 245.80,
    "average_loss": -132.50,
    "profit_factor": 1.85,
    "largest_win": 1250.00,
    "largest_loss": -450.00,
    "current_streak": {
      "type": "win",
      "count": 3
    },
    "longest_win_streak": 7,
    "longest_loss_streak": 4,
    "total_fees": 342.50,
    "sharpe_ratio": 1.45
  }
}
```

---

### Get Metrics by Symbol

**Endpoint:** `GET /api/metrics/by-symbol`

**Authentication:** Required

**Query Parameters:**
- `from` (optional): Start date
- `to` (optional): End date
- `limit` (optional, default: 10): Number of symbols to return

**Response:**
```json
{
  "success": true,
  "data": {
    "symbols": [
      {
        "symbol": "AAPL",
        "total_trades": 23,
        "winning_trades": 15,
        "win_rate": 65.22,
        "total_pnl": 2340.50,
        "average_pnl": 101.76
      },
      {
        "symbol": "TSLA",
        "total_trades": 18,
        "winning_trades": 10,
        "win_rate": 55.56,
        "total_pnl": 1890.25,
        "average_pnl": 105.01
      }
    ]
  }
}
```

---

### Get Daily Performance

**Endpoint:** `GET /api/metrics/daily`

**Authentication:** Required

**Query Parameters:**
- `from` (optional): Start date
- `to` (optional): End date

**Response:**
```json
{
  "success": true,
  "data": {
    "daily_performance": [
      {
        "date": "2024-01-15",
        "trades": 3,
        "pnl": 450.75,
        "win_rate": 66.67
      },
      {
        "date": "2024-01-16",
        "trades": 2,
        "pnl": -120.50,
        "win_rate": 0
      }
    ]
  }
}
```

---

## WebSocket Notifications

### Connect to Notifications

**Endpoint:** `GET /api/ws`

**Authentication:** Required (JWT token)

**Description:** Establishes a WebSocket connection for real-time notifications.

**Connection:**
```javascript
const ws = new WebSocket('wss://api.tradepulse.drivenw.com/api/ws?token=<jwt_token>');
```

**Message Format:**
All messages from server are JSON:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "trade.created",
  "user_id": "660e8400-e29b-41d4-a716-446655440001",
  "title": "Trade Created",
  "message": "Successfully created trade for AAPL",
  "data": {
    "trade_id": "770e8400-e29b-41d4-a716-446655440002",
    "symbol": "AAPL",
    "pnl": 125.50
  },
  "timestamp": "2024-01-20T15:30:00Z"
}
```

**Notification Types:**
- `trade.created` - New trade created
- `trade.updated` - Trade modified
- `trade.deleted` - Trade removed
- `journal.created` - New journal entry
- `journal.updated` - Journal entry modified
- `csv.import` - CSV import completed
- `error` - Error occurred
- `info` - Informational message
- `success` - Success message

**Connection Lifecycle:**
- Send ping/pong messages for keepalive (60s timeout)
- Automatically reconnect on disconnect
- Connection closed on logout

**Status Codes:**
- `1000`: Normal closure
- `1001`: Going away
- `1006`: Abnormal closure (network issue)

---

### Get Notification Stats

**Endpoint:** `GET /api/notifications/stats`

**Authentication:** Required

**Description:** Get statistics about connected WebSocket clients.

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

---

## Error Codes

| Code | Description |
|------|-------------|
| `INVALID_EMAIL` | Email format is invalid |
| `MAGIC_LINK_EXPIRED` | Magic link has expired (15 min limit) |
| `MAGIC_LINK_USED` | Magic link already used |
| `INVALID_TOKEN` | JWT token is invalid or malformed |
| `TOKEN_EXPIRED` | JWT token has expired |
| `UNAUTHORIZED` | Authentication required |
| `FORBIDDEN` | User doesn't have access to resource |
| `NOT_FOUND` | Resource not found |
| `VALIDATION_ERROR` | Request validation failed |
| `DUPLICATE_ENTRY` | Resource already exists |
| `FILE_TOO_LARGE` | Uploaded file exceeds size limit |
| `INVALID_FILE_TYPE` | File type not allowed |
| `RATE_LIMIT_EXCEEDED` | Too many requests |
| `INTERNAL_ERROR` | Server error |

---

## Rate Limiting

- Authentication endpoints: 5 requests per 15 minutes per IP
- All other endpoints: 100 requests per minute per user
- File upload endpoints: 20 requests per hour per user

Rate limit headers included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1705843200
```

---

## CORS

Allowed origins (development):
- `https://tradepulse.drivenw.com`

Allowed methods:
- GET, POST, PUT, DELETE, OPTIONS

Allowed headers:
- Authorization, Content-Type

---

## Versioning

API version is included in the URL path. Current version: `v1`

Future versions will be accessible at `/api/v2`, `/api/v3`, etc.
