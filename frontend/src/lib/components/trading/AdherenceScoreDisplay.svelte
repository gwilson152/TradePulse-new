<script lang="ts">
	import type { RuleAdherence, Rule, AdherenceScore } from '$lib/types';
	import Icon from '@iconify/svelte';

	interface Props {
		adherences: RuleAdherence[];
		rules: Rule[];
		showDetails?: boolean;
	}

	let { adherences = [], rules = [], showDetails = false }: Props = $props();

	// Calculate adherence score
	const score = $derived((): AdherenceScore => {
		if (adherences.length === 0 || rules.length === 0) {
			return {
				overall_score: 0,
				weighted_score: 0,
				color: 'red',
				phase_scores: {
					pre_trade: 0,
					during_trade: 0,
					post_trade: 0
				}
			};
		}

		// Calculate overall score (simple average)
		const totalScore = adherences.reduce((sum, a) => sum + a.score, 0);
		const overall_score = totalScore / adherences.length;

		// Calculate weighted score
		let weightedSum = 0;
		let totalWeight = 0;

		adherences.forEach((adherence) => {
			const rule = rules.find((r) => r.id === adherence.rule_id);
			if (rule) {
				weightedSum += adherence.score * rule.weight;
				totalWeight += rule.weight;
			}
		});

		const weighted_score = totalWeight > 0 ? weightedSum / totalWeight : 0;

		// Determine color
		let color: 'green' | 'yellow' | 'red' = 'red';
		if (weighted_score >= 80) color = 'green';
		else if (weighted_score >= 50) color = 'yellow';

		// Calculate phase scores
		const phaseScores = {
			PRE_TRADE: [] as number[],
			DURING_TRADE: [] as number[],
			POST_TRADE: [] as number[]
		};

		adherences.forEach((adherence) => {
			const rule = rules.find((r) => r.id === adherence.rule_id);
			if (rule) {
				phaseScores[rule.phase].push(adherence.score);
			}
		});

		const phase_scores = {
			pre_trade:
				phaseScores.PRE_TRADE.length > 0
					? phaseScores.PRE_TRADE.reduce((a, b) => a + b, 0) / phaseScores.PRE_TRADE.length
					: 0,
			during_trade:
				phaseScores.DURING_TRADE.length > 0
					? phaseScores.DURING_TRADE.reduce((a, b) => a + b, 0) / phaseScores.DURING_TRADE.length
					: 0,
			post_trade:
				phaseScores.POST_TRADE.length > 0
					? phaseScores.POST_TRADE.reduce((a, b) => a + b, 0) / phaseScores.POST_TRADE.length
					: 0
		};

		return {
			overall_score,
			weighted_score,
			color,
			phase_scores
		};
	});

	function getColorClasses(color: 'green' | 'yellow' | 'red'): string {
		if (color === 'green')
			return 'bg-profit-500 text-white border-profit-600';
		if (color === 'yellow')
			return 'bg-yellow-500 text-white border-yellow-600';
		return 'bg-loss-500 text-white border-loss-600';
	}

	function getPhaseColor(score: number): string {
		if (score >= 80) return 'text-profit-600';
		if (score >= 50) return 'text-yellow-600';
		return 'text-loss-600';
	}

	function getTrafficLight(score: number): string {
		if (score >= 80) return '🟢';
		if (score >= 50) return '🟡';
		return '🔴';
	}
</script>

<div class="adherence-score-display">
	<div class="flex items-center gap-4">
		<!-- Traffic light circle with score -->
		<div class="relative">
			<div
				class="w-20 h-20 rounded-full flex flex-col items-center justify-center border-4 {getColorClasses(
					score().color
				)}"
			>
				<div class="text-2xl font-bold">{Math.round(score().weighted_score)}</div>
				<div class="text-xs opacity-90">Score</div>
			</div>
			<div class="absolute -top-1 -right-1 text-2xl">
				{getTrafficLight(score().weighted_score)}
			</div>
		</div>

		<!-- Score details -->
		<div class="flex-1">
			<h3 class="text-lg font-semibold text-surface-900 dark:text-surface-100 mb-2">
				Rule Adherence
			</h3>
			<div class="grid grid-cols-2 gap-2 text-sm">
				<div>
					<span class="text-surface-600 dark:text-surface-400">Overall Score:</span>
					<span class="font-semibold text-surface-900 dark:text-surface-100 ml-2">
						{Math.round(score().overall_score)}%
					</span>
				</div>
				<div>
					<span class="text-surface-600 dark:text-surface-400">Weighted Score:</span>
					<span class="font-semibold text-surface-900 dark:text-surface-100 ml-2">
						{Math.round(score().weighted_score)}%
					</span>
				</div>
			</div>
			<div class="mt-2 text-xs text-surface-600 dark:text-surface-400">
				<Icon icon="mdi:information" class="inline mr-1" />
				Based on {adherences.length} rule{adherences.length !== 1 ? 's' : ''}
			</div>
		</div>
	</div>

	{#if showDetails}
		<div class="mt-4 pt-4 border-t border-surface-200 dark:border-surface-700">
			<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">
				Score by Phase
			</h4>
			<div class="grid grid-cols-3 gap-3">
				<div class="bg-surface-100 dark:bg-surface-800 rounded-lg p-3 text-center">
					<div class="text-xs text-surface-600 dark:text-surface-400 mb-1">Pre-Trade</div>
					<div class="text-xl font-bold {getPhaseColor(score().phase_scores.pre_trade)}">
						{score().phase_scores.pre_trade > 0
							? Math.round(score().phase_scores.pre_trade)
							: '-'}
					</div>
					{#if score().phase_scores.pre_trade > 0}
						<div class="text-sm mt-1">{getTrafficLight(score().phase_scores.pre_trade)}</div>
					{/if}
				</div>
				<div class="bg-surface-100 dark:bg-surface-800 rounded-lg p-3 text-center">
					<div class="text-xs text-surface-600 dark:text-surface-400 mb-1">During Trade</div>
					<div class="text-xl font-bold {getPhaseColor(score().phase_scores.during_trade)}">
						{score().phase_scores.during_trade > 0
							? Math.round(score().phase_scores.during_trade)
							: '-'}
					</div>
					{#if score().phase_scores.during_trade > 0}
						<div class="text-sm mt-1">{getTrafficLight(score().phase_scores.during_trade)}</div>
					{/if}
				</div>
				<div class="bg-surface-100 dark:bg-surface-800 rounded-lg p-3 text-center">
					<div class="text-xs text-surface-600 dark:text-surface-400 mb-1">Post-Trade</div>
					<div class="text-xl font-bold {getPhaseColor(score().phase_scores.post_trade)}">
						{score().phase_scores.post_trade > 0
							? Math.round(score().phase_scores.post_trade)
							: '-'}
					</div>
					{#if score().phase_scores.post_trade > 0}
						<div class="text-sm mt-1">{getTrafficLight(score().phase_scores.post_trade)}</div>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.adherence-score-display {
		padding: 1rem;
		background-color: rgb(var(--color-surface-50));
		border-radius: 0.5rem;
		border: 1px solid rgb(var(--color-surface-200));
	}

	:global(.dark) .adherence-score-display {
		background-color: rgb(var(--color-surface-800));
		border-color: rgb(var(--color-surface-700));
	}
</style>
