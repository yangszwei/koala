import { useMemo, useState } from 'react';
import { useSearchParams } from 'react-router';

import type { Dispatch, SetStateAction } from 'react';
import type { Query } from '@/models/search';

/**
 * Custom hook that extracts search query parameters from the URL and provides state management with merge support.
 *
 * @returns A tuple of the current Query state and a merge-capable setter function.
 */
export default function useQuery(): [Query, Dispatch<SetStateAction<Partial<Query>>>] {
	const [params] = useSearchParams();

	const initialQuery = useMemo((): Query => {
		/**
		 * Helper function to extract an array of values from a given search param key.
		 *
		 * @param key - The query parameter key.
		 * @returns An array of values for that key, or an empty array if none found.
		 */
		const getArray = (key: string): string[] => {
			const values = params.getAll(key);
			return values.length > 0 ? values : [];
		};

		return {
			...{
				search: params.get('q') ?? '',
				type: params.get('type') ?? undefined,
				modality: params.get('modality') ?? undefined,
				patientId: params.get('patientId') ?? undefined,
				patientName: params.get('patientName') ?? undefined,
				fromDate: params.get('fromDate') ?? undefined,
				toDate: params.get('toDate') ?? undefined,
				gender: getArray('gender'),
				category: getArray('category'),
				limit: params.get('limit') ? parseInt(params.get('limit') as string) : undefined,
				offset: params.get('offset') ? parseInt(params.get('offset') as string) : undefined,
			},
		};
	}, [params]);

	const [query, setQuery] = useState<Query>(initialQuery);

	/**
	 * Merges new query parameters into the current state.
	 *
	 * @param update - A partial query update or a function that returns a partial update.
	 */
	const mergeQuery: Dispatch<SetStateAction<Partial<Query>>> = (update) => {
		setQuery((prev) => {
			const next = typeof update === 'function' ? update(prev) : update;
			return { ...prev, ...next } as Query;
		});
	};

	return [query, mergeQuery];
}
