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
	import { apiClient } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import type { Trade } from '$lib/types';

	let trades = $state<Trade[]>([]);
	let isLoading = $state(true);
	let showAddModal = $state(false);
	let showImportModal = $state(false);
	let showDetailSlideOver = $state(false);
	let selectedTrade = $state<Trade | null>(null);
	let searchQuery = $state('');
	let filterStatus = $state('all'); // all, open, closed, profitable, loss

	onMount(async () => {
		await loadTrades();
	});

	async function loadTrades() {
		try {
			isLoading = true;
			trades = await apiClient.getTrades();
		} catch (err) {
			console.error('Failed to load trades:', err);
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

	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric'
		});
	}

	function formatCurrency(value: number) {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(value);
	}

	let filteredTrades = $derived(() => {
		let filtered = trades;

		// Apply status filter
		if (filterStatus === 'open') {
			filtered = filtered.filter((t) => !t.exit_price);
		} else if (filterStatus === 'closed') {
			filtered = filtered.filter((t) => t.exit_price);
		} else if (filterStatus === 'profitable') {
			filtered = filtered.filter((t) => t.pnl && t.pnl > 0);
		} else if (filterStatus === 'loss') {
			filtered = filtered.filter((t) => t.pnl && t.pnl < 0);
		}

		// Apply search
		if (searchQuery) {
			filtered = filtered.filter((t) =>
				t.symbol.toLowerCase().includes(searchQuery.toLowerCase())
			);
		}

		return filtered;
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
			<p class="text-surface-600 dark:text-surface-400">Import from CSV or manually add trades</p>
		</div>
		<div class="flex gap-3">
			<Button color="soft" variant="soft" onclick={() => (showAddModal = true)}>
				<Icon icon="mdi:plus" width="20" class="mr-2" />
				Add Trade
			</Button>
			<Button color="primary" onclick={() => (showImportModal = true)}>
				<Icon icon="mdi:file-upload" width="20" class="mr-2" />
				Import CSV
			</Button>
		</div>
	</div>

	<!-- Filters -->
	<Card>
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

			<!-- Filter Chips -->
			<div class="flex gap-2 flex-wrap">
				<button
					class="px-4 py-2 rounded-xl text-sm font-medium transition-all hover:scale-105"
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
					<Icon icon="mdi:view-grid" width="16" class="inline mr-1.5" />
					All ({trades.length})
				</button>
				<button
					class="px-4 py-2 rounded-xl text-sm font-medium transition-all hover:scale-105"
					class:bg-blue-500={filterStatus === 'open'}
					class:text-white={filterStatus === 'open'}
					class:shadow-md={filterStatus === 'open'}
					class:bg-slate-100={filterStatus !== 'open'}
					class:dark:bg-slate-800={filterStatus !== 'open'}
					class:text-slate-700={filterStatus !== 'open'}
					class:dark:text-slate-300={filterStatus !== 'open'}
					onclick={() => (filterStatus = 'open')}
				>
					<Icon icon="mdi:clock-outline" width="16" class="inline mr-1.5" />
					Open
				</button>
				<button
					class="px-4 py-2 rounded-xl text-sm font-medium transition-all hover:scale-105"
					class:bg-slate-600={filterStatus === 'closed'}
					class:text-white={filterStatus === 'closed'}
					class:shadow-md={filterStatus === 'closed'}
					class:bg-slate-100={filterStatus !== 'closed'}
					class:dark:bg-slate-800={filterStatus !== 'closed'}
					class:text-slate-700={filterStatus !== 'closed'}
					class:dark:text-slate-300={filterStatus !== 'closed'}
					onclick={() => (filterStatus = 'closed')}
				>
					<Icon icon="mdi:check-circle" width="16" class="inline mr-1.5" />
					Closed
				</button>
				<button
					class="px-4 py-2 rounded-xl text-sm font-medium transition-all hover:scale-105"
					class:bg-emerald-500={filterStatus === 'profitable'}
					class:text-white={filterStatus === 'profitable'}
					class:shadow-md={filterStatus === 'profitable'}
					class:bg-slate-100={filterStatus !== 'profitable'}
					class:dark:bg-slate-800={filterStatus !== 'profitable'}
					class:text-slate-700={filterStatus !== 'profitable'}
					class:dark:text-slate-300={filterStatus !== 'profitable'}
					onclick={() => (filterStatus = 'profitable')}
				>
					<Icon icon="mdi:trending-up" width="16" class="inline mr-1.5" />
					Winners
				</button>
				<button
					class="px-4 py-2 rounded-xl text-sm font-medium transition-all hover:scale-105"
					class:bg-red-500={filterStatus === 'loss'}
					class:text-white={filterStatus === 'loss'}
					class:shadow-md={filterStatus === 'loss'}
					class:bg-slate-100={filterStatus !== 'loss'}
					class:dark:bg-slate-800={filterStatus !== 'loss'}
					class:text-slate-700={filterStatus !== 'loss'}
					class:dark:text-slate-300={filterStatus !== 'loss'}
					onclick={() => (filterStatus === 'loss')}
				>
					<Icon icon="mdi:trending-down" width="16" class="inline mr-1.5" />
					Losers
				</button>
			</div>
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
	{:else if filteredTrades().length === 0}
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
						Import your trade history from CSV or manually add your first trade to begin tracking your performance
					{/if}
				</p>
				{#if !searchQuery && filterStatus === 'all'}
					<div class="flex flex-col sm:flex-row gap-3 justify-center">
						<Button color="primary" size="lg" onclick={() => (showImportModal = true)}>
							<Icon icon="mdi:file-upload" width="24" class="mr-2" />
							Import from CSV
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
		<div class="space-y-3">
			{#each filteredTrades() as trade}
				<div
					role="button"
					tabindex="0"
					onclick={() => handleViewTrade(trade)}
					onkeydown={(e) => e.key === 'Enter' && handleViewTrade(trade)}
					class="group bg-white/80 dark:bg-slate-800/80 backdrop-blur-xl rounded-2xl border border-slate-200/50 dark:border-slate-700/50 p-5 hover:shadow-xl hover:border-blue-300 dark:hover:border-blue-600 transition-all duration-200 cursor-pointer hover:scale-[1.01]"
				>
					<div class="flex items-center justify-between gap-4">
						<!-- Left: Symbol & Time -->
						<div class="flex items-center gap-4 flex-1">
							<div class="flex-shrink-0">
								<div class="w-12 h-12 rounded-xl bg-gradient-to-br {trade.trade_type === 'LONG' ? 'from-emerald-400 to-emerald-600' : 'from-red-400 to-red-600'} flex items-center justify-center shadow-md">
									<Icon icon={trade.trade_type === 'LONG' ? 'mdi:arrow-up-bold' : 'mdi:arrow-down-bold'} width="24" class="text-white" />
								</div>
							</div>
							<div>
								<div class="flex items-center gap-2 mb-1">
									<h3 class="text-lg font-bold text-slate-800 dark:text-slate-100">{trade.symbol}</h3>
									<span class="px-2.5 py-0.5 rounded-full text-xs font-semibold {trade.trade_type === 'LONG' ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400' : 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'}">
										{trade.trade_type}
									</span>
									{#if !trade.exit_price}
										<span class="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 flex items-center gap-1">
											<Icon icon="mdi:clock-outline" width="12" />
											Open
										</span>
									{/if}
								</div>
								<p class="text-sm text-slate-600 dark:text-slate-400 flex items-center gap-2">
									<Icon icon="mdi:calendar" width="14" />
									{formatDate(trade.opened_at)}
								</p>
							</div>
						</div>

						<!-- Middle: Trade Details -->
						<div class="hidden md:flex items-center gap-6">
							<div class="text-center">
								<p class="text-xs text-slate-500 dark:text-slate-400 mb-1">Quantity</p>
								<p class="text-sm font-semibold text-slate-800 dark:text-slate-200">{trade.quantity}</p>
							</div>
							<div class="text-center">
								<p class="text-xs text-slate-500 dark:text-slate-400 mb-1">Entry</p>
								<p class="text-sm font-semibold text-slate-800 dark:text-slate-200">{formatCurrency(trade.entry_price)}</p>
							</div>
							<div class="text-center">
								<p class="text-xs text-slate-500 dark:text-slate-400 mb-1">Exit</p>
								<p class="text-sm font-semibold text-slate-800 dark:text-slate-200">
									{trade.exit_price ? formatCurrency(trade.exit_price) : '-'}
								</p>
							</div>
						</div>

						<!-- Right: P&L & Action -->
						<div class="flex items-center gap-4">
							<div class="text-right">
								{#if trade.pnl !== null}
									<div class="px-4 py-2 rounded-xl font-bold text-lg {trade.pnl >= 0 ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400' : 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'}">
										{trade.pnl >= 0 ? '+' : ''}{formatCurrency(trade.pnl)}
									</div>
								{:else}
									<div class="px-4 py-2 rounded-xl bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400 font-semibold">
										Pending
									</div>
								{/if}
							</div>
							<button
								onclick={(e) => {
									e.stopPropagation();
									handleViewTrade(trade);
								}}
								class="p-3 rounded-xl bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400 hover:bg-blue-100 dark:hover:bg-blue-900/30 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
							>
								<Icon icon="mdi:chevron-right" width="24" />
							</button>
						</div>
					</div>

					<!-- Mobile: Extra Details -->
					<div class="md:hidden mt-4 pt-4 border-t border-slate-200 dark:border-slate-700 grid grid-cols-3 gap-4 text-center">
						<div>
							<p class="text-xs text-slate-500 dark:text-slate-400 mb-1">Quantity</p>
							<p class="text-sm font-semibold text-slate-800 dark:text-slate-200">{trade.quantity}</p>
						</div>
						<div>
							<p class="text-xs text-slate-500 dark:text-slate-400 mb-1">Entry</p>
							<p class="text-sm font-semibold text-slate-800 dark:text-slate-200">{formatCurrency(trade.entry_price)}</p>
						</div>
						<div>
							<p class="text-xs text-slate-500 dark:text-slate-400 mb-1">Exit</p>
							<p class="text-sm font-semibold text-slate-800 dark:text-slate-200">
								{trade.exit_price ? formatCurrency(trade.exit_price) : '-'}
							</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

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
/>
