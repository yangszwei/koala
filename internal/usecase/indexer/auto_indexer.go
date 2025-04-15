package indexer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yangszwei/go-micala/internal/infrastructure/datasource"
	"github.com/yangszwei/go-micala/internal/usecase/search"
)

// ScanPolicy defines the configuration for how and when a data source should be scanned for indexing.
type ScanPolicy struct {
	FullScanInterval time.Duration
	PageSize         int
	MaxPagesPerCycle int // kept for future partial scan support
}

// AutoIndexer manages scheduled background indexing of multiple data sources based on their respective scan policies.
type AutoIndexer struct {
	clients  map[string]datasource.Client
	svc      search.Service
	policies map[string]ScanPolicy
	state    map[string]*indexerState
	running  map[string]bool
	mu       sync.Mutex
}

// indexerState stores runtime information for a data source, such as when the last full scan occurred.
type indexerState struct {
	lastFullScan time.Time
}

// New returns an initialized AutoIndexer with default state and client mappings.
func New(svc search.Service) *AutoIndexer {
	return &AutoIndexer{
		clients:  make(map[string]datasource.Client),
		svc:      svc,
		policies: make(map[string]ScanPolicy),
		state:    make(map[string]*indexerState),
		running:  make(map[string]bool),
	}
}

// Register adds a client and its associated scan policy to the AutoIndexer.
func (ai *AutoIndexer) Register(client datasource.Client, policy ScanPolicy) {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	ai.clients[client.Name()] = client
	ai.policies[client.Name()] = policy
	ai.state[client.Name()] = &indexerState{}
}

// markRunning marks a client as currently running if it's not already.
// Returns true if marking was successful; false if the client is already running.
func (ai *AutoIndexer) markRunning(name string) bool {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	if ai.running[name] {
		return false
	}
	ai.running[name] = true
	return true
}

// markDone marks a client as no longer running after a scan is complete.
func (ai *AutoIndexer) markDone(name string) {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	ai.running[name] = false
}

// Start launches background goroutines that perform periodic indexing for each registered client.
func (ai *AutoIndexer) Start(ctx context.Context) {
	for name, client := range ai.clients {
		go ai.runClient(ctx, name, client)
	}
}

// runClient periodically attempts to trigger a scan for the given client, ensuring only one concurrent run.
func (ai *AutoIndexer) runClient(ctx context.Context, name string, client datasource.Client) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ai.markRunning(name) {
				go func() {
					defer ai.markDone(name)
					ai.runOnce(ctx, name, client)
				}()
			}
		}
	}
}

// runOnce checks if a full scan should be performed for a data source and executes it if necessary.
func (ai *AutoIndexer) runOnce(ctx context.Context, name string, client datasource.Client) {
	ai.mu.Lock()
	state := ai.state[name]
	policy := ai.policies[name]
	ai.mu.Unlock()

	now := time.Now()
	if now.Sub(state.lastFullScan) > policy.FullScanInterval {
		log.Printf("[%s] Running full scan...", name)
		ai.runScanAll(ctx, name, client)
		state.lastFullScan = now
	}
}

// runScanAll performs a full scan using either streaming or paginated retrieval, depending on client capabilities.
func (ai *AutoIndexer) runScanAll(ctx context.Context, name string, client datasource.Client) {
	if streamable, ok := client.(datasource.Streamer); ok {
		ai.streamSummaries(ctx, name, streamable)
	} else if pager, ok := client.(datasource.Pager); ok {
		ai.pageSummaries(ctx, name, client, pager)
	} else {
		log.Printf("[%s] Skipping: no paging or streaming method implemented", name)
	}
}

// streamSummaries fetches documents from a streaming data source and processes them for indexing.
func (ai *AutoIndexer) streamSummaries(ctx context.Context, name string, client datasource.Streamer) {
	log.Printf("[%s] Starting indexing...", name)

	ai.mu.Lock()
	policy := ai.policies[name]
	ai.mu.Unlock()

	stream, err := client.Stream(ctx, policy.PageSize)
	if err != nil {
		log.Printf("[%s] Stream failed: %v", name, err)
		return
	}

	ai.processSummaries(ctx, name, client, stream)

	log.Printf("[%s] Finished indexing", name)
}

// pageSummaries fetches documents using offset-based pagination and processes them for indexing.
func (ai *AutoIndexer) pageSummaries(ctx context.Context, name string, client datasource.Client, pager datasource.Pager) {
	log.Printf("[%s] Starting indexing...", name)

	ai.mu.Lock()
	policy := ai.policies[name]
	ai.mu.Unlock()

	page := 0
	backoff := time.Second

	stream := make(chan datasource.DataSummary, policy.PageSize)
	go func() {
		defer close(stream)
		for policy.MaxPagesPerCycle <= 0 || page < policy.MaxPagesPerCycle {
			summaries, err := pager.List(ctx, page*policy.PageSize, policy.PageSize)
			if err != nil {
				log.Printf("[%s] List failed: %v", name, err)
				time.Sleep(backoff)
				if backoff < 30*time.Second {
					backoff *= 2
				}
				continue
			}
			backoff = time.Second

			if len(summaries) == 0 {
				break
			}
			for _, s := range summaries {
				stream <- s
			}
			page++
		}
	}()

	ai.processSummaries(ctx, name, client, stream)

	log.Printf("[%s] Finished indexing", name)
}

// processSummaries reads document summaries from a channel and concurrently indexes documents that don't already exist.
func (ai *AutoIndexer) processSummaries(ctx context.Context, name string, client datasource.Client, summaries <-chan datasource.DataSummary) {
	const maxWorkers = 5
	const slowThreshold = 500 * time.Millisecond

	var wg sync.WaitGroup
	tasks := make(chan datasource.DataSummary, 100)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Each worker pulls summaries from the task queue, checks for existence, fetches full document, and indexes it.
			// If processing is slow, it uses an exponential backoff to avoid overwhelming the system.
			wait := 250 * time.Millisecond
			for summary := range tasks {
				docID := summary.DocID()

				log.Printf("[%s] Checking ID %s", name, docID)
				start := time.Now()

				exists, err := ai.svc.Exists(ctx, docID)
				if err != nil {
					log.Printf("[%s] Exists check failed for ID %s: %v", name, docID, err)
					continue
				}
				if exists {
					log.Printf("[%s] Scan found for ID %s", name, docID)
					continue
				}

				doc, err := client.Fetch(ctx, summary)
				if err != nil {
					log.Printf("[%s] Fetch failed for ID %s: %v", name, docID, err)
					continue
				}
				if err := ai.svc.Index(ctx, *doc); err != nil {
					log.Printf("[%s] Index failed for ID %s: %v", name, doc.ID, err)
				} else {
					log.Printf("[%s] Successfully indexed ID %s", name, doc.ID)
				}

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
	}

	for summary := range summaries {
		tasks <- summary
	}

	close(tasks)

	wg.Wait()
}
