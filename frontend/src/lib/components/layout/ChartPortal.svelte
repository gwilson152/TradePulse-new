<script lang="ts">
	import Icon from '@iconify/svelte';
	import { chartPortal } from '$lib/stores/chartPortal';

	let chartContainerId = $state(`tradingview_portal_${Math.random().toString(36).substr(2, 9)}`);
	let isOpen = $state(false);
	let data = $state<any>(null);

	// Subscribe to store changes
	chartPortal.subscribe(state => {
		isOpen = state.isOpen;
		data = state.data;
	});

	// Format time for display
	function formatTime(timestamp: string): string {
		return new Date(timestamp).toLocaleTimeString('en-US', {
			hour: 'numeric',
			minute: '2-digit',
			second: '2-digit',
			hour12: true
		});
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

	// Get time range for the trade
	function getTimeRange() {
		if (!data) return null;

		const allTimestamps = [
			...data.entries.map(e => new Date(e.executed_at).getTime()),
			...data.exits.map(e => new Date(e.executed_at).getTime())
		];

		if (allTimestamps.length === 0) return null;

		const min = Math.min(...allTimestamps);
		const max = Math.max(...allTimestamps);

		return { start: min, end: max };
	}

	function handleClose() {
		chartPortal.close();
	}

	// Close on Escape key
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && isOpen) {
			handleClose();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen && data}
	<div
		class="chart-portal-overlay"
		onclick={handleClose}
		role="button"
		tabindex="0"
		onkeydown={(e) => e.key === 'Enter' && handleClose()}
	>
		<div
			class="chart-portal-container"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
			tabindex="-1"
			onkeydown={(e) => e.key === 'Enter' && e.stopPropagation()}
		>
			<!-- Header -->
			<div class="flex items-center justify-between p-4 border-b border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800">
				<div class="flex items-center gap-3">
					<div class="bg-blue-500/10 p-2 rounded-lg">
						<Icon icon="mdi:chart-line" class="text-xl text-blue-500" />
					</div>
					<div>
						<h3 class="font-semibold text-slate-900 dark:text-slate-100">Trade Chart - {data.symbol}</h3>
						<p class="text-sm text-slate-500 dark:text-slate-400">1 Minute Timeframe</p>
					</div>
				</div>
				<div class="flex items-center gap-4">
					<div class="flex items-center gap-4 text-xs">
						<div class="flex items-center gap-2">
							<div class="w-3 h-3 bg-emerald-500 rounded-full"></div>
							<span class="text-slate-600 dark:text-slate-400">
								{data.entries.length} {data.entries.length === 1 ? 'Entry' : 'Entries'}
							</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="w-3 h-3 bg-red-500 rounded-full"></div>
							<span class="text-slate-600 dark:text-slate-400">
								{data.exits.length} {data.exits.length === 1 ? 'Exit' : 'Exits'}
							</span>
						</div>
					</div>
					<button
						onclick={handleClose}
						class="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors"
						title="Close (Esc)"
					>
						<Icon icon="mdi:close" class="text-xl text-slate-600 dark:text-slate-400" />
					</button>
				</div>
			</div>

			<!-- TradingView Chart -->
			<div class="flex-1 bg-white dark:bg-slate-800">
				<iframe
					src="https://s.tradingview.com/widgetembed/?frameElementId={chartContainerId}&symbol={data.symbol}&interval=1&hidesidetoolbar=0&symboledit=0&saveimage=0&toolbarbg=f1f3f6&studies=%5B%5D&theme={document.documentElement.classList.contains('dark') ? 'dark' : 'light'}&style=1&timezone=Etc%2FUTC&studies_overrides=%7B%7D&overrides=%7B%22mainSeriesProperties.candleStyle.upColor%22%3A%22%2310b981%22%2C%22mainSeriesProperties.candleStyle.downColor%22%3A%22%23ef4444%22%2C%22mainSeriesProperties.candleStyle.borderUpColor%22%3A%22%2310b981%22%2C%22mainSeriesProperties.candleStyle.borderDownColor%22%3A%22%23ef4444%22%2C%22mainSeriesProperties.candleStyle.wickUpColor%22%3A%22%2310b981%22%2C%22mainSeriesProperties.candleStyle.wickDownColor%22%3A%22%23ef4444%22%7D&enabled_features=%5B%5D&disabled_features=%5B%5D&locale=en&utm_source=localhost&utm_medium=widget_new&utm_campaign=chart&utm_term={data.symbol}"
					style="width: 100%; height: calc(100vh - 80px); margin: 0 !important; padding: 0 !important;"
					frameborder="0"
					allowtransparency="true"
					scrolling="no"
					allowfullscreen
					title="TradingView Chart - {data.symbol}"
				></iframe>
			</div>
		</div>
	</div>
{/if}

<style>
	.chart-portal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 99999;
		background: rgba(0, 0, 0, 0.95);
		backdrop-filter: blur(8px);
		animation: fadeIn 0.3s ease-out;
	}

	.chart-portal-container {
		width: 100%;
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.tradingview-chart {
		transition: height 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	}
</style>
