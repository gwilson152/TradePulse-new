# TradePulse - Trading Journal & Analytics Platform

## Project Overview

TradePulse is a modern, full-stack trading journal and analytics platform designed to help traders track their performance, maintain discipline, and improve their trading psychology through systematic journaling and data-driven insights.

**Current Version:** 2.0
**Last Updated:** January 2025
**Status:** Frontend Complete | Backend In Progress

---

## Vision & Differentiation

### What Makes TradePulse Unique

1. **Position Lifecycle Management**

   - Track multiple entries and exits per position
   - Scale in/out with precision tracking
   - Visual timeline of all position events
   - Support for FIFO, LIFO, and Average cost basis methods
   - Real-time P&L calculation (realized vs unrealized)

2. **Rule-Based Trading Discipline**

   - Create custom rule sets for different strategies
   - Traffic light scoring system (0-100% with 5 levels)
   - Weighted rule importance (1-5 stars)
   - Phase-based organization (Pre-trade, During-trade, Post-trade)
   - Visual adherence tracking and correlation analysis

3. **Comprehensive Journaling**

   - Rich text reflection entries
   - Emotional state tracking (Confidence, Stress, Discipline)
   - Screenshot uploads with lightbox viewer
   - Voice note recording and playback
   - Rule adherence scoring per entry
   - Automatic weighted adherence score calculation

4. **Advanced Analytics**

   - 6 interactive ECharts visualizations
   - Correlation analysis (emotions vs P&L, rules vs performance)
   - Time-series P&L tracking
   - Win rate and trade distribution analysis
   - Customizable time ranges (1W, 1M, 3M, 6M, 1Y, All)

5. **Fresh macOS-Inspired Design**
   - Unique dock-style bottom navigation
   - Glassmorphism UI with backdrop blur effects
   - No traditional sidebar - completely fresh layout
   - Gradient-based color system with colored shadows
   - Smooth animations and micro-interactions
   - Full dark mode support
   - WCAG AA accessibility compliant

---

## Technology Stack

### Frontend

**Framework & Core:**

- **Svelte 5.0** - Latest version with runes ($state, $derived, $props)
- **SvelteKit 2.0** - Full-stack framework with SSR and routing
- **TypeScript 5.x** - Full type safety throughout
- **Vite 5.0** - Lightning-fast build tool

**Styling:**

- **Tailwind CSS 3.4** - Utility-first CSS framework
- **Skeleton UI 2.11** - Component foundation (minimal usage)
- **Custom Design System** - Glassmorphism with gradients

**Libraries:**

- **Apache ECharts 5.x** - Interactive charts and visualizations
- **Iconify** - Icon system with Material Design Icons (outline style)
- **MediaRecorder API** - Native browser voice recording

**Key Frontend Features:**

- Glassmorphism design with backdrop blur
- Responsive breakpoints (mobile-first approach)
- Full accessibility (WCAG AA compliant)
- Keyboard navigation throughout
- Drag-and-drop file uploads
- Image lightbox with keyboard controls
- Voice note recording/playback

### Backend

**Runtime & Framework:**

- **Go 1.21+** - High-performance compiled language
- **chi v5** - Lightweight, composable HTTP router
- **database/sql** - Standard library database interface
- **lib/pq** - PostgreSQL driver

**Database:**

- **PostgreSQL 15+** - Primary relational database
- **Database Host:** postgres1.drivenw.local:5432
- **Database Name:** tradepulse
- **Connection pooling** - 25 max open, 5 max idle connections
- **Automatic migrations** - Embedded migration system

**Authentication:**

- **Passwordless magic links** via email
- **JWT tokens** (golang-jwt/jwt) for session management
- **Cryptographically secure** token generation
- **15-minute expiry** on magic links
- **24-hour expiry** on JWTs (configurable)

**Real-time Communication:**

- **WebSocket** (gorilla/websocket) - Live notifications
- **Event-driven architecture** - Central notification bus
- **Multi-client support** - Per-user notification channels

**Email Delivery:**

- **Nodemailer** - Email client
- **SMTP Support** - Multiple provider support
- **Trusted relay** for internal networks

**File Storage:**

- **Local file system** for uploads
- **URL-based access** for screenshots and voice notes
- **10MB max file size** limit

---

## Infrastructure

### Deployment Architecture

**Frontend:**

- URL: https://tradepulse.drivenw.com
- Internal Port: 4000
- External Proxy: HTTPS (443) → Internal (4000)
- Framework: SvelteKit (Node adapter)

**Backend API:**

- URL: https://api.tradepulse.drivenw.com
- Internal Port: 9000
- External Proxy: HTTPS (443) → Internal (9000)
- Framework: Express.js

**Database:**

- Host: postgres1.drivenw.local:5432
- Database: tradepulse
- User: tradepulse
- Version: PostgreSQL 15+

**Email Options:**

1. SMTP with trusted relay
2. Microsoft 365 Graph API
3. Google Workspace API

---

## Database Schema

### Core Tables

**users**

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  last_login TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NOW()
);
```

**trades**

```sql
CREATE TABLE trades (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  symbol VARCHAR(20) NOT NULL,
  trade_type VARCHAR(10) NOT NULL CHECK (trade_type IN ('LONG', 'SHORT')),

  -- Position lifecycle fields
  entries JSONB NOT NULL DEFAULT '[]',
  exits JSONB NOT NULL DEFAULT '[]',
  current_position_size DECIMAL(20, 8) NOT NULL DEFAULT 0,
  average_entry_price DECIMAL(20, 8),
  total_fees DECIMAL(20, 2) DEFAULT 0,
  realized_pnl DECIMAL(20, 2) DEFAULT 0,
  unrealized_pnl DECIMAL(20, 2),
  cost_basis_method VARCHAR(10) DEFAULT 'FIFO' CHECK (cost_basis_method IN ('FIFO', 'LIFO', 'AVERAGE')),

  -- Timestamps
  opened_at TIMESTAMP NOT NULL,
  closed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  -- Additional fields
  notes TEXT,
  tags VARCHAR(50)[]
);
```

**journal_entries**

```sql
CREATE TABLE journal_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  trade_id UUID REFERENCES trades(id) ON DELETE SET NULL,
  entry_date TIMESTAMP NOT NULL DEFAULT NOW(),
  content TEXT NOT NULL,

  -- Emotional state tracking
  emotional_state JSONB,  -- {confidence, stress, discipline, notes}

  -- Rule adherence
  rule_adherence JSONB DEFAULT '[]',  -- Array of {rule_id, rule_title, score, notes, timestamp}
  adherence_score DECIMAL(5, 2),  -- Calculated weighted average

  -- Media attachments
  screenshots TEXT[] DEFAULT '{}',
  voice_notes TEXT[] DEFAULT '{}',

  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
```

**rule_sets**

```sql
CREATE TABLE rule_sets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  rules JSONB NOT NULL DEFAULT '[]',  -- Array of Rule objects
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
```

**magic_links**

```sql
CREATE TABLE magic_links (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(255) NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  used_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW()
);
```

### JSONB Structures

**Entry (in trades.entries)**

```typescript
{
  id: string;
  price: number;
  quantity: number;
  timestamp: string;  // ISO 8601
  notes?: string;
  fees?: number;
}
```

**Exit (in trades.exits)**

```typescript
{
  id: string;
  price: number;
  quantity: number;
  timestamp: string;  // ISO 8601
  notes?: string;
  fees?: number;
  pnl: number;  // Calculated
}
```

**Rule (in rule_sets.rules)**

```typescript
{
  id: string;
  title: string;
  description: string;
  weight: number; // 1-5
  phase: "PRE_TRADE" | "DURING_TRADE" | "POST_TRADE";
  category: "RISK_MANAGEMENT" |
    "ENTRY" |
    "EXIT" |
    "POSITION_SIZING" |
    "TIMING" |
    "PSYCHOLOGY" |
    "GENERAL";
  created_at: string;
}
```

**EmotionalState (in journal_entries.emotional_state)**

```typescript
{
  confidence: number; // 1-10
  stress: number; // 1-10
  discipline: number; // 1-10
  notes: string;
}
```

**RuleAdherence (in journal_entries.rule_adherence)**

```typescript
{
  rule_id: string;
  rule_title: string;
  score: number; // 0, 25, 50, 75, 100
  notes: string;
  timestamp: string;
}
```

---

## API Endpoints

### Authentication

- `POST /api/auth/request-magic-link` - Send magic link to email
- `GET /api/auth/verify?token=...` - Verify magic link and return JWT
- `POST /api/auth/logout` - Invalidate current JWT token
- `POST /api/auth/refresh` - Refresh JWT before expiry

### Trades (Position Management)

- `GET /api/trades` - List trades with pagination and filters
- `GET /api/trades/:id` - Get single trade with full details
- `POST /api/trades` - Create new trade/position
- `PUT /api/trades/:id` - Update trade
- `DELETE /api/trades/:id` - Delete trade

### Position Lifecycle

- `POST /api/trades/:id/entries` - Add entry to position (scale in)
- `POST /api/trades/:id/exits` - Add exit to position (scale out)

### Journal Entries

- `GET /api/journal` - List journal entries with pagination
- `GET /api/journal/:id` - Get single journal entry
- `POST /api/journal` - Create journal entry (supports multipart/form-data)
- `PUT /api/journal/:id` - Update journal entry
- `DELETE /api/journal/:id` - Delete journal entry

### Rule Sets

- `GET /api/rulesets` - List all rule sets for user
- `GET /api/rulesets/:id` - Get single rule set
- `POST /api/rulesets` - Create new rule set
- `PUT /api/rulesets/:id` - Update rule set
- `DELETE /api/rulesets/:id` - Delete rule set
- `POST /api/rulesets/:id/rules` - Add rule to rule set
- `PUT /api/rulesets/:rulesetId/rules/:ruleId` - Update specific rule
- `DELETE /api/rulesets/:rulesetId/rules/:ruleId` - Delete specific rule

### Analytics & Metrics

- `GET /api/metrics/summary` - Overall performance summary
- `GET /api/metrics/analytics?range=1M` - Time-series data for charts
- `GET /api/metrics/by-symbol` - Performance breakdown by symbol
- `GET /api/metrics/daily` - Daily performance data

### File Upload

- `POST /api/upload` - Upload screenshot or voice note (multipart/form-data)

### WebSocket

- `GET /api/ws?token=...` - Real-time notifications connection

See [api-spec.md](./api-spec.md) for complete API documentation.

---

## UI Components Library

### Base UI Components (13)

1. **Button.svelte** - 4 variants (filled, gradient, soft, ghost) × 6 colors
2. **Card.svelte** - 3 variants (glass, solid, elevated) with glassmorphism
3. **Input.svelte** - Text input with validation and accessibility
4. **Select.svelte** - Dropdown select with custom styling
5. **Textarea.svelte** - Multi-line input with character counter
6. **Badge.svelte** - Status indicators in multiple colors
7. **Modal.svelte** - Accessible dialog with backdrop and escape handling
8. **FileUpload.svelte** - Drag-and-drop upload with preview
9. **ImageGallery.svelte** - Responsive grid with lightbox viewer
10. **AudioRecorder.svelte** - Voice note recording with MediaRecorder API
11. **AudioPlayer.svelte** - Audio playback with controls
12. **ChartCard.svelte** - ECharts wrapper with loading states
13. **NotificationBell.svelte** - WebSocket notification indicator

### Trading-Specific Components (9)

1. **MetricCard.svelte** - Dashboard stat cards with gradient icon backgrounds
2. **PnLBadge.svelte** - Color-coded profit/loss display
3. **TradeModal.svelte** - Create/edit trade positions
4. **TradeDetailSlideOver.svelte** - Full position details panel
5. **PositionTimeline.svelte** - Chronological entry/exit visualization
6. **PositionSizeBar.svelte** - Visual position size tracking
7. **AddToPositionModal.svelte** - Scale in/out modal
8. **RuleCard.svelte** - Individual rule display card
9. **RuleAdherenceInput.svelte** - Traffic light scoring interface
10. **AdherenceScoreDisplay.svelte** - Overall adherence widget
11. **JournalEntryModal.svelte** - 4-tab comprehensive entry form

See [component-library.md](./frontend/component-library.md) for detailed component documentation.

---

## Design System

### Layout Architecture

**macOS-Inspired Navigation:**

- Top menu bar (44px) - App name, live clock, notifications, user avatar
- Bottom dock (floating) - 5 main navigation items with colored icons
- Main content area - Max-width 1800px, centered, scrollable
- No traditional sidebar - unique approach

**Color-Coded Sections:**

- 🔵 **Overview** - Blue (`text-blue-500`)
- 🟢 **Trades** - Emerald (`text-emerald-500`)
- 🟣 **Journal** - Purple (`text-purple-500`)
- 🟠 **Analytics** - Orange (`text-orange-500`)
- ⚫ **Settings** - Slate (`text-slate-500`)

### Visual Design

**Glassmorphism:**

```css
bg-white/60 dark:bg-slate-800/60
backdrop-blur-xl
border border-slate-200/50 dark:border-slate-700/50
rounded-2xl
```

**Gradient System:**

- Icon backgrounds with matching colored shadows
- Button gradients for emphasis
- Text gradients for hero headers
- Smooth color transitions

**Icon Strategy:**

- Outline style only (e.g., `mdi:chart-line-variant`)
- Single color per icon (no multi-color)
- Color-coded by section
- 24px standard size in navigation

**Shadows & Depth:**

```css
/* Subtle */
shadow-sm shadow-blue-500/30

/* Medium */
shadow-lg shadow-blue-500/30

/* Large/Glow */
shadow-xl shadow-blue-500/40
```

See [design-system.md](./frontend/design-system.md) for complete design guidelines.

---

## File Structure

```
TradePulse/
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   │   ├── ui/              # 13 base UI components
│   │   │   │   ├── trading/         # 9 trading components
│   │   │   │   └── notifications/   # Notification components
│   │   │   ├── api/
│   │   │   │   └── client.ts        # Centralized API client
│   │   │   ├── types/
│   │   │   │   └── index.ts         # TypeScript definitions
│   │   │   └── utils/               # Helper functions
│   │   ├── routes/
│   │   │   ├── +layout.svelte       # Root layout
│   │   │   ├── +page.svelte         # Landing page
│   │   │   ├── app/
│   │   │   │   ├── +layout.svelte   # App layout (dock nav)
│   │   │   │   ├── dashboard/       # Overview page
│   │   │   │   ├── trades/          # Trade management
│   │   │   │   ├── journal/         # Journal entries
│   │   │   │   ├── analytics/       # Charts & insights
│   │   │   │   └── settings/
│   │   │   │       └── rules/       # Rule set management
│   │   │   └── auth/
│   │   │       ├── login/           # Magic link request
│   │   │       └── verify/          # Magic link verification
│   │   ├── app.html                 # HTML template
│   │   ├── app.css                  # Global styles
│   │   └── hooks.server.ts          # Server hooks
│   ├── static/                      # Static assets
│   ├── svelte.config.js             # SvelteKit config
│   ├── tailwind.config.js           # Tailwind config
│   └── package.json
│
├── backend/
│   ├── src/
│   │   ├── routes/
│   │   │   ├── auth.ts              # Authentication endpoints
│   │   │   ├── trades.ts            # Trade endpoints
│   │   │   ├── journal.ts           # Journal endpoints
│   │   │   ├── rulesets.ts          # Rule set endpoints
│   │   │   ├── metrics.ts           # Analytics endpoints
│   │   │   ├── upload.ts            # File upload
│   │   │   └── websocket.ts         # WebSocket handler
│   │   ├── middleware/
│   │   │   ├── auth.ts              # JWT verification
│   │   │   ├── validation.ts        # Request validation
│   │   │   └── errorHandler.ts      # Global error handling
│   │   ├── services/
│   │   │   ├── authService.ts       # Auth logic
│   │   │   ├── tradeService.ts      # Trade logic
│   │   │   ├── journalService.ts    # Journal logic
│   │   │   ├── emailService.ts      # Email sending
│   │   │   └── fileService.ts       # File operations
│   │   ├── database/
│   │   │   ├── client.ts            # PostgreSQL connection
│   │   │   ├── migrations/          # SQL migrations
│   │   │   └── queries/             # SQL queries
│   │   ├── types/
│   │   │   └── index.ts             # Shared types
│   │   └── server.ts                # Express app entry
│   ├── uploads/                     # File storage
│   │   ├── screenshots/
│   │   └── voice/
│   ├── .env                         # Environment variables
│   ├── .env.example                 # Environment template
│   └── package.json
│
└── docs/
    ├── project.md                   # This file
    ├── architecture.md              # System design
    ├── api-spec.md                  # Complete API docs
    ├── websocket-notifications.md   # WebSocket protocol
    ├── frontend/
    │   ├── design-system.md         # Design guidelines
    │   ├── component-library.md     # Component docs
    │   ├── implementation-status.md # Feature tracking
    │   ├── getting-started.md       # Frontend setup
    │   └── ui-wireframes.md         # UI mockups
    └── backend/
        ├── database.md              # Schema & migrations
        ├── authentication.md        # Auth flow details
        ├── api-patterns.md          # API conventions
        ├── structure.md             # Code organization
        ├── email.md                 # Email setup
        └── getting-started.md       # Backend setup
```

---

## Development Setup

### Prerequisites

- **Node.js** 18+ with npm
- **PostgreSQL** 15+
- **Git** for version control

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

Runs on: http://localhost:5173
Hot reload enabled

### Backend Setup

```bash
cd backend
npm install

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
nano .env

# Run development server
npm run dev
```

Runs on: http://localhost:9000
Auto-restart with nodemon

### Database Setup

```bash
# Create database
createdb tradepulse

# Connect and run initial schema
psql tradepulse < backend/src/database/migrations/001_initial_schema.sql
```

### Environment Variables

**Frontend** (`frontend/.env`):

```env
PUBLIC_API_URL=http://localhost:9000
```

**Backend** (`backend/.env`):

```env
# Server
PORT=9000
NODE_ENV=development

# Database
DATABASE_URL=postgresql://tradepulse:password@postgres1.drivenw.local:5432/tradepulse

# Authentication
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRY=24h

# Email (SMTP)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_SECURE=false
SMTP_USER=your-email@example.com
SMTP_PASSWORD=your-password

# Email (M365 Graph API) - Alternative
MICROSOFT_TENANT_ID=
MICROSOFT_CLIENT_ID=
MICROSOFT_CLIENT_SECRET=

# URLs
FRONTEND_URL=http://localhost:5173
API_URL=http://localhost:9000

# File Upload
MAX_FILE_SIZE=10485760  # 10MB in bytes
UPLOAD_DIR=./uploads
```

---

## Implementation Status

### ✅ Completed

**Frontend (100%):**

- [x] macOS-inspired dock navigation layout
- [x] Glassmorphism design system
- [x] 22 accessible UI components
- [x] Position lifecycle management UI
- [x] Rule set creation and management
- [x] Traffic light scoring interface
- [x] Comprehensive journal entry form (4 tabs)
- [x] 6 interactive analytics charts
- [x] File upload with drag-drop
- [x] Voice note recording/playback
- [x] Image gallery with lightbox
- [x] Dark mode support
- [x] Responsive design (mobile/tablet/desktop)
- [x] Full TypeScript type definitions
- [x] Complete API client
- [x] WCAG AA accessibility compliance
- [x] All accessibility warnings fixed

**Documentation:**

- [x] Design system documentation
- [x] Component library documentation
- [x] API specification
- [x] Project documentation (this file)

### 🔄 In Progress

**Backend (0%):**

- [ ] Database schema implementation
- [ ] Authentication endpoints
- [ ] Trade CRUD operations
- [ ] Position lifecycle endpoints
- [ ] Journal entry endpoints
- [ ] Rule set endpoints
- [ ] Analytics/metrics endpoints
- [ ] File upload handling
- [ ] WebSocket implementation
- [ ] Email service integration

### 📋 Planned

**Features:**

- [ ] CSV import/export for trades
- [ ] PDF report generation
- [ ] Calendar view for trades
- [ ] Tag management system
- [ ] Advanced filtering
- [ ] Trade playbook templates
- [ ] Mobile responsive improvements
- [ ] Progressive Web App (PWA)
- [ ] Mobile app (React Native)

**Technical:**

- [ ] Unit tests (frontend)
- [ ] Integration tests (backend)
- [ ] E2E tests (Playwright)
- [ ] CI/CD pipeline
- [ ] Docker containerization
- [ ] Kubernetes deployment configs

---

## Key Architectural Decisions

### 1. Position Lifecycle Model

**Decision:** Use arrays of entries and exits instead of single entry/exit fields.

**Rationale:**

- Supports real-world trading (scaling in/out)
- Maintains complete position history
- Enables timeline visualization
- Allows accurate P&L tracking per exit

**Trade-offs:**

- More complex data structure
- Requires calculation logic for averages
- Higher storage for active traders

### 2. Rule Adherence Scoring

**Decision:** 5-level scoring (0, 25, 50, 75, 100) with weighted averaging.

**Rationale:**

- Simple enough for quick input
- Granular enough to measure progress
- Weighted system accounts for rule importance
- Color coding provides instant visual feedback

**Trade-offs:**

- Requires discipline to score honestly
- Subjective scoring (not automated)

### 3. macOS-Inspired Dock Navigation

**Decision:** Bottom dock instead of sidebar or top navigation.

**Rationale:**

- Unique, fresh UI that stands out
- More screen space for content
- Familiar pattern for macOS users
- Modern, app-like feel
- Accessible on all screen sizes

**Trade-offs:**

- Less conventional than sidebar
- Requires user adaptation
- Limited to ~5 primary nav items

### 4. Glassmorphism Design

**Decision:** Semi-transparent cards with backdrop blur.

**Rationale:**

- Modern, premium aesthetic
- Depth perception without heavy shadows
- Works well in light and dark mode
- Differentiates from competitors

**Trade-offs:**

- Browser compatibility (backdrop-filter)
- Performance on low-end devices
- May not be timeless (trend-dependent)

### 5. JSONB for Nested Data

**Decision:** Use PostgreSQL JSONB for entries, exits, and rules instead of separate tables.

**Rationale:**

- Simpler schema (fewer joins)
- Faster reads (single query)
- Flexible structure
- Atomic updates for position lifecycle

**Trade-offs:**

- Harder to query individual entries
- No foreign key constraints on nested data
- Potential for data inconsistency

---

## Performance Targets

| Metric                       | Target         | Status       |
| ---------------------------- | -------------- | ------------ |
| **Lighthouse Performance**   | 90+            | TBD          |
| **Lighthouse Accessibility** | 100            | ✅ (WCAG AA) |
| **First Contentful Paint**   | < 1.5s         | TBD          |
| **Time to Interactive**      | < 3s           | TBD          |
| **API Response (p95)**       | < 200ms        | TBD          |
| **Database Query (p95)**     | < 100ms        | TBD          |
| **Bundle Size**              | < 300KB (gzip) | TBD          |

---

## Security Considerations

### Authentication

- ✅ Passwordless (magic links only)
- ✅ Magic links expire after 15 minutes
- ✅ Single-use tokens (marked as used)
- ✅ JWTs expire after 24 hours
- 🔄 Refresh token rotation
- 🔄 Rate limiting on auth endpoints

### Data Protection

- ✅ All API requests over HTTPS in production
- 🔄 JWT in httpOnly cookies (currently localStorage)
- ✅ Input validation on all endpoints
- ✅ Parameterized SQL queries (prevent injection)
- ✅ CORS configuration
- 🔄 Content Security Policy headers

### File Uploads

- ✅ File type validation (MIME check)
- ✅ File size limits (10MB max)
- 🔄 Virus scanning (ClamAV integration planned)
- 🔄 Secure storage with signed URLs
- ✅ Unique filenames (UUID-based)

### Database

- ✅ Separate database user (tradepulse)
- ✅ Row-level security policies (user_id checks)
- 🔄 Encrypted backups
- 🔄 Automated backup retention

---

## Browser Support

| Browser       | Version     | Status       |
| ------------- | ----------- | ------------ |
| Chrome        | 90+         | ✅ Supported |
| Firefox       | 88+         | ✅ Supported |
| Safari        | 14+         | ✅ Supported |
| Edge          | 90+         | ✅ Supported |
| Mobile Safari | iOS 14+     | ✅ Supported |
| Chrome Mobile | Android 90+ | ✅ Supported |

**Required Features:**

- ES2020 JavaScript
- CSS Grid & Flexbox
- backdrop-filter (graceful degradation)
- MediaRecorder API (voice notes)
- Drag and Drop API

---

## Accessibility Standards

**WCAG 2.1 Level AA Compliance:**

- ✅ Color contrast ratios > 4.5:1
- ✅ Keyboard navigation for all interactive elements
- ✅ Screen reader support (ARIA labels)
- ✅ Focus indicators on all focusable elements
- ✅ Form labels associated with controls
- ✅ Semantic HTML structure
- ✅ Skip links for navigation
- ✅ Accessible modals (focus trapping)

**Keyboard Shortcuts:**

- `Escape` - Close modals/lightboxes
- `Tab` - Navigate between elements
- `Enter` - Activate buttons/links
- `Arrow keys` - Navigate image gallery

---

## Testing Strategy

### Frontend

- **Unit Tests:** Vitest for component logic
- **Component Tests:** Testing Library for UI
- **E2E Tests:** Playwright for user flows
- **Accessibility:** axe-core automated checks
- **Visual Regression:** Percy or Chromatic

### Backend

- **Unit Tests:** Jest for business logic
- **Integration Tests:** Supertest for API endpoints
- **Database Tests:** Test database with migrations
- **Load Tests:** Artillery or k6

### Coverage Targets

- **Frontend:** 80% coverage
- **Backend:** 90% coverage (critical paths 100%)

---

## Deployment

### Production Environment

**Frontend (Vercel/Netlify):**

```bash
npm run build
# Deploy build/ folder
# Set PUBLIC_API_URL environment variable
```

**Backend (Railway/Render/AWS):**

```bash
npm run build
npm start
# Set all environment variables
# Configure database connection
```

**Database (Managed PostgreSQL):**

- Supabase (recommended for indie projects)
- Railway (easy integration)
- Neon (serverless PostgreSQL)
- AWS RDS (enterprise)

### Monitoring & Logging

**Frontend:**

- Sentry for error tracking
- Google Analytics or Plausible
- Vercel Analytics

**Backend:**

- Winston or Pino for logging
- Sentry for error tracking
- Database slow query logs

---

## Contributing

### Code Style

- **Linting:** ESLint with TypeScript rules
- **Formatting:** Prettier (single quotes, tabs, 100 char width)
- **Commits:** Conventional Commits format

### Git Workflow

- `main` - Production-ready code
- `develop` - Integration branch
- `feature/*` - New features
- `fix/*` - Bug fixes
- `docs/*` - Documentation updates

### Pull Request Process

1. Create feature branch from `develop`
2. Write tests for new features
3. Ensure all tests pass
4. Update documentation
5. Submit PR with description
6. Code review (1 approval required)
7. Merge to `develop`

---

## Roadmap

### Q1 2025

- [x] Frontend UI complete
- [ ] Backend API implementation
- [ ] Database migrations
- [ ] Authentication flow
- [ ] Core CRUD operations

### Q2 2025

- [ ] File upload implementation
- [ ] WebSocket real-time updates
- [ ] Analytics calculations
- [ ] Testing suite
- [ ] Beta launch

### Q3 2025

- [ ] CSV import/export
- [ ] PDF report generation
- [ ] Mobile app (React Native)
- [ ] Advanced analytics
- [ ] Public launch

### Q4 2025

- [ ] API for third-party integrations
- [ ] Trading platform integrations
- [ ] Collaborative features
- [ ] Premium tier features

---

## License

**Proprietary** - All rights reserved
Copyright © 2025 TradePulse

---

## Changelog

### Version 2.0.0 (January 2025)

**Major Redesign:**

- ✅ Complete UI overhaul with macOS-inspired dock navigation
- ✅ Glassmorphism design system implementation
- ✅ Position lifecycle management (multiple entries/exits)
- ✅ Rule-based trading discipline system
- ✅ Comprehensive journaling with media attachments
- ✅ 6 interactive analytics visualizations
- ✅ Voice note recording and playback
- ✅ Image gallery with lightbox
- ✅ Full dark mode support
- ✅ WCAG AA accessibility compliance
- ✅ Complete TypeScript type safety
- ✅ 22 reusable UI components

**Technical Improvements:**

- ✅ Upgraded to Svelte 5 with runes
- ✅ Modern Tailwind CSS 3.4 setup
- ✅ Removed sidebar, added dock navigation
- ✅ Accessibility fixes (all warnings resolved)
- ✅ Responsive design improvements
- ✅ Performance optimizations

### Version 1.0.0 (Deprecated)

- Basic trade tracking
- Simple journaling
- Sidebar navigation
- Limited analytics

---

## Contact & Support

**Project Repository:** [GitHub URL]
**Documentation:** https://tradepulse.drivenw.com/docs
**API Documentation:** https://api.tradepulse.drivenw.com/docs
**Support Email:** support@tradepulse.drivenw.com

---

## Acknowledgments

- **Svelte Team** - For the incredible reactive framework
- **Tailwind Labs** - For utility-first CSS paradigm
- **Apache ECharts** - For powerful charting library
- **Iconify** - For comprehensive icon ecosystem
- **Skeleton UI** - For component foundation
- **PostgreSQL Community** - For robust database system

---

**Last Updated:** January 13, 2025
**Document Version:** 2.0
**Status:** Frontend Complete | Backend In Progress
