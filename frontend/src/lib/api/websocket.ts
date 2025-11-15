import { notificationStore } from '$lib/stores/notifications';
import type { Notification } from '$lib/stores/notifications';

class WebSocketClient {
	private ws: WebSocket | null = null;
	private reconnectInterval: number = 5000;
	private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	private url: string;
	private token: string | null = null;
	private shouldReconnect: boolean = true;

	constructor(baseUrl: string) {
		// Convert https:// to wss:// and http:// to ws://
		this.url = baseUrl.replace(/^https:/, 'wss:').replace(/^http:/, 'ws:');
	}

	connect(token: string) {
		if (this.ws?.readyState === WebSocket.OPEN) {
			console.log('WebSocket already connected');
			return;
		}

		this.token = token;
		this.shouldReconnect = true;

		try {
			// Connect to WebSocket with token in URL
			const wsUrl = `${this.url}/api/ws?token=${token}`;
			console.log('Connecting to WebSocket:', wsUrl);

			this.ws = new WebSocket(wsUrl);

			this.ws.onopen = () => {
				console.log('WebSocket connected');
				notificationStore.setConnected(true);
				if (this.reconnectTimer) {
					clearTimeout(this.reconnectTimer);
					this.reconnectTimer = null;
				}
			};

			this.ws.onmessage = (event) => {
				try {
					const notification: Notification = JSON.parse(event.data);
					console.log('Received notification:', notification);
					notificationStore.add(notification);

					// Show browser notification if permitted
					this.showBrowserNotification(notification);
				} catch (error) {
					console.error('Failed to parse notification:', error);
				}
			};

			this.ws.onerror = (error) => {
				console.error('WebSocket error:', error);
			};

			this.ws.onclose = (event) => {
				console.log('WebSocket disconnected:', event.code, event.reason);
				notificationStore.setConnected(false);
				this.ws = null;

				// Attempt to reconnect if not intentionally closed
				if (this.shouldReconnect && this.token) {
					console.log(`Reconnecting in ${this.reconnectInterval / 1000}s...`);
					this.reconnectTimer = setTimeout(() => {
						this.connect(this.token!);
					}, this.reconnectInterval);
				}
			};
		} catch (error) {
			console.error('Failed to create WebSocket connection:', error);
			notificationStore.setConnected(false);
		}
	}

	disconnect() {
		this.shouldReconnect = false;
		if (this.reconnectTimer) {
			clearTimeout(this.reconnectTimer);
			this.reconnectTimer = null;
		}
		if (this.ws) {
			this.ws.close(1000, 'Client disconnecting');
			this.ws = null;
		}
		notificationStore.setConnected(false);
	}

	isConnected(): boolean {
		return this.ws?.readyState === WebSocket.OPEN;
	}

	private showBrowserNotification(notification: Notification) {
		// Check if browser supports notifications
		if (!('Notification' in window)) {
			return;
		}

		// Check permission
		if (Notification.permission === 'granted') {
			new Notification(notification.title, {
				body: notification.message,
				icon: '/favicon.png',
				tag: notification.id
			});
		} else if (Notification.permission !== 'denied') {
			Notification.requestPermission().then((permission) => {
				if (permission === 'granted') {
					new Notification(notification.title, {
						body: notification.message,
						icon: '/favicon.png',
						tag: notification.id
					});
				}
			});
		}
	}
}

// Export singleton instance
export const wsClient = new WebSocketClient(
	import.meta.env.PUBLIC_API_URL || 'https://api.tradepulse.drivenw.com'
);
