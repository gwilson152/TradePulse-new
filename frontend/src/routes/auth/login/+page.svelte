<script lang="ts">
	import Icon from '@iconify/svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { apiClient } from '$lib/api/client';

	let email = $state('');
	let isLoading = $state(false);
	let success = $state(false);
	let error = $state('');

	async function handleSubmit(event: Event) {
		event.preventDefault();
		error = '';
		success = false;
		isLoading = true;

		try {
			await apiClient.requestMagicLink(email);
			success = true;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to send magic link';
		} finally {
			isLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Login - TradePulse</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center p-4">
	<div class="w-full max-w-md">
		<div class="text-center mb-8">
			<h1 class="text-3xl font-bold text-white mb-2">TradePulse</h1>
			<p class="text-white/90">Track your trading journey</p>
		</div>

		<Card>
			{#if success}
				<div class="text-center space-y-4">
					<div class="w-16 h-16 bg-success-100 dark:bg-success-900/30 rounded-full flex items-center justify-center mx-auto">
						<Icon icon="mdi:email" width="32" class="text-success-600 dark:text-success-400" />
					</div>
					<h2 class="text-xl font-semibold">Check your email</h2>
					<p class="text-surface-600 dark:text-surface-400">
						We've sent a magic link to <strong>{email}</strong>
					</p>
					<p class="text-sm text-surface-500">
						Click the link in the email to sign in. The link expires in 15 minutes.
					</p>
					<button
						class="text-primary-600 hover:text-primary-700 text-sm"
						onclick={() => (success = false)}
					>
						Use a different email
					</button>
				</div>
			{:else}
				<form onsubmit={handleSubmit} class="space-y-4">
					<div>
						<h2 class="text-2xl font-bold mb-2">Sign in</h2>
						<p class="text-surface-600 dark:text-surface-400">
							Enter your email to receive a magic link
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
						<div class="p-3 bg-error-100 dark:bg-error-900/30 border border-error-500 rounded-lg">
							<p class="text-error-700 dark:text-error-400 text-sm">{error}</p>
						</div>
					{/if}

					<Button type="submit" disabled={isLoading || !email} variant="filled" color="primary">
						{#if isLoading}
							<span>Sending magic link...</span>
						{:else}
							<span>Send magic link</span>
						{/if}
					</Button>

					<p class="text-xs text-surface-500 text-center">
						By signing in, you agree to our Terms of Service and Privacy Policy
					</p>
				</form>
			{/if}
		</Card>

		<div class="text-center mt-6">
			<a href="/" class="text-sm text-white/80 hover:text-white">
				← Back to home
			</a>
		</div>
	</div>
</div>

<style>
	:global(body) {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}
</style>
