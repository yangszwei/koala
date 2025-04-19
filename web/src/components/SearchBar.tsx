import { useEffect, useRef, useState } from 'react';
import { Icon } from '@iconify/react';
import mdiMagnify from '@iconify-icons/mdi/magnify';
import { useTermSuggestions } from '@/hooks/terms';

import type { FormEvent } from 'react';

interface Props {
	className?: string;
	query: string;
	setQuery: (value: string) => void;
	onSubmit: (e: FormEvent) => void;
}

export default function SearchBar({ query, setQuery, onSubmit }: Props) {
	const { suggestions } = useTermSuggestions(query);
	const [selectedIndex, setSelectedIndex] = useState(-1);
	const containerRef = useRef<HTMLDivElement>(null);
	const [showSuggestions, setShowSuggestions] = useState(false);

	useEffect(() => {
		setSelectedIndex(-1);
	}, [query]);

	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			if (containerRef.current) {
				if (!containerRef.current.contains(event.target as Node)) {
					setShowSuggestions(false);
				}
			}
		};
		document.addEventListener('mousedown', handleClickOutside);
		return () => {
			document.removeEventListener('mousedown', handleClickOutside);
		};
	}, [query]);

	return (
		<div ref={containerRef} className="relative w-full max-w-2xl">
			<div className="overflow-hidden rounded-full border border-gray-300 bg-white shadow-md">
				<form
					onSubmit={(e) => {
						onSubmit(e);
						setShowSuggestions(false);
					}}
					className="flex h-12 items-stretch"
				>
					<input
						type="text"
						value={query}
						onChange={(e) => setQuery(e.target.value)}
						onFocus={() => setShowSuggestions(true)}
						onClick={() => setShowSuggestions(true)}
						onKeyDown={(e) => {
							if (e.key === 'ArrowDown') {
								e.preventDefault();
								setSelectedIndex((prev) => (prev + 1) % suggestions.length);
							} else if (e.key === 'ArrowUp') {
								e.preventDefault();
								setSelectedIndex((prev) => (prev - 1 + suggestions.length) % suggestions.length);
							} else if (e.key === 'Enter' && selectedIndex >= 0) {
								e.preventDefault();
								const selected = suggestions[selectedIndex];
								if (selected) {
									setQuery(selected);
								}
							}
						}}
						placeholder="Search medical images or reports..."
						className="h-full flex-grow px-6 text-base placeholder:select-none focus:outline-none"
					/>
					<button
						type="submit"
						className="h-full cursor-pointer bg-[#e9f1f2] px-6 text-sm font-medium text-[#24808B] transition hover:bg-[#dce9ea]"
					>
						<Icon icon={mdiMagnify} className="text-xl" />
					</button>
				</form>
			</div>

			{showSuggestions && query.trim() !== '' && suggestions.length > 0 && (
				<ul className="absolute top-full left-0 z-10 mt-1 w-full divide-y divide-gray-200 rounded-md border border-gray-300 bg-white shadow-md">
					{suggestions.map((suggestion, index) => (
						<li
							key={index}
							className={`cursor-pointer px-4 py-2 hover:bg-gray-100 ${index === selectedIndex ? 'bg-gray-200' : ''}`}
							onMouseEnter={() => setSelectedIndex(index)}
							onClick={() => {
								setQuery(suggestion);
								setSelectedIndex(-1);
								setShowSuggestions(true);
							}}
						>
							{suggestion}
						</li>
					))}
				</ul>
			)}
		</div>
	);
}
