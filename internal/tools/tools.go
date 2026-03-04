// Package tools provides MCP tool implementations that proxy OGC EDR API calls.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mrauhala/mcp-ogc-edr/internal/edr"
	"github.com/mrauhala/mcp-ogc-edr/pkg/types"
)

// Registry holds all MCP tools and the EDR client.
type Registry struct {
	client *edr.Client
}

// NewRegistry creates a new tool registry.
func NewRegistry(client *edr.Client) *Registry {
	return &Registry{client: client}
}

// --- Tool Definitions ---

// ListCollectionsTool returns the tool definition for listing EDR collections.
func (r *Registry) ListCollectionsTool() mcp.Tool {
	return mcp.NewTool(
		"list_collections",
		mcp.WithDescription("List all available OGC EDR collections (datasets). Returns collection IDs, titles, descriptions, spatial/temporal extents, and available query types."),
	)
}

// GetCollectionTool returns the tool definition for getting collection metadata.
func (r *Registry) GetCollectionTool() mcp.Tool {
	return mcp.NewTool(
		"get_collection",
		mcp.WithDescription("Get detailed metadata for a specific OGC EDR collection, including available parameters, output formats, and supported query types."),
		mcp.WithString("collection_id",
			mcp.Required(),
			mcp.Description("The ID of the collection to retrieve"),
		),
	)
}

// GetLocationsTool returns the tool definition for listing named locations.
func (r *Registry) GetLocationsTool() mcp.Tool {
	return mcp.NewTool(
		"get_locations",
		mcp.WithDescription("List named locations (stations, sites, etc.) available in an EDR collection."),
		mcp.WithString("collection_id",
			mcp.Required(),
			mcp.Description("The ID of the collection"),
		),
	)
}

// QueryPositionTool returns the tool for point position queries.
func (r *Registry) QueryPositionTool() mcp.Tool {
	return mcp.NewTool(
		"query_position",
		mcp.WithDescription("Query EDR data at a geographic point (position query). Returns observational or forecast data at a specific coordinate."),
		mcp.WithString("collection_id",
			mcp.Required(),
			mcp.Description("The ID of the EDR collection to query"),
		),
		mcp.WithString("coords",
			mcp.Required(),
			mcp.Description("WKT POINT geometry, e.g. 'POINT(lon lat)' or 'POINT(-1.5 52.0)'"),
		),
		mcp.WithString("datetime",
			mcp.Description("Datetime filter in ISO 8601. Single instant or interval: '2024-01-01T00:00:00Z' or '2024-01-01T00:00:00Z/2024-01-02T00:00:00Z'"),
		),
		mcp.WithString("parameter_name",
			mcp.Description("Comma-separated list of parameter names to include, e.g. 'temperature,wind_speed'"),
		),
		mcp.WithString("crs",
			mcp.Description("Coordinate reference system for the response, e.g. 'CRS84'"),
		),
		mcp.WithString("z",
			mcp.Description("Vertical level(s), e.g. '500' or '100/500' or 'R5/100/100'"),
		),
		mcp.WithString("f",
			mcp.Description("Output format, e.g. 'CoverageJSON' or 'GeoJSON'"),
		),
	)
}

// QueryRadiusTool returns the tool for radius queries.
func (r *Registry) QueryRadiusTool() mcp.Tool {
	return mcp.NewTool(
		"query_radius",
		mcp.WithDescription("Query EDR data within a radius of a geographic point."),
		mcp.WithString("collection_id",
			mcp.Required(),
			mcp.Description("The ID of the EDR collection to query"),
		),
		mcp.WithString("coords",
			mcp.Required(),
			mcp.Description("WKT POINT geometry, e.g. 'POINT(-1.5 52.0)'"),
		),
		mcp.WithString("within",
			mcp.Required(),
			mcp.Description("Radius distance, e.g. '50'"),
		),
		mcp.WithString("within_units",
			mcp.Required(),
			mcp.Description("Units for radius: 'km' or 'mi'"),
		),
		mcp.WithString("datetime",
			mcp.Description("Datetime filter in ISO 8601"),
		),
		mcp.WithString("parameter_name",
			mcp.Description("Comma-separated parameter names"),
		),
		mcp.WithString("f",
			mcp.Description("Output format"),
		),
	)
}

// QueryAreaTool returns the tool for area queries.
func (r *Registry) QueryAreaTool() mcp.Tool {
	return mcp.NewTool(
		"query_area",
		mcp.WithDescription("Query EDR data within a polygon area."),
		mcp.WithString("collection_id",
			mcp.Required(),
			mcp.Description("The ID of the EDR collection to query"),
		),
		mcp.WithString("coords",
			mcp.Required(),
			mcp.Description("WKT POLYGON geometry, e.g. 'POLYGON((-2 51,-2 52,0 52,0 51,-2 51))'"),
		),
		mcp.WithString("datetime",
			mcp.Description("Datetime filter in ISO 8601"),
		),
		mcp.WithString("parameter_name",
			mcp.Description("Comma-separated parameter names"),
		),
		mcp.WithString("z",
			mcp.Description("Vertical level(s)"),
		),
		mcp.WithString("f",
			mcp.Description("Output format"),
		),
	)
}

// QueryLocationTool returns the tool for named location queries.
func (r *Registry) QueryLocationTool() mcp.Tool {
	return mcp.NewTool(
		"query_location",
		mcp.WithDescription("Query EDR data for a specific named location (e.g. weather station, observation site)."),
		mcp.WithString("collection_id",
			mcp.Required(),
			mcp.Description("The ID of the EDR collection to query"),
		),
		mcp.WithString("location_id",
			mcp.Required(),
			mcp.Description("The ID of the named location"),
		),
		mcp.WithString("datetime",
			mcp.Description("Datetime filter in ISO 8601"),
		),
		mcp.WithString("parameter_name",
			mcp.Description("Comma-separated parameter names"),
		),
		mcp.WithString("f",
			mcp.Description("Output format"),
		),
	)
}

// QueryTrajectoryTool returns the tool for trajectory queries.
func (r *Registry) QueryTrajectoryTool() mcp.Tool {
	return mcp.NewTool(
		"query_trajectory",
		mcp.WithDescription("Query EDR data along a trajectory (LINESTRING path)."),
		mcp.WithString("collection_id",
			mcp.Required(),
			mcp.Description("The ID of the EDR collection to query"),
		),
		mcp.WithString("coords",
			mcp.Required(),
			mcp.Description("WKT LINESTRING geometry, e.g. 'LINESTRING(-2 51,-1 51.5,0 52)'"),
		),
		mcp.WithString("datetime",
			mcp.Description("Datetime filter in ISO 8601"),
		),
		mcp.WithString("parameter_name",
			mcp.Description("Comma-separated parameter names"),
		),
		mcp.WithString("f",
			mcp.Description("Output format"),
		),
	)
}

// --- Tool Handlers ---

// requireString extracts a required string argument from a tool request.
func requireString(req mcp.CallToolRequest, key string) (string, error) {
	v, ok := req.Params.Arguments[key]
	if !ok {
		return "", fmt.Errorf("missing required argument %q", key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("argument %q must be a string", key)
	}
	return s, nil
}

// HandleListCollections handles the list_collections tool call.
func (r *Registry) HandleListCollections(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collections, err := r.client.GetCollections(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list collections: %v", err)), nil
	}

	// Return a summarized, LLM-friendly version
	type summary struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
		QueryTypes  []string `json:"query_types,omitempty"`
		Parameters  []string `json:"parameters,omitempty"`
	}

	summaries := make([]summary, 0, len(collections.Collections))
	for _, c := range collections.Collections {
		s := summary{
			ID:          c.ID,
			Title:       c.Title,
			Description: c.Description,
			QueryTypes:  queryTypes(c.DataQueries),
		}
		for k := range c.ParameterNames {
			s.Parameters = append(s.Parameters, k)
		}
		summaries = append(summaries, s)
	}

	return jsonResult(summaries)
}

// HandleGetCollection handles the get_collection tool call.
func (r *Registry) HandleGetCollection(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionID, err := requireString(req,"collection_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	col, err := r.client.GetCollection(ctx, collectionID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get collection %q: %v", collectionID, err)), nil
	}

	return jsonResult(col)
}

// HandleGetLocations handles the get_locations tool call.
func (r *Registry) HandleGetLocations(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionID, err := requireString(req,"collection_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	locs, err := r.client.GetLocations(ctx, collectionID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get locations: %v", err)), nil
	}

	return jsonResult(locs)
}

// HandleQueryPosition handles the query_position tool call.
func (r *Registry) HandleQueryPosition(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionID, err := requireString(req,"collection_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	params, err := extractQueryParams(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, err := r.client.QueryPosition(ctx, collectionID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Position query failed: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleQueryRadius handles the query_radius tool call.
func (r *Registry) HandleQueryRadius(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionID, err := requireString(req,"collection_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	params, err := extractQueryParams(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	params.Within, _ = requireString(req,"within")
	params.WithinUnits, _ = requireString(req,"within_units")

	data, err := r.client.QueryRadius(ctx, collectionID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Radius query failed: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleQueryArea handles the query_area tool call.
func (r *Registry) HandleQueryArea(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionID, err := requireString(req,"collection_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	params, err := extractQueryParams(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, err := r.client.QueryArea(ctx, collectionID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Area query failed: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleQueryLocation handles the query_location tool call.
func (r *Registry) HandleQueryLocation(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionID, err := requireString(req,"collection_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	locationID, err := requireString(req,"location_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	params, err := extractQueryParams(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, err := r.client.QueryLocation(ctx, collectionID, locationID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Location query failed: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleQueryTrajectory handles the query_trajectory tool call.
func (r *Registry) HandleQueryTrajectory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionID, err := requireString(req,"collection_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	params, err := extractQueryParams(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, err := r.client.QueryTrajectory(ctx, collectionID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Trajectory query failed: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// --- Helpers ---

func extractQueryParams(req mcp.CallToolRequest) (types.EDRQueryParams, error) {
	p := types.EDRQueryParams{}
	p.Coords, _ = requireString(req,"coords")
	// Optional fields — ignore errors (they won't be set if missing)
	if v, ok := req.Params.Arguments["datetime"].(string); ok {
		p.Datetime = v
	}
	if v, ok := req.Params.Arguments["parameter_name"].(string); ok {
		p.ParameterName = v
	}
	if v, ok := req.Params.Arguments["crs"].(string); ok {
		p.CRS = v
	}
	if v, ok := req.Params.Arguments["z"].(string); ok {
		p.Z = v
	}
	if v, ok := req.Params.Arguments["f"].(string); ok {
		p.F = v
	}
	return p, nil
}

func jsonResult(v interface{}) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}

func queryTypes(dq types.DataQueryLinks) []string {
	var qt []string
	if dq.Position != nil {
		qt = append(qt, "position")
	}
	if dq.Radius != nil {
		qt = append(qt, "radius")
	}
	if dq.Area != nil {
		qt = append(qt, "area")
	}
	if dq.Cube != nil {
		qt = append(qt, "cube")
	}
	if dq.Trajectory != nil {
		qt = append(qt, "trajectory")
	}
	if dq.Corridor != nil {
		qt = append(qt, "corridor")
	}
	if dq.Locations != nil {
		qt = append(qt, "locations")
	}
	if dq.Items != nil {
		qt = append(qt, "items")
	}
	return qt
}
