import type { Trade } from './index';

export interface CSVRow {
	[key: string]: string;
}

export interface ImportResult {
	success: boolean;
	trades: Partial<Trade>[];
	errors: ImportError[];
	warnings: ImportWarning[];
	statistics: ImportStatistics;
}

export interface ImportError {
	row: number;
	column?: string;
	message: string;
	data?: CSVRow;
}

export interface ImportWarning {
	row: number;
	type: 'duplicate' | 'missing_data' | 'unusual_value';
	message: string;
	data?: CSVRow;
}

export interface ImportStatistics {
	totalRows: number;
	validTrades: number;
	duplicates: number;
	errors: number;
	warnings: number;
}

export interface PlatformSchema {
	id: string;
	name: string;
	description: string;
	requiresDate: boolean; // DAS Trader needs date input
	columns: ColumnMapping;
	transformations: TransformFunctions;
	groupExecutions: boolean; // True for platforms that export executions, not positions
}

export interface ColumnMapping {
	symbol: string | string[]; // Column name(s) for symbol
	side: string | string[]; // "Side", "Action", "B/S"
	quantity: string | string[]; // "Qty", "Quantity", "Shares"
	price: string | string[]; // "Price", "Fill Price", "Exec Price"
	timestamp: string | string[]; // "Time", "Date/Time", "Timestamp"
	fees?: string | string[]; // Optional commission column
	account?: string | string[];
	orderType?: string | string[];
}

export interface TransformFunctions {
	side: (value: string) => 'B' | 'S'; // Normalize to B/S
	timestamp: (value: string, date?: Date) => Date;
	price: (value: string) => number;
	quantity: (value: string) => number;
	fees?: (value: string) => number;
}

export interface Execution {
	symbol: string;
	side: 'B' | 'S';
	price: number;
	quantity: number;
	timestamp: Date;
	fees?: number;
	account?: string;
	orderType?: string;
	rawData: CSVRow;
}

export interface GroupedPosition {
	symbol: string;
	buys: Execution[];
	sells: Execution[];
	firstExecution: Execution;
	tradeType: 'LONG' | 'SHORT';
}
