<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		title?: string;
		text: string;
		type?: 'info' | 'tip' | 'warning';
		collapsible?: boolean;
	}

	let { title, text, type = 'info', collapsible = false }: Props = $props();

	let isExpanded = $state(!collapsible);

	const typeStyles = {
		info: {
			bg: 'bg-blue-50 dark:bg-blue-900/20',
			border: 'border-blue-200 dark:border-blue-800',
			icon: 'mdi:information-outline',
			iconColor: 'text-blue-600 dark:text-blue-400'
		},
		tip: {
			bg: 'bg-emerald-50 dark:bg-emerald-900/20',
			border: 'border-emerald-200 dark:border-emerald-800',
			icon: 'mdi:lightbulb-outline',
			iconColor: 'text-emerald-600 dark:text-emerald-400'
		},
		warning: {
			bg: 'bg-amber-50 dark:bg-amber-900/20',
			border: 'border-amber-200 dark:border-amber-800',
			icon: 'mdi:alert-outline',
			iconColor: 'text-amber-600 dark:text-amber-400'
		}
	};

	const style = typeStyles[type];
</script>

<div class="rounded-lg border {style.border} {style.bg} p-4">
	{#if collapsible && title}
		<button
			type="button"
			onclick={() => (isExpanded = !isExpanded)}
			class="flex items-center justify-between w-full text-left"
		>
			<div class="flex items-center gap-2">
				<Icon icon={style.icon} width="20" class={style.iconColor} />
				<span class="text-sm font-semibold text-slate-900 dark:text-slate-100">
					{title}
				</span>
			</div>
			<Icon
				icon={isExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'}
				width="20"
				class="text-slate-600 dark:text-slate-400"
			/>
		</button>
	{:else if title}
		<div class="flex items-center gap-2 mb-2">
			<Icon icon={style.icon} width="20" class={style.iconColor} />
			<span class="text-sm font-semibold text-slate-900 dark:text-slate-100">
				{title}
			</span>
		</div>
	{:else}
		<div class="flex items-start gap-2">
			<Icon icon={style.icon} width="20" class={style.iconColor + ' mt-0.5'} />
			<p class="text-sm text-slate-700 dark:text-slate-300 flex-1">
				{text}
			</p>
		</div>
	{/if}

	{#if (collapsible && isExpanded) || (!collapsible && title)}
		<p class="text-sm text-slate-700 dark:text-slate-300 {title ? 'mt-2 ml-7' : ''}">
			{text}
		</p>
	{/if}
</div>
