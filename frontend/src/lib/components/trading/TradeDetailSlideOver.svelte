<script lang="ts">
	import Icon from '@iconify/svelte';
	import Button from '../ui/Button.svelte';
	import Badge from '../ui/Badge.svelte';
	import PnLBadge from './PnLBadge.svelte';
	import PositionTimeline from './PositionTimeline.svelte';
	import PositionSizeBar from './PositionSizeBar.svelte';
	import AddToPositionModal from './AddToPositionModal.svelte';
	import type { Trade, Entry, Exit } from '$lib/types';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
		trade: Trade;
		onEdit?: () => void;
		onDelete?: () => void;
		onAddEntry?: (data: Partial<Entry>) => Promise<void>;
		onAddExit?: (data: Partial<Exit>) => Promise<void>;
	}

	let { isOpen, onClose, trade, onEdit, onDelete, onAddEntry, onAddExit }: Props = $props();

	let showAddModal = $state(false);
	let addModalAction = $state<'entry' | 'exit'>('entry');
	let currentTab = $state<'overview' | 'timeline'>('overview');

	function formatDate(date: string) {
		return new Date(date).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatCurrency(value: number) {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(value);
	}

	function calculateReturn() {
		if (trade.exit_price) {
			const diff = trade.trade_type === 'LONG'
				? trade.exit_price - trade.entry_price
				: trade.entry_price - trade.exit_price;
			return ((diff / trade.entry_price) * 100).toFixed(2);
		}
		return null;
	}

	function openAddModal(action: 'entry' | 'exit') {
		addModalAction = action;
		showAddModal = true;
	}

	async function handleAddToPosition(data: Partial<Entry> | Partial<Exit>) {
		if (addModalAction === 'entry' && onAddEntry) {
			await onAddEntry(data as Partial<Entry>);
		} else if (addModalAction === 'exit' && onAddExit) {
			await onAddExit(data as Partial<Exit>);
		}
	}

	const isOpen_position = $derived(
		trade.current_position_size > 0 || (!trade.closed_at && trade.entries && trade.entries.length > 0)
	);
</script>

{#if isOpen && trade}
	<div class="fixed inset-0 z-50 overflow-hidden">
		<div class="absolute inset-0 bg-black/50" onclick={onClose}></div>

		<div class="absolute inset-y-0 right-0 max-w-md w-full">
			<div
				class="h-full bg-white dark:bg-surface-800 shadow-xl flex flex-col"
			>
				<!-- Header -->
				<div class="p-6 border-b border-surface-200 dark:border-surface-700">
					<div class="flex items-start justify-between mb-4">
						<div>
							<div class="flex items-center gap-3 mb-2">
								<h2 class="text-2xl font-bold">{trade.symbol}</h2>
								{#if isOpen_position}
									<Badge color="success" variant="soft" size="sm" icon="mdi:circle">
										Open
									</Badge>
								{:else}
									<Badge color="neutral" variant="soft" size="sm">
										Closed
									</Badge>
								{/if}
							</div>
							<p class="text-surface-600 dark:text-surface-400">
								{trade.trade_type} Position
							</p>
						</div>
						<button
							onclick={onClose}
							class="p-2 hover:bg-surface-100 dark:hover:bg-surface-700 rounded"
						>
							<Icon icon="mdi:close" width="24" />
						</button>
					</div>

					<!-- Tabs -->
					<div class="flex gap-2 border-b border-surface-200 dark:border-surface-700 -mb-px">
						<button
							onclick={() => (currentTab = 'overview')}
							class="px-4 py-2 border-b-2 transition-colors {currentTab === 'overview'
								? 'border-primary-500 text-primary-600 dark:text-primary-400'
								: 'border-transparent text-surface-600 dark:text-surface-400 hover:text-surface-800 dark:hover:text-surface-200'}"
						>
							Overview
						</button>
						<button
							onclick={() => (currentTab = 'timeline')}
							class="px-4 py-2 border-b-2 transition-colors {currentTab === 'timeline'
								? 'border-primary-500 text-primary-600 dark:text-primary-400'
								: 'border-transparent text-surface-600 dark:text-surface-400 hover:text-surface-800 dark:hover:text-surface-200'}"
						>
							Timeline
						</button>
					</div>
				</div>

				<!-- Content -->
				<div class="flex-1 overflow-y-auto p-6 space-y-6">
					{#if currentTab === 'overview'}
					<!-- Position Size -->
					{#if trade.entries && trade.entries.length > 0}
						<PositionSizeBar
							entries={trade.entries}
							exits={trade.exits || []}
							symbol={trade.symbol}
						/>
					{/if}

					<!-- Quick Actions (for open positions) -->
					{#if isOpen_position && (onAddEntry || onAddExit)}
						<div class="grid grid-cols-2 gap-3">
							{#if onAddEntry}
								<Button
									color="primary"
									variant="soft"
									size="sm"
									onclick={() => openAddModal('entry')}
								>
									<Icon icon="mdi:plus" class="mr-1" />
									Add Entry
								</Button>
							{/if}
							{#if onAddExit && trade.current_position_size > 0}
								<Button
									color="secondary"
									variant="soft"
									size="sm"
									onclick={() => openAddModal('exit')}
								>
									<Icon icon="mdi:minus" class="mr-1" />
									Add Exit
								</Button>
							{/if}
						</div>
					{/if}

					<!-- Trade Information -->
					<div>
						<h3 class="text-lg font-semibold mb-3">Position Summary</h3>
						<div class="space-y-2">
							<div class="flex justify-between">
								<span class="text-surface-600 dark:text-surface-400">Average Entry:</span>
								<span class="font-medium">
									{formatCurrency(trade.average_entry_price || trade.entry_price)}
								</span>
							</div>
							<div class="flex justify-between">
								<span class="text-surface-600 dark:text-surface-400">Current Size:</span>
								<span class="font-medium">{trade.current_position_size || trade.quantity} shares</span>
							</div>
							<div class="flex justify-between">
								<span class="text-surface-600 dark:text-surface-400">Total Entries:</span>
								<span class="font-medium">
									{trade.entries?.length || 1} × {trade.entries?.reduce((sum, e) => sum + e.quantity, 0) || trade.quantity} shares
								</span>
							</div>
							{#if trade.exits && trade.exits.length > 0}
								<div class="flex justify-between">
									<span class="text-surface-600 dark:text-surface-400">Total Exits:</span>
									<span class="font-medium">
										{trade.exits.length} × {trade.exits.reduce((sum, e) => sum + e.quantity, 0)} shares
									</span>
								</div>
							{/if}
						</div>
					</div>

					<!-- Entry -->
					<div>
						<h3 class="text-lg font-semibold mb-3">Entry</h3>
						<div class="space-y-2">
							<div class="flex justify-between">
								<span class="text-surface-600 dark:text-surface-400">Price:</span>
								<span class="font-medium">{formatCurrency(trade.entry_price)}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-surface-600 dark:text-surface-400">Date:</span>
								<span class="font-medium">{formatDate(trade.entry_date)}</span>
							</div>
						</div>
					</div>

					<!-- Exit -->
					{#if trade.exit_price}
						<div>
							<h3 class="text-lg font-semibold mb-3">Exit</h3>
							<div class="space-y-2">
								<div class="flex justify-between">
									<span class="text-surface-600 dark:text-surface-400">Price:</span>
									<span class="font-medium">{formatCurrency(trade.exit_price)}</span>
								</div>
								<div class="flex justify-between">
									<span class="text-surface-600 dark:text-surface-400">Date:</span>
									<span class="font-medium">{formatDate(trade.exit_date)}</span>
								</div>
							</div>
						</div>

						<!-- Results -->
						<div>
							<h3 class="text-lg font-semibold mb-3">Results</h3>
							<div class="space-y-2">
								<div class="flex justify-between">
									<span class="text-surface-600 dark:text-surface-400">P&L:</span>
									<PnLBadge value={trade.pnl} showSign={true} />
								</div>
								<div class="flex justify-between">
									<span class="text-surface-600 dark:text-surface-400">Fees:</span>
									<span class="font-medium">{formatCurrency(trade.fees)}</span>
								</div>
								<div class="flex justify-between">
									<span class="text-surface-600 dark:text-surface-400">Net:</span>
									<span class="font-medium">{formatCurrency(trade.pnl - trade.fees)}</span>
								</div>
								<div class="flex justify-between">
									<span class="text-surface-600 dark:text-surface-400">Return:</span>
									<span
										class="font-medium"
										class:text-profit-600={parseFloat(calculateReturn() || '0') > 0}
										class:text-loss-600={parseFloat(calculateReturn() || '0') < 0}
									>
										{calculateReturn()}%
									</span>
								</div>
							</div>
						</div>
					{:else}
						<div class="p-4 bg-primary-100 dark:bg-primary-900/30 rounded-lg">
							<p class="text-sm text-primary-700 dark:text-primary-300">
								This trade is still open
							</p>
						</div>
					{/if}

					<!-- Tags -->
					{#if trade.tags && trade.tags.length > 0}
						<div>
							<h3 class="text-lg font-semibold mb-3">Tags</h3>
							<div class="flex flex-wrap gap-2">
								{#each trade.tags as tag}
									<Badge color="primary" variant="soft" size="sm">
										{tag}
									</Badge>
								{/each}
							</div>
						</div>
					{/if}

					<!-- Notes -->
					{#if trade.notes}
						<div>
							<h3 class="text-lg font-semibold mb-3">Notes</h3>
							<p class="text-surface-600 dark:text-surface-400 whitespace-pre-wrap">
								{trade.notes}
							</p>
						</div>
					{/if}

					{:else if currentTab === 'timeline'}
						<!-- Timeline Tab -->
						{#if trade.entries && trade.entries.length > 0}
							<PositionTimeline
								entries={trade.entries}
								exits={trade.exits || []}
								tradeType={trade.trade_type}
							/>
						{:else}
							<div class="text-center py-12 text-surface-600 dark:text-surface-400">
								<Icon icon="mdi:timeline" class="text-4xl mb-3 mx-auto" />
								<p>No timeline data available for this trade</p>
							</div>
						{/if}
					{/if}
				</div>

				<!-- Actions -->
				<div class="p-6 border-t border-surface-200 dark:border-surface-700 flex gap-4">
					{#if onEdit}
						<Button color="primary" onclick={onEdit}>
							<Icon icon="mdi:pencil" width="20" class="mr-2" />
							Edit
						</Button>
					{/if}
					{#if onDelete}
						<Button color="error" onclick={onDelete}>
							<Icon icon="mdi:delete" width="20" class="mr-2" />
							Delete
						</Button>
					{/if}
				</div>
			</div>
		</div>
	</div>

	<!-- Add to Position Modal -->
	<AddToPositionModal
		open={showAddModal}
		{trade}
		action={addModalAction}
		onClose={() => (showAddModal = false)}
		onSubmit={handleAddToPosition}
	/>
{/if}
