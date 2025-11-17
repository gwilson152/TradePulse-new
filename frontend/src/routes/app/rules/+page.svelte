<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import RuleCard from '$lib/components/trading/RuleCard.svelte';
	import Icon from '@iconify/svelte';
	import type { RuleSet, Rule, RulePhase, RuleCategory } from '$lib/types';
	import { onMount } from 'svelte';
	import { apiClient } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';

	let ruleSets = $state<RuleSet[]>([]);
	let activeRuleSet = $state<RuleSet | null>(null);
	let showRuleSetModal = $state(false);
	let showRuleModal = $state(false);
	let editingRuleSet = $state<RuleSet | null>(null);
	let editingRule = $state<Rule | null>(null);
	let loading = $state(true);

	// Rule Set Form
	let ruleSetForm = $state({
		name: '',
		description: '',
		is_active: true
	});

	// Rule Form
	let ruleForm = $state({
		title: '',
		description: '',
		weight: 3,
		phase: 'PRE_TRADE' as RulePhase,
		category: 'GENERAL' as RuleCategory
	});

	const phaseOptions = [
		{ value: 'PRE_TRADE', label: 'Pre-Trade (Planning)' },
		{ value: 'DURING_TRADE', label: 'During Trade (Management)' },
		{ value: 'POST_TRADE', label: 'Post-Trade (Review)' }
	];

	const categoryOptions = [
		{ value: 'RISK_MANAGEMENT', label: 'Risk Management' },
		{ value: 'ENTRY', label: 'Entry' },
		{ value: 'EXIT', label: 'Exit' },
		{ value: 'POSITION_SIZING', label: 'Position Sizing' },
		{ value: 'TIMING', label: 'Timing' },
		{ value: 'PSYCHOLOGY', label: 'Psychology' },
		{ value: 'GENERAL', label: 'General' }
	];

	const ruleTemplates = [
		{
			title: 'Risk no more than 2% per trade',
			description: 'Position size should be calculated to risk maximum 2% of account',
			weight: 5,
			phase: 'PRE_TRADE' as RulePhase,
			category: 'RISK_MANAGEMENT' as RuleCategory
		},
		{
			title: 'Wait for confirmation candle',
			description: 'Do not enter until there is a confirmation candle in the direction of the trade',
			weight: 4,
			phase: 'PRE_TRADE' as RulePhase,
			category: 'ENTRY' as RuleCategory
		},
		{
			title: 'Set stop loss immediately',
			description: 'Stop loss must be set as soon as position is opened',
			weight: 5,
			phase: 'DURING_TRADE' as RulePhase,
			category: 'RISK_MANAGEMENT' as RuleCategory
		},
		{
			title: 'Move stop to breakeven at 1:1',
			description: 'Once trade reaches 1:1 risk/reward, move stop to breakeven',
			weight: 3,
			phase: 'DURING_TRADE' as RulePhase,
			category: 'EXIT' as RuleCategory
		},
		{
			title: 'Journal every trade',
			description: 'Create a journal entry for every trade documenting the setup and emotions',
			weight: 4,
			phase: 'POST_TRADE' as RulePhase,
			category: 'PSYCHOLOGY' as RuleCategory
		},
		{
			title: 'Only trade during market hours',
			description: 'Avoid trading during pre-market or after-hours sessions',
			weight: 3,
			phase: 'PRE_TRADE' as RulePhase,
			category: 'TIMING' as RuleCategory
		}
	];

	onMount(async () => {
		await loadRuleSets();
	});

	async function loadRuleSets() {
		try {
			loading = true;
			ruleSets = await apiClient.getRuleSets();
			// Set first active rule set or first rule set as active
			if (ruleSets.length > 0) {
				const active = ruleSets.find(rs => rs.is_active);
				activeRuleSet = active || ruleSets[0];
			}
		} catch (err) {
			console.error('Failed to load rule sets:', err);
			toast.error('Failed to load rule sets');
		} finally {
			loading = false;
		}
	}

	function openRuleSetModal(ruleSet?: RuleSet) {
		if (ruleSet) {
			editingRuleSet = ruleSet;
			ruleSetForm = {
				name: ruleSet.name,
				description: ruleSet.description,
				is_active: ruleSet.is_active
			};
		} else {
			editingRuleSet = null;
			ruleSetForm = { name: '', description: '', is_active: true };
		}
		showRuleSetModal = true;
	}

	function openRuleModal(rule?: Rule) {
		if (!activeRuleSet) return;

		if (rule) {
			editingRule = rule;
			ruleForm = {
				title: rule.title,
				description: rule.description,
				weight: rule.weight,
				phase: rule.phase,
				category: rule.category
			};
		} else {
			editingRule = null;
			ruleForm = {
				title: '',
				description: '',
				weight: 3,
				phase: 'PRE_TRADE',
				category: 'GENERAL'
			};
		}
		showRuleModal = true;
	}

	function useTemplate(template: typeof ruleTemplates[0]) {
		ruleForm = { ...template };
		showRuleModal = true;
	}

	async function saveRuleSet() {
		try {
			if (editingRuleSet) {
				await apiClient.updateRuleSet(editingRuleSet.id, ruleSetForm);
				toast.success('Rule set updated successfully');
			} else {
				await apiClient.createRuleSet(ruleSetForm);
				toast.success('Rule set created successfully');
			}
			showRuleSetModal = false;
			await loadRuleSets();
		} catch (err) {
			console.error('Failed to save rule set:', err);
			toast.error('Failed to save rule set');
		}
	}

	async function saveRule() {
		if (!activeRuleSet) return;

		try {
			if (editingRule) {
				await apiClient.updateRule(activeRuleSet.id, editingRule.id, ruleForm);
				toast.success('Rule updated successfully');
			} else {
				await apiClient.addRule(activeRuleSet.id, ruleForm);
				toast.success('Rule added successfully');
			}
			showRuleModal = false;
			await loadRuleSets();
		} catch (err) {
			console.error('Failed to save rule:', err);
			toast.error('Failed to save rule');
		}
	}

	async function deleteRule(rule: Rule) {
		if (!activeRuleSet || !confirm(`Delete rule "${rule.title}"?`)) return;

		try {
			await apiClient.deleteRule(activeRuleSet.id, rule.id);
			toast.success('Rule deleted successfully');
			await loadRuleSets();
		} catch (err) {
			console.error('Failed to delete rule:', err);
			toast.error('Failed to delete rule');
		}
	}

	async function deleteRuleSet(ruleSet: RuleSet) {
		if (!confirm(`Delete rule set "${ruleSet.name}" and all its rules?`)) return;

		try {
			await apiClient.deleteRuleSet(ruleSet.id);
			toast.success('Rule set deleted successfully');
			activeRuleSet = null;
			await loadRuleSets();
		} catch (err) {
			console.error('Failed to delete rule set:', err);
			toast.error('Failed to delete rule set');
		}
	}

	function setActiveRuleSet(ruleSet: RuleSet) {
		activeRuleSet = ruleSet;
	}
</script>

<svelte:head>
	<title>Rule Sets - TradePulse</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold mb-2">Trading Rule Sets</h1>
			<p class="text-surface-600 dark:text-surface-400">
				Define and manage your trading rules for consistent adherence tracking
			</p>
		</div>
		<Button color="primary" onclick={() => openRuleSetModal()}>
			<Icon icon="mdi:plus" class="mr-2" />
			New Rule Set
		</Button>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<Icon icon="mdi:loading" class="animate-spin text-4xl text-primary-500" />
		</div>
	{:else if ruleSets.length === 0}
		<Card>
			<div class="text-center py-12">
				<Icon icon="mdi:checkbox-marked-circle" class="text-6xl text-surface-400 mb-4 mx-auto" />
				<h3 class="text-xl font-semibold mb-2">No Rule Sets Yet</h3>
				<p class="text-surface-600 dark:text-surface-400 mb-6">
					Create your first rule set to start tracking adherence in your journal
				</p>
				<Button color="primary" onclick={() => openRuleSetModal()}>
					<Icon icon="mdi:plus" class="mr-2" />
					Create First Rule Set
				</Button>
			</div>
		</Card>
	{:else}
		<div class="grid grid-cols-12 gap-6">
			<!-- Rule Sets List (Sidebar) -->
			<div class="col-span-4 space-y-3">
				<h3 class="text-sm font-semibold text-surface-700 dark:text-surface-300 uppercase tracking-wide">
					Your Rule Sets
				</h3>
				{#each ruleSets as ruleSet}
					<button
						onclick={() => setActiveRuleSet(ruleSet)}
						class="w-full text-left p-4 rounded-lg border-2 transition-all {activeRuleSet?.id ===
						ruleSet.id
							? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
							: 'border-surface-200 dark:border-surface-700 hover:border-surface-300 dark:hover:border-surface-600 bg-surface-50 dark:bg-surface-800'}"
					>
						<div class="flex items-start justify-between mb-2">
							<h4 class="font-semibold text-surface-900 dark:text-surface-100">
								{ruleSet.name}
							</h4>
							{#if ruleSet.is_active}
								<span
									class="px-2 py-0.5 bg-success-100 dark:bg-success-900/30 text-success-700 dark:text-success-300 rounded text-xs font-medium"
								>
									Active
								</span>
							{/if}
						</div>
						<p class="text-sm text-surface-600 dark:text-surface-400 mb-2">
							{ruleSet.description}
						</p>
						<p class="text-xs text-surface-500">
							{ruleSet.rules.length} rule{ruleSet.rules.length !== 1 ? 's' : ''}
						</p>
					</button>
				{/each}
			</div>

			<!-- Rule Details (Main Content) -->
			<div class="col-span-8">
				{#if activeRuleSet}
					<Card>
						<div class="space-y-6">
							<!-- Header -->
							<div class="flex items-start justify-between">
								<div>
									<h2 class="text-2xl font-bold mb-2">{activeRuleSet.name}</h2>
									<p class="text-surface-600 dark:text-surface-400">
										{activeRuleSet.description}
									</p>
								</div>
								<div class="flex gap-2">
									<Button variant="ghost" size="sm" onclick={() => openRuleSetModal(activeRuleSet)}>
										<Icon icon="mdi:pencil" class="mr-1" />
										Edit
									</Button>
									<Button
										variant="ghost"
										size="sm"
										onclick={() => deleteRuleSet(activeRuleSet!)}
										class="text-error-600 hover:bg-error-100 dark:hover:bg-error-900/30"
									>
										<Icon icon="mdi:delete" class="mr-1" />
										Delete
									</Button>
								</div>
							</div>

							<!-- Rules -->
							<div>
								<div class="flex items-center justify-between mb-4">
									<h3 class="text-lg font-semibold">Rules ({activeRuleSet.rules.length})</h3>
									<Button color="primary" size="sm" onclick={() => openRuleModal()}>
										<Icon icon="mdi:plus" class="mr-1" />
										Add Rule
									</Button>
								</div>

								{#if activeRuleSet.rules.length === 0}
									<div class="text-center py-8 bg-surface-50 dark:bg-surface-900 rounded-lg border border-surface-200 dark:border-surface-700">
										<Icon icon="mdi:text-box-check" class="text-4xl text-surface-400 mb-2 mx-auto" />
										<p class="text-sm text-surface-600 dark:text-surface-400 mb-4">
											No rules in this set yet
										</p>
										<Button color="primary" size="sm" onclick={() => openRuleModal()}>
											Add First Rule
										</Button>
									</div>
								{:else}
									<div class="space-y-3">
										{#each activeRuleSet.rules as rule}
											<RuleCard
												{rule}
												editable={true}
												onEdit={() => openRuleModal(rule)}
												onDelete={() => deleteRule(rule)}
											/>
										{/each}
									</div>
								{/if}
							</div>

							<!-- Rule Templates -->
							<div class="pt-6 border-t border-surface-200 dark:border-surface-700">
								<h3 class="text-lg font-semibold mb-4">Rule Templates</h3>
								<p class="text-sm text-surface-600 dark:text-surface-400 mb-4">
									Quick-start with common trading rules
								</p>
								<div class="grid grid-cols-2 gap-3">
									{#each ruleTemplates as template}
										<button
											onclick={() => useTemplate(template)}
											class="text-left p-3 rounded-lg border border-surface-200 dark:border-surface-700 hover:border-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 transition-all"
										>
											<h4 class="font-semibold text-sm mb-1 text-surface-900 dark:text-surface-100">
												{template.title}
											</h4>
											<p class="text-xs text-surface-600 dark:text-surface-400">
												{template.description}
											</p>
										</button>
									{/each}
								</div>
							</div>
						</div>
					</Card>
				{:else}
					<Card>
						<div class="text-center py-12">
							<Icon icon="mdi:arrow-left" class="text-4xl text-surface-400 mb-4 mx-auto" />
							<p class="text-surface-600 dark:text-surface-400">
								Select a rule set from the left to view and manage its rules
							</p>
						</div>
					</Card>
				{/if}
			</div>
		</div>
	{/if}
</div>

<!-- Rule Set Modal -->
<Modal
	open={showRuleSetModal}
	onClose={() => (showRuleSetModal = false)}
	title={editingRuleSet ? 'Edit Rule Set' : 'New Rule Set'}
	description="Define a new set of trading rules to track adherence in your journal"
	size="md"
>
	{#snippet children()}
		<form
			onsubmit={(e) => {
				e.preventDefault();
				saveRuleSet();
			}}
			id="ruleset-form"
			class="space-y-4"
		>
			<Input label="Name *" bind:value={ruleSetForm.name} placeholder="e.g., Day Trading Rules" required />

			<Textarea
				label="Description"
				bind:value={ruleSetForm.description}
				placeholder="Describe when to use this rule set..."
				rows={3}
			/>

			<div class="flex items-center gap-2">
				<input type="checkbox" id="is_active" bind:checked={ruleSetForm.is_active} class="rounded" />
				<label for="is_active" class="text-sm text-surface-700 dark:text-surface-300">
					Set as active rule set
				</label>
			</div>
		</form>
	{/snippet}

	{#snippet actionsSlot()}
		<div class="flex gap-3">
			<Button type="button" variant="ghost" onclick={() => (showRuleSetModal = false)} class="flex-1">
				Cancel
			</Button>
			<Button type="submit" form="ruleset-form" color="primary" class="flex-1">
				{editingRuleSet ? 'Update' : 'Create'} Rule Set
			</Button>
		</div>
	{/snippet}
</Modal>

<!-- Rule Modal -->
<Modal
	open={showRuleModal}
	onClose={() => (showRuleModal = false)}
	title={editingRule ? 'Edit Rule' : 'New Rule'}
	description="Add a specific trading rule to track your adherence"
	size="lg"
>
	{#snippet children()}
		<form
			onsubmit={(e) => {
				e.preventDefault();
				saveRule();
			}}
			id="rule-form"
			class="space-y-4"
		>
			<Input label="Rule Title *" bind:value={ruleForm.title} placeholder="e.g., Risk no more than 2% per trade" required />

			<Textarea
				label="Description *"
				bind:value={ruleForm.description}
				placeholder="Describe what this rule means and how to follow it..."
				rows={3}
				required
			/>

			<div class="grid grid-cols-2 gap-4">
				<Select
					label="Phase *"
					options={phaseOptions}
					bind:value={ruleForm.phase}
					required
				/>

				<Select
					label="Category *"
					options={categoryOptions}
					bind:value={ruleForm.category}
					required
				/>
			</div>

			<div>
				<label class="block text-sm font-medium mb-2 text-surface-700 dark:text-surface-300">
					Importance (Weight) *
				</label>
				<div class="flex items-center gap-4">
					<input
						type="range"
						min="1"
						max="5"
						bind:value={ruleForm.weight}
						class="flex-1 h-2 bg-surface-200 dark:bg-surface-700 rounded-lg appearance-none cursor-pointer"
					/>
					<span class="text-yellow-500 text-xl font-bold min-w-[100px]">
						{'★'.repeat(ruleForm.weight)}{'☆'.repeat(5 - ruleForm.weight)}
					</span>
				</div>
				<p class="text-xs text-surface-500 mt-1">
					Higher weight means this rule is more important in adherence calculations
				</p>
			</div>
		</form>
	{/snippet}

	{#snippet actionsSlot()}
		<div class="flex gap-3">
			<Button type="button" variant="ghost" onclick={() => (showRuleModal = false)} class="flex-1">
				Cancel
			</Button>
			<Button type="submit" form="rule-form" color="primary" class="flex-1">
				{editingRule ? 'Update' : 'Add'} Rule
			</Button>
		</div>
	{/snippet}
</Modal>
