import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 4000,
		host: true,
		allowedHosts: [
			'tradepulse.drivenw.com',
			'localhost',
			'127.0.0.1'
		]
	}
});
