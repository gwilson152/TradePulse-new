<script lang="ts">
	import Icon from '@iconify/svelte';
	import PricingCard from '$lib/components/ui/PricingCard.svelte';
	import BetaBadge from '$lib/components/ui/BetaBadge.svelte';
	import { goto } from '$app/navigation';

	const plans = [
		{
			type: 'starter' as const,
			title: 'Starter',
			description: 'Perfect for beginners',
			price: '$2.99',
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
		goto(`/auth/signup?plan=${planType}`);
	}
</script>

<div class="landing">
	<header>
		<div class="container">
			<h1 class="logo">TradePulse</h1>
			<nav>
				<a href="/auth/login" class="btn-primary">Get Started</a>
			</nav>
		</div>
	</header>

	<main>
		<section class="hero">
			<div class="container">
				<div class="hero-content">
					<div class="beta-badge-wrapper">
						<BetaBadge size="lg" />
					</div>
					<h2>Track Your Trading Journey</h2>
					<p>Advanced journaling and metrics for serious traders</p>
					<p class="beta-notice">Now in Beta - All plans free!</p>
					<a href="/auth/signup" class="btn-large">Start Free Trial</a>
				</div>
			</div>
		</section>

		<section class="features">
			<div class="container">
				<h3>Why TradePulse?</h3>
				<div class="feature-grid">
					<div class="feature">
						<div class="feature-icon">
							<Icon icon="mdi:chart-line" width="40" />
						</div>
						<h4>Comprehensive Metrics</h4>
						<p>Track P&L, win rate, and performance across all your trades</p>
					</div>
					<div class="feature">
						<div class="feature-icon">
							<Icon icon="mdi:notebook-outline" width="40" />
						</div>
						<h4>Advanced Journaling</h4>
						<p>Document your thoughts, emotions, and trade decisions with rich media</p>
					</div>
					<div class="feature">
						<div class="feature-icon">
							<Icon icon="mdi:image-multiple-outline" width="40" />
						</div>
						<h4>Visual Documentation</h4>
						<p>Upload screenshots and voice notes to capture the full context</p>
					</div>
					<div class="feature">
						<div class="feature-icon">
							<Icon icon="mdi:emoticon-outline" width="40" />
						</div>
						<h4>Emotional Tracking</h4>
						<p>Monitor your emotional state before and after trades</p>
					</div>
				</div>
			</div>
		</section>

		<section class="pricing">
			<div class="container">
				<div class="pricing-header">
					<h3>Choose Your Plan</h3>
					<div class="beta-banner">
						<BetaBadge size="md" />
						<span class="beta-text">All plans are free during Beta - No credit card required</span>
					</div>
				</div>
				<div class="pricing-grid">
					{#each plans as plan}
						<PricingCard
							planType={plan.type}
							title={plan.title}
							description={plan.description}
							price={plan.price}
							features={plan.features}
							recommended={plan.recommended}
							onSelect={() => selectPlan(plan.type)}
						/>
					{/each}
				</div>
			</div>
		</section>
	</main>

	<footer>
		<div class="container">
			<p>&copy; 2024 TradePulse. All rights reserved.</p>
		</div>
	</footer>
</div>

<style>
	.landing {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		overflow-y: auto;
	}

	main {
		flex: 1;
	}

	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 2rem;
	}

	header {
		background: white;
		border-bottom: 1px solid #e5e7eb;
		padding: 1rem 0;
	}

	header .container {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.logo {
		font-size: 1.5rem;
		font-weight: 700;
		color: #3b82f6;
		margin: 0;
	}

	.btn-primary {
		background: #3b82f6;
		color: white;
		padding: 0.5rem 1.5rem;
		border-radius: 0.5rem;
		text-decoration: none;
		font-weight: 500;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background: #2563eb;
	}

	.hero {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		padding: 6rem 0;
		text-align: center;
	}

	.hero-content {
		max-width: 800px;
		margin: 0 auto;
	}

	.beta-badge-wrapper {
		display: flex;
		justify-content: center;
		margin-bottom: 1.5rem;
	}

	.hero h2 {
		font-size: 3rem;
		margin: 0 0 1rem 0;
	}

	.hero p {
		font-size: 1.25rem;
		margin: 0 0 1rem 0;
		opacity: 0.9;
	}

	.beta-notice {
		font-size: 1.125rem;
		font-weight: 600;
		margin: 0 0 2rem 0;
		opacity: 1;
		color: #fbbf24;
	}

	.btn-large {
		background: white;
		color: #667eea;
		padding: 1rem 2rem;
		border-radius: 0.5rem;
		text-decoration: none;
		font-weight: 600;
		font-size: 1.125rem;
		display: inline-block;
		transition: transform 0.2s;
	}

	.btn-large:hover {
		transform: translateY(-2px);
	}

	.features {
		padding: 4rem 0;
	}

	.features h3 {
		text-align: center;
		font-size: 2rem;
		margin: 0 0 3rem 0;
	}

	.feature-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 2rem;
	}

	.feature {
		background: white;
		padding: 2rem;
		border-radius: 0.5rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		text-align: center;
	}

	.feature-icon {
		margin-bottom: 1rem;
		color: #667eea;
		display: flex;
		justify-content: center;
	}

	.feature h4 {
		margin: 0 0 0.5rem 0;
		color: #1f2937;
	}

	.feature p {
		margin: 0;
		color: #6b7280;
		line-height: 1.6;
	}

	.pricing {
		padding: 4rem 0;
		background: #f9fafb;
	}

	.pricing-header {
		text-align: center;
		margin-bottom: 3rem;
	}

	.pricing-header h3 {
		font-size: 2.5rem;
		margin: 0 0 1rem 0;
	}

	.beta-banner {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
		flex-wrap: wrap;
	}

	.beta-text {
		font-size: 1.125rem;
		color: #6b7280;
		font-weight: 500;
	}

	.pricing-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: 2rem;
		max-width: 1200px;
		margin: 0 auto;
	}

	@media (max-width: 768px) {
		.hero h2 {
			font-size: 2rem;
		}

		.pricing-grid {
			grid-template-columns: 1fr;
			max-width: 400px;
		}
	}

	footer {
		background: #1f2937;
		color: white;
		padding: 2rem 0;
		margin-top: auto;
		text-align: center;
	}

	footer p {
		margin: 0;
		opacity: 0.8;
	}
</style>
