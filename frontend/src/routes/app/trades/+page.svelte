<script lang="ts">
	import { onMount } from 'svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Table from '$lib/components/ui/Table.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Icon from '@iconify/svelte';
	import PnLBadge from '$lib/components/trading/PnLBadge.svelte';
	import TradeFormSlideOver from '$lib/components/trading/TradeFormSlideOver.svelte';
	import TradeDetailSlideOver from '$lib/components/trading/TradeDetailSlideOver.svelte';
	import CSVImportSlideOver from '$lib/components/trading/CSVImportSlideOver.svelte';
	import JournalFormSlideOver from '$lib/components/trading/JournalFormSlideOver.svelte';
	import { apiClient } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import type { Trade, Rule } from '$lib/types';

	let trades = $state<Trade[]>([]);
	let rules = $state<Rule[]>([]);
	let isLoading = $state(true);
	let showAddModal = $state(false);
	let showImportModal = $state(false);
	let showDetailSlideOver = $state(false);
	let showJournalModal = $state(false);
	let selectedTrade = $state<Trade | null>(null);
	let searchQuery = $state('');
	let filterStatus = $state('all'); // all, open, closed, profitable, loss
	let filterTradeType = $state('all'); // all, LONG, SHORT
	let filterDateRange = $state('all'); // all, today, week, month, quarter, year
	let filterStrategy = $state('all'); // all, or specific strategy
	let minPnL = $state<number | null>(null);
	let maxPnL = $state<number | null>(null);
	let showAdvancedFilters = $state(false);
	let hoveredTrade = $state<Trade | null>(null);
	let tooltipX = $state(0);
	let tooltipY = $state(0);
	let longPressTimer: number | null = null;

	// Server-side pagination metadata
	let totalTrades = $state(0);
	let totalPages = $state(1);
	let currentPage = $state(1);
	let pageSize = $state(25); // Default 25 trades per page
	const pageSizeOptions = [10, 25, 50, 100];

	onMount(async () => {
		await loadTrades();
		await loadRules();
	});

	async function loadRules() {
		try {
			const ruleSets = await apiClient.getRuleSets();
			// Get all rules from active rule sets
			rules = ruleSets
				.filter(rs => rs.is_active)
				.flatMap(rs => rs.rules || []);
		} catch (err) {
			console.error('Failed to load rules:', err);
			// Don't show error toast for rules, it's optional
		}
	}

	async function loadJournalEntries() {
		try {
			// Fetch journal entries for each trade
			const journalPromises = trades.map(async (trade) => {
				try {
					const entries = await apiClient.getJournalEntriesByTradeID(trade.id);
					return { tradeId: trade.id, entries };
				} catch (err) {
					// Silently fail for individual trades
					return { tradeId: trade.id, entries: [] };
				}
			});

			const results = await Promise.all(journalPromises);

			// Update trades with journal entries
			trades = trades.map(trade => {
				const result = results.find(r => r.tradeId === trade.id);
				return {
					...trade,
					journal_entries: result?.entries || [],
					has_journal: (result?.entries?.length || 0) > 0
				};
			});
		} catch (err) {
			console.error('Failed to load journal entries:', err);
			// Don't show error toast, this is optional
		}
	}

	async function loadTrades() {
		try {
			isLoading = true;

			// Build API parameters
			const params: any = {
				limit: pageSize,
				offset: (currentPage - 1) * pageSize
			};

			// Add search
			if (searchQuery) {
				params.symbol = searchQuery;
			}

			// Add filters
			if (filterStatus === 'open') {
				params.status = 'open';
			} else if (filterStatus === 'closed') {
				params.status = 'closed';
			} else if (filterStatus === 'profitable') {
				params.status = 'closed';
				params.min_pnl = 0.01;
			} else if (filterStatus === 'loss') {
				params.status = 'closed';
				params.max_pnl = -0.01;
			}

			if (filterTradeType !== 'all') {
				params.trade_type = filterTradeType;
			}

			if (filterStrategy !== 'all') {
				params.strategy = filterStrategy;
			}

			if (filterDateRange !== 'all') {
				const rangeStart = getDateRangeStart(filterDateRange);
				if (rangeStart) {
					params.start_date = rangeStart.toISOString();
				}
			}

			if (minPnL !== null) {
				params.min_pnl = minPnL;
			}

			if (maxPnL !== null) {
				params.max_pnl = maxPnL;
			}

			const response = await apiClient.getTrades(params);
			trades = response.data;
			totalTrades = response.pagination.total;
			totalPages = response.pagination.total_pages;

			// Fetch journal entries for each trade
			await loadJournalEntries();
		} catch (err) {
			console.error('Failed to load trades:', err);
			toast.error('Failed to load trades');
		} finally {
			isLoading = false;
		}
	}

	async function handleAddTrade(tradeData: any) {
		console.log('Adding trade:', tradeData);
		// TODO: Call API to create trade
		showAddModal = false;
		await loadTrades();
	}

	async function handleImportTrades(importedTrades: Partial<Trade>[]) {
		try {
			const result = await apiClient.importTrades(importedTrades);
			console.log('Import result:', result);
			toast.success(`Imported ${result.imported} trades successfully`);
			await loadTrades();
		} catch (err) {
			console.error('Failed to import trades:', err);
			toast.error('Failed to import trades');
			throw err;
		}
	}

	function handleViewTrade(trade: Trade) {
		selectedTrade = trade;
		showDetailSlideOver = true;
	}

	function handleEditTrade() {
		showDetailSlideOver = false;
		showAddModal = true;
	}

	function handleDeleteTrade() {
		console.log('Deleting trade:', selectedTrade);
		// TODO: Call API to delete trade
		showDetailSlideOver = false;
	}

	function handleCreateJournal(trade: Trade) {
		selectedTrade = trade;
		showJournalModal = true;
	}

	async function handleJournalSubmit(data: any, screenshots: File[], voiceNotes: Blob[]) {
		try {
			// Upload files first
			const screenshotUrls: string[] = [];
			const voiceNoteUrls: string[] = [];

			for (const file of screenshots) {
				const result = await apiClient.uploadFile(file, 'screenshot');
				screenshotUrls.push(result.url);
			}

			for (const blob of voiceNotes) {
				const file = new File([blob], `voice-${Date.now()}.webm`, { type: 'audio/webm' });
				const result = await apiClient.uploadFile(file, 'voice');
				voiceNoteUrls.push(result.url);
			}

			// Create journal entry with file URLs
			const journalEntry = {
				...data,
				screenshots: screenshotUrls,
				voice_notes: voiceNoteUrls
			};

			await apiClient.createJournalEntry(journalEntry);
			toast.success('Journal entry created successfully');
		} catch (err) {
			console.error('Failed to create journal entry:', err);
			toast.error('Failed to create journal entry');
			throw err;
		}
	}

	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function formatDateTime(dateString: string) {
		const date = new Date(dateString);
		return date.toLocaleString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			hour12: true
		});
	}

	function formatCurrency(value: number) {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(value);
	}

	function handleMouseMove(e: MouseEvent, trade: Trade) {
		hoveredTrade = trade;
		tooltipX = e.clientX;
		tooltipY = e.clientY;
	}

	function handleMouseLeave() {
		hoveredTrade = null;
	}

	function handleTouchStart(e: TouchEvent, trade: Trade) {
		longPressTimer = window.setTimeout(() => {
			const touch = e.touches[0];
			tooltipX = touch.clientX;
			tooltipY = touch.clientY;
			hoveredTrade = trade;
		}, 500); // 500ms long press
	}

	function handleTouchEnd() {
		if (longPressTimer) {
			clearTimeout(longPressTimer);
			longPressTimer = null;
		}
	}

	function handleTouchMove() {
		if (longPressTimer) {
			clearTimeout(longPressTimer);
			longPressTimer = null;
		}
	}

	// Get unique strategies from trades (for filter dropdown)
	const strategies = $derived(() => {
		if (!trades || trades.length === 0) return [];
		const uniqueStrategies = new Set(trades.filter(t => t.strategy).map(t => t.strategy));
		return Array.from(uniqueStrategies).sort();
	});

	function clearFilters() {
		searchQuery = '';
		filterStatus = 'all';
		filterTradeType = 'all';
		filterDateRange = 'all';
		filterStrategy = 'all';
		minPnL = null;
		maxPnL = null;
		currentPage = 1;
		loadTrades();
	}

	function getDateRangeStart(range: string): Date | null {
		const now = new Date();
		switch (range) {
			case 'today':
				return new Date(now.setHours(0, 0, 0, 0));
			case 'week':
				const weekAgo = new Date(now);
				weekAgo.setDate(weekAgo.getDate() - 7);
				return weekAgo;
			case 'month':
				const monthAgo = new Date(now);
				monthAgo.setMonth(monthAgo.getMonth() - 1);
				return monthAgo;
			case 'quarter':
				const quarterAgo = new Date(now);
				quarterAgo.setMonth(quarterAgo.getMonth() - 3);
				return quarterAgo;
			case 'year':
				const yearAgo = new Date(now);
				yearAgo.setFullYear(yearAgo.getFullYear() - 1);
				return yearAgo;
			default:
				return null;
		}
	}

	// Count active filters
	const activeFilterCount = $derived(() => {
		let count = 0;
		if (filterStatus !== 'all') count++;
		if (filterTradeType !== 'all') count++;
		if (filterDateRange !== 'all') count++;
		if (filterStrategy !== 'all') count++;
		if (minPnL !== null) count++;
		if (maxPnL !== null) count++;
		if (searchQuery) count++;
		return count;
	});

	// Watch for filter changes and reload
	let lastFilterState = $state('');

	$effect(() => {
		// Create a string representation of filter state
		const filterState = JSON.stringify({
			filterStatus,
			filterTradeType,
			filterDateRange,
			filterStrategy,
			minPnL,
			maxPnL,
			searchQuery,
			pageSize
		});

		// Skip initial load (handled by onMount)
		if (lastFilterState === '') {
			lastFilterState = filterState;
			return;
		}

		// If filters changed, reset to page 1 and reload
		if (filterState !== lastFilterState) {
			lastFilterState = filterState;
			currentPage = 1;
			loadTrades();
		}
	});

	// Watch for page changes and reload
	let lastPage = $state(0);

	$effect(() => {
		// Skip initial load
		if (lastPage === 0) {
			lastPage = currentPage;
			return;
		}

		// If page changed, reload
		if (currentPage !== lastPage) {
			lastPage = currentPage;
			loadTrades();
		}
	});

	function goToPage(page: number) {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
		}
	}

	function changePageSize(newSize: number) {
		pageSize = newSize;
		currentPage = 1;
	}

	// Generate page numbers for pagination controls
	const pageNumbers = $derived(() => {
		const total = totalPages;
		const current = currentPage;
		const pages: (number | string)[] = [];

		if (total <= 7) {
			// Show all pages if 7 or fewer
			for (let i = 1; i <= total; i++) {
				pages.push(i);
			}
		} else {
			// Always show first page
			pages.push(1);

			if (current > 3) {
				pages.push('...');
			}

			// Show pages around current page
			const start = Math.max(2, current - 1);
			const end = Math.min(total - 1, current + 1);

			for (let i = start; i <= end; i++) {
				pages.push(i);
			}

			if (current < total - 2) {
				pages.push('...');
			}

			// Always show last page
			pages.push(total);
		}

		return pages;
	});
</script>

<svelte:head>
	<title>Trades - TradePulse</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold mb-2">Trades</h1>
			<p class="text-surface-600 dark:text-surface-400">Import trades or add them manually</p>
		</div>
		<div class="flex gap-3">
			<Button color="soft" variant="soft" onclick={() => (showAddModal = true)}>
				<Icon icon="mdi:plus" width="20" class="mr-2" />
				Add Trade
			</Button>
			<Button color="primary" onclick={() => (showImportModal = true)}>
				<Icon icon="mdi:file-upload" width="20" class="mr-2" />
				Import
			</Button>
		</div>
	</div>

	<!-- Filters -->
	<Card>
		<div class="space-y-4">
			<!-- Top Row: Search and Quick Filters -->
			<div class="flex flex-col md:flex-row gap-4">
				<!-- Search -->
				<div class="flex-1">
					<div class="relative">
						<Icon
							icon="mdi:magnify"
							width="20"
							class="absolute left-3 top-1/2 -translate-y-1/2 text-surface-400"
						/>
						<input
							type="text"
							class="input pl-10 w-full"
							placeholder="Search by symbol..."
							bind:value={searchQuery}
						/>
					</div>
				</div>

				<!-- Quick Action Buttons -->
				<div class="flex gap-2">
					<button
						onclick={() => showAdvancedFilters = !showAdvancedFilters}
						class="px-4 py-2 rounded-xl text-sm font-medium transition-all hover:scale-105 flex items-center gap-2
							{showAdvancedFilters ? 'bg-blue-500 text-white shadow-md' : 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300'}"
					>
						<Icon icon="mdi:filter-cog" width="16" />
						Advanced
						{#if activeFilterCount() > 0}
							<span class="px-1.5 py-0.5 rounded-full bg-white/20 text-xs font-bold">{activeFilterCount()}</span>
						{/if}
					</button>
					{#if activeFilterCount() > 0}
						<button
							onclick={clearFilters}
							class="px-4 py-2 rounded-xl text-sm font-medium transition-all hover:scale-105 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400"
						>
							<Icon icon="mdi:filter-off" width="16" class="inline mr-1.5" />
							Clear
						</button>
					{/if}
				</div>
			</div>

			<!-- Status Filter Chips -->
			<div class="flex gap-2 flex-wrap">
				<button
					class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all"
					class:bg-gradient-to-r={filterStatus === 'all'}
					class:from-blue-500={filterStatus === 'all'}
					class:to-purple-600={filterStatus === 'all'}
					class:text-white={filterStatus === 'all'}
					class:shadow-md={filterStatus === 'all'}
					class:bg-slate-100={filterStatus !== 'all'}
					class:dark:bg-slate-800={filterStatus !== 'all'}
					class:text-slate-700={filterStatus !== 'all'}
					class:dark:text-slate-300={filterStatus !== 'all'}
					onclick={() => (filterStatus = 'all')}
				>
					<Icon icon="mdi:view-grid" width="16" class="inline mr-1" />
					All ({totalTrades})
				</button>
				<button
					class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all"
					class:bg-blue-500={filterStatus === 'open'}
					class:text-white={filterStatus === 'open'}
					class:shadow-md={filterStatus === 'open'}
					class:bg-slate-100={filterStatus !== 'open'}
					class:dark:bg-slate-800={filterStatus !== 'open'}
					class:text-slate-700={filterStatus !== 'open'}
					class:dark:text-slate-300={filterStatus !== 'open'}
					onclick={() => (filterStatus = 'open')}
				>
					<Icon icon="mdi:clock-outline" width="16" class="inline mr-1" />
					Open
				</button>
				<button
					class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all"
					class:bg-slate-600={filterStatus === 'closed'}
					class:text-white={filterStatus === 'closed'}
					class:shadow-md={filterStatus === 'closed'}
					class:bg-slate-100={filterStatus !== 'closed'}
					class:dark:bg-slate-800={filterStatus !== 'closed'}
					class:text-slate-700={filterStatus !== 'closed'}
					class:dark:text-slate-300={filterStatus !== 'closed'}
					onclick={() => (filterStatus = 'closed')}
				>
					<Icon icon="mdi:check-circle" width="16" class="inline mr-1" />
					Closed
				</button>
				<button
					class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all"
					class:bg-emerald-500={filterStatus === 'profitable'}
					class:text-white={filterStatus === 'profitable'}
					class:shadow-md={filterStatus === 'profitable'}
					class:bg-slate-100={filterStatus !== 'profitable'}
					class:dark:bg-slate-800={filterStatus !== 'profitable'}
					class:text-slate-700={filterStatus !== 'profitable'}
					class:dark:text-slate-300={filterStatus !== 'profitable'}
					onclick={() => (filterStatus = 'profitable')}
				>
					<Icon icon="mdi:trending-up" width="16" class="inline mr-1" />
					Winners
				</button>
				<button
					class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all"
					class:bg-red-500={filterStatus === 'loss'}
					class:text-white={filterStatus === 'loss'}
					class:shadow-md={filterStatus === 'loss'}
					class:bg-slate-100={filterStatus !== 'loss'}
					class:dark:bg-slate-800={filterStatus !== 'loss'}
					class:text-slate-700={filterStatus !== 'loss'}
					class:dark:text-slate-300={filterStatus !== 'loss'}
					onclick={() => (filterStatus = 'loss')}
				>
					<Icon icon="mdi:trending-down" width="16" class="inline mr-1" />
					Losers
				</button>
			</div>

			<!-- Advanced Filters Panel -->
			{#if showAdvancedFilters}
				<div class="pt-4 border-t border-slate-200 dark:border-slate-700 space-y-4 animate-in slide-in-from-top">
					<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
						<!-- Trade Type Filter -->
						<div>
							<label for="trade-type-filter" class="block text-xs font-semibold text-slate-600 dark:text-slate-400 mb-2 uppercase tracking-wide">
								Trade Type
							</label>
							<select
								id="trade-type-filter"
								bind:value={filterTradeType}
								class="w-full px-3 py-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 text-sm focus:ring-2 focus:ring-blue-500"
							>
								<option value="all">All Types</option>
								<option value="LONG">LONG</option>
								<option value="SHORT">SHORT</option>
							</select>
						</div>

						<!-- Date Range Filter -->
						<div>
							<label for="date-range-filter" class="block text-xs font-semibold text-slate-600 dark:text-slate-400 mb-2 uppercase tracking-wide">
								Date Range
							</label>
							<select
								id="date-range-filter"
								bind:value={filterDateRange}
								class="w-full px-3 py-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 text-sm focus:ring-2 focus:ring-blue-500"
							>
								<option value="all">All Time</option>
								<option value="today">Today</option>
								<option value="week">Last 7 Days</option>
								<option value="month">Last 30 Days</option>
								<option value="quarter">Last 90 Days</option>
								<option value="year">Last Year</option>
							</select>
						</div>

						<!-- Strategy Filter -->
						<div>
							<label for="strategy-filter" class="block text-xs font-semibold text-slate-600 dark:text-slate-400 mb-2 uppercase tracking-wide">
								Strategy
							</label>
							<select
								id="strategy-filter"
								bind:value={filterStrategy}
								class="w-full px-3 py-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 text-sm focus:ring-2 focus:ring-blue-500"
							>
								<option value="all">All Strategies</option>
								{#each strategies() as strategy}
									<option value={strategy}>{strategy}</option>
								{/each}
							</select>
						</div>

						<!-- Results Count -->
						<div class="flex items-end">
							<div class="w-full px-4 py-2 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
								<div class="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">
									Showing Results
								</div>
								<div class="text-2xl font-bold text-blue-700 dark:text-blue-300">
									{totalTrades}
								</div>
							</div>
						</div>
					</div>

					<!-- P&L Range Filter -->
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div>
							<label for="min-pnl" class="block text-xs font-semibold text-slate-600 dark:text-slate-400 mb-2 uppercase tracking-wide">
								Min P&L ($)
							</label>
							<input
								id="min-pnl"
								type="number"
								placeholder="e.g. -1000"
								bind:value={minPnL}
								class="w-full px-3 py-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 text-sm focus:ring-2 focus:ring-blue-500"
							/>
						</div>
						<div>
							<label for="max-pnl" class="block text-xs font-semibold text-slate-600 dark:text-slate-400 mb-2 uppercase tracking-wide">
								Max P&L ($)
							</label>
							<input
								id="max-pnl"
								type="number"
								placeholder="e.g. 5000"
								bind:value={maxPnL}
								class="w-full px-3 py-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 text-sm focus:ring-2 focus:ring-blue-500"
							/>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</Card>

	<!-- Trades Table -->
	{#if isLoading}
		<Card>
			<div class="text-center py-12">
				<Icon icon="mdi:loading" width="48" class="animate-spin text-primary-600 mx-auto mb-4" />
				<p class="text-surface-600 dark:text-surface-400">Loading trades...</p>
			</div>
		</Card>
	{:else if !trades || trades.length === 0}
		<div class="bg-white/80 dark:bg-slate-800/80 backdrop-blur-xl rounded-2xl border border-slate-200/50 dark:border-slate-700/50 p-12">
			<div class="max-w-md mx-auto text-center">
				<div class="mb-6 relative">
					<div class="w-24 h-24 mx-auto rounded-full bg-gradient-to-br from-blue-100 to-purple-100 dark:from-blue-900/30 dark:to-purple-900/30 flex items-center justify-center">
						<Icon icon="mdi:chart-line-variant" width="48" class="text-blue-600 dark:text-blue-400" />
					</div>
					<div class="absolute top-0 right-1/2 translate-x-12 -translate-y-2">
						<div class="w-6 h-6 rounded-full bg-yellow-400 flex items-center justify-center">
							<Icon icon="mdi:star" width="14" class="text-white" />
						</div>
					</div>
				</div>
				<h3 class="text-2xl font-bold text-slate-800 dark:text-slate-100 mb-3">
					{#if searchQuery || filterStatus !== 'all'}
						No trades match your filters
					{:else}
						Start Your Trading Journey
					{/if}
				</h3>
				<p class="text-slate-600 dark:text-slate-400 mb-8">
					{#if searchQuery || filterStatus !== 'all'}
						Try adjusting your search or filter criteria
					{:else}
						Import your trade history or manually add your first trade to begin tracking your performance
					{/if}
				</p>
				{#if !searchQuery && filterStatus === 'all'}
					<div class="flex flex-col sm:flex-row gap-3 justify-center">
						<Button color="primary" size="lg" onclick={() => (showImportModal = true)}>
							<Icon icon="mdi:file-upload" width="24" class="mr-2" />
							Import Trades
						</Button>
						<Button variant="soft" color="secondary" size="lg" onclick={() => (showAddModal = true)}>
							<Icon icon="mdi:plus" width="20" class="mr-2" />
							Add Manually
						</Button>
					</div>
				{:else}
					<Button variant="soft" color="secondary" onclick={() => { searchQuery = ''; filterStatus = 'all'; }}>
						<Icon icon="mdi:filter-off" width="20" class="mr-2" />
						Clear Filters
					</Button>
				{/if}
			</div>
		</div>
	{:else}
		<Card>
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead>
						<tr class="border-b border-slate-200 dark:border-slate-700">
							<th class="text-left py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">Symbol</th>
							<th class="text-left py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">Type</th>
							<th class="text-left py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">Date & Time</th>
							<th class="text-right py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">Qty</th>
							<th class="text-right py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">Entry</th>
							<th class="text-right py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">Exit</th>
							<th class="text-right py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">P&L</th>
							<th class="text-center py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">Status</th>
							<th class="text-center py-3 px-4 text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase tracking-wider">Journal</th>
							<th class="py-3 px-4"></th>
						</tr>
					</thead>
					<tbody>
						{#each trades as trade}
							<tr
								role="button"
								tabindex="0"
								onclick={() => handleViewTrade(trade)}
								onkeydown={(e) => e.key === 'Enter' && handleViewTrade(trade)}
								onmousemove={(e) => handleMouseMove(e, trade)}
								onmouseleave={handleMouseLeave}
								ontouchstart={(e) => handleTouchStart(e, trade)}
								ontouchend={handleTouchEnd}
								ontouchmove={handleTouchMove}
								class="relative border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 cursor-pointer transition-colors group"
							>
								<td class="py-2.5 px-4">
									<div class="flex items-center gap-2">
										<div class="w-8 h-8 rounded-md bg-gradient-to-br {trade.trade_type === 'LONG' ? 'from-emerald-400 to-emerald-600' : 'from-red-400 to-red-600'} flex items-center justify-center flex-shrink-0">
											<Icon icon={trade.trade_type === 'LONG' ? 'mdi:arrow-up-bold' : 'mdi:arrow-down-bold'} width="16" class="text-white" />
										</div>
										<span class="font-semibold text-slate-800 dark:text-slate-100">{trade.symbol}</span>
									</div>
								</td>
								<td class="py-2.5 px-4">
									<span class="px-2 py-1 rounded-md text-xs font-medium {trade.trade_type === 'LONG' ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400' : 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'}">
										{trade.trade_type}
									</span>
								</td>
								<td class="py-2.5 px-4 text-xs text-slate-600 dark:text-slate-400">
									<div>{formatDate(trade.opened_at)}</div>
									<div class="text-slate-500 dark:text-slate-500">{new Date(trade.opened_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })}</div>
								</td>
								<td class="py-2.5 px-4 text-right text-sm font-medium text-slate-800 dark:text-slate-200">
									{trade.quantity}
								</td>
								<td class="py-2.5 px-4 text-right text-sm font-medium text-slate-800 dark:text-slate-200">
									{formatCurrency(trade.entry_price)}
								</td>
								<td class="py-2.5 px-4 text-right text-sm font-medium text-slate-800 dark:text-slate-200">
									{trade.exit_price ? formatCurrency(trade.exit_price) : '-'}
								</td>
								<td class="py-2.5 px-4 text-right">
									{#if trade.pnl !== null}
										<span class="font-bold {trade.pnl >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}">
											{trade.pnl >= 0 ? '+' : ''}{formatCurrency(trade.pnl)}
										</span>
									{:else}
										<span class="text-slate-500 dark:text-slate-400">-</span>
									{/if}
								</td>
								<td class="py-2.5 px-4 text-center">
									{#if !trade.exit_price}
										<span class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400">
											<Icon icon="mdi:clock-outline" width="12" />
											Open
										</span>
									{:else}
										<span class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400">
											<Icon icon="mdi:check" width="12" />
											Closed
										</span>
									{/if}
								</td>
								<td class="py-2.5 px-4 text-center">
									<button
										onclick={(e) => {
											e.stopPropagation();
											handleCreateJournal(trade);
										}}
										class="p-2 rounded-lg transition-colors group/journal relative
											{trade.has_journal
												? 'hover:bg-emerald-100 dark:hover:bg-emerald-900/30'
												: 'hover:bg-blue-100 dark:hover:bg-blue-900/30'}"
										title={trade.has_journal ? `${trade.journal_entries?.length || 0} journal ${trade.journal_entries?.length === 1 ? 'entry' : 'entries'}` : 'Create journal entry'}
									>
										<Icon
											icon={trade.has_journal ? 'mdi:book-check' : 'mdi:book-edit'}
											width="20"
											class="transition-colors
												{trade.has_journal
													? 'text-emerald-600 dark:text-emerald-400 group-hover/journal:text-emerald-700 dark:group-hover/journal:text-emerald-300'
													: 'text-slate-400 group-hover/journal:text-blue-500'}"
										/>
										{#if trade.has_journal && (trade.journal_entries?.length || 0) > 1}
											<span class="absolute -top-1 -right-1 w-4 h-4 bg-emerald-600 dark:bg-emerald-500 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
												{trade.journal_entries?.length}
											</span>
										{/if}
									</button>
								</td>
								<td class="py-2.5 px-4">
									<Icon icon="mdi:chevron-right" width="20" class="text-slate-400 group-hover:text-blue-500 transition-colors" />
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<!-- Pagination Controls -->
			{#if totalPages > 1}
				<div class="flex items-center justify-between px-4 py-3 border-t border-slate-200 dark:border-slate-700">
					<!-- Page Size Selector -->
					<div class="flex items-center gap-2">
						<span class="text-sm text-slate-600 dark:text-slate-400">Show</span>
						<select
							bind:value={pageSize}
							onchange={() => changePageSize(pageSize)}
							class="px-2 py-1 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 text-sm focus:ring-2 focus:ring-blue-500"
						>
							{#each pageSizeOptions as size}
								<option value={size}>{size}</option>
							{/each}
						</select>
						<span class="text-sm text-slate-600 dark:text-slate-400">per page</span>
					</div>

					<!-- Pagination Info and Controls -->
					<div class="flex items-center gap-4">
						<!-- Results info -->
						<span class="text-sm text-slate-600 dark:text-slate-400">
							{((currentPage - 1) * pageSize) + 1}-{Math.min(currentPage * pageSize, totalTrades)} of {totalTrades}
						</span>

						<!-- Page controls -->
						<div class="flex items-center gap-1">
							<!-- Previous button -->
							<button
								onclick={() => goToPage(currentPage - 1)}
								disabled={currentPage === 1}
								class="p-2 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed
									{currentPage === 1 ? 'text-slate-400 dark:text-slate-600' : 'text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'}"
							>
								<Icon icon="mdi:chevron-left" width="20" />
							</button>

							<!-- Page numbers -->
							{#each pageNumbers() as page}
								{#if page === '...'}
									<span class="px-2 text-slate-400">...</span>
								{:else}
									<button
										onclick={() => goToPage(page)}
										class="min-w-[2rem] px-3 py-1.5 rounded-lg text-sm font-medium transition-colors
											{currentPage === page
												? 'bg-blue-500 text-white shadow-md'
												: 'text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'}"
									>
										{page}
									</button>
								{/if}
							{/each}

							<!-- Next button -->
							<button
								onclick={() => goToPage(currentPage + 1)}
								disabled={currentPage === totalPages}
								class="p-2 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed
									{currentPage === totalPages ? 'text-slate-400 dark:text-slate-600' : 'text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'}"
							>
								<Icon icon="mdi:chevron-right" width="20" />
							</button>
						</div>
					</div>
				</div>
			{/if}
		</Card>
	{/if}
</div>

<!-- Tooltip Card at Mouse Location -->
{#if hoveredTrade}
	<div
		role="button"
		tabindex="-1"
		class="fixed z-50 pointer-events-none"
		style="left: {tooltipX + 15}px; top: {tooltipY + 15}px;"
		onclick={() => hoveredTrade = null}
		onkeydown={(e) => e.key === 'Escape' && (hoveredTrade = null)}
	>
		<div class="bg-slate-900 dark:bg-slate-950 text-white rounded-lg shadow-2xl p-3 border border-slate-700 text-xs max-w-md">
			<div class="font-semibold text-sm mb-2 border-b border-slate-700 pb-2">
				{hoveredTrade.symbol} - Trade Details
			</div>
			<div class="grid grid-cols-2 gap-x-4 gap-y-2">
				{#if hoveredTrade.stop_loss}
					<div>
						<span class="text-slate-400">Stop Loss:</span>
						<span class="font-semibold ml-1.5">{formatCurrency(hoveredTrade.stop_loss)}</span>
					</div>
				{/if}
				{#if hoveredTrade.take_profit}
					<div>
						<span class="text-slate-400">Take Profit:</span>
						<span class="font-semibold ml-1.5">{formatCurrency(hoveredTrade.take_profit)}</span>
					</div>
				{/if}
				{#if hoveredTrade.closed_at}
					<div>
						<span class="text-slate-400">Closed:</span>
						<span class="font-semibold ml-1.5">{formatDate(hoveredTrade.closed_at)}</span>
					</div>
				{/if}
				{#if hoveredTrade.pnl !== null && hoveredTrade.pnl !== undefined}
					<div>
						<span class="text-slate-400">P&L %:</span>
						<span class="font-semibold ml-1.5 {hoveredTrade.pnl >= 0 ? 'text-emerald-400' : 'text-red-400'}">
							{((hoveredTrade.pnl / (hoveredTrade.entry_price * hoveredTrade.quantity)) * 100).toFixed(2)}%
						</span>
					</div>
				{/if}
				{#if hoveredTrade.strategy}
					<div class="col-span-2">
						<span class="text-slate-400">Strategy:</span>
						<span class="font-semibold ml-1.5">{hoveredTrade.strategy}</span>
					</div>
				{/if}
				{#if hoveredTrade.notes}
					<div class="col-span-2">
						<span class="text-slate-400">Notes:</span>
						<span class="ml-1.5">{hoveredTrade.notes.substring(0, 150)}{hoveredTrade.notes.length > 150 ? '...' : ''}</span>
					</div>
				{/if}
			</div>

			<!-- Journal Entries Summary -->
			{#if hoveredTrade.has_journal && hoveredTrade.journal_entries && hoveredTrade.journal_entries.length > 0}
				<div class="mt-3 pt-3 border-t border-slate-700">
					<div class="flex items-center gap-2 mb-2">
						<Icon icon="mdi:book-check" width="16" class="text-emerald-400" />
						<span class="font-semibold text-sm text-emerald-400">
							{hoveredTrade.journal_entries.length} Journal {hoveredTrade.journal_entries.length === 1 ? 'Entry' : 'Entries'}
						</span>
					</div>
					{#each hoveredTrade.journal_entries.slice(0, 2) as entry}
						<div class="mb-2 p-2 bg-slate-800 rounded">
							<div class="text-xs text-slate-500 mb-1">
								{new Date(entry.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
							</div>
							<div class="text-xs text-slate-300 line-clamp-2">
								{entry.content.substring(0, 100)}{entry.content.length > 100 ? '...' : ''}
							</div>
							{#if entry.emotional_state}
								<div class="flex gap-2 mt-1 text-[10px]">
									<span class="text-blue-400">Conf: {JSON.parse(entry.emotional_state).confidence}/10</span>
									<span class="text-amber-400">Stress: {JSON.parse(entry.emotional_state).stress}/10</span>
									<span class="text-emerald-400">Disc: {JSON.parse(entry.emotional_state).discipline}/10</span>
								</div>
							{/if}
						</div>
					{/each}
					{#if hoveredTrade.journal_entries.length > 2}
						<div class="text-xs text-slate-500 mt-1">
							+{hoveredTrade.journal_entries.length - 2} more {hoveredTrade.journal_entries.length - 2 === 1 ? 'entry' : 'entries'}
						</div>
					{/if}
				</div>
			{/if}

			<div class="text-xs text-slate-500 mt-2 pt-2 border-t border-slate-700">
				Tap anywhere to close
			</div>
		</div>
	</div>
{/if}

<!-- Slide-overs -->
<CSVImportSlideOver
	open={showImportModal}
	onClose={() => (showImportModal = false)}
	onImport={handleImportTrades}
/>
<TradeFormSlideOver isOpen={showAddModal} onClose={() => (showAddModal = false)} onSave={handleAddTrade} trade={selectedTrade} />
<TradeDetailSlideOver
	isOpen={showDetailSlideOver}
	onClose={() => (showDetailSlideOver = false)}
	trade={selectedTrade}
	onEdit={handleEditTrade}
	onDelete={handleDeleteTrade}
	onCreateJournal={() => {
		showDetailSlideOver = false;
		handleCreateJournal(selectedTrade);
	}}
/>
<JournalFormSlideOver
	open={showJournalModal}
	trade={selectedTrade}
	{rules}
	onClose={() => {
		showJournalModal = false;
		selectedTrade = null;
	}}
	onSubmit={handleJournalSubmit}
/>
