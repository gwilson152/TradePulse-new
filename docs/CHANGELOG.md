# TradePulse Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- Signup with plan selection (Starter, Pro, Premium)
- Beta free status for all users (all plans free during Beta)
- Plan management system with database fields
- Migration 003 for user plan fields
- POST /api/auth/signup endpoint
- Plan type and status validation
- PropReports API integration for trade data
- Server-side pagination for trades API
- Advanced filtering (strategy, P&L range)
- Timezone support with user settings store
- Mouse-following tooltips for trade details
- Mobile long-press support for tooltips
- Compact data table view for trades
- Live result counts in filter panel
- Settings store with localStorage persistence

### Changed
- Trades list now uses server-side pagination
- All filters processed at database level for performance
- Date/time display supports multiple timezones
- Pagination controls show server-provided metadata
- User model includes plan_type, plan_status, plan_selected_at fields

### Fixed
- Svelte 5 `{@const}` tag placement issues
- Accessibility warnings for labels and interactive elements
- Reactive effect infinite loops
- Null reference errors in derived values
- Losers filter not working (was comparison instead of assignment)

## [2.0.0] - 2025-01-17

### Major Features Completed

#### Frontend (SvelteKit + Svelte 5)
- ✅ Authentication flow with magic links
- ✅ Dashboard with 6 interactive charts (ECharts)
- ✅ Trade management (list, create, edit, delete)
- ✅ CSV import with DAS Trader support
- ✅ Journal entry system with tabs
- ✅ Rule-based trading system
- ✅ Emotional state tracking
- ✅ Screenshot uploads with lightbox
- ✅ Voice note recording
- ✅ Settings management
- ✅ Real-time WebSocket notifications
- ✅ Toast notification system
- ✅ Notification center
- ✅ Modern macOS-inspired design
- ✅ Dark mode support
- ✅ Responsive mobile layout

#### Backend (Go + PostgreSQL)
- ✅ RESTful API with Chi router
- ✅ JWT authentication
- ✅ Magic link authentication via Resend
- ✅ PostgreSQL database with migrations
- ✅ Trade CRUD operations
- ✅ Journal entry management
- ✅ Rule set management
- ✅ Tag system
- ✅ Analytics endpoints
- ✅ WebSocket notification system
- ✅ Email service integration
- ✅ CSV import processing
- ✅ Server-side pagination
- ✅ Advanced filtering

### Technical Improvements
- Migrated to Svelte 5 with Runes API
- Implemented server-side pagination for scalability
- Added comprehensive error handling
- Improved type safety with TypeScript
- Optimized database queries
- Added proper indexes for performance
- Implemented timezone support
- Created reusable component library

### Documentation
- ✅ Architecture documentation
- ✅ API specification
- ✅ Frontend component library
- ✅ Backend structure guide
- ✅ Authentication flow docs
- ✅ CSV import guide
- ✅ WebSocket setup guide
- ✅ Getting started guides

## [1.0.0] - 2024-12-01

### Initial Release
- Basic trade tracking
- Simple journal entries
- Basic analytics
- Manual trade entry
- User authentication
- Dashboard layout

---

## Upcoming Features

### High Priority
- [ ] User profile management
- [ ] Advanced analytics (correlation analysis)
- [ ] Export functionality (PDF reports)
- [ ] Multi-account support
- [ ] Trade performance metrics
- [ ] Custom date range selection
- [ ] Bulk edit trades
- [ ] Settings page with timezone UI

### Medium Priority
- [ ] Social sharing features
- [ ] Trade templates
- [ ] Automated trade sync with brokers
- [ ] Mobile app (React Native)
- [ ] Advanced charting
- [ ] Custom indicators
- [ ] Backtesting support

### Low Priority
- [ ] Community features
- [ ] Trading education content
- [ ] AI-powered insights
- [ ] Integration with trading platforms
- [ ] API for third-party integrations

---

## Version History

| Version | Date | Key Changes |
|---------|------|-------------|
| 2.0.0 | 2025-01-17 | Complete rewrite with Svelte 5, server-side pagination, timezone support |
| 1.0.0 | 2024-12-01 | Initial release |

---

## Breaking Changes

### 2.0.0
- **API Response Format**: List endpoints now return pagination metadata when `limit` is specified
- **Frontend**: Migrated to Svelte 5 Runes (incompatible with Svelte 4)
- **Database**: Added new columns for strategy, stop_loss, take_profit
- **Authentication**: Changed to JWT-only (removed session-based auth)

---

## Migration Guides

### Upgrading to 2.0.0

**Frontend:**
1. Update to Svelte 5: `npm install svelte@5`
2. Update API client calls to handle paginated responses
3. Update settings store import for timezone support

**Backend:**
1. Run new database migrations
2. Update API calls to include pagination parameters
3. Configure timezone settings in environment

**Database:**
```sql
-- Add new columns
ALTER TABLE trades ADD COLUMN strategy VARCHAR(255);
ALTER TABLE trades ADD COLUMN stop_loss DECIMAL(10,2);
ALTER TABLE trades ADD COLUMN take_profit DECIMAL(10,2);

-- Add indexes for performance
CREATE INDEX idx_trades_strategy ON trades(strategy);
CREATE INDEX idx_trades_pnl ON trades(pnl);
CREATE INDEX idx_trades_opened_at ON trades(opened_at);
```

---

For detailed information about specific features, see the main documentation in `/docs/`.
