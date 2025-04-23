import { useEffect, useRef, useState } from 'react';
import useCategories from '@/hooks/useCategories';

export interface CategoryFilterProps {
	selected: string[];
	onChange: (selected: string[]) => void;
}

export default function CategoryFilter({ selected, onChange }: CategoryFilterProps) {
	const [inputValue, setInputValue] = useState('');
	const [filterText, setFilterText] = useState('');
	const { categories, loading } = useCategories(filterText);
	const selectedCountsRef = useRef<Map<string, number>>(new Map());

	categories.forEach((cat) => {
		selectedCountsRef.current.set(cat.key, cat.doc_count);
	});

	// Debounce effect
	useEffect(() => {
		const timer = setTimeout(() => setFilterText(inputValue), 300);
		return () => clearTimeout(timer);
	}, [inputValue]);

	const handleToggle = (key: string) => {
		if (selected.includes(key)) {
			onChange(selected.filter((k) => k !== key));
		} else {
			onChange([...selected, key]);
		}
	};

	const categoryMap = new Map(categories.map((c) => [c.key, c.doc_count]));
	selected.forEach((key) => {
		if (!categoryMap.has(key)) {
			const preservedCount = selectedCountsRef.current.get(key) ?? 0;
			categoryMap.set(key, preservedCount);
		}
	});
	const allItems = Array.from(categoryMap.entries()).map(([key, doc_count]) => ({ key, doc_count }));

	return (
		<div>
			<label className="mb-0.5 block text-xs font-medium text-gray-600 select-none">Categories</label>
			<div className="space-y-2">
				<input
					type="text"
					value={inputValue}
					onChange={(e) => setInputValue(e.target.value)}
					className="w-full rounded border border-gray-300 bg-white px-2 py-1 text-xs placeholder:select-none"
					placeholder="Search categories..."
				/>
				{loading ? (
					<div className="text-sm text-gray-500">Loading...</div>
				) : (
					<ul className="space-y-1">
						{allItems.map((cat) => (
							<li key={cat.key} className="flex items-center gap-2 text-sm">
								<input
									id={`cat-${cat.key}`}
									type="checkbox"
									checked={selected.includes(cat.key)}
									onChange={() => handleToggle(cat.key)}
									className="h-3.5 w-3.5"
								/>
								<label htmlFor={`cat-${cat.key}`} className="cursor-pointer text-gray-700 select-none">
									{cat.key} <span className="text-gray-400">({cat.doc_count})</span>
								</label>
							</li>
						))}
					</ul>
				)}
			</div>
		</div>
	);
}
