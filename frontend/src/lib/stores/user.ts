import { writable } from 'svelte/store';
import type { User } from '$lib/types';

interface UserStore {
	user: User | null;
	isLoading: boolean;
}

function createUserStore() {
	const { subscribe, set, update } = writable<UserStore>({
		user: null,
		isLoading: false
	});

	return {
		subscribe,
		setUser: (user: User) => update(state => ({ ...state, user })),
		clearUser: () => update(state => ({ ...state, user: null })),
		setLoading: (isLoading: boolean) => update(state => ({ ...state, isLoading }))
	};
}

export const userStore = createUserStore();

// Helper function to get user initials
export function getUserInitials(user: User | null): string {
	if (!user || !user.email) return 'U';
	const email = user.email;
	const name = email.split('@')[0];
	return name.substring(0, 1).toUpperCase();
}

// Helper function to get display name
export function getDisplayName(user: User | null): string {
	if (!user || !user.email) return 'User';
	const email = user.email;
	const name = email.split('@')[0];
	// Capitalize first letter
	return name.charAt(0).toUpperCase() + name.slice(1);
}
