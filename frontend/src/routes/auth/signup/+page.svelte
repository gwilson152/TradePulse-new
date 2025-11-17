<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Icon from '@iconify/svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import PricingCard from '$lib/components/ui/PricingCard.svelte';
	import BetaBadge from '$lib/components/ui/BetaBadge.svelte';
	import { apiClient } from '$lib/api/client';

	// Check if plan is specified in query params
	const planParam = $page.url.searchParams.get('plan');
	const initialPlan = (planParam === 'starter' || planParam === 'pro' || planParam === 'premium') ? planParam : 'pro';
	const initialStep = planParam ? 'email' : 'plan';

	let step = $state<'plan' | 'email' | 'success'>(initialStep);
	let selectedPlan = $state<'starter' | 'pro' | 'premium'>(initialPlan);
	let email = $state('');
	let isLoading = $state(false);
	let error = $state('');

	const plans = [
		{
			type: 'starter' as const,
			title: 'Starter',
			description: 'Perfect for beginners',
			price: '$2.99',
			originalPrice: undefined,
			features: [
				'Up to 100 trades/month',
				'Basic analytics',
				'1 trading account',
				'CSV import',
				'Mobile access',
				'Email support'
			],
			recommended: false
		},
		{
			type: 'pro' as const,
			title: 'Pro',
			description: 'For serious traders',
			price: '$9.99',
			originalPrice: undefined,
			features: [
				'Unlimited trades',
				'Advanced analytics',
				'Up to 5 trading accounts',
				'CSV + API import',
				'Auto-sync with brokers',
				'Priority support',
				'Custom reports',
				'Trade replay'
			],
			recommended: true
		},
		{
			type: 'premium' as const,
			title: 'Premium',
			description: 'For professional traders',
			price: '$14.99',
			originalPrice: undefined,
			features: [
				'Everything in Pro',
				'Unlimited accounts',
				'AI-powered insights',
				'Dedicated support',
				'White-label reports',
				'Team collaboration'
			],
			recommended: false
		}
	];

	function selectPlan(planType: 'starter' | 'pro' | 'premium') {
		selectedPlan = planType;
		step = 'email';
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();
		error = '';
		isLoading = true;

		try {
			await apiClient.signup(email, selectedPlan);
			step = 'success';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create account';
		} finally {
			isLoading = false;
		}
	}

	function backToPlanSelection() {
		step = 'plan';
		error = '';
	}
</script>

<svelte:head>
	<title>Sign Up - TradePulse</title>
</svelte:head>

<div class="signup-page">
	<div class="w-full max-w-6xl mx-auto py-8 px-4">
		{#if step === 'plan'}
			<!-- Plan Selection Step -->
			<div class="text-center mb-12">
				<div class="flex items-center justify-center gap-3 mb-4">
					<h1 class="text-4xl font-bold text-white">Choose Your Plan</h1>
					<BetaBadge size="lg" />
				</div>
				<p class="text-xl text-white/90 mb-2">All plans are free during Beta</p>
				<p class="text-white/70">No credit card required</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
				{#each plans as plan}
					<PricingCard
						planType={plan.type}
						title={plan.title}
						description={plan.description}
						price={plan.price}
						originalPrice={plan.originalPrice}
						features={plan.features}
						recommended={plan.recommended}
						selected={selectedPlan === plan.type}
						onSelect={() => selectPlan(plan.type)}
					/>
				{/each}
			</div>

			<div class="text-center mt-8">
				<p class="text-white/80 mb-2">Already have an account?</p>
				<a
					href="/auth/login"
					class="text-white hover:text-white/80 font-medium underline"
				>
					Sign in instead
				</a>
			</div>
		{:else if step === 'email'}
			<!-- Email Step -->
			<div class="max-w-md mx-auto">
				<div class="text-center mb-8">
					<h1 class="text-3xl font-bold text-white mb-2">TradePulse</h1>
					<div class="flex items-center justify-center gap-2 mb-4">
						<p class="text-white/90">Creating your</p>
						<span class="px-3 py-1 bg-white/20 backdrop-blur-sm rounded-full text-white font-semibold">
							{plans.find((p) => p.type === selectedPlan)?.title}
						</span>
						<p class="text-white/90">account</p>
					</div>
				</div>

				<Card>
					<form onsubmit={handleSubmit} class="space-y-4">
						<div>
							<h2 class="text-2xl font-bold mb-2">Enter your email</h2>
							<p class="text-surface-600 dark:text-surface-400">
								We'll send you a magic link to verify your account
							</p>
						</div>

						<Input
							type="email"
							bind:value={email}
							label="Email address"
							placeholder="you@example.com"
							required={true}
							disabled={isLoading}
						/>

						{#if error}
							<div
								class="p-3 bg-error-100 dark:bg-error-900/30 border border-error-500 rounded-lg"
							>
								<p class="text-error-700 dark:text-error-400 text-sm">{error}</p>
							</div>
						{/if}

						<div class="flex gap-3">
							<Button
								type="button"
								onclick={backToPlanSelection}
								variant="soft"
								color="neutral"
								disabled={isLoading}
								class="flex-1"
							>
								<Icon icon="mdi:arrow-left" width="20" class="mr-2" />
								Back
							</Button>

							<Button
								type="submit"
								disabled={isLoading || !email}
								variant="filled"
								color="primary"
								class="flex-1"
							>
								{#if isLoading}
									<Icon icon="mdi:loading" width="20" class="animate-spin mr-2" />
									<span>Creating account...</span>
								{:else}
									<span>Create account</span>
								{/if}
							</Button>
						</div>

						<p class="text-xs text-surface-500 text-center">
							By creating an account, you agree to our Terms of Service and Privacy Policy
						</p>
					</form>
				</Card>

				<div class="text-center mt-6">
					<a href="/" class="text-sm text-white/80 hover:text-white"> ← Back to home </a>
				</div>
			</div>
		{:else if step === 'success'}
			<!-- Success Step -->
			<div class="max-w-md mx-auto">
				<Card>
					<div class="text-center space-y-4">
						<div
							class="w-16 h-16 bg-success-100 dark:bg-success-900/30 rounded-full flex items-center justify-center mx-auto"
						>
							<Icon icon="mdi:email" width="32" class="text-success-600 dark:text-success-400" />
						</div>
						<h2 class="text-xl font-semibold">Check your email</h2>
						<p class="text-surface-600 dark:text-surface-400">
							We've sent a magic link to <strong>{email}</strong>
						</p>
						<p class="text-sm text-surface-500">
							Click the link in the email to verify your account and start using TradePulse. The
							link expires in 15 minutes.
						</p>

						<div class="pt-4">
							<a href="/auth/login" class="text-primary-600 hover:text-primary-700 text-sm">
								Go to login page
							</a>
						</div>
					</div>
				</Card>

				<div class="text-center mt-6">
					<a href="/" class="text-sm text-white/80 hover:text-white"> ← Back to home </a>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.signup-page {
		min-height: 100vh;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		background-attachment: fixed;
		overflow-y: auto;
	}

	:global(body) {
		overflow-y: auto !important;
		height: auto !important;
	}

	:global(html) {
		overflow-y: auto !important;
	}
</style>
