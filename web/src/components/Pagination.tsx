export interface PaginationProps {
	current: number;
	total: number;
	onChange: (page: number) => void;
}

export default function Pagination({ current, total, onChange }: PaginationProps) {
	const getPageRange = (delta = 2) => {
		const start = Math.max(1, current - delta);
		const end = Math.min(total, current + delta);

		return [
			...(start > 2 ? [1, '...'] : start > 1 ? [1] : []),
			...Array.from({ length: end - start + 1 }, (_, i) => start + i),
			...(end < total - 1 ? ['...', total] : end < total ? [total] : []),
		];
	};

	const visiblePages = getPageRange();

	if (total < 1) {
		return null; // No pagination needed if there's no data
	}

	return (
		<div className="flex items-center justify-center text-sm text-gray-700">
			<div className="flex max-w-full flex-wrap justify-center gap-2 px-2">
				<button
					onClick={() => onChange(Math.max(current - 1, 1))}
					disabled={current === 1}
					className="cursor-pointer rounded px-2 py-1 hover:bg-gray-100 disabled:text-gray-300"
				>
					←
				</button>
				{visiblePages.map((page, idx) =>
					typeof page === 'number' ? (
						<button
							key={page}
							onClick={() => onChange(page)}
							className={`cursor-pointer rounded px-3 py-1 hover:bg-gray-100 ${
								page === current ? 'bg-gray-200 font-bold' : ''
							}`}
						>
							{page}
						</button>
					) : (
						<span key={`ellipsis-${idx}`} className="px-2 py-1 text-gray-400">
							...
						</span>
					),
				)}
				<button
					onClick={() => onChange(Math.min(current + 1, total))}
					disabled={current === total}
					className="cursor-pointer rounded px-2 py-1 hover:bg-gray-100 disabled:text-gray-300"
				>
					→
				</button>
			</div>
		</div>
	);
}
