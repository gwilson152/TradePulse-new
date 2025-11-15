import type { PlatformSchema } from '$lib/types/import';

export const dasTraderSchema: PlatformSchema = {
	id: 'das-trader',
	name: 'DAS Trader Pro',
	description: 'Import from DAS Trader Pro CSV export (Trade Log)',
	requiresDate: true, // DAS doesn't include date, only time
	groupExecutions: true, // Need to group executions into positions
	columns: {
		symbol: ['Symb', 'Symbol'],
		side: ['Side'],
		quantity: ['Qty', 'Quantity'],
		price: ['Price', 'Exec Price'],
		timestamp: ['Time'],
		fees: ['Commission', 'Comm'], // Optional
		account: ['Account'],
		orderType: ['Type']
	},
	transformations: {
		side: (value: string): 'B' | 'S' => {
			const normalized = value.trim().toUpperCase();
			if (normalized === 'B' || normalized === 'BUY' || normalized.startsWith('BOT')) {
				return 'B';
			}
			if (normalized === 'S' || normalized === 'SELL' || normalized.startsWith('SOLD')) {
				return 'S';
			}
			throw new Error(`Invalid side value: ${value}`);
		},
		timestamp: (time: string, date?: Date): Date => {
			if (!date) {
				throw new Error('Date is required for DAS Trader imports');
			}

			// Parse time (format: HH:MM:SS or HH:MM:SS.mmm)
			const timeParts = time.split(':');
			if (timeParts.length < 2) {
				throw new Error(`Invalid time format: ${time}`);
			}

			const hours = parseInt(timeParts[0]);
			const minutes = parseInt(timeParts[1]);
			const seconds = timeParts[2] ? parseFloat(timeParts[2]) : 0;

			// Create timestamp using provided date + parsed time
			const timestamp = new Date(date);
			timestamp.setHours(hours, minutes, Math.floor(seconds), (seconds % 1) * 1000);

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
