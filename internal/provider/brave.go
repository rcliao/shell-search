package provider

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Brave implements Provider using the Brave Search API.
type Brave struct {
	APIKey string
	Client *http.Client
}

func (b *Brave) Name() string { return "brave" }

func (b *Brave) Search(ctx context.Context, opts Options) (*SearchResponse, error) {
	u, _ := url.Parse("https://api.search.brave.com/res/v1/web/search")
	q := u.Query()
	q.Set("q", opts.Query)
	if opts.Count > 0 {
		q.Set("count", strconv.Itoa(opts.Count))
	}
	if opts.Freshness != "" {
		q.Set("freshness", opts.Freshness)
	}
	if opts.Country != "" {
		q.Set("country", opts.Country)
	}
	q.Set("extra_snippets", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("X-Subscription-Token", b.APIKey)

	client := b.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("brave search request: %w", err)
	}
	defer resp.Body.Close()

	// Decompress gzip if needed.
	var body io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip reader: %w", err)
		}
		defer gr.Close()
		body = gr
	}

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(body)
		return nil, fmt.Errorf("brave API error (status %d): %s", resp.StatusCode, string(data))
	}

	var raw braveResponse
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding brave response: %w", err)
	}

	sr := &SearchResponse{Query: opts.Query}
	for _, r := range raw.Web.Results {
		extra := ""
		for _, s := range r.ExtraSnippets {
			if extra != "" {
				extra += "\n"
			}
			extra += s
		}
		sr.Results = append(sr.Results, Result{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Description,
			Age:         r.Age,
			ExtraText:   extra,
		})
	}
	return sr, nil
}

// Brave API response types (subset of fields we care about).

type braveResponse struct {
	Web braveWebResults `json:"web"`
}

type braveWebResults struct {
	Results []braveWebResult `json:"results"`
}

type braveWebResult struct {
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	Description   string   `json:"description"`
	Age           string   `json:"age"`
	ExtraSnippets []string `json:"extra_snippets"`
}
