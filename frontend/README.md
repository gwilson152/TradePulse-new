# TradePulse Frontend

A modern trading journal application built with Svelte 5, SvelteKit 2, and a custom macOS-inspired design system featuring glassmorphism.

## 🎨 Design Philosophy

TradePulse features a unique **macOS-inspired interface** with:
- **Bottom Dock Navigation** - 5 color-coded sections in a floating glassmorphism dock
- **Top Menu Bar** - 44px bar with app branding, live clock, and notifications
- **Glassmorphism Cards** - Semi-transparent cards with backdrop blur effects
- **Single-Color Outline Icons** - Material Design Icons with section-specific colors
- **Full Accessibility** - WCAG AA compliant with zero compiler warnings

## 🚀 Quick Start

### Prerequisites

- Node.js 18+ or Bun
- npm, pnpm, yarn, or bun

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd TradePulse/frontend

# Install dependencies
npm install
# or
bun install

# Start development server
npm run dev
# or
bun run dev

# Open browser
# Navigate to http://localhost:5173
```

### Build for Production

```bash
# Build the application
npm run build

# Preview production build
npm run preview
```

## 📁 Project Structure

```
frontend/
├── src/
│   ├── routes/                    # SvelteKit routes
│   │   ├── app/                   # Protected app routes
│   │   │   ├── dashboard/         # Dashboard page
│   │   │   ├── trades/            # Trades management
│   │   │   ├── journal/           # Trading journal
│   │   │   ├── analytics/         # Analytics & charts
│   │   │   ├── settings/          # Settings & rules
│   │   │   └── +layout.svelte     # macOS-inspired layout
│   │   ├── auth/                  # Authentication pages
│   │   └── +page.svelte           # Landing page
│   │
│   ├── lib/
│   │   ├── components/
│   │   │   ├── ui/                # 13 Base UI components
│   │   │   ├── trading/           # 9 Trading components
│   │   │   └── notifications/     # Notification components
│   │   ├── api/
│   │   │   └── client.ts          # API client with all endpoints
│   │   ├── types/
│   │   │   └── index.ts           # TypeScript type definitions
│   │   └── stores/                # Svelte stores
│   │
│   ├── app.css                    # Global styles & Tailwind imports
│   └── app.html                   # HTML template
│
├── static/                        # Static assets
├── tailwind.config.js             # Tailwind configuration
├── svelte.config.js               # SvelteKit configuration
├── vite.config.ts                 # Vite configuration
└── package.json
```

## 🎨 Design System

### Color-Coded Sections

Each section of the app has a unique color for visual organization:

| Section | Icon | Color | Gradient |
|---------|------|-------|----------|
| **Overview** | `mdi:view-dashboard-outline` | Blue 500 | `from-blue-500 to-cyan-500` |
| **Trades** | `mdi:chart-line-variant` | Emerald 500 | `from-emerald-500 to-teal-500` |
| **Journal** | `mdi:book-open-page-variant-outline` | Purple 500 | `from-purple-500 to-pink-500` |
| **Analytics** | `mdi:chart-box-outline` | Orange 500 | `from-amber-500 to-orange-500` |
| **Settings** | `mdi:cog-outline` | Slate 500 | N/A |

### Glassmorphism Pattern

Our signature glass effect for cards:

```css
bg-white/60 dark:bg-slate-800/60
backdrop-blur-xl
border border-slate-200/50 dark:border-slate-700/50
```

### Component Variants

#### Card
- **glass** - Semi-transparent with backdrop blur (default)
- **solid** - Traditional solid background
- **elevated** - With shadow elevation

#### Button
- **filled** - Solid color with shadow
- **gradient** - Gradient background with glow
- **soft** - Light background with colored text
- **ghost** - Transparent with hover state

### Icon Guidelines

- **Use outline-style icons only** (`-outline` suffix)
- **Single color per icon** (no multi-color icons)
- **24px size** for navigation
- **Color-code by section**

Example:
```svelte
<Icon icon="mdi:view-dashboard-outline" width="24" class="text-blue-500" />
```

## 🧩 Component Library

### Base UI Components (13)

| Component | Description | Accessibility |
|-----------|-------------|---------------|
| `Button.svelte` | 4 variants × 6 colors | ✅ Full keyboard support |
| `Card.svelte` | 3 variants with glassmorphism | ✅ Semantic HTML |
| `Modal.svelte` | Self-contained dialog | ✅ ARIA attributes, Escape key |
| `Input.svelte` | Form input with validation | ✅ Label association |
| `Select.svelte` | Custom dropdown | ✅ Label association |
| `Textarea.svelte` | Multi-line input | ✅ Label association, char counter |
| `FileUpload.svelte` | Drag-drop with preview | ✅ Keyboard support, ARIA role |
| `ImageGallery.svelte` | Lightbox viewer | ✅ Keyboard navigation |
| `AudioRecorder.svelte` | Voice note recorder | ✅ Visual feedback |
| `AudioPlayer.svelte` | Audio playback | ✅ Native controls |
| `Badge.svelte` | Status indicators | ✅ Semantic colors |
| `ChartCard.svelte` | ECharts wrapper | ✅ Responsive resize |
| `MetricCard.svelte` | Dashboard stats | ✅ Semantic structure |

### Trading Components (9)

| Component | Description | Features |
|-----------|-------------|----------|
| `TradeModal.svelte` | Create/edit trades | Cost basis method selection |
| `PositionTimeline.svelte` | Entry/exit timeline | Visual chronology |
| `PositionSizeBar.svelte` | Position composition | Progress bar visualization |
| `AddToPositionModal.svelte` | Add entries/exits | Entry type selection |
| `RuleCard.svelte` | Rule display | Phase & category badges |
| `RuleAdherenceInput.svelte` | 5-level scoring | Traffic light colors |
| `AdherenceScoreDisplay.svelte` | Score widget | Weighted average |
| `JournalEntryModal.svelte` | 4-tab journal | Emotions, rules, media |
| `TradeDetailSlideOver.svelte` | Trade details panel | Slide-over animation |

### Using Components

```svelte
<script>
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
</script>

<!-- Glassmorphism card -->
<Card variant="glass" padding="lg" hover={true}>
  <h2>Total P&L</h2>
  <p class="text-3xl font-bold text-emerald-600">$5,234.56</p>

  <Button variant="gradient" color="primary" size="sm">
    View Details
  </Button>
</Card>
```

## 🔌 API Integration

The API client is located at `src/lib/api/client.ts` and provides typed methods for all backend endpoints.

### Example Usage

```typescript
import { apiClient } from '$lib/api/client';

// Get trades
const trades = await apiClient.getTrades({ limit: 20, offset: 0 });

// Create trade
const newTrade = await apiClient.createTrade({
  symbol: 'AAPL',
  direction: 'LONG',
  quantity: 100,
  entry_price: 150.25,
  cost_basis_method: 'FIFO'
});

// Add position entry
await apiClient.addEntry(tradeId, {
  price: 148.50,
  quantity: 50,
  timestamp: new Date().toISOString(),
  notes: 'Adding to winner'
});

// Create journal entry
await apiClient.createJournalEntry({
  trade_id: tradeId,
  content: 'Followed my rules perfectly',
  emotional_state: {
    confidence: 8,
    stress: 3,
    discipline: 9
  }
});
```

### Available Methods

**Trades:**
- `getTrades(params)` - List trades with filtering
- `getTradeById(id)` - Get single trade
- `createTrade(data)` - Create new trade
- `updateTrade(id, data)` - Update trade

**Position Lifecycle:**
- `addEntry(tradeId, data)` - Add entry to position
- `addExit(tradeId, data)` - Add exit from position

**Rule Sets:**
- `getRuleSets()` - List all rule sets
- `createRuleSet(data)` - Create rule set
- `updateRuleSet(id, data)` - Update rule set
- `deleteRuleSet(id)` - Delete rule set
- `addRule(ruleSetId, data)` - Add rule to set
- `updateRule(ruleId, data)` - Update rule
- `deleteRule(ruleId)` - Delete rule

**Journal:**
- `getJournalEntries(params)` - List journal entries
- `createJournalEntry(data)` - Create entry

**File Upload:**
- `uploadFile(file)` - Upload screenshot or audio file

## 📊 Features

### Position Lifecycle Tracking

Track multiple entries and exits for a single position with three cost basis methods:

- **FIFO** (First In, First Out)
- **LIFO** (Last In, First Out)
- **AVERAGE** (Average Cost)

```typescript
interface Entry {
  id: string;
  price: number;
  quantity: number;
  timestamp: string;
  notes?: string;
  fees?: number;
}

interface Exit {
  id: string;
  price: number;
  quantity: number;
  timestamp: string;
  notes?: string;
  fees?: number;
  pnl: number;
}
```

### Rule Adherence Tracking

Traffic light scoring system with 5 levels:

| Score | Label | Color | Use Case |
|-------|-------|-------|----------|
| 100% | Perfect | Emerald | Rule fully followed |
| 75% | Good | Emerald | Minor deviation |
| 50% | Partial | Amber | Significant deviation |
| 25% | Poor | Red | Major violation |
| 0% | Failed | Red | Complete violation |

**Rule Phases:**
- `PRE_TRADE` - Before entering
- `DURING_TRADE` - While position is open
- `POST_TRADE` - After closing

**Rule Categories:**
- RISK_MANAGEMENT
- ENTRY
- EXIT
- POSITION_SIZING
- TIMING
- PSYCHOLOGY
- GENERAL

### Comprehensive Journaling

4-tab journal entry system:

1. **Reflection** - Rich text content area
2. **Emotions** - Confidence, stress, discipline sliders (1-10)
3. **Rules** - Adherence scoring for all active rules
4. **Media** - Screenshots (PNG, JPG) + voice notes (MP3, WAV)

```typescript
interface EmotionalState {
  confidence: number; // 1-10
  stress: number; // 1-10
  discipline: number; // 1-10
  notes: string;
}

interface JournalEntry {
  trade_id?: string;
  content: string;
  emotional_state?: EmotionalState;
  rule_adherence?: RuleAdherence[];
  screenshots?: string[];
  voice_notes?: string[];
}
```

### Interactive Analytics

6 ECharts visualizations with time range selector:

1. **Cumulative P&L Line Chart** - P&L over time
2. **Win Rate by Day Bar Chart** - Daily win percentages
3. **Trade Outcomes Pie Chart** - Win/Loss distribution
4. **Avg Win vs Avg Loss** - Comparison bars
5. **Rule Adherence Correlation** - Scatter plot (adherence vs P&L)
6. **Emotional State Analysis** - Scatter plot (confidence/discipline vs P&L)

## ♿ Accessibility

All components are **WCAG AA compliant** with:

- ✅ Proper ARIA labels and roles
- ✅ Full keyboard navigation
- ✅ Screen reader compatibility
- ✅ Form label associations
- ✅ Focus management
- ✅ Color contrast ratios
- ✅ Skip links and landmarks

### Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Escape` | Close modal/lightbox |
| `Enter` | Activate file upload zone |
| `Tab` | Navigate through interactive elements |
| `←/→` | Navigate images in gallery |

## 🌙 Dark Mode

Full dark mode support with automatic detection:

```typescript
// Manual toggle (future feature)
const toggleDarkMode = () => {
  document.documentElement.classList.toggle('dark');
};
```

Dark mode uses Tailwind's `dark:` prefix:

```svelte
<div class="bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100">
  <!-- Content -->
</div>
```

## 🧪 Testing

```bash
# Run unit tests (future)
npm run test

# Run e2e tests (future)
npm run test:e2e

# Type checking
npm run check

# Linting
npm run lint
```

## 📦 Dependencies

### Core
- `svelte@^5.0.0` - Reactive UI framework
- `@sveltejs/kit@^2.0.0` - Full-stack framework
- `tailwindcss@^3.4.18` - Utility-first CSS

### UI & Icons
- `@iconify/svelte@^5.1.0` - Material Design Icons
- `echarts@^5.x` - Interactive charts

### Development
- `vite` - Build tool
- `typescript` - Type safety
- `autoprefixer` - CSS vendor prefixes
- `postcss` - CSS processing

## 🛠 Configuration

### Tailwind Configuration

Custom color palette defined in `tailwind.config.js`:

```javascript
module.exports = {
  theme: {
    extend: {
      colors: {
        profit: { /* Emerald shades */ },
        loss: { /* Red shades */ },
        surface: { /* Slate shades */ }
      }
    }
  }
}
```

### SvelteKit Configuration

Adapter and prerendering settings in `svelte.config.js`:

```javascript
import adapter from '@sveltejs/adapter-auto';

export default {
  kit: {
    adapter: adapter()
  }
};
```

## 📖 Documentation

- **[Project Overview](../docs/project.md)** - Complete project documentation (1047 lines)
- **[Design System Guide](../docs/frontend/design-system.md)** - Design patterns and guidelines (400+ lines)
- **[Component Library](../docs/frontend/component-library.md)** - All 22 components documented (650+ lines)
- **[Implementation Status](../docs/frontend/implementation-status.md)** - Feature completion tracking
- **[API Specification](../docs/api-spec.md)** - Backend API endpoints

## 🎯 Development Guidelines

### Code Style

- Use **Svelte 5 Runes** (`$state`, `$derived`, `$props`)
- **TypeScript** for all new code
- **Tailwind classes** for styling (no custom CSS unless necessary)
- **Functional components** with clear prop interfaces

### Component Structure

```svelte
<script lang="ts">
  import Icon from '@iconify/svelte';

  interface Props {
    variant?: 'glass' | 'solid' | 'elevated';
    children?: any;
  }

  let { variant = 'glass', children }: Props = $props();

  // Reactive state
  let isHovered = $state(false);

  // Derived values
  const classes = $derived(`card card-${variant}`);
</script>

<div class={classes} onmouseenter={() => isHovered = true}>
  {@render children?.()}
</div>

<style>
  /* Only if Tailwind is insufficient */
  .card { /* custom styles */ }
</style>
```

### Accessibility Checklist

Before committing:

- [ ] All form inputs have associated labels (`id`/`for` attributes)
- [ ] Click handlers have keyboard equivalents (`onkeydown`)
- [ ] Interactive elements have ARIA roles and labels
- [ ] Color contrast meets WCAG AA (4.5:1 for text)
- [ ] Focus states are visible
- [ ] Screen reader tested (if possible)

## 🚀 Deployment

### Environment Variables

Create `.env` file:

```bash
VITE_API_URL=http://localhost:8000/api
VITE_WS_URL=ws://localhost:8000/ws
```

### Production Build

```bash
# Build for production
npm run build

# Preview build locally
npm run preview

# Deploy to Vercel/Netlify/etc
# Follow platform-specific instructions
```

## 🤝 Contributing

1. Follow the existing code style
2. Use TypeScript for new components
3. Ensure accessibility compliance
4. Test in multiple browsers
5. Update documentation

## 📄 License

[Your License Here]

## 🔗 Links

- **Documentation:** [/docs](/docs)
- **Design System:** [design-system.md](../docs/frontend/design-system.md)
- **API Spec:** [api-spec.md](../docs/api-spec.md)
- **Issue Tracker:** [GitHub Issues]

---

**Built with ❤️ using Svelte 5 + SvelteKit 2**
