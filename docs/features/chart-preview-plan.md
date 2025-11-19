# 1-Minute Chart Preview Feature

**Status:** ✅ Implemented (November 18, 2025)
**Priority:** High
**Target:** Journal Detail View

> **Note:** This document originally outlined the implementation plan. The feature has been completed using TradingView iframe embeds instead of the originally planned ECharts approach. See [Backend Implementation Status](../backend/implementation-status.md#-november-18-2025-session-4---chart-integration--journal-enhancements) for details on the actual implementation.

## Overview

Add an interactive 1-minute candlestick chart to the journal detail panel that visualizes trade entries and exits with markers, helping traders review their execution timing and decision-making.

## Goals

1. **Visual Trade Analysis** - See exactly when entries/exits occurred relative to price action
2. **Execution Quality** - Identify if entries were at good prices vs chasing
3. **Exit Timing** - Understand if exits were premature or well-timed
4. **Pattern Recognition** - Spot recurring execution patterns across multiple trades

## Technical Approach

### Option 1: Apache ECharts (Recommended)
**Why:** Already used in the dashboard, powerful candlestick support, lightweight

```typescript
// Chart configuration
const chartOptions = {
  xAxis: { type: 'time' },
  yAxis: { type: 'value', scale: true },
  series: [
    {
      type: 'candlestick',
      data: priceData // [[timestamp, open, close, low, high], ...]
    },
    {
      type: 'scatter',  // Entry markers
      symbol: 'triangle',
      symbolSize: 12,
      itemStyle: { color: '#10b981' }, // Green
      data: entryPoints
    },
    {
      type: 'scatter',  // Exit markers
      symbol: 'triangle',
      symbolRotate: 180,
      symbolSize: 12,
      itemStyle: { color: '#ef4444' }, // Red
      data: exitPoints
    }
  ]
}
```

### Option 2: TradingView Lightweight Charts
**Why:** Specialized for financial charts, better performance for large datasets

**Pros:**
- Highly optimized for trading charts
- Built-in candlestick rendering
- Professional look and feel
- Touch-friendly for mobile

**Cons:**
- Additional dependency (~100KB)
- Learning curve for new library

### Option 3: Chart.js with Financial Plugin
**Why:** Popular, well-documented

**Pros:**
- Large community
- Many examples
- Flexible customization

**Cons:**
- Not specialized for trading
- Less performant for large datasets
- Requires financial plugin

## Data Requirements

### Backend API Endpoint
```
GET /api/trades/{id}/chart-data?timeframe=1m&range=auto
```

**Query Parameters:**
- `timeframe`: `1m`, `5m`, `15m`, etc.
- `range`: `auto` (trade duration + buffer), `full-day`, or custom `YYYY-MM-DD`
- `buffer_minutes`: Extra time before/after trade (default: 30)

**Response:**
```json
{
  "success": true,
  "data": {
    "symbol": "SPY",
    "timeframe": "1m",
    "candles": [
      {
        "timestamp": "2025-01-18T09:30:00Z",
        "open": 580.50,
        "high": 580.75,
        "low": 580.25,
        "close": 580.60,
        "volume": 125000
      }
    ],
    "entries": [
      {
        "id": "entry-1",
        "timestamp": "2025-01-18T09:35:22Z",
        "price": 580.65,
        "quantity": 100,
        "notes": "Entry at support level"
      }
    ],
    "exits": [
      {
        "id": "exit-1",
        "timestamp": "2025-01-18T10:15:45Z",
        "price": 581.45,
        "quantity": 100,
        "pnl": 80.00,
        "notes": "Exit at resistance"
      }
    ],
    "trade_info": {
      "opened_at": "2025-01-18T09:35:22Z",
      "closed_at": "2025-01-18T10:15:45Z",
      "duration_minutes": 40,
      "realized_pnl": 80.00
    }
  }
}
```

### Data Sources

**Phase 1 - Manual Entry:**
- Use existing entry/exit timestamps from database
- No historical price data
- Show timeline with entry/exit markers only
- Message: "Chart data available when historical prices are imported"

**Phase 2 - CSV Import:**
- Allow users to upload 1-minute bar CSV data
- Store in separate `chart_data` table
- Link to trades by symbol + timestamp range

**Phase 3 - Live Data Provider (Future):**
- Integration with Polygon.io, Alpaca, or similar
- Automatic chart data fetching
- Real-time updates for open positions

## Database Schema

```sql
-- New table for storing historical chart data
CREATE TABLE chart_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) NOT NULL,
    timeframe VARCHAR(10) NOT NULL, -- '1m', '5m', '15m', etc.
    timestamp TIMESTAMPTZ NOT NULL,
    open DECIMAL(12, 4) NOT NULL,
    high DECIMAL(12, 4) NOT NULL,
    low DECIMAL(12, 4) NOT NULL,
    close DECIMAL(12, 4) NOT NULL,
    volume BIGINT,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    -- Indexes for fast querying
    CONSTRAINT unique_symbol_timeframe_timestamp UNIQUE (symbol, timeframe, timestamp)
);

CREATE INDEX idx_chart_data_symbol_timeframe ON chart_data(symbol, timeframe);
CREATE INDEX idx_chart_data_timestamp ON chart_data(timestamp);
CREATE INDEX idx_chart_data_symbol_timestamp ON chart_data(symbol, timestamp);
```

## UI Components

### ChartPreview.svelte
```svelte
<script lang="ts">
  import ECharts from '$lib/components/charts/ECharts.svelte';

  interface Props {
    tradeId: string;
    symbol: string;
    entries: Entry[];
    exits: Exit[];
  }

  let { tradeId, symbol, entries, exits }: Props = $props();

  let chartData = $state<ChartData | null>(null);
  let loading = $state(true);
  let error = $state<string | null>(null);

  // Load chart data
  async function loadChartData() {
    try {
      loading = true;
      const response = await apiClient.getTradeChartData(tradeId, {
        timeframe: '1m',
        range: 'auto',
        buffer_minutes: 30
      });
      chartData = response.data;
    } catch (err) {
      error = 'Failed to load chart data';
      console.error(err);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadChartData();
  });
</script>

<div class="chart-container">
  {#if loading}
    <LoadingSpinner />
  {:else if error}
    <ErrorMessage message={error} />
  {:else if chartData}
    <ECharts options={buildChartOptions(chartData)} height={400} />
  {:else}
    <EmptyState message="No chart data available" />
  {/if}
</div>
```

## Implementation Phases

### Phase 1: Basic Timeline (Week 1)
- [ ] Display entry/exit markers on a timeline
- [ ] Show timestamps and prices
- [ ] No candlestick data yet
- [ ] Basic hover tooltips

**Deliverable:** Visual timeline showing trade execution points

### Phase 2: Chart Data Storage (Week 2)
- [ ] Create `chart_data` table
- [ ] CSV upload endpoint for historical data
- [ ] Parse and store 1-minute bars
- [ ] API endpoint to fetch chart data

**Deliverable:** Backend infrastructure for chart data

### Phase 3: Candlestick Chart (Week 3)
- [ ] Implement ECharts candlestick component
- [ ] Overlay entry/exit markers
- [ ] Add zoom and pan controls
- [ ] Price scale and time axis formatting

**Deliverable:** Full interactive candlestick chart

### Phase 4: Advanced Features (Week 4)
- [ ] Multiple timeframe support (1m, 5m, 15m)
- [ ] Volume bars
- [ ] Drawing tools (trendlines, support/resistance)
- [ ] Screenshot capture
- [ ] Annotations on chart

**Deliverable:** Production-ready chart with advanced features

### Phase 5: Data Provider Integration (Future)
- [ ] Polygon.io integration
- [ ] Alpaca integration
- [ ] Automatic data fetching
- [ ] Real-time updates

## Chart Features

### Core Features
- **Candlestick Chart** - OHLC data visualization
- **Entry Markers** - Green triangles pointing up
- **Exit Markers** - Red triangles pointing down
- **Hover Tooltips** - Show exact price, time, quantity, P&L
- **Zoom Controls** - Pinch zoom, mouse wheel zoom
- **Pan** - Drag to navigate
- **Crosshair** - Precise price/time reading

### Advanced Features
- **Volume Bars** - Show volume below price chart
- **Position Size Overlay** - Show how position changed over time
- **P&L Line** - Real-time P&L during trade
- **Time Markers** - Highlight key moments (journal entries, rule checks)
- **Pattern Highlights** - Support/resistance levels
- **Export** - Download chart as image
- **Print** - Print-friendly view

## Mobile Considerations

- **Touch Gestures** - Pinch zoom, swipe pan
- **Simplified Markers** - Larger touch targets
- **Responsive Size** - Adjust height based on screen
- **Minimal UI** - Hide advanced controls on small screens
- **Performance** - Limit candle count on mobile

## Performance Optimization

1. **Data Limits**
   - Max 500 candles in viewport
   - Lazy load additional data on zoom out
   - Aggregate to higher timeframes for large ranges

2. **Caching**
   - Cache chart data in localStorage (5 min TTL)
   - Debounce zoom/pan events
   - Use virtual rendering for large datasets

3. **Progressive Loading**
   - Show markers first
   - Load candle data in background
   - Display progressively as data arrives

## Testing Strategy

1. **Unit Tests**
   - Data transformation functions
   - Price scaling calculations
   - Marker positioning logic

2. **Integration Tests**
   - API data fetching
   - Chart rendering with real data
   - Entry/exit marker accuracy

3. **Visual Regression Tests**
   - Screenshot comparison
   - Cross-browser rendering
   - Mobile vs desktop layouts

4. **Performance Tests**
   - Render time with 500 candles
   - Memory usage over time
   - Zoom/pan responsiveness

## User Stories

**As a trader, I want to:**
1. See my entry points on a 1-minute chart to understand if I entered at good prices
2. Compare my exits to subsequent price action to improve my exit strategy
3. Identify patterns in my execution (e.g., always chasing, poor timing)
4. Review multiple trades to spot consistent execution issues
5. Share chart screenshots with my trading mentor

## Success Metrics

- **Adoption:** 80%+ of journal entries viewed with chart preview
- **Engagement:** Users spend 30s+ analyzing chart per entry
- **Feedback:** 4.5+ star rating for chart feature
- **Performance:** Chart loads in < 2 seconds
- **Accuracy:** 100% alignment of entry/exit markers with trade data

## Future Enhancements

1. **AI Insights**
   - "You entered near the top of this candle"
   - "Exit was well-timed, price dropped 2% after"
   - Pattern recognition: "This looks like a failed breakout"

2. **Comparison View**
   - Compare execution across similar setups
   - See all SPY entries overlaid
   - Best vs worst executions side-by-side

3. **Replay Mode**
   - Watch the trade unfold bar-by-bar
   - Pause at entry/exit points
   - Review decision-making in real-time context

4. **Strategy Backtesting**
   - Test "what if" scenarios
   - Different entry/exit timing
   - Calculate hypothetical P&L

## Resources Needed

- **Development:** 4 weeks (1 developer)
- **Data Provider:** Polygon.io Starter ($29/mo) or Alpaca (Free tier)
- **Storage:** ~100MB per year of 1-minute data per symbol
- **Infrastructure:** No significant changes needed

## Risks & Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Data provider costs | High | Medium | Start with manual CSV, add provider as premium feature |
| Performance issues | Medium | Low | Implement data limits and caching |
| Chart library complexity | Low | Medium | Use ECharts (already familiar) |
| Mobile usability | Medium | Medium | Extensive mobile testing, touch gestures |
| Data accuracy | High | Low | Validate timestamps, automated tests |

## Decision Log

**2025-01-18:** Chose Apache ECharts over TradingView for consistency with existing dashboard
**TBD:** Data provider selection (Phase 5)
**TBD:** Timeframe options (start with 1m only)

---

## Actual Implementation (November 18, 2025)

### What Was Built

Instead of the complex ECharts implementation with backend chart data storage, we implemented a simpler, faster solution:

**TradingView Iframe Embed:**
- Used TradingView's free embed widget (iframe-based)
- No backend chart data storage needed
- No API endpoint for historical price data
- Shows current market charts for symbols
- Entry/exit points displayed in collapsible timeline below chart

**Components Created:**
- `ChartPreview.svelte` - Embedded in journal detail view
- `ChartPortal.svelte` - Global fullscreen chart overlay
- `chartPortal` store - Global state management for chart portal

**Features Delivered:**
- 1-minute TradingView charts for each trade
- Maximize chart button for fullscreen view
- Execution timeline with timestamps, prices, and P&L
- Collapsible timeline to save space
- Trade summary stats (avg entry, avg exit, quantity, executions)
- Dark mode support
- Mobile responsive

### Technical Decisions

**Why TradingView instead of ECharts:**
1. No need to source/store historical price data
2. No backend development required
3. Professional-grade charts with built-in features
4. Zero cost (using free embed)
5. Faster time to market (2 hours vs 4 weeks)

**Limitations vs Original Plan:**
- ❌ No programmatic entry/exit markers on chart (not supported by free TradingView)
- ❌ No historical price data storage
- ❌ Charts show current market data, not trade-time data
- ✅ Timeline view compensates by showing all execution details
- ✅ Users can manually draw on chart using TradingView tools

### What Was NOT Built

- Chart data table (`chart_data`)
- CSV upload for historical price data
- `/api/trades/{id}/chart-data` endpoint
- Data provider integration (Polygon.io, Alpaca)
- Cost basis overlay on chart
- Automated marker placement

### Future Enhancements

If we need the originally planned features:
1. Add backend chart data storage
2. Integrate with data provider API
3. Switch to TradingView advanced charts or ECharts
4. Implement programmatic marker overlay

For now, the TradingView iframe solution provides 80% of the value with 5% of the complexity.

---

**Original Status:** Planning → **Current Status:** ✅ Implemented (Simplified)
**Implementation Date:** November 18, 2025
**Time to Implement:** ~2 hours (vs estimated 4 weeks)
**Files:** See `docs/backend/implementation-status.md` Session 4
