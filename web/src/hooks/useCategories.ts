import { useEffect, useState } from 'react';
import { apiBase } from '@/configs/path';

import type { CategoryBucket } from '@/models/search';

/**
 * Custom React hook to fetch category suggestions from `/api/search/categories`.
 *
 * @param prefix - Optional string to filter categories by prefix.
 * @returns An object containing:
 *
 *   - `categories`: Array of category buckets (name and count).
 *   - `loading`: Whether the request is in progress.
 */
export default function useCategories(prefix?: string) {
	const [categories, setCategories] = useState<CategoryBucket[]>([]);
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		const url = prefix ? `/search/categories?prefix=${encodeURIComponent(prefix)}` : `/search/categories`;
		setLoading(true);
		fetch(apiBase + url)
			.then((res) => res.json())
			.then((categories: CategoryBucket[]) => {
				setCategories(categories || []);
				setLoading(false);
			})
			.catch(() => setLoading(false));
	}, [prefix]);

	return { categories, loading };
}
