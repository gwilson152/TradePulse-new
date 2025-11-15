<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import FileUpload from '$lib/components/ui/FileUpload.svelte';
	import AudioRecorder from '$lib/components/ui/AudioRecorder.svelte';
	import RuleAdherenceInput from './RuleAdherenceInput.svelte';
	import AdherenceScoreDisplay from './AdherenceScoreDisplay.svelte';
	import HelpText from '$lib/components/ui/HelpText.svelte';
	import type { JournalEntry, EmotionalState, Trade, Rule, RuleAdherence } from '$lib/types';
	import Icon from '@iconify/svelte';

	interface Props {
		open: boolean;
		trade?: Trade | null;
		rules: Rule[];
		onClose: () => void;
		onSubmit: (data: Partial<JournalEntry>, screenshots: File[], voiceNotes: Blob[]) => Promise<void>;
	}

	let { open = false, trade = null, rules = [], onClose, onSubmit }: Props = $props();

	let currentTab = $state<'reflection' | 'emotions' | 'rules' | 'media'>('reflection');
	let loading = $state(false);
	let error = $state('');

	// Form data
	let content = $state('');
	let emotionalState = $state<EmotionalState>({
		confidence: 5,
		stress: 5,
		discipline: 5,
		notes: ''
	});
	let ruleAdherences = $state<Map<string, RuleAdherence>>(new Map());
	let screenshots = $state<File[]>([]);
	let voiceNoteBlobs = $state<Blob[]>([]);

	// Calculate overall adherence score
	const adherenceScore = $derived(() => {
		const adherenceArray = Array.from(ruleAdherences.values());
		if (adherenceArray.length === 0) return null;

		let weightedSum = 0;
		let totalWeight = 0;

		adherenceArray.forEach((adherence) => {
			const rule = rules.find((r) => r.id === adherence.rule_id);
			if (rule) {
				weightedSum += adherence.score * rule.weight;
				totalWeight += rule.weight;
			}
		});

		return totalWeight > 0 ? weightedSum / totalWeight : 0;
	});

	function handleRuleAdherence(adherence: RuleAdherence) {
		ruleAdherences.set(adherence.rule_id, adherence);
	}

	function handleScreenshotsChange(files: File[]) {
		screenshots = files;
	}

	function handleVoiceNoteComplete(blob: Blob) {
		voiceNoteBlobs = [...voiceNoteBlobs, blob];
	}

	function removeVoiceNote(index: number) {
		voiceNoteBlobs = voiceNoteBlobs.filter((_, i) => i !== index);
	}

	async function handleSubmit() {
		error = '';

		// Validation
		if (!content.trim()) {
			error = 'Please add some reflection notes';
			currentTab = 'reflection';
			return;
		}

		// Check for required rule adherence explanations
		const adherenceArray = Array.from(ruleAdherences.values());
		const missingExplanations = adherenceArray.filter(
			(a) => a.score < 100 && !a.notes.trim()
		);

		if (missingExplanations.length > 0) {
			error = 'Please provide explanations for all rules with scores below 100%';
			currentTab = 'rules';
			return;
		}

		loading = true;

		try {
			const journalData: Partial<JournalEntry> = {
				trade_id: trade?.id || null,
				entry_date: new Date().toISOString(),
				content,
				emotional_state: emotionalState,
				rule_adherence: adherenceArray,
				adherence_score: adherenceScore()
			};

			await onSubmit(journalData, screenshots, voiceNoteBlobs);
			resetForm();
			onClose();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create journal entry';
		} finally {
			loading = false;
		}
	}

	function resetForm() {
		content = '';
		emotionalState = {
			confidence: 5,
			stress: 5,
			discipline: 5,
			notes: ''
		};
		ruleAdherences.clear();
		screenshots = [];
		voiceNoteBlobs = [];
		currentTab = 'reflection';
		error = '';
	}

	function handleClose() {
		resetForm();
		onClose();
	}

	function getTabIcon(tab: string): string {
		switch (tab) {
			case 'reflection':
				return 'mdi:note-text';
			case 'emotions':
				return 'mdi:emoticon-happy';
			case 'rules':
				return 'mdi:checkbox-marked-circle';
			case 'media':
				return 'mdi:image-multiple';
			default:
				return 'mdi:circle';
		}
	}

	function isTabComplete(tab: string): boolean {
		switch (tab) {
			case 'reflection':
				return content.trim().length > 0;
			case 'emotions':
				return true; // Always complete with defaults
			case 'rules':
				return (
					ruleAdherences.size === rules.length &&
					Array.from(ruleAdherences.values()).every((a) => a.score === 100 || a.notes.trim().length > 0)
				);
			case 'media':
				return true; // Optional
			default:
				return false;
		}
	}
</script>

<Modal
	{open}
	onClose={handleClose}
	title="New Journal Entry{trade ? ` - ${trade.symbol}` : ''}"
	size="large"
>
	<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
		{#if error}
			<div class="bg-error-50 dark:bg-error-900/20 border border-error-200 dark:border-error-800 text-error-700 dark:text-error-300 px-4 py-3 rounded">
				<Icon icon="mdi:alert-circle" class="inline mr-2" />
				{error}
			</div>
		{/if}

		<!-- Help Section -->
		<HelpText
			type="info"
			title="Why Journal Your Trades?"
			text="Consistent journaling helps you identify patterns, learn from mistakes, and improve your trading psychology. Complete all tabs for comprehensive insights."
			collapsible={true}
		/>

		<!-- Tabs -->
		<div class="border-b border-surface-200 dark:border-surface-700">
			<div class="flex gap-1">
				{#each ['reflection', 'emotions', 'rules', 'media'] as tab}
					<button
						type="button"
						onclick={() => (currentTab = tab as typeof currentTab)}
						class="flex items-center gap-2 px-4 py-2 border-b-2 transition-colors {currentTab ===
						tab
							? 'border-primary-500 text-primary-600 dark:text-primary-400'
							: 'border-transparent text-surface-600 dark:text-surface-400 hover:text-surface-800 dark:hover:text-surface-200'}"
					>
						<Icon icon={getTabIcon(tab)} />
						<span class="capitalize">{tab}</span>
						{#if isTabComplete(tab)}
							<Icon icon="mdi:check-circle" class="text-success-500 text-sm" />
						{/if}
					</button>
				{/each}
			</div>
		</div>

		<!-- Tab content -->
		<div class="min-h-[400px]">
			{#if currentTab === 'reflection'}
				<div class="space-y-4">
					<HelpText
						type="tip"
						text="Describe what happened, your thought process, and what you learned. Focus on decisions you made and why. This tab is required."
					/>

					<div>
						<label class="block text-sm font-medium mb-2 text-surface-700 dark:text-surface-300">
							Reflection & Notes *
						</label>
						<textarea
							bind:value={content}
							class="w-full px-3 py-2 border border-surface-300 dark:border-surface-600 rounded-lg bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100"
							rows="12"
							placeholder="What happened during this trade? What did you learn? What would you do differently?"
							required
						></textarea>
					</div>

					{#if trade}
						<div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
							<h4 class="font-semibold text-blue-900 dark:text-blue-100 mb-2">
								Trade Summary
							</h4>
							<div class="grid grid-cols-2 gap-2 text-sm text-blue-700 dark:text-blue-300">
								<div>Symbol: <strong>{trade.symbol}</strong></div>
								<div>Type: <strong>{trade.trade_type}</strong></div>
								<div>Entry: <strong>${trade.entry_price.toFixed(2)}</strong></div>
								<div>
									P&L: <strong class="{(trade.pnl ?? 0) >= 0 ? 'text-profit-600' : 'text-loss-600'}">
										{(trade.pnl ?? 0) >= 0 ? '+' : ''}${(trade.pnl ?? 0).toFixed(2)}
									</strong>
								</div>
							</div>
						</div>
					{/if}
				</div>
			{:else if currentTab === 'emotions'}
				<div class="space-y-6">
					<HelpText
						type="info"
						title="Track Your Trading Psychology"
						text="Emotional awareness is key to consistent trading. Rate your confidence, stress, and discipline honestly. Over time, you'll see patterns between your emotions and trading performance."
						collapsible={true}
					/>

					<p class="text-sm text-surface-600 dark:text-surface-400">
						Rate your emotional state during this trade (1-10)
					</p>

					<!-- Confidence -->
					<div>
						<div class="flex justify-between items-center mb-2">
							<label class="text-sm font-medium text-surface-700 dark:text-surface-300">
								<Icon icon="mdi:account-check" class="inline mr-1" />
								Confidence
							</label>
							<span class="text-lg font-bold text-primary-600">{emotionalState.confidence}</span>
						</div>
						<input
							type="range"
							min="1"
							max="10"
							bind:value={emotionalState.confidence}
							class="w-full h-2 bg-surface-200 dark:bg-surface-700 rounded-lg appearance-none cursor-pointer"
						/>
						<div class="flex justify-between text-xs text-surface-500 mt-1">
							<span>Low</span>
							<span>High</span>
						</div>
					</div>

					<!-- Stress -->
					<div>
						<div class="flex justify-between items-center mb-2">
							<label class="text-sm font-medium text-surface-700 dark:text-surface-300">
								<Icon icon="mdi:alert" class="inline mr-1" />
								Stress Level
							</label>
							<span class="text-lg font-bold text-warning-600">{emotionalState.stress}</span>
						</div>
						<input
							type="range"
							min="1"
							max="10"
							bind:value={emotionalState.stress}
							class="w-full h-2 bg-surface-200 dark:bg-surface-700 rounded-lg appearance-none cursor-pointer"
						/>
						<div class="flex justify-between text-xs text-surface-500 mt-1">
							<span>Calm</span>
							<span>High Stress</span>
						</div>
					</div>

					<!-- Discipline -->
					<div>
						<div class="flex justify-between items-center mb-2">
							<label class="text-sm font-medium text-surface-700 dark:text-surface-300">
								<Icon icon="mdi:shield-check" class="inline mr-1" />
								Discipline
							</label>
							<span class="text-lg font-bold text-success-600">{emotionalState.discipline}</span>
						</div>
						<input
							type="range"
							min="1"
							max="10"
							bind:value={emotionalState.discipline}
							class="w-full h-2 bg-surface-200 dark:bg-surface-700 rounded-lg appearance-none cursor-pointer"
						/>
						<div class="flex justify-between text-xs text-surface-500 mt-1">
							<span>Poor</span>
							<span>Excellent</span>
						</div>
					</div>

					<!-- Emotional notes -->
					<div>
						<label class="block text-sm font-medium mb-2 text-surface-700 dark:text-surface-300">
							Emotional Notes
						</label>
						<textarea
							bind:value={emotionalState.notes}
							class="w-full px-3 py-2 border border-surface-300 dark:border-surface-600 rounded-lg bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100"
							rows="4"
							placeholder="How were you feeling? Any specific emotional reactions or concerns?"
						></textarea>
					</div>
				</div>
			{:else if currentTab === 'rules'}
				<div class="space-y-4">
					{#if rules.length === 0}
						<div class="text-center py-12 text-surface-600 dark:text-surface-400">
							<Icon icon="mdi:checkbox-blank-circle-outline" class="text-5xl mb-3 mx-auto" />
							<p class="font-medium mb-2">No Rules Configured</p>
							<p class="text-sm">Create rule sets in Settings to track adherence</p>
						</div>
					{:else}
						<HelpText
							type="warning"
							title="Rule Adherence Scoring"
							text="Rate how well you followed each rule (0-100%). If you scored below 100%, EXPLAIN WHY in the notes field. This helps you identify which rules you struggle with and why."
							collapsible={true}
						/>

						<div class="mb-4">
							<AdherenceScoreDisplay adherences={Array.from(ruleAdherences.values())} {rules} showDetails={true} />
						</div>

						<div class="space-y-3">
							{#each rules as rule}
								<RuleAdherenceInput
									{rule}
									value={ruleAdherences.get(rule.id)}
									onChange={handleRuleAdherence}
								/>
							{/each}
						</div>
					{/if}
				</div>
			{:else if currentTab === 'media'}
				<div class="space-y-6">
					<HelpText
						type="tip"
						title="Visual Documentation (Optional)"
						text="Screenshots of your charts, setups, and indicators help you review trades later. Voice notes let you quickly capture thoughts when typing isn't convenient. Both are optional but valuable."
						collapsible={true}
					/>

					<!-- Screenshots -->
					<div>
						<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">
							Screenshots
						</h4>
						<FileUpload
							accept="image/*"
							multiple={true}
							maxSize={10}
							label=""
							value={screenshots}
							onChange={handleScreenshotsChange}
						/>
						<p class="text-xs text-surface-500 dark:text-surface-400 mt-2">
							Upload chart screenshots, entry/exit points, or trade setup images
						</p>
					</div>

					<!-- Voice Notes -->
					<div>
						<h4 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">
							Voice Notes
						</h4>
						<AudioRecorder onRecordingComplete={handleVoiceNoteComplete} maxDuration={300} />
						<p class="text-xs text-surface-500 dark:text-surface-400 mt-2">
							Record quick thoughts or observations (max 5 minutes per recording)
						</p>

						{#if voiceNoteBlobs.length > 0}
							<div class="mt-4 space-y-2">
								<p class="text-sm text-surface-600 dark:text-surface-400">
									{voiceNoteBlobs.length} voice note{voiceNoteBlobs.length !== 1 ? 's' : ''} recorded
								</p>
								{#each voiceNoteBlobs as blob, index}
									<div class="flex items-center gap-2 p-2 bg-surface-100 dark:bg-surface-800 rounded">
										<Icon icon="mdi:microphone" class="text-primary-500" />
										<span class="text-sm flex-1">Voice Note {index + 1}</span>
										<button
											type="button"
											onclick={() => removeVoiceNote(index)}
											class="p-1 hover:bg-loss-100 dark:hover:bg-loss-900/30 rounded"
										>
											<Icon icon="mdi:close" class="text-loss-600" />
										</button>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			{/if}
		</div>

		<!-- Footer buttons -->
		<div class="flex gap-2 pt-4 border-t border-surface-200 dark:border-surface-700">
			<Button type="button" variant="ghost" onclick={handleClose} class="flex-1">
				Cancel
			</Button>
			<Button type="submit" color="primary" disabled={loading} class="flex-1">
				{#if loading}
					<Icon icon="mdi:loading" class="animate-spin mr-2" />
					Saving...
				{:else}
					<Icon icon="mdi:content-save" class="mr-2" />
					Save Entry
				{/if}
			</Button>
		</div>
	</form>
</Modal>
