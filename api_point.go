package dvb

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// GetPointParams contains the parameters for finding a point/stop using the DVB point finder API.
// The point finder allows searching for public transport stops, stations, and points of interest.
type GetPointParams struct {
	// Query is the search term for finding a point. This is required and cannot be empty.
	// Examples: "Hauptbahnhof", "Dresden Altstadt", "Neustadt Bahnhof"
	Query string

	// Format specifies the response format. Optional parameter.
	// Supported values depend on the DVB API implementation.
	Format *string

	// StopsOnly when set to true, limits the search results to public transport stops only.
	// When false or nil, includes other points of interest as well.
	StopsOnly *bool

	// AssignedStops when set to true, includes only stops that are assigned to specific lines.
	// When false or nil, includes all matching stops regardless of line assignments.
	AssignedStops *bool

	// Limit restricts the maximum number of results returned.
	// Must be a positive integer. If nil or 0, uses the API's default limit.
	Limit *int

	// Dvb when set to true, includes only DVB (Dresden public transport) stops.
	// When false or nil, may include stops from other transport providers.
	Dvb *bool
}

// GetPointResponse represents the response from the DVB point finder API.
// It contains information about the search results and status.
type GetPointResponse struct {
	// PointStatus indicates the status of the point search operation
	PointStatus string `json:"PointStatus"`

	// Status contains the API response status including error codes and messages
	Status Status `json:"Status"`

	// Points is an array of point identifiers that match the search query
	Points []string `json:"Points"`

	// ExpirationTime indicates when this response data expires and should be refreshed
	ExpirationTime string `json:"ExpirationTime"`
}

// GetPoint searches for public transport stops, stations, and points of interest
// using the DVB point finder API. This is typically the first step when looking
// up public transport information, as you need stop IDs for other API calls.
//
// The function performs a search based on the provided query string and returns
// matching points. The search can be refined using various optional parameters
// to filter results by type, provider, or limit the number of results.
//
// Parameters:
//   - ctx: Context for the request, allowing for cancellation and timeouts
//   - options: Search parameters including the required query string and optional filters
//
// Returns:
//   - *GetPointResponse: Contains the search results and metadata
//   - error: Returns an error if the query is empty or if the API request fails
//
// Example usage:
//
//	params := &GetPointParams{
//		Query: "Hauptbahnhof",
//		StopsOnly: &[]bool{true}[0],
//		Limit: &[]int{5}[0],
//	}
//	response, err := client.GetPoint(ctx, params)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, point := range response.Points {
//		fmt.Println("Found point:", point)
//	}
func (c *Client) GetPoint(ctx context.Context, options *GetPointParams) (*GetPointResponse, error) {
	query := url.Values{}

	if options != nil {
		if options.Query != "" {
			query.Set("query", options.Query)
		} else {
			return nil, errors.New("query can not be empty")
		}
		if options.Format != nil && *options.Format != "" {
			query.Set("format", *options.Format)
		}
		if options.Limit != nil && *options.Limit > 0 {
			query.Set("limit", strconv.Itoa(*options.Limit))
		}
		if options.StopsOnly != nil {
			query.Set("stopsOnly", strconv.FormatBool(*options.StopsOnly))
		}
		if options.AssignedStops != nil {
			query.Set("assignedStops", strconv.FormatBool(*options.AssignedStops))
		}
		if options.Dvb != nil {
			query.Set("dvb", strconv.FormatBool(*options.Dvb))
		}
	}

	opts := requestOptions{
		Method: http.MethodGet,
		Path:   "/tr/pointfinder",
		Query:  query,
	}

	resp, err := c.doRequest(ctx, opts)
	if err != nil {
		return nil, err
	}

	var resource GetPointResponse
	if err := c.handleResponse(resp, &resource); err != nil {
		return nil, err
	}

	return &resource, nil
}
