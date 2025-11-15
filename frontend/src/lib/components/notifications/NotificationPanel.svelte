<script lang="ts">
	import Icon from '@iconify/svelte';
	import { notificationStore } from '$lib/stores/notifications';
	import NotificationItem from './NotificationItem.svelte';

	interface Props {
		onClose: () => void;
	}

	let { onClose }: Props = $props();

	function markAllAsRead() {
		notificationStore.markAllAsRead();
	}

	function clearAll() {
		notificationStore.clear();
	}
</script>

<div class="notification-panel">
	<div class="panel-header">
		<h3>Notifications</h3>
		<div class="panel-actions">
			{#if $notificationStore.unreadCount > 0}
				<button class="action-btn" onclick={markAllAsRead}>Mark all read</button>
			{/if}
			{#if $notificationStore.notifications.length > 0}
				<button class="action-btn" onclick={clearAll}>Clear all</button>
			{/if}
		</div>
	</div>

	<div class="panel-body">
		{#if $notificationStore.notifications.length === 0}
			<div class="empty-state">
				<Icon icon="mdi:bell" width="48" style="opacity: 0.5" />
				<p>No notifications</p>
			</div>
		{:else}
			<div class="notification-list">
				{#each $notificationStore.notifications as notification (notification.id)}
					<NotificationItem {notification} />
				{/each}
			</div>
		{/if}
	</div>

	<div class="panel-footer">
		<span class="connection-status">
			<span class="status-dot" class:connected={$notificationStore.connected}></span>
			{$notificationStore.connected ? 'Connected' : 'Disconnected'}
		</span>
	</div>
</div>

<style>
	.notification-panel {
		position: absolute;
		top: calc(100% + 0.5rem);
		right: 0;
		width: 380px;
		max-height: 600px;
		background: white;
		border-radius: 0.5rem;
		box-shadow:
			0 10px 15px -3px rgba(0, 0, 0, 0.1),
			0 4px 6px -2px rgba(0, 0, 0, 0.05);
		z-index: 50;
		display: flex;
		flex-direction: column;
	}

	.panel-header {
		padding: 1rem;
		border-bottom: 1px solid #e5e7eb;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.panel-header h3 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: #111827;
	}

	.panel-actions {
		display: flex;
		gap: 0.5rem;
	}

	.action-btn {
		background: transparent;
		border: none;
		color: #3b82f6;
		font-size: 0.75rem;
		cursor: pointer;
		padding: 0.25rem 0.5rem;
		border-radius: 0.25rem;
		transition: background 0.2s;
	}

	.action-btn:hover {
		background: #eff6ff;
	}

	.panel-body {
		flex: 1;
		overflow-y: auto;
		min-height: 200px;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem 1rem;
		color: #9ca3af;
	}

	.empty-state p {
		margin: 0;
		font-size: 0.875rem;
	}

	.notification-list {
		display: flex;
		flex-direction: column;
	}

	.panel-footer {
		padding: 0.75rem 1rem;
		border-top: 1px solid #e5e7eb;
		display: flex;
		justify-content: center;
	}

	.connection-status {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.75rem;
		color: #6b7280;
	}

	.status-dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: #ef4444;
	}

	.status-dot.connected {
		background: #10b981;
	}
</style>
