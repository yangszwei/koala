import 'dotenv/config';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import tailwindcss from '@tailwindcss/vite';

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
	base: mode === 'production' ? '/{{.BaseName}}' : '/',
	plugins: [tailwindcss(), react()],
	server: {
		proxy: {
			'/api': {
				target: process.env.VITE_BACKEND_TARGET || 'http://localhost:8080',
				changeOrigin: true,
			},
		},
	},
}));
