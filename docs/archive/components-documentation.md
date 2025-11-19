# Frontend Components Documentation

This document provides comprehensive documentation for all components created and updated in the TradePulse frontend application.

**Last Updated:** January 2025

---

## 📦 New Components

### 1. EntryExitManager.svelte
**Location:** `frontend/src/lib/components/trading/EntryExitManager.svelte`

**Purpose:** Manages multiple entries and exits for a single trade, supporting the advanced position management system.

**Props:**
- `tradeId: string` - ID of the trade
- `tradeType: 'LONG' | 'SHORT'` - Type of trade
- `entries?: any[]` - Array of entry executions
- `exits?: any[]` - Array of exit executions
- `onUpdate?: () => void` - Callback when entries/exits are modified

**Features:**
- Add new entries with price, quantity, timestamp, fees, and notes
- Add new exits (calculates P&L automatically on backend)
- Delete entries and exits with confirmation
- Separate UI sections for entries vs exits
- Color-coded by trade type (blue for entries, emerald for exits)
- Empty states when no entries/exits exist
- Inline forms that expand when clicking "Add Entry" or "Add Exit"
- Shows P&L for each exit execution

**Key Functions:**
- `handleAddEntry()` - Creates new entry via API
- `handleAddExit()` - Creates new exit via API
- `handleDeleteEntry(entryId)` - Deletes entry with confirmation
- `handleDeleteExit(exitId)` - Deletes exit with confirmation

**Usage Example:**
```svelte
<EntryExitManager
	tradeId={trade.id}
	tradeType={trade.trade_type}
	entries={trade.entries}
	exits={trade.exits}
	onUpdate={handleTradeUpdate}
/>
```

---

### 2. TradeReviewWizard.svelte
**Location:** `frontend/src/lib/components/trading/TradeReviewWizard.svelte`

**Purpose:** Comprehensive 5-tab wizard for reviewing completed trades with rule adherence tracking, emotional state logging, and journal entries.

**Props:**
- `trade: Trade` - The trade to review
- `onComplete: () => void` - Callback when review is completed
- `onSkip: () => void` - Callback when review is skipped
- `onCancel: () => void` - Callback when wizard is cancelled

**Tabs:**

1. **Overview Tab**
   - Trade summary metrics (P&L, position size, avg entry, fees)
   - Review checklist showing completion status
   - Quick reference of what needs to be done

2. **Entries & Exits Tab**
   - Full entry/exit management using EntryExitManager component
   - Add, view, delete entries and exits
   - Verify all executions are recorded correctly

3. **Rule Adherence Tab**
   - Select rule set from available rule sets
   - Rate adherence to each rule (0-100% with traffic light indicators)
   - Weighted adherence score calculation
   - Required notes for scores below 100%
   - Uses RuleAdherenceInput component for each rule

4. **Emotional State Tab**
   - Pre-trade confidence (1-10 slider)
   - Pre-trade clarity (1-10 slider)
   - Post-trade discipline (1-10 slider)
   - Post-trade emotion (text description)

5. **Journal Tab**
   - Trade Setup - thesis and confluence factors
   - Execution - trade management and adjustments
   - Reflection - what went well/could improve
   - Lessons Learned - key takeaways
   - Screenshot upload with FileUpload component
   - ImageGallery for viewing uploaded screenshots

**Features:**
- Tracks completion status for each tab
- Auto-loads existing journal data if trade was previously reviewed
- Saves all data as primary journal entry
- Marks trade as `is_reviewed: true` on completion
- Marks trade as `review_skipped: true` if skipped
- Full keyboard navigation support

**Key Functions:**
- `loadRuleSets()` - Fetches available rule sets
- `loadTradeData()` - Loads existing journal entries
- `checkTabCompletion()` - Updates completion status
- `handleComplete()` - Saves journal and marks reviewed
- `handleSkipReview()` - Marks trade as skipped

**Usage Example:**
```svelte
<TradeReviewWizard
	trade={currentTrade}
	onComplete={handleReviewComplete}
	onSkip={handleReviewSkip}
	onCancel={handleCancelWizard}
/>
```

---

### 3. ProfileSetupWizard.svelte
**Location:** `frontend/src/lib/components/onboarding/ProfileSetupWizard.svelte`

**Purpose:** 4-step wizard for new users to complete their profile after initial signup.

**Props:**
- `onComplete: () => void` - Callback when setup is completed

**Steps:**

1. **Personal Info**
   - First name (required)
   - Last name (required)
   - Phone (optional)
   - Company (optional)

2. **Address (Optional)**
   - Address Line 1
   - Address Line 2
   - City
   - State/Province
   - Postal Code

3. **Preferences**
   - Country selection
   - Timezone selection (11 major timezones)
   - Info about why timezone matters

4. **Complete**
   - Summary of entered information (including address if provided)
   - Confirmation before saving

**Features:**
- Cannot skip - required for new users
- Progress indicator with step icons
- Form validation (names required)
- Smooth animations between steps
- Auto-saves profile on completion
- Updates user store with `profile_completed: true`

**Key Functions:**
- `nextStep()` - Advances to next step with validation
- `prevStep()` - Returns to previous step
- `handleComplete()` - Saves profile via API

**Integration:**
- Automatically shown in app layout if `!user.profile_completed && !user.first_name`
- Modal overlay prevents access to app until completed
- Z-index set to `z-[100]` to appear above navigation (which uses `z-50`)

**API Integration:**
- Calls `apiClient.updateProfile()` with profile data
- Marks `profile_completed: true` on save
- Refreshes user data from server after save

---

## 📝 Updated Components

### 1. Settings Page (+page.svelte)
**Location:** `frontend/src/routes/app/settings/+page.svelte`

**Major Changes:**

**Profile Tab - Now Comprehensive:**
- Personal Info section (first/last name, email, phone, company)
- Address section (line 1, line 2, city, state, postal code, country)
- Preferences section (timezone with 11 major options)
- Save button to update all profile data

**Account Tab - New Features:**
- Reset Account Data option in danger zone
  - Deletes all trades, journal entries, rule sets
  - Keeps user account and authentication
  - Requires typing "DELETE ALL DATA" to confirm
  - Reloads page after reset

**New Functions:**
- `handleProfileUpdate()` - Saves profile changes via `apiClient.updateProfile()`
- `handleResetAccountData()` - Resets all user data via API

**Form State:**
```typescript
profileForm = {
	email, first_name, last_name, phone, company,
	address_line1, address_line2, city, state, postal_code,
	country, timezone
}
```

**API Integration:**
- Profile update now fully implemented with backend endpoint
- All profile fields saved to database
- Updates user store with fresh data after save

---

### 2. Trades Page (+page.svelte)
**Location:** `frontend/src/routes/app/trades/+page.svelte`

**Schema Updates:**
- Changed from old schema (quantity, entry_price, exit_price, pnl, fees)
- To new schema (current_position_size, average_entry_price, realized_pnl, etc.)

**New Column:**
- "Review" column added to trades table
- Shows review status with badges:
  - ✅ **Reviewed** (green) - `is_reviewed: true`
  - ⏭️ **Skipped** (amber) - `review_skipped: true`
  - ⚠️ **Pending** (red) - closed but not reviewed
  - ➖ **N/A** (gray) - trade still open

**Fixed:**
- P&L percentage calculation in tooltip
- Now uses `realized_pnl / (average_entry_price * current_position_size)`

---

### 3. App Layout (+layout.svelte)
**Location:** `frontend/src/routes/app/+layout.svelte`

**New Features:**
- Imports ProfileSetupWizard component
- Checks if profile is complete on mount
- Shows wizard if `!user.profile_completed && !user.first_name`
- Adds `clearAuthToken()` method to API client

**Profile Check Logic:**
```typescript
if (!user.profile_completed && !user.first_name) {
	showProfileWizard = true;
}
```

---

### 4. API Client (client.ts)
**Location:** `frontend/src/lib/api/client.ts`

**New Methods:**

```typescript
// Position lifecycle - Entries
getEntries(tradeId: string): Promise<any[]>
addEntry(tradeId: string, entry: any): Promise<any>
deleteEntry(tradeId: string, entryId: string): Promise<void>

// Position lifecycle - Exits
getExits(tradeId: string): Promise<any[]>
addExit(tradeId: string, exit: any): Promise<any>
deleteExit(tradeId: string, exitId: string): Promise<void>

// Account management
resetAccountData(): Promise<{ message: string }>
clearAuthToken(): void

// User profile
updateProfile(profileData: any): Promise<any>
```

---

## 📄 New Pages

### 1. Review Queue Page
**Location:** `frontend/src/routes/app/review/+page.svelte`

**Purpose:** Dedicated page for reviewing unreviewed closed trades.

**Features:**
- Shows count of pending reviews
- Stats cards (pending count, completion rate, avg time)
- List of unreviewed trades with key metrics
- "Start Reviewing" button to begin wizard
- Auto-advances to next trade after completing each review
- Empty state when all caught up
- Tips card with review best practices

**Key Functions:**
- `loadPendingTrades()` - Fetches closed trades where `!is_reviewed && !review_skipped`
- `startReview(trade, index)` - Opens wizard for specific trade
- `startNextReview()` - Opens wizard for first pending trade
- `handleCompleteReview()` - Auto-advances to next trade
- `handleSkipReview()` - Auto-advances to next trade

**Trade Card Display:**
- Trade symbol and type
- Close date
- P&L (color-coded)
- Position size
- Average entry price
- Entry/exit count
- Notes (if any)
- "Review Trade" button

---

## 🎨 Utility Files

### 1. Formatting Utilities
**Location:** `frontend/src/lib/utils/formatting.ts`

**Functions:**
```typescript
formatCurrency(value: number, currency = 'USD'): string
formatNumber(value: number, decimals = 2): string
formatPercentage(value: number, decimals = 2): string
formatDate(date: Date | string): string
formatDateTime(date: Date | string): string
formatTime(date: Date | string): string
```

---

### 2. Type Definitions
**Location:** `frontend/src/lib/types.ts`

**Interfaces:**
```typescript
Trade - Full trade with entries, exits, review status
Entry - Single entry execution
Exit - Single exit execution with P&L
RuleSet - Collection of trading rules
Rule - Individual rule with weight/phase/category
RuleAdherence - Rule adherence rating
JournalEntry - Journal entry with emotional state
EmotionalState - Pre/post trade emotional ratings
Attachment - File attachment for journal entries
```

---

## 🎨 Component Design Patterns

### Modal Pattern
Used by: TradeReviewWizard, ProfileSetupWizard
- Fixed positioning with backdrop blur
- Prevents body scroll
- ESC key to close
- Click outside to close
- Max height with scroll

### Wizard Pattern
Used by: TradeReviewWizard, ProfileSetupWizard
- Step-by-step progression
- Progress indicators
- Back/Next navigation
- Form validation
- Completion summary

### Inline Form Pattern
Used by: EntryExitManager
- Forms expand inline when clicking add button
- Show/hide with state variable
- Cancel button collapses form
- Success clears form and refreshes data

### Empty State Pattern
Used by: EntryExitManager, Review Queue
- Icon + message
- Call-to-action button
- Friendly messaging

---

## 🔄 Data Flow

### Trade Review Flow
1. User closes trade → appears in review queue
2. User clicks "Review Trade" → TradeReviewWizard opens
3. User completes all tabs → data saved as journal entry
4. Trade marked `is_reviewed: true`
5. Trade removed from queue
6. Next trade auto-opens (if any)

### Entry/Exit Management Flow
1. User adds entry → POST /api/trades/{id}/entries
2. Backend creates entry → recalculates metrics
3. Frontend refreshes trade data
4. User adds exit → POST /api/trades/{id}/exits
5. Backend calculates P&L → recalculates metrics
6. Frontend refreshes trade data
7. Position auto-closes when size reaches 0

### Profile Setup Flow
1. New user signs up → email verification
2. User lands on dashboard → layout checks profile
3. No profile → ProfileSetupWizard shows
4. User completes 3 steps → profile saved
5. Wizard closes → user accesses app

---

## 🧪 Testing Considerations

### EntryExitManager
- Test with LONG and SHORT trades
- Test with no entries/exits (empty state)
- Test with multiple entries/exits
- Test delete confirmation
- Test form validation

### TradeReviewWizard
- Test with trades that have/don't have existing reviews
- Test with no rule sets
- Test with multiple rule sets
- Test tab completion logic
- Test skip vs complete

### ProfileSetupWizard
- Test required field validation
- Test cannot skip
- Test back navigation
- Test summary accuracy

---

## 📱 Responsive Design

All components are fully responsive:
- Mobile: Single column layout, stacked forms
- Tablet: 2-column grids where appropriate
- Desktop: Full multi-column layouts

Breakpoints follow Tailwind defaults:
- `sm`: 640px
- `md`: 768px
- `lg`: 1024px
- `xl`: 1280px

---

## ♿ Accessibility

All components include:
- Semantic HTML elements
- ARIA labels and roles
- Keyboard navigation
- Focus management
- Screen reader support
- Color contrast compliance
- Icon + text labels

---

## 🎨 Theming

All components support:
- Light mode
- Dark mode
- Auto mode (system preference)

Theme classes use Tailwind's `dark:` variant:
```svelte
class="bg-white dark:bg-slate-900"
```

---

## 🚀 Performance

Optimizations:
- Lazy loading of modals/wizards
- Debounced form inputs where appropriate
- Minimal re-renders with Svelte 5 runes
- Efficient list rendering with keyed each blocks
- API request caching where applicable

---

## 📚 Dependencies

Component-specific dependencies:
- `@iconify/svelte` - Icons (bundled with @iconify-json/mdi)
- `$lib/api/client` - API client for data fetching
- `$lib/stores/toast` - Toast notifications
- `$lib/stores/user` - User state management

All UI components are built with:
- SvelteKit 5
- TypeScript
- TailwindCSS
- Svelte 5 runes ($state, $props, $derived, $effect)
