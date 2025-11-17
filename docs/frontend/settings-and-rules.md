# TradePulse Settings Page & Rules Management

**Date:** 2025-01-13
**Status:** ✅ Complete

## Overview

Complete redesign of the Settings page using TabContainer component with comprehensive trading rules management system, enhanced preferences, and improved account management.

---

## 🎯 Key Features

### 1. Tab-Based Navigation

**Component Used:** `TabContainer.svelte`
- 4 main tabs: Profile, Trading Rules, Preferences, Account
- Icons for each tab (mdi:account, mdi:gavel, mdi:cog, mdi:shield-account)
- Active tab highlighting
- Smooth transitions between tabs

### 2. Trading Rules Management ⭐ NEW

**Purpose:** Allow traders to create, edit, and manage personal trading rules that will be tracked in journal entries.

**Features:**
- Create custom trading rules
- Edit existing rules
- Delete rules
- Color-coded categories and phases
- Weighted importance (1-5)
- Empty state with call-to-action

**Rule Structure:**
```typescript
interface Rule {
	id: string;
	title: string;
	description: string;
	weight: number;  // 1-5 (importance)
	phase: RulePhase;  // PRE_TRADE | DURING_TRADE | POST_TRADE
	category: RuleCategory;  // RISK_MANAGEMENT | ENTRY | EXIT | etc.
	created_at: string;
}
```

### 3. Enhanced Profile Tab

**Features:**
- Email display (read-only)
- Help text explaining how to update profile
- FormSection wrapper for consistent styling

### 4. Enhanced Preferences Tab

**Sections:**
- **Display Preferences:** Theme selection (Light, Dark, Auto)
- **Notification Preferences:** Email notifications, trade reminders
- Collapsible sections using FormSection

### 5. Enhanced Account Tab

**Features:**
- Sign out button with icon
- Danger zone for account deletion
- Visual separation between sections
- Toast notification for unimplemented features

---

## 📋 Trading Rules Management Details

### Rule Categories

1. **RISK_MANAGEMENT** (Red badge)
   - Position sizing rules
   - Stop loss requirements
   - Max daily loss limits

2. **ENTRY** (Green badge)
   - Entry criteria
   - Setup requirements
   - Confirmation signals

3. **EXIT** (Yellow badge)
   - Take profit rules
   - Stop loss management
   - Scaling out strategies

4. **POSITION_SIZING** (Neutral badge)
   - Size calculation methods
   - Leverage limits
   - Account percentage rules

5. **TIMING** (Neutral badge)
   - Time of day restrictions
   - Holding period guidelines
   - Market condition rules

6. **PSYCHOLOGY** (Blue badge)
   - Emotional state requirements
   - Confidence level guidelines
   - Revenge trading prevention

7. **GENERAL** (Neutral badge)
   - Miscellaneous rules
   - General guidelines

### Rule Phases

1. **PRE_TRADE** (Blue badge)
   - Rules to check before entering a trade
   - Planning and preparation
   - Example: "Risk no more than 2% per trade"

2. **DURING_TRADE** (Yellow badge)
   - Rules to follow while in a position
   - Management and monitoring
   - Example: "Don't move stop loss closer to entry"

3. **POST_TRADE** (Green badge)
   - Rules for after closing a position
   - Review and learning
   - Example: "Journal every trade within 1 hour of closing"

### Rule Weight (Importance)

- **1 - Low Importance:** Nice to follow, but not critical
- **2 - Below Average:** Somewhat important
- **3 - Average:** Standard rule, should be followed
- **4 - High Importance:** Very important, critical for success
- **5 - Critical:** Must follow, non-negotiable

**Impact on Adherence Score:**
Higher weights mean the rule has more impact on the overall adherence score. A weight-5 rule that's violated will hurt the score more than a weight-1 rule.

### Add/Edit Rule Form

**Fields:**
1. **Rule Title** (required)
   - Short, descriptive title
   - Example: "Risk no more than 2% per trade"

2. **Description** (required)
   - Detailed explanation of the rule
   - Why it matters
   - How to follow it

3. **Phase** (required)
   - Select: Pre-Trade, During Trade, or Post-Trade

4. **Category** (required)
   - Select from 7 categories

5. **Importance** (required)
   - Weight 1-5 via dropdown

**Validation:**
- Title and description are required
- Toast error if fields are empty
- All fields must be filled before submission

**Actions:**
- **Add Rule:** Creates new rule, adds to list
- **Update Rule:** Saves changes to existing rule
- **Cancel:** Resets form, closes editor

### Rule Display Cards

Each rule is displayed in a Card with:

**Header:**
- Rule title (large, bold)
- Category badge (color-coded)
- Phase badge (color-coded)

**Body:**
- Rule description
- Metadata:
  - Importance (stars icon)
  - Creation date (calendar icon)

**Actions:**
- **Edit button:** Opens rule in edit mode
- **Delete button:** Prompts for confirmation, then deletes

**Visual Example:**
```
┌──────────────────────────────────────────────────┐
│ Risk no more than 2% per trade                   │
│ [Risk Management] [Pre-Trade]                    │
│                                                   │
│ Position size should be calculated so maximum    │
│ loss is 2% of account                            │
│                                                   │
│ ⭐ Importance: 5/5  📅 Created 01/13/2025         │
│                                         [Edit] [Delete] │
└──────────────────────────────────────────────────┘
```

### Empty State

When no rules exist:
- Large gavel icon (48px, gray)
- "No Trading Rules Yet" heading
- Descriptive text
- "Add Your First Rule" button (gradient, primary)

### Color Coding System

**Category Colors (getCategoryColor function):**
```typescript
RISK_MANAGEMENT → 'error' (red)
ENTRY → 'success' (green)
EXIT → 'warning' (yellow)
PSYCHOLOGY → 'primary' (blue)
Default → 'neutral' (gray)
```

**Phase Colors (getPhaseColor function):**
```typescript
PRE_TRADE → 'primary' (blue)
DURING_TRADE → 'warning' (yellow)
POST_TRADE → 'success' (green)
```

---

## 🎨 UI Components Used

### TabContainer
```svelte
<TabContainer {tabs} bind:activeTab onTabChange={(tabId) => (activeTab = tabId)}>
	<!-- Tab content -->
</TabContainer>
```

### FormSection
```svelte
<FormSection
	title="Profile Information"
	icon="mdi:account-circle"
	helpText="View and manage your account information"
>
	<!-- Section content -->
</FormSection>
```

### Badges
```svelte
<Badge color={getCategoryColor(rule.category)} variant="soft" size="sm">
	{categoryOptions.find((c) => c.value === rule.category)?.label}
</Badge>
```

### Buttons
```svelte
<Button variant="gradient" color="primary" onclick={() => (showAddRule = true)}>
	<Icon icon="mdi:plus" width="20" />
	Add New Rule
</Button>
```

---

## 🔄 State Management

### Component State

```typescript
let activeTab = $state('profile');  // Current tab
let rules = $state<Rule[]>([...]);  // Rules array
let showAddRule = $state(false);   // Show add form
let editingRule: Rule | null = $state(null);  // Rule being edited
let newRule = $state({ ... });  // Form data
```

### CRUD Operations

#### Create Rule
```typescript
function handleAddRule() {
	if (!newRule.title || !newRule.description) {
		toast.error('Please fill in all required fields');
		return;
	}

	const rule: Rule = {
		id: crypto.randomUUID(),
		...newRule,
		created_at: new Date().toISOString()
	};

	rules = [...rules, rule];
	toast.success('Rule added successfully');
	resetRuleForm();
}
```

#### Read Rules
```typescript
// Displayed via reactive {#each} loop
{#each rules as rule (rule.id)}
	<Card>...</Card>
{/each}
```

#### Update Rule
```typescript
function handleUpdateRule() {
	rules = rules.map((r) =>
		r.id === editingRule!.id ? { ...r, ...newRule } : r
	);
	toast.success('Rule updated successfully');
	resetRuleForm();
}
```

#### Delete Rule
```typescript
function handleDeleteRule(id: string) {
	if (confirm('Are you sure you want to delete this rule?')) {
		rules = rules.filter((r) => r.id !== id);
		toast.success('Rule deleted');
	}
}
```

---

## 📊 Integration with Journal Entries

### How Rules Are Used

1. **Rule Selection in Journal:**
   - When creating a journal entry, traders select which rules to evaluate
   - Score each rule: 0, 25, 50, 75, or 100
   - Add notes explaining adherence

2. **Adherence Score Calculation:**
   - Weighted average based on rule importance
   - Overall score (0-100)
   - Phase-specific scores (pre-trade, during-trade, post-trade)

3. **Analytics Correlation:**
   - Chart showing adherence % vs. P&L
   - Identify which rules matter most
   - See impact of following/breaking rules

### Example Journal Entry Rule Adherence

```typescript
rule_adherence: [
	{
		rule_id: '1',
		rule_title: 'Risk no more than 2% per trade',
		score: 100,  // Fully followed
		notes: 'Calculated exact position size for 2% risk',
		timestamp: '2025-01-13T10:30:00Z'
	},
	{
		rule_id: '2',
		rule_title: 'Set stop loss before entry',
		score: 75,  // Mostly followed
		notes: 'Set stop loss 30 seconds after entry',
		timestamp: '2025-01-13T10:30:00Z'
	}
]
```

---

## 🎓 User Guidance

### Help Text Examples

**Trading Rules Tab:**
> "Create and manage your personal trading rules. These rules will be tracked in your journal entries to measure adherence and correlate with performance."

**Rule Weight:**
> "Weight (1-5) determines how much this rule affects your adherence score. Higher weights mean the rule is more critical to follow."

**Profile Tab:**
> "Contact support to change your email address or update profile information."

**Preferences Tab:**
> "Theme changes will apply immediately. Auto mode follows your system preferences."

### Suggested Default Rules

When users first access the Rules tab, they could be offered common rules:

1. **Risk Management:**
   - Risk no more than 2% per trade
   - Set stop loss before entry
   - No adding to losing positions

2. **Entry:**
   - Wait for confirmation signal
   - Check multiple timeframes
   - Ensure 1.5:1 R:R minimum

3. **Exit:**
   - Don't move stop closer to entry
   - Take partial profits at target 1
   - Journal trade within 1 hour

4. **Psychology:**
   - No revenge trading after loss
   - Max 3 trades per day
   - Don't trade when stressed (>7/10)

---

## 📱 Responsive Design

### Mobile Layout

- Tabs stack vertically on small screens
- Form fields go full-width
- Buttons stack vertically
- Cards maintain padding

### Desktop Layout

- Tabs horizontal in row
- Forms use max-width constraints (max-w-md)
- Grid layouts for rule fields (3 columns)
- Cards maintain consistent spacing

---

## 🔧 Technical Implementation

### Tab Structure

```svelte
const tabs = [
	{ id: 'profile', label: 'Profile', icon: 'mdi:account' },
	{ id: 'rules', label: 'Trading Rules', icon: 'mdi:gavel' },
	{ id: 'preferences', label: 'Preferences', icon: 'mdi:cog' },
	{ id: 'account', label: 'Account', icon: 'mdi:shield-account' }
];
```

### Form Validation

Client-side validation before submission:
- Required fields checked
- Toast error if validation fails
- Focus returns to form

### Confirmation Dialogs

Native `confirm()` for destructive actions:
```typescript
if (confirm('Are you sure you want to delete this rule?')) {
	// Delete
}
```

---

## 🚀 Future Enhancements

### Phase 1 (Ready for Implementation)

1. **Rule Templates:**
   - Pre-built rule library
   - One-click add from templates
   - Community-shared rules

2. **Rule Reordering:**
   - Drag-and-drop to reorder
   - Priority sorting
   - Group by category/phase

3. **Rule Sets:**
   - Multiple rule sets (Day Trading, Swing Trading)
   - Switch between sets
   - Export/import rules

### Phase 2 (Advanced Features)

1. **Smart Rules:**
   - Conditional logic ("If X, then Y")
   - Automated checking
   - Real-time violation alerts

2. **Rule Analytics:**
   - Per-rule performance correlation
   - Identify most impactful rules
   - Trend analysis over time

3. **Collaboration:**
   - Share rules with other traders
   - Community rating system
   - Expert-curated rule sets

4. **Gamification:**
   - Streak tracking (days followed all rules)
   - Achievements for adherence
   - Visual progress indicators

---

## 📁 Files Modified

### Modified Files

**`src/routes/app/settings/+page.svelte`**
- Complete redesign using TabContainer
- Added comprehensive rules management
- Enhanced preferences section
- Improved account management
- Added toast notifications
- Integrated FormSection components

**Changes:**
- Replaced manual tab buttons with TabContainer
- Added rules CRUD operations
- Added color-coding functions
- Added form validation
- Added empty states
- Enhanced accessibility (fieldsets, legends)

---

## 📊 Before vs After Comparison

### Before

**Layout:**
- Manual tab buttons (inline styles)
- Basic cards
- Minimal functionality
- Only 3 tabs (Profile, Preferences, Account)

**Features:**
- Email display (read-only)
- Theme selection (radio buttons)
- Logout button
- Delete account button (placeholder)

**Styling:**
- Inconsistent styling
- No icons
- No help text
- Basic card layout

### After

**Layout:**
- TabContainer component (reusable)
- FormSection wrappers
- Professional UI
- 4 tabs (added Trading Rules)

**Features:**
- All previous features +
- **Rules Management:**
  - Create, edit, delete rules
  - Category and phase classification
  - Weighted importance
  - Color-coded badges
  - Empty state handling
- **Enhanced Preferences:**
  - Notification settings
  - Collapsible sections
  - Help text
- **Toast Notifications:**
  - Success messages
  - Error validation
  - Warning for unimplemented features

**Styling:**
- Consistent glassmorphism theme
- Icons throughout
- Color-coded badges
- Help text with collapsible sections
- FormSection with icons and tooltips
- Responsive grid layouts

---

## 🎯 Design Goals Achieved

✅ **Comprehensive Rules Management** - Full CRUD for trading rules
✅ **Tab-Based Navigation** - Clean, organized settings
✅ **Consistent Styling** - Uses FormSection and TabContainer
✅ **User Guidance** - Help text throughout
✅ **Visual Feedback** - Toast notifications for all actions
✅ **Color Coding** - Categories and phases clearly identified
✅ **Empty States** - Helpful CTAs when no data exists
✅ **Accessibility** - Proper fieldsets, legends, labels
✅ **Responsive Design** - Works on mobile and desktop
✅ **Future-Ready** - Structure supports easy expansion

---

## 🧪 Testing Scenarios

### Rules Management

**Create Rule:**
1. Click "Add New Rule"
2. Fill in all fields
3. Click "Add Rule"
4. ✅ Rule appears in list
5. ✅ Toast success notification
6. ✅ Form resets

**Edit Rule:**
1. Click "Edit" on a rule
2. Modify fields
3. Click "Update Rule"
4. ✅ Changes saved
5. ✅ Toast success notification
6. ✅ Form closes

**Delete Rule:**
1. Click "Delete" on a rule
2. Confirm deletion
3. ✅ Rule removed from list
4. ✅ Toast success notification

**Validation:**
1. Click "Add New Rule" without filling fields
2. ✅ Toast error shown
3. ✅ Form stays open

**Empty State:**
1. Delete all rules
2. ✅ Empty state card shows
3. ✅ CTA button displays
4. Click CTA
5. ✅ Add form opens

### Navigation

**Tab Switching:**
1. Click each tab
2. ✅ Active tab highlights
3. ✅ Content updates
4. ✅ Smooth transition

### Accessibility

**Keyboard Navigation:**
1. Tab through form fields
2. ✅ Proper focus order
3. ✅ Visual focus indicators

**Screen Readers:**
1. Use screen reader
2. ✅ Labels announced correctly
3. ✅ Fieldset/legend for radio groups
4. ✅ Help text accessible

---

## 💡 Usage Examples

### Adding a Risk Management Rule

```typescript
Title: "Risk no more than 2% per trade"
Description: "Position size should be calculated so maximum loss is 2% of account equity"
Phase: "Pre-Trade"
Category: "Risk Management"
Importance: 5 (Critical)
```

### Adding a Psychology Rule

```typescript
Title: "No revenge trading after loss"
Description: "Wait at least 30 minutes after a losing trade before entering another position"
Phase: "Post-Trade"
Category: "Psychology"
Importance: 4 (High Importance)
```

### Adding a Timing Rule

```typescript
Title: "No trading in first 30 minutes"
Description: "Wait for market open volatility to settle before taking any trades"
Phase: "Pre-Trade"
Category: "Timing"
Importance: 3 (Average)
```

---

## 📚 Related Documentation

- **theme-improvements.md** - Form component styling
- **analytics-enhancements.md** - Rule adherence analytics
- **Types (types/index.ts)** - Rule, RulePhase, RuleCategory interfaces

---

**Status:** Ready for production use
**Next Steps:** Integrate with backend API for persistent storage
**API Endpoints Needed:**
- `GET /api/rules` - Fetch user's rules
- `POST /api/rules` - Create new rule
- `PUT /api/rules/:id` - Update rule
- `DELETE /api/rules/:id` - Delete rule
