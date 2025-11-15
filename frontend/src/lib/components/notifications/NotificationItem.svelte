<script lang="ts">
	import Icon from '@iconify/svelte';
	import { notificationStore } from '$lib/stores/notifications';
	import type { Notification } from '$lib/stores/notifications';

	interface Props {
		notification: Notification;
	}

	let { notification }: Props = $props();

	function getNotificationIcon(type: string): string {
		switch (type) {
			case 'success':
			case 'trade.created':
			case 'journal.created':
				return 'mdi:check';
			case 'error':
				return 'mdi:close';
			case 'info':
				return 'mdi:information';
			case 'trade.updated':
			case 'journal.updated':
				return 'mdi:file-edit';
			case 'trade.deleted':
				return 'mdi:delete';
			case 'csv.import':
				return 'mdi:file-upload';
			default:
				return 'mdi:trending-up';
		}
	}

	function getNotificationColor(type: string) {
		switch (type) {
			case 'success':
			case 'trade.created':
			case 'journal.created':
				return '#10b981';
			case 'error':
				return '#ef4444';
			case 'info':
			case 'trade.updated':
			case 'journal.updated':
				return '#3b82f6';
			case 'csv.import':
				return '#8b5cf6';
			default:
				return '#6b7280';
		}
	}

	function handleClick() {
		if (!notification.read) {
			notificationStore.markAsRead(notification.id);
		}
	}

	function handleDismiss(event: MouseEvent) {
		event.stopPropagation();
		notificationStore.remove(notification.id);
	}

	function formatTime(timestamp: string) {
		const date = new Date(timestamp);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);

		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h ago`;
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays}d ago`;
	}
</script>

<div
	class="notification-item"
	class:unread={!notification.read}
	onclick={handleClick}
	role="button"
	tabindex="0"
	onkeydown={(e) => e.key === 'Enter' && handleClick()}
>
	<div class="notification-icon" style="background-color: {getNotificationColor(notification.type)}">
		<Icon icon={getNotificationIcon(notification.type)} width="16" />
	</div>
	<div class="notification-content">
		<div class="notification-title">{notification.title}</div>
		<div class="notification-message">{notification.message}</div>
		<div class="notification-time">{formatTime(notification.timestamp)}</div>
	</div>
	<button class="dismiss-btn" onclick={handleDismiss} aria-label="Dismiss">
		<Icon icon="mdi:close" width="16" />
	</button>
</div>

<style>
	.notification-item {
		display: flex;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		border: none;
		background: transparent;
		width: 100%;
		text-align: left;
		cursor: pointer;
		border-bottom: 1px solid #f3f4f6;
		transition: background 0.15s;
		position: relative;
	}

	.notification-item:hover {
		background: #f9fafb;
	}

	.notification-item.unread {
		background: #eff6ff;
	}

	.notification-item.unread:hover {
		background: #dbeafe;
	}

	.notification-icon {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		font-size: 1rem;
		flex-shrink: 0;
	}

	.notification-content {
		flex: 1;
		min-width: 0;
	}

	.notification-title {
		font-weight: 600;
		font-size: 0.875rem;
		color: #111827;
		margin-bottom: 0.125rem;
	}

	.notification-message {
		font-size: 0.8125rem;
		color: #6b7280;
		margin-bottom: 0.25rem;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.notification-time {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.dismiss-btn {
		background: transparent;
		border: none;
		color: #9ca3af;
		cursor: pointer;
		padding: 0.25rem;
		border-radius: 0.25rem;
		flex-shrink: 0;
		opacity: 0;
		transition: all 0.15s;
	}

	.notification-item:hover .dismiss-btn {
		opacity: 1;
	}

	.dismiss-btn:hover {
		background: #f3f4f6;
		color: #6b7280;
	}
</style>
