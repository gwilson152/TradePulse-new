<script lang="ts">
	import Card from './Card.svelte';
	import Icon from '@iconify/svelte';
	import * as echarts from 'echarts';
	import { onMount, onDestroy } from 'svelte';

	interface Props {
		title: string;
		subtitle?: string;
		options: echarts.EChartsOption;
		loading?: boolean;
		height?: string;
	}

	let { title, subtitle, options, loading = false, height = '400px' }: Props = $props();

	let chartContainer: HTMLDivElement;
	let chartInstance: echarts.ECharts | null = null;

	onMount(() => {
		if (chartContainer) {
			chartInstance = echarts.init(chartContainer);
			chartInstance.setOption(options);

			// Handle window resize
			const resizeObserver = new ResizeObserver(() => {
				chartInstance?.resize();
			});
			resizeObserver.observe(chartContainer);

			return () => {
				resizeObserver.disconnect();
			};
		}
	});

	// Update chart when options change
	$effect(() => {
		if (chartInstance && options) {
			chartInstance.setOption(options, { notMerge: true });
		}
	});

	onDestroy(() => {
		chartInstance?.dispose();
	});
</script>

<Card>
	<div class="mb-4">
		<h3 class="text-lg font-semibold text-surface-900 dark:text-surface-100">{title}</h3>
		{#if subtitle}
			<p class="text-sm text-surface-600 dark:text-surface-400">{subtitle}</p>
		{/if}
	</div>

	<div class="relative" style="height: {height}">
		{#if loading}
			<div class="absolute inset-0 flex items-center justify-center bg-surface-50 dark:bg-surface-900 rounded">
				<Icon icon="mdi:loading" class="animate-spin text-4xl text-primary-500" />
			</div>
		{/if}
		<div bind:this={chartContainer} style="width: 100%; height: 100%;"></div>
	</div>
</Card>
