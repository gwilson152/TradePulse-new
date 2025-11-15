<script lang="ts">
	import Icon from '@iconify/svelte';
	import { notificationStore, unreadNotifications } from '$lib/stores/notifications';
	import NotificationPanel from './NotificationPanel.svelte';

	let showPanel = $state(false);
	let unreadCount = $derived($unreadNotifications.length);

	function togglePanel() {
		showPanel = !showPanel;
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.notification-bell') && !target.closest('.notification-panel')) {
			showPanel = false;
		}
	}
</script>

<svelte:window on:click={handleClickOutside} />

<div class="notification-bell">
	<button class="bell-button" onclick={togglePanel} aria-label="Notifications">
		<Icon icon="mdi:bell" width="20" />
		{#if unreadCount > 0}
			<span class="badge">{unreadCount > 99 ? '99+' : unreadCount}</span>
		{/if}
		{#if !$notificationStore.connected}
			<span class="disconnected-indicator" title="Disconnected"></span>
		{/if}
	</button>

	{#if showPanel}
		<NotificationPanel onClose={() => (showPanel = false)} />
	{/if}
</div>

<style>
	.notification-bell {
		position: relative;
	}

	.bell-button {
		position: relative;
		background: transparent;
		border: none;
		cursor: pointer;
		padding: 0.5rem;
		border-radius: 0.5rem;
		color: #6b7280;
		transition: all 0.2s;
	}

	.bell-button:hover {
		background: #f3f4f6;
		color: #111827;
	}

	.badge {
		position: absolute;
		top: 0;
		right: 0;
		background: #ef4444;
		color: white;
		font-size: 0.625rem;
		font-weight: 600;
		padding: 0.125rem 0.375rem;
		border-radius: 999px;
		min-width: 1.25rem;
		text-align: center;
	}

	.disconnected-indicator {
		position: absolute;
		bottom: 0.25rem;
		right: 0.25rem;
		width: 8px;
		height: 8px;
		background: #fbbf24;
		border-radius: 50%;
		border: 2px solid white;
	}
</style>
