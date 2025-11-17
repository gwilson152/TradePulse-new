import { PUBLIC_API_URL } from '$env/static/public';

const BASE_URL = PUBLIC_API_URL || 'https://api.tradepulse.drivenw.com:9000';

interface APIResponse<T> {
	success: boolean;
	data?: T;
	pagination?: {
		total: number;
		page: number;
		page_size: number;
		total_pages: number;
	};
	error?: {
		code: string;
		message: string;
	};
}

export interface PaginatedResponse<T> {
	data: T;
	pagination: {
		total: number;
		page: number;
		page_size: number;
		total_pages: number;
	};
}

class APIClient {
	private getAuthToken(): string | null {
		if (typeof window === 'undefined') return null;
		return localStorage.getItem('auth_token');
	}

	private setAuthToken(token: string): void {
		if (typeof window === 'undefined') return;
		localStorage.setItem('auth_token', token);
	}

	private removeAuthToken(): void {
		if (typeof window === 'undefined') return;
		localStorage.removeItem('auth_token');
	}

	// Public method to check if user is authenticated
	public getToken(): string | null {
		return this.getAuthToken();
	}

	// Decode JWT token to get user info (without verification)
	public getTokenPayload(): any | null {
		const token = this.getAuthToken();
		if (!token) return null;

		try {
			// JWT format: header.payload.signature
			const parts = token.split('.');
			if (parts.length !== 3) return null;

			// Decode the payload (base64url)
			const payload = parts[1];
			const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'));
			return JSON.parse(decoded);
		} catch (error) {
			console.error('Failed to decode JWT:', error);
			return null;
		}
	}

	async request<T>(
		endpoint: string,
		options?: RequestInit
	): Promise<T> {
		const token = this.getAuthToken();

		const response = await fetch(`${BASE_URL}${endpoint}`, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				...(token && { Authorization: `Bearer ${token}` }),
				...options?.headers
			}
		});

		const result: APIResponse<T> = await response.json();

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data as T;
	}

	async requestWithPagination<T>(
		endpoint: string,
		options?: RequestInit
	): Promise<PaginatedResponse<T>> {
		const token = this.getAuthToken();

		const response = await fetch(`${BASE_URL}${endpoint}`, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				...(token && { Authorization: `Bearer ${token}` }),
				...options?.headers
			}
		});

		const result: APIResponse<T> = await response.json();

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return {
			data: result.data as T,
			pagination: result.pagination || { total: 0, page: 1, page_size: 0, total_pages: 0 }
		};
	}

	// Auth methods
	async requestMagicLink(email: string): Promise<{ message: string }> {
		return this.request('/api/auth/request-magic-link', {
			method: 'POST',
			body: JSON.stringify({ email })
		});
	}

	async signup(email: string, planType: 'starter' | 'pro' | 'premium'): Promise<{ message: string }> {
		return this.request('/api/auth/signup', {
			method: 'POST',
			body: JSON.stringify({ email, plan_type: planType })
		});
	}

	async verifyMagicLink(token: string): Promise<{ jwt: string; user: any }> {
		const result = await this.request<{ jwt: string; user: any }>(
			`/api/auth/verify?token=${token}`
		);
		this.setAuthToken(result.jwt);
		return result;
	}

	async logout(): Promise<void> {
		await this.request('/api/auth/logout', { method: 'POST' });
		this.removeAuthToken();
	}

	async getCurrentUser(): Promise<any> {
		return this.request('/api/auth/me');
	}

	async setPassword(password: string): Promise<{ message: string }> {
		return this.request('/api/auth/set-password', {
			method: 'POST',
			body: JSON.stringify({ password })
		});
	}

	async loginWithPassword(email: string, password: string): Promise<{ jwt: string; user: any }> {
		const result = await this.request<{ jwt: string; user: any }>(
			'/api/auth/login',
			{
				method: 'POST',
				body: JSON.stringify({ email, password })
			}
		);
		this.setAuthToken(result.jwt);
		return result;
	}

	// Trade methods
	async getTrades(params?: {
		limit?: number;
		offset?: number;
		symbol?: string;
		trade_type?: string;
		status?: string;
		start_date?: string;
		end_date?: string;
		strategy?: string;
		min_pnl?: number;
		max_pnl?: number;
	}): Promise<PaginatedResponse<any[]>> {
		const cleanParams: any = {};
		if (params) {
			Object.entries(params).forEach(([key, value]) => {
				if (value !== undefined && value !== null && value !== '') {
					cleanParams[key] = value.toString();
				}
			});
		}
		const query = new URLSearchParams(cleanParams).toString();
		return this.requestWithPagination<any[]>(`/api/trades${query ? '?' + query : ''}`);
	}

	async getTradesForDateRange(from: string, to: string): Promise<any[]> {
		const params = new URLSearchParams({ from, to });
		const result = await this.request<any[]>(`/api/trades?${params.toString()}`);
		return result || [];
	}

	async createTrade(trade: any): Promise<any> {
		return this.request('/api/trades', {
			method: 'POST',
			body: JSON.stringify(trade)
		});
	}

	async updateTrade(id: string, trade: any): Promise<any> {
		return this.request(`/api/trades/${id}`, {
			method: 'PUT',
			body: JSON.stringify(trade)
		});
	}

	async deleteTrade(id: string): Promise<void> {
		return this.request(`/api/trades/${id}`, {
			method: 'DELETE'
		});
	}

	async importTrades(trades: any[]): Promise<{ imported: number; errors: any[] }> {
		return this.request('/api/trades/import-csv', {
			method: 'POST',
			body: JSON.stringify({ trades })
		});
	}

	async fetchPropReportsTrades(
		site: string,
		username: string,
		password: string,
		fromDate?: string,
		toDate?: string
	): Promise<any[]> {
		return this.request('/api/integrations/propreports/fetch', {
			method: 'POST',
			body: JSON.stringify({
				site,
				username,
				password,
				from_date: fromDate,
				to_date: toDate
			})
		});
	}

	// Position lifecycle methods
	async addEntry(tradeId: string, entry: any): Promise<any> {
		return this.request(`/api/trades/${tradeId}/entries`, {
			method: 'POST',
			body: JSON.stringify(entry)
		});
	}

	async addExit(tradeId: string, exit: any): Promise<any> {
		return this.request(`/api/trades/${tradeId}/exits`, {
			method: 'POST',
			body: JSON.stringify(exit)
		});
	}

	// Journal methods
	async getJournalEntries(params?: {
		limit?: number;
		offset?: number;
	}): Promise<{ entries: any[]; total: number }> {
		const query = new URLSearchParams(params as any).toString();
		return this.request(`/api/journal?${query}`);
	}

	async getJournalEntriesByTradeID(tradeId: string): Promise<any[]> {
		const params = new URLSearchParams({ tradeId });
		return this.request(`/api/trades/${tradeId}/journal?${params.toString()}`);
	}

	async createJournalEntry(entry: any): Promise<any> {
		return this.request('/api/journal', {
			method: 'POST',
			body: JSON.stringify(entry)
		});
	}

	// Metrics methods
	async getSummaryMetrics(): Promise<any> {
		return this.request('/api/metrics/summary');
	}

	async getAnalytics(timeRange?: string): Promise<any> {
		const query = timeRange ? `?range=${timeRange}` : '';
		return this.request(`/api/metrics/analytics${query}`);
	}

	// Rule Set methods
	async getRuleSets(): Promise<any[]> {
		return this.request('/api/rulesets');
	}

	async createRuleSet(ruleSet: any): Promise<any> {
		return this.request('/api/rulesets', {
			method: 'POST',
			body: JSON.stringify(ruleSet)
		});
	}

	async updateRuleSet(id: string, ruleSet: any): Promise<any> {
		return this.request(`/api/rulesets/${id}`, {
			method: 'PUT',
			body: JSON.stringify(ruleSet)
		});
	}

	async deleteRuleSet(id: string): Promise<void> {
		return this.request(`/api/rulesets/${id}`, {
			method: 'DELETE'
		});
	}

	async addRule(ruleSetId: string, rule: any): Promise<any> {
		return this.request(`/api/rulesets/${ruleSetId}/rules`, {
			method: 'POST',
			body: JSON.stringify(rule)
		});
	}

	async updateRule(ruleSetId: string, ruleId: string, rule: any): Promise<any> {
		return this.request(`/api/rulesets/${ruleSetId}/rules/${ruleId}`, {
			method: 'PUT',
			body: JSON.stringify(rule)
		});
	}

	async deleteRule(ruleSetId: string, ruleId: string): Promise<void> {
		return this.request(`/api/rulesets/${ruleSetId}/rules/${ruleId}`, {
			method: 'DELETE'
		});
	}

	// File upload methods
	async uploadFile(file: File, type: 'screenshot' | 'voice'): Promise<{ url: string }> {
		const formData = new FormData();
		formData.append('file', file);
		formData.append('type', type);

		const token = this.getAuthToken();
		const response = await fetch(`${BASE_URL}/api/upload`, {
			method: 'POST',
			headers: {
				...(token && { Authorization: `Bearer ${token}` })
			},
			body: formData
		});

		const result = await response.json();
		if (!result.success) {
			throw new Error(result.error?.message || 'Upload failed');
		}

		return result.data;
	}
}

export const apiClient = new APIClient();
