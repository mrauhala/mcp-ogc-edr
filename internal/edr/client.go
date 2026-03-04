// Package edr provides a client for OGC EDR (Environmental Data Retrieval) APIs.
package edr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mrauhala/mcp-ogc-edr/pkg/types"
)

// Client is an OGC EDR API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new EDR client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetLandingPage fetches the API landing page.
func (c *Client) GetLandingPage(ctx context.Context) (*types.LandingPage, error) {
	var result types.LandingPage
	if err := c.get(ctx, "/", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCollections fetches all available collections.
func (c *Client) GetCollections(ctx context.Context) (*types.Collections, error) {
	var result types.Collections
	if err := c.get(ctx, "/collections", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCollection fetches metadata for a specific collection.
func (c *Client) GetCollection(ctx context.Context, collectionID string) (*types.Collection, error) {
	var result types.Collection
	if err := c.get(ctx, fmt.Sprintf("/collections/%s", collectionID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetLocations fetches named locations for a collection.
func (c *Client) GetLocations(ctx context.Context, collectionID string) (*types.Locations, error) {
	var result types.Locations
	if err := c.get(ctx, fmt.Sprintf("/collections/%s/locations", collectionID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// QueryPosition fetches data at a point position.
// coords should be WKT POINT, e.g. "POINT(1.0 51.0)"
func (c *Client) QueryPosition(ctx context.Context, collectionID string, params types.EDRQueryParams) (json.RawMessage, error) {
	return c.queryData(ctx, fmt.Sprintf("/collections/%s/position", collectionID), params)
}

// QueryRadius fetches data within a radius of a point.
func (c *Client) QueryRadius(ctx context.Context, collectionID string, params types.EDRQueryParams) (json.RawMessage, error) {
	return c.queryData(ctx, fmt.Sprintf("/collections/%s/radius", collectionID), params)
}

// QueryArea fetches data within a polygon area.
// coords should be WKT POLYGON.
func (c *Client) QueryArea(ctx context.Context, collectionID string, params types.EDRQueryParams) (json.RawMessage, error) {
	return c.queryData(ctx, fmt.Sprintf("/collections/%s/area", collectionID), params)
}

// QueryCube fetches data within a bounding box cube.
// coords should be WKT POLYGON representing the bbox.
func (c *Client) QueryCube(ctx context.Context, collectionID string, params types.EDRQueryParams) (json.RawMessage, error) {
	return c.queryData(ctx, fmt.Sprintf("/collections/%s/cube", collectionID), params)
}

// QueryTrajectory fetches data along a trajectory.
// coords should be WKT LINESTRING or LINESTRINGZ.
func (c *Client) QueryTrajectory(ctx context.Context, collectionID string, params types.EDRQueryParams) (json.RawMessage, error) {
	return c.queryData(ctx, fmt.Sprintf("/collections/%s/trajectory", collectionID), params)
}

// QueryLocation fetches data for a named location.
func (c *Client) QueryLocation(ctx context.Context, collectionID, locationID string, params types.EDRQueryParams) (json.RawMessage, error) {
	return c.queryData(ctx, fmt.Sprintf("/collections/%s/locations/%s", collectionID, locationID), params)
}

// queryData is the common data query implementation.
func (c *Client) queryData(ctx context.Context, path string, params types.EDRQueryParams) (json.RawMessage, error) {
	q := url.Values{}
	if params.Coords != "" {
		q.Set("coords", params.Coords)
	}
	if params.Datetime != "" {
		q.Set("datetime", params.Datetime)
	}
	if params.ParameterName != "" {
		q.Set("parameter-name", params.ParameterName)
	}
	if params.CRS != "" {
		q.Set("crs", params.CRS)
	}
	if params.F != "" {
		q.Set("f", params.F)
	}
	if params.Z != "" {
		q.Set("z", params.Z)
	}
	if params.Within != "" {
		q.Set("within", params.Within)
	}
	if params.WithinUnits != "" {
		q.Set("within-units", params.WithinUnits)
	}

	var raw json.RawMessage
	if err := c.get(ctx, path, q, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// get performs a GET request to the EDR API and decodes the JSON response.
func (c *Client) get(ctx context.Context, path string, query url.Values, out interface{}) error {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	slog.Debug("EDR request", "url", u)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("EDR API error %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}
