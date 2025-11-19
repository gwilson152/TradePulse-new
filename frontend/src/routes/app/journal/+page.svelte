<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import ImageGallery from '$lib/components/ui/ImageGallery.svelte';
	import AudioPlayer from '$lib/components/ui/AudioPlayer.svelte';
	import JournalFormSlideOver from '$lib/components/trading/JournalFormSlideOver.svelte';
	import AdherenceScoreDisplay from '$lib/components/trading/AdherenceScoreDisplay.svelte';
	import ChartPreview from '$lib/components/charts/ChartPreview.svelte';
	import Icon from '@iconify/svelte';
	import type { JournalEntry, Rule } from '$lib/types';
	import { onMount } from 'svelte';
	import { apiClient } from '$lib/api/client';

	let entries = $state<JournalEntry[]>([]);
	let rules = $state<Rule[]>([]);
	let showAddModal = $state(false);
	let loading = $state(true);
	let selectedEntry = $state<JournalEntry | null>(null);
	let selectedTrade = $state<any | null>(null);
	let editingEntry = $state<JournalEntry | null>(null);
	let loadingTrade = $state(false);
	let isListCollapsed = $state(false);
	let entryTrades = $state<Map<string, any>>(new Map());

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			const result = await apiClient.getJournalEntries();
			entries = result.entries || [];

			// Load trade data for all entries with trade_id
			const tradePromises = entries
				.filter(entry => entry.trade_id)
				.map(async (entry) => {
					try {
						const trade = await apiClient.getTrade(entry.trade_id!);
						return { entryId: entry.id, trade };
					} catch (err) {
						console.error(`Failed to load trade for entry ${entry.id}:`, err);
						return null;
					}
				});

			const tradeResults = await Promise.all(tradePromises);
			const newEntryTrades = new Map();
			tradeResults.forEach(result => {
				if (result) {
					newEntryTrades.set(result.entryId, result.trade);
				}
			});
			entryTrades = newEntryTrades;

			if (entries.length > 0 && !selectedEntry) {
				await selectEntry(entries[0]);
			}

			const ruleSets = await apiClient.getRuleSets();
			rules = ruleSets.flatMap(rs => rs.rules || []);
		} catch (err) {
			console.error('Failed to load journal data:', err);
			entries = [];
			rules = [];
		} finally {
			loading = false;
		}
	}

	async function handleCreateEntry(data: Partial<JournalEntry>, screenshots: File[], voiceNotes: Blob[]) {
		try {
			console.log('Creating entry:', data, screenshots, voiceNotes);
			await loadData();
		} catch (err) {
			throw new Error('Failed to create journal entry');
		}
	}

	function formatDate(date: string | null | undefined) {
		if (!date) return 'Unknown Date';
		try {
			return new Date(date).toLocaleDateString('en-US', {
				month: 'short',
				day: 'numeric',
				year: 'numeric'
			});
		} catch {
			return 'Invalid Date';
		}
	}

	function formatDateTime(date: string | null | undefined) {
		if (!date) return 'Unknown Date';
		try {
			return new Date(date).toLocaleDateString('en-US', {
				month: 'long',
				day: 'numeric',
				year: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		} catch {
			return 'Invalid Date';
		}
	}

	function formatCurrency(value: number | null | undefined): string {
		if (value === null || value === undefined) return '$0.00';
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 2,
			maximumFractionDigits: 2
		}).format(value);
	}

	async function selectEntry(entry: JournalEntry) {
		selectedEntry = entry;
		selectedTrade = null;

		// Load linked trade if available
		if (entry.trade_id) {
			try {
				loadingTrade = true;
				selectedTrade = await apiClient.getTrade(entry.trade_id);
			} catch (err) {
				console.error('Failed to load trade:', err);
			} finally {
				loadingTrade = false;
			}
		}
	}

	function getAdherenceColor(score: number | undefined): 'success' | 'warning' | 'error' {
		if (!score) return 'error';
		if (score >= 80) return 'success';
		if (score >= 50) return 'warning';
		return 'error';
	}

	function handleEdit(entry: JournalEntry) {
		editingEntry = entry;
		showAddModal = true;
	}

	async function handleDelete(entryId: string) {
		if (!confirm('Are you sure you want to delete this journal entry? This action cannot be undone.')) {
			return;
		}

		try {
			await apiClient.deleteJournalEntry(entryId);
			if (selectedEntry?.id === entryId) {
				selectedEntry = null;
			}
			await loadData();
		} catch (err) {
			console.error('Failed to delete journal entry:', err);
			alert('Failed to delete journal entry. Please try again.');
		}
	}
</script>

<svelte:head>
	<title>Journal - TradePulse</title>
</svelte:head>

<!-- Header -->
<div class="flex items-center justify-between mb-6">
	<div>
		<h1 class="text-3xl font-bold mb-2">Trading Journal</h1>
		<p class="text-slate-600 dark:text-slate-400">
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
		<Icon icon="mdi:loading" class="animate-spin text-4xl text-blue-600" />
	</div>
{:else if entries.length === 0}
	<Card>
		<div class="text-center py-12">
			<Icon icon="mdi:notebook" width="64" class="mx-auto mb-4 text-slate-400" />
			<h3 class="text-xl font-semibold mb-2">No journal entries yet</h3>
			<p class="text-slate-600 dark:text-slate-400 mb-6">
				Start documenting your trading journey, emotional state, and rule adherence
			</p>
			<Button color="primary" onclick={() => (showAddModal = true)}>
				<Icon icon="mdi:pencil" width="20" class="mr-2" />
				Create First Entry
			</Button>
		</div>
	</Card>
{:else}
	<!-- Floating Toggle Button (shown when collapsed) -->
	{#if isListCollapsed}
		<button
			onclick={() => isListCollapsed = false}
			class="hidden lg:flex fixed left-6 top-24 z-10 items-center gap-2 p-3 rounded-lg bg-blue-500 hover:bg-blue-600 text-white shadow-lg transition-all"
			title="Show Journal List"
		>
			<Icon icon="mdi:menu" class="text-xl" />
			<span class="text-sm font-medium">Show List</span>
		</button>
	{/if}

	<!-- Split View Layout -->
	<div class="grid grid-cols-1 gap-6" class:lg:grid-cols-3={!isListCollapsed}>
		<!-- Left Panel: Entry List -->
		<div class="lg:col-span-1 lg:sticky lg:top-6 lg:self-start lg:max-h-[calc(100vh-12rem)] transition-all" class:hidden={isListCollapsed}>
			<!-- Collapse Toggle -->
			<button
				onclick={() => isListCollapsed = !isListCollapsed}
				class="hidden lg:flex items-center gap-2 w-full p-3 mb-3 rounded-lg bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700 transition-colors"
			>
				<Icon icon="mdi:chevron-left" class="text-lg" />
				<span class="text-sm font-medium">Hide List</span>
			</button>

			<div class="space-y-3 lg:overflow-y-auto lg:pr-2">
			{#each entries as entry}
				{@const trade = entryTrades.get(entry.id)}
				<button
					onclick={() => selectEntry(entry)}
					class="w-full text-left p-4 rounded-xl border-2 transition-all {selectedEntry?.id === entry.id
						? 'border-blue-500 dark:border-blue-400 bg-blue-50 dark:bg-blue-900/20 shadow-md'
						: 'border-slate-200 dark:border-slate-700 hover:border-slate-300 dark:hover:border-slate-600 bg-white dark:bg-slate-800'}"
				>
					<div class="flex items-start justify-between mb-2">
						<div class="flex-1 min-w-0">
							<div class="text-sm font-semibold text-slate-900 dark:text-slate-100 truncate">
								{formatDate(entry.entry_date || entry.created_at)}
							</div>
							{#if trade}
								<div class="flex items-center gap-2 mt-1">
									<span class="text-lg font-bold text-slate-900 dark:text-slate-100">{trade.symbol}</span>
									<Badge color={trade.trade_type === 'LONG' ? 'success' : 'error'} variant="soft" size="sm">
										{trade.trade_type}
									</Badge>
								</div>
							{:else if entry.trade_id}
								<div class="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1 mt-1">
									<Icon icon="mdi:chart-line" width="12" />
									Trade linked
								</div>
							{/if}
						</div>
						{#if entry.adherence_score !== undefined && entry.adherence_score !== null}
							<Badge color={getAdherenceColor(entry.adherence_score)} variant="soft" size="sm">
								{Math.round(entry.adherence_score)}%
							</Badge>
						{/if}
					</div>

					{#if trade}
						<div class="grid grid-cols-2 gap-2 mt-2">
							<div class="bg-slate-50 dark:bg-slate-700/50 rounded p-2">
								<div class="text-[10px] text-slate-500 dark:text-slate-400 mb-0.5">Shares</div>
								<div class="text-sm font-bold text-slate-900 dark:text-slate-100">
									{trade.total_entry_quantity || trade.quantity || 0}
								</div>
							</div>
							<div class="bg-slate-50 dark:bg-slate-700/50 rounded p-2">
								<div class="text-[10px] text-slate-500 dark:text-slate-400 mb-0.5">P&L</div>
								<div class="text-sm font-bold {trade.realized_pnl >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}">
									{trade.realized_pnl >= 0 ? '+' : ''}{formatCurrency(trade.realized_pnl)}
								</div>
							</div>
						</div>
						<div class="mt-2">
							<Badge color={trade.realized_pnl >= 0 ? 'success' : 'error'} variant="filled" size="sm">
								{trade.realized_pnl >= 0 ? 'WIN' : 'LOSS'}
							</Badge>
						</div>
					{:else if entry.content}
						{@const parsedContent = (() => {
							try {
								return typeof entry.content === 'string' ? JSON.parse(entry.content) : entry.content;
							} catch {
								return { reflection: entry.content };
							}
						})()}
						<p class="text-xs text-slate-600 dark:text-slate-400 line-clamp-2">
							{parsedContent.reflection || parsedContent.setup || parsedContent.execution || parsedContent.lessonsLearned || 'No content'}
						</p>
					{/if}
				</button>
			{/each}
			</div>
		</div>

		<!-- Right Panel: Detail View (hidden on mobile) -->
		{#if selectedEntry}
			<div class="hidden lg:block lg:col-span-2">
				<Card>
					<div class="space-y-6">
						<!-- Header -->
						<div class="flex items-start justify-between pb-4 border-b border-slate-200 dark:border-slate-700">
							<div class="flex-1">
								<h2 class="text-2xl font-bold text-slate-900 dark:text-slate-100 mb-2">
									{formatDateTime(selectedEntry.entry_date || selectedEntry.created_at)}
								</h2>
								{#if selectedEntry.trade_id}
									<p class="text-sm text-slate-600 dark:text-slate-400 flex items-center gap-1">
										<Icon icon="mdi:chart-line" width="16" />
										Linked to trade
									</p>
								{/if}
							</div>
							<div class="flex gap-2">
								<button
									onclick={() => handleEdit(selectedEntry!)}
									class="p-2 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
									title="Edit entry"
								>
									<Icon icon="mdi:pencil" width="20" />
								</button>
								<button
									onclick={() => handleDelete(selectedEntry!.id)}
									class="p-2 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-lg transition-colors text-red-600 dark:text-red-400"
									title="Delete entry"
								>
									<Icon icon="mdi:delete" width="20" />
								</button>
							</div>
						</div>

						<!-- Trade Details -->
						{#if selectedTrade}
							<div class="bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-800/50 dark:to-slate-900/50 rounded-xl p-6 border border-slate-200 dark:border-slate-700">
								<div class="flex items-start justify-between mb-4">
									<div>
										<h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-1">Trade Details</h3>
										<div class="flex items-center gap-2">
											<span class="text-2xl font-bold text-slate-900 dark:text-slate-100">{selectedTrade.symbol}</span>
											<Badge color={selectedTrade.trade_type === 'LONG' ? 'success' : 'error'} variant="soft">
												{selectedTrade.trade_type}
											</Badge>
											<Badge color={selectedTrade.realized_pnl >= 0 ? 'success' : 'error'} variant="filled">
												{selectedTrade.realized_pnl >= 0 ? '+' : ''}{formatCurrency(selectedTrade.realized_pnl)}
											</Badge>
										</div>
									</div>
									<a
										href="/app/trades?id={selectedTrade.id}"
										class="text-sm text-blue-600 dark:text-blue-400 hover:underline flex items-center gap-1"
									>
										<Icon icon="mdi:open-in-new" width="16" />
										View Trade
									</a>
								</div>

								<div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
									<div>
										<div class="text-xs text-slate-500 dark:text-slate-400 mb-1">Position Size</div>
										<div class="text-lg font-bold text-slate-900 dark:text-slate-100">
											{selectedTrade.total_entry_quantity || selectedTrade.quantity} shares
										</div>
									</div>
									<div>
										<div class="text-xs text-slate-500 dark:text-slate-400 mb-1">Avg Entry</div>
										<div class="text-lg font-bold text-slate-900 dark:text-slate-100">
											{formatCurrency(selectedTrade.average_entry_price)}
										</div>
									</div>
									<div>
										<div class="text-xs text-slate-500 dark:text-slate-400 mb-1">Avg Exit</div>
										<div class="text-lg font-bold text-slate-900 dark:text-slate-100">
											{selectedTrade.exit_price ? formatCurrency(selectedTrade.exit_price) : 'N/A'}
										</div>
									</div>
									<div>
										<div class="text-xs text-slate-500 dark:text-slate-400 mb-1">Total Fees</div>
										<div class="text-lg font-bold text-slate-900 dark:text-slate-100">
											{formatCurrency(selectedTrade.total_fees)}
										</div>
									</div>
								</div>

								<!-- Entries and Exits -->
								<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
									<!-- Entries -->
									<div class="bg-white dark:bg-slate-800 rounded-lg p-4 border border-emerald-200 dark:border-emerald-900/30">
										<div class="flex items-center gap-2 mb-3">
											<Icon icon="mdi:arrow-down-bold" class="text-emerald-600 dark:text-emerald-400" width="20" />
											<h4 class="font-semibold text-slate-900 dark:text-slate-100">
												Entries ({selectedTrade.entries?.length || 0})
											</h4>
										</div>
										{#if selectedTrade.entries && selectedTrade.entries.length > 0}
											<div class="space-y-2 max-h-48 overflow-y-auto">
												{#each selectedTrade.entries as entry}
													<div class="flex items-center justify-between text-sm border-b border-slate-200 dark:border-slate-700 pb-2">
														<div>
															<div class="font-medium text-slate-900 dark:text-slate-100">
																{entry.quantity} @ {formatCurrency(entry.price)}
															</div>
															<div class="text-xs text-slate-500 dark:text-slate-400">
																{formatDateTime(entry.timestamp)}
															</div>
														</div>
														{#if entry.fees}
															<div class="text-xs text-slate-500">
																Fee: {formatCurrency(entry.fees)}
															</div>
														{/if}
													</div>
												{/each}
											</div>
										{:else}
											<p class="text-sm text-slate-500 dark:text-slate-400 italic">No entries recorded</p>
										{/if}
									</div>

									<!-- Exits -->
									<div class="bg-white dark:bg-slate-800 rounded-lg p-4 border border-red-200 dark:border-red-900/30">
										<div class="flex items-center gap-2 mb-3">
											<Icon icon="mdi:arrow-up-bold" class="text-red-600 dark:text-red-400" width="20" />
											<h4 class="font-semibold text-slate-900 dark:text-slate-100">
												Exits ({selectedTrade.exits?.length || 0})
											</h4>
										</div>
										{#if selectedTrade.exits && selectedTrade.exits.length > 0}
											<div class="space-y-2 max-h-48 overflow-y-auto">
												{#each selectedTrade.exits as exit}
													<div class="flex items-center justify-between text-sm border-b border-slate-200 dark:border-slate-700 pb-2">
														<div>
															<div class="font-medium text-slate-900 dark:text-slate-100">
																{exit.quantity} @ {formatCurrency(exit.price)}
															</div>
															<div class="text-xs text-slate-500 dark:text-slate-400">
																{formatDateTime(exit.timestamp)}
															</div>
														</div>
														<div class="text-right">
															<div class="text-xs font-semibold {exit.pnl >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}">
																{exit.pnl >= 0 ? '+' : ''}{formatCurrency(exit.pnl)}
															</div>
															{#if exit.fees}
																<div class="text-xs text-slate-500">
																	Fee: {formatCurrency(exit.fees)}
																</div>
															{/if}
														</div>
													</div>
												{/each}
											</div>
										{:else}
											<p class="text-sm text-slate-500 dark:text-slate-400 italic">No exits recorded</p>
										{/if}
									</div>
								</div>

								<!-- Chart Preview -->
								{#if selectedTrade.entries && selectedTrade.entries.length > 0}
									<div class="mt-4">
										<ChartPreview
											tradeId={selectedTrade.id}
											symbol={selectedTrade.symbol}
											entries={selectedTrade.entries}
											exits={selectedTrade.exits || []}
											openedAt={selectedTrade.opened_at}
											closedAt={selectedTrade.closed_at}
										/>
									</div>
								{/if}
							</div>
						{:else if selectedEntry.trade_id && loadingTrade}
							<div class="bg-slate-50 dark:bg-slate-800/50 rounded-xl p-6 border border-slate-200 dark:border-slate-700">
								<div class="flex items-center justify-center py-8">
									<Icon icon="mdi:loading" class="animate-spin text-3xl text-blue-600" />
									<span class="ml-3 text-slate-600 dark:text-slate-400">Loading trade details...</span>
								</div>
							</div>
						{/if}

						<!-- Emotional State -->
						{#if selectedEntry.emotional_state}
							{@const emotionalState = typeof selectedEntry.emotional_state === 'string' ? JSON.parse(selectedEntry.emotional_state) : selectedEntry.emotional_state}
							<div>
								<h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-4">Emotional State</h3>
								<div class="grid grid-cols-3 gap-4">
									<div>
										<p class="text-sm text-slate-600 dark:text-slate-400 mb-2 flex items-center gap-1">
											<Icon icon="mdi:account-check" class="text-blue-500" />
											Pre-Trade Confidence
										</p>
										<div class="flex items-center gap-2">
											<div class="flex-1 h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
												<div class="h-full bg-blue-600 transition-all" style="width: {(emotionalState.pre_trade_confidence || 0) * 10}%"></div>
											</div>
											<span class="text-sm font-medium">{emotionalState.pre_trade_confidence || 0}/10</span>
										</div>
									</div>
									<div>
										<p class="text-sm text-slate-600 dark:text-slate-400 mb-2 flex items-center gap-1">
											<Icon icon="mdi:lightbulb-outline" class="text-amber-500" />
											Pre-Trade Clarity
										</p>
										<div class="flex items-center gap-2">
											<div class="flex-1 h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
												<div class="h-full bg-amber-600 transition-all" style="width: {(emotionalState.pre_trade_clarity || 0) * 10}%"></div>
											</div>
											<span class="text-sm font-medium">{emotionalState.pre_trade_clarity || 0}/10</span>
										</div>
									</div>
									<div>
										<p class="text-sm text-slate-600 dark:text-slate-400 mb-2 flex items-center gap-1">
											<Icon icon="mdi:shield-check" class="text-emerald-500" />
											Post-Trade Discipline
										</p>
										<div class="flex items-center gap-2">
											<div class="flex-1 h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
												<div class="h-full bg-emerald-600 transition-all" style="width: {(emotionalState.post_trade_discipline || 0) * 10}%"></div>
											</div>
											<span class="text-sm font-medium">{emotionalState.post_trade_discipline || 0}/10</span>
										</div>
									</div>
								</div>
								{#if emotionalState.post_trade_emotion}
									<div class="mt-4 p-3 bg-slate-50 dark:bg-slate-800/50 rounded-lg">
										<p class="text-sm text-slate-600 dark:text-slate-400 mb-1 flex items-center gap-1">
											<Icon icon="mdi:emoticon-outline" class="text-purple-500" />
											Post-Trade Emotion
										</p>
										<p class="text-sm text-slate-700 dark:text-slate-300">{emotionalState.post_trade_emotion}</p>
									</div>
								{/if}
							</div>
						{/if}

						<!-- Journal Content -->
						{#if selectedEntry.content}
							{@const parsedContent = (() => {
								try {
									return typeof selectedEntry.content === 'string' ? JSON.parse(selectedEntry.content) : selectedEntry.content;
								} catch {
									return { reflection: selectedEntry.content };
								}
							})()}
							<div class="space-y-4">
								{#if parsedContent.setup}
									<div>
										<h4 class="text-sm font-semibold text-slate-700 dark:text-slate-300 mb-2">Trade Setup</h4>
										<p class="text-slate-700 dark:text-slate-300 whitespace-pre-wrap">{parsedContent.setup}</p>
									</div>
								{/if}
								{#if parsedContent.execution}
									<div>
										<h4 class="text-sm font-semibold text-slate-700 dark:text-slate-300 mb-2">Execution</h4>
										<p class="text-slate-700 dark:text-slate-300 whitespace-pre-wrap">{parsedContent.execution}</p>
									</div>
								{/if}
								{#if parsedContent.reflection}
									<div>
										<h4 class="text-sm font-semibold text-slate-700 dark:text-slate-300 mb-2">Reflection</h4>
										<p class="text-slate-700 dark:text-slate-300 whitespace-pre-wrap">{parsedContent.reflection}</p>
									</div>
								{/if}
								{#if parsedContent.lessonsLearned}
									<div>
										<h4 class="text-sm font-semibold text-slate-700 dark:text-slate-300 mb-2">Lessons Learned</h4>
										<p class="text-slate-700 dark:text-slate-300 whitespace-pre-wrap">{parsedContent.lessonsLearned}</p>
									</div>
								{/if}
							</div>
						{/if}

						<!-- Rule Adherence -->
						{#if selectedEntry.rule_adherence && selectedEntry.rule_adherence.length > 0}
							<div>
								<h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-3">Rule Adherence</h3>
								<AdherenceScoreDisplay adherences={selectedEntry.rule_adherence} {rules} showDetails={true} />
							</div>
						{/if}

						<!-- Screenshots -->
						{#if selectedEntry.screenshots && selectedEntry.screenshots.length > 0}
							<div>
								<h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-3">Screenshots</h3>
								<ImageGallery images={selectedEntry.screenshots} alt="Trade screenshot" />
							</div>
						{/if}

						<!-- Voice Notes -->
						{#if selectedEntry.voice_notes && selectedEntry.voice_notes.length > 0}
							<div>
								<h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-3">Voice Notes</h3>
								<div class="space-y-2">
									{#each selectedEntry.voice_notes as voiceNote, index}
										<AudioPlayer src={voiceNote} label="Voice Note {index + 1}" />
									{/each}
								</div>
							</div>
						{/if}
					</div>
				</Card>
			</div>
		{/if}
	</div>
{/if}

<!-- Journal Entry Slide-over -->
<JournalFormSlideOver
	open={showAddModal}
	{rules}
	onClose={() => (showAddModal = false)}
	onSubmit={handleCreateEntry}
/>
