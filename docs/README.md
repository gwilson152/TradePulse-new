# TradePulse Documentation

**Version:** 2.0
**Last Updated:** January 2025
**Status:** Active Development

## Quick Links

- [Project Overview](./project.md)
- [Architecture](./architecture.md)
- [API Specification](./api-spec.md)
- [Authentication](./authentication.md)

### Frontend
- [Frontend Documentation Index](./frontend/readme.md)
- [Getting Started](./frontend/getting-started.md)
- [Design System](./frontend/design-system.md)
- [Component Library](./frontend/component-library.md)
- [Settings & Rules](./frontend/settings-and-rules.md)
- [Implementation Status](./frontend/implementation-status.md)

### Backend
- [Getting Started](./backend/getting-started.md)
- [Database Schema](./backend/database.md)
- [API Patterns](./backend/api-patterns.md)
- [Trades API](./backend/trades-api.md)
- [Notifications](./backend/notifications.md)
- [Implementation Status](./backend/implementation-status.md)

## What is TradePulse?

TradePulse is a modern, full-stack trading journal and analytics platform designed to help traders:
- Track trading performance with detailed analytics
- Maintain discipline through rule-based trading systems
- Improve trading psychology through systematic journaling
- Analyze emotional patterns and their correlation to P&L

## Technology Stack

### Frontend
- **Framework:** SvelteKit 2.0 with Svelte 5 Runes
- **Styling:** Tailwind CSS 3.4+
- **Charts:** Apache ECharts
- **Icons:** Iconify
- **State Management:** Svelte Stores
- **Build Tool:** Vite

### Backend
- **Language:** Go 1.21+
- **Framework:** Chi Router
- **Database:** PostgreSQL 15+
- **Authentication:** JWT + Magic Links
- **Email:** Resend API
- **Real-time:** WebSockets

## Recent Major Updates (November 2025)

### Plan Selection & Signup Flow
- ✅ Implemented signup with plan selection (Starter, Pro, Premium)
- ✅ Beta free status - all plans free during Beta period
- ✅ New POST /api/auth/signup endpoint
- ✅ Plan validation and database constraints
- ✅ Migration 003 for user plan fields
- ✅ PropReports API integration completed

### Server-Side Pagination & Filtering
- ✅ Implemented database-level pagination for trades
- ✅ Added advanced filtering (status, type, date range, strategy, P&L range)
- ✅ Pagination metadata (total, page, page_size, total_pages)
- ✅ Efficient SQL queries with proper indexing

### Timezone Support
- ✅ Created user settings store with timezone preferences
- ✅ Support for local time vs market time display
- ✅ Configurable date/time formats (12h/24h, short/medium/long)
- ✅ Common timezone presets (ET, CT, MT, PT, London, Tokyo, etc.)

### UI Improvements
- ✅ Data table format for trades list
- ✅ Mouse-following tooltips with trade details
- ✅ Mobile long-press support for tooltips
- ✅ Compact pagination controls
- ✅ Advanced filter panel with live result counts

### Bug Fixes
- ✅ Fixed Svelte 5 `{@const}` tag placement issues
- ✅ Fixed accessibility warnings (labels, interactive elements)
- ✅ Fixed reactive effect infinite loops
- ✅ Fixed null reference errors in derived values

## Key Features

### 1. Trade Management
- Import trades from CSV (DAS Trader, Interactive Brokers, etc.)
- Manual trade entry with detailed fields
- Position lifecycle tracking (multiple entries/exits)
- Real-time P&L calculations
- Server-side pagination for large datasets
- Advanced filtering and search

### 2. Journaling System
- Rich text journal entries per trade
- Emotional state tracking (confidence, stress, discipline)
- Screenshot uploads with lightbox viewer
- Voice note recording (up to 5 minutes)
- Rule adherence scoring
- Tabbed interface for organized data entry

### 3. Rule-Based Trading
- Custom rule sets per trading strategy
- Weighted rule importance (1-5 stars)
- Traffic light scoring (Perfect, Good, Partial, Poor, Failed)
- Phase-based organization (Pre-trade, During, Post-trade)
- Visual adherence tracking
- Correlation analysis with performance

### 4. Analytics & Insights
- 6 interactive ECharts visualizations
- P&L over time (line chart)
- Win rate analysis (donut chart)
- Trade distribution by symbol (bar chart)
- Emotional state correlation (scatter plot)
- Rule adherence trends (radar chart)
- Customizable time ranges

### 5. Real-Time Notifications
- WebSocket-based notification system
- Toast notifications for important events
- Notification center with history
- Mark as read/unread functionality
- Event types: trade_created, trade_updated, etc.

## Documentation Structure

```
docs/
├── README.md                    # This file
├── project.md                   # Comprehensive project overview
├── architecture.md              # System architecture
├── api-spec.md                  # API endpoints and schemas
├── authentication.md            # Auth flow documentation
├── csv-import.md               # CSV import guide
├── websocket-notifications.md  # WebSocket setup
│
├── frontend/
│   ├── getting-started.md      # Frontend setup guide
│   ├── design-system.md        # Design tokens and theming
│   ├── component-library.md   # Reusable components
│   ├── components.md           # Component details
│   ├── settings-and-rules.md  # Settings & rules features
│   ├── analytics-enhancements.md # Analytics implementation
│   └── implementation-status.md # Frontend progress
│
└── backend/
    ├── getting-started.md      # Backend setup guide
    ├── structure.md            # Project structure
    ├── database.md             # Database schema
    ├── migrations.md           # Migration guide
    ├── api-patterns.md         # API design patterns
    ├── trades-api.md           # Trades API details
    ├── authentication.md       # Auth implementation
    ├── notifications.md        # Notification system
    ├── email.md                # Email service setup
    └── implementation-status.md # Backend progress
```

## Getting Started

### Prerequisites
- Node.js 18+ (for frontend)
- Go 1.21+ (for backend)
- PostgreSQL 15+
- Git

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd TradePulse
   ```

2. **Setup Frontend**
   ```bash
   cd frontend
   npm install
   cp .env.example .env
   npm run dev
   ```

3. **Setup Backend**
   ```bash
   cd backend
   go mod download
   cp .env.example .env
   # Configure database connection in .env
   go run cmd/api/main.go
   ```

4. **Setup Database**
   ```bash
   cd backend
   # Run migrations
   psql -U postgres -d tradepulse < migrations/*.sql
   ```

See detailed setup guides:
- [Frontend Getting Started](./frontend/getting-started.md)
- [Backend Getting Started](./backend/getting-started.md)

## API Overview

### Base URL
```
Production: https://api.tradepulse.drivenw.com:9000
Development: http://localhost:9000
```

### Authentication
All API endpoints (except auth endpoints) require a JWT token:
```
Authorization: Bearer <jwt_token>
```

### Core Endpoints

**Authentication**
- `POST /api/auth/signup` - Signup with plan selection
- `POST /api/auth/request-magic-link` - Request magic link
- `GET /api/auth/verify` - Verify magic link token
- `POST /api/auth/login` - Email/password login
- `GET /api/auth/me` - Get current user
- `POST /api/auth/set-password` - Set/update password
- `POST /api/auth/logout` - Logout

**Trades**
- `GET /api/trades` - List trades (paginated)
- `POST /api/trades` - Create trade
- `GET /api/trades/{id}` - Get trade details
- `PUT /api/trades/{id}` - Update trade
- `DELETE /api/trades/{id}` - Delete trade
- `POST /api/trades/import` - Bulk import trades

**Journal**
- `GET /api/journal` - List journal entries
- `POST /api/journal` - Create entry
- `GET /api/journal/{id}` - Get entry
- `PUT /api/journal/{id}` - Update entry
- `DELETE /api/journal/{id}` - Delete entry

**Rules**
- `GET /api/rules` - List rule sets
- `POST /api/rules` - Create rule set
- `PUT /api/rules/{id}` - Update rule set
- `DELETE /api/rules/{id}` - Delete rule set

**Analytics**
- `GET /api/analytics/summary` - Get summary stats
- `GET /api/analytics/pnl-over-time` - P&L time series
- `GET /api/analytics/correlation` - Correlation analysis

See [API Specification](./api-spec.md) for complete details.

## Pagination

All list endpoints support pagination with the following parameters:

**Query Parameters:**
- `limit` - Number of items per page (default: 25, max: 100)
- `offset` - Number of items to skip (default: 0)

**Response Format:**
```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "total": 150,
    "page": 1,
    "page_size": 25,
    "total_pages": 6
  }
}
```

## Contributing

1. Create a feature branch from `main`
2. Make your changes
3. Update relevant documentation
4. Submit a pull request

## Support

For questions or issues:
- Check existing documentation
- Review [Implementation Status](./frontend/implementation-status.md)
- Create an issue in the repository

## License

[Add license information]

---

**Note:** This is an active development project. Documentation is updated regularly as features are implemented.
