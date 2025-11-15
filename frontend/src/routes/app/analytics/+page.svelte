<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import ChartCard from '$lib/components/ui/ChartCard.svelte';
	import MetricCard from '$lib/components/trading/MetricCard.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Icon from '@iconify/svelte';
	import HelpText from '$lib/components/ui/HelpText.svelte';
	import type * as echarts from 'echarts';
	import { onMount } from 'svelte';
	import {
		calculateMetrics,
		analyzeTimeOfDay,
		analyzeHoldTimes,
		analyzePriceRanges,
		analyzeTradeSizes,
		analyzeStreaks
	} from '$lib/utils/analytics';

	let timeRange = $state('30d');
	let loading = $state(true);

	// Mock data - replace with actual API calls
	let metrics = $state({
		winRate: 65.5,
		avgWin: 245.32,
		avgLoss: -128.45,
		riskReward: 1.91,
		totalTrades: 42,
		profitFactor: 2.1,
		avgHoldTime: '2h 35m',
		bestHoldTime: '1-3 hours',
		bestPriceRange: '$50-$150',
		bestTimeOfDay: '10:00-11:00 AM'
	});

	onMount(async () => {
		// TODO: Load analytics data from API
		loading = false;
	});

	// P&L Over Time Chart
	const pnlChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'axis',
				axisPointer: { type: 'cross' }
			},
			xAxis: {
				type: 'category',
				data: ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Week 5'],
				axisLabel: { color: '#888' }
			},
			yAxis: {
				type: 'value',
				axisLabel: {
					color: '#888',
					formatter: '${value}'
				}
			},
			series: [
				{
					name: 'Cumulative P&L',
					type: 'line',
					data: [0, 450, 320, 850, 1240],
					smooth: true,
					lineStyle: { color: '#10b981', width: 3 },
					areaStyle: {
						color: {
							type: 'linear',
							x: 0,
							y: 0,
							x2: 0,
							y2: 1,
							colorStops: [
								{ offset: 0, color: 'rgba(16, 185, 129, 0.3)' },
								{ offset: 1, color: 'rgba(16, 185, 129, 0)' }
							]
						}
					}
				}
			],
			grid: { left: 60, right: 20, top: 20, bottom: 40 }
		};
	});

	// Win Rate by Day of Week
	const dayOfWeekChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: { trigger: 'axis' },
			xAxis: {
				type: 'category',
				data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri'],
				axisLabel: { color: '#888' }
			},
			yAxis: {
				type: 'value',
				max: 100,
				axisLabel: {
					color: '#888',
					formatter: '{value}%'
				}
			},
			series: [
				{
					name: 'Win Rate',
					type: 'bar',
					data: [72, 65, 58, 70, 68],
					itemStyle: {
						color: (params: any) => {
							const value = params.value;
							return value >= 70 ? '#10b981' : value >= 50 ? '#f59e0b' : '#ef4444';
						}
					},
					label: {
						show: true,
						position: 'top',
						formatter: '{c}%',
						color: '#888'
					}
				}
			],
			grid: { left: 60, right: 20, top: 20, bottom: 40 }
		};
	});

	// Trade Outcomes Distribution
	const outcomeDistributionOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'item',
				formatter: '{b}: {c} ({d}%)'
			},
			legend: {
				orient: 'vertical',
				right: 10,
				top: 'center',
				textStyle: { color: '#888' }
			},
			series: [
				{
					name: 'Trade Outcomes',
					type: 'pie',
					radius: ['40%', '70%'],
					avoidLabelOverlap: false,
					label: {
						show: true,
						formatter: '{b}\n{d}%'
					},
					data: [
						{ value: 27, name: 'Wins', itemStyle: { color: '#10b981' } },
						{ value: 15, name: 'Losses', itemStyle: { color: '#ef4444' } }
					]
				}
			]
		};
	});

	// Rule Adherence Correlation
	const ruleAdherenceChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'axis',
				axisPointer: { type: 'shadow' }
			},
			legend: {
				data: ['Avg P&L', 'Adherence %'],
				textStyle: { color: '#888' }
			},
			xAxis: {
				type: 'category',
				data: ['0-20%', '21-40%', '41-60%', '61-80%', '81-100%'],
				axisLabel: { color: '#888' }
			},
			yAxis: [
				{
					type: 'value',
					name: 'Avg P&L',
					position: 'left',
					axisLabel: {
						color: '#888',
						formatter: '${value}'
					}
				},
				{
					type: 'value',
					name: 'Count',
					position: 'right',
					axisLabel: { color: '#888' }
				}
			],
			series: [
				{
					name: 'Avg P&L',
					type: 'bar',
					data: [-50, 20, 150, 280, 420],
					itemStyle: {
						color: (params: any) => params.value >= 0 ? '#10b981' : '#ef4444'
					}
				},
				{
					name: 'Trade Count',
					type: 'line',
					yAxisIndex: 1,
					data: [2, 5, 10, 15, 10],
					lineStyle: { color: '#3b82f6', width: 2 },
					itemStyle: { color: '#3b82f6' }
				}
			],
			grid: { left: 60, right: 60, top: 60, bottom: 40 }
		};
	});

	// Emotional State vs Performance
	const emotionalStateChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'item',
				formatter: 'Confidence: {c[0]}<br/>Avg P&L: ${c[1]}'
			},
			xAxis: {
				type: 'value',
				name: 'Confidence Level',
				min: 0,
				max: 10,
				axisLabel: { color: '#888' }
			},
			yAxis: {
				type: 'value',
				name: 'Avg P&L',
				axisLabel: {
					color: '#888',
					formatter: '${value}'
				}
			},
			series: [
				{
					name: 'Trades',
					type: 'scatter',
					symbolSize: (data: number[]) => Math.abs(data[1]) / 10,
					data: [
						[3, -50], [5, 100], [6, 200], [7, 350], [8, 420],
						[4, -30], [7, 280], [9, 450], [6, 150], [8, 380]
					],
					itemStyle: {
						color: (params: any) => params.value[1] >= 0 ? '#10b981' : '#ef4444',
						opacity: 0.6
					}
				}
			],
			grid: { left: 80, right: 20, top: 40, bottom: 60 }
		};
	});

	// Time of Day Performance Chart
	const timeOfDayChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'axis',
				axisPointer: { type: 'shadow' }
			},
			legend: {
				data: ['Win Rate %', 'Avg P&L'],
				textStyle: { color: '#888' }
			},
			xAxis: {
				type: 'category',
				data: ['9-10 AM', '10-11 AM', '11-12 PM', '12-1 PM', '1-2 PM', '2-3 PM', '3-4 PM'],
				axisLabel: {
					color: '#888',
					rotate: 30
				}
			},
			yAxis: [
				{
					type: 'value',
					name: 'Win Rate %',
					position: 'left',
					max: 100,
					axisLabel: {
						color: '#888',
						formatter: '{value}%'
					}
				},
				{
					type: 'value',
					name: 'Avg P&L',
					position: 'right',
					axisLabel: {
						color: '#888',
						formatter: '${value}'
					}
				}
			],
			series: [
				{
					name: 'Win Rate %',
					type: 'bar',
					data: [55, 72, 68, 45, 62, 58, 50],
					itemStyle: {
						color: (params: any) => {
							const value = params.value;
							return value >= 70 ? '#10b981' : value >= 50 ? '#f59e0b' : '#ef4444';
						}
					}
				},
				{
					name: 'Avg P&L',
					type: 'line',
					yAxisIndex: 1,
					data: [120, 280, 230, 80, 180, 150, 90],
					lineStyle: { color: '#3b82f6', width: 3 },
					itemStyle: { color: '#3b82f6' },
					smooth: true
				}
			],
			grid: { left: 60, right: 60, top: 60, bottom: 80 }
		};
	});

	// Hold Time Distribution Chart
	const holdTimeChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'axis',
				axisPointer: { type: 'shadow' }
			},
			legend: {
				data: ['Trades', 'Avg P&L'],
				textStyle: { color: '#888' }
			},
			xAxis: {
				type: 'category',
				data: ['<30m', '30m-1h', '1-3h', '3-6h', '6h-1d', '1-3d', '>3d'],
				axisLabel: { color: '#888' }
			},
			yAxis: [
				{
					type: 'value',
					name: 'Trade Count',
					position: 'left',
					axisLabel: { color: '#888' }
				},
				{
					type: 'value',
					name: 'Avg P&L',
					position: 'right',
					axisLabel: {
						color: '#888',
						formatter: '${value}'
					}
				}
			],
			series: [
				{
					name: 'Trades',
					type: 'bar',
					data: [8, 12, 18, 15, 10, 6, 3],
					itemStyle: { color: '#8b5cf6' }
				},
				{
					name: 'Avg P&L',
					type: 'line',
					yAxisIndex: 1,
					data: [85, 180, 320, 280, 195, 120, -50],
					lineStyle: { color: '#10b981', width: 3 },
					itemStyle: { color: '#10b981' },
					smooth: true,
					areaStyle: {
						color: {
							type: 'linear',
							x: 0,
							y: 0,
							x2: 0,
							y2: 1,
							colorStops: [
								{ offset: 0, color: 'rgba(16, 185, 129, 0.3)' },
								{ offset: 1, color: 'rgba(16, 185, 129, 0)' }
							]
						}
					}
				}
			],
			grid: { left: 60, right: 60, top: 60, bottom: 40 }
		};
	});

	// Price Range Performance Chart
	const priceRangeChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'axis',
				axisPointer: { type: 'shadow' }
			},
			xAxis: {
				type: 'category',
				data: ['<$20', '$20-$50', '$50-$100', '$100-$200', '$200-$500', '>$500'],
				axisLabel: { color: '#888' }
			},
			yAxis: {
				type: 'value',
				name: 'Win Rate %',
				max: 100,
				axisLabel: {
					color: '#888',
					formatter: '{value}%'
				}
			},
			series: [
				{
					name: 'Win Rate',
					type: 'bar',
					data: [52, 58, 72, 68, 61, 55],
					itemStyle: {
						color: (params: any) => {
							const value = params.value;
							return value >= 70 ? '#10b981' : value >= 50 ? '#f59e0b' : '#ef4444';
						}
					},
					label: {
						show: true,
						position: 'top',
						formatter: '{c}%',
						color: '#888'
					}
				}
			],
			grid: { left: 80, right: 20, top: 40, bottom: 40 }
		};
	});

	// Trade Size vs Performance
	const tradeSizeChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'item',
				formatter: 'Size: ${c[0]}<br/>P&L: ${c[1]}<br/>Win: {c[2]}'
			},
			xAxis: {
				type: 'value',
				name: 'Position Size ($)',
				axisLabel: {
					color: '#888',
					formatter: '${value}'
				}
			},
			yAxis: {
				type: 'value',
				name: 'P&L ($)',
				axisLabel: {
					color: '#888',
					formatter: '${value}'
				}
			},
			series: [
				{
					name: 'Trades',
					type: 'scatter',
					symbolSize: 15,
					data: [
						[500, 120, 'Win'], [1000, 280, 'Win'], [1500, -150, 'Loss'],
						[2000, 420, 'Win'], [2500, 350, 'Win'], [3000, -200, 'Loss'],
						[800, 180, 'Win'], [1200, 240, 'Win'], [1800, -120, 'Loss'],
						[2200, 380, 'Win'], [2800, 450, 'Win'], [3500, 520, 'Win'],
						[700, 150, 'Win'], [1400, -180, 'Loss'], [2100, 310, 'Win']
					],
					itemStyle: {
						color: (params: any) => params.value[2] === 'Win' ? '#10b981' : '#ef4444',
						opacity: 0.7
					}
				}
			],
			grid: { left: 80, right: 20, top: 40, bottom: 60 }
		};
	});

	// Win/Loss Streak Analysis
	const streakChartOptions = $derived((): echarts.EChartsOption => {
		return {
			tooltip: {
				trigger: 'axis',
				axisPointer: { type: 'shadow' }
			},
			xAxis: {
				type: 'category',
				data: ['Single', '2 in a row', '3 in a row', '4 in a row', '5+ in a row'],
				axisLabel: { color: '#888' }
			},
			yAxis: {
				type: 'value',
				name: 'Frequency',
				axisLabel: { color: '#888' }
			},
			series: [
				{
					name: 'Win Streaks',
					type: 'bar',
					stack: 'streaks',
					data: [15, 8, 4, 2, 1],
					itemStyle: { color: '#10b981' },
					label: {
						show: true,
						position: 'inside',
						color: '#fff'
					}
				},
				{
					name: 'Loss Streaks',
					type: 'bar',
					stack: 'streaks',
					data: [10, 4, 2, 0, 0],
					itemStyle: { color: '#ef4444' },
					label: {
						show: true,
						position: 'inside',
						color: '#fff'
					}
				}
			],
			grid: { left: 60, right: 20, top: 40, bottom: 40 }
		};
	});

	const timeRangeOptions = [
		{ value: '7d', label: 'Last 7 Days' },
		{ value: '30d', label: 'Last 30 Days' },
		{ value: '90d', label: 'Last 90 Days' },
		{ value: 'ytd', label: 'Year to Date' },
		{ value: 'all', label: 'All Time' }
	];
</script>

<svelte:head>
	<title>Analytics - TradePulse</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold mb-2">Analytics</h1>
			<p class="text-surface-600 dark:text-surface-400">
				Deep dive into your trading performance
			</p>
		</div>
		<div class="w-48">
			<Select
				options={timeRangeOptions}
				bind:value={timeRange}
				placeholder="Select range"
			/>
		</div>
	</div>

	<!-- Key Insights -->
	<HelpText
		type="info"
		title="Advanced Analytics"
		text="These charts analyze your trading patterns to identify optimal times, price ranges, and hold durations. Use these insights to refine your strategy and focus on what works best for you."
		collapsible={true}
	/>

	<!-- Key Metrics -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
		<MetricCard
			title="Win Rate"
			value="{metrics.winRate}%"
			icon="mdi:bullseye-arrow"
			color="success"
			trend="up"
			trendValue="+3.2%"
		/>
		<MetricCard
			title="Avg Win"
			value="${metrics.avgWin}"
			icon="mdi:trending-up"
			color="success"
		/>
		<MetricCard
			title="Avg Loss"
			value="${Math.abs(metrics.avgLoss)}"
			icon="mdi:trending-down"
			color="error"
		/>
		<MetricCard
			title="Risk/Reward"
			value={metrics.riskReward.toFixed(2)}
			icon="mdi:scale-balance"
			color="primary"
		/>
	</div>

	<!-- P&L Over Time -->
	<ChartCard
		title="Cumulative P&L"
		subtitle="Your profit and loss over time"
		options={pnlChartOptions()}
		{loading}
		height="350px"
	/>

	<!-- Two Column Charts -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
		<!-- Win Rate by Day -->
		<ChartCard
			title="Win Rate by Day of Week"
			subtitle="Performance breakdown by trading day"
			options={dayOfWeekChartOptions()}
			{loading}
			height="300px"
		/>

		<!-- Trade Outcomes -->
		<ChartCard
			title="Trade Outcomes"
			subtitle="Distribution of wins vs losses"
			options={outcomeDistributionOptions()}
			{loading}
			height="300px"
		/>
	</div>

	<!-- Rule Adherence Correlation -->
	<ChartCard
		title="Rule Adherence vs Performance"
		subtitle="How following your rules affects profitability"
		options={ruleAdherenceChartOptions()}
		{loading}
		height="350px"
	/>

	<!-- Emotional State Correlation -->
	<ChartCard
		title="Emotional State vs Performance"
		subtitle="Correlation between confidence and trade outcomes"
		options={emotionalStateChartOptions()}
		{loading}
		height="350px"
	/>

	<!-- Time of Day Analysis -->
	<ChartCard
		title="Performance by Time of Day"
		subtitle="Win rate and P&L across different trading hours"
		options={timeOfDayChartOptions()}
		{loading}
		height="350px"
	/>

	<!-- Hold Time Analysis -->
	<ChartCard
		title="Hold Time Distribution & Performance"
		subtitle="How long you hold positions affects profitability"
		options={holdTimeChartOptions()}
		{loading}
		height="350px"
	/>

	<!-- Price Range Performance -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
		<ChartCard
			title="Win Rate by Price Range"
			subtitle="Performance across different stock prices"
			options={priceRangeChartOptions()}
			{loading}
			height="300px"
		/>

		<!-- Trade Size vs Performance -->
		<ChartCard
			title="Position Size vs P&L"
			subtitle="Relationship between position size and outcomes"
			options={tradeSizeChartOptions()}
			{loading}
			height="300px"
		/>
	</div>

	<!-- Win/Loss Streaks -->
	<ChartCard
		title="Win/Loss Streak Analysis"
		subtitle="Frequency of consecutive wins and losses"
		options={streakChartOptions()}
		{loading}
		height="300px"
	/>

	<!-- Additional Insights -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
		<Card>
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-surface-600 dark:text-surface-400 mb-1">Total Trades</p>
					<p class="text-3xl font-bold">{metrics.totalTrades}</p>
				</div>
				<Icon icon="mdi:chart-line" class="text-4xl text-primary-500" />
			</div>
		</Card>
		<Card>
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-surface-600 dark:text-surface-400 mb-1">Profit Factor</p>
					<p class="text-3xl font-bold">{metrics.profitFactor}</p>
				</div>
				<Icon icon="mdi:finance" class="text-4xl text-success-500" />
			</div>
		</Card>
		<Card>
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-surface-600 dark:text-surface-400 mb-1">Avg Hold Time</p>
					<p class="text-3xl font-bold">{metrics.avgHoldTime}</p>
					<p class="text-xs text-surface-500 dark:text-surface-400 mt-1">
						Best: {metrics.bestHoldTime}
					</p>
				</div>
				<Icon icon="mdi:clock-outline" class="text-4xl text-blue-500" />
			</div>
		</Card>
		<Card>
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-surface-600 dark:text-surface-400 mb-1">Best Time</p>
					<p class="text-2xl font-bold">{metrics.bestTimeOfDay}</p>
					<p class="text-xs text-surface-500 dark:text-surface-400 mt-1">
						Price: {metrics.bestPriceRange}
					</p>
				</div>
				<Icon icon="mdi:trophy" class="text-4xl text-warning-500" />
			</div>
		</Card>
	</div>
</div>
