package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yangszwei/go-micala/internal/usecase/search"
	"github.com/yangszwei/go-micala/pkg/elasticutil"
)

// dicomwebClient implements the Client interface for DICOMweb data sources.
type dicomwebClient struct {
	name   string
	base   string
	client *http.Client
}

// NewDICOMwebClient creates a new DICOMweb client.
func NewDICOMwebClient(name, base string) Client {
	return &dicomwebClient{
		name:   name,
		base:   strings.TrimRight(base, "/"),
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (d *dicomwebClient) Name() string {
	return elasticutil.EscapeQueryString(d.name)
}

func (d *dicomwebClient) List(ctx context.Context, offset, limit int) ([]DataSummary, error) {
	url := fmt.Sprintf("%s/studies?offset=%d&limit=%d", d.base, offset, limit)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var studies []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&studies); err != nil {
		return nil, err
	}

	summaries := make([]DataSummary, 0, len(studies))
	for _, s := range studies {
		id, _ := s["0020000D"].(map[string]interface{}) // StudyInstanceUID
		uid, _ := id["Value"].([]interface{})
		if len(uid) > 0 {
			summaries = append(summaries, DataSummary{
				ID:     uid[0].(string),
				Source: d.Name(),
				Type:   "study",
				Raw:    s,
			})
		}
	}

	return summaries, nil
}

func (d *dicomwebClient) Fetch(ctx context.Context, summary DataSummary) (*search.Document, error) {
	studyUID := summary.ID

	url := fmt.Sprintf("%s/studies/%s/metadata", d.base, studyUID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var metadata []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, err
	}

	return mapMetadataToDocument(summary.DocID(), metadata), nil
}

func (d *dicomwebClient) Count(_ context.Context) (int, error) {
	return -1, nil // not supported
}

// mapMetadataToDocument converts DICOM metadata to a search.Document.
func mapMetadataToDocument(id string, metadata []map[string]interface{}) *search.Document {
	doc := &search.Document{
		ID:         id, // Elasticsearch document ID
		Type:       "image",
		Modality:   "UNKNOWN",
		Categories: []string{"unsorted"},
	}

	for _, elem := range metadata {
		if tag, ok := elem["00080060"].(map[string]interface{}); ok {
			if val, ok := tag["Value"].([]interface{}); ok && len(val) > 0 {
				doc.Modality = val[0].(string)
			}
		}
		if tag, ok := elem["00100020"].(map[string]interface{}); ok {
			if val, ok := tag["Value"].([]interface{}); ok && len(val) > 0 {
				doc.PatientID = val[0].(string)
			}
		}
		if tag, ok := elem["00100010"].(map[string]interface{}); ok {
			if val, ok := tag["Value"].([]interface{}); ok && len(val) > 0 {
				doc.PatientName = val[0].(string)
			}
		}
		if tag, ok := elem["00100040"].(map[string]interface{}); ok {
			if val, ok := tag["Value"].([]interface{}); ok && len(val) > 0 {
				doc.Gender = mapDICOMGender(val[0].(string))
			}
		}
		if tag, ok := elem["00080020"].(map[string]interface{}); ok {
			if val, ok := tag["Value"].([]interface{}); ok && len(val) > 0 {
				raw := val[0].(string)
				if len(raw) == 8 {
					doc.StudyDate = fmt.Sprintf("%s-%s-%s", raw[:4], raw[4:6], raw[6:])
				} else {
					doc.StudyDate = raw
				}
			}
		}
		// No standard DICOM tags for report text or impression in image metadata.
		// These fields are typically extracted from Structured Reports or external systems.
		doc.ReportText = ""
		doc.Impression = ""
	}
	return doc
}

// mapDICOMGender maps DICOM gender values (e.g., "M", "F", "O", etc.) to "male", "female", or "other".
func mapDICOMGender(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "M":
		return "male"
	case "F":
		return "female"
	case "O":
		return "other"
	default:
		return "unknown"
	}
}
