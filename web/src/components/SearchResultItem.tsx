import { Icon } from '@iconify/react';
import mdiImage from '@iconify-icons/mdi/image';
import mdiImageBrokenVariant from '@iconify-icons/mdi/image-broken-variant';
import mdiImageOff from '@iconify-icons/mdi/image-off';
import { useState } from 'react';

import type { Result } from '@/models/search';

export interface SearchResultItemProps {
	result: Result;
}

export default function SearchResultItem({ result }: SearchResultItemProps) {
	const { document, score } = result;
	const { impression, reportText, studyDate, modality, gender, patientName } = document;

	/** State to manage thumbnail loading (true = loaded, false = failed, null = loading) */
	const [isThumbnailLoaded, setIsThumbnailLoaded] = useState<boolean | null>(null);

	return (
		<div className="flex cursor-pointer gap-4 p-3 transition hover:bg-gray-100">
			<div className="relative flex h-32 w-28 shrink-0 items-center justify-center overflow-hidden rounded bg-gray-200 text-gray-400">
				{document.type === 'report' ? (
					<>
						{/* TODO: determine thumbnail availability here */}
						{!isThumbnailLoaded && (
							<div className="absolute inset-0 flex items-center justify-center bg-transparent">
								<Icon
									icon={isThumbnailLoaded === false ? mdiImageBrokenVariant : mdiImage}
									className={`h-8 w-8 ${isThumbnailLoaded === null ? 'animate-pulse' : ''}`}
								/>
							</div>
						)}
						<img
							src={`https://picsum.photos/seed/${result.document.id}/200`}
							alt="thumbnail"
							draggable="false"
							className={`h-full w-full object-cover transition-opacity ${isThumbnailLoaded ? 'opacity-100' : 'opacity-0'}`}
							onLoad={() => setIsThumbnailLoaded(true)}
							onError={() => setIsThumbnailLoaded(false)}
						/>
					</>
				) : (
					<Icon icon={mdiImageOff} className="h-8 w-8" />
				)}
			</div>
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
