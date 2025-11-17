import { writable } from 'svelte/store';

export interface SlideOverPanel {
	id: string;
	component: any;
	props?: Record<string, any>;
	size?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full';
	onClose?: () => void;
}

interface SlideOverState {
	panels: SlideOverPanel[];
	activeIndex: number;
}

function createSlideOverStore() {
	const { subscribe, set, update } = writable<SlideOverState>({
		panels: [],
		activeIndex: -1
	});

	return {
		subscribe,
		open: (panel: SlideOverPanel) => {
			update(state => ({
				panels: [...state.panels, panel],
				activeIndex: state.panels.length
			}));
		},
		close: (id?: string) => {
			update(state => {
				if (id) {
					// Close specific panel
					const index = state.panels.findIndex(p => p.id === id);
					if (index === -1) return state;

					const panel = state.panels[index];
					panel.onClose?.();

					return {
						panels: state.panels.filter(p => p.id !== id),
						activeIndex: Math.max(0, state.activeIndex - (index <= state.activeIndex ? 1 : 0))
					};
				} else {
					// Close top panel
					if (state.panels.length === 0) return state;

					const topPanel = state.panels[state.panels.length - 1];
					topPanel.onClose?.();

					return {
						panels: state.panels.slice(0, -1),
						activeIndex: Math.max(-1, state.activeIndex - 1)
					};
				}
			});
		},
		closeAll: () => {
			update(state => {
				state.panels.forEach(panel => panel.onClose?.());
				return { panels: [], activeIndex: -1 };
			});
		},
		setActive: (index: number) => {
			update(state => ({ ...state, activeIndex: index }));
		}
	};
}

export const slideOverStore = createSlideOverStore();
