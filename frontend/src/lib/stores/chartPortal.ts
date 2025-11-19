import { writable } from 'svelte/store';
import type { Entry, Exit } from '$lib/types';

export interface ChartPortalData {
	tradeId?: string;
	symbol: string;
	entries: Entry[];
	exits: Exit[];
	openedAt?: string;
	closedAt?: string | null;
}

interface ChartPortalState {
	isOpen: boolean;
	data: ChartPortalData | null;
}

function createChartPortalStore() {
	const { subscribe, set, update } = writable<ChartPortalState>({
		isOpen: false,
		data: null
	});

	return {
		subscribe,
		open: (data: ChartPortalData) => {
			set({ isOpen: true, data });
		},
		close: () => {
			set({ isOpen: false, data: null });
		},
		toggle: () => {
			update(state => ({ ...state, isOpen: !state.isOpen }));
		}
	};
}

export const chartPortal = createChartPortalStore();
