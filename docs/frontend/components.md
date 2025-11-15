# Components

## Component System

TradePulse uses **Skeleton UI** for the base component library with custom wrappers for consistency, and **Lucide Icons** for professional, consistent iconography.

## Component Organization

```
src/lib/components/
├── ui/                    # Base UI wrappers
│   ├── Button.svelte
│   ├── Card.svelte
│   ├── Input.svelte
│   ├── Modal.svelte
│   └── Table.svelte
├── trading/              # Trading-specific
│   ├── PnLBadge.svelte
│   └── MetricCard.svelte
├── journal/              # Journal components
├── notifications/        # Notification system
└── layout/               # Layout components
```

## Base UI Components

### Button

```svelte
<script>
  import Button from '$lib/components/ui/Button.svelte';
</script>

<Button variant="filled" color="primary" size="md">
  Click me
</Button>

<Button variant="ghost" color="secondary">
  Ghost button
</Button>
```

**Props:**
- `variant`: 'filled' | 'ghost' | 'soft'
- `color`: 'primary' | 'secondary' | 'success' | 'warning' | 'error'
- `size`: 'sm' | 'md' | 'lg'
- `disabled`: boolean
- `type`: 'button' | 'submit' | 'reset'

### Card

```svelte
<script>
  import Card from '$lib/components/ui/Card.svelte';
</script>

<Card padding="md" hover={true}>
  Card content here
</Card>
```

**Props:**
- `padding`: 'none' | 'sm' | 'md' | 'lg'
- `hover`: boolean (adds hover effect)

### Input

```svelte
<script>
  import Input from '$lib/components/ui/Input.svelte';

  let email = $state('');
</script>

<Input
  type="email"
  bind:value={email}
  label="Email Address"
  placeholder="you@example.com"
  required={true}
/>
```

**Props:**
- `type`: string
- `value`: string | number (bindable)
- `label`: string
- `placeholder`: string
- `error`: string (shows error message)
- `required`: boolean
- `disabled`: boolean

### Table

```svelte
<script>
  import Table from '$lib/components/ui/Table.svelte';
</script>

<Table hover={true} interactive={false}>
  <thead>
    <tr>
      <th>Symbol</th>
      <th>P&L</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>AAPL</td>
      <td>+$125.50</td>
    </tr>
  </tbody>
</Table>
```

**Props:**
- `hover`: boolean (row hover effect)
- `interactive`: boolean (clickable rows)

### Modal

```svelte
<script>
  import Modal from '$lib/components/ui/Modal.svelte';
</script>

<Modal title="Confirm Action">
  <p>Are you sure you want to delete this trade?</p>
  <div class="flex gap-2 mt-4">
    <Button>Cancel</Button>
    <Button color="error">Delete</Button>
  </div>
</Modal>
```

## Trading Components

### PnLBadge

Displays profit/loss with color coding.

```svelte
<script>
  import PnLBadge from '$lib/components/trading/PnLBadge.svelte';
</script>

<PnLBadge value={125.50} showSign={true} size="md" />
<PnLBadge value={-45.20} size="sm" />
```

**Props:**
- `value`: number (P&L amount)
- `showSign`: boolean (show + for positive)
- `size`: 'sm' | 'md' | 'lg'

**Behavior:**
- Green for positive values
- Red for negative values
- Auto-formats with 2 decimal places

### MetricCard

Dashboard metric card with optional trend indicator.

```svelte
<script>
  import MetricCard from '$lib/components/trading/MetricCard.svelte';
</script>

<MetricCard
  title="Total P&L"
  value="$1,234.56"
  subtitle="+12.5% this month"
  trend="up"
  icon="💰"
/>
```

**Props:**
- `title`: string (metric name)
- `value`: string | number (main value)
- `subtitle`: string (optional description)
- `trend`: 'up' | 'down' | 'neutral' (colors subtitle)
- `icon`: ComponentType (Lucide icon component)

**Example:**
```svelte
<script>
  import { DollarSign } from 'lucide-svelte';
  import MetricCard from '$lib/components/trading/MetricCard.svelte';
</script>

<MetricCard
  title="Total P&L"
  value="$1,234.56"
  subtitle="+12.5% this month"
  trend="up"
  icon={DollarSign}
/>
```

## Icons

TradePulse uses **Lucide Icons** - a professional, consistent icon set with CSS styling support.

### Usage

```svelte
<script>
  import { TrendingUp, DollarSign, Calendar, AlertCircle } from 'lucide-svelte';
</script>

<!-- Basic icon -->
<TrendingUp />

<!-- Custom size and stroke -->
<DollarSign size={24} strokeWidth={2} />

<!-- With CSS classes -->
<Calendar class="text-profit-600" />
<AlertCircle class="text-loss-500" size={20} />
```

### Common Icons for Trading

```svelte
import {
  TrendingUp,      // Profit, positive trend
  TrendingDown,    // Loss, negative trend
  DollarSign,      // Money, P&L
  Calendar,        // Dates, schedules
  BarChart3,       // Metrics, analytics
  FileText,        // Journal, notes
  Upload,          // CSV import
  Download,        // Export
  Settings,        // Configuration
  Bell,            // Notifications
  Check,           // Success
  X,               // Error, close
  Info,            // Information
  AlertCircle,     // Warning
  Trash2,          // Delete
  Edit,            // Edit
  Plus,            // Add/Create
  Filter           // Filtering
} from 'lucide-svelte';
```

### Styling Icons

Icons inherit color from parent text color:

```svelte
<div class="text-profit-600">
  <TrendingUp size={20} />
</div>

<div class="text-loss-500 dark:text-loss-400">
  <TrendingDown size={20} />
</div>
```

## Theme & Styling

### Dark Mode

Skeleton UI includes built-in dark mode support.

```svelte
<script>
  import { modeCurrent } from '@skeletonlabs/skeleton';

  function toggleDarkMode() {
    modeCurrent.set($modeCurrent === 'dark' ? 'light' : 'dark');
  }
</script>

<button onclick={toggleDarkMode}>
  Toggle Dark Mode
</button>
```

### Trading Colors

Custom colors defined in `tailwind.config.js`:

```css
/* Profit colors */
text-profit-500, bg-profit-500
text-profit-600, bg-profit-600

/* Loss colors */
text-loss-500, bg-loss-500
text-loss-600, bg-loss-600
```

### Utility Classes

```svelte
<!-- Profit/Loss text -->
<span class="text-profit">+$125.50</span>
<span class="text-loss">-$45.20</span>

<!-- Profit/Loss background -->
<div class="bg-profit text-white p-4">Winning trade</div>
<div class="bg-loss text-white p-4">Losing trade</div>
```

## Component Patterns

### Using Svelte 5 Runes

```svelte
<script lang="ts">
  // Props
  let { value, onChange } = $props<Props>();

  // State
  let isLoading = $state(false);

  // Derived
  let displayValue = $derived(value.toFixed(2));

  // Effect
  $effect(() => {
    console.log('Value changed:', value);
  });
</script>
```

### Event Handlers

```svelte
<script lang="ts">
  function handleClick(event: MouseEvent) {
    console.log('Clicked!', event);
  }
</script>

<button onclick={handleClick}>Click me</button>
```

### Two-way Binding

```svelte
<script lang="ts">
  let value = $state('');
</script>

<Input bind:value />
<p>You typed: {value}</p>
```

## Best Practices

1. **Use Skeleton UI classes** - Leverage built-in styles
2. **Wrap base components** - Create consistent wrappers in `ui/`
3. **Domain-specific components** - Put in `trading/` or `journal/`
4. **Keep components small** - Single responsibility
5. **Use TypeScript** - Type all props and events
6. **Follow Svelte 5 patterns** - Use runes ($state, $props, $derived)
7. **Responsive design** - Use Tailwind responsive classes
8. **Accessibility** - Include labels, aria attributes

## Adding New Components

1. **Create component file** - In appropriate directory
2. **Define props interface** - TypeScript for type safety
3. **Use Skeleton classes** - For consistency
4. **Export from index** - For easier imports
5. **Document usage** - Add to this file

## See Also

- [Skeleton UI Docs](https://www.skeleton.dev/)
- [Tailwind CSS Docs](https://tailwindcss.com/)
- [Svelte 5 Docs](https://svelte-5-preview.vercel.app/)
