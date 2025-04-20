package completion

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/yangszwei/koala/internal/domain"
	"github.com/yangszwei/koala/pkg/elasticutil"
	"github.com/yangszwei/koala/pkg/iox"
)

// Service defines autocomplete term operations.
type Service interface {
	// Upload uploads terms to the service.
	Upload(ctx context.Context, terms io.Reader) error
	// Remove removes a term by its ID.
	Remove(ctx context.Context, id string) error
	// Suggest retrieves suggestions based on a query.
	Suggest(ctx context.Context, query string, size int) ([]string, error)
}

const indexName = "terms_completion"

// service implements the autocomplete term operations using Elasticsearch.
type service struct {
	es *elasticsearch.Client
}

// NewService returns a new instance of the completion Service.
func NewService(es *elasticsearch.Client) Service {
	return &service{es: es}
}

// Upload indexes terms from a CSV reader into the Elasticsearch completion index.
func (s *service) Upload(ctx context.Context, terms io.Reader) error {
	terms = iox.StripBOM(terms)
	reader := csv.NewReader(terms)

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %w", err)
	}

	termIdx, err := findTermColumnIndex(headers)
	if err != nil {
		return err
	}

	ch := make(chan domain.CompletionTerm, 1000)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errCh <- fmt.Errorf("failed to read record: %w", err)
				return
			}
			ch <- domain.NewCompletionTerm(record[termIdx])
		}
	}()

	err = elasticutil.BulkInsertChan(ctx, s.es, indexName, ch, func(doc domain.CompletionTerm) string {
		return doc.Term
	}, 1000)
	if err != nil {
		return fmt.Errorf("bulk insert failed: %w", err)
	}

	if readErr := <-errCh; readErr != nil {
		return readErr
	}

	res, err := s.es.Indices.Refresh(
		s.es.Indices.Refresh.WithContext(ctx),
		s.es.Indices.Refresh.WithIndex(indexName),
	)
	if err != nil {
		return fmt.Errorf("failed to refresh index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("refresh index error: %s", res.String())
	}

	return nil
}

// Remove deletes a completion term by ID.
func (s *service) Remove(ctx context.Context, id string) error {
	res, err := s.es.Delete(
		indexName,
		id,
		s.es.Delete.WithContext(ctx),
		s.es.Delete.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && !strings.Contains(res.String(), "not_found") {
		return fmt.Errorf("deleting term failed: %s", res.String())
	}

	return nil
}

// Suggest retrieves term suggestions for a given prefix.
func (s *service) Suggest(ctx context.Context, prefix string, size int) ([]string, error) {
	body := map[string]interface{}{
		"suggest": map[string]interface{}{
			"term-suggest": map[string]interface{}{
				"prefix": prefix,
				"completion": map[string]interface{}{
					"field": "term",
					"size":  size / 2,
				},
			},
			"term-suggest-fuzzy": map[string]interface{}{
				"prefix": prefix,
				"completion": map[string]interface{}{
					"field": "term",
					"fuzzy": true,
					"size":  size / 2,
				},
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(indexName),
		s.es.Search.WithBody(bytes.NewReader(data)),
		s.es.Search.WithTrackTotalHits(false),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("suggestion failed: %s", res.String())
	}

	var parsed struct {
		Suggest map[string][]struct {
			Options []struct {
				Text string `json:"text"`
			} `json:"options"`
		} `json:"suggest"`
	}

	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	var suggestions []string
	seen := make(map[string]struct{})
	for _, key := range []string{"term-suggest", "term-suggest-fuzzy"} {
		for _, result := range parsed.Suggest[key] {
			for _, opt := range result.Options {
				if _, exists := seen[opt.Text]; !exists {
					suggestions = append(suggestions, opt.Text)
					seen[opt.Text] = struct{}{}
				}
			}
		}
	}

	return suggestions, nil
}

// findTermColumnIndex returns the index of the "term" or "terms" column in the header row.
// It returns an error if no such column is found.
func findTermColumnIndex(headers []string) (int, error) {
	for i, h := range headers {
		if strings.EqualFold(h, "term") || strings.EqualFold(h, "terms") {
			return i, nil
		}
	}
	return -1, fmt.Errorf("CSV does not contain a 'term' column")
}
