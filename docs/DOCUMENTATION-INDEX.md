# TradePulse Documentation Index

**Last Updated:** November 19, 2025

## 📚 Documentation Structure

### Core Documentation
- **[README](README.md)** - Main entry point with quick links
- **[Project Overview](project.md)** - Comprehensive project description, goals, and features
- **[Architecture](architecture.md)** - System architecture, tech stack, and design patterns
- **[CHANGELOG](CHANGELOG.md)** - Release history and version changes

### API Documentation
- **[API Specification](api-spec.md)** - Complete API endpoint reference
- **[Authentication](authentication.md)** - Auth flows (Google OAuth, password-based, JWT)

### Feature Documentation
- **[CSV Import Guide](csv-import.md)** - Platform-specific CSV import instructions (DAS Trader, PropReports)
- **[Chart Preview Feature](features/chart-preview-plan.md)** - ✅ Implemented - TradingView chart integration

### Backend Documentation (`backend/`)
- **[Getting Started](backend/getting-started.md)** - Setup, development, and build instructions
- **[Database Schema](backend/database.md)** - Complete schema with relationships
- **[Migrations](backend/migrations.md)** - Migration system and history
- **[Implementation Status](backend/implementation-status.md)** - ✅ **START HERE** - Current feature status and recent changes
- **[API Patterns](backend/api-patterns.md)** - Code patterns, error handling, validation
- **[Trades API](backend/trades-api.md)** - Detailed trades endpoint documentation
- **[Notifications](backend/notifications.md)** - WebSocket notification system
- **[Email System](backend/email.md)** - Email templates and sending

### Frontend Documentation (`frontend/`)
- **[Getting Started](frontend/getting-started.md)** - Setup and development
- **[Implementation Status](frontend/implementation-status.md)** - Frontend feature status
- **[Design System](frontend/design-system.md)** - Colors, typography, spacing
- **[Component Library](frontend/component-library.md)** - Reusable UI components
- **[Settings & Rules](frontend/settings-and-rules.md)** - Configuration pages

---

## 🗄️ Archived Files (November 19, 2025)

The following files have been moved to `docs/archive/` for historical reference:

### Archived Documentation
- `archive/api-endpoints.md` → Use `api-spec.md` and `backend/trades-api.md` instead
- `archive/structure.md` → Use `architecture.md` instead
- `archive/components-documentation.md` → Use `frontend/component-library.md` instead
- `archive/readme.md` → Use main `README.md` and `frontend/getting-started.md` instead

### Previously Deleted
- ~~`websocket-notifications.md`~~ → Consolidated into `backend/notifications.md` (November 18, 2025)

---

## 🚀 Quick Start Guides

### For Developers
1. Read [README](README.md) for overview
2. Check [Backend Getting Started](backend/getting-started.md) or [Frontend Getting Started](frontend/getting-started.md)
3. Review [Implementation Status](backend/implementation-status.md) for current state
4. Reference [Architecture](architecture.md) for system design

### For API Users
1. Start with [API Specification](api-spec.md)
2. Review [Authentication](authentication.md)
3. Check specific endpoints in [Trades API](backend/trades-api.md)

### For Contributors
1. Check [Implementation Status](backend/implementation-status.md) first
2. Review [API Patterns](backend/api-patterns.md) for coding standards
3. Check [Database Schema](backend/database.md) before modifying data structures
4. Update this index when adding new documentation

---

## 📝 Documentation Guidelines

### When to Update Documentation

**Always update:**
- `backend/implementation-status.md` - After completing any backend work
- `frontend/implementation-status.md` - After completing any frontend work
- `CHANGELOG.md` - For version releases

**Update as needed:**
- API documentation when adding/changing endpoints
- Schema documentation when running migrations
- Component documentation when creating reusable components

### Documentation Standards

1. **Date all updates** - Include "Last Updated" at the top
2. **Use status indicators** - ✅ Implemented, 🚧 In Progress, 📋 Planned, ❌ Not Planned
3. **Be concise** - Focus on "what" and "why", not detailed "how"
4. **Link between docs** - Reference related documentation
5. **Keep examples minimal** - Show patterns, not full implementations
6. **Remove outdated info** - Better to delete than maintain wrong information

### File Organization

```
docs/
├── README.md                    # Main entry point
├── DOCUMENTATION-INDEX.md       # This file (updated)
├── CLEANUP-SUMMARY.md          # Documentation cleanup history
├── project.md                   # Project overview
├── architecture.md              # System architecture
├── CHANGELOG.md                 # Release history
├── api-spec.md                  # API reference
├── authentication.md            # Auth documentation
├── csv-import.md               # CSV import guide
├── backend/
│   ├── getting-started.md      # Backend setup
│   ├── implementation-status.md # ✅ Current backend state (Session 4)
│   ├── database.md             # Schema reference
│   ├── migrations.md           # Migration guide
│   ├── api-patterns.md         # Code patterns
│   ├── trades-api.md           # Trades endpoints
│   ├── notifications.md        # WebSocket system
│   └── email.md                # Email system
├── frontend/
│   ├── getting-started.md      # Frontend setup
│   ├── implementation-status.md # Frontend state
│   ├── design-system.md        # Design tokens
│   ├── component-library.md    # UI components
│   └── settings-and-rules.md   # Config pages
├── features/
│   └── chart-preview-plan.md   # ✅ Chart implementation (updated)
└── archive/                     # 🗄️ Deprecated files (historical reference)
    ├── api-endpoints.md
    ├── structure.md
    ├── components-documentation.md
    └── readme.md
```

---

## 🔍 Finding Information

| Looking for... | Check here |
|----------------|------------|
| What's implemented? | `backend/implementation-status.md` or `frontend/implementation-status.md` |
| How to set up development? | `backend/getting-started.md` or `frontend/getting-started.md` |
| API endpoint details? | `api-spec.md` or `backend/trades-api.md` |
| Database tables? | `backend/database.md` |
| How to add a migration? | `backend/migrations.md` |
| Component usage? | `frontend/component-library.md` |
| Design tokens? | `frontend/design-system.md` |
| Authentication flow? | `authentication.md` |
| Import CSV files? | `csv-import.md` |
| Recent changes? | `backend/implementation-status.md` (top section) |
| System architecture? | `architecture.md` |

---

## 🎯 Most Important Documents

**Start here if you're new:**
1. [README](README.md) - Overview
2. [Backend Implementation Status](backend/implementation-status.md) - Current state
3. [Architecture](architecture.md) - System design

**Reference frequently:**
- [API Specification](api-spec.md) - Endpoint reference
- [Database Schema](backend/database.md) - Data structure
- [Component Library](frontend/component-library.md) - UI components
