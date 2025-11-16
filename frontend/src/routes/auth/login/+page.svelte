<script lang="ts">
	import { goto } from '$app/navigation';
	import Icon from '@iconify/svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { apiClient } from '$lib/api/client';

	let loginMode = $state<'magic-link' | 'password'>('password');
	let email = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let success = $state(false);
	let error = $state('');

	async function handleMagicLinkSubmit(event: Event) {
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

	async function handlePasswordSubmit(event: Event) {
		event.preventDefault();
		error = '';
		isLoading = true;

		try {
			await apiClient.loginWithPassword(email, password);
			goto('/app/dashboard');
		} catch (err: any) {
			if (err.message?.includes('No password set')) {
				error = 'No password set for this account. Please use magic link login.';
			} else {
				error = err instanceof Error ? err.message : 'Invalid email or password';
			}
		} finally {
			isLoading = false;
		}
	}

	function toggleLoginMode() {
		loginMode = loginMode === 'magic-link' ? 'password' : 'magic-link';
		error = '';
		password = '';
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
				<div class="space-y-4">
					<div>
						<h2 class="text-2xl font-bold mb-2">Sign in</h2>
						<p class="text-surface-600 dark:text-surface-400">
							{#if loginMode === 'magic-link'}
								Enter your email to receive a magic link
							{:else}
								Enter your email and password
							{/if}
						</p>
					</div>

					{#if loginMode === 'password'}
						<form onsubmit={handlePasswordSubmit} class="space-y-4">
							<Input
								type="email"
								bind:value={email}
								label="Email address"
								placeholder="you@example.com"
								required={true}
								disabled={isLoading}
							/>

							<Input
								type="password"
								bind:value={password}
								label="Password"
								placeholder="••••••••"
								required={true}
								disabled={isLoading}
							/>

							{#if error}
								<div class="p-3 bg-error-100 dark:bg-error-900/30 border border-error-500 rounded-lg">
									<p class="text-error-700 dark:text-error-400 text-sm">{error}</p>
								</div>
							{/if}

							<Button type="submit" disabled={isLoading || !email || !password} variant="filled" color="primary">
								{#if isLoading}
									<Icon icon="mdi:loading" width="20" class="animate-spin mr-2" />
									<span>Signing in...</span>
								{:else}
									<span>Sign in with password</span>
								{/if}
							</Button>
						</form>
					{:else}
						<form onsubmit={handleMagicLinkSubmit} class="space-y-4">
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
									<Icon icon="mdi:loading" width="20" class="animate-spin mr-2" />
									<span>Sending magic link...</span>
								{:else}
									<span>Send magic link</span>
								{/if}
							</Button>
						</form>
					{/if}

					<!-- Toggle between login modes -->
					<div class="relative">
						<div class="absolute inset-0 flex items-center">
							<div class="w-full border-t border-surface-300 dark:border-surface-700"></div>
						</div>
						<div class="relative flex justify-center text-sm">
							<span class="px-2 bg-white dark:bg-surface-900 text-surface-500">or</span>
						</div>
					</div>

					<button
						type="button"
						onclick={toggleLoginMode}
						class="w-full text-center text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium"
					>
						{#if loginMode === 'password'}
							Sign in with magic link instead
						{:else}
							Sign in with password instead
						{/if}
					</button>

					<p class="text-xs text-surface-500 text-center">
						By signing in, you agree to our Terms of Service and Privacy Policy
					</p>
				</div>
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
