# mcp-ogc-edr

A Model Context Protocol (MCP) server that exposes OGC EDR (Environmental Data Retrieval) APIs as MCP tools, enabling LLMs to query geospatial and environmental datasets.

## Architecture

```
LLM / MCP Client
     │
     ▼
MCP Server (stdio or SSE)
     │
     ▼
OGC EDR API Backend
(weather, climate, oceanography, etc.)
```

## Tools Exposed

| Tool | Description |
|---|---|
| `list_collections` | Discover all available datasets |
| `get_collection` | Get metadata, parameters, query types for a dataset |
| `get_locations` | List named locations (stations, sites) |
| `query_position` | Data at a geographic point |
| `query_radius` | Data within radius of a point |
| `query_area` | Data within a polygon |
| `query_trajectory` | Data along a linestring path |
| `query_location` | Data at a named location |

## Quick Start

```bash
# Install dependencies
go mod tidy

# Build
make build

# Run against a real EDR server (stdio transport for MCP clients)
EDR_BASE_URL=https://your-edr-server.com/edr ./bin/mcp-ogc-edr

# Run as SSE server (for web-based MCP clients)
EDR_BASE_URL=https://your-edr-server.com/edr \
MCP_TRANSPORT=sse \
SSE_ADDR=:8080 \
./bin/mcp-ogc-edr
```

## Configuration

| Env Var | Flag | Default | Description |
|---|---|---|---|
| `EDR_BASE_URL` | `-edr-url` | (required) | OGC EDR API base URL |
| `MCP_TRANSPORT` | `-transport` | `stdio` | `stdio` or `sse` |
| `SSE_ADDR` | `-sse-addr` | `:8080` | SSE listen address |
| `LOG_LEVEL` | `-log-level` | `info` | `debug`, `info`, `warn`, `error` |

## Claude Desktop Configuration

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "ogc-edr": {
      "command": "/path/to/mcp-ogc-edr",
      "env": {
        "EDR_BASE_URL": "https://your-edr-server.com/edr"
      }
    }
  }
}
```

## Public EDR Servers to Test With

- **NOAA Big Data**: `https://edr.ioos.us/api`
- **UK Met Office**: `https://api.meteomatics.com` (auth required)
- **OGC Reference**: `https://ogc.api.weather.gc.ca/` (Environment Canada)
- **USGS**: various at `https://labs.waterdata.usgs.gov`

## Extending

To add a new tool:

1. Add a method pair in `internal/tools/tools.go`:
   - `MyNewTool() mcp.Tool` — defines the tool schema
   - `HandleMyNewTool(ctx, req) (*mcp.CallToolResult, error)` — implements it

2. Register in `internal/mcp/server.go`:
   ```go
   s.mcp.AddTool(r.MyNewTool(), r.HandleMyNewTool)
   ```

3. Add EDR client method in `internal/edr/client.go` if needed.

## Project Structure

```
.
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── edr/client.go           # OGC EDR HTTP client
│   ├── mcp/server.go           # MCP server wiring
│   └── tools/tools.go          # MCP tool definitions & handlers
├── pkg/types/edr.go            # OGC EDR type definitions
├── Dockerfile
├── Makefile
└── go.mod
```
