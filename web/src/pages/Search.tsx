import { useEffect, useState } from 'react';
import CategoryFilter from '@/components/CategoryFilter';
import FilterInput from '@/components/FilterInput';
import Pagination from '@/components/Pagination';
import SearchBar from '@/components/SearchBar';
import SearchResultItem from '@/components/SearchResultItem';
import { basename } from '@/configs/path';
import textLogo from '@/assets/text-logo.png';
import useQuery from '@/hooks/useQuery';
import useSearch from '@/hooks/useSearch';
import { useSearchParams } from 'react-router';

import type { FormEvent } from 'react';

export default function Search() {
	const [query, setQuery] = useQuery();
	const [submittedQuery, setSubmittedQuery] = useState(query);
	const [page, setPage] = useState(1);
	const [showSidebar, setShowSidebar] = useState(false);
	const [, setParams] = useSearchParams();

	const { data, loading, error } = useSearch(submittedQuery);

	const handleSearch = (e: FormEvent) => {
		e.preventDefault();
		setSubmittedQuery(query);
	};

	// Update the URL parameters when the query changes
	useEffect(() => {
		const applyQueryToParams = (query: typeof submittedQuery) => {
			if (!query.search) return;
			const newParams: Record<string, string | string[]> = {};
			if (query.search) newParams.q = query.search;
			if (query.type) newParams.type = query.type;
			if (query.modality) newParams.modality = query.modality;
			if (query.patientId) newParams.patientId = query.patientId;
			if (query.patientName) newParams.patientName = query.patientName;
			if (query.fromDate) newParams.fromDate = query.fromDate;
			if (query.toDate) newParams.toDate = query.toDate;
			if (query.gender?.length) newParams.gender = query.gender;
			if (query.category?.length) newParams.category = query.category;
			if (query.limit !== undefined) newParams.limit = String(query.limit);
			if (query.offset !== undefined) newParams.offset = String(query.offset);
			setParams(newParams);
		};

		applyQueryToParams(submittedQuery);
	}, [setParams, submittedQuery]);

	return (
		<div className="flex h-screen flex-col bg-white">
			<header className="w-full border-b border-gray-200 bg-gray-50">
				<div className="mx-auto flex shrink-0 flex-col items-stretch gap-4 py-3 md:flex-row md:items-center md:gap-0">
					<div className="mt-1 flex w-full min-w-[5rem] justify-center px-4 py-0 md:mt-0 md:w-60 md:justify-start md:py-2">
						<a href={basename.trim() || '/'} target="_self" className="select-none">
							<img src={textLogo} alt="KOALA 🐨" draggable="false" className="h-8 md:h-10" />
						</a>
					</div>
					<div className="flex-grow px-2">
						<SearchBar query={query.search} setQuery={(search) => setQuery({ search })} onSubmit={handleSearch} />
					</div>
				</div>
			</header>
			<div className="flex h-0 flex-1 flex-col overflow-auto md:flex-row md:overflow-hidden">
				<aside
					className={`w-full min-w-[5rem] shrink-0 overflow-y-auto border-b border-gray-200 bg-gray-100 px-4 py-2 ring-gray-200 md:mb-0 md:w-60 md:bg-white ${showSidebar ? '' : 'hidden md:block'}`}
				>
					<h2 className="mb-2 text-base font-medium text-gray-700 select-none">Filters</h2>
					<div className="space-y-2">
						<FilterInput
							label="Patient ID"
							value={query.patientId ?? ''}
							onChange={(v) => setQuery({ patientId: v })}
							onEnter={() => setSubmittedQuery(query)}
						/>
						<FilterInput
							label="Patient Name"
							value={query.patientName ?? ''}
							onChange={(v) => setQuery({ patientName: v })}
							onEnter={() => setSubmittedQuery(query)}
						/>
						<FilterInput
							label="From Date"
							type="date"
							value={query.fromDate ?? ''}
							onChange={(v) => setQuery({ fromDate: v })}
							onEnter={() => setSubmittedQuery(query)}
						/>
						<FilterInput
							label="To Date"
							type="date"
							value={query.toDate ?? ''}
							onChange={(v) => setQuery({ toDate: v })}
							onEnter={() => setSubmittedQuery(query)}
						/>
						<CategoryFilter selected={query.category} onChange={(category) => setQuery({ category })} />
					</div>
				</aside>
				<div className="flex-1 px-2 py-2 md:overflow-y-auto">
					<main className="w-full space-y-2 md:w-fit md:min-w-4/5">
						<div className="flex items-center justify-center md:hidden">
							<button
								type="button"
								className="text-sm font-medium text-gray-600 underline"
								onClick={() => setShowSidebar((prev) => !prev)}
							>
								{showSidebar ? 'Hide Filters' : 'Show Filters'}
							</button>
						</div>
						<div className="space-y-2">
							{loading ? (
								<div className="bg-white py-8 text-center text-gray-500 select-none">Loading results...</div>
							) : error ? (
								<div className="bg-red-50 py-8 text-center text-red-600 select-none">
									Failed to load results. Please try again later.
								</div>
							) : data?.length === 0 ? (
								<div className="bg-white py-8 text-center text-gray-500 select-none">
									No results found for your search.
								</div>
							) : (
								data.map((result, i) => <SearchResultItem key={i} result={result} />)
							)}
						</div>
						<div className="hidden w-full px-2 py-10 md:py-2">
							<Pagination current={page} total={0} onChange={setPage} />
						</div>
					</main>
				</div>
			</div>
		</div>
	);
}
