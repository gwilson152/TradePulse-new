<script lang="ts">
	import type { Rule, RulePhase, RuleCategory } from '$lib/types';
	import Icon from '@iconify/svelte';

	interface Props {
		rule: Rule;
		editable?: boolean;
		onEdit?: () => void;
		onDelete?: () => void;
	}

	let { rule, editable = false, onEdit, onDelete }: Props = $props();

	function getPhaseColor(phase: RulePhase): string {
		switch (phase) {
			case 'PRE_TRADE':
				return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300';
			case 'DURING_TRADE':
				return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-300';
			case 'POST_TRADE':
				return 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300';
			default:
				return 'bg-surface-100 text-surface-700 dark:bg-surface-800 dark:text-surface-300';
		}
	}

	function getPhaseLabel(phase: RulePhase): string {
		switch (phase) {
			case 'PRE_TRADE':
				return 'Pre-Trade';
			case 'DURING_TRADE':
				return 'During Trade';
			case 'POST_TRADE':
				return 'Post-Trade';
			default:
				return phase;
		}
	}

	function getCategoryIcon(category: RuleCategory): string {
		switch (category) {
			case 'RISK_MANAGEMENT':
				return 'mdi:shield-check';
			case 'ENTRY':
				return 'mdi:login';
			case 'EXIT':
				return 'mdi:logout';
			case 'POSITION_SIZING':
				return 'mdi:chart-box';
			case 'TIMING':
				return 'mdi:clock-outline';
			case 'PSYCHOLOGY':
				return 'mdi:brain';
			case 'GENERAL':
				return 'mdi:text-box-check';
			default:
				return 'mdi:text-box-check';
		}
	}

	function getCategoryLabel(category: RuleCategory): string {
		return category
			.split('_')
			.map((word) => word.charAt(0) + word.slice(1).toLowerCase())
			.join(' ');
	}

	function renderStars(weight: number): string {
		return '★'.repeat(weight) + '☆'.repeat(5 - weight);
	}
</script>

<div
	class="rule-card group bg-surface-50 dark:bg-surface-800 rounded-lg border border-surface-200 dark:border-surface-700 hover:border-surface-300 dark:hover:border-surface-600 transition-colors p-4"
>
	<div class="flex items-start justify-between mb-2">
		<div class="flex items-start gap-3 flex-1">
			<div class="mt-1">
				<Icon icon={getCategoryIcon(rule.category)} class="text-2xl text-surface-600 dark:text-surface-400" />
			</div>
			<div class="flex-1">
				<h4 class="font-semibold text-base text-surface-900 dark:text-surface-100 mb-1">
					{rule.title}
				</h4>
				<p class="text-sm text-surface-600 dark:text-surface-400">
					{rule.description}
				</p>
			</div>
		</div>

		{#if editable}
			<div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
				<button
					onclick={onEdit}
					class="p-1.5 hover:bg-surface-200 dark:hover:bg-surface-700 rounded transition-colors"
					title="Edit rule"
				>
					<Icon icon="mdi:pencil" class="text-surface-600 dark:text-surface-400" />
				</button>
				<button
					onclick={onDelete}
					class="p-1.5 hover:bg-error-100 dark:hover:bg-error-900/30 rounded transition-colors"
					title="Delete rule"
				>
					<Icon icon="mdi:delete" class="text-error-600 dark:text-error-400" />
				</button>
			</div>
		{/if}
	</div>

	<div class="flex items-center gap-3 mt-3">
		<!-- Phase badge -->
		<span class="text-xs font-medium px-2 py-1 rounded {getPhaseColor(rule.phase)}">
			{getPhaseLabel(rule.phase)}
		</span>

		<!-- Category badge -->
		<span
			class="text-xs font-medium px-2 py-1 rounded bg-surface-200 text-surface-700 dark:bg-surface-700 dark:text-surface-300"
		>
			{getCategoryLabel(rule.category)}
		</span>

		<!-- Weight (importance) -->
		<div class="flex items-center gap-1 ml-auto">
			<span class="text-xs text-surface-600 dark:text-surface-400">Importance:</span>
			<span class="text-yellow-500 text-sm" title="Importance: {rule.weight}/5">
				{renderStars(rule.weight)}
			</span>
		</div>
	</div>
</div>

<style>
	.rule-card {
		transition: all 0.2s ease;
	}

	.rule-card:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
	}
</style>
