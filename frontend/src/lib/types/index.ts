// Position Lifecycle Types
export interface Entry {
	id: string;
	price: number;
	quantity: number;
	timestamp: string;
	notes?: string;
	fees?: number;
}

export interface Exit {
	id: string;
	price: number;
	quantity: number;
	timestamp: string;
	notes?: string;
	fees?: number;
	pnl: number;
}

export type CostBasisMethod = 'FIFO' | 'LIFO' | 'AVERAGE';

export interface Trade {
	id: string;
	user_id: string;
	symbol: string;
	trade_type: 'LONG' | 'SHORT';
	// Legacy fields (for backward compatibility)
	quantity: number;
	entry_price: number;
	exit_price: number | null;
	fees: number;
	pnl: number | null;
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
	// Journal entries
	journal_entries?: JournalEntry[];
	has_journal?: boolean;
}

export interface JournalEntry {
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

export interface EmotionalState {
	confidence: number; // 1-10
	stress: number; // 1-10
	discipline: number; // 1-10
	notes: string;
}

export interface User {
	id: string;
	email: string;
	has_password: boolean;
	plan_type?: string;
	plan_status?: string;
	plan_selected_at?: string;
	created_at: string;
	last_login: string;
}

export interface AuthResponse {
	jwt: string;
	user: User;
}

export interface Notification {
	id: string;
	type: string;
	user_id: string;
	title: string;
	message: string;
	data?: any;
	timestamp: string;
	read: boolean;
}

// Rule Set Types
export type RulePhase = 'PRE_TRADE' | 'DURING_TRADE' | 'POST_TRADE';
export type RuleCategory = 'RISK_MANAGEMENT' | 'ENTRY' | 'EXIT' | 'POSITION_SIZING' | 'TIMING' | 'PSYCHOLOGY' | 'GENERAL';

export interface Rule {
	id: string;
	title: string;
	description: string;
	weight: number; // 1-5 (importance)
	phase: RulePhase;
	category: RuleCategory;
	created_at: string;
}

export interface RuleSet {
	id: string;
	user_id: string;
	name: string;
	description: string;
	rules: Rule[];
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface RuleAdherence {
	rule_id: string;
	rule_title: string;
	score: number; // 0, 25, 50, 75, 100
	notes: string;
	timestamp: string;
}

export interface AdherenceScore {
	overall_score: number; // 0-100
	weighted_score: number; // 0-100
	color: 'green' | 'yellow' | 'red';
	phase_scores: {
		pre_trade: number;
		during_trade: number;
		post_trade: number;
	};
}
