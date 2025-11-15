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
			// TODO: Call API to bulk create trades
			console.log('Importing trades:', importedTrades);
			toast.success(`Imported ${importedTrades.length} trades successfully`);
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
					class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
					class:bg-primary-600={filterStatus === 'all'}
					class:text-white={filterStatus === 'all'}
					class:bg-surface-100={filterStatus !== 'all'}
					class:text-surface-700={filterStatus !== 'all'}
					onclick={() => (filterStatus = 'all')}
				>
					All
				</button>
				<button
					class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
					class:bg-primary-600={filterStatus === 'open'}
					class:text-white={filterStatus === 'open'}
					class:bg-surface-100={filterStatus !== 'open'}
					class:text-surface-700={filterStatus !== 'open'}
					onclick={() => (filterStatus = 'open')}
				>
					Open
				</button>
				<button
					class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
					class:bg-primary-600={filterStatus === 'closed'}
					class:text-white={filterStatus === 'closed'}
					class:bg-surface-100={filterStatus !== 'closed'}
					class:text-surface-700={filterStatus !== 'closed'}
					onclick={() => (filterStatus = 'closed')}
				>
					Closed
				</button>
				<button
					class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
					class:bg-primary-600={filterStatus === 'profitable'}
					class:text-white={filterStatus === 'profitable'}
					class:bg-surface-100={filterStatus !== 'profitable'}
					class:text-surface-700={filterStatus !== 'profitable'}
					onclick={() => (filterStatus = 'profitable')}
				>
					Profitable
				</button>
				<button
					class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
					class:bg-primary-600={filterStatus === 'loss'}
					class:text-white={filterStatus === 'loss'}
					class:bg-surface-100={filterStatus !== 'loss'}
					class:text-surface-700={filterStatus !== 'loss'}
					onclick={() => (filterStatus = 'loss')}
				>
					Loss
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
		<Card>
			<div class="text-center py-12">
				<Icon icon="mdi:chart-line" width="64" class="mx-auto mb-4 text-surface-400" />
				<h3 class="text-xl font-semibold mb-2">No trades yet</h3>
				<p class="text-surface-600 dark:text-surface-400 mb-6">
					Import your trade history from CSV or manually add trades
				</p>
				<div class="flex gap-4 justify-center">
					<Button color="primary" size="lg" onclick={() => (showImportModal = true)}>
						<Icon icon="mdi:file-upload" width="24" class="mr-2" />
						Import from CSV
					</Button>
					<Button color="soft" variant="soft" onclick={() => (showAddModal = true)}>
						<Icon icon="mdi:plus" width="20" class="mr-2" />
						Add Manually
					</Button>
				</div>
			</div>
		</Card>
	{:else}
		<Card padding="none">
			<Table hover={true}>
				<thead>
					<tr>
						<th>Date</th>
						<th>Symbol</th>
						<th>Type</th>
						<th>Qty</th>
						<th>Entry</th>
						<th>Exit</th>
						<th>P&L</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					{#each filteredTrades() as trade}
						<tr onclick={() => handleViewTrade(trade)} class="cursor-pointer">
							<td class="text-sm">{formatDate(trade.opened_at)}</td>
							<td class="font-medium">{trade.symbol}</td>
							<td>
								<span
									class="px-2 py-1 rounded text-xs font-medium"
									class:bg-success-100={trade.trade_type === 'LONG'}
									class:text-success-700={trade.trade_type === 'LONG'}
									class:bg-error-100={trade.trade_type === 'SHORT'}
									class:text-error-700={trade.trade_type === 'SHORT'}
								>
									{trade.trade_type}
								</span>
							</td>
							<td>{trade.quantity}</td>
							<td class="text-sm">{formatCurrency(trade.entry_price)}</td>
							<td class="text-sm">
								{trade.exit_price ? formatCurrency(trade.exit_price) : '-'}
							</td>
							<td>
								{#if trade.pnl !== null}
									<PnLBadge value={trade.pnl} showSign={true} size="sm" />
								{:else}
									<span class="text-surface-500 text-sm">Open</span>
								{/if}
							</td>
							<td>
								<button
									class="p-2 hover:bg-surface-100 dark:hover:bg-surface-700 rounded"
									onclick={(e) => {
										e.stopPropagation();
										handleViewTrade(trade);
									}}
								>
									<Icon icon="mdi:chevron-right" width="20" />
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</Table>
		</Card>
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
