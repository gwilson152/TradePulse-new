<script lang="ts">
	import { goto } from '$app/navigation';
	import Icon from '@iconify/svelte';
	import { apiClient } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';

	interface Props {
		userName?: string;
		userEmail?: string;
		userInitials?: string;
	}

	let { userName = 'User', userEmail = '', userInitials = 'U' }: Props = $props();

	let isOpen = $state(false);

	function toggleMenu() {
		isOpen = !isOpen;
	}

	function closeMenu() {
		isOpen = false;
	}

	async function handleLogout() {
		try {
			await apiClient.logout();
			toast.success('Logged out successfully');
			goto('/auth/login');
		} catch (error) {
			toast.error('Failed to logout');
			console.error('Logout error:', error);
		}
	}

	function handleSettings() {
		closeMenu();
		goto('/app/settings');
	}

	// Close menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		if (isOpen && !(event.target as Element).closest('.user-menu-container')) {
			closeMenu();
		}
	}

	// Close menu on Escape key
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && isOpen) {
			closeMenu();
		}
	}
</script>

<svelte:window onclick={handleClickOutside} onkeydown={handleKeydown} />

<div class="user-menu-container relative">
	<!-- Avatar Button -->
	<button
		onclick={toggleMenu}
		class="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-xs font-semibold shadow-sm hover:shadow-md transition-all hover:scale-105"
		aria-label="User menu"
		aria-expanded={isOpen}
		aria-haspopup="true"
	>
		{userInitials}
	</button>

	<!-- Dropdown Menu -->
	{#if isOpen}
		<div
			class="absolute right-0 top-full mt-2 w-64 bg-white/95 dark:bg-slate-800/95 backdrop-blur-xl rounded-xl shadow-2xl border border-slate-200/50 dark:border-slate-700/50 overflow-hidden z-[60]"
			role="menu"
			aria-label="User menu"
		>
			<!-- User Info Section -->
			<div class="px-4 py-3 border-b border-slate-200/50 dark:border-slate-700/50">
				<div class="flex items-center gap-3">
					<div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold shadow-sm">
						{userInitials}
					</div>
					<div class="flex-1 min-w-0">
						<p class="text-sm font-semibold text-slate-800 dark:text-slate-200 truncate">
							{userName}
						</p>
						{#if userEmail}
							<p class="text-xs text-slate-600 dark:text-slate-400 truncate">
								{userEmail}
							</p>
						{/if}
					</div>
				</div>
			</div>

			<!-- Menu Items -->
			<div class="py-2">
				<button
					onclick={handleSettings}
					class="w-full px-4 py-2.5 flex items-center gap-3 hover:bg-slate-100/80 dark:hover:bg-slate-700/80 transition-colors text-left"
					role="menuitem"
				>
					<Icon icon="mdi:cog-outline" width="18" class="text-slate-600 dark:text-slate-400" />
					<span class="text-sm text-slate-700 dark:text-slate-300">Settings</span>
				</button>

				<button
					onclick={handleLogout}
					class="w-full px-4 py-2.5 flex items-center gap-3 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors text-left border-t border-slate-200/50 dark:border-slate-700/50 mt-1 pt-3"
					role="menuitem"
				>
					<Icon icon="mdi:logout" width="18" class="text-red-600 dark:text-red-400" />
					<span class="text-sm text-red-700 dark:text-red-300 font-medium">Log out</span>
				</button>
			</div>
		</div>
	{/if}
</div>

<style>
	/* Ensure dropdown appears above other elements */
	.user-menu-container {
		position: relative;
		z-index: 60;
	}
</style>
