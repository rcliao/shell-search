package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

// DDGR implements Provider by shelling out to the ddgr CLI.
type DDGR struct{}

func (d *DDGR) Name() string { return "ddgr" }

func (d *DDGR) Search(ctx context.Context, opts Options) (*SearchResponse, error) {
	count := opts.Count
	if count <= 0 {
		count = 5
	}

	args := []string{"--json", "-n", strconv.Itoa(count)}

	if opts.Freshness != "" {
		// ddgr uses -t with d/w/m/y (not pd/pw/pm/py)
		f := opts.Freshness
		switch f {
		case "pd":
			f = "d"
		case "pw":
			f = "w"
		case "pm":
			f = "m"
		case "py":
			f = "y"
		}
		args = append(args, "-t", f)
	}

	args = append(args, opts.Query)

	cmd := exec.CommandContext(ctx, "ddgr", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ddgr command failed: %w", err)
	}

	var raw []ddgrResult
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("parsing ddgr output: %w", err)
	}

	sr := &SearchResponse{Query: opts.Query}
	for _, r := range raw {
		sr.Results = append(sr.Results, Result{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Abstract,
		})
	}
	return sr, nil
}

type ddgrResult struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Abstract string `json:"abstract"`
}
