import type {
	CSVRow,
	ImportResult,
	ImportError,
	ImportWarning,
	PlatformSchema,
	Execution,
	GroupedPosition
} from '$lib/types/import';
import type { Trade } from '$lib/types';

/**
 * Parse CSV file content into rows
 */
export function parseCSV(content: string): CSVRow[] {
	const lines = content.trim().split(/\r?\n/);
	if (lines.length < 2) {
		throw new Error('CSV file must contain headers and at least one data row');
	}

	const headers = lines[0].split(',').map((h) => h.trim());
	const rows: CSVRow[] = [];

	for (let i = 1; i < lines.length; i++) {
		const values = parseCSVLine(lines[i]);
		if (values.length !== headers.length) {
			continue; // Skip malformed rows
		}

		const row: CSVRow = {};
		headers.forEach((header, index) => {
			row[header] = values[index]?.trim() || '';
		});
		rows.push(row);
	}

	return rows;
}

/**
 * Parse a CSV line handling quoted fields
 */
function parseCSVLine(line: string): string[] {
	const result: string[] = [];
	let current = '';
	let inQuotes = false;

	for (let i = 0; i < line.length; i++) {
		const char = line[i];

		if (char === '"') {
			inQuotes = !inQuotes;
		} else if (char === ',' && !inQuotes) {
			result.push(current);
			current = '';
		} else {
			current += char;
		}
	}
	result.push(current);

	return result;
}

/**
 * Find column name from schema (supports multiple possible names)
 */
function findColumn(row: CSVRow, possibleNames: string | string[]): string | null {
	const names = Array.isArray(possibleNames) ? possibleNames : [possibleNames];

	for (const name of names) {
		// Case-insensitive match
		const found = Object.keys(row).find((key) => key.toLowerCase() === name.toLowerCase());
		if (found) return found;
	}

	return null;
}

/**
 * Import CSV using platform schema
 */
export async function importCSV(
	file: File,
	platform: PlatformSchema,
	tradingDate?: Date
): Promise<ImportResult> {
	const content = await file.text();
	const rows = parseCSV(content);

	const errors: ImportError[] = [];
	const warnings: ImportWarning[] = [];
	const executions: Execution[] = [];

	// Parse each row into an execution
	for (let i = 0; i < rows.length; i++) {
		const row = rows[i];
		const rowNum = i + 2; // +2 for header row and 0-index

		// Apply row filter if provided (e.g., to filter only "Execute" events for DAS Trader)
		if (platform.rowFilter && !platform.rowFilter(row)) {
			continue; // Skip this row
		}

		try {
			const execution = parseExecution(row, platform, tradingDate, rowNum);
			executions.push(execution);
		} catch (error) {
			errors.push({
				row: rowNum,
				message: error instanceof Error ? error.message : 'Unknown error',
				data: row
			});
		}
	}

	// Group executions into trades if needed
	let trades: Partial<Trade>[];
	if (platform.groupExecutions) {
		trades = groupExecutionsIntoTrades(executions);
	} else {
		trades = executions.map(executionToTrade);
	}

	// Detect duplicates (simple check by timestamp + symbol)
	const duplicateWarnings = detectDuplicates(trades);
	warnings.push(...duplicateWarnings);

	const statistics = {
		totalRows: rows.length,
		validTrades: trades.length,
		duplicates: duplicateWarnings.length,
		errors: errors.length,
		warnings: warnings.length
	};

	return {
		success: errors.length === 0,
		trades,
		errors,
		warnings,
		statistics
	};
}

/**
 * Parse a single CSV row into an Execution
 */
function parseExecution(
	row: CSVRow,
	platform: PlatformSchema,
	tradingDate: Date | undefined,
	rowNum: number
): Execution {
	const { columns, transformations } = platform;

	// Find and extract required fields
	const symbolCol = findColumn(row, columns.symbol);
	const sideCol = findColumn(row, columns.side);
	const priceCol = findColumn(row, columns.price);
	const quantityCol = findColumn(row, columns.quantity);
	const timestampCol = findColumn(row, columns.timestamp);

	if (!symbolCol || !sideCol || !priceCol || !quantityCol || !timestampCol) {
		throw new Error(
			`Missing required columns. Found: ${Object.keys(row).join(', ')}`
		);
	}

	// Transform values
	const symbol = row[symbolCol].toUpperCase();
	const side = transformations.side(row[sideCol]);
	const price = transformations.price(row[priceCol]);
	const quantity = transformations.quantity(row[quantityCol]);
	const timestamp = transformations.timestamp(row[timestampCol], tradingDate);

	// Optional fields
	let fees = 0;
	if (columns.fees && transformations.fees) {
		const feesCol = findColumn(row, columns.fees);
		if (feesCol && row[feesCol]) {
			fees = transformations.fees(row[feesCol]);
		}
	}

	let account: string | undefined;
	if (columns.account) {
		const accountCol = findColumn(row, columns.account);
		if (accountCol) {
			account = row[accountCol];
		}
	}

	return {
		symbol,
		side,
		price,
		quantity,
		timestamp,
		fees,
		account,
		rawData: row
	};
}

/**
 * Group executions by symbol into positions (for DAS Trader, etc.)
 */
function groupExecutionsIntoTrades(executions: Execution[]): Partial<Trade>[] {
	// Sort by timestamp
	const sorted = [...executions].sort((a, b) => a.timestamp.getTime() - b.timestamp.getTime());

	// Group by symbol
	const bySymbol: { [symbol: string]: Execution[] } = {};
	sorted.forEach((exec) => {
		if (!bySymbol[exec.symbol]) {
			bySymbol[exec.symbol] = [];
		}
		bySymbol[exec.symbol].push(exec);
	});

	// Convert each symbol's executions into a trade
	const trades: Partial<Trade>[] = [];

	for (const symbol in bySymbol) {
		const execs = bySymbol[symbol];
		const grouped = groupPositionExecutions(execs);

		grouped.forEach((position) => {
			const trade = positionToTrade(position);
			trades.push(trade);
		});
	}

	return trades;
}

/**
 * Group executions for a single symbol into positions
 * Handles multiple round trips (open/close/open/close)
 */
function groupPositionExecutions(executions: Execution[]): GroupedPosition[] {
	const positions: GroupedPosition[] = [];
	let currentBuys: Execution[] = [];
	let currentSells: Execution[] = [];
	let positionOpen = false;
	let isLong = false;

	for (const exec of executions) {
		if (!positionOpen) {
			// Opening a new position
			positionOpen = true;
			isLong = exec.side === 'B';

			if (isLong) {
				currentBuys.push(exec);
			} else {
				currentSells.push(exec);
			}
		} else {
			// Adding to existing position
			if (exec.side === 'B') {
				currentBuys.push(exec);
			} else {
				currentSells.push(exec);
			}

			// Check if position is closed
			const totalBought = currentBuys.reduce((sum, e) => sum + e.quantity, 0);
			const totalSold = currentSells.reduce((sum, e) => sum + e.quantity, 0);

			if (totalBought === totalSold) {
				// Position closed
				positions.push({
					symbol: executions[0].symbol,
					buys: [...currentBuys],
					sells: [...currentSells],
					firstExecution: currentBuys[0] || currentSells[0],
					tradeType: isLong ? 'LONG' : 'SHORT'
				});

				// Reset for next position
				currentBuys = [];
				currentSells = [];
				positionOpen = false;
			}
		}
	}

	// If position still open, add it
	if (positionOpen && (currentBuys.length > 0 || currentSells.length > 0)) {
		positions.push({
			symbol: executions[0].symbol,
			buys: currentBuys,
			sells: currentSells,
			firstExecution: currentBuys[0] || currentSells[0],
			tradeType: isLong ? 'LONG' : 'SHORT'
		});
	}

	return positions;
}

/**
 * Convert grouped position to Trade object
 */
function positionToTrade(position: GroupedPosition): Partial<Trade> {
	const { symbol, buys, sells, tradeType } = position;

	// Calculate average entry price and total quantity
	const entryExecs = tradeType === 'LONG' ? buys : sells;
	const exitExecs = tradeType === 'LONG' ? sells : buys;

	const totalEntryQty = entryExecs.reduce((sum, e) => sum + e.quantity, 0);
	const totalEntryValue = entryExecs.reduce((sum, e) => sum + e.price * e.quantity, 0);
	const avgEntryPrice = totalEntryValue / totalEntryQty;

	const totalExitQty = exitExecs.reduce((sum, e) => sum + e.quantity, 0);
	const totalExitValue = exitExecs.reduce((sum, e) => sum + e.price * e.quantity, 0);
	const avgExitPrice = totalExitQty > 0 ? totalExitValue / totalExitQty : null;

	// Calculate fees
	const totalFees = [...entryExecs, ...exitExecs].reduce((sum, e) => sum + (e.fees || 0), 0);

	// Calculate P&L
	let pnl: number | null = null;
	if (avgExitPrice !== null && totalExitQty > 0) {
		const priceDiff = tradeType === 'LONG' ? avgExitPrice - avgEntryPrice : avgEntryPrice - avgExitPrice;
		pnl = priceDiff * totalExitQty - totalFees;
	}

	// Build Trade object
	const trade: Partial<Trade> = {
		symbol,
		trade_type: tradeType,
		quantity: totalEntryQty,
		entry_price: avgEntryPrice,
		exit_price: avgExitPrice,
		fees: totalFees,
		pnl,
		entries: entryExecs.map((e) => ({
			id: crypto.randomUUID(),
			price: e.price,
			quantity: e.quantity,
			timestamp: e.timestamp.toISOString(),
			fees: e.fees || 0
		})),
		exits: exitExecs.map((e) => ({
			id: crypto.randomUUID(),
			price: e.price,
			quantity: e.quantity,
			timestamp: e.timestamp.toISOString(),
			fees: e.fees || 0,
			pnl: 0 // TODO: Calculate per-exit P&L
		})),
		current_position_size: totalEntryQty - totalExitQty,
		average_entry_price: avgEntryPrice,
		total_fees: totalFees,
		realized_pnl: pnl || 0,
		unrealized_pnl: null,
		opened_at: entryExecs[0].timestamp.toISOString(),
		closed_at: totalExitQty === totalEntryQty ? exitExecs[exitExecs.length - 1].timestamp.toISOString() : null
	};

	return trade;
}

/**
 * Convert single execution to Trade (for platforms that export positions)
 */
function executionToTrade(execution: Execution): Partial<Trade> {
	const tradeType = execution.side === 'B' ? 'LONG' : 'SHORT';

	return {
		symbol: execution.symbol,
		trade_type: tradeType,
		quantity: execution.quantity,
		entry_price: execution.price,
		exit_price: null,
		fees: execution.fees || 0,
		pnl: null,
		entries: [
			{
				id: crypto.randomUUID(),
				price: execution.price,
				quantity: execution.quantity,
				timestamp: execution.timestamp.toISOString(),
				fees: execution.fees || 0
			}
		],
		exits: [],
		current_position_size: execution.quantity,
		average_entry_price: execution.price,
		total_fees: execution.fees || 0,
		realized_pnl: 0,
		unrealized_pnl: null,
		opened_at: execution.timestamp.toISOString(),
		closed_at: null
	};
}

/**
 * Detect duplicate trades
 */
function detectDuplicates(trades: Partial<Trade>[]): ImportWarning[] {
	const warnings: ImportWarning[] = [];
	const seen = new Set<string>();

	trades.forEach((trade, index) => {
		const key = `${trade.symbol}-${trade.opened_at}-${trade.entry_price}`;
		if (seen.has(key)) {
			warnings.push({
				row: index + 2,
				type: 'duplicate',
				message: `Possible duplicate trade: ${trade.symbol} at ${trade.entry_price}`
			});
		}
		seen.add(key);
	});

	return warnings;
}
