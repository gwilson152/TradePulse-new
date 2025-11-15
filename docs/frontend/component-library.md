# TradePulse Frontend Component Library

## Overview

TradePulse uses a comprehensive component library built with Svelte 5, Tailwind CSS, and Skeleton UI. The library is organized into three main categories: UI components, trading-specific components, and page layouts.

## Component Architecture

### Design System
- **Framework**: Svelte 5.0 with runes ($state, $derived, $props)
- **Styling**: Tailwind CSS 3.4 + Skeleton UI 2.11
- **Icons**: Iconify with Material Design Icons
- **Charts**: Apache ECharts 5.x

### Color Palette
```css
/* Trading-specific colors */
--color-profit-50 through --color-profit-700  /* Green shades */
--color-loss-50 through --color-loss-700      /* Red shades */

/* Semantic colors */
--color-primary-*    /* Brand primary */
--color-success-*    /* Success states */
--color-warning-*    /* Warning states */
--color-error-*      /* Error states */
--color-surface-*    /* Backgrounds and surfaces */
```

## UI Components

### Base Components

#### Button.svelte
Multi-variant button component with loading states.

**Props:**
- `variant`: 'filled' | 'ghost' | 'soft'
- `color`: 'primary' | 'secondary' | 'success' | 'warning' | 'error'
- `size`: 'sm' | 'md' | 'lg'
- `disabled`: boolean
- `onclick`: function

**Example:**
```svelte
<Button color="primary" size="md" onclick={handleClick}>
  <Icon icon="mdi:save" class="mr-2" />
  Save Changes
</Button>
```

#### Input.svelte
Styled input field with label and validation support.

**Props:**
- `label`: string
- `type`: 'text' | 'number' | 'email' | 'password' | 'date' | 'time' | 'datetime-local'
- `value`: string (bindable)
- `placeholder`: string
- `required`: boolean
- `disabled`: boolean
- `error`: string
- `step`: string (for number inputs)
- `min`, `max`: string

#### Select.svelte
Dropdown select component with custom styling.

**Props:**
- `label`: string
- `options`: Array<{value: string, label: string, disabled?: boolean}>
- `value`: string (bindable)
- `placeholder`: string
- `required`: boolean
- `error`: string

**Example:**
```svelte
<Select
  label="Cost Basis Method"
  options={[
    { value: 'FIFO', label: 'First In First Out' },
    { value: 'LIFO', label: 'Last In First Out' }
  ]}
  bind:value={method}
/>
```

#### Textarea.svelte
Multi-line text input with character count.

**Props:**
- `label`: string
- `value`: string (bindable)
- `rows`: number (default: 4)
- `maxLength`: number
- `placeholder`: string
- `helperText`: string
- `error`: string
- `resize`: 'none' | 'vertical' | 'horizontal' | 'both'

#### Badge.svelte
Flexible badge component for status indicators and tags.

**Props:**
- `color`: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info' | 'neutral'
- `variant`: 'filled' | 'outlined' | 'soft'
- `size`: 'sm' | 'md' | 'lg'
- `icon`: string (Iconify icon name)
- `removable`: boolean
- `onRemove`: function

**Example:**
```svelte
<Badge color="success" variant="soft" size="sm">
  Open Position
</Badge>

<Badge color="primary" removable onRemove={() => removeTag(tag)}>
  {tag}
</Badge>
```

#### Card.svelte
Container component with hover effects.

**Props:**
- `padding`: 'sm' | 'md' | 'lg'
- `hover`: boolean

#### Modal.svelte
Full-screen modal with backdrop.

**Props:**
- `open`: boolean
- `title`: string
- `size`: 'sm' | 'md' | 'lg' | 'xl'
- `onClose`: function

### Media Components

#### FileUpload.svelte
Drag-and-drop file upload with preview.

**Props:**
- `accept`: string (MIME types, default: '*')
- `multiple`: boolean
- `maxSize`: number (in MB)
- `label`: string
- `preview`: boolean
- `value`: File[]
- `onChange`: (files: File[]) => void

**Features:**
- Drag-and-drop zone
- File size validation
- Preview generation for images
- Multi-file support
- Remove individual files

**Example:**
```svelte
<FileUpload
  accept="image/*"
  multiple={true}
  maxSize={10}
  label="Upload Screenshots"
  onChange={handleFiles}
/>
```

#### ImageGallery.svelte
Responsive image grid with lightbox viewer.

**Props:**
- `images`: string[] (URLs)
- `alt`: string
- `onDelete`: (index: number) => void (optional)

**Features:**
- Grid layout (responsive columns)
- Lightbox viewer with navigation
- Keyboard controls (Escape, Arrow keys)
- Delete functionality

#### AudioRecorder.svelte
Voice note recorder using MediaRecorder API.

**Props:**
- `onRecordingComplete`: (blob: Blob) => void
- `maxDuration`: number (seconds, default: 300)

**Features:**
- Record/pause/resume controls
- Visual timer with progress bar
- Animated recording indicator
- Duration limit enforcement
- Cancel functionality

#### AudioPlayer.svelte
Audio playback component with controls.

**Props:**
- `src`: string (audio URL)
- `label`: string
- `onDelete`: () => void (optional)

**Features:**
- Play/pause controls
- Seek bar with visual feedback
- Time display (current/total)
- Responsive design

### Advanced UI Components

#### ChartCard.svelte
Wrapper for ECharts visualizations.

**Props:**
- `title`: string
- `subtitle`: string
- `options`: EChartsOption
- `loading`: boolean
- `height`: string (default: '400px')

**Features:**
- Automatic resize handling
- Loading state overlay
- Reactive option updates
- Proper cleanup on unmount

**Example:**
```svelte
<ChartCard
  title="Cumulative P&L"
  subtitle="Your profit and loss over time"
  options={chartOptions()}
  loading={dataLoading}
  height="350px"
/>
```

## Trading Components

### Position Lifecycle Components

#### PositionTimeline.svelte
Visual timeline showing all entries and exits for a position.

**Props:**
- `entries`: Entry[]
- `exits`: Exit[]
- `tradeType`: 'LONG' | 'SHORT'

**Features:**
- Chronological ordering
- Color-coded events (blue for entries, purple for exits)
- P&L display for exits
- Notes and fees display
- Empty state handling

#### PositionSizeBar.svelte
Visual progress bar showing position size changes.

**Props:**
- `entries`: Entry[]
- `exits`: Exit[]
- `symbol`: string

**Features:**
- Visual bar showing remaining vs exited portions
- Entry/exit totals
- Current position size display
- Percentage remaining
- Active position indicator

#### AddToPositionModal.svelte
Modal for adding entries or exits to existing positions.

**Props:**
- `open`: boolean
- `trade`: Trade
- `action`: 'entry' | 'exit'
- `onClose`: () => void
- `onSubmit`: (data: Partial<Entry> | Partial<Exit>) => Promise<void>

**Features:**
- Price, quantity, timestamp, notes, fees inputs
- Estimated P&L calculation for exits
- Validation (prevent over-exiting)
- Current position size display

### Rule Management Components

#### RuleCard.svelte
Display card for individual trading rules.

**Props:**
- `rule`: Rule
- `editable`: boolean
- `onEdit`: () => void
- `onDelete`: () => void

**Features:**
- Phase badge (Pre-trade/During/Post-trade)
- Category badge
- Importance stars (1-5)
- Hover actions for edit/delete

#### RuleAdherenceInput.svelte
Interactive traffic light scoring component.

**Props:**
- `rule`: Rule
- `value`: RuleAdherence
- `onChange`: (adherence: RuleAdherence) => void

**Features:**
- 5-level scoring grid (Perfect/Good/Partial/Poor/Failed)
- Color-coded buttons
- Traffic light indicator
- Mandatory notes for scores <100%
- Auto-expand notes section

**Scoring Levels:**
- **Perfect (100%)**: Green - Rule fully followed
- **Good (75%)**: Green - Minor deviation
- **Partial (50%)**: Yellow - Significant deviation
- **Poor (25%)**: Red - Major violation
- **Failed (0%)**: Red - Complete violation

#### AdherenceScoreDisplay.svelte
Dashboard widget showing overall adherence score.

**Props:**
- `adherences`: RuleAdherence[]
- `rules`: Rule[]
- `showDetails`: boolean

**Features:**
- Large circular score display with traffic light
- Overall and weighted scores
- Phase-based breakdown (Pre/During/Post)
- Color-coded indicators

### Trading UI Components

#### MetricCard.svelte
Stat card for dashboard metrics.

**Props:**
- `title`: string
- `value`: string | number
- `icon`: string
- `color`: 'primary' | 'success' | 'error' | 'warning'
- `trend`: 'up' | 'down' | 'neutral'
- `trendValue`: string

#### PnLBadge.svelte
Color-coded P&L display.

**Props:**
- `value`: number
- `showSign`: boolean

**Features:**
- Auto color (green for profit, red for loss)
- Optional +/- sign prefix
- Currency formatting

#### TradeModal.svelte
Modal for opening new positions or editing trades.

**Props:**
- `isOpen`: boolean
- `trade`: Trade (optional, for editing)
- `onClose`: () => void
- `onSave`: (trade: Partial<Trade>) => void

**Features:**
- Symbol and trade type selection
- Entry details (price, quantity, date, time)
- Exit details (optional)
- Fees and cost basis method
- Tags with badge UI
- Notes textarea
- Validation and error handling

#### TradeDetailSlideOver.svelte
Slide-over panel with detailed position information.

**Props:**
- `isOpen`: boolean
- `trade`: Trade
- `onClose`: () => void
- `onEdit`: () => void
- `onDelete`: () => void
- `onAddEntry`: (data: Partial<Entry>) => Promise<void>
- `onAddExit`: (data: Partial<Exit>) => Promise<void>

**Features:**
- Two-tab interface (Overview + Timeline)
- Position size visualization
- Quick add entry/exit buttons
- Detailed metrics display
- Timeline view with PositionTimeline component

#### JournalEntryModal.svelte
Comprehensive 4-tab modal for journal entries.

**Props:**
- `open`: boolean
- `trade`: Trade (optional)
- `rules`: Rule[]
- `onClose`: () => void
- `onSubmit`: (data: Partial<JournalEntry>, screenshots: File[], voiceNotes: Blob[]) => Promise<void>

**Tabs:**
1. **Reflection**: Rich textarea for trade narratives
2. **Emotions**: Sliders for confidence, stress, discipline
3. **Rules**: Rule adherence scoring for all active rules
4. **Media**: Screenshot upload and voice recorder

**Features:**
- Tab completion indicators
- Form validation
- Adherence score display
- Trade summary (if linked)
- Error handling

## Page Layouts

### Dashboard (`/app/dashboard`)
- 4 metric cards (Total P&L, Win Rate, Avg Win, Active Trades)
- Recent trades table
- Quick stats

### Trades (`/app/trades`)
- Trade list with search and filters
- Add/edit/delete functionality
- Trade detail slide-over
- Position lifecycle support

### Journal (`/app/journal`)
- Journal entry cards (expandable)
- Rule adherence scores
- Emotional state visualization
- Media galleries
- Create/edit modals

### Analytics (`/app/analytics`)
- 6 interactive ECharts visualizations
- Time range selector
- Key metrics grid
- Correlation analysis

### Rules (`/app/settings/rules`)
- Rule set sidebar
- Rule management panel
- Template library
- CRUD operations

## TypeScript Types

### Core Types

```typescript
// Position Lifecycle
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

type CostBasisMethod = 'FIFO' | 'LIFO' | 'AVERAGE';

interface Trade {
  id: string;
  user_id: string;
  symbol: string;
  trade_type: 'LONG' | 'SHORT';
  // Position lifecycle fields
  entries: Entry[];
  exits: Exit[];
  current_position_size: number;
  average_entry_price: number;
  total_fees: number;
  realized_pnl: number;
  unrealized_pnl: number | null;
  cost_basis_method: CostBasisMethod;
  // Timestamps
  opened_at: string;
  closed_at: string | null;
  created_at: string;
  updated_at: string;
}

// Rule Set System
type RulePhase = 'PRE_TRADE' | 'DURING_TRADE' | 'POST_TRADE';
type RuleCategory = 'RISK_MANAGEMENT' | 'ENTRY' | 'EXIT' | 'POSITION_SIZING' | 'TIMING' | 'PSYCHOLOGY' | 'GENERAL';

interface Rule {
  id: string;
  title: string;
  description: string;
  weight: number; // 1-5
  phase: RulePhase;
  category: RuleCategory;
  created_at: string;
}

interface RuleSet {
  id: string;
  user_id: string;
  name: string;
  description: string;
  rules: Rule[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface RuleAdherence {
  rule_id: string;
  rule_title: string;
  score: number; // 0, 25, 50, 75, 100
  notes: string;
  timestamp: string;
}

// Journal System
interface EmotionalState {
  confidence: number; // 1-10
  stress: number; // 1-10
  discipline: number; // 1-10
  notes: string;
}

interface JournalEntry {
  id: string;
  user_id: string;
  trade_id: string | null;
  entry_date: string;
  content: string;
  emotional_state: EmotionalState | null;
  rule_adherence?: RuleAdherence[];
  adherence_score?: number;
  screenshots: string[];
  voice_notes: string[];
  created_at: string;
  updated_at: string;
}
```

## Best Practices

### Component Development

1. **Use Svelte 5 Runes**:
   ```svelte
   let value = $state(0);
   let doubled = $derived(value * 2);
   let { prop1, prop2 } = $props();
   ```

2. **Type Safety**:
   - Always define Props interface
   - Use strict TypeScript types
   - Import types from `$lib/types`

3. **Accessibility**:
   - Include ARIA labels
   - Keyboard navigation support
   - Focus management in modals

4. **Error Handling**:
   - Display error messages
   - Loading states
   - Validation feedback

5. **Responsive Design**:
   - Mobile-first approach
   - Use Tailwind responsive utilities
   - Test on multiple screen sizes

### Styling Guidelines

1. **Use Tailwind Utilities**:
   ```svelte
   <div class="flex items-center gap-4 p-6 rounded-lg">
   ```

2. **Component-specific Styles**:
   ```svelte
   <style>
     .custom-component {
       /* Only for truly custom needs */
     }
   </style>
   ```

3. **Dark Mode Support**:
   ```svelte
   <div class="bg-surface-50 dark:bg-surface-800">
   ```

4. **Color Usage**:
   - Profit/Loss: Use `text-profit-600` / `text-loss-600`
   - Status: Use semantic colors (success, warning, error)
   - Surfaces: Use surface-* scales for backgrounds

## Component Status

### Completed ✅
- All 22 components fully implemented
- Type-safe with TypeScript
- Responsive design
- Dark mode support
- Loading and error states
- Form validation

### Pending 🔄
- Toast notification component
- Loading skeleton screens
- DatePicker component (currently using native input)
- Pagination component

## Dependencies

```json
{
  "@skeletonlabs/skeleton": "^2.11.0",
  "@iconify/svelte": "^5.1.0",
  "echarts": "^5.x",
  "svelte": "^5.0.0",
  "@sveltejs/kit": "^2.0.0",
  "tailwindcss": "^3.4.18"
}
```
