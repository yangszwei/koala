import './index.css';
import { BrowserRouter, Route, Routes } from 'react-router';
import App from './App.tsx';
import { createRoot } from 'react-dom/client';

createRoot(document.getElementById('root')!).render(
	<BrowserRouter basename={import.meta.env.DEV ? '/' : '/{{.BaseName}}'}>
		<Routes>
			<Route path="/" element={<App />} />
		</Routes>
	</BrowserRouter>,
);
