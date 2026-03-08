package provider

import "context"

// Result represents a single search result.
type Result struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Age         string `json:"age,omitempty"`
	ExtraText   string `json:"extra_text,omitempty"`
}

// SearchResponse holds the full search response.
type SearchResponse struct {
	Query   string   `json:"query"`
	Results []Result `json:"results"`
}

// Options configures a search request.
type Options struct {
	Query     string
	Count     int
	Freshness string // pd=24h, pw=7d, pm=31d, py=1yr
	Country   string
}

// Provider is the interface that search backends implement.
type Provider interface {
	Name() string
	Search(ctx context.Context, opts Options) (*SearchResponse, error)
}
