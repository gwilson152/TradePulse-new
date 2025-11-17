<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		size?: 'sm' | 'md' | 'lg';
		showIcon?: boolean;
	}

	let { size = 'md', showIcon = true }: Props = $props();

	const sizeClasses = {
		sm: 'text-xs px-2 py-0.5',
		md: 'text-sm px-2.5 py-1',
		lg: 'text-base px-3 py-1.5'
	};
</script>

<span
	class="inline-flex items-center gap-1 font-semibold rounded-full bg-gradient-to-r from-purple-500 to-pink-500 text-white {sizeClasses[size]}"
>
	{#if showIcon}
		<Icon icon="mdi:flask-outline" class="text-current" />
	{/if}
	<span>BETA</span>
</span>

<style>
	span {
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
		box-shadow: 0 2px 8px rgba(168, 85, 247, 0.3);
	}
</style>
