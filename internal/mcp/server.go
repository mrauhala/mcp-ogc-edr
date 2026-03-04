// Package mcp wires together the MCP server with EDR tools.
package mcp

import (
	"context"
	"fmt"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mrauhala/mcp-ogc-edr/internal/edr"
	"github.com/mrauhala/mcp-ogc-edr/internal/prompts"
	"github.com/mrauhala/mcp-ogc-edr/internal/tools"
)

const (
	serverName    = "mcp-ogc-edr"
	serverVersion = "0.1.0"
)

// Server wraps the MCP server and EDR tools.
type Server struct {
	mcp      *server.MCPServer
	registry *tools.Registry
	prompts  *prompts.Registry
}

// NewServer creates and configures the MCP server.
func NewServer(edrBaseURL string) (*Server, error) {
	edrClient := edr.NewClient(edrBaseURL)
	registry := tools.NewRegistry(edrClient)

	s := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	srv := &Server{
		mcp:      s,
		registry: registry,
		prompts:  prompts.NewRegistry(),
	}

	srv.registerTools()
	srv.registerPrompts()
	srv.registerResources(edrBaseURL)

	return srv, nil
}

// registerTools wires all EDR tools to the MCP server.
func (s *Server) registerTools() {
	r := s.registry

	s.mcp.AddTool(r.ListCollectionsTool(), r.HandleListCollections)
	s.mcp.AddTool(r.GetCollectionTool(), r.HandleGetCollection)
	s.mcp.AddTool(r.GetLocationsTool(), r.HandleGetLocations)
	s.mcp.AddTool(r.QueryPositionTool(), r.HandleQueryPosition)
	s.mcp.AddTool(r.QueryRadiusTool(), r.HandleQueryRadius)
	s.mcp.AddTool(r.QueryAreaTool(), r.HandleQueryArea)
	s.mcp.AddTool(r.QueryLocationTool(), r.HandleQueryLocation)
	s.mcp.AddTool(r.QueryTrajectoryTool(), r.HandleQueryTrajectory)
}

// registerPrompts wires all EDR prompt templates to the MCP server.
func (s *Server) registerPrompts() {
	p := s.prompts
	s.mcp.AddPrompt(p.CurrentWeatherPrompt(), p.HandleCurrentWeather)
	s.mcp.AddPrompt(p.WeatherForecastPrompt(), p.HandleWeatherForecast)
	s.mcp.AddPrompt(p.ExploreCollectionsPrompt(), p.HandleExploreCollections)
	s.mcp.AddPrompt(p.MarineConditionsPrompt(), p.HandleMarineConditions)
}

// registerResources registers static MCP resources describing the EDR API.
func (s *Server) registerResources(edrBaseURL string) {
	s.mcp.AddResource(
		mcpgo.NewResource(
			"edr://api-info",
			"OGC EDR API Information",
			mcpgo.WithResourceDescription("Information about the connected OGC EDR backend, query types, and usage patterns"),
			mcpgo.WithMIMEType("text/plain"),
		),
		func(ctx context.Context, req mcpgo.ReadResourceRequest) ([]mcpgo.ResourceContents, error) {
			content := fmt.Sprintf(`OGC EDR (Environmental Data Retrieval) API
==========================================
Backend URL: %s

Supported Query Types
---------------------
- position  : Data at a geographic point (WKT POINT)
- radius    : Data within a radius of a point
- area      : Data within a polygon (WKT POLYGON)
- trajectory: Data along a path (WKT LINESTRING)
- location  : Data at a named location (station/site)

Common Parameters
-----------------
- coords         : WKT geometry string
- datetime       : ISO 8601 instant or interval (e.g. 2024-01-01T00:00:00Z/2024-01-02T00:00:00Z)
- parameter-name : Comma-separated list of variables (e.g. temperature,wind_speed)
- crs            : Coordinate reference system (default: CRS84)
- z              : Vertical level or range
- f              : Output format (CoverageJSON, GeoJSON, etc.)

Workflow
--------
1. Use list_collections to discover available datasets
2. Use get_collection to see parameters and query types for a dataset
3. Use get_locations to find named stations/sites (if supported)
4. Query data using query_position, query_area, query_location, etc.

WKT Examples
------------
Point:      POINT(-1.5 52.0)
Polygon:    POLYGON((-2 51,-2 52,0 52,0 51,-2 51))
LineString: LINESTRING(-2 51,-1 51.5,0 52)
`, edrBaseURL)

			return []mcpgo.ResourceContents{
				mcpgo.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "text/plain",
					Text:     content,
				},
			}, nil
		},
	)
}

// Run starts the MCP server with the specified transport.
func (s *Server) Run(ctx context.Context, transport, sseAddr string) error {
	switch transport {
	case "stdio":
		return server.ServeStdio(s.mcp)
	case "sse":
		sseServer := server.NewSSEServer(s.mcp,
			server.WithBaseURL(fmt.Sprintf("http://%s", sseAddr)),
		)
		return sseServer.Start(sseAddr)
	case "streamable-http":
		httpServer := server.NewStreamableHTTPServer(s.mcp)
		return httpServer.Start(sseAddr)
	default:
		return fmt.Errorf("unknown transport %q: must be 'stdio', 'sse', or 'streamable-http'", transport)
	}
}
