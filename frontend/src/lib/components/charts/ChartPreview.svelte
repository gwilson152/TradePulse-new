<script lang="ts">
	import Icon from '@iconify/svelte';
	import { chartPortal } from '$lib/stores/chartPortal';
	import type { Entry, Exit } from '$lib/types';

	interface Props {
		tradeId?: string;
		symbol: string;
		entries: Entry[];
		exits: Exit[];
		openedAt?: string;
		closedAt?: string | null;
	}

	let { tradeId, symbol, entries, exits, openedAt, closedAt }: Props = $props();

	let chartContainerId = $state(`tradingview_${Math.random().toString(36).substr(2, 9)}`);
	let isTimelineExpanded = $state(true);

	// Format time for display
	function formatTime(timestamp: string | null | undefined): string {
		if (!timestamp) return 'Invalid Date';
		try {
			const date = new Date(timestamp);
			if (isNaN(date.getTime())) return 'Invalid Date';
			return date.toLocaleTimeString('en-US', {
				hour: 'numeric',
				minute: '2-digit',
				second: '2-digit',
				hour12: true
			});
		} catch {
			return 'Invalid Date';
		}
	}

	// Format price
	function formatPrice(price: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 2,
			maximumFractionDigits: 2
		}).format(price);
	}

	// Calculate P&L for exit
	function calculateExitPnL(exit: Exit): number {
		const avgEntryPrice = entries.reduce((sum, e) => sum + e.price, 0) / entries.length;
		const direction = entries[0]?.side === 'buy' ? 1 : -1;
		return direction * (exit.price - avgEntryPrice) * exit.quantity;
	}

	// Get time range for the trade
	function getTimeRange() {
		const allTimestamps = [
			...entries.map(e => new Date(e.executed_at).getTime()),
			...exits.map(e => new Date(e.executed_at).getTime())
		];

		if (allTimestamps.length === 0) return null;

		const min = Math.min(...allTimestamps);
		const max = Math.max(...allTimestamps);

		return {
			start: min,
			end: max
		};
	}

	function openFullscreen() {
		chartPortal.open({
			tradeId,
			symbol,
			entries,
			exits,
			openedAt,
			closedAt
		});
	}

</script>

<div class="chart-preview-container">
	<div class="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 overflow-hidden">
		<!-- Header -->
		<div class="flex items-center justify-between p-4 border-b border-slate-200 dark:border-slate-700">
			<div class="flex items-center gap-3">
				<div class="bg-blue-500/10 p-2 rounded-lg">
					<Icon icon="mdi:chart-line" class="text-xl text-blue-500" />
				</div>
				<div>
					<h3 class="font-semibold text-slate-900 dark:text-slate-100">Trade Chart</h3>
					<p class="text-sm text-slate-500 dark:text-slate-400">{symbol} - 1 Minute</p>
				</div>
			</div>
			<div class="flex items-center gap-4">
				<div class="flex items-center gap-4 text-xs">
					<div class="flex items-center gap-2">
						<div class="w-3 h-3 bg-emerald-500 rounded-full"></div>
						<span class="text-slate-600 dark:text-slate-400">{entries.length} {entries.length === 1 ? 'Entry' : 'Entries'}</span>
					</div>
					<div class="flex items-center gap-2">
						<div class="w-3 h-3 bg-red-500 rounded-full"></div>
						<span class="text-slate-600 dark:text-slate-400">{exits.length} {exits.length === 1 ? 'Exit' : 'Exits'}</span>
					</div>
				</div>
				<button
					onclick={openFullscreen}
					class="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors"
					title="Maximize Chart"
				>
					<Icon icon="mdi:fullscreen" class="text-xl text-slate-600 dark:text-slate-400" />
				</button>
			</div>
		</div>

		<!-- TradingView Chart -->
		<div class="tradingview-chart" style="height: 500px;">
			<iframe
				src="https://s.tradingview.com/widgetembed/?frameElementId={chartContainerId}&symbol={symbol}&interval=1&hidesidetoolbar=0&symboledit=0&saveimage=0&toolbarbg=f1f3f6&studies=%5B%5D&theme={document.documentElement.classList.contains('dark') ? 'dark' : 'light'}&style=1&timezone=Etc%2FUTC&studies_overrides=%7B%7D&overrides=%7B%22mainSeriesProperties.candleStyle.upColor%22%3A%22%2310b981%22%2C%22mainSeriesProperties.candleStyle.downColor%22%3A%22%23ef4444%22%2C%22mainSeriesProperties.candleStyle.borderUpColor%22%3A%22%2310b981%22%2C%22mainSeriesProperties.candleStyle.borderDownColor%22%3A%22%23ef4444%22%2C%22mainSeriesProperties.candleStyle.wickUpColor%22%3A%22%2310b981%22%2C%22mainSeriesProperties.candleStyle.wickDownColor%22%3A%22%23ef4444%22%7D&enabled_features=%5B%5D&disabled_features=%5B%5D&locale=en&utm_source=localhost&utm_medium=widget_new&utm_campaign=chart&utm_term={symbol}"
				style="width: 100%; height: 100%; margin: 0 !important; padding: 0 !important;"
				frameborder="0"
				allowtransparency="true"
				scrolling="no"
				allowfullscreen
				title="TradingView Chart"
			></iframe>
		</div>

		<!-- Trade Execution Timeline -->
		<div class="bg-slate-50 dark:bg-slate-900/50 border-t border-slate-200 dark:border-slate-700">
			<button
				onclick={() => isTimelineExpanded = !isTimelineExpanded}
				class="w-full flex items-center justify-between p-4 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
			>
				<h4 class="text-sm font-semibold text-slate-700 dark:text-slate-300">Execution Timeline</h4>
				<Icon icon={isTimelineExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'} class="text-lg text-slate-600 dark:text-slate-400" />
			</button>
			{#if isTimelineExpanded}
			<div class="px-4 pb-4 space-y-2">
				{#each entries as entry, idx}
					<div class="flex items-center gap-3 text-sm">
						<div class="flex-shrink-0 w-20 text-xs text-slate-500 dark:text-slate-400">
							{formatTime(entry.executed_at)}
						</div>
						<div class="flex-shrink-0">
							<div class="w-3 h-3 bg-emerald-500 rounded-full"></div>
						</div>
						<div class="flex-1 flex items-center gap-2">
							<span class="font-medium text-emerald-600 dark:text-emerald-400">Entry {idx + 1}</span>
							<span class="text-slate-600 dark:text-slate-400">
								{entry.quantity} @ {formatPrice(entry.price)}
							</span>
							{#if entry.notes}
								<span class="text-xs text-slate-500 dark:text-slate-400 italic">- {entry.notes}</span>
							{/if}
						</div>
					</div>
				{/each}

				{#each exits as exit, idx}
					{@const pnl = calculateExitPnL(exit)}
					<div class="flex items-center gap-3 text-sm">
						<div class="flex-shrink-0 w-20 text-xs text-slate-500 dark:text-slate-400">
							{formatTime(exit.executed_at)}
						</div>
						<div class="flex-shrink-0">
							<div class="w-3 h-3 bg-red-500 rounded-full"></div>
						</div>
						<div class="flex-1 flex items-center gap-2">
							<span class="font-medium text-red-600 dark:text-red-400">Exit {idx + 1}</span>
							<span class="text-slate-600 dark:text-slate-400">
								{exit.quantity} @ {formatPrice(exit.price)}
							</span>
							<span class="font-semibold {pnl >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}">
								{pnl >= 0 ? '+' : ''}{formatPrice(pnl)}
							</span>
							{#if exit.notes}
								<span class="text-xs text-slate-500 dark:text-slate-400 italic">- {exit.notes}</span>
							{/if}
						</div>
					</div>
				{/each}
			</div>
			{/if}
		</div>

		<!-- Summary stats -->
		<div class="p-4 bg-white dark:bg-slate-800 border-t border-slate-200 dark:border-slate-700">
			<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
				<div class="bg-emerald-50 dark:bg-emerald-900/10 p-3 rounded-lg">
					<div class="text-xs text-emerald-600 dark:text-emerald-400 font-medium mb-1">Avg Entry</div>
					<div class="text-lg font-bold text-emerald-700 dark:text-emerald-300">
						{formatPrice(entries.reduce((sum, e) => sum + e.price, 0) / entries.length)}
					</div>
				</div>
				<div class="bg-red-50 dark:bg-red-900/10 p-3 rounded-lg">
					<div class="text-xs text-red-600 dark:text-red-400 font-medium mb-1">Avg Exit</div>
					<div class="text-lg font-bold text-red-700 dark:text-red-300">
						{exits.length > 0 ? formatPrice(exits.reduce((sum, e) => sum + e.price, 0) / exits.length) : 'N/A'}
					</div>
				</div>
				<div class="bg-blue-50 dark:bg-blue-900/10 p-3 rounded-lg">
					<div class="text-xs text-blue-600 dark:text-blue-400 font-medium mb-1">Total Quantity</div>
					<div class="text-lg font-bold text-blue-700 dark:text-blue-300">
						{entries.reduce((sum, e) => sum + e.quantity, 0)}
					</div>
				</div>
				<div class="bg-purple-50 dark:bg-purple-900/10 p-3 rounded-lg">
					<div class="text-xs text-purple-600 dark:text-purple-400 font-medium mb-1">Executions</div>
					<div class="text-lg font-bold text-purple-700 dark:text-purple-300">
						{entries.length + exits.length}
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.chart-preview-container {
		width: 100%;
	}

	.tradingview-chart {
		transition: height 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	}
</style>
