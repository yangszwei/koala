package datasource

import (
	"context"
	"fmt"

	"github.com/yangszwei/koala/internal/usecase/search"
)

// DataSummary represents a summary returned by a data source.
type DataSummary struct {
	ID     string // Data ID from the data source
	Source string // Data source name (e.g., "dicomweb", "fhir")
	Type   string // study/report/etc.
	Raw    any    // Backend-specific context for Fetch()
}

// DocID builds the document ID of the data summary, used in Elasticsearch.
func (d *DataSummary) DocID() string {
	return fmt.Sprintf("%s:%s:%s", d.Type, d.Source, d.ID)
}

// Client defines the interface for external data sources (e.g., DICOMweb, FHIR).
// It provides methods for fetching and counting documents.
type Client interface {
	// Name returns a sanitized name of the data source client (e.g., "dicomweb", "fhir").
	Name() string
	// Fetch retrieves a full document from a given summary.
	Fetch(ctx context.Context, summary DataSummary) (*search.Document, error)
	// Count returns the total number of documents available.
	// Returns -1 if the operation is not supported.
	Count(ctx context.Context) (int, error)
}

// Pager represents clients that support offset-based pagination.
type Pager interface {
	Client

	// List retrieves a slice of DataSummary entries with offset-based pagination.
	// The result contains at most `limit` items, starting from the specified `offset`.
	List(ctx context.Context, offset, limit int) ([]DataSummary, error)
}

// Streamer interface for cursor-based pagination.
type Streamer interface {
	Client

	// Stream returns a channel of DataSummary entries using cursor-based pagination.
	// This allows clients to consume entries as they become available.
	Stream(ctx context.Context, pageSize int) (<-chan DataSummary, error)
}

// New creates a Client implementation based on the provided type string.
// Supported types include "dicomweb" and "fhir".
func New(typ, name, url string) (Client, error) {
	switch typ {
	case "dicomweb":
		return NewDICOMwebClient(name, url), nil
	case "fhir":
		return NewFHIRClient(name, url), nil
	default:
		return nil, fmt.Errorf("unsupported datasource type: %s", typ)
	}
}
