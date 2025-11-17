<script lang="ts">
	import { slideOverStore } from '$lib/stores/slideOver';
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';

	let state = $state($slideOverStore);

	$effect(() => {
		state = $slideOverStore;
	});

	function handleEscape(event: KeyboardEvent) {
		if (event.key === 'Escape' && state.panels.length > 0) {
			slideOverStore.close();
		}
	}

	onMount(() => {
		window.addEventListener('keydown', handleEscape);
		return () => window.removeEventListener('keydown', handleEscape);
	});

	const sizeClasses = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-xl',
		'2xl': 'max-w-2xl',
		full: 'max-w-full'
	};
</script>

{#if state.panels.length > 0}
	<!-- Backdrop -->
	<div
		role="button"
		tabindex="0"
		class="fixed inset-0 bg-black/50 backdrop-blur-sm z-[60] transition-opacity"
		onclick={() => slideOverStore.close()}
		onkeydown={(e) => e.key === 'Enter' && slideOverStore.close()}
	></div>

	<!-- Panels Container -->
	<div class="fixed inset-y-0 right-0 z-[70] flex pointer-events-none">
		{#each state.panels as panel, index}
			{@const isActive = index === state.activeIndex}
			{@const offset = (state.panels.length - 1 - index) * 20}
			<div
				role="button"
				tabindex="0"
				class="absolute inset-0 transition-all duration-300 pointer-events-auto {sizeClasses[
					panel.size || 'xl'
				]}"
				style="transform: translateX({isActive ? 0 : offset}px); opacity: {isActive ? 1 : 0.7}; z-index: {index}"
				onclick={() => slideOverStore.setActive(index)}
				onkeydown={(e) => e.key === 'Enter' && slideOverStore.setActive(index)}
			>
				<div
					class="h-full bg-white dark:bg-slate-900 shadow-2xl overflow-hidden flex flex-col"
				>
					<!-- Header -->
					<div
						class="flex items-center justify-between px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-gradient-to-r from-slate-50 to-slate-100 dark:from-slate-800 dark:to-slate-900"
					>
						<div class="flex items-center gap-2">
							{#if state.panels.length > 1}
								<span
									class="text-xs font-semibold text-slate-500 dark:text-slate-400 bg-slate-200 dark:bg-slate-700 px-2 py-1 rounded"
								>
									{index + 1} / {state.panels.length}
								</span>
							{/if}
						</div>
						<button
							onclick={() => slideOverStore.close(panel.id)}
							class="p-2 hover:bg-slate-200 dark:hover:bg-slate-800 rounded-xl transition-colors"
						>
							<Icon icon="mdi:close" width="24" class="text-slate-600 dark:text-slate-400" />
						</button>
					</div>

					<!-- Content -->
					<div class="flex-1 overflow-y-auto">
						{@const Component = panel.component}
						<Component {...panel.props || {}} />
					</div>
				</div>
			</div>
		{/each}
	</div>
{/if}
