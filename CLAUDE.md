# shell-search

Agent-friendly web search CLI. Returns search results as markdown (default) or JSON.

## Architecture

- `cmd/shell-search/main.go` — Cobra CLI entrypoint
- `internal/provider/` — Search backend implementations (Brave, Tavily, ddgr)
- `internal/formatter/` — Output formatting (Markdown, JSON)

## Provider Auto-Detection

Checks env vars in order: `BRAVE_SEARCH_API_KEY` > `TAVILY_API_KEY` > ddgr (no key needed).

## Build & Test

```bash
make build    # Build binary
make test     # Run tests
make vet      # Run go vet
```

## Usage

```bash
shell-search "query"                       # Markdown output, auto-detect provider
shell-search "query" --json                # JSON output
shell-search "query" -p brave -n 10       # Brave provider, 10 results
shell-search "query" -f pw                # Results from past week
```
