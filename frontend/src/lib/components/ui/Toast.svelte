<script lang="ts">
	import { toast, type Toast } from '$lib/stores/toast';
	import Icon from '@iconify/svelte';
	import { fade, fly } from 'svelte/transition';

	let toasts = $state<Toast[]>([]);

	toast.subscribe((value) => {
		toasts = value;
	});

	function getIcon(type: Toast['type']): string {
		switch (type) {
			case 'success':
				return 'mdi:check-circle';
			case 'error':
				return 'mdi:alert-circle';
			case 'warning':
				return 'mdi:alert';
			case 'info':
				return 'mdi:information';
		}
	}

	function getStyles(type: Toast['type']): string {
		switch (type) {
			case 'success':
				return 'bg-emerald-500 dark:bg-emerald-600 text-white';
			case 'error':
				return 'bg-red-500 dark:bg-red-600 text-white';
			case 'warning':
				return 'bg-amber-500 dark:bg-amber-600 text-white';
			case 'info':
				return 'bg-blue-500 dark:bg-blue-600 text-white';
		}
	}
</script>

<div class="fixed top-16 right-4 z-[200] space-y-2 pointer-events-none">
	{#each toasts as t (t.id)}
		<div
			transition:fly={{ x: 300, duration: 300 }}
			class="pointer-events-auto {getStyles(t.type)} rounded-xl px-4 py-3 shadow-xl shadow-black/20 flex items-center gap-3 min-w-[320px] max-w-md backdrop-blur-xl"
		>
			<Icon icon={getIcon(t.type)} width="24" />
			<p class="flex-1 text-sm font-medium">{t.message}</p>
			<button
				onclick={() => toast.remove(t.id)}
				class="p-1 hover:bg-white/20 rounded transition-colors"
				aria-label="Close"
			>
				<Icon icon="mdi:close" width="18" />
			</button>
		</div>
	{/each}
</div>
