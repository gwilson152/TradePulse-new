<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Tab {
		id: string;
		label: string;
		icon: string;
		isComplete?: boolean;
	}

	interface Props {
		tabs: Tab[];
		activeTab: string;
		onTabChange: (tabId: string) => void;
		children?: any;
	}

	let { tabs, activeTab = $bindable(), onTabChange, children }: Props = $props();

	function handleTabClick(tabId: string) {
		activeTab = tabId;
		if (onTabChange) {
			onTabChange(tabId);
		}
	}
</script>

<div class="tab-container">
	<!-- Tab Headers -->
	<div class="border-b border-slate-200 dark:border-slate-700 -mx-8 px-8">
		<div class="flex gap-1 -mb-px">
			{#each tabs as tab}
				<button
					type="button"
					onclick={() => handleTabClick(tab.id)}
					class="group relative flex items-center gap-2.5 px-5 py-3 border-b-2 transition-all duration-200
						{activeTab === tab.id
							? 'border-blue-500 text-blue-600 dark:text-blue-400 bg-blue-50/50 dark:bg-blue-900/20'
							: 'border-transparent text-slate-600 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-800/30'}
						rounded-t-lg font-medium text-sm"
				>
					<Icon icon={tab.icon} width="18" />
					<span class="capitalize">{tab.label}</span>
					{#if tab.isComplete}
						<Icon icon="mdi:check-circle" width="16" class="text-emerald-500" />
					{/if}
				</button>
			{/each}
		</div>
	</div>

	<!-- Tab Content -->
	<div class="tab-content py-6">
		{@render children?.()}
	</div>
</div>
