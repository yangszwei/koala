/** Represents a unified study document stored in Elasticsearch. */
export interface Document {
	id: string;
	type: 'image' | 'report' | 'report_image';
	studyDate: string;
	modality: string;
	patientId: string;
	patientName: string;
	gender: string;
	categories: string[];
	reportText: string;
	impression: string;
}

/** Represents a single search result with relevance score. */
export interface Result {
	document: Document;
	score: number;
}

/** Represents a category and the number of documents in that category. */
export interface CategoryBucket {
	key: string;
	doc_count: number;
}

/** Parameters used to query the search API. */
export interface Query {
	search: string;
	type?: string;
	modality?: string;
	patientId?: string;
	patientName?: string;
	fromDate?: string;
	toDate?: string;
	gender?: string[];
	category: string[];
	limit?: number;
	offset?: number;
}
