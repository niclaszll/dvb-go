package dvb

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// GetLinesParams contains the parameters for retrieving available public transport lines for a stop.
// This API provides information about which bus and tram lines serve a particular stop.
type GetLinesParams struct {
	// StopId is the unique identifier for the stop. This is required and cannot be empty.
	// Use the GetPoint API to find stop IDs based on stop names or locations.
	StopId string

	// Format specifies the response format. Optional parameter.
	// Supported values depend on the DVB API implementation.
	Format *string
}

// GetLinesResponse represents the response from the DVB lines API.
// It contains information about all public transport lines that serve a specific stop.
type GetLinesResponse struct {
	// Lines is an array of public transport lines that serve the specified stop
	Lines []Line `json:"Lines"`

	// Status contains the API response status including error codes and messages
	Status Status `json:"Status"`

	// ExpirationTime indicates when this response data expires and should be refreshed
	ExpirationTime string `json:"ExpirationTime"`
}

// Line represents a single public transport line that serves a stop.
// It contains information about the line, its directions, and available timetables.
type Line struct {
	// Name is the display name of the line (e.g., "11", "85", "S1")
	Name string `json:"Name"`

	// Mot indicates the mode of transport (e.g., "Tram", "Bus", "S-Bahn")
	Mot string `json:"Mot"`

	// Changes contains information about any service changes or disruptions for this line
	Changes []string `json:"Changes,omitempty"`

	// Directions lists all directions this line travels from the current stop
	Directions []Direction `json:"Directions"`

	// Diva contains DVB-specific identifiers for the line
	Diva Diva `json:"Diva"`
}

// Direction represents a specific direction or destination for a public transport line.
// Each line typically has two directions (e.g., inbound and outbound).
type Direction struct {
	// Name is the destination name for this direction (e.g., "Striesen", "Löbtau")
	Name string `json:"Name"`

	// TimeTables lists available timetables for this direction
	TimeTables []TimeTable `json:"TimeTables"`
}

// TimeTable represents a specific timetable variant for a line direction.
// Different timetables may apply during different times or days.
type TimeTable struct {
	// Id is the unique identifier for this timetable
	Id string `json:"Id"`

	// Name is the display name for this timetable (e.g., "Werktag", "Samstag")
	Name string `json:"Name"`
}

// GetLines retrieves a list of all public transport lines that serve a specific stop.
// This function is useful for discovering which buses, trams, and trains are available
// at a particular stop, along with their destinations and timetable information.
//
// The response includes comprehensive information about each line including
// the mode of transport, available directions, and timetable variants.
// This information is essential for journey planning and understanding
// transport options at any given stop.
//
// Parameters:
//   - ctx: Context for the request, allowing for cancellation and timeouts
//   - options: Parameters including the required stop ID and optional format specification
//
// Returns:
//   - *GetLinesResponse: Contains the list of lines and metadata
//   - error: Returns an error if the stop ID is empty or if the API request fails
//
// Example usage:
//
//	params := &GetLinesParams{
//		StopId: "33000037",
//	}
//	response, err := client.GetLines(ctx, params)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Lines serving this stop:\n")
//	for _, line := range response.Lines {
//		fmt.Printf("Line %s (%s):\n", line.Name, line.Mot)
//		for _, direction := range line.Directions {
//			fmt.Printf("  → %s\n", direction.Name)
//		}
//	}
func (c *Client) GetLines(ctx context.Context, options *GetLinesParams) (*GetLinesResponse, error) {
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

	}

	opts := requestOptions{
		Method: http.MethodGet,
		Path:   "/stt/lines",
		Query:  query,
	}

	resp, err := c.doRequest(ctx, opts)
	if err != nil {
		return nil, err
	}

	var resource GetLinesResponse
	if err := c.handleResponse(resp, &resource); err != nil {
		return nil, err
	}

	return &resource, nil
}
