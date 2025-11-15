<script lang="ts">
	// Landing page component
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
				<h2>Track Your Trading Journey</h2>
				<p>Advanced journaling and metrics for serious traders</p>
				<a href="/auth/login" class="btn-large">Start Tracking Free</a>
			</div>
		</section>

		<section class="features">
			<div class="container">
				<h3>Why TradePulse?</h3>
				<div class="feature-grid">
					<div class="feature">
						<h4>📊 Comprehensive Metrics</h4>
						<p>Track P&L, win rate, and performance across all your trades</p>
					</div>
					<div class="feature">
						<h4>📝 Advanced Journaling</h4>
						<p>Document your thoughts, emotions, and trade decisions with rich media</p>
					</div>
					<div class="feature">
						<h4>📸 Visual Documentation</h4>
						<p>Upload screenshots and voice notes to capture the full context</p>
					</div>
					<div class="feature">
						<h4>🎯 Emotional Tracking</h4>
						<p>Monitor your emotional state before and after trades</p>
					</div>
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

	.hero h2 {
		font-size: 3rem;
		margin: 0 0 1rem 0;
	}

	.hero p {
		font-size: 1.25rem;
		margin: 0 0 2rem 0;
		opacity: 0.9;
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
