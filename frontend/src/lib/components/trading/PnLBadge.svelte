<script lang="ts">
	interface Props {
		value: number;
		showSign?: boolean;
		size?: 'sm' | 'md' | 'lg';
	}

	let { value, showSign = true, size = 'md' }: Props = $props();

	const isProfit = value >= 0;
	const displayValue = showSign && value > 0 ? `+${value.toFixed(2)}` : value.toFixed(2);

	const sizeClasses = {
		sm: 'text-sm px-2 py-1',
		md: 'text-base px-3 py-1.5',
		lg: 'text-lg px-4 py-2'
	};
</script>

<span
	class={`badge font-semibold rounded-full ${isProfit ? 'bg-profit-100 text-profit-700 dark:bg-profit-900/30 dark:text-profit-400' : 'bg-loss-100 text-loss-700 dark:bg-loss-900/30 dark:text-loss-400'} ${sizeClasses[size]}`}
>
	${displayValue}
</span>
