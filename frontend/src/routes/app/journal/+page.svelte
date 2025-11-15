<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import ImageGallery from '$lib/components/ui/ImageGallery.svelte';
	import AudioPlayer from '$lib/components/ui/AudioPlayer.svelte';
	import JournalFormSlideOver from '$lib/components/trading/JournalFormSlideOver.svelte';
	import AdherenceScoreDisplay from '$lib/components/trading/AdherenceScoreDisplay.svelte';
	import Icon from '@iconify/svelte';
	import type { JournalEntry, Rule } from '$lib/types';
	import { onMount } from 'svelte';
	import { apiClient } from '$lib/api/client';

	let entries = $state<JournalEntry[]>([]);
	let rules = $state<Rule[]>([]);
	let showAddModal = $state(false);
	let loading = $state(true);
	let expandedEntry = $state<string | null>(null);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			// Load journal entries and rules
			// TODO: Implement actual API calls when backend is ready
			entries = [];
			rules = [];
		} catch (err) {
			console.error('Failed to load journal data:', err);
		} finally {
			loading = false;
		}
	}

	async function handleCreateEntry(
		data: Partial<JournalEntry>,
		screenshots: File[],
		voiceNotes: Blob[]
	) {
		try {
			// TODO: Upload files and create entry via API
			console.log('Creating entry:', data, screenshots, voiceNotes);
			await loadData();
		} catch (err) {
			throw new Error('Failed to create journal entry');
		}
	}

	function formatDate(date: string) {
		return new Date(date).toLocaleDateString('en-US', {
			month: 'long',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function toggleExpand(entryId: string) {
		expandedEntry = expandedEntry === entryId ? null : entryId;
	}

	function getAdherenceColor(score: number | undefined): 'success' | 'warning' | 'error' {
		if (!score) return 'error';
		if (score >= 80) return 'success';
		if (score >= 50) return 'warning';
		return 'error';
	}
</script>

<svelte:head>
	<title>Journal - TradePulse</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold mb-2">Trading Journal</h1>
			<p class="text-surface-600 dark:text-surface-400">
				Track your emotional state and trade reflections
			</p>
		</div>
		<Button color="primary" onclick={() => (showAddModal = true)}>
			<Icon icon="mdi:plus" width="20" class="mr-2" />
			New Entry
		</Button>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<Icon icon="mdi:loading" class="animate-spin text-4xl text-primary-500" />
		</div>
	{:else if entries.length === 0}
		<Card>
			<div class="text-center py-12">
				<Icon icon="mdi:notebook" width="64" class="mx-auto mb-4 text-surface-400" />
				<h3 class="text-xl font-semibold mb-2">No journal entries yet</h3>
				<p class="text-surface-600 dark:text-surface-400 mb-6">
					Start documenting your trading journey, emotional state, and rule adherence
				</p>
				<Button color="primary" onclick={() => (showAddModal = true)}>
					<Icon icon="mdi:pencil" width="20" class="mr-2" />
					Create First Entry
				</Button>
			</div>
		</Card>
	{:else}
		<div class="space-y-4">
			{#each entries as entry}
				<Card hover={true}>
					<div class="space-y-4">
						<!-- Header -->
						<div class="flex items-start justify-between">
							<div class="flex-1">
								<div class="flex items-center gap-3 mb-2">
									<h3 class="text-lg font-semibold">{formatDate(entry.entry_date)}</h3>
									{#if entry.adherence_score !== undefined}
										<Badge color={getAdherenceColor(entry.adherence_score)} variant="soft" size="sm">
											{Math.round(entry.adherence_score)}% Adherence
										</Badge>
									{/if}
								</div>
								{#if entry.trade_id}
									<p class="text-sm text-surface-600 dark:text-surface-400">
										<Icon icon="mdi:chart-line" class="inline mr-1" />
										Linked to trade
									</p>
								{/if}
							</div>
							<div class="flex gap-2">
								<button
									onclick={() => toggleExpand(entry.id)}
									class="p-2 hover:bg-surface-100 dark:hover:bg-surface-700 rounded transition-colors"
									title={expandedEntry === entry.id ? 'Collapse' : 'Expand'}
								>
									<Icon
										icon={expandedEntry === entry.id ? 'mdi:chevron-up' : 'mdi:chevron-down'}
										width="20"
									/>
								</button>
								<button class="p-2 hover:bg-surface-100 dark:hover:bg-surface-700 rounded transition-colors">
									<Icon icon="mdi:pencil" width="20" />
								</button>
								<button class="p-2 hover:bg-error-100 dark:hover:bg-error-900/30 rounded transition-colors text-error-600">
									<Icon icon="mdi:delete" width="20" />
								</button>
							</div>
						</div>

						<!-- Emotional State -->
						{#if entry.emotional_state}
							<div class="grid grid-cols-3 gap-4">
								<div>
									<p class="text-sm text-surface-600 dark:text-surface-400 mb-1 flex items-center gap-1">
										<Icon icon="mdi:account-check" class="text-primary-500" />
										Confidence
									</p>
									<div class="flex items-center gap-2">
										<div class="flex-1 h-2 bg-surface-200 dark:bg-surface-700 rounded-full overflow-hidden">
											<div
												class="h-full bg-primary-600 transition-all"
												style="width: {entry.emotional_state.confidence * 10}%"
											></div>
										</div>
										<span class="text-sm font-medium">{entry.emotional_state.confidence}/10</span>
									</div>
								</div>
								<div>
									<p class="text-sm text-surface-600 dark:text-surface-400 mb-1 flex items-center gap-1">
										<Icon icon="mdi:alert" class="text-warning-500" />
										Stress
									</p>
									<div class="flex items-center gap-2">
										<div class="flex-1 h-2 bg-surface-200 dark:bg-surface-700 rounded-full overflow-hidden">
											<div
												class="h-full bg-warning-600 transition-all"
												style="width: {entry.emotional_state.stress * 10}%"
											></div>
										</div>
										<span class="text-sm font-medium">{entry.emotional_state.stress}/10</span>
									</div>
								</div>
								<div>
									<p class="text-sm text-surface-600 dark:text-surface-400 mb-1 flex items-center gap-1">
										<Icon icon="mdi:shield-check" class="text-success-500" />
										Discipline
									</p>
									<div class="flex items-center gap-2">
										<div class="flex-1 h-2 bg-surface-200 dark:bg-surface-700 rounded-full overflow-hidden">
											<div
												class="h-full bg-success-600 transition-all"
												style="width: {entry.emotional_state.discipline * 10}%"
											></div>
										</div>
										<span class="text-sm font-medium">{entry.emotional_state.discipline}/10</span>
									</div>
								</div>
							</div>
						{/if}

						<!-- Content Preview -->
						<p class="text-surface-700 dark:text-surface-300 line-clamp-3">
							{entry.content}
						</p>

						<!-- Expanded Content -->
						{#if expandedEntry === entry.id}
							<div class="pt-4 border-t border-surface-200 dark:border-surface-700 space-y-4">
								<!-- Full Content -->
								<div>
									<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-2">
										Full Reflection
									</h4>
									<p class="text-surface-700 dark:text-surface-300 whitespace-pre-wrap">
										{entry.content}
									</p>
								</div>

								<!-- Emotional Notes -->
								{#if entry.emotional_state?.notes}
									<div>
										<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-2">
											Emotional Notes
										</h4>
										<p class="text-surface-700 dark:text-surface-300">
											{entry.emotional_state.notes}
										</p>
									</div>
								{/if}

								<!-- Rule Adherence -->
								{#if entry.rule_adherence && entry.rule_adherence.length > 0}
									<div>
										<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">
											Rule Adherence
										</h4>
										<AdherenceScoreDisplay adherences={entry.rule_adherence} {rules} showDetails={true} />
									</div>
								{/if}

								<!-- Screenshots -->
								{#if entry.screenshots && entry.screenshots.length > 0}
									<div>
										<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">
											Screenshots
										</h4>
										<ImageGallery images={entry.screenshots} alt="Trade screenshot" />
									</div>
								{/if}

								<!-- Voice Notes -->
								{#if entry.voice_notes && entry.voice_notes.length > 0}
									<div>
										<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">
											Voice Notes
										</h4>
										<div class="space-y-2">
											{#each entry.voice_notes as voiceNote, index}
												<AudioPlayer src={voiceNote} label="Voice Note {index + 1}" />
											{/each}
										</div>
									</div>
								{/if}
							</div>
						{/if}
					</div>
				</Card>
			{/each}
		</div>
	{/if}
</div>

<!-- Journal Entry Slide-over -->
<JournalFormSlideOver
	open={showAddModal}
	{rules}
	onClose={() => (showAddModal = false)}
	onSubmit={handleCreateEntry}
/>
