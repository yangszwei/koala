package search

// Document represents a unified study document stored in Elasticsearch.
type Document struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"` //  image | report | report_image
	StudyDate   string   `json:"studyDate"`
	Modality    string   `json:"modality"`
	PatientID   string   `json:"patientId"`
	PatientName string   `json:"patientName"`
	Gender      string   `json:"gender"`
	Categories  []string `json:"categories"`
	ReportText  string   `json:"reportText"`
	Impression  string   `json:"impression"`
}

// Query defines search parameters.
type Query struct {
	Search      string   `form:"search"`
	Type        string   `form:"type"`
	Modality    string   `form:"modality"`
	PatientID   string   `form:"patientId"`
	PatientName string   `form:"patientName"`
	FromDate    string   `form:"fromDate"`
	ToDate      string   `form:"toDate"`
	Gender      []string `form:"gender"`
	Category    []string `form:"category"`
	Limit       int      `form:"limit,default=10"`
	Offset      int      `form:"offset,default=0"`
}

// Result is a single returned hit.
type Result struct {
	Document Document `json:"document"`
	Score    float64  `json:"score"`
}

// CategoryBucket represents a category and its document count.
type CategoryBucket struct {
	Key      string `json:"key"`
	DocCount int64  `json:"doc_count"`
}
