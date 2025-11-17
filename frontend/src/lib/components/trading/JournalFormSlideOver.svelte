<script lang="ts">
	import FormSlideOver from '../ui/FormSlideOver.svelte';
	import Button from '../ui/Button.svelte';
	import FileUpload from '../ui/FileUpload.svelte';
	import AudioRecorder from '../ui/AudioRecorder.svelte';
	import RuleAdherenceInput from './RuleAdherenceInput.svelte';
	import AdherenceScoreDisplay from './AdherenceScoreDisplay.svelte';
	import HelpText from '../ui/HelpText.svelte';
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

	// Tab configuration with reactive completion status
	const tabs = $derived([
		{ id: 'reflection', label: 'Reflection', icon: 'mdi:note-text', isComplete: isTabComplete('reflection') },
		{ id: 'emotions', label: 'Emotions', icon: 'mdi:emoticon-happy', isComplete: isTabComplete('emotions') },
		{ id: 'rules', label: 'Rules', icon: 'mdi:checkbox-marked-circle', isComplete: isTabComplete('rules') },
		{ id: 'media', label: 'Media', icon: 'mdi:image-multiple', isComplete: isTabComplete('media') }
	]);

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
</script>

<FormSlideOver
	{open}
	title="New Journal Entry{trade ? ` - ${trade.symbol}` : ''}"
	subtitle="Document your trade and learnings"
	size="2xl"
	{error}
	{loading}
	showHelp={true}
	helpTitle="Why Journal Your Trades?"
	helpText="Consistent journaling helps you identify patterns, learn from mistakes, and improve your trading psychology. Complete all tabs for comprehensive insights."
	{tabs}
	activeTab={currentTab}
	onTabChange={(tabId) => (currentTab = tabId as typeof currentTab)}
	onClose={handleClose}
	onSubmit={handleSubmit}
	submitText="Save Entry"
>
	{#if currentTab === 'reflection'}
		<div class="space-y-4">
			<HelpText
				type="tip"
				text="Describe what happened, your thought process, and what you learned. Focus on decisions you made and why. This tab is required."
			/>

			<div>
				<label for="reflection-notes" class="block text-sm font-semibold mb-2 text-slate-700 dark:text-slate-300">
					Reflection & Notes *
				</label>
				<textarea
					id="reflection-notes"
					bind:value={content}
					class="w-full px-4 py-3 rounded-xl border-2 transition-all duration-200
						border-slate-200 dark:border-slate-700
						focus:border-blue-500 dark:focus:border-blue-400
						focus:ring-4 focus:ring-blue-500/20
						bg-white/80 dark:bg-slate-800/80 backdrop-blur-sm
						shadow-sm hover:shadow-md focus:shadow-lg
						text-slate-900 dark:text-slate-100
						placeholder-slate-400 dark:placeholder-slate-500"
					rows="12"
					placeholder="What happened during this trade? What did you learn? What would you do differently?"
					required
				></textarea>
			</div>

			{#if trade}
				<div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl border-2 border-blue-200 dark:border-blue-800">
					<h4 class="font-semibold text-blue-900 dark:text-blue-100 mb-2 flex items-center gap-2">
						<Icon icon="mdi:chart-line" width="20" />
						Trade Summary
					</h4>
					<div class="grid grid-cols-2 gap-3 text-sm text-blue-700 dark:text-blue-300">
						<div class="flex items-center gap-2">
							<span class="text-blue-500">Symbol:</span>
							<strong>{trade.symbol}</strong>
						</div>
						<div class="flex items-center gap-2">
							<span class="text-blue-500">Type:</span>
							<strong>{trade.trade_type}</strong>
						</div>
						<div class="flex items-center gap-2">
							<span class="text-blue-500">Entry:</span>
							<strong>${trade.entry_price.toFixed(2)}</strong>
						</div>
						<div class="flex items-center gap-2">
							<span class="text-blue-500">P&L:</span>
							<strong class="{(trade.pnl ?? 0) >= 0 ? 'text-emerald-600' : 'text-red-600'}">
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

			<p class="text-sm text-slate-600 dark:text-slate-400">
				Rate your emotional state during this trade (1-10)
			</p>

			<!-- Confidence -->
			<div>
				<div class="flex justify-between items-center mb-3">
					<label class="text-sm font-semibold text-slate-700 dark:text-slate-300 flex items-center gap-2">
						<Icon icon="mdi:account-check" width="20" class="text-blue-500" />
						Confidence
					</label>
					<span class="text-2xl font-bold text-blue-600">{emotionalState.confidence}</span>
				</div>
				<input
					type="range"
					min="1"
					max="10"
					bind:value={emotionalState.confidence}
					class="w-full h-3 bg-slate-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer
						accent-blue-600"
				/>
				<div class="flex justify-between text-xs text-slate-500 mt-2">
					<span>Low</span>
					<span>High</span>
				</div>
			</div>

			<!-- Stress -->
			<div>
				<div class="flex justify-between items-center mb-3">
					<label class="text-sm font-semibold text-slate-700 dark:text-slate-300 flex items-center gap-2">
						<Icon icon="mdi:alert" width="20" class="text-amber-500" />
						Stress Level
					</label>
					<span class="text-2xl font-bold text-amber-600">{emotionalState.stress}</span>
				</div>
				<input
					type="range"
					min="1"
					max="10"
					bind:value={emotionalState.stress}
					class="w-full h-3 bg-slate-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer
						accent-amber-600"
				/>
				<div class="flex justify-between text-xs text-slate-500 mt-2">
					<span>Calm</span>
					<span>High Stress</span>
				</div>
			</div>

			<!-- Discipline -->
			<div>
				<div class="flex justify-between items-center mb-3">
					<label class="text-sm font-semibold text-slate-700 dark:text-slate-300 flex items-center gap-2">
						<Icon icon="mdi:shield-check" width="20" class="text-emerald-500" />
						Discipline
					</label>
					<span class="text-2xl font-bold text-emerald-600">{emotionalState.discipline}</span>
				</div>
				<input
					type="range"
					min="1"
					max="10"
					bind:value={emotionalState.discipline}
					class="w-full h-3 bg-slate-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer
						accent-emerald-600"
				/>
				<div class="flex justify-between text-xs text-slate-500 mt-2">
					<span>Poor</span>
					<span>Excellent</span>
				</div>
			</div>

			<!-- Emotional notes -->
			<div>
				<label for="emotional-notes" class="block text-sm font-semibold mb-2 text-slate-700 dark:text-slate-300">
					Emotional Notes
				</label>
				<textarea
					id="emotional-notes"
					bind:value={emotionalState.notes}
					class="w-full px-4 py-3 rounded-xl border-2 transition-all duration-200
						border-slate-200 dark:border-slate-700
						focus:border-blue-500 dark:focus:border-blue-400
						focus:ring-4 focus:ring-blue-500/20
						bg-white/80 dark:bg-slate-800/80 backdrop-blur-sm
						shadow-sm hover:shadow-md focus:shadow-lg
						text-slate-900 dark:text-slate-100
						placeholder-slate-400 dark:placeholder-slate-500"
					rows="4"
					placeholder="How were you feeling? Any specific emotional reactions or concerns?"
				></textarea>
			</div>
		</div>
	{:else if currentTab === 'rules'}
		<div class="space-y-4">
			{#if rules.length === 0}
				<div class="text-center py-12 text-slate-600 dark:text-slate-400">
					<Icon icon="mdi:checkbox-blank-circle-outline" width="64" class="mx-auto mb-4 text-slate-400" />
					<p class="font-semibold mb-2 text-lg">No Rules Configured</p>
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
				<h4 class="text-sm font-semibold text-slate-700 dark:text-slate-300 mb-3 flex items-center gap-2">
					<Icon icon="mdi:image-multiple" width="20" />
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
				<p class="text-xs text-slate-500 dark:text-slate-400 mt-2">
					Upload chart screenshots, entry/exit points, or trade setup images (max 10MB each)
				</p>
			</div>

			<!-- Voice Notes -->
			<div>
				<h4 class="text-sm font-semibold text-slate-700 dark:text-slate-300 mb-3 flex items-center gap-2">
					<Icon icon="mdi:microphone" width="20" />
					Voice Notes
				</h4>
				<AudioRecorder onRecordingComplete={handleVoiceNoteComplete} maxDuration={300} />
				<p class="text-xs text-slate-500 dark:text-slate-400 mt-2">
					Record quick thoughts or observations (max 5 minutes per recording)
				</p>

				{#if voiceNoteBlobs.length > 0}
					<div class="mt-4 space-y-2">
						<p class="text-sm font-medium text-slate-600 dark:text-slate-400">
							{voiceNoteBlobs.length} voice note{voiceNoteBlobs.length !== 1 ? 's' : ''} recorded
						</p>
						{#each voiceNoteBlobs as blob, index}
							<div class="flex items-center gap-3 p-3 bg-slate-100 dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700">
								<Icon icon="mdi:microphone" width="24" class="text-blue-500" />
								<span class="text-sm flex-1 font-medium">Voice Note {index + 1}</span>
								<button
									type="button"
									onclick={() => removeVoiceNote(index)}
									class="p-2 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-lg transition-colors"
								>
									<Icon icon="mdi:close" width="20" class="text-red-600" />
								</button>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	{/if}
</FormSlideOver>
