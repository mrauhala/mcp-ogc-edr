// Package prompts provides MCP prompt templates for the FMI EDR server.
package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// Registry holds all MCP prompt definitions.
type Registry struct{}

// NewRegistry creates a new prompt registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// --- Prompt Definitions ---

// CurrentWeatherPrompt defines the current-weather prompt.
func (r *Registry) CurrentWeatherPrompt() mcp.Prompt {
	return mcp.NewPrompt(
		"current-weather",
		mcp.WithPromptDescription("Get current weather observations for a location in Finland from FMI stations"),
		mcp.WithArgument("location",
			mcp.ArgumentDescription("City or place name in Finland, e.g. 'Helsinki', 'Tampere', 'Oulu'"),
			mcp.RequiredArgument(),
		),
	)
}

// WeatherForecastPrompt defines the weather-forecast prompt.
func (r *Registry) WeatherForecastPrompt() mcp.Prompt {
	return mcp.NewPrompt(
		"weather-forecast",
		mcp.WithPromptDescription("Get a weather forecast for a location in Finland using ECMWF or HARMONIE model data"),
		mcp.WithArgument("location",
			mcp.ArgumentDescription("City or place name in Finland, e.g. 'Helsinki', 'Rovaniemi'"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("hours",
			mcp.ArgumentDescription("Forecast horizon in hours (default: 24, max: 72)"),
		),
	)
}

// ExploreCollectionsPrompt defines the explore-collections prompt.
func (r *Registry) ExploreCollectionsPrompt() mcp.Prompt {
	return mcp.NewPrompt(
		"explore-collections",
		mcp.WithPromptDescription("Discover and explore the available FMI weather and environmental datasets"),
	)
}

// MarineConditionsPrompt defines the marine-conditions prompt.
func (r *Registry) MarineConditionsPrompt() mcp.Prompt {
	return mcp.NewPrompt(
		"marine-conditions",
		mcp.WithPromptDescription("Get sea state, water temperature and water level data from FMI marine observations"),
		mcp.WithArgument("location",
			mcp.ArgumentDescription("Coastal location or sea area, e.g. 'Helsinki', 'Gulf of Finland', 'Turku'"),
		),
	)
}

// --- Prompt Handlers ---

// HandleCurrentWeather handles the current-weather prompt.
func (r *Registry) HandleCurrentWeather(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	location := req.Params.Arguments["location"]
	if location == "" {
		location = "Helsinki"
	}

	text := fmt.Sprintf(`You have access to FMI (Finnish Meteorological Institute) weather observation tools.

Find the current weather conditions at %s, Finland.

Steps:
1. Call query_position with:
   - collection_id: "opendata_hourly"
   - coords: WKT POINT for %s (look up approximate coordinates if needed; Helsinki is POINT(24.9384 60.1699), Tampere POINT(23.7610 61.4978), Oulu POINT(25.4667 65.0167), Turku POINT(22.2666 60.4518), Rovaniemi POINT(25.7333 66.5000))
   - parameter_name: "ta_pt1h_avg,ws_pt1h_avg,wg_pt1h_max,ri_pt1h_sum,rh_pt1h_avg,pa_pt1h_avg"
   - datetime: last 2 hours interval (e.g. 2024-01-01T10:00:00Z/2024-01-01T12:00:00Z)
2. Use the most recent hour's values
3. Present a concise weather summary with:
   - Temperature (°C)
   - Wind speed (m/s) and gusts (m/s)
   - Precipitation last hour (mm)
   - Relative humidity (%%)
   - Air pressure (hPa)`, location, location)

	return mcp.NewGetPromptResult(
		fmt.Sprintf("Current weather observations at %s", location),
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(text)),
		},
	), nil
}

// HandleWeatherForecast handles the weather-forecast prompt.
func (r *Registry) HandleWeatherForecast(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	location := req.Params.Arguments["location"]
	if location == "" {
		location = "Helsinki"
	}
	hours := req.Params.Arguments["hours"]
	if hours == "" {
		hours = "24"
	}

	text := fmt.Sprintf(`You have access to FMI weather forecast model data via EDR tools.

Provide a %s-hour weather forecast for %s, Finland using ECMWF model data.

Steps:
1. Call query_position with:
   - collection_id: "ecmwf"
   - coords: WKT POINT for %s (Helsinki: POINT(24.9384 60.1699), Tampere: POINT(23.7610 61.4978), Oulu: POINT(25.4667 65.0167), Turku: POINT(22.2666 60.4518), Rovaniemi: POINT(25.7333 66.5000))
   - parameter_name: "temperature,precipitation1h,windums,windvms,humidity,pressure"
   - datetime: from now for the next %s hours as an interval (e.g. 2024-01-01T12:00:00Z/2024-01-02T12:00:00Z)
   - f: "GeoJSON"
2. Parse the forecast time series
3. Present a readable forecast summary highlighting:
   - Temperature range (min/max °C)
   - Precipitation (total mm and timing)
   - Wind conditions
   - Any notable weather changes`, hours, location, location, hours)

	return mcp.NewGetPromptResult(
		fmt.Sprintf("%s-hour weather forecast for %s", hours, location),
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(text)),
		},
	), nil
}

// HandleExploreCollections handles the explore-collections prompt.
func (r *Registry) HandleExploreCollections(_ context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	text := `You have access to the FMI (Finnish Meteorological Institute) OGC EDR API with dozens of weather and environmental datasets.

Help the user discover what data is available by doing the following:

1. Call list_collections to retrieve all available datasets
2. Group the collections by theme and present them clearly:
   - **Forecast models**: ecmwf, harmonie_scandinavia_surface, harmonie_scandinavia_hybrid, harmonie_scandinavia_pressure, pal_skandinavia
   - **Observations**: opendata, opendata_hourly, opendata_daily, opendata_minute, hourly, monthly
   - **Marine**: hbm, nemo, wam, sealevel, opendata_mareograph, opendata_buoy, meripalvelut
   - **Air quality**: airquality, airquality_fmi, airquality_urban, enfuser_helsinki_metropolitan
   - **Climate scenarios**: scenario_a1b, scenario_a2, scenario_b1, scenario_b2, scenario_1km
   - **Specialised**: flash, magneto, solar, external_radiation, air_radionuclides, road, sounding

3. Ask the user which dataset they are interested in
4. Call get_collection for their chosen collection to show available parameters and query types
5. Offer to run a sample query to demonstrate the data`

	return mcp.NewGetPromptResult(
		"Explore FMI EDR weather and environmental datasets",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(text)),
		},
	), nil
}

// HandleMarineConditions handles the marine-conditions prompt.
func (r *Registry) HandleMarineConditions(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	location := req.Params.Arguments["location"]

	var locationPart string
	if location != "" {
		locationPart = fmt.Sprintf(" near %s", location)
	}

	text := fmt.Sprintf(`You have access to FMI marine and sea observation data via EDR tools.

Retrieve current marine conditions%s from Finnish waters.

Steps:
1. For sea water temperature and currents, call query_position with:
   - collection_id: "nemo" (Baltic Sea model) or "hbm"
   - coords: appropriate POINT coordinates%s
   - f: "GeoJSON"

2. For water level / sea level, call query_position or query_location with:
   - collection_id: "sealevel" or "opendata_mareograph"
   - Use get_locations first to find the nearest tide gauge station

3. For wave data, use:
   - collection_id: "wam"
   - parameter_name relevant wave height and period parameters

4. Summarise the marine conditions:
   - Sea surface temperature (°C)
   - Wave height (m) and period (s)
   - Water level relative to mean sea level (cm)
   - Any notable conditions`, locationPart, locationPart)

	return mcp.NewGetPromptResult(
		"Marine conditions from FMI",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(text)),
		},
	), nil
}
