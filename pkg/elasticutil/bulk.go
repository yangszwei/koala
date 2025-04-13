package elasticutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

// BulkInsertChan reads from a channel of docs, batching and sending to ES.
func BulkInsertChan[T any](
	ctx context.Context,
	es *elasticsearch.Client,
	indexName string,
	ch <-chan T,
	idFunc func(T) string,
	batchSize int,
) error {
	var batch []T

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case doc, ok := <-ch:
			if !ok {
				if len(batch) > 0 {
					return sendBatch(ctx, es, indexName, batch, idFunc)
				}
				return nil
			}

			batch = append(batch, doc)
			if len(batch) >= batchSize {
				if err := sendBatch(ctx, es, indexName, batch, idFunc); err != nil {
					return err
				}
				batch = batch[:0]
			}
		}
	}
}

// sendBatch converts a slice of documents into a bulk API request.
func sendBatch[T any](
	ctx context.Context,
	es *elasticsearch.Client,
	index string,
	batch []T,
	idFunc func(T) string,
) error {
	var buf bytes.Buffer

	for _, doc := range batch {
		meta := map[string]map[string]string{
			"index": {
				"_index": index,
				"_id":    idFunc(doc),
			},
		}
		metaLine, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("marshal meta: %w", err)
		}
		docLine, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("marshal doc: %w", err)
		}
		buf.Write(metaLine)
		buf.WriteByte('\n')
		buf.Write(docLine)
		buf.WriteByte('\n')
	}

	res, err := es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk insert error: %s", res.String())
	}

	return nil
}
