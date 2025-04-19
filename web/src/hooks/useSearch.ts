import { useEffect, useState } from 'react';
import { apiBase } from '@/configs/path';

import type { Query, Result } from '@/models/search';

/**
 * Custom React hook to perform search queries against the `/api/search` endpoint.
 *
 * @param query - The search query string.
 * @returns An object containing:
 *
 *   - `data`: Array of search results (Document and score).
 *   - `loading`: Whether the request is in progress.
 *   - `error`: Any error encountered during the fetch.
 */
export default function useSearch(query: Query) {
	const [data, setData] = useState<Result[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<Error | null>(null);

	useEffect(() => {
		if (!query || !query.search) return;
		setLoading(true);
		const params = new URLSearchParams();
		Object.entries(query).forEach(([key, value]) => {
			if (value !== undefined && value !== null) {
				if (Array.isArray(value)) {
					value.forEach((v) => params.append(key, v));
				} else {
					params.append(key, String(value));
				}
			}
		});
		fetch(`${apiBase}/search?${params.toString()}`)
			.then((res) => res.json())
			.then((data) => {
				setData(data || []);
				setLoading(false);
			})
			.catch((err) => {
				setError(err);
				setLoading(false);
			});
	}, [query]);

	return { data, loading, error };
}
