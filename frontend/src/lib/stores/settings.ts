import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface UserSettings {
	timezone: string; // IANA timezone string (e.g., 'America/New_York')
	useMarketTime: boolean; // If true, show times in market timezone
	dateFormat: 'short' | 'medium' | 'long';
	timeFormat: '12h' | '24h';
}

const DEFAULT_SETTINGS: UserSettings = {
	timezone: browser ? Intl.DateTimeFormat().resolvedOptions().timeZone : 'America/New_York',
	useMarketTime: false,
	dateFormat: 'medium',
	timeFormat: '12h'
};

const STORAGE_KEY = 'tradepulse_settings';

function createSettingsStore() {
	// Load from localStorage if available
	const stored = browser ? localStorage.getItem(STORAGE_KEY) : null;
	const initial = stored ? { ...DEFAULT_SETTINGS, ...JSON.parse(stored) } : DEFAULT_SETTINGS;

	const { subscribe, set, update } = writable<UserSettings>(initial);

	return {
		subscribe,
		update: (settings: Partial<UserSettings>) =>
			update(state => {
				const newState = { ...state, ...settings };
				if (browser) {
					localStorage.setItem(STORAGE_KEY, JSON.stringify(newState));
				}
				return newState;
			}),
		reset: () => {
			set(DEFAULT_SETTINGS);
			if (browser) {
				localStorage.removeItem(STORAGE_KEY);
			}
		}
	};
}

export const settingsStore = createSettingsStore();

// Market timezone constants
export const MARKET_TIMEZONES = {
	NYSE: 'America/New_York',
	NASDAQ: 'America/New_York',
	LSE: 'Europe/London',
	TSE: 'Asia/Tokyo',
	HKEX: 'Asia/Hong_Kong',
	SSE: 'Asia/Shanghai'
} as const;

// Common timezones for selection
export const COMMON_TIMEZONES = [
	{ value: 'America/New_York', label: 'Eastern Time (ET)' },
	{ value: 'America/Chicago', label: 'Central Time (CT)' },
	{ value: 'America/Denver', label: 'Mountain Time (MT)' },
	{ value: 'America/Los_Angeles', label: 'Pacific Time (PT)' },
	{ value: 'Europe/London', label: 'London (GMT/BST)' },
	{ value: 'Europe/Paris', label: 'Paris (CET/CEST)' },
	{ value: 'Asia/Tokyo', label: 'Tokyo (JST)' },
	{ value: 'Asia/Hong_Kong', label: 'Hong Kong (HKT)' },
	{ value: 'Asia/Shanghai', label: 'Shanghai (CST)' },
	{ value: 'Australia/Sydney', label: 'Sydney (AEDT/AEST)' }
];

/**
 * Format a date/time string according to user settings
 */
export function formatDateTime(
	dateString: string,
	settings: UserSettings,
	options: {
		showDate?: boolean;
		showTime?: boolean;
		dateOnly?: boolean;
		timeOnly?: boolean;
	} = {}
): string {
	const { showDate = true, showTime = true, dateOnly = false, timeOnly = false } = options;
	const date = new Date(dateString);

	// Determine which timezone to use
	const timezone = settings.useMarketTime ? MARKET_TIMEZONES.NYSE : settings.timezone;

	// Date formatting options
	const dateFormatOptions: Intl.DateTimeFormatOptions = {
		timeZone: timezone
	};

	if (dateOnly || showDate) {
		if (settings.dateFormat === 'short') {
			dateFormatOptions.month = 'numeric';
			dateFormatOptions.day = 'numeric';
			dateFormatOptions.year = '2-digit';
		} else if (settings.dateFormat === 'medium') {
			dateFormatOptions.month = 'short';
			dateFormatOptions.day = 'numeric';
			dateFormatOptions.year = 'numeric';
		} else {
			dateFormatOptions.month = 'long';
			dateFormatOptions.day = 'numeric';
			dateFormatOptions.year = 'numeric';
		}
	}

	if (timeOnly || (showTime && !dateOnly)) {
		dateFormatOptions.hour = '2-digit';
		dateFormatOptions.minute = '2-digit';
		dateFormatOptions.hour12 = settings.timeFormat === '12h';
	}

	return date.toLocaleString('en-US', dateFormatOptions);
}

/**
 * Get timezone abbreviation (e.g., "EST", "PST")
 */
export function getTimezoneAbbr(timezone: string, date: Date = new Date()): string {
	const formatter = new Intl.DateTimeFormat('en-US', {
		timeZone: timezone,
		timeZoneName: 'short'
	});

	const parts = formatter.formatToParts(date);
	const tzPart = parts.find(part => part.type === 'timeZoneName');
	return tzPart?.value || '';
}
