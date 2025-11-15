<script lang="ts">
	import type { Entry, Exit } from '$lib/types';
	import Icon from '@iconify/svelte';

	interface Props {
		entries: Entry[];
		exits: Exit[];
		tradeType: 'LONG' | 'SHORT';
	}

	let { entries = [], exits = [], tradeType }: Props = $props();

	// Combine entries and exits into timeline events
	const timelineEvents = $derived(() => {
		const events: Array<{
			type: 'entry' | 'exit';
			data: Entry | Exit;
			timestamp: string;
		}> = [];

		entries.forEach((entry) => {
			events.push({
				type: 'entry',
				data: entry,
				timestamp: entry.timestamp
			});
		});

		exits.forEach((exit) => {
			events.push({
				type: 'exit',
				data: exit,
				timestamp: exit.timestamp
			});
		});

		// Sort by timestamp
		return events.sort(
			(a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
		);
	});

	function formatDate(timestamp: string): string {
		const date = new Date(timestamp);
		return date.toLocaleString('en-US', {
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatPrice(price: number): string {
		return `$${price.toFixed(2)}`;
	}
</script>

<div class="position-timeline">
	<h3 class="text-lg font-semibold mb-4">Position Timeline</h3>

	{#if timelineEvents().length === 0}
		<p class="text-surface-600 dark:text-surface-400 text-sm">No entries or exits recorded yet.</p>
	{:else}
		<div class="relative">
			<!-- Timeline line -->
			<div class="absolute left-6 top-0 bottom-0 w-0.5 bg-surface-300 dark:bg-surface-700"></div>

			<!-- Timeline events -->
			<div class="space-y-6">
				{#each timelineEvents() as event}
					<div class="relative flex gap-4">
						<!-- Icon -->
						<div
							class="relative z-10 flex items-center justify-center w-12 h-12 rounded-full border-4 border-surface-50 dark:border-surface-900 {event.type ===
							'entry'
								? 'bg-blue-500'
								: 'bg-purple-500'}"
						>
							<Icon
								icon={event.type === 'entry'
									? 'mdi:arrow-down-bold'
									: 'mdi:arrow-up-bold'}
								class="text-white text-xl"
							/>
						</div>

						<!-- Event content -->
						<div class="flex-1 pb-6">
							<div class="bg-surface-100 dark:bg-surface-800 rounded-lg p-4">
								<div class="flex justify-between items-start mb-2">
									<div>
										<h4 class="font-semibold text-base">
											{event.type === 'entry' ? 'Entry' : 'Exit'}
										</h4>
										<p class="text-sm text-surface-600 dark:text-surface-400">
											{formatDate(event.timestamp)}
										</p>
									</div>
									<div class="text-right">
										<p class="text-lg font-bold">
											{formatPrice(event.data.price)}
										</p>
										<p class="text-sm text-surface-600 dark:text-surface-400">
											{event.data.quantity} shares
										</p>
									</div>
								</div>

								{#if event.type === 'exit' && 'pnl' in event.data}
									<div class="mt-2 pt-2 border-t border-surface-300 dark:border-surface-700">
										<div class="flex justify-between items-center">
											<span class="text-sm font-medium">P&L:</span>
											<span
												class="text-sm font-bold {event.data.pnl >= 0
													? 'text-profit-600'
													: 'text-loss-600'}"
											>
												{event.data.pnl >= 0 ? '+' : ''}{formatPrice(event.data.pnl)}
											</span>
										</div>
									</div>
								{/if}

								{#if event.data.notes}
									<div class="mt-2 pt-2 border-t border-surface-300 dark:border-surface-700">
										<p class="text-sm text-surface-700 dark:text-surface-300">
											<Icon icon="mdi:note-text" class="inline mr-1" />
											{event.data.notes}
										</p>
									</div>
								{/if}

								{#if event.data.fees && event.data.fees > 0}
									<div class="mt-1">
										<p class="text-xs text-surface-500 dark:text-surface-500">
											Fees: {formatPrice(event.data.fees)}
										</p>
									</div>
								{/if}
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	.position-timeline {
		padding: 1rem 0;
	}
</style>
