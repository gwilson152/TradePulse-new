import type { Trade, JournalEntry } from '$lib/types';

export interface AnalyticsMetrics {
	winRate: number;
	avgWin: number;
	avgLoss: number;
	riskReward: number;
	totalTrades: number;
	profitFactor: number;
	avgHoldTime: string;
	bestHoldTime: string;
	bestPriceRange: string;
	bestTimeOfDay: string;
}

export interface TimeOfDayData {
	hour: string;
	winRate: number;
	avgPnL: number;
	trades: number;
}

export interface HoldTimeData {
	range: string;
	trades: number;
	avgPnL: number;
	winRate: number;
}

export interface PriceRangeData {
	range: string;
	winRate: number;
	trades: number;
	avgPnL: number;
}

export interface TradeSizeData {
	size: number;
	pnl: number;
	outcome: 'Win' | 'Loss';
}

export interface StreakData {
	streakLength: string;
	winStreaks: number;
	lossStreaks: number;
}

/**
 * Calculate basic performance metrics
 */
export function calculateMetrics(trades: Trade[]): AnalyticsMetrics {
	const closedTrades = trades.filter((t) => t.closed_at !== null);
	const winningTrades = closedTrades.filter((t) => (t.pnl ?? 0) > 0);
	const losingTrades = closedTrades.filter((t) => (t.pnl ?? 0) < 0);

	const totalWins = winningTrades.reduce((sum, t) => sum + (t.pnl ?? 0), 0);
	const totalLosses = Math.abs(losingTrades.reduce((sum, t) => sum + (t.pnl ?? 0), 0));

	const winRate = closedTrades.length > 0 ? (winningTrades.length / closedTrades.length) * 100 : 0;
	const avgWin = winningTrades.length > 0 ? totalWins / winningTrades.length : 0;
	const avgLoss = losingTrades.length > 0 ? totalLosses / losingTrades.length : 0;
	const riskReward = avgLoss > 0 ? avgWin / avgLoss : 0;
	const profitFactor = totalLosses > 0 ? totalWins / totalLosses : totalWins > 0 ? 999 : 0;

	// Calculate average hold time
	const holdTimes = closedTrades
		.map((t) => {
			const opened = new Date(t.opened_at).getTime();
			const closed = new Date(t.closed_at!).getTime();
			return closed - opened;
		})
		.filter((time) => !isNaN(time));

	const avgHoldTimeMs = holdTimes.length > 0 ? holdTimes.reduce((a, b) => a + b, 0) / holdTimes.length : 0;
	const avgHoldTime = formatDuration(avgHoldTimeMs);

	// Find best hold time range
	const holdTimeAnalysis = analyzeHoldTimes(closedTrades);
	const bestHoldTime = findBestHoldTimeRange(holdTimeAnalysis);

	// Find best price range
	const priceRangeAnalysis = analyzePriceRanges(closedTrades);
	const bestPriceRange = findBestPriceRange(priceRangeAnalysis);

	// Find best time of day
	const timeOfDayAnalysis = analyzeTimeOfDay(closedTrades);
	const bestTimeOfDay = findBestTimeOfDay(timeOfDayAnalysis);

	return {
		winRate: Math.round(winRate * 10) / 10,
		avgWin: Math.round(avgWin * 100) / 100,
		avgLoss: -Math.round(avgLoss * 100) / 100,
		riskReward: Math.round(riskReward * 100) / 100,
		totalTrades: closedTrades.length,
		profitFactor: Math.round(profitFactor * 100) / 100,
		avgHoldTime,
		bestHoldTime,
		bestPriceRange,
		bestTimeOfDay
	};
}

/**
 * Analyze performance by time of day (hourly buckets)
 */
export function analyzeTimeOfDay(trades: Trade[]): TimeOfDayData[] {
	const closedTrades = trades.filter((t) => t.closed_at !== null);
	const hourBuckets: { [key: string]: { wins: number; total: number; pnl: number } } = {};

	// Market hours: 9 AM - 4 PM
	for (let hour = 9; hour < 16; hour++) {
		const nextHour = hour + 1;
		const key = `${hour}-${nextHour} ${hour < 12 ? 'AM' : 'PM'}`;
		hourBuckets[key] = { wins: 0, total: 0, pnl: 0 };
	}

	closedTrades.forEach((trade) => {
		const entryHour = new Date(trade.opened_at).getHours();
		if (entryHour >= 9 && entryHour < 16) {
			const nextHour = entryHour + 1;
			const key = `${entryHour}-${nextHour} ${entryHour < 12 ? 'AM' : 'PM'}`;
			if (hourBuckets[key]) {
				hourBuckets[key].total++;
				hourBuckets[key].pnl += trade.pnl ?? 0;
				if ((trade.pnl ?? 0) > 0) hourBuckets[key].wins++;
			}
		}
	});

	return Object.entries(hourBuckets).map(([hour, data]) => ({
		hour,
		winRate: data.total > 0 ? Math.round((data.wins / data.total) * 100) : 0,
		avgPnL: data.total > 0 ? Math.round((data.pnl / data.total) * 100) / 100 : 0,
		trades: data.total
	}));
}

/**
 * Analyze performance by hold time
 */
export function analyzeHoldTimes(trades: Trade[]): HoldTimeData[] {
	const closedTrades = trades.filter((t) => t.closed_at !== null);
	const ranges = [
		{ key: '<30m', min: 0, max: 30 * 60 * 1000 },
		{ key: '30m-1h', min: 30 * 60 * 1000, max: 60 * 60 * 1000 },
		{ key: '1-3h', min: 60 * 60 * 1000, max: 3 * 60 * 60 * 1000 },
		{ key: '3-6h', min: 3 * 60 * 60 * 1000, max: 6 * 60 * 60 * 1000 },
		{ key: '6h-1d', min: 6 * 60 * 60 * 1000, max: 24 * 60 * 60 * 1000 },
		{ key: '1-3d', min: 24 * 60 * 60 * 1000, max: 3 * 24 * 60 * 60 * 1000 },
		{ key: '>3d', min: 3 * 24 * 60 * 60 * 1000, max: Infinity }
	];

	const buckets = ranges.map((range) => ({
		range: range.key,
		trades: 0,
		wins: 0,
		pnl: 0
	}));

	closedTrades.forEach((trade) => {
		const holdTime = new Date(trade.closed_at!).getTime() - new Date(trade.opened_at).getTime();
		const bucket = buckets.find((b, i) => holdTime >= ranges[i].min && holdTime < ranges[i].max);
		if (bucket) {
			bucket.trades++;
			bucket.pnl += trade.pnl ?? 0;
			if ((trade.pnl ?? 0) > 0) bucket.wins++;
		}
	});

	return buckets.map((b) => ({
		range: b.range,
		trades: b.trades,
		avgPnL: b.trades > 0 ? Math.round((b.pnl / b.trades) * 100) / 100 : 0,
		winRate: b.trades > 0 ? Math.round((b.wins / b.trades) * 100) : 0
	}));
}

/**
 * Analyze performance by price range
 */
export function analyzePriceRanges(trades: Trade[]): PriceRangeData[] {
	const closedTrades = trades.filter((t) => t.closed_at !== null);
	const ranges = [
		{ key: '<$20', min: 0, max: 20 },
		{ key: '$20-$50', min: 20, max: 50 },
		{ key: '$50-$100', min: 50, max: 100 },
		{ key: '$100-$200', min: 100, max: 200 },
		{ key: '$200-$500', min: 200, max: 500 },
		{ key: '>$500', min: 500, max: Infinity }
	];

	const buckets = ranges.map((range) => ({
		range: range.key,
		trades: 0,
		wins: 0,
		pnl: 0
	}));

	closedTrades.forEach((trade) => {
		const price = trade.entry_price;
		const bucket = buckets.find((b, i) => price >= ranges[i].min && price < ranges[i].max);
		if (bucket) {
			bucket.trades++;
			bucket.pnl += trade.pnl ?? 0;
			if ((trade.pnl ?? 0) > 0) bucket.wins++;
		}
	});

	return buckets.map((b) => ({
		range: b.range,
		trades: b.trades,
		winRate: b.trades > 0 ? Math.round((b.wins / b.trades) * 100) : 0,
		avgPnL: b.trades > 0 ? Math.round((b.pnl / b.trades) * 100) / 100 : 0
	}));
}

/**
 * Analyze trade sizes vs performance
 */
export function analyzeTradeSizes(trades: Trade[]): TradeSizeData[] {
	return trades
		.filter((t) => t.closed_at !== null)
		.map((trade) => ({
			size: Math.round(trade.entry_price * trade.quantity),
			pnl: Math.round((trade.pnl ?? 0) * 100) / 100,
			outcome: (trade.pnl ?? 0) > 0 ? ('Win' as const) : ('Loss' as const)
		}));
}

/**
 * Analyze win/loss streaks
 */
export function analyzeStreaks(trades: Trade[]): StreakData[] {
	const closedTrades = trades
		.filter((t) => t.closed_at !== null)
		.sort((a, b) => new Date(a.closed_at!).getTime() - new Date(b.closed_at!).getTime());

	const streaks: { type: 'win' | 'loss'; length: number }[] = [];
	let currentStreak: { type: 'win' | 'loss'; length: number } | null = null;

	closedTrades.forEach((trade) => {
		const isWin = (trade.pnl ?? 0) > 0;
		const type = isWin ? 'win' : 'loss';

		if (!currentStreak || currentStreak.type !== type) {
			if (currentStreak) streaks.push(currentStreak);
			currentStreak = { type, length: 1 };
		} else {
			currentStreak.length++;
		}
	});

	if (currentStreak) streaks.push(currentStreak);

	// Bucket by streak length
	const ranges = ['Single', '2 in a row', '3 in a row', '4 in a row', '5+ in a row'];
	return ranges.map((range, i) => {
		const minLength = i + 1;
		const maxLength = i === 4 ? Infinity : i + 1;

		return {
			streakLength: range,
			winStreaks: streaks.filter((s) => s.type === 'win' && s.length >= minLength && s.length <= maxLength).length,
			lossStreaks: streaks.filter((s) => s.type === 'loss' && s.length >= minLength && s.length <= maxLength).length
		};
	});
}

/**
 * Helper function to format duration
 */
function formatDuration(ms: number): string {
	const hours = Math.floor(ms / (1000 * 60 * 60));
	const minutes = Math.floor((ms % (1000 * 60 * 60)) / (1000 * 60));

	if (hours === 0) return `${minutes}m`;
	return `${hours}h ${minutes}m`;
}

/**
 * Find best hold time range based on highest win rate
 */
function findBestHoldTimeRange(data: HoldTimeData[]): string {
	const best = data.reduce((prev, curr) => (curr.winRate > prev.winRate ? curr : prev), data[0]);
	return best?.range ?? '1-3h';
}

/**
 * Find best price range based on highest win rate
 */
function findBestPriceRange(data: PriceRangeData[]): string {
	const best = data.reduce((prev, curr) => (curr.winRate > prev.winRate ? curr : prev), data[0]);
	return best?.range ?? '$50-$100';
}

/**
 * Find best time of day based on highest win rate
 */
function findBestTimeOfDay(data: TimeOfDayData[]): string {
	const best = data.reduce((prev, curr) => (curr.winRate > prev.winRate && curr.trades > 0 ? curr : prev), data[0]);
	return best?.hour ?? '10-11 AM';
}
