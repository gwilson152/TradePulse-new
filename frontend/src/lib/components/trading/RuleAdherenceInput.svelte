<script lang="ts">
	import type { Rule, RuleAdherence } from '$lib/types';
	import Icon from '@iconify/svelte';

	interface Props {
		rule: Rule;
		value?: RuleAdherence;
		onChange: (adherence: RuleAdherence) => void;
	}

	let { rule, value, onChange }: Props = $props();

	// Initialize with existing value or defaults
	let score = $state(value?.score ?? 100);
	let notes = $state(value?.notes ?? '');
	let expanded = $state(false);

	const scoreOptions = [
		{ value: 100, label: 'Perfect', emoji: '✓', color: 'green' },
		{ value: 75, label: 'Good', emoji: '○', color: 'green' },
		{ value: 50, label: 'Partial', emoji: '△', color: 'yellow' },
		{ value: 25, label: 'Poor', emoji: '◇', color: 'red' },
		{ value: 0, label: 'Failed', emoji: '✗', color: 'red' }
	] as const;

	const currentOption = $derived(scoreOptions.find((opt) => opt.value === score) || scoreOptions[0]);

	function getTrafficLightColor(): string {
		if (score >= 80) return 'green';
		if (score >= 50) return 'yellow';
		return 'red';
	}

	function getColorClasses(color: string, isSelected: boolean = false): string {
		const baseClasses = 'border-2 transition-all duration-200';

		if (isSelected) {
			if (color === 'green')
				return `${baseClasses} bg-profit-100 border-profit-500 dark:bg-profit-900/30 dark:border-profit-400 ring-2 ring-profit-300 dark:ring-profit-600`;
			if (color === 'yellow')
				return `${baseClasses} bg-yellow-100 border-yellow-500 dark:bg-yellow-900/30 dark:border-yellow-400 ring-2 ring-yellow-300 dark:ring-yellow-600`;
			if (color === 'red')
				return `${baseClasses} bg-loss-100 border-loss-500 dark:bg-loss-900/30 dark:border-loss-400 ring-2 ring-loss-300 dark:ring-loss-600`;
		}

		if (color === 'green')
			return `${baseClasses} border-surface-300 dark:border-surface-600 hover:border-profit-400 dark:hover:border-profit-500`;
		if (color === 'yellow')
			return `${baseClasses} border-surface-300 dark:border-surface-600 hover:border-yellow-400 dark:hover:border-yellow-500`;
		if (color === 'red')
			return `${baseClasses} border-surface-300 dark:border-surface-600 hover:border-loss-400 dark:hover:border-loss-500`;

		return baseClasses;
	}

	function selectScore(value: number) {
		score = value;
		updateAdherence();
		// Auto-expand notes if score < 100
		if (value < 100) {
			expanded = true;
		}
	}

	function updateAdherence() {
		const adherence: RuleAdherence = {
			rule_id: rule.id,
			rule_title: rule.title,
			score,
			notes,
			timestamp: new Date().toISOString()
		};
		onChange(adherence);
	}

	// Watch notes for changes
	$effect(() => {
		if (notes !== (value?.notes ?? '')) {
			updateAdherence();
		}
	});
</script>

<div class="rule-adherence-input bg-surface-50 dark:bg-surface-800 rounded-lg border border-surface-200 dark:border-surface-700 p-4">
	<!-- Rule header -->
	<div class="flex items-start gap-3 mb-4">
		<!-- Traffic light indicator -->
		<div class="mt-1">
			<div
				class="w-8 h-8 rounded-full flex items-center justify-center {getTrafficLightColor() ===
				'green'
					? 'bg-profit-500 text-white'
					: getTrafficLightColor() === 'yellow'
						? 'bg-yellow-500 text-white'
						: 'bg-loss-500 text-white'}"
			>
				<span class="text-lg font-bold">{currentOption.emoji}</span>
			</div>
		</div>

		<div class="flex-1">
			<h4 class="font-semibold text-sm text-surface-900 dark:text-surface-100 mb-1">
				{rule.title}
			</h4>
			<p class="text-xs text-surface-600 dark:text-surface-400">
				{rule.description}
			</p>
			<div class="flex items-center gap-2 mt-2">
				<span class="text-xs text-surface-500">Weight:</span>
				<span class="text-yellow-500 text-xs">{'★'.repeat(rule.weight)}{'☆'.repeat(5 - rule.weight)}</span>
			</div>
		</div>
	</div>

	<!-- Score selection -->
	<div class="space-y-2 mb-3">
		<div class="block text-xs font-semibold text-surface-700 dark:text-surface-300 uppercase tracking-wide">
			Adherence Score
		</div>
		<div class="grid grid-cols-5 gap-2">
			{#each scoreOptions as option}
				<button
					type="button"
					onclick={() => selectScore(option.value)}
					class="p-3 rounded-lg text-center cursor-pointer {getColorClasses(
						option.color,
						score === option.value
					)}"
				>
					<div class="text-2xl mb-1">{option.emoji}</div>
					<div class="text-xs font-medium text-surface-900 dark:text-surface-100">
						{option.label}
					</div>
					<div class="text-xs text-surface-600 dark:text-surface-400">{option.value}%</div>
				</button>
			{/each}
		</div>
	</div>

	<!-- Notes (auto-expand for non-perfect scores) -->
	{#if expanded || score < 100}
		<div class="mt-3 pt-3 border-t border-surface-200 dark:border-surface-700">
			<label for="rule-notes-{rule.id}" class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-2 uppercase tracking-wide">
				Explanation {score < 100 ? '(Required)' : '(Optional)'}
			</label>
			<textarea
				id="rule-notes-{rule.id}"
				bind:value={notes}
				onchange={updateAdherence}
				class="w-full px-3 py-2 border border-surface-300 dark:border-surface-600 rounded-lg bg-white dark:bg-surface-900 text-surface-900 dark:text-surface-100 text-sm"
				rows="3"
				placeholder={score < 100
					? 'Please explain what happened and why you scored this way...'
					: 'Optional notes about adherence...'}
				required={score < 100}
			></textarea>
			{#if score < 100 && !notes.trim()}
				<p class="text-xs text-loss-600 dark:text-loss-400 mt-1">
					<Icon icon="mdi:alert" class="inline mr-1" />
					Please provide an explanation for scores below 100%
				</p>
			{/if}
		</div>
	{:else}
		<button
			type="button"
			onclick={() => (expanded = true)}
			class="text-xs text-surface-600 dark:text-surface-400 hover:text-surface-800 dark:hover:text-surface-200 flex items-center gap-1 mt-2"
		>
			<Icon icon="mdi:note-plus" />
			Add notes
		</button>
	{/if}
</div>

<style>
	.rule-adherence-input {
		transition: all 0.2s ease;
	}
</style>
