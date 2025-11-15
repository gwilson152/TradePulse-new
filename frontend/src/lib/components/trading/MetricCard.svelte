<script lang="ts">
	import Icon from '@iconify/svelte';
	import Card from '../ui/Card.svelte';

	interface Props {
		title: string;
		value: string | number;
		subtitle?: string;
		trend?: 'up' | 'down' | 'neutral';
		icon?: string;
	}

	let { title, value, subtitle, trend, icon }: Props = $props();

	const trendColors = {
		up: 'text-profit-600 dark:text-profit-400',
		down: 'text-loss-600 dark:text-loss-400',
		neutral: 'text-surface-600 dark:text-surface-400'
	};
</script>

<Card hover={true}>
	<div class="flex items-start justify-between">
		<div class="flex-1">
			<p class="text-sm font-medium text-surface-600 dark:text-surface-400">{title}</p>
			<p class="text-2xl font-bold mt-2">{value}</p>
			{#if subtitle}
				<p class="text-sm mt-1 {trend ? trendColors[trend] : 'text-surface-500'}">{subtitle}</p>
			{/if}
		</div>
		{#if icon}
			<div class="text-surface-500 dark:text-surface-400">
				<Icon icon={icon} width="32" />
			</div>
		{/if}
	</div>
</Card>
