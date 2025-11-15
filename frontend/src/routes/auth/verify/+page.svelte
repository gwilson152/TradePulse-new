<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Icon from '@iconify/svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import { apiClient } from '$lib/api/client';
	import { wsClient } from '$lib/api/websocket';

	let status = $state<'loading' | 'success' | 'error'>('loading');
	let errorMessage = $state('');

	onMount(async () => {
		const token = $page.url.searchParams.get('token');

		if (!token) {
			status = 'error';
			errorMessage = 'No verification token provided';
			return;
		}

		try {
			const result = await apiClient.verifyMagicLink(token);

			// Connect WebSocket with the new JWT
			wsClient.connect(result.jwt);

			status = 'success';

			// Redirect to dashboard after a short delay
			setTimeout(() => {
				goto('/app/dashboard');
			}, 1500);
		} catch (err) {
			status = 'error';
			errorMessage = err instanceof Error ? err.message : 'Failed to verify magic link';
		}
	});
</script>

<svelte:head>
	<title>Verifying - TradePulse</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center p-4">
	<div class="w-full max-w-md">
		<div class="text-center mb-8">
			<h1 class="text-3xl font-bold text-white mb-2">TradePulse</h1>
		</div>

		<Card>
			<div class="text-center space-y-4 py-8">
				{#if status === 'loading'}
					<Icon icon="mdi:loading" width="48" class="animate-spin text-primary-600 mx-auto" />
					<h2 class="text-xl font-semibold">Verifying your link...</h2>
					<p class="text-surface-600 dark:text-surface-400">Please wait a moment</p>
				{:else if status === 'success'}
					<div class="w-16 h-16 bg-success-100 dark:bg-success-900/30 rounded-full flex items-center justify-center mx-auto">
						<Icon icon="mdi:check-circle" width="32" class="text-success-600 dark:text-success-400" />
					</div>
					<h2 class="text-xl font-semibold">Success!</h2>
					<p class="text-surface-600 dark:text-surface-400">
						Redirecting to your dashboard...
					</p>
				{:else}
					<div class="w-16 h-16 bg-error-100 dark:bg-error-900/30 rounded-full flex items-center justify-center mx-auto">
						<Icon icon="mdi:alert-circle" width="32" class="text-error-600 dark:text-error-400" />
					</div>
					<h2 class="text-xl font-semibold">Verification failed</h2>
					<p class="text-surface-600 dark:text-surface-400">{errorMessage}</p>
					<div class="pt-4">
						<a
							href="/auth/login"
							class="inline-block px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
						>
							Request a new link
						</a>
					</div>
				{/if}
			</div>
		</Card>
	</div>
</div>

<style>
	:global(body) {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}
</style>
