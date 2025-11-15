# TradePulse Architecture

## System Overview

TradePulse is a trading metrics and journaling platform with real-time notifications.

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser                              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │           SvelteKit 5 Frontend                        │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │  │
│  │  │  Pages   │  │  Stores  │  │   Components     │   │  │
│  │  │ (Routes) │  │  (State) │  │  (Skeleton UI)   │   │  │
│  │  └──────────┘  └──────────┘  └──────────────────┘   │  │
│  │       │              │                  │            │  │
│  │       └──────────────┴──────────────────┘            │  │
│  │                      │                                │  │
│  │              ┌───────▼────────┐                       │  │
│  │              │   API Client   │                       │  │
│  │              └────────────────┘                       │  │
│  └─────────────────────│──────────────────────────────────┘
│                        │                                    │
└────────────────────────┼────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │ HTTPS          │ WSS            │
        │ (REST API)     │ (WebSocket)    │
        └────────────────┼────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│              External Proxy (Nginx)                          │
│  • Routes HTTPS (443) → Internal Ports                       │
│  • SSL/TLS Termination                                       │
│  • WebSocket Upgrade Support                                 │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
   Port 4000        Port 9000             │
   Frontend         Backend API           │
        │                │                │
┌───────▼─────┐   ┌──────▼──────────────────────────────┐
│  Vite Dev   │   │      Go API Server                  │
│   Server    │   │  ┌────────────────────────────────┐ │
│             │   │  │        Chi Router              │ │
│             │   │  │  ┌──────────┐  ┌────────────┐ │ │
│             │   │  │  │ Handlers │  │ Middleware │ │ │
│             │   │  │  └──────────┘  └────────────┘ │ │
│             │   │  └────────────────────────────────┘ │
│             │   │  ┌────────────────────────────────┐ │
│             │   │  │    Notification Bus            │ │
│             │   │  │  • WebSocket Manager           │ │
│             │   │  │  • Client Connections          │ │
│             │   │  │  • Message Broadcasting        │ │
│             │   │  └────────────────────────────────┘ │
│             │   │  ┌────────────────────────────────┐ │
│             │   │  │    Email Service               │ │
│             │   │  │  • SMTP / M365 / Gmail         │ │
│             │   │  │  • Magic Link Sender           │ │
│             │   │  └────────────────────────────────┘ │
└─────────────┘   └──────┬──────────────────────────────┘
                         │
                         │ SQL
                         │
                  ┌──────▼──────┐
                  │ PostgreSQL  │
                  │   Database  │
                  │             │
                  │ postgres1   │
                  │ .drivenw    │
                  │ .local:5432 │
                  └─────────────┘
```

## Tech Stack

### Frontend
- **Framework**: SvelteKit 5 (TypeScript)
- **Styling**: Tailwind CSS + Skeleton UI
- **State**: Svelte stores
- **Build**: Vite

### Backend
- **Language**: Go 1.21+
- **Router**: Chi
- **WebSocket**: Gorilla WebSocket
- **Auth**: JWT tokens

### Database
- **System**: PostgreSQL 15+
- **Host**: postgres1.drivenw.local:5432
- **Migrations**: SQL files

### Infrastructure
- **Proxy**: Nginx (external)
- **Protocol**: HTTPS/WSS
- **Ports**: 4000 (frontend), 9000 (backend)

## Data Flow

### Authentication Flow

```
User → Email Address
         ↓
Frontend → POST /api/auth/request-magic-link
         ↓
Backend → Generate Token + Store in DB
         ↓
Email Service → Send Magic Link
         ↓
User Clicks Link → GET /api/auth/verify?token=...
         ↓
Backend → Validate Token + Generate JWT
         ↓
Frontend → Store JWT + Redirect to Dashboard
         ↓
WebSocket → Connect with JWT
```

### Trade Creation Flow

```
User → Fill Trade Form
         ↓
Frontend → POST /api/trades + JWT
         ↓
Backend Handler → Validate Data
         ↓
Database → INSERT INTO trades
         ↓
Notification Bus → Publish "trade.created"
         ↓
WebSocket → Broadcast to User
         ↓
Frontend → Update UI + Show Notification
```

### Real-time Notification Flow

```
Backend Action (any handler)
         ↓
bus.Publish(type, userID, title, message, data)
         ↓
Notification Bus → Find User's WebSocket Clients
         ↓
WebSocket → Send JSON Message
         ↓
Frontend Store → Add Notification
         ↓
UI → Update Bell Badge + Show Toast/Browser Notification
```

## Key Components

### Backend Packages

**`cmd/api`** - Application entry point
- Initialize database, bus, router
- Start HTTP server

**`internal/handlers`** - HTTP request handlers
- Auth, trades, journal, metrics
- WebSocket upgrade

**`internal/middleware`** - HTTP middleware
- JWT authentication
- CORS handling

**`internal/notifications`** - WebSocket system
- Notification bus (central hub)
- Client manager (per-user connections)

**`internal/models`** - Data structures
- User, Trade, Journal, Attachment

**`internal/database`** - Database layer
- Connection management
- Query execution

**`internal/email`** - Email providers
- SMTP, M365 Graph, Gmail
- Abstracted interface

### Frontend Structure

**`src/routes`** - File-based routing
- `+page.svelte` - Landing page
- `auth/` - Login/verify pages
- `app/` - Authenticated routes (dashboard, trades, journal)

**`src/lib/stores`** - Svelte stores
- `notifications.ts` - Notification state
- `auth.ts` - User session
- `trades.ts` - Trade data

**`src/lib/api`** - API communication
- `client.ts` - HTTP client with JWT
- `websocket.ts` - WebSocket client

**`src/lib/components`** - Reusable components
- `ui/` - Base UI wrappers (Skeleton UI)
- `trading/` - Trade-specific components
- `journal/` - Journal components
- `notifications/` - Notification UI

## Security

### Authentication
- Magic link tokens (15-minute expiry, one-time use)
- JWT tokens for API requests (24-hour expiry)
- Secure, HTTP-only cookies for session

### Authorization
- User ID extracted from JWT
- All queries filtered by user ID
- No cross-user data access

### Database
- Parameterized queries only
- UUID primary keys
- Foreign key constraints
- Triggers for auto-calculations

### WebSocket
- JWT required for connection
- User-specific message routing
- No broadcast to all users

## Performance

### Backend
- Connection pooling (25 max, 5 idle)
- Notification buffer (256 messages)
- Concurrent request handling (Go routines)

### Frontend
- Code splitting (route-based)
- Lazy loading (components)
- WebSocket auto-reconnect
- Notification history limit (50)

### Database
- Indexes on frequently queried columns
- Triggers for computed fields (P&L)
- Pagination for large result sets

## Deployment

### Development
- Frontend: https://tradepulse.drivenw.com → :4000
- API: https://api.tradepulse.drivenw.com → :9000
- External proxy handles SSL and routing

### Production
- Standard HTTPS ports (443)
- Database connection pooling
- Automated backups
- Monitoring and logging

## Scalability Considerations

### Current Architecture
- Single server deployment
- Direct database connections
- In-memory notification bus

### Future Scaling
- **Horizontal**: Load balancer + multiple API servers
- **WebSocket**: Redis pub/sub for cross-server notifications
- **Database**: Read replicas for queries
- **Cache**: Redis for session/frequently accessed data
- **Storage**: S3 for attachments (screenshots, voice notes)
- **CDN**: Static asset delivery

## File Organization

```
TradePulse/
├── frontend/              # SvelteKit application
│   ├── src/
│   │   ├── routes/       # Pages
│   │   └── lib/          # Reusable code
│   └── static/           # Public assets
├── backend/              # Go API
│   ├── cmd/api/          # Entry point
│   ├── internal/         # Private packages
│   └── migrations/       # SQL migrations
├── docs/                 # Documentation
│   ├── architecture.md   # This file
│   ├── backend/          # Backend guides
│   └── frontend/         # Frontend guides
└── deployment/           # Config files
```

## Development Workflow

1. **Start Database** - Ensure PostgreSQL is accessible
2. **Start Backend** - `go run cmd/api/main.go` (port 9000)
3. **Start Frontend** - `npm run dev` (port 4000)
4. **Access App** - https://tradepulse.drivenw.com
5. **Monitor** - Check logs, WebSocket stats

## Testing Strategy

### Backend
- Unit tests for business logic
- Integration tests for handlers
- Database tests with test fixtures

### Frontend
- Component tests (Vitest)
- E2E tests for critical flows (Playwright)
- Visual regression tests

### System
- Manual testing via proxy URLs
- WebSocket connection testing
- Load testing for concurrent users
