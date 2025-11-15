import type { PlatformSchema } from '$lib/types/import';

export const propReportsSchema: PlatformSchema = {
	id: 'prop-reports',
	name: 'PropReports',
	description: 'Import from PropReports detailed trade export',
	requiresDate: false, // PropReports includes full date/time
	groupExecutions: false, // PropReports exports positions, not individual executions
	columns: {
		symbol: ['Symbol', 'Ticker', 'Instrument'],
		side: ['Side', 'Direction', 'Type'],
		quantity: ['Quantity', 'Qty', 'Shares', 'Size'],
		price: ['Price', 'Entry Price', 'Avg Price', 'Average Price'],
		timestamp: ['Date', 'Time', 'DateTime', 'Timestamp', 'Entry Time'],
		fees: ['Commission', 'Comm', 'Fees', 'Total Fees'],
		account: ['Account', 'Account Number']
	},
	transformations: {
		side: (value: string): 'B' | 'S' => {
			const normalized = value.trim().toUpperCase();
			if (normalized === 'LONG' || normalized === 'BUY' || normalized === 'B') {
				return 'B';
			}
			if (normalized === 'SHORT' || normalized === 'SELL' || normalized === 'S') {
				return 'S';
			}
			throw new Error(`Invalid side value: ${value}`);
		},
		timestamp: (value: string, date?: Date): Date => {
			// Try parsing various date formats
			const timestamp = new Date(value);
			if (isNaN(timestamp.getTime())) {
				throw new Error(`Invalid timestamp: ${value}`);
			}
			return timestamp;
		},
		price: (value: string): number => {
			const price = parseFloat(value.replace(/[$,]/g, ''));
			if (isNaN(price)) {
				throw new Error(`Invalid price: ${value}`);
			}
			return price;
		},
		quantity: (value: string): number => {
			const qty = parseInt(value.replace(/,/g, ''));
			if (isNaN(qty) || qty <= 0) {
				throw new Error(`Invalid quantity: ${value}`);
			}
			return qty;
		},
		fees: (value: string): number => {
			const fees = parseFloat(value.replace(/[$,]/g, ''));
			return isNaN(fees) ? 0 : Math.abs(fees);
		}
	}
};
