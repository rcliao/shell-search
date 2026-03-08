package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rcliao/shell-search/internal/provider"
)

// Markdown formats a SearchResponse as agent-friendly markdown.
func Markdown(resp *provider.SearchResponse) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## Search Results for %q\n\n", resp.Query)

	if len(resp.Results) == 0 {
		b.WriteString("No results found.\n")
		return b.String()
	}

	for i, r := range resp.Results {
		fmt.Fprintf(&b, "### %d. %s\n", i+1, r.Title)
		fmt.Fprintf(&b, "**URL:** %s\n", r.URL)
		if r.Age != "" {
			fmt.Fprintf(&b, "**Age:** %s\n", r.Age)
		}
		b.WriteString("\n")
		b.WriteString(r.Description)
		b.WriteString("\n")
		if r.ExtraText != "" {
			b.WriteString("\n")
			b.WriteString(r.ExtraText)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("---\n\n**Sources:**\n")
	for _, r := range resp.Results {
		fmt.Fprintf(&b, "- %s\n", r.URL)
	}

	return b.String()
}

// JSON formats a SearchResponse as JSON.
func JSON(resp *provider.SearchResponse) (string, error) {
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling JSON: %w", err)
	}
	return string(data), nil
}
