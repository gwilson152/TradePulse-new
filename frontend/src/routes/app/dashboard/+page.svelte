<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import PnLBadge from '$lib/components/trading/PnLBadge.svelte';
	import { apiClient } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import type { Trade } from '$lib/types';

	let metrics = $state({
		totalPnL: 0,
		winRate: 0,
		totalTrades: 0,
		avgWin: 0,
		activeTrades: 0
	});

	let recentTrades = $state<Trade[]>([]);
	let isLoading = $state(true);

	onMount(async () => {
		try {
			// Fetch trades and calculate metrics
			const trades = await apiClient.getTrades();
			recentTrades = trades.slice(0, 5); // Show last 5 trades

			// Calculate metrics
			const closedTrades = trades.filter((t) => t.closed_at !== null);
			metrics.totalTrades = trades.length;
			metrics.activeTrades = trades.length - closedTrades.length;
			metrics.totalPnL = closedTrades.reduce((sum, t) => sum + (t.realized_pnl || 0), 0);

			const winningTrades = closedTrades.filter((t) => (t.realized_pnl || 0) > 0);
			metrics.winRate = closedTrades.length > 0 ? (winningTrades.length / closedTrades.length) * 100 : 0;
			metrics.avgWin = winningTrades.length > 0 ? winningTrades.reduce((sum, t) => sum + (t.realized_pnl || 0), 0) / winningTrades.length : 0;
		} catch (err) {
			console.error('Failed to load dashboard data:', err);
		} finally {
			isLoading = false;
		}
	});

	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric'
		});
	}

	function formatCurrency(value: number) {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 0,
			maximumFractionDigits: 0
		}).format(value);
	}
</script>

<svelte:head>
	<title>Overview - TradePulse</title>
</svelte:head>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-4xl font-bold bg-gradient-to-r from-slate-800 to-slate-600 dark:from-slate-100 dark:to-slate-400 bg-clip-text text-transparent mb-2">
				Good {new Date().getHours() < 12 ? 'morning' : new Date().getHours() < 18 ? 'afternoon' : 'evening'}
			</h1>
			<p class="text-slate-600 dark:text-slate-400">Here's what's happening with your trading today</p>
		</div>
		<Button variant="gradient" color="primary" size="lg" onclick={() => toast.info('Trade entry modal coming soon!')}>
			<Icon icon="mdi:plus-circle-outline" width="20" />
			New Trade
		</Button>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-20">
			<div class="flex flex-col items-center gap-4">
				<div class="w-12 h-12 rounded-full border-4 border-blue-200 border-t-blue-500 animate-spin"></div>
				<p class="text-slate-600 dark:text-slate-400">Loading your portfolio...</p>
			</div>
		</div>
	{:else}
		<!-- Metrics Grid -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
			<!-- Total P&L -->
			<Card variant="glass" padding="lg" hover={true}>
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center shadow-lg shadow-blue-500/30">
							<Icon icon="mdi:cash-multiple" width="24" class="text-white" />
						</div>
						<div class="px-3 py-1 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 text-xs font-semibold">
							{metrics.totalPnL >= 0 ? '↑' : '↓'} {Math.abs(metrics.totalPnL / 1000).toFixed(1)}K
						</div>
					</div>
					<div>
						<p class="text-sm text-slate-600 dark:text-slate-400 mb-1">Total P&L</p>
						<p class="text-3xl font-bold {metrics.totalPnL >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}">
							{formatCurrency(metrics.totalPnL)}
						</p>
					</div>
				</div>
			</Card>

			<!-- Win Rate -->
			<Card variant="glass" padding="lg" hover={true}>
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-500 to-teal-500 flex items-center justify-center shadow-lg shadow-emerald-500/30">
							<Icon icon="mdi:target-variant" width="24" class="text-white" />
						</div>
						<div class="px-3 py-1 rounded-full bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400 text-xs font-semibold">
							{metrics.totalTrades} trades
						</div>
					</div>
					<div>
						<p class="text-sm text-slate-600 dark:text-slate-400 mb-1">Win Rate</p>
						<p class="text-3xl font-bold text-slate-800 dark:text-slate-200">
							{metrics.winRate.toFixed(1)}%
						</p>
					</div>
				</div>
			</Card>

			<!-- Average Win -->
			<Card variant="glass" padding="lg" hover={true}>
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center shadow-lg shadow-purple-500/30">
							<Icon icon="mdi:trending-up" width="24" class="text-white" />
						</div>
						<div class="px-3 py-1 rounded-full bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 text-xs font-semibold">
							Per win
						</div>
					</div>
					<div>
						<p class="text-sm text-slate-600 dark:text-slate-400 mb-1">Average Win</p>
						<p class="text-3xl font-bold text-slate-800 dark:text-slate-200">
							{formatCurrency(metrics.avgWin)}
						</p>
					</div>
				</div>
			</Card>

			<!-- Active Trades -->
			<Card variant="glass" padding="lg" hover={true}>
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-orange-500 to-yellow-500 flex items-center justify-center shadow-lg shadow-orange-500/30">
							<Icon icon="mdi:chart-line-variant" width="24" class="text-white" />
						</div>
						<div class="px-3 py-1 rounded-full bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-400 text-xs font-semibold">
							Open
						</div>
					</div>
					<div>
						<p class="text-sm text-slate-600 dark:text-slate-400 mb-1">Active Trades</p>
						<p class="text-3xl font-bold text-slate-800 dark:text-slate-200">
							{metrics.activeTrades}
						</p>
					</div>
				</div>
			</Card>
		</div>

		<!-- Recent Activity -->
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- Recent Trades -->
			<div class="lg:col-span-2">
				<Card variant="glass" padding="lg">
					<div class="flex items-center justify-between mb-6">
						<h2 class="text-xl font-bold text-slate-800 dark:text-slate-200">Recent Trades</h2>
						<Button variant="ghost" color="primary" size="sm" onclick={() => toast.info('Navigate to Trades page to view all')}>
							View all
							<Icon icon="mdi:arrow-right" width="16" />
						</Button>
					</div>

					{#if recentTrades.length === 0}
						<div class="text-center py-12">
							<Icon icon="mdi:chart-line-variant" width="64" class="mx-auto mb-4 text-slate-300 dark:text-slate-700" />
							<p class="text-slate-600 dark:text-slate-400 mb-4">No trades yet</p>
							<Button variant="soft" color="primary" onclick={() => toast.info('Trade entry modal coming soon!')}>
								<Icon icon="mdi:plus" width="18" />
								Add Your First Trade
							</Button>
						</div>
					{:else}
						<div class="space-y-3">
							{#each recentTrades as trade}
								<div class="flex items-center justify-between p-4 rounded-xl bg-slate-50/50 dark:bg-slate-900/30 hover:bg-slate-100/70 dark:hover:bg-slate-800/50 transition-colors cursor-pointer">
									<div class="flex items-center gap-4">
										<div class="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-purple-500 flex items-center justify-center">
											<span class="text-white font-bold text-sm">{trade.symbol.slice(0, 2)}</span>
										</div>
										<div>
											<p class="font-semibold text-slate-800 dark:text-slate-200">{trade.symbol}</p>
											<p class="text-xs text-slate-600 dark:text-slate-400">
												{trade.trade_type} • {formatDate(trade.opened_at)}
											</p>
										</div>
									</div>
									<div class="text-right">
										<PnLBadge value={trade.realized_pnl || 0} showSign={true} />
										<p class="text-xs text-slate-600 dark:text-slate-400 mt-1">
											{trade.current_position_size > 0 ? 'Open' : 'Closed'}
										</p>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</Card>
			</div>

			<!-- Quick Actions -->
			<div class="space-y-6">
				<Card variant="glass" padding="lg">
					<h2 class="text-xl font-bold text-slate-800 dark:text-slate-200 mb-4">Quick Actions</h2>
					<div class="space-y-3">
						<Button variant="soft" color="primary" size="md" class="w-full justify-start" onclick={() => toast.info('Trade entry modal coming soon!')}>
							<Icon icon="mdi:plus-circle-outline" width="20" />
							New Trade
						</Button>
						<Button variant="soft" color="secondary" size="md" class="w-full justify-start" onclick={() => toast.info('Journal entry modal coming soon!')}>
							<Icon icon="mdi:notebook-outline" width="20" />
							Journal Entry
						</Button>
						<Button variant="soft" color="success" size="md" class="w-full justify-start" onclick={() => toast.info('Navigate to Analytics page')}>
							<Icon icon="mdi:chart-box-outline" width="20" />
							View Analytics
						</Button>
					</div>
				</Card>

				<Card variant="glass" padding="lg">
					<h2 class="text-xl font-bold text-slate-800 dark:text-slate-200 mb-4">Market Status</h2>
					<div class="space-y-3">
						<div class="flex items-center justify-between">
							<span class="text-sm text-slate-600 dark:text-slate-400">Market Hours</span>
							<div class="flex items-center gap-2">
								<div class="w-2 h-2 rounded-full bg-red-500 animate-pulse"></div>
								<span class="text-sm font-medium text-slate-800 dark:text-slate-200">Closed</span>
							</div>
						</div>
						<div class="flex items-center justify-between">
							<span class="text-sm text-slate-600 dark:text-slate-400">Next Open</span>
							<span class="text-sm font-medium text-slate-800 dark:text-slate-200">Tomorrow 9:30 AM</span>
						</div>
					</div>
				</Card>
			</div>
		</div>
	{/if}
</div>
