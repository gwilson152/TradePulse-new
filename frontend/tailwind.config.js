import { skeleton } from '@skeletonlabs/tw-plugin';
import { join } from 'path';

/** @type {import('tailwindcss').Config} */
export default {
	darkMode: 'class',
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		join(require.resolve('@skeletonlabs/skeleton'), '../**/*.{html,js,svelte,ts}')
	],
	theme: {
		extend: {
			colors: {
				// Trading-specific colors
				profit: {
					50: '#f0fdf4',
					100: '#dcfce7',
					500: '#10b981',
					600: '#059669',
					700: '#047857'
				},
				loss: {
					50: '#fef2f2',
					100: '#fee2e2',
					500: '#ef4444',
					600: '#dc2626',
					700: '#b91c1c'
				}
			}
		}
	},
	plugins: [
		skeleton({
			themes: {
				preset: [
					{
						name: 'skeleton',
						enhancements: true
					},
					{
						name: 'modern',
						enhancements: true
					}
				]
			}
		})
	]
};
