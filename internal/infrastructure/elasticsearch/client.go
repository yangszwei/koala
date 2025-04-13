package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
)

// Client wraps the official Elasticsearch Go client to simplify interaction with the cluster.
type Client struct {
	Client *elasticsearch.Client
}

// NewClient initializes and returns a new Client instance connected to the given address.
// It returns an error if the connection cannot be established.
func NewClient(addr string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{addr},
	}

	// Create a new Elasticsearch client
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Check if the client is working
	if _, err = client.Ping(); err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil
}
