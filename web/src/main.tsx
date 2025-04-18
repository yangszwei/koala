import './index.css';
import './global.css';
import { BrowserRouter, Route, Routes } from 'react-router';
import App from './App.tsx';
import { basename } from '@/configs/path';
import { createRoot } from 'react-dom/client';

createRoot(document.getElementById('root')!).render(
	<BrowserRouter basename={basename}>
		<Routes>
			<Route path="/" element={<App />} />
		</Routes>
	</BrowserRouter>,
);
