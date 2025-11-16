<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Icon from '@iconify/svelte';
	import NotificationBell from '$lib/components/notifications/NotificationBell.svelte';
	import UserMenu from '$lib/components/ui/UserMenu.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';
	import { apiClient } from '$lib/api/client';

	let { children } = $props();

	let isAuthenticated = $state(false);
	let currentTime = $state(new Date());

	onMount(async () => {
		// Check if user has valid JWT token
		const token = apiClient.getToken();
		if (!token) {
			goto('/auth/login');
			return;
		}
		isAuthenticated = true;

		// Update time every minute
		const interval = setInterval(() => {
			currentTime = new Date();
		}, 60000);

		return () => clearInterval(interval);
	});

	const navItems = [
		{ href: '/app/dashboard', icon: 'mdi:view-dashboard-outline', label: 'Overview', color: 'text-blue-500' },
		{ href: '/app/trades', icon: 'mdi:chart-line-variant', label: 'Trades', color: 'text-emerald-500' },
		{ href: '/app/journal', icon: 'mdi:book-open-page-variant-outline', label: 'Journal', color: 'text-purple-500' },
		{ href: '/app/analytics', icon: 'mdi:chart-box-outline', label: 'Analytics', color: 'text-orange-500' },
		{ href: '/app/settings', icon: 'mdi:cog-outline', label: 'Settings', color: 'text-slate-500' }
	];

	function formatTime() {
		return currentTime.toLocaleTimeString('en-US', {
			hour: 'numeric',
			minute: '2-digit',
			hour12: true
		});
	}

	function isActive(href: string): boolean {
		return $page.url.pathname.startsWith(href);
	}
</script>

<svelte:head>
	<title>TradePulse</title>
</svelte:head>

{#if isAuthenticated}
	<Toast />
	<div class="min-h-screen bg-gradient-to-br from-slate-50 via-slate-100 to-slate-200 dark:from-slate-950 dark:via-slate-900 dark:to-slate-950">
		<!-- Menu Bar (macOS-style) -->
		<div class="fixed top-0 left-0 right-0 h-11 bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl border-b border-slate-200/50 dark:border-slate-800/50 z-50 flex items-center justify-between px-4">
			<!-- Left: App Name & Time -->
			<div class="flex items-center gap-6">
				<div class="flex items-center gap-2">
					<div class="w-2 h-2 rounded-full bg-gradient-to-br from-emerald-400 to-emerald-600"></div>
					<span class="text-sm font-semibold bg-gradient-to-r from-slate-700 to-slate-900 dark:from-slate-200 dark:to-slate-400 bg-clip-text text-transparent">
						TradePulse
					</span>
				</div>
				<span class="text-xs text-slate-600 dark:text-slate-400 font-medium">
					{formatTime()}
				</span>
			</div>

			<!-- Right: Notifications & User -->
			<div class="flex items-center gap-3">
				<NotificationBell />
				<UserMenu userName="User" userInitials="U" />
			</div>
		</div>

		<!-- Main Content Area -->
		<main class="pt-11 pb-24 px-8 h-screen overflow-y-auto">
			<div class="max-w-[1800px] mx-auto py-8">
				{@render children()}
			</div>
		</main>

		<!-- Dock Navigation (macOS-style bottom dock) -->
		<nav class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50">
			<div class="bg-white/70 dark:bg-slate-800/70 backdrop-blur-2xl rounded-2xl border border-slate-300/50 dark:border-slate-700/50 shadow-2xl px-3 py-2.5">
				<div class="flex items-center gap-2">
					{#each navItems as item}
						<a
							href={item.href}
							class="group relative flex flex-col items-center gap-1.5 px-4 py-2 rounded-xl transition-all duration-200 hover:bg-slate-100/80 dark:hover:bg-slate-700/80"
							class:bg-slate-100={isActive(item.href)}
							class:dark:bg-slate-700={isActive(item.href)}
						>
							<!-- Icon with color -->
							<div class="relative">
								<Icon
									icon={item.icon}
									width="24"
									class="{item.color} transition-transform duration-200 group-hover:scale-110 {isActive(item.href) ? 'scale-110' : ''}"
								/>
								{#if isActive(item.href)}
									<div class="absolute -bottom-3 left-1/2 -translate-x-1/2 w-1 h-1 rounded-full {item.color.replace('text-', 'bg-')}"></div>
								{/if}
							</div>
							<!-- Label -->
							<span class="text-[10px] font-medium text-slate-600 dark:text-slate-400 whitespace-nowrap">
								{item.label}
							</span>
						</a>
					{/each}
				</div>
			</div>
		</nav>
	</div>
{/if}

<style>
	:global(body) {
		margin: 0;
		padding: 0;
	}
</style>
