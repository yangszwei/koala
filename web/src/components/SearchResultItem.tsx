import type { Result } from '@/models/search';

export interface SearchResultItemProps {
	result: Result;
}

export default function SearchResultItem({ result }: SearchResultItemProps) {
	const { document, score } = result;
	const { impression, reportText, studyDate, modality, gender, patientName } = document;

	return (
		<div className="flex cursor-pointer gap-4 p-3 transition hover:bg-gray-100">
			{document.type === 'report' ? (
				<img
					src={`https://picsum.photos/seed/${result.document.id}/200`}
					alt="thumbnail"
					draggable="false"
					className="h-32 w-28 flex-shrink-0 rounded object-cover"
				/>
			) : (
				<div className="flex h-32 w-28 flex-shrink-0 items-center justify-center rounded bg-gray-200 text-[10px] text-gray-400 italic">
					No Image
				</div>
			)}
			<div className="flex flex-col gap-2">
				<h3 className="text-lg font-semibold text-gray-800">{impression || 'Untitled'}</h3>
				<p className="line-clamp-2 text-sm text-gray-600">{reportText}</p>
				<div className="mt-1 flex flex-wrap gap-4 text-xs text-gray-500">
					{modality && <span>{modality}</span>}
					{studyDate && <span>{studyDate}</span>}
					{gender && <span>{gender}</span>}
					{patientName && <span>{patientName}</span>}
					<span>Score: {score.toFixed(2)}</span>
				</div>
			</div>
		</div>
	);
}
