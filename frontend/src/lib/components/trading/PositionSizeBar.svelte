<script lang="ts">
	import type { Entry, Exit } from '$lib/types';

	interface Props {
		entries: Entry[];
		exits: Exit[];
		symbol: string;
	}

	let { entries = [], exits = [], symbol }: Props = $props();

	// Calculate position metrics
	const metrics = $derived(() => {
		const totalEntered = entries.reduce((sum, entry) => sum + entry.quantity, 0);
		const totalExited = exits.reduce((sum, exit) => sum + exit.quantity, 0);
		const currentSize = totalEntered - totalExited;
		const percentRemaining = totalEntered > 0 ? (currentSize / totalEntered) * 100 : 0;

		return {
			totalEntered,
			totalExited,
			currentSize,
			percentRemaining,
			isClosed: currentSize === 0 && totalEntered > 0
		};
	});

	function formatShares(qty: number): string {
		return qty.toLocaleString();
	}
</script>

<div class="position-size-bar">
	<div class="flex justify-between items-center mb-2">
		<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300">Position Size</h4>
		<div class="text-right">
			<span class="text-lg font-bold text-surface-900 dark:text-surface-100">
				{formatShares(metrics().currentSize)}
			</span>
			<span class="text-sm text-surface-600 dark:text-surface-400 ml-1">shares</span>
		</div>
	</div>

	<!-- Visual bar -->
	<div class="relative h-8 bg-surface-200 dark:bg-surface-700 rounded-lg overflow-hidden mb-2">
		<!-- Remaining position (green/blue) -->
		<div
			class="absolute left-0 top-0 bottom-0 {metrics().isClosed
				? 'bg-surface-400 dark:bg-surface-600'
				: 'bg-blue-500 dark:bg-blue-600'} transition-all duration-300"
			style="width: {metrics().percentRemaining}%"
		></div>

		<!-- Exited portion (gray) -->
		{#if metrics().totalExited > 0}
			<div
				class="absolute right-0 top-0 bottom-0 bg-surface-300 dark:bg-surface-600 transition-all duration-300"
				style="width: {100 - metrics().percentRemaining}%"
			></div>
		{/if}

		<!-- Center text overlay -->
		<div class="absolute inset-0 flex items-center justify-center">
			<span class="text-xs font-semibold text-white drop-shadow-md">
				{#if metrics().isClosed}
					Position Closed
				{:else if metrics().percentRemaining === 100}
					Full Position
				{:else}
					{metrics().percentRemaining.toFixed(0)}% Remaining
				{/if}
			</span>
		</div>
	</div>

	<!-- Details -->
	<div class="grid grid-cols-2 gap-4 text-sm">
		<div>
			<p class="text-surface-600 dark:text-surface-400">Total Entered</p>
			<p class="font-semibold text-surface-900 dark:text-surface-100">
				{formatShares(metrics().totalEntered)}
			</p>
		</div>
		<div>
			<p class="text-surface-600 dark:text-surface-400">Total Exited</p>
			<p class="font-semibold text-surface-900 dark:text-surface-100">
				{formatShares(metrics().totalExited)}
			</p>
		</div>
	</div>

	{#if !metrics().isClosed && metrics().currentSize > 0}
		<div class="mt-3 p-2 bg-blue-50 dark:bg-blue-900/20 rounded border border-blue-200 dark:border-blue-800">
			<p class="text-xs text-blue-700 dark:text-blue-300">
				<strong>Active Position:</strong> You currently hold {formatShares(metrics().currentSize)} shares
				of {symbol}
			</p>
		</div>
	{/if}
</div>

<style>
	.position-size-bar {
		padding: 1rem;
		background-color: rgb(var(--color-surface-50));
		border-radius: 0.5rem;
	}

	:global(.dark) .position-size-bar {
		background-color: rgb(var(--color-surface-800));
	}
</style>
