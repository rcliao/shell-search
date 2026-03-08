package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Tavily implements Provider using the Tavily Search API.
type Tavily struct {
	APIKey string
	Client *http.Client
}

func (t *Tavily) Name() string { return "tavily" }

func (t *Tavily) Search(ctx context.Context, opts Options) (*SearchResponse, error) {
	depth := "basic"
	if opts.Count > 5 {
		depth = "advanced"
	}

	reqBody := tavilyRequest{
		APIKey:      t.APIKey,
		Query:       opts.Query,
		SearchDepth: depth,
		MaxResults:  opts.Count,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling tavily request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.tavily.com/search", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := t.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tavily search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tavily API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var raw tavilyResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding tavily response: %w", err)
	}

	sr := &SearchResponse{Query: opts.Query}
	for _, r := range raw.Results {
		sr.Results = append(sr.Results, Result{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Content,
		})
	}
	return sr, nil
}

type tavilyRequest struct {
	APIKey      string `json:"api_key"`
	Query       string `json:"query"`
	SearchDepth string `json:"search_depth"`
	MaxResults  int    `json:"max_results"`
}

type tavilyResponse struct {
	Results []tavilyResult `json:"results"`
}

type tavilyResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}
