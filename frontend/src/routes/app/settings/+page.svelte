<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import TabContainer from '$lib/components/ui/TabContainer.svelte';
	import FormSection from '$lib/components/ui/FormSection.svelte';
	import HelpText from '$lib/components/ui/HelpText.svelte';
	import Icon from '@iconify/svelte';
	import { apiClient } from '$lib/api/client';
	import { goto } from '$app/navigation';
	import { toast } from '$lib/stores/toast';
	import { userStore } from '$lib/stores/user';
	import SetPasswordModal from '$lib/components/settings/SetPasswordModal.svelte';
	import type { Rule, RulePhase, RuleCategory } from '$lib/types';

	let activeTab = $state('profile');
	let showPasswordModal = $state(false);
	let currentUser = $state($userStore.user);

	// Mock rules data - replace with API call
	let rules = $state<Rule[]>([
		{
			id: '1',
			title: 'Risk no more than 2% per trade',
			description: 'Position size should be calculated so maximum loss is 2% of account',
			weight: 5,
			phase: 'PRE_TRADE',
			category: 'RISK_MANAGEMENT',
			created_at: new Date().toISOString()
		},
		{
			id: '2',
			title: 'Set stop loss before entry',
			description: 'Always define your stop loss level before entering a position',
			weight: 5,
			phase: 'PRE_TRADE',
			category: 'ENTRY',
			created_at: new Date().toISOString()
		},
		{
			id: '3',
			title: 'Take profit at 1.5:1 R:R minimum',
			description: 'Ensure reward is at least 1.5x the risk before entering',
			weight: 4,
			phase: 'PRE_TRADE',
			category: 'EXIT',
			created_at: new Date().toISOString()
		},
		{
			id: '4',
			title: 'No trading in first 30 minutes',
			description: 'Wait for market open volatility to settle',
			weight: 3,
			phase: 'PRE_TRADE',
			category: 'TIMING',
			created_at: new Date().toISOString()
		}
	]);

	let showAddRule = $state(false);
	let editingRule: Rule | null = $state(null);
	let newRule = $state({
		title: '',
		description: '',
		weight: 3,
		phase: 'PRE_TRADE' as RulePhase,
		category: 'GENERAL' as RuleCategory
	});

	const tabs = [
		{ id: 'profile', label: 'Profile', icon: 'mdi:account' },
		{ id: 'rules', label: 'Trading Rules', icon: 'mdi:gavel' },
		{ id: 'preferences', label: 'Preferences', icon: 'mdi:cog' },
		{ id: 'account', label: 'Account', icon: 'mdi:shield-account' }
	];

	const phaseOptions = [
		{ value: 'PRE_TRADE', label: 'Pre-Trade' },
		{ value: 'DURING_TRADE', label: 'During Trade' },
		{ value: 'POST_TRADE', label: 'Post-Trade' }
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

	const weightOptions = [
		{ value: 1, label: '1 - Low Importance' },
		{ value: 2, label: '2 - Below Average' },
		{ value: 3, label: '3 - Average' },
		{ value: 4, label: '4 - High Importance' },
		{ value: 5, label: '5 - Critical' }
	];

	function handleAddRule() {
		if (!newRule.title || !newRule.description) {
			toast.error('Please fill in all required fields');
			return;
		}

		const rule: Rule = {
			id: crypto.randomUUID(),
			...newRule,
			created_at: new Date().toISOString()
		};

		rules = [...rules, rule];
		toast.success('Rule added successfully');
		resetRuleForm();
	}

	function handleEditRule(rule: Rule) {
		editingRule = rule;
		newRule = {
			title: rule.title,
			description: rule.description,
			weight: rule.weight,
			phase: rule.phase,
			category: rule.category
		};
		showAddRule = true;
	}

	function handleUpdateRule() {
		if (!editingRule) return;

		rules = rules.map((r) =>
			r.id === editingRule!.id
				? { ...r, ...newRule }
				: r
		);

		toast.success('Rule updated successfully');
		resetRuleForm();
	}

	function handleDeleteRule(id: string) {
		if (confirm('Are you sure you want to delete this rule?')) {
			rules = rules.filter((r) => r.id !== id);
			toast.success('Rule deleted');
		}
	}

	function resetRuleForm() {
		showAddRule = false;
		editingRule = null;
		newRule = {
			title: '',
			description: '',
			weight: 3,
			phase: 'PRE_TRADE',
			category: 'GENERAL'
		};
	}

	function getCategoryColor(category: RuleCategory): 'primary' | 'success' | 'warning' | 'error' | 'neutral' {
		switch (category) {
			case 'RISK_MANAGEMENT':
				return 'error';
			case 'ENTRY':
				return 'success';
			case 'EXIT':
				return 'warning';
			case 'PSYCHOLOGY':
				return 'primary';
			default:
				return 'neutral';
		}
	}

	function getPhaseColor(phase: RulePhase): 'primary' | 'success' | 'warning' {
		switch (phase) {
			case 'PRE_TRADE':
				return 'primary';
			case 'DURING_TRADE':
				return 'warning';
			case 'POST_TRADE':
				return 'success';
		}
	}

	async function handleLogout() {
		try {
			await apiClient.logout();
		} catch (err) {
			console.error('Logout error:', err);
		} finally {
			goto('/auth/login');
		}
	}

	async function handlePasswordUpdate() {
		// Refresh user data after password is set
		try {
			const user = await apiClient.getCurrentUser();
			currentUser = user;
			userStore.setUser(user);
		} catch (err) {
			console.error('Failed to refresh user:', err);
		}
	}
</script>

<svelte:head>
	<title>Settings - TradePulse</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold mb-2">Settings</h1>
		<p class="text-surface-600 dark:text-surface-400">
			Manage your account, trading rules, and preferences
		</p>
	</div>

	<TabContainer {tabs} bind:activeTab onTabChange={(tabId) => (activeTab = tabId)}>
		{#if activeTab === 'profile'}
			<Card>
				<FormSection
					title="Profile Information"
					icon="mdi:account-circle"
					helpText="View and manage your account information"
				>
					<div class="space-y-4 max-w-md">
						<Input label="Email" type="email" value="user@example.com" disabled={true} />
						<HelpText
							type="info"
							text="Contact support to change your email address or update profile information."
						/>
					</div>
				</FormSection>
			</Card>
		{/if}

		{#if activeTab === 'rules'}
			<div class="space-y-6">
				<HelpText
					type="info"
					title="Trading Rules Management"
					text="Create and manage your personal trading rules. These rules will be tracked in your journal entries to measure adherence and correlate with performance."
					collapsible={true}
				/>

				<!-- Add/Edit Rule Form -->
				{#if showAddRule}
					<Card>
						<FormSection
							title={editingRule ? 'Edit Rule' : 'Add New Rule'}
							icon="mdi:plus-circle"
						>
							<div class="space-y-4">
								<Input
									label="Rule Title"
									bind:value={newRule.title}
									placeholder="e.g., Risk no more than 2% per trade"
									required={true}
								/>
								<Textarea
									label="Description"
									bind:value={newRule.description}
									placeholder="Explain the rule in detail..."
									rows={3}
									required={true}
								/>
								<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
									<Select
										label="Phase"
										bind:value={newRule.phase}
										options={phaseOptions}
										required={true}
									/>
									<Select
										label="Category"
										bind:value={newRule.category}
										options={categoryOptions}
										required={true}
									/>
									<Select
										label="Importance"
										bind:value={newRule.weight}
										options={weightOptions}
										required={true}
									/>
								</div>
								<HelpText
									type="tip"
									text="Weight (1-5) determines how much this rule affects your adherence score. Higher weights mean the rule is more critical to follow."
								/>
								<div class="flex gap-3">
									{#if editingRule}
										<Button variant="gradient" color="primary" onclick={handleUpdateRule}>
											<Icon icon="mdi:content-save" width="20" />
											Update Rule
										</Button>
									{:else}
										<Button variant="gradient" color="success" onclick={handleAddRule}>
											<Icon icon="mdi:plus" width="20" />
											Add Rule
										</Button>
									{/if}
									<Button variant="soft" color="neutral" onclick={resetRuleForm}>
										Cancel
									</Button>
								</div>
							</div>
						</FormSection>
					</Card>
				{:else}
					<div class="flex justify-end">
						<Button variant="gradient" color="primary" onclick={() => (showAddRule = true)}>
							<Icon icon="mdi:plus" width="20" />
							Add New Rule
						</Button>
					</div>
				{/if}

				<!-- Rules List -->
				<div class="grid grid-cols-1 gap-4">
					{#each rules as rule (rule.id)}
						<Card>
							<div class="flex items-start justify-between gap-4">
								<div class="flex-1">
									<div class="flex items-center gap-2 mb-2">
										<h3 class="text-lg font-semibold text-slate-800 dark:text-slate-200">
											{rule.title}
										</h3>
										<Badge color={getCategoryColor(rule.category)} variant="soft" size="sm">
											{categoryOptions.find((c) => c.value === rule.category)?.label}
										</Badge>
										<Badge color={getPhaseColor(rule.phase)} variant="soft" size="sm">
											{phaseOptions.find((p) => p.value === rule.phase)?.label}
										</Badge>
									</div>
									<p class="text-slate-600 dark:text-slate-400 mb-3">
										{rule.description}
									</p>
									<div class="flex items-center gap-4 text-sm text-slate-500 dark:text-slate-400">
										<div class="flex items-center gap-1">
											<Icon icon="mdi:star" width="16" />
											<span>Importance: {rule.weight}/5</span>
										</div>
										<div class="flex items-center gap-1">
											<Icon icon="mdi:calendar" width="16" />
											<span>Created {new Date(rule.created_at).toLocaleDateString()}</span>
										</div>
									</div>
								</div>
								<div class="flex gap-2">
									<Button
										variant="soft"
										color="primary"
										size="sm"
										onclick={() => handleEditRule(rule)}
									>
										<Icon icon="mdi:pencil" width="16" />
										Edit
									</Button>
									<Button
										variant="soft"
										color="error"
										size="sm"
										onclick={() => handleDeleteRule(rule.id)}
									>
										<Icon icon="mdi:delete" width="16" />
										Delete
									</Button>
								</div>
							</div>
						</Card>
					{/each}

					{#if rules.length === 0}
						<Card>
							<div class="text-center py-8">
								<Icon icon="mdi:gavel" width="48" class="mx-auto mb-4 text-slate-400" />
								<h3 class="text-lg font-semibold mb-2">No Trading Rules Yet</h3>
								<p class="text-slate-600 dark:text-slate-400 mb-4">
									Create your first trading rule to start tracking adherence
								</p>
								<Button variant="gradient" color="primary" onclick={() => (showAddRule = true)}>
									<Icon icon="mdi:plus" width="20" />
									Add Your First Rule
								</Button>
							</div>
						</Card>
					{/if}
				</div>
			</div>
		{/if}

		{#if activeTab === 'preferences'}
			<Card>
				<FormSection
					title="Display Preferences"
					icon="mdi:palette"
					helpText="Customize how TradePulse looks and feels"
				>
					<div class="space-y-6 max-w-md">
						<div>
							<fieldset>
								<legend class="block text-sm font-semibold mb-2 text-slate-700 dark:text-slate-300">
									Theme
								</legend>
								<div class="flex gap-4">
									<label class="flex items-center gap-2 cursor-pointer">
										<input
											type="radio"
											name="theme"
											value="light"
											class="w-4 h-4 text-blue-600 focus:ring-blue-500"
										/>
										<span class="text-sm font-medium text-slate-700 dark:text-slate-300">Light</span>
									</label>
									<label class="flex items-center gap-2 cursor-pointer">
										<input
											type="radio"
											name="theme"
											value="dark"
											checked
											class="w-4 h-4 text-blue-600 focus:ring-blue-500"
										/>
										<span class="text-sm font-medium text-slate-700 dark:text-slate-300">Dark</span>
									</label>
									<label class="flex items-center gap-2 cursor-pointer">
										<input
											type="radio"
											name="theme"
											value="auto"
											class="w-4 h-4 text-blue-600 focus:ring-blue-500"
										/>
										<span class="text-sm font-medium text-slate-700 dark:text-slate-300">Auto</span>
									</label>
								</div>
							</fieldset>
						</div>
						<HelpText
							type="info"
							text="Theme changes will apply immediately. Auto mode follows your system preferences."
						/>
					</div>
				</FormSection>

				<FormSection
					title="Notification Preferences"
					icon="mdi:bell"
					helpText="Manage how you receive notifications"
					collapsible={true}
				>
					<div class="space-y-4 max-w-md">
						<label class="flex items-center gap-3 cursor-pointer">
							<input type="checkbox" class="w-5 h-5 text-blue-600 focus:ring-blue-500 rounded" checked />
							<div>
								<p class="text-sm font-medium text-slate-700 dark:text-slate-300">Email Notifications</p>
								<p class="text-xs text-slate-500 dark:text-slate-400">Receive updates via email</p>
							</div>
						</label>
						<label class="flex items-center gap-3 cursor-pointer">
							<input type="checkbox" class="w-5 h-5 text-blue-600 focus:ring-blue-500 rounded" />
							<div>
								<p class="text-sm font-medium text-slate-700 dark:text-slate-300">Trade Reminders</p>
								<p class="text-xs text-slate-500 dark:text-slate-400">Get reminded to journal your trades</p>
							</div>
						</label>
					</div>
				</FormSection>
			</Card>
		{/if}

		{#if activeTab === 'account'}
			<Card>
				<FormSection title="Security" icon="mdi:shield-lock">
					<div class="space-y-6 max-w-md">
						<div>
							<h3 class="text-lg font-semibold mb-2 text-slate-800 dark:text-slate-200">Password</h3>
							<p class="text-slate-600 dark:text-slate-400 mb-4">
								{#if currentUser?.has_password}
									Change your password for signing in
								{:else}
									Set a password to enable quick sign-in without magic links
								{/if}
							</p>
							<Button variant="soft" color="primary" onclick={() => (showPasswordModal = true)}>
								<Icon icon="mdi:key" width="20" />
								{currentUser?.has_password ? 'Change Password' : 'Set Password'}
							</Button>
						</div>
					</div>
				</FormSection>
			</Card>

			<Card>
				<FormSection title="Account Actions" icon="mdi:account-cog">
					<div class="space-y-6 max-w-md">
						<div>
							<h3 class="text-lg font-semibold mb-2 text-slate-800 dark:text-slate-200">Sign Out</h3>
							<p class="text-slate-600 dark:text-slate-400 mb-4">
								Sign out of your account on this device
							</p>
							<Button variant="soft" color="secondary" onclick={handleLogout}>
								<Icon icon="mdi:logout" width="20" />
								Sign Out
							</Button>
						</div>

						<div class="pt-6 border-t border-slate-200 dark:border-slate-700">
							<h3 class="text-lg font-semibold mb-2 text-red-600 dark:text-red-400">Danger Zone</h3>
							<p class="text-slate-600 dark:text-slate-400 mb-4">
								Permanently delete your account and all data. This action cannot be undone.
							</p>
							<Button variant="soft" color="error" onclick={() => toast.warning('Account deletion is not yet implemented')}>
								<Icon icon="mdi:delete-forever" width="20" />
								Delete Account
							</Button>
						</div>
					</div>
				</FormSection>
			</Card>
		{/if}
	</TabContainer>
</div>

<!-- Password Setup Modal -->
<SetPasswordModal
	open={showPasswordModal}
	hasPassword={currentUser?.has_password || false}
	onClose={() => (showPasswordModal = false)}
	onSuccess={handlePasswordUpdate}
/>
