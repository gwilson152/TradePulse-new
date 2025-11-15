<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		text: string;
		position?: 'top' | 'bottom' | 'left' | 'right';
		maxWidth?: string;
		children?: any;
	}

	let { text, position = 'top', maxWidth = '200px', children }: Props = $props();

	let isVisible = $state(false);

	const positionClasses = {
		top: 'bottom-full left-1/2 -translate-x-1/2 mb-2',
		bottom: 'top-full left-1/2 -translate-x-1/2 mt-2',
		left: 'right-full top-1/2 -translate-y-1/2 mr-2',
		right: 'left-full top-1/2 -translate-y-1/2 ml-2'
	};

	const arrowClasses = {
		top: 'top-full left-1/2 -translate-x-1/2 border-l-transparent border-r-transparent border-b-transparent border-t-slate-800 dark:border-t-slate-700',
		bottom: 'bottom-full left-1/2 -translate-x-1/2 border-l-transparent border-r-transparent border-t-transparent border-b-slate-800 dark:border-b-slate-700',
		left: 'left-full top-1/2 -translate-y-1/2 border-t-transparent border-b-transparent border-r-transparent border-l-slate-800 dark:border-l-slate-700',
		right: 'right-full top-1/2 -translate-y-1/2 border-t-transparent border-b-transparent border-l-transparent border-r-slate-800 dark:border-r-slate-700'
	};
</script>

<div class="relative inline-block">
	<div
		onmouseenter={() => (isVisible = true)}
		onmouseleave={() => (isVisible = false)}
		onfocus={() => (isVisible = true)}
		onblur={() => (isVisible = false)}
		role="button"
		tabindex="0"
		aria-label="Show tooltip"
	>
		{@render children?.()}
	</div>

	{#if isVisible}
		<div
			class="absolute z-50 {positionClasses[position]} pointer-events-none"
			style="max-width: {maxWidth}"
			role="tooltip"
		>
			<div
				class="bg-slate-800 dark:bg-slate-700 text-white text-xs rounded-lg px-3 py-2 shadow-xl"
			>
				{text}
			</div>
			<div class="absolute w-0 h-0 border-4 {arrowClasses[position]}"></div>
		</div>
	{/if}
</div>
