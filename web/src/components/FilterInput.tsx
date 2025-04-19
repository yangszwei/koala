export interface FilterInputProps {
	label: string;
	value: string;
	onChange: (value: string) => void;
	type?: 'text' | 'date';
	onEnter?: () => void;
}

export default function FilterInput({ label, value, onChange, type = 'text', onEnter }: FilterInputProps) {
	return (
		<div>
			<label className="mb-0.5 block text-xs font-medium text-gray-600 select-none">{label}</label>
			<input
				type={type}
				value={value}
				onChange={(e) => onChange(e.target.value)}
				onKeyDown={(e) => e.key === 'Enter' && onEnter?.()}
				className="w-full rounded border border-gray-300 bg-white px-2 py-0.5 text-sm placeholder:select-none"
			/>
		</div>
	);
}
