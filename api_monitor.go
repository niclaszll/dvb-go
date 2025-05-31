package dvb

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MonitorStopParams contains the parameters for monitoring departures from a specific stop.
// This is used to get real-time departure information for public transport vehicles.
type MonitorStopParams struct {
	// StopId is the unique identifier for the stop to monitor. This is required and cannot be empty.
	// Use the GetPoint API to find stop IDs based on stop names or locations.
	StopId string

	// Format specifies the response format. Optional parameter.
	// Supported values depend on the DVB API implementation.
	Format *string

	// Time specifies the time for which to get departures. Optional parameter.
	// Format should be compatible with the DVB API time format.
	// If not specified, uses the current time.
	Time *string

	// IsArrival when set to true, shows arrivals instead of departures.
	// When false or nil, shows departures (default behavior).
	IsArrival *bool

	// Limit restricts the maximum number of departures/arrivals returned.
	// Must be a positive integer. If nil or 0, uses the API's default limit.
	Limit *int

	// ShortTermChanges when set to true, includes short-term changes like delays or cancellations.
	// When false or nil, may not include the most recent real-time updates.
	ShortTermChanges *bool

	// MentzOnly when set to true, includes only data from the Mentz system.
	// When false or nil, includes data from all available systems.
	MentzOnly *bool
}

// MonitorStopResponse represents the response from the DVB stop monitoring API.
// It contains real-time departure/arrival information for a specific stop.
type MonitorStopResponse struct {
	// Name is the official name of the monitored stop
	Name string `json:"Name"`

	// Status contains the API response status including error codes and messages
	Status Status `json:"Status"`

	// Place indicates the city or area where the stop is located
	Place string `json:"Place"`

	// ExpirationTime indicates when this response data expires and should be refreshed
	ExpirationTime string `json:"ExpirationTime"`

	// Departures is an array of upcoming departures/arrivals from this stop
	Departures []Departure `json:"Departures"`
}

// Departure represents a single departure or arrival at a monitored stop.
// It contains detailed information about the vehicle, timing, and any disruptions.
type Departure struct {
	// Id is the unique identifier for this departure
	Id string `json:"Id"`

	// DlId is the DVB line identifier
	DlId string `json:"DlId"`

	// LineName is the display name of the public transport line (e.g., "11", "85", "S1")
	LineName string `json:"LineName"`

	// Direction indicates the destination or direction of the vehicle
	Direction string `json:"Direction"`

	// Platform contains information about the platform or stop position
	Platform Platform `json:"Platform"`

	// Mot indicates the mode of transport (e.g., "Tram", "Bus", "S-Bahn")
	Mot string `json:"Mot"`

	// RealTime is the actual departure/arrival time including delays
	RealTime string `json:"RealTime"`

	// ScheduledTime is the originally planned departure/arrival time
	ScheduledTime string `json:"ScheduledTime"`

	// State indicates the current status of the departure (e.g., "InTime", "Delayed", "Cancelled")
	State string `json:"State"`

	// RouteChanges contains information about any route diversions or changes
	RouteChanges []string `json:"RouteChanges"`

	// Diva contains DVB-specific identifiers for the vehicle/line
	Diva Diva `json:"Diva"`

	// CancelReasons contains reasons if the departure is cancelled
	CancelReasons []string `json:"CancelReasons"`

	// Occupancy indicates how crowded the vehicle is (e.g., "Low", "Medium", "High")
	Occupancy string `json:"Occupancy"`
}

// MonitorStop retrieves real-time departure and arrival information for a specific stop.
// This function is used to monitor when buses, trams, and trains will leave or arrive
// at a particular stop, including real-time updates about delays and cancellations.
//
// The function provides comprehensive information about each departure including
// line details, timing, platform information, and any service disruptions.
// This is essential for passengers who want to know when their next transport will arrive.
//
// Parameters:
//   - ctx: Context for the request, allowing for cancellation and timeouts
//   - options: Monitoring parameters including the required stop ID and optional filters
//
// Returns:
//   - *MonitorStopResponse: Contains the departure/arrival information and metadata
//   - error: Returns an error if the stop ID is empty or if the API request fails
//
// Example usage:
//
//	params := &MonitorStopParams{
//		StopId: "33000037",
//		Limit: &[]int{10}[0],
//		ShortTermChanges: &[]bool{true}[0],
//	}
//	response, err := client.MonitorStop(ctx, params)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Departures from %s:\n", response.Name)
//	for _, dep := range response.Departures {
//		fmt.Printf("Line %s to %s: %s\n", dep.LineName, dep.Direction, dep.RealTime)
//	}
func (c *Client) MonitorStop(ctx context.Context, options *MonitorStopParams) (*MonitorStopResponse, error) {
	query := url.Values{}

	if options != nil {
		if options.StopId != "" {
			query.Set("stopid", options.StopId)
		} else {
			return nil, errors.New("stopid can not be empty")
		}
		if options.Format != nil && *options.Format != "" {
			query.Set("format", *options.Format)
		}
		if options.Time != nil && *options.Time != "" {
			query.Set("time", *options.Time)
		}
		if options.IsArrival != nil {
			query.Set("isarrival", strconv.FormatBool(*options.IsArrival))
		}
		if options.Limit != nil && *options.Limit > 0 {
			query.Set("limit", strconv.Itoa(*options.Limit))
		}
		if options.ShortTermChanges != nil {
			query.Set("shorttermchanges", strconv.FormatBool(*options.ShortTermChanges))
		}
		if options.MentzOnly != nil {
			query.Set("mentzonly", strconv.FormatBool(*options.MentzOnly))
		}
	}

	opts := requestOptions{
		Method: http.MethodGet,
		Path:   "/dm",
		Query:  query,
	}

	resp, err := c.doRequest(ctx, opts)
	if err != nil {
		return nil, err
	}

	var resource MonitorStopResponse
	if err := c.handleResponse(resp, &resource); err != nil {
		return nil, err
	}

	return &resource, nil
}
