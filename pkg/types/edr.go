package types

import "encoding/json"

// I18nString unmarshals either a plain string or a {"lang": "value"} i18n map.
// The OGC EDR spec allows both forms for description/label fields.
type I18nString map[string]string

func (s *I18nString) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		*s = I18nString{"": str}
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*s = I18nString(m)
	return nil
}

// LandingPage is the OGC EDR API root response.
type LandingPage struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Links       []Link  `json:"links"`
}

// Link is an OGC API link object.
type Link struct {
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
}

// Collections is the /collections response.
type Collections struct {
	Links       []Link       `json:"links"`
	Collections []Collection `json:"collections"`
}

// Collection represents a single EDR collection (dataset).
type Collection struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description,omitempty"`
	Links       []Link           `json:"links,omitempty"`
	Extent      *Extent          `json:"extent,omitempty"`
	DataQueries DataQueryLinks   `json:"data_queries,omitempty"`
	CRS         []string         `json:"crs,omitempty"`
	OutputFormats []string       `json:"output_formats,omitempty"`
	ParameterNames map[string]Parameter `json:"parameter_names,omitempty"`
}

// Extent describes the spatial and temporal extent.
type Extent struct {
	Spatial  *SpatialExtent  `json:"spatial,omitempty"`
	Temporal *TemporalExtent `json:"temporal,omitempty"`
	Vertical *VerticalExtent `json:"vertical,omitempty"`
}

// SpatialExtent holds the bounding box.
type SpatialExtent struct {
	Bbox [][]float64 `json:"bbox"`
	CRS  string      `json:"crs,omitempty"`
}

// TemporalExtent holds time intervals.
type TemporalExtent struct {
	Interval [][]string `json:"interval"`
	TRS      string     `json:"trs,omitempty"`
}

// VerticalExtent holds vertical levels.
type VerticalExtent struct {
	Interval [][]string `json:"interval"`
	VRS      string     `json:"vrs,omitempty"`
	Values   []string   `json:"values,omitempty"`
}

// DataQueryLinks holds links to each EDR query type.
type DataQueryLinks struct {
	Position  *DataQueryLink `json:"position,omitempty"`
	Radius    *DataQueryLink `json:"radius,omitempty"`
	Area      *DataQueryLink `json:"area,omitempty"`
	Cube      *DataQueryLink `json:"cube,omitempty"`
	Trajectory *DataQueryLink `json:"trajectory,omitempty"`
	Corridor  *DataQueryLink `json:"corridor,omitempty"`
	Items     *DataQueryLink `json:"items,omitempty"`
	Locations *DataQueryLink `json:"locations,omitempty"`
}

// DataQueryLink is a link to an EDR query endpoint.
type DataQueryLink struct {
	Link  Link         `json:"link"`
	Variables QueryVars `json:"variables,omitempty"`
}

// QueryVars describes supported variables for a query type.
type QueryVars struct {
	QueryType     string   `json:"query_type"`
	OutputFormats []string `json:"output_formats,omitempty"`
	CRSDetails    []CRSDetail `json:"crs_details,omitempty"`
}

// CRSDetail is a CRS descriptor.
type CRSDetail struct {
	CRS  string `json:"crs"`
	WKT  string `json:"wkt,omitempty"`
}

// Parameter is an observable parameter (variable) in a collection.
type Parameter struct {
	ID               string       `json:"id,omitempty"`
	Type             string       `json:"type"`
	Label            I18nString   `json:"label,omitempty"`
	Description      I18nString   `json:"description,omitempty"`
	Unit             *Unit        `json:"unit,omitempty"`
	ObservedProperty ObservedProp `json:"observedProperty"`
	MeasurementType  *MeasType    `json:"measurement_type,omitempty"`
}

// Unit describes the unit of measurement.
type Unit struct {
	Label  I18nString  `json:"label,omitempty"`
	Symbol *UnitSymbol `json:"symbol,omitempty"`
}

// UnitSymbol is the unit symbol value and type.
type UnitSymbol struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

// ObservedProp is the observed property.
type ObservedProp struct {
	ID    string     `json:"id,omitempty"`
	Label I18nString `json:"label"`
}

// MeasType is the measurement type.
type MeasType struct {
	Method   string `json:"method"`
	Duration string `json:"duration,omitempty"`
}

// Locations is the /collections/{id}/locations response.
type Locations struct {
	Type     string     `json:"type"` // FeatureCollection
	Features []Location `json:"features"`
	Links    []Link     `json:"links,omitempty"`
}

// Location is a GeoJSON Feature for a named location.
type Location struct {
	Type       string                 `json:"type"` // Feature
	ID         string                 `json:"id"`
	Geometry   interface{}            `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// EDRQueryParams holds common query parameters for EDR data queries.
type EDRQueryParams struct {
	// Coordinate string (WKT POINT, POLYGON, LINESTRING, etc.)
	Coords string
	// Datetime filter: e.g. "2024-01-01T00:00:00Z/2024-01-02T00:00:00Z"
	Datetime string
	// Parameter names to include (comma-separated)
	ParameterName string
	// CRS for the response
	CRS string
	// Output format (e.g. "CoverageJSON", "GeoJSON")
	F string
	// Vertical level
	Z string
	// Radius query specific
	Within     string
	WithinUnits string
}
