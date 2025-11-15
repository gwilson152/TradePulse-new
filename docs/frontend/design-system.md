# TradePulse Design System

## Overview

TradePulse features a fresh, macOS-inspired design system with glassmorphism effects, vibrant gradients, and a unique dock-style navigation. The design prioritizes clarity, elegance, and modern aesthetics while maintaining excellent usability.

**Design Philosophy:**
- **Clean & Minimal** - Generous whitespace, clear hierarchy
- **Glassmorphism** - Translucent cards with backdrop blur
- **Vibrant Gradients** - Colorful, engaging visual elements
- **Smooth Interactions** - Subtle animations and transitions
- **Icon-First** - Colored icons for visual navigation and recognition

---

## Layout System

### macOS-Inspired Navigation

**Menu Bar** (Top, 44px height)
- Fixed position, backdrop blur
- Translucent background (white/80 dark:slate-900/80)
- App name with status indicator (green dot)
- Live clock display
- Notifications and user avatar

**Dock Navigation** (Bottom, floating)
- Fixed bottom center position
- Glassmorphism container with blur
- 5 main navigation items
- Colored icons for each section
- Active state indicators
- Hover scale animations

**Color-Coded Navigation:**
- Overview: Blue (`text-blue-500`)
- Trades: Emerald (`text-emerald-500`)
- Journal: Purple (`text-purple-500`)
- Analytics: Orange (`text-orange-500`)
- Settings: Slate (`text-slate-500`)

---

## Color System

### Primary Palette

**Background Gradients:**
```css
/* Light mode */
bg-gradient-to-br from-slate-50 via-slate-100 to-slate-200

/* Dark mode */
bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950
```

**Functional Colors:**

| Purpose | Light | Dark | Usage |
|---------|-------|------|-------|
| **Primary** | `blue-500` | `blue-400` | Main actions, links |
| **Success** | `emerald-500` | `emerald-400` | Positive P&L, wins |
| **Error** | `red-500` | `red-400` | Negative P&L, losses |
| **Warning** | `orange-500` | `orange-400` | Alerts, cautions |
| **Secondary** | `purple-500` | `purple-400` | Accents, highlights |
| **Neutral** | `slate-700` | `slate-300` | Text, general UI |

### Gradient Combinations

**Icon Backgrounds (with shadows):**
```css
/* Blue */
bg-gradient-to-br from-blue-500 to-cyan-500
shadow-lg shadow-blue-500/30

/* Emerald */
bg-gradient-to-br from-emerald-500 to-teal-500
shadow-lg shadow-emerald-500/30

/* Purple */
bg-gradient-to-br from-purple-500 to-pink-500
shadow-lg shadow-purple-500/30

/* Orange */
bg-gradient-to-br from-orange-500 to-yellow-500
shadow-lg shadow-orange-500/30
```

**Button Gradients:**
```css
bg-gradient-to-r from-blue-500 to-cyan-500 hover:from-blue-600 hover:to-cyan-600
```

### Text Gradients

**Hero Headers:**
```css
bg-gradient-to-r from-slate-800 to-slate-600 dark:from-slate-100 dark:to-slate-400 bg-clip-text text-transparent
```

---

## Typography

### Font Stack
```css
font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, sans-serif
```

### Scale

| Element | Size | Weight | Line Height |
|---------|------|--------|-------------|
| **Hero (h1)** | `text-4xl` (36px) | `font-bold` | Auto |
| **Heading (h2)** | `text-xl` (20px) | `font-bold` | Auto |
| **Subheading (h3)** | `text-lg` (18px) | `font-semibold` | Auto |
| **Body** | `text-sm` (14px) | `font-normal` | Auto |
| **Caption** | `text-xs` (12px) | `font-medium` | Auto |
| **Tiny** | `text-[10px]` | `font-medium` | Auto |

### Color Usage
- **Primary text:** `text-slate-800 dark:text-slate-200`
- **Secondary text:** `text-slate-600 dark:text-slate-400`
- **Muted text:** `text-slate-500 dark:text-slate-500`

---

## Components

### Card Component

**Variants:**

1. **Glass (Default)**
   ```svelte
   <Card variant="glass" padding="lg">
   ```
   - Background: `bg-white/60 dark:bg-slate-800/60`
   - Backdrop: `backdrop-blur-xl`
   - Border: `border border-slate-200/50 dark:border-slate-700/50`
   - Rounded: `rounded-2xl`

2. **Solid**
   ```svelte
   <Card variant="solid" padding="md">
   ```
   - Background: `bg-white dark:bg-slate-800`
   - Border: `border border-slate-200 dark:border-slate-700`

3. **Elevated**
   ```svelte
   <Card variant="elevated" padding="lg">
   ```
   - Background: `bg-white dark:bg-slate-800`
   - Shadow: `shadow-xl shadow-slate-200/50`

**Hover Effect:**
```svelte
<Card hover={true}>
```
- Scale: `hover:scale-[1.02]`
- Shadow: `hover:shadow-2xl`
- Border: `hover:border-slate-300`
- Transition: `transition-all duration-300`

**Padding Options:**
- `none`: No padding
- `sm`: `p-4`
- `md`: `p-6` (default)
- `lg`: `p-8`

---

### Button Component

**Variants:**

1. **Filled**
   ```svelte
   <Button variant="filled" color="primary">
   ```
   - Solid background with color
   - Box shadow matching color
   - Hover: Darker shade + larger shadow

2. **Gradient**
   ```svelte
   <Button variant="gradient" color="primary">
   ```
   - Gradient background (2 colors)
   - Large shadow with color glow
   - Hover: Intensified gradient + shadow

3. **Soft**
   ```svelte
   <Button variant="soft" color="primary">
   ```
   - Light background (color/10)
   - Colored text
   - No shadow

4. **Ghost**
   ```svelte
   <Button variant="ghost" color="primary">
   ```
   - Transparent background
   - Colored text
   - Hover: Light background

**Sizes:**
- `sm`: `text-xs px-3 py-1.5`
- `md`: `text-sm px-4 py-2.5` (default)
- `lg`: `text-base px-6 py-3`

**Colors:**
- `primary` - Blue
- `secondary` - Purple
- `success` - Emerald
- `warning` - Orange
- `error` - Red
- `neutral` - Slate

---

## Icon Strategy

### Icon Library
**Iconify with Material Design Icons** (`@iconify/svelte`)

### Icon Usage Rules

1. **Use Outline Icons**
   - Prefer `-outline` variants for cleaner look
   - Examples: `mdi:chart-line-variant`, `mdi:book-open-page-variant-outline`

2. **Single Color Per Icon**
   - NO multi-color icons
   - Use solid colors from the palette
   - Apply color via Tailwind classes

3. **Size Standards**
   - **Dock navigation:** 24px
   - **Card icons:** 24px (in 48px gradient circle)
   - **Button icons:** 16-20px depending on button size
   - **Decorative:** 64px

4. **Color Mapping**

   ```svelte
   <!-- Navigation -->
   <Icon icon="mdi:view-dashboard-outline" class="text-blue-500" />
   <Icon icon="mdi:chart-line-variant" class="text-emerald-500" />
   <Icon icon="mdi:book-open-page-variant-outline" class="text-purple-500" />

   <!-- Metric Cards -->
   <div class="bg-gradient-to-br from-blue-500 to-cyan-500">
     <Icon icon="mdi:cash-multiple" class="text-white" />
   </div>
   ```

5. **Icon Backgrounds**
   - Use gradient circles for metric cards
   - Size: `w-12 h-12` for cards
   - Rounded: `rounded-xl`
   - Shadow: Colored shadow matching gradient

---

## Spacing System

**Scale:**
```
1 = 0.25rem = 4px
2 = 0.5rem = 8px
3 = 0.75rem = 12px
4 = 1rem = 16px
6 = 1.5rem = 24px
8 = 2rem = 32px
12 = 3rem = 48px
```

**Common Patterns:**

- **Page padding:** `px-8 py-8`
- **Card spacing:** `space-y-8` between sections
- **Grid gaps:** `gap-6` for cards
- **Component gaps:** `gap-2` for related items, `gap-4` for buttons
- **Section spacing:** `mb-6` for section titles

---

## Effects & Animations

### Glassmorphism
```css
bg-white/60 dark:bg-slate-800/60
backdrop-blur-xl
border border-slate-200/50 dark:border-slate-700/50
```

### Shadows

**Subtle:**
```css
shadow-sm shadow-blue-500/30
```

**Medium:**
```css
shadow-lg shadow-blue-500/30
```

**Large:**
```css
shadow-xl shadow-slate-200/50
```

**Glow (hover):**
```css
hover:shadow-2xl hover:shadow-blue-500/40
```

### Transitions

**Standard:**
```css
transition-all duration-200
```

**Long:**
```css
transition-all duration-300
```

**Specific:**
```css
transition-transform duration-200
transition-colors
transition-shadow
```

### Hover Effects

**Cards:**
```css
hover:scale-[1.02] hover:shadow-2xl
```

**Buttons:**
```css
hover:bg-blue-600 hover:shadow-md
```

**Icons:**
```css
group-hover:scale-110
```

---

## Responsive Design

### Breakpoints
```css
sm: 640px   /* Mobile landscape */
md: 768px   /* Tablet */
lg: 1024px  /* Desktop */
xl: 1280px  /* Large desktop */
2xl: 1536px /* Extra large */
```

### Grid Patterns

**Metric Cards:**
```svelte
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
```

**Content Sections:**
```svelte
<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
  <div class="lg:col-span-2"><!-- Main content --></div>
  <div><!-- Sidebar --></div>
</div>
```

---

## Examples

### Metric Card Pattern
```svelte
<Card variant="glass" padding="lg" hover={true}>
  <div class="space-y-4">
    <!-- Icon & Badge -->
    <div class="flex items-center justify-between">
      <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center shadow-lg shadow-blue-500/30">
        <Icon icon="mdi:cash-multiple" width="24" class="text-white" />
      </div>
      <div class="px-3 py-1 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 text-xs font-semibold">
        +5.2K
      </div>
    </div>

    <!-- Label & Value -->
    <div>
      <p class="text-sm text-slate-600 dark:text-slate-400 mb-1">Total P&L</p>
      <p class="text-3xl font-bold text-emerald-600 dark:text-emerald-400">
        $12,450
      </p>
    </div>
  </div>
</Card>
```

### Action Button Group
```svelte
<div class="space-y-3">
  <Button variant="soft" color="primary" size="md" class="w-full justify-start">
    <Icon icon="mdi:plus-circle-outline" width="20" />
    New Trade
  </Button>
  <Button variant="soft" color="secondary" size="md" class="w-full justify-start">
    <Icon icon="mdi:notebook-outline" width="20" />
    Journal Entry
  </Button>
</div>
```

### Page Header
```svelte
<div class="flex items-center justify-between">
  <div>
    <h1 class="text-4xl font-bold bg-gradient-to-r from-slate-800 to-slate-600 dark:from-slate-100 dark:to-slate-400 bg-clip-text text-transparent mb-2">
      Page Title
    </h1>
    <p class="text-slate-600 dark:text-slate-400">Subtitle description</p>
  </div>
  <Button variant="gradient" color="primary" size="lg">
    <Icon icon="mdi:plus" width="20" />
    Action
  </Button>
</div>
```

---

## Dark Mode

### Strategy
- Automatic with `dark:` prefix
- Inverted luminosity, not just darker colors
- Maintain contrast ratios
- Use opacity for depth

### Color Adjustments

| Light | Dark |
|-------|------|
| `slate-50` | `slate-950` |
| `slate-100` | `slate-900` |
| `slate-200` | `slate-800` |
| `slate-600` | `slate-400` |
| `slate-800` | `slate-200` |
| `white/60` | `slate-800/60` |

### Testing
- Test all components in both modes
- Ensure text contrast is WCAG AA compliant
- Verify gradient visibility

---

## Accessibility

### Contrast
- Minimum 4.5:1 for normal text
- Minimum 3:1 for large text
- Use contrast checker for all color combinations

### Focus States
- Visible focus rings on all interactive elements
- Use `focus:ring-2 focus:ring-blue-500 focus:ring-offset-2`

### Keyboard Navigation
- All interactive elements accessible via tab
- Logical tab order
- Escape key closes modals

### ARIA
- Proper labels on icons
- Screen reader text where needed
- Semantic HTML structure

---

## Best Practices

### Do's ✅
- Use outline icons for consistency
- Apply single colors to icons (no gradients on icons themselves)
- Use glassmorphism for main content cards
- Apply colored shadows to gradient backgrounds
- Use subtle animations (200-300ms)
- Maintain generous whitespace
- Use gradient text for major headings
- Keep dock navigation visible

### Don'ts ❌
- Don't use multi-color icons
- Don't apply gradients directly to text (use bg-clip instead)
- Don't overuse animations
- Don't create deep nesting (max 3 levels)
- Don't use colors without considering dark mode
- Don't forget hover states
- Don't hide the dock on any page

---

## Implementation Checklist

When creating a new component or page:

- [ ] Use appropriate Card variant (glass for main content)
- [ ] Add hover effects where interactive
- [ ] Include dark mode styles
- [ ] Use outline icons with single colors
- [ ] Apply proper spacing (gap-6 for cards)
- [ ] Add gradient backgrounds to icon containers
- [ ] Include colored shadows on elevated elements
- [ ] Test responsive breakpoints
- [ ] Verify keyboard accessibility
- [ ] Check contrast ratios

---

## Resources

- **Icons:** [Iconify - Material Design Icons](https://icon-sets.iconify.design/mdi/)
- **Colors:** [Tailwind Colors](https://tailwindcss.com/docs/customizing-colors)
- **Gradients:** [UI Gradients](https://uigradients.com/)
- **Glassmorphism:** [CSS Glass](https://css.glass/)

---

## Complete Component Examples

### Dashboard Metric Grid

```svelte
<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Icon from '@iconify/svelte';

  const metrics = [
    { label: 'Total P&L', value: '$12,450', change: '+5.2K', icon: 'mdi:cash-multiple', gradient: 'from-blue-500 to-cyan-500', color: 'blue' },
    { label: 'Win Rate', value: '68%', change: '+3%', icon: 'mdi:chart-line', gradient: 'from-emerald-500 to-teal-500', color: 'emerald' },
    { label: 'Total Trades', value: '142', change: '+12', icon: 'mdi:chart-box', gradient: 'from-purple-500 to-pink-500', color: 'purple' },
    { label: 'Avg Trade', value: '$87.68', change: '-$2.10', icon: 'mdi:trending-up', gradient: 'from-orange-500 to-yellow-500', color: 'orange' }
  ];
</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
  {#each metrics as metric}
    <Card variant="glass" padding="lg" hover={true}>
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br {metric.gradient} flex items-center justify-center shadow-lg shadow-{metric.color}-500/30">
            <Icon icon={metric.icon} width="24" class="text-white" />
          </div>
          <div class="px-3 py-1 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 text-xs font-semibold">
            {metric.change}
          </div>
        </div>
        <div>
          <p class="text-sm text-slate-600 dark:text-slate-400 mb-1">{metric.label}</p>
          <p class="text-3xl font-bold text-emerald-600 dark:text-emerald-400">{metric.value}</p>
        </div>
      </div>
    </Card>
  {/each}
</div>
```

### Trade List with Glassmorphism

```svelte
<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import PnLBadge from '$lib/components/ui/PnLBadge.svelte';

  const trades = [
    { symbol: 'AAPL', direction: 'LONG', pnl: 245.50, status: 'CLOSED', date: '2025-01-12' },
    { symbol: 'TSLA', direction: 'SHORT', pnl: -125.30, status: 'CLOSED', date: '2025-01-12' },
    { symbol: 'NVDA', direction: 'LONG', pnl: 0, status: 'OPEN', date: '2025-01-13' }
  ];
</script>

<Card variant="glass" padding="none">
  <div class="overflow-hidden">
    {#each trades as trade, index}
      <div class="flex items-center justify-between p-4 border-b border-slate-200/50 dark:border-slate-700/50 last:border-b-0 hover:bg-white/30 dark:hover:bg-slate-700/30 transition-colors">
        <div class="flex items-center gap-4">
          <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center shadow-sm shadow-blue-500/30">
            <span class="text-white font-bold text-sm">{trade.symbol.charAt(0)}</span>
          </div>
          <div>
            <p class="font-semibold text-slate-900 dark:text-slate-100">{trade.symbol}</p>
            <p class="text-xs text-slate-600 dark:text-slate-400">{trade.date}</p>
          </div>
        </div>

        <div class="flex items-center gap-4">
          <Badge variant={trade.direction === 'LONG' ? 'success' : 'error'}>
            {trade.direction}
          </Badge>
          <Badge variant={trade.status === 'OPEN' ? 'warning' : 'neutral'}>
            {trade.status}
          </Badge>
          {#if trade.status === 'CLOSED'}
            <PnLBadge value={trade.pnl} />
          {/if}
        </div>
      </div>
    {/each}
  </div>
</Card>
```

### Modal with Glassmorphism

```svelte
<script lang="ts">
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Select from '$lib/components/ui/Select.svelte';

  let isOpen = $state(false);
  let symbol = $state('');
  let direction = $state('');
</script>

<Button variant="gradient" color="primary" onclick={() => isOpen = true}>
  New Trade
</Button>

<Modal open={isOpen} title="Create New Trade" onClose={() => isOpen = false}>
  <form class="space-y-4">
    <Input
      label="Symbol"
      bind:value={symbol}
      placeholder="AAPL"
      required={true}
    />

    <Select
      label="Direction"
      bind:value={direction}
      options={[
        { value: 'LONG', label: 'Long' },
        { value: 'SHORT', label: 'Short' }
      ]}
      required={true}
    />

    <div class="flex gap-3 pt-4">
      <Button variant="ghost" color="neutral" onclick={() => isOpen = false} class="flex-1">
        Cancel
      </Button>
      <Button variant="gradient" color="primary" type="submit" class="flex-1">
        Create Trade
      </Button>
    </div>
  </form>
</Modal>
```

### Stat Card with Trend

```svelte
<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Icon from '@iconify/svelte';

  interface Props {
    label: string;
    value: string;
    trend: number; // percentage
    icon: string;
    gradient: string;
    color: string;
  }

  let { label, value, trend, icon, gradient, color }: Props = $props();

  const isPositive = $derived(trend >= 0);
  const trendIcon = $derived(isPositive ? 'mdi:trending-up' : 'mdi:trending-down');
  const trendColor = $derived(isPositive ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400');
</script>

<Card variant="glass" padding="lg" hover={true}>
  <div class="flex items-start justify-between mb-4">
    <div class="w-12 h-12 rounded-xl bg-gradient-to-br {gradient} flex items-center justify-center shadow-lg shadow-{color}-500/30">
      <Icon icon={icon} width="24" class="text-white" />
    </div>
    <div class="flex items-center gap-1 px-2 py-1 rounded-full bg-{isPositive ? 'emerald' : 'red'}-100 dark:bg-{isPositive ? 'emerald' : 'red'}-900/30">
      <Icon icon={trendIcon} width="14" class={trendColor} />
      <span class="text-xs font-semibold {trendColor}">
        {Math.abs(trend)}%
      </span>
    </div>
  </div>

  <div>
    <p class="text-sm text-slate-600 dark:text-slate-400 mb-1">{label}</p>
    <p class="text-3xl font-bold text-slate-900 dark:text-slate-100">{value}</p>
  </div>
</Card>
```

### Navigation Dock (Reference Implementation)

```svelte
<script lang="ts">
  import { page } from '$app/stores';
  import Icon from '@iconify/svelte';

  const navItems = [
    { href: '/app/dashboard', icon: 'mdi:view-dashboard-outline', label: 'Overview', color: 'text-blue-500' },
    { href: '/app/trades', icon: 'mdi:chart-line-variant', label: 'Trades', color: 'text-emerald-500' },
    { href: '/app/journal', icon: 'mdi:book-open-page-variant-outline', label: 'Journal', color: 'text-purple-500' },
    { href: '/app/analytics', icon: 'mdi:chart-box-outline', label: 'Analytics', color: 'text-orange-500' },
    { href: '/app/settings', icon: 'mdi:cog-outline', label: 'Settings', color: 'text-slate-500' }
  ];

  function isActive(href: string): boolean {
    return $page.url.pathname.startsWith(href);
  }
</script>

<nav class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50">
  <div class="bg-white/70 dark:bg-slate-800/70 backdrop-blur-2xl rounded-2xl border border-slate-300/50 dark:border-slate-700/50 shadow-2xl px-3 py-2.5">
    <div class="flex items-center gap-2">
      {#each navItems as item}
        <a
          href={item.href}
          class="group relative flex flex-col items-center gap-1.5 px-4 py-2 rounded-xl transition-all duration-200 hover:bg-slate-100/80 dark:hover:bg-slate-700/80"
          class:bg-slate-100={isActive(item.href)}
          class:dark:bg-slate-700={isActive(item.href)}
        >
          <div class="relative">
            <Icon
              icon={item.icon}
              width="24"
              class="{item.color} transition-transform duration-200 group-hover:scale-110 {isActive(item.href) ? 'scale-110' : ''}"
            />
            {#if isActive(item.href)}
              <div class="absolute -bottom-3 left-1/2 -translate-x-1/2 w-1 h-1 rounded-full {item.color.replace('text-', 'bg-')}"></div>
            {/if}
          </div>
          <span class="text-[10px] font-medium text-slate-600 dark:text-slate-400 whitespace-nowrap">
            {item.label}
          </span>
        </a>
      {/each}
    </div>
  </div>
</nav>
```

### Empty State Pattern

```svelte
<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Icon from '@iconify/svelte';
</script>

<Card variant="glass" padding="lg">
  <div class="text-center py-12">
    <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-slate-400 to-slate-500 flex items-center justify-center mx-auto mb-4 shadow-lg shadow-slate-400/30">
      <Icon icon="mdi:notebook-outline" width="32" class="text-white" />
    </div>
    <h3 class="text-xl font-bold text-slate-900 dark:text-slate-100 mb-2">
      No journal entries yet
    </h3>
    <p class="text-sm text-slate-600 dark:text-slate-400 mb-6 max-w-md mx-auto">
      Start documenting your trading journey by creating your first journal entry. Track your emotions, rule adherence, and attach screenshots.
    </p>
    <Button variant="gradient" color="primary">
      <Icon icon="mdi:plus-circle-outline" width="20" />
      Create First Entry
    </Button>
  </div>
</Card>
```

### Loading State Pattern

```svelte
<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';

  let isLoading = $state(true);
</script>

{#if isLoading}
  <Card variant="glass" padding="lg">
    <div class="animate-pulse space-y-4">
      <!-- Icon skeleton -->
      <div class="flex items-center justify-between">
        <div class="w-12 h-12 rounded-xl bg-slate-300 dark:bg-slate-700"></div>
        <div class="w-16 h-6 rounded-full bg-slate-300 dark:bg-slate-700"></div>
      </div>

      <!-- Text skeleton -->
      <div class="space-y-2">
        <div class="w-24 h-4 rounded bg-slate-300 dark:bg-slate-700"></div>
        <div class="w-32 h-8 rounded bg-slate-300 dark:bg-slate-700"></div>
      </div>
    </div>
  </Card>
{:else}
  <!-- Actual content -->
{/if}
```

### Form Layout Pattern

```svelte
<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Select from '$lib/components/ui/Select.svelte';
  import Textarea from '$lib/components/ui/Textarea.svelte';
  import Button from '$lib/components/ui/Button.svelte';
</script>

<Card variant="glass" padding="lg">
  <form class="space-y-6">
    <div>
      <h2 class="text-xl font-bold text-slate-900 dark:text-slate-100 mb-1">
        Trade Details
      </h2>
      <p class="text-sm text-slate-600 dark:text-slate-400">
        Enter the details of your trade
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <Input label="Symbol" placeholder="AAPL" required={true} />
      <Select
        label="Direction"
        options={[
          { value: 'LONG', label: 'Long' },
          { value: 'SHORT', label: 'Short' }
        ]}
        required={true}
      />
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Input label="Quantity" type="number" placeholder="100" required={true} />
      <Input label="Entry Price" type="number" placeholder="150.25" required={true} />
      <Input label="Stop Loss" type="number" placeholder="145.00" />
    </div>

    <Textarea
      label="Notes"
      placeholder="Why did you take this trade?"
      rows={4}
    />

    <div class="flex gap-3 pt-4 border-t border-slate-200/50 dark:border-slate-700/50">
      <Button variant="ghost" color="neutral" type="button" class="flex-1">
        Cancel
      </Button>
      <Button variant="gradient" color="primary" type="submit" class="flex-1">
        Create Trade
      </Button>
    </div>
  </form>
</Card>
```

---

## Version History

- **v2.0** (2025-01-13) - Complete redesign with macOS-inspired dock navigation, glassmorphism, and gradient system. Added comprehensive component examples.
- **v1.0** - Original sidebar design (deprecated)
