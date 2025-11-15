import { writable } from 'svelte/store';

export interface Toast {
	id: string;
	message: string;
	type: 'success' | 'error' | 'warning' | 'info';
	duration?: number;
}

function createToastStore() {
	const { subscribe, update } = writable<Toast[]>([]);

	return {
		subscribe,
		success: (message: string, duration = 3000) => {
			const id = crypto.randomUUID();
			update((toasts) => [...toasts, { id, message, type: 'success', duration }]);
			setTimeout(() => remove(id), duration);
		},
		error: (message: string, duration = 4000) => {
			const id = crypto.randomUUID();
			update((toasts) => [...toasts, { id, message, type: 'error', duration }]);
			setTimeout(() => remove(id), duration);
		},
		warning: (message: string, duration = 3500) => {
			const id = crypto.randomUUID();
			update((toasts) => [...toasts, { id, message, type: 'warning', duration }]);
			setTimeout(() => remove(id), duration);
		},
		info: (message: string, duration = 3000) => {
			const id = crypto.randomUUID();
			update((toasts) => [...toasts, { id, message, type: 'info', duration }]);
			setTimeout(() => remove(id), duration);
		},
		remove: (id: string) => {
			update((toasts) => toasts.filter((t) => t.id !== id));
		}
	};

	function remove(id: string) {
		update((toasts) => toasts.filter((t) => t.id !== id));
	}
}

export const toast = createToastStore();
