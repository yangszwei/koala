import 'dotenv/config';
import * as path from 'node:path';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import tailwindcss from '@tailwindcss/vite';

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
	base: mode === 'production' ? '/__KOALA_BASE_PATH__' : '/',
	plugins: [tailwindcss(), react()],
	resolve: {
		alias: {
			'@': path.resolve(__dirname, './src'),
			'index.css': path.resolve(__dirname, './src/index.css'),
		},
	},
	server: {
		proxy: {
			'/api': {
				target: process.env.VITE_BACKEND_TARGET || 'http://localhost:8080',
				changeOrigin: true,
			},
		},
	},
}));
