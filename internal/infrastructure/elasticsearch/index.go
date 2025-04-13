package elasticsearch

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed indices/*.json
var indexFS embed.FS

// EnsureIndices loads embedded index definition files and ensures each index exists in Elasticsearch.
// If an index does not exist, it will be created based on its corresponding JSON definition.
func (c *Client) EnsureIndices() error {
	entries, err := indexFS.ReadDir("indices")
	if err != nil {
		return fmt.Errorf("failed to read embedded indices: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".json")

		// Check if index already exists
		exists, err := c.Client.Indices.Exists([]string{name})
		if err != nil {
			return fmt.Errorf("checking if index %q exists: %w", name, err)
		}
		if exists.StatusCode == 200 {
			continue
		}

		defBytes, err := indexFS.ReadFile("indices/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read index definition %q: %w", entry.Name(), err)
		}

		var def map[string]interface{}
		if err := json.Unmarshal(defBytes, &def); err != nil {
			return fmt.Errorf("invalid JSON in %q: %w", entry.Name(), err)
		}

		if err := c.CreateIndex(name, def); err != nil {
			return fmt.Errorf("creating index %q: %w", name, err)
		}
	}

	return nil
}

// CreateIndex creates a new index in Elasticsearch with the specified name and settings.
func (c *Client) CreateIndex(name string, body map[string]interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	res, err := c.Client.Indices.Create(name, c.Client.Indices.Create.WithBody(bytes.NewReader(data)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error creating index %s: %s", name, res.String())
	}
	return nil
}
