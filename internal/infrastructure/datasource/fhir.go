package datasource

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yangszwei/go-micala/internal/usecase/search"
	"github.com/yangszwei/go-micala/pkg/elasticutil"
)

// fhirClient implements the Client interface for FHIR data sources.
type fhirClient struct {
	name   string
	base   string
	client *http.Client
}

// NewFHIRClient creates a new FHIR client.
func NewFHIRClient(name, base string) Client {
	return &fhirClient{
		name:   name,
		base:   strings.TrimSuffix(base, "/"),
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (f *fhirClient) Name() string {
	return elasticutil.EscapeQueryString(f.name)
}

func (f *fhirClient) Count(ctx context.Context) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.base+"/DiagnosticReport?_summary=count", nil)
	if err != nil {
		return -1, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var res struct {
		Total int `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return -1, err
	}
	return res.Total, nil
}

func (f *fhirClient) Stream(ctx context.Context, pageSize int) (<-chan DataSummary, error) {
	out := make(chan DataSummary, 100)

	go func() {
		defer close(out)
		url := fmt.Sprintf("%s/DiagnosticReport?_count=%d", f.base, pageSize)

		wait := 250 * time.Millisecond
		slowThreshold := 500 * time.Millisecond

		for {
			start := time.Now()
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return
			}
			resp, err := f.client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			var bundle struct {
				Entry []struct {
					Resource map[string]interface{} `json:"resource"`
				} `json:"entry"`
				Link []struct {
					Relation string `json:"relation"`
					URL      string `json:"url"`
				} `json:"link"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&bundle); err != nil {
				return
			}

			for _, e := range bundle.Entry {
				if id, ok := e.Resource["id"].(string); ok {
					out <- DataSummary{
						ID:     id,
						Source: f.Name(),
						Type:   "report",
						Raw:    e.Resource,
					}
				}
			}

			nextURL := ""
			for _, l := range bundle.Link {
				if l.Relation == "next" {
					nextURL = l.URL
					break
				}
			}
			if nextURL == "" {
				break
			}
			url = nextURL

			elapsed := time.Since(start)
			if elapsed > slowThreshold {
				time.Sleep(wait)
				if wait < 30*time.Second {
					wait *= 2
				}
			} else {
				wait = 250 * time.Millisecond
			}
		}
	}()

	return out, nil
}

func (f *fhirClient) Fetch(_ context.Context, summary DataSummary) (*search.Document, error) {
	res := summary.Raw.(map[string]interface{})
	doc := &search.Document{
		ID:         summary.DocID(),
		Type:       "report",
		Categories: []string{},
	}

	if ts, ok := res["effectiveDateTime"].(string); ok && len(ts) >= 10 {
		doc.StudyDate = ts[:10]
	}

	if cat, ok := res["category"].([]interface{}); ok {
		for _, item := range cat {
			if coding, ok := item.(map[string]interface{})["coding"].([]interface{}); ok {
				for _, c := range coding {
					if cd, ok := c.(map[string]interface{})["display"].(string); ok {
						doc.Categories = append(doc.Categories, cd)
					}
				}
			}
		}
	}

	if concl, ok := res["conclusion"].(string); ok {
		doc.Impression = concl
	}

	if forms, ok := res["presentedForm"].([]interface{}); ok && len(forms) > 0 {
		if data, ok := forms[0].(map[string]interface{})["data"].(string); ok {
			if decoded, err := base64.StdEncoding.DecodeString(data); err == nil {
				doc.ReportText = string(decoded)
			}
		}
	}

	if subj, ok := res["subject"].(map[string]interface{}); ok {
		if ref, ok := subj["reference"].(string); ok && strings.HasPrefix(ref, "Patient/") {
			doc.PatientID = strings.TrimPrefix(ref, "Patient/")
			f.populatePatientInfo(doc)
		}
	}

	return doc, nil
}

// populatePatientInfo fetches and populates patient information from the FHIR server.
func (f *fhirClient) populatePatientInfo(doc *search.Document) {
	url := fmt.Sprintf("%s/Patient/%s", f.base, doc.PatientID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(context.Background())
	resp, err := f.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	defer resp.Body.Close()

	var patient map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return
	}

	if nameList, ok := patient["name"].([]interface{}); ok && len(nameList) > 0 {
		if nameMap, ok := nameList[0].(map[string]interface{}); ok {
			given := ""
			if g, ok := nameMap["given"].([]interface{}); ok && len(g) > 0 {
				given, _ = g[0].(string)
			}
			family, _ := nameMap["family"].(string)
			doc.PatientName = strings.TrimSpace(given + " " + family)
		}
	}
	if gender, ok := patient["gender"].(string); ok {
		doc.Gender = gender
	}
}
