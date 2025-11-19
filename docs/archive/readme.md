# TradePulse Frontend Documentation

Complete documentation for the TradePulse frontend application.

**Version:** 2.0
**Framework:** SvelteKit 2.0 + Svelte 5 (Runes)
**Last Updated:** January 2025

## 📚 Documentation Index

### Core Documentation
- [**Getting Started**](getting-started.md) - Setup, installation, and development
- [**Design System**](design-system.md) - macOS-inspired layout, glassmorphism, color system
- [**Component Library**](component-library.md) - Complete component reference (22 components)
- [**Implementation Status**](implementation-status.md) - Feature status and recent updates
- [**Settings & Rules**](settings-and-rules.md) - Rule-based trading system documentation

### Quick Reference

**Start Development:**
```bash
cd frontend
npm install
cp .env.example .env
npm run dev
```

**Access Application:**
- Development: http://localhost:5173
- Production: https://tradepulse.drivenw.com

## 🎨 Design System (v2.0)

TradePulse features a unique macOS-inspired design with glassmorphism effects:

- **Layout:** Bottom dock navigation + top menu bar (no sidebar)
- **Style:** Glassmorphism with backdrop blur effects
- **Icons:** Material Design Icons (outline style, color-coded by section)
- **Accessibility:** WCAG AA compliant with full keyboard navigation
- **Dark Mode:** Full dark mode support with proper contrast

See [Design System](design-system.md) for complete details.

## 🧩 Component Library

22 production-ready components built with Svelte 5:

**Base UI Components:**
- Button, Input, Select, Textarea, Card, Badge, Modal
- FileUpload, ImageGallery, AudioRecorder, AudioPlayer
- ChartCard (ECharts wrapper)

**Trading Components:**
- TradeModal, TradeDetailSlideOver, PnLBadge
- PositionTimeline, PositionSizeBar, AddToPositionModal
- JournalEntryModal, RuleCard, RuleAdherenceInput
- CSVImportSlideOver

See [Component Library](component-library.md) for props, variants, and examples.

## 🏗️ Technology Stack

- **Framework:** SvelteKit 2.0 with Svelte 5 Runes API
- **Styling:** Tailwind CSS 3.4
- **UI Base:** Skeleton UI 2.11 (minimal usage, mostly custom components)
- **Charts:** Apache ECharts 5.x
- **Icons:** Iconify with Material Design Icons
- **Build Tool:** Vite
- **Language:** TypeScript

## 📊 Recent Updates (January 2025)

### Server-Side Pagination & Filtering
- ✅ Database-level pagination for trades API
- ✅ Advanced filters (trade type, date range, strategy, P&L range)
- ✅ Live result counts and pagination controls
- ✅ Reactive filter system with state tracking

### Timezone Support
- ✅ User settings store with timezone preferences
- ✅ Market time vs local time display options
- ✅ Configurable date/time formats (12h/24h, short/medium/long)
- ✅ Common timezone presets (ET, CT, MT, PT, London, Tokyo, etc.)

### Trades Page Redesign
- ✅ Data table format for compact display
- ✅ Mouse-following tooltips with trade details
- ✅ Mobile long-press support
- ✅ Show date and time in all trade listings

See [Implementation Status](implementation-status.md) for complete changelog.

## 🎯 Key Features

### 1. Position Lifecycle Tracking
- Multiple entries/exits per position
- Three cost basis methods: FIFO, LIFO, Average
- Visual timeline of position events
- Realized vs unrealized P&L

### 2. Rule-Based Trading System
- Custom rule sets with weighted scoring
- Traffic light adherence system (5 levels)
- Phase-based organization (Pre-trade, During, Post-trade)
- Correlation analysis with performance

### 3. Comprehensive Journaling
- Rich text entries with emotional state tracking
- Screenshot uploads with lightbox viewer
- Voice note recording (up to 5 minutes)
- Rule adherence scoring per trade

### 4. Interactive Analytics
- 6 ECharts visualizations
- Time range selector (1W, 1M, 3M, 6M, 1Y, All)
- Real-time data updates
- Responsive charts with dark mode support

### 5. CSV Import
- DAS Trader and PropReports format support
- 3-step import workflow
- Validation and preview
- Bulk trade creation

## 📖 Additional Resources

- [Main Project Documentation](../README.md)
- [Architecture Overview](../architecture.md)
- [API Specification](../api-spec.md)
- [Backend Documentation](../backend/)

## 🚀 Development Status

**Current Phase:** Production Ready (v2.0)

- ✅ All core features implemented
- ✅ Full accessibility compliance (WCAG AA)
- ✅ Zero TypeScript errors
- ✅ Zero accessibility warnings
- ✅ Responsive design (mobile, tablet, desktop)
- ✅ Dark mode support
- ✅ Performance optimized

See [Implementation Status](implementation-status.md) for detailed progress tracking.
