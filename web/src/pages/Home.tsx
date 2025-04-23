import SearchBar from '@/components/SearchBar';
import textLogo from '@/assets/text-logo.png';
import { useNavigate } from 'react-router';
import { useState } from 'react';

export default function Home() {
	const [query, setQuery] = useState('');
	const navigate = useNavigate();

	const handleSearch = (e: React.FormEvent) => {
		e.preventDefault();
		if (query.trim()) {
			navigate(`/search?q=${encodeURIComponent(query.trim())}`);
		}
	};

	return (
		<div className="flex h-full flex-col items-center justify-start px-4 pt-36">
			{/* Logo or App Title */}
			<h1 className="mb-8 flex items-center justify-center md:mb-12">
				<img src={textLogo} className="h-14 select-none md:h-20" draggable="false" alt="KOALA" />
			</h1>

			{/* Search Bar */}
			<div className="w-full flex-1">
				<div className="flex w-full justify-center">
					<SearchBar className="flex-1" query={query} setQuery={setQuery} onSubmit={handleSearch} />
				</div>
			</div>

			{/* Optional Footer */}
			<footer className="p-3 text-sm text-gray-500 select-none">A smarter search for medical imaging</footer>
		</div>
	);
}
