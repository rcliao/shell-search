# shell-search

Agent-friendly web search CLI. Returns search results as markdown or JSON.

Part of the [Ghost in the Shell](https://github.com/rcliao?tab=repositories&q=shell-) ecosystem.

## Install

```bash
go install github.com/rcliao/shell-search/cmd/shell-search@latest
```

## Usage

```bash
shell-search "query"                        # Markdown output, auto-detect provider
shell-search "query" --json                 # JSON output
shell-search "query" -p brave -n 10         # Brave provider, 10 results
shell-search "query" -f pw                  # Results from past week
shell-search "query" -f pd --country us     # Past 24h, US results
```

## Provider Auto-Detection

Checks environment variables in order:
1. `BRAVE_SEARCH_API_KEY` → Brave Search
2. `TAVILY_API_KEY` → Tavily
3. Falls back to `ddgr` (no key needed, must be installed)

Override with `--provider` / `-p` flag.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--num` | `-n` | 5 | Number of results |
| `--freshness` | `-f` | | Time filter: `pd` (24h), `pw` (7d), `pm` (31d), `py` (1yr) |
| `--country` | | | Country code (e.g., `us`, `jp`) |
| `--provider` | `-p` | auto | Provider: `brave`, `tavily`, `ddgr` |
| `--timeout` | | 30 | Request timeout in seconds |
| `--json` | | false | Output as JSON |

## Build

```bash
make build    # Build binary
make test     # Run tests
make vet      # Run go vet
```

## License

MIT
