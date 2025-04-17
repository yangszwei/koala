package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/yangszwei/koala/pkg/elasticutil"
)

// Service defines indexing and search operations for study documents.
type Service interface {
	// Index adds or updates a Document in the search backend.
	Index(ctx context.Context, doc Document) error
	// Search performs a fulltext + metadata search across indexed studies.
	Search(ctx context.Context, query Query) ([]Result, error)
	// ListCategories returns categories that optionally match a given prefix.
	ListCategories(ctx context.Context, prefix string) ([]CategoryBucket, error)
	// Exists checks if a document with the given ID already exists in the index.
	Exists(ctx context.Context, id string) (bool, error)
}

var indexName = "search_documents"

type service struct {
	es *elasticsearch.Client // Elasticsearch client
}

// NewService creates a new search service instance using Elasticsearch and the provided index name.
func NewService(es *elasticsearch.Client) Service {
	return &service{
		es: es,
	}
}

// Index adds or updates the given Document into the Elasticsearch index.
func (s *service) Index(ctx context.Context, doc Document) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal document: %w", err)
	}

	res, err := s.es.Index(
		indexName,
		bytes.NewReader(data),
		s.es.Index.WithDocumentID(doc.ID),
		s.es.Index.WithRefresh("wait_for"),
		s.es.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("index request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("indexing error: %s", res.String())
	}

	return nil
}

// Exists checks if a document with the given ID already exists in the index.
func (s *service) Exists(ctx context.Context, id string) (bool, error) {
	res, err := s.es.Exists(indexName, id, s.es.Exists.WithContext(ctx))
	if err != nil {
		return false, fmt.Errorf("exists request: %w", err)
	}
	defer res.Body.Close()
	return res.StatusCode == 200, nil
}

// Search performs a query on the Elasticsearch index with support for progressive fuzziness.
// It attempts the query using increasing levels of fuzziness ("AUTO", "1", "2") until results are found.
func (s *service) Search(ctx context.Context, q Query) ([]Result, error) {
	var results []Result
	var lastErr error

	for _, fuzziness := range []string{"AUTO", "1", "2"} {
		must := []map[string]interface{}{}
		if q.Search != "" {
			safeQuery := elasticutil.EscapeQueryString(q.Search)
			must = append(must, map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":     safeQuery,
					"fields":    []string{"reportText", "reportText.autocomplete", "reportText.edge_ngram"},
					"fuzziness": fuzziness,
					"operator":  "or",
				},
			})
		}

		filter := []map[string]interface{}{}

		if q.Type != "" {
			escapedType := elasticutil.EscapeQueryString(q.Type)
			fmt.Println(escapedType)
			filter = append(filter, map[string]interface{}{"term": map[string]interface{}{"type": escapedType}})
		}
		if q.Modality != "" {
			escapedModality := elasticutil.EscapeQueryString(q.Modality)
			filter = append(filter, map[string]interface{}{"term": map[string]interface{}{"modality": escapedModality}})
		}
		if q.PatientID != "" {
			escapedPatientID := elasticutil.EscapeQueryString(q.PatientID)
			filter = append(filter, map[string]interface{}{"term": map[string]interface{}{"patientId": escapedPatientID}})
		}
		if len(q.Gender) > 0 {
			escapedGender := elasticutil.EscapeQueryStrings(q.Gender)
			filter = append(filter, map[string]interface{}{"terms": map[string]interface{}{"gender": escapedGender}})
		}
		if len(q.Category) > 0 {
			escapedCategory := elasticutil.EscapeQueryStrings(q.Category)
			filter = append(filter, map[string]interface{}{"terms": map[string]interface{}{"categories": escapedCategory}})
		}
		if q.FromDate != "" || q.ToDate != "" {
			dateRange := map[string]interface{}{}
			if q.FromDate != "" {
				dateRange["gte"] = q.FromDate
			}
			if q.ToDate != "" {
				dateRange["lte"] = q.ToDate
			}
			filter = append(filter, map[string]interface{}{
				"range": map[string]interface{}{
					"studyDate": dateRange,
				},
			})
		}

		queryBody := map[string]interface{}{
			"from": q.Offset,
			"size": q.Limit,
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must":   must,
					"filter": filter,
				},
			},
		}

		fmt.Println(queryBody)

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(queryBody); err != nil {
			return nil, fmt.Errorf("encode query body: %w", err)
		}

		res, err := s.es.Search(
			s.es.Search.WithContext(ctx),
			s.es.Search.WithIndex(indexName),
			s.es.Search.WithBody(&buf),
			s.es.Search.WithTrackTotalHits(true),
		)
		if err != nil {
			return nil, fmt.Errorf("search request: %w", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return nil, fmt.Errorf("search error: %s", res.String())
		}

		var parsed struct {
			Hits struct {
				Hits []struct {
					Source Document `json:"_source"`
					Score  float64  `json:"_score"`
				} `json:"hits"`
			} `json:"hits"`
		}

		if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
			lastErr = fmt.Errorf("decode response: %w", err)
			continue
		}

		if len(parsed.Hits.Hits) > 0 {
			results = make([]Result, len(parsed.Hits.Hits))
			for i, hit := range parsed.Hits.Hits {
				results[i] = Result{
					Document: hit.Source,
					Score:    hit.Score,
				}
			}
			return results, nil
		}
	}

	return results, lastErr
}

// ListCategories returns all unique categories with their document counts,
// optionally filtering by a prefix.
func (s *service) ListCategories(ctx context.Context, prefix string) ([]CategoryBucket, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"categories": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "categories",
					"size":  100,
				},
			},
		},
	}

	if prefix != "" {
		query["query"] = map[string]interface{}{
			"prefix": map[string]interface{}{
				"categories": prefix,
			},
		}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("encode category query: %w", err)
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(indexName),
		s.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("category aggregation request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("category aggregation error: %s", res.String())
	}

	var parsed struct {
		Aggregations struct {
			Categories struct {
				Buckets []CategoryBucket `json:"buckets"`
			} `json:"categories"`
		} `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode category aggregation response: %w", err)
	}

	return parsed.Aggregations.Categories.Buckets, nil
}
