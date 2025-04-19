import './index.css';
import './global.css';
import { BrowserRouter, Route, Routes } from 'react-router';
import Home from '@/pages/Home';
import Search from '@/pages/Search';
import { basename } from '@/configs/path';
import { createRoot } from 'react-dom/client';

createRoot(document.getElementById('root')!).render(
	<BrowserRouter basename={basename}>
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/search" element={<Search />} />
		</Routes>
	</BrowserRouter>,
);
