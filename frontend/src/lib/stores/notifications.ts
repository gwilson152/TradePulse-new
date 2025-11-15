import { writable, derived } from 'svelte/store';
import type { Writable } from 'svelte/store';

export interface Notification {
	id: string;
	type: string;
	user_id: string;
	title: string;
	message: string;
	data?: any;
	timestamp: string;
	read?: boolean;
}

interface NotificationStore {
	notifications: Notification[];
	unreadCount: number;
	connected: boolean;
}

const initialState: NotificationStore = {
	notifications: [],
	unreadCount: 0,
	connected: false
};

function createNotificationStore() {
	const { subscribe, set, update }: Writable<NotificationStore> = writable(initialState);

	return {
		subscribe,
		add: (notification: Notification) => {
			update((state) => ({
				...state,
				notifications: [notification, ...state.notifications].slice(0, 50), // Keep last 50
				unreadCount: state.unreadCount + 1
			}));
		},
		markAsRead: (id: string) => {
			update((state) => ({
				...state,
				notifications: state.notifications.map((n) =>
					n.id === id ? { ...n, read: true } : n
				),
				unreadCount: Math.max(0, state.unreadCount - 1)
			}));
		},
		markAllAsRead: () => {
			update((state) => ({
				...state,
				notifications: state.notifications.map((n) => ({ ...n, read: true })),
				unreadCount: 0
			}));
		},
		remove: (id: string) => {
			update((state) => {
				const notification = state.notifications.find((n) => n.id === id);
				const wasUnread = notification && !notification.read;
				return {
					...state,
					notifications: state.notifications.filter((n) => n.id !== id),
					unreadCount: wasUnread ? Math.max(0, state.unreadCount - 1) : state.unreadCount
				};
			});
		},
		clear: () => {
			set(initialState);
		},
		setConnected: (connected: boolean) => {
			update((state) => ({
				...state,
				connected
			}));
		}
	};
}

export const notificationStore = createNotificationStore();

// Derived store for unread notifications
export const unreadNotifications = derived(
	notificationStore,
	($store) => $store.notifications.filter((n) => !n.read)
);
