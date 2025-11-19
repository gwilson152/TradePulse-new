<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Icon from '@iconify/svelte';
	import NotificationBell from '$lib/components/notifications/NotificationBell.svelte';
	import UserMenu from '$lib/components/ui/UserMenu.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';
	import UniversalSlideOver from '$lib/components/layout/UniversalSlideOver.svelte';
	import ChartPortal from '$lib/components/layout/ChartPortal.svelte';
	import ProfileSetupWizard from '$lib/components/onboarding/ProfileSetupWizard.svelte';
	import { apiClient } from '$lib/api/client';
	import { userStore, getUserInitials, getDisplayName } from '$lib/stores/user';
	import type { User } from '$lib/types';

	let { children } = $props();

	let isAuthenticated = $state(false);
	let currentTime = $state(new Date());
	let currentUser = $state<User | null>(null);
	let showProfileWizard = $state(false);
	let pendingReviewCount = $state(0);
	let showReviewBanner = $state(false);
	let reviewBannerDismissed = $state(false);

	onMount(async () => {
		// Check if review banner was dismissed
		const dismissed = localStorage.getItem('reviewBannerDismissed');
		if (dismissed === 'true') {
			reviewBannerDismissed = true;
		}

		// Check if user has valid JWT token
		const token = apiClient.getToken();
		if (!token) {
			goto('/auth/login');
			return;
		}

		// Fetch current user (if endpoint is available)
		try {
			const user = await apiClient.getCurrentUser();
			currentUser = user;
			userStore.setUser(user);

			// Check if profile is complete
			if (!user.profile_completed && !user.first_name) {
				showProfileWizard = true;
			}
		} catch (error: any) {
			// If authentication fails (401, 403, or user not found), redirect to login
			if (error?.message?.includes('Failed to get user') ||
			    error?.message?.includes('Unauthorized') ||
			    error?.message?.includes('401') ||
			    error?.message?.includes('403')) {
				console.error('Session invalid or expired, redirecting to login');
				apiClient.clearAuthToken();
				goto('/auth/login');
				return;
			}

			console.warn('Could not fetch user details from API, using JWT token data:', error);
			// Fallback: extract user info from JWT token
			const payload = apiClient.getTokenPayload();
			if (payload && payload.email) {
				currentUser = {
					id: payload.user_id || '',
					email: payload.email,
					created_at: '',
					last_login: ''
				};
				userStore.setUser(currentUser);
			} else {
				// No valid token payload, redirect to login
				console.error('No valid session found, redirecting to login');
				apiClient.clearAuthToken();
				goto('/auth/login');
				return;
			}
		}

		isAuthenticated = true;

		// Check for pending reviews
		checkPendingReviews();

		// Update time every minute
		const interval = setInterval(() => {
			currentTime = new Date();
		}, 60000);

		// Check for pending reviews every 5 minutes
		const reviewInterval = setInterval(() => {
			checkPendingReviews();
		}, 300000);

		return () => {
			clearInterval(interval);
			clearInterval(reviewInterval);
		};
	});

	async function checkPendingReviews() {
		try {
			const response = await apiClient.getTrades({
				status: 'closed',
				limit: 100
			});
			const pending = response.data.filter(
				(trade: any) => !trade.is_reviewed && !trade.review_skipped
			);
			pendingReviewCount = pending.length;

			// Show banner if there are pending reviews, not dismissed, and we're not already on the review page
			if (pendingReviewCount > 0 && !reviewBannerDismissed && !$page.url.pathname.includes('/review')) {
				showReviewBanner = true;
			} else {
				showReviewBanner = false;
			}
		} catch (err) {
			console.error('Failed to check pending reviews:', err);
		}
	}

	function dismissReviewBanner() {
		showReviewBanner = false;
		reviewBannerDismissed = true;
		localStorage.setItem('reviewBannerDismissed', 'true');
	}

	const navItems = [
		{ href: '/app/dashboard', icon: 'mdi:view-dashboard-outline', label: 'Overview', color: 'text-blue-500' },
		{ href: '/app/accounts', icon: 'mdi:bank-outline', label: 'Accounts', color: 'text-cyan-500' },
		{ href: '/app/trades', icon: 'mdi:chart-line-variant', label: 'Trades', color: 'text-emerald-500' },
		{ href: '/app/review', icon: 'mdi:clipboard-check-outline', label: 'Review', color: 'text-indigo-500' },
		{ href: '/app/journal', icon: 'mdi:book-open-page-variant-outline', label: 'Journal', color: 'text-purple-500' },
		{ href: '/app/rules', icon: 'mdi:checkbox-marked-circle-outline', label: 'Rules', color: 'text-pink-500' },
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
	<UniversalSlideOver />
	<ChartPortal />

	<!-- Profile Setup Wizard -->
	{#if showProfileWizard}
		<ProfileSetupWizard onComplete={() => (showProfileWizard = false)} />
	{/if}

	<div class="min-h-screen bg-gradient-to-br from-slate-50 via-slate-100 to-slate-200 dark:from-slate-950 dark:via-slate-900 dark:to-slate-950">
		<!-- Menu Bar (macOS-style) -->
		<div class="fixed top-0 left-0 right-0 h-11 bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl border-b border-slate-200/50 dark:border-slate-800/50 z-50 flex items-center justify-between px-4">
			<!-- Left: App Name, Plan & Time -->
			<div class="flex items-center gap-6">
				<div class="flex items-center gap-2">
					<div class="w-2 h-2 rounded-full bg-gradient-to-br from-emerald-400 to-emerald-600"></div>
					<span class="text-sm font-semibold bg-gradient-to-r from-slate-700 to-slate-900 dark:from-slate-200 dark:to-slate-400 bg-clip-text text-transparent">
						TradePulse
					</span>
				</div>
				{#if currentUser?.plan_type}
					<span class="text-xs px-2 py-0.5 rounded-full font-semibold bg-gradient-to-r {currentUser.plan_type === 'premium' ? 'from-amber-500 to-orange-500' : currentUser.plan_type === 'pro' ? 'from-purple-500 to-pink-500' : 'from-blue-500 to-cyan-500'} text-white">
						{currentUser.plan_type.charAt(0).toUpperCase() + currentUser.plan_type.slice(1)}
					</span>
				{/if}
				<span class="text-xs text-slate-600 dark:text-slate-400 font-medium">
					{formatTime()}
				</span>
			</div>

			<!-- Right: Notifications & User -->
			<div class="flex items-center gap-3">
				<NotificationBell />
				<UserMenu
					userName={getDisplayName(currentUser)}
					userEmail={currentUser?.email}
					userInitials={getUserInitials(currentUser)}
				/>
			</div>
		</div>

		<!-- Main Content Area -->
		<main class="pt-11 pb-24 px-8 h-screen overflow-y-auto">
			<!-- Pending Reviews Banner -->
			{#if showReviewBanner && pendingReviewCount > 0}
				<div class="max-w-[1800px] mx-auto pt-6 px-0">
					<div class="bg-gradient-to-r from-blue-500 to-indigo-600 text-white rounded-xl shadow-lg p-4 flex items-center justify-between">
						<div class="flex items-center gap-4">
							<div class="bg-white/20 p-3 rounded-lg">
								<Icon icon="mdi:clipboard-check-outline" width="28" />
							</div>
							<div>
								<h3 class="font-bold text-lg">Trade Reviews Pending</h3>
								<p class="text-sm text-blue-50">
									You have {pendingReviewCount} {pendingReviewCount === 1 ? 'trade' : 'trades'} waiting to be reviewed
								</p>
							</div>
						</div>
						<div class="flex items-center gap-3">
							<button
								onclick={() => goto('/app/review')}
								class="bg-white text-blue-600 hover:bg-blue-50 px-6 py-2.5 rounded-lg font-semibold transition-colors flex items-center gap-2 shadow-md"
							>
								<Icon icon="mdi:play-circle" width="20" />
								Start Reviewing
							</button>
							<button
								onclick={dismissReviewBanner}
								class="text-white/80 hover:text-white p-2 rounded-lg hover:bg-white/10 transition-colors"
								title="Dismiss"
							>
								<Icon icon="mdi:close" width="24" />
							</button>
						</div>
					</div>
				</div>
			{/if}

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
