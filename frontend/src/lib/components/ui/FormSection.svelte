<script lang="ts">
	import Icon from '@iconify/svelte';
	import Tooltip from './Tooltip.svelte';

	interface Props {
		title: string;
		subtitle?: string;
		icon?: string;
		helpText?: string;
		collapsible?: boolean;
		defaultOpen?: boolean;
		children?: any;
	}

	let {
		title,
		subtitle,
		icon,
		helpText,
		collapsible = false,
		defaultOpen = true,
		children
	}: Props = $props();

	let isOpen = $state(defaultOpen);
</script>

<div class="form-section">
	{#if collapsible}
		<button
			type="button"
			onclick={() => (isOpen = !isOpen)}
			class="w-full flex items-center justify-between py-3 hover:bg-slate-50 dark:hover:bg-slate-800/50 rounded-lg px-2 -mx-2 transition-colors"
		>
			<div class="flex items-center gap-3">
				{#if icon}
					<div class="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center shadow-sm">
						<Icon {icon} width="18" class="text-white" />
					</div>
				{/if}
				<div class="text-left">
					<h3 class="text-lg font-semibold text-slate-800 dark:text-slate-200">
						{title}
					</h3>
					{#if subtitle}
						<p class="text-xs text-slate-600 dark:text-slate-400">
							{subtitle}
						</p>
					{/if}
				</div>
			</div>
			<div class="flex items-center gap-2">
				{#if helpText}
					<Tooltip text={helpText} position="left" maxWidth="300px">
						<Icon
							icon="mdi:help-circle-outline"
							width="18"
							class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 cursor-help"
						/>
					</Tooltip>
				{/if}
				<Icon
					icon={isOpen ? 'mdi:chevron-up' : 'mdi:chevron-down'}
					width="20"
					class="text-slate-400"
				/>
			</div>
		</button>
	{:else}
		<div class="flex items-center gap-3 mb-4">
			{#if icon}
				<div class="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center shadow-sm">
					<Icon {icon} width="18" class="text-white" />
				</div>
			{/if}
			<div class="flex-1">
				<div class="flex items-center gap-2">
					<h3 class="text-lg font-semibold text-slate-800 dark:text-slate-200">
						{title}
					</h3>
					{#if helpText}
						<Tooltip text={helpText} position="right" maxWidth="300px">
							<Icon
								icon="mdi:help-circle-outline"
								width="18"
								class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 cursor-help"
							/>
						</Tooltip>
					{/if}
				</div>
				{#if subtitle}
					<p class="text-xs text-slate-600 dark:text-slate-400 mt-0.5">
						{subtitle}
					</p>
				{/if}
			</div>
		</div>
	{/if}

	{#if (!collapsible || isOpen)}
		<div class="space-y-4 {collapsible ? 'mt-4' : ''}">
			{@render children?.()}
		</div>
	{/if}
</div>
