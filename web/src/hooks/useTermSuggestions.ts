import { useEffect, useState } from 'react';
import { apiBase } from '@/configs/path';

/**
 * Custom React hook to fetch term suggestions for a given query.
 *
 * @param {string} query - The search term used to fetch suggestions.
 * @returns {{ suggestions: string[]; loading: boolean }} An object containing the fetched suggestions and loading
 *   state.
 */
export default function useTermSuggestions(query: string): { suggestions: string[]; loading: boolean } {
	const [suggestions, setSuggestions] = useState<string[]>([]);
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		if (!query) return;
		setLoading(true);
		fetch(`${apiBase}/terms/suggest?q=${encodeURIComponent(query)}`)
			.then((res) => res.json())
			.then((res) => {
				setSuggestions(res.results || []);
				setLoading(false);
			})
			.catch(() => setLoading(false));
	}, [query]);

	return { suggestions, loading };
}
