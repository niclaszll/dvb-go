package dvb

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// GetRouteParams contains the parameters for trip planning between two locations.
// This API provides journey planning with public transport connections, including
// walking routes, transfers, and timing information.
type GetRouteParams struct {
	// Origin is the starting point for the journey. This is required and cannot be empty.
	// Can be a stop ID (from GetPoint API) or a location name.
	Origin string

	// Destination is the end point for the journey. This is required and cannot be empty.
	// Can be a stop ID (from GetPoint API) or a location name.
	Destination string

	// Format specifies the response format. Optional parameter.
	// Supported values depend on the DVB API implementation.
	Format *string

	// IsArrivalTime when set to true, treats the Time parameter as arrival time.
	// When false or nil, treats Time as departure time (default behavior).
	IsArrivalTime *bool

	// ShortTermChanges when set to true, includes short-term changes like delays or cancellations.
	// When false or nil, uses scheduled times without real-time updates.
	ShortTermChanges *bool

	// Time specifies the departure or arrival time for the journey. Optional parameter.
	// Format should be compatible with the DVB API time format.
	// If not specified, uses the current time.
	Time *string

	// Via specifies an intermediate stop that the route should pass through.
	// Optional parameter for more specific route planning.
	Via *string
}

// GetRouteResponse represents the response from the DVB trip planning API.
// It contains multiple route options with detailed journey information.
type GetRouteResponse struct {
	// SessionId is a unique identifier for this trip planning session
	SessionId string `json:"SessionId"`

	// Status contains the API response status including error codes and messages
	Status Status `json:"Status"`

	// Routes is an array of possible journey options from origin to destination
	Routes []Route `json:"Routes"`
}

// Route represents a single journey option from origin to destination.
// It contains comprehensive information about the trip including timing, cost, and detailed steps.
type Route struct {
	// PriceLevel indicates the fare zone level for this route
	PriceLevel int `json:"PriceLevel"`

	// Price is the cost of a single ticket for this journey
	Price string `json:"Price"`

	// PriceDayTicket is the cost of a day ticket that covers this journey
	PriceDayTicket string `json:"PriceDayTicket"`

	// Net indicates the transport network (e.g., "VVO" for Dresden area)
	Net string `json:"Net"`

	// Duration is the total journey time in minutes
	Duration int `json:"Duration"`

	// Interchanges is the number of transfers required for this route
	Interchanges int `json:"Interchanges"`

	// MotChain lists all modes of transport used in this journey
	MotChain []MotChain `json:"MotChain"`

	// NumberOfFareZones indicates how many fare zones this journey crosses
	NumberOfFareZones string `json:"NumberOfFareZones"`

	// NumberOfFareZonesDayTicket indicates fare zones for day ticket pricing
	NumberOfFareZonesDayTicket string `json:"NumberOfFareZonesDayTicket"`

	// FareZoneNames lists the names of fare zones crossed during the journey
	FareZoneNames string `json:"FareZoneNames"`

	// FareZoneNamesDayTicket lists fare zone names for day ticket calculation
	FareZoneNamesDayTicket string `json:"FareZoneNamesDayTicket"`

	// FareZoneOrigin is the fare zone number of the starting point
	FareZoneOrigin int `json:"FareZoneOrigin"`

	// FareZoneDestination is the fare zone number of the destination
	FareZoneDestination int `json:"FareZoneDestination"`

	// RouteId is a unique identifier for this specific route
	RouteId int `json:"RouteId"`

	// PartialRoutes contains detailed step-by-step journey information
	PartialRoutes []PartialRoute `json:"PartialRoutes"`

	// MapData contains coordinate information for mapping the route
	MapData []string `json:"MapData"`

	// Tickets lists all available ticket types for this journey
	Tickets []Ticket `json:"Tickets"`
}

// MotChain represents a mode of transport used in the journey.
// This provides summary information about each transport line used.
type MotChain struct {
	// DlId is the DVB line identifier
	DlId string `json:"DlId"`

	// StatelessId is an alternative identifier for the transport line
	StatelessId string `json:"StatelessId"`

	// Type indicates the mode of transport (e.g., "Tram", "Bus", "S-Bahn", "Walking")
	Type string `json:"Type"`

	// Name is the line name or number (e.g., "11", "85", "S1")
	Name string `json:"Name"`

	// Direction indicates the destination or direction of the vehicle
	Direction string `json:"Direction"`

	// Changes contains information about any service changes or disruptions
	Changes []string `json:"Changes"`

	// Diva contains DVB-specific identifiers
	Diva Diva `json:"Diva"`

	// TransportationCompany is the name of the transport operator
	TransportationCompany string `json:"TransportationCompany"`

	// OperatorCode is the code identifying the transport operator
	OperatorCode string `json:"OperatorCode"`

	// ProductName describes the type of service (e.g., "Straßenbahn", "Stadtbus")
	ProductName string `json:"ProductName"`

	// TrainNumber is the specific train number for rail services
	TrainNumber string `json:"TrainNumber"`
}

// PartialRoute represents a single segment of the overall journey.
// This could be a walking segment, a ride on public transport, or a transfer.
type PartialRoute struct {
	// PartialRouteId is a unique identifier for this route segment
	PartialRouteId *int `json:"PartialRouteId,omitempty"`

	// Duration is the time required for this segment in minutes
	Duration int `json:"Duration"`

	// Mot contains detailed information about the mode of transport for this segment
	Mot Mot `json:"Mot"`

	// MapDataIndex points to the relevant map data for this segment
	MapDataIndex *int `json:"MapDataIndex,omitempty"`

	// Shift indicates any timing adjustments for this segment
	Shift string `json:"Shift"`

	// RegularStops lists all stops visited during this segment (for public transport)
	RegularStops []RegularStop `json:"RegularStops,omitempty"`

	// ChangeoverEndangered indicates if a transfer connection might be at risk
	ChangeoverEndangered *bool `json:"ChangeoverEndangered,omitempty"`

	// NextDepartureTimes lists alternative departure times for this segment
	NextDepartureTimes []string `json:"NextDepartureTimes,omitempty"`

	// PreviousDepartureTimes lists earlier departure options for this segment
	PreviousDepartureTimes []string `json:"PreviousDepartureTimes,omitempty"`
}

// Mot represents detailed mode of transport information for a route segment.
// This contains more specific information than the summary in MotChain.
type Mot struct {
	// DlId is the DVB line identifier
	DlId *string `json:"DlId,omitempty"`

	// StatelessId is an alternative identifier for the transport line
	StatelessId *string `json:"StatelessId,omitempty"`

	// Type indicates the mode of transport (e.g., "Tram", "Bus", "S-Bahn", "Walking")
	Type string `json:"Type"`

	// Name is the line name or number (e.g., "11", "85", "S1")
	Name *string `json:"Name,omitempty"`

	// Direction indicates the destination or direction of the vehicle
	Direction *string `json:"Direction,omitempty"`

	// Changes contains information about any service changes or disruptions
	Changes []string `json:"Changes,omitempty"`

	// Diva contains DVB-specific identifiers
	Diva *Diva `json:"Diva,omitempty"`

	// TransportationCompany is the name of the transport operator
	TransportationCompany *string `json:"TransportationCompany,omitempty"`

	// OperatorCode is the code identifying the transport operator
	OperatorCode *string `json:"OperatorCode,omitempty"`

	// ProductName describes the type of service (e.g., "Straßenbahn", "Stadtbus")
	ProductName *string `json:"ProductName,omitempty"`

	// TrainNumber is the specific train number for rail services
	TrainNumber *string `json:"TrainNumber,omitempty"`
}

// RegularStop represents a stop visited during a journey segment.
// This provides detailed timing and location information for each stop.
type RegularStop struct {
	// ArrivalTime is the scheduled arrival time at this stop
	ArrivalTime string `json:"ArrivalTime"`

	// DepartureTime is the scheduled departure time from this stop
	DepartureTime string `json:"DepartureTime"`

	// ArrivalRealTime is the real-time arrival time including delays
	ArrivalRealTime *string `json:"ArrivalRealTime,omitempty"`

	// DepartureRealTime is the real-time departure time including delays
	DepartureRealTime *string `json:"DepartureRealTime,omitempty"`

	// Place indicates the city or area where this stop is located
	Place string `json:"Place"`

	// Name is the official name of the stop
	Name string `json:"Name"`

	// Type indicates the type of stop (e.g., "Platform", "Stop")
	Type string `json:"Type"`

	// DataId is a unique identifier for this stop in the data system
	DataId string `json:"DataId"`

	// DhId is an alternative identifier for this stop
	DhId string `json:"DhId"`

	// Platform contains information about the platform or stop position
	Platform Platform `json:"Platform"`

	// Latitude is the geographical latitude coordinate of the stop
	Latitude int `json:"Latitude"`

	// Longitude is the geographical longitude coordinate of the stop
	Longitude int `json:"Longitude"`

	// DepartureState indicates the current status of departures from this stop
	DepartureState *string `json:"DepartureState,omitempty"`

	// ArrivalState indicates the current status of arrivals at this stop
	ArrivalState *string `json:"ArrivalState,omitempty"`

	// CancelReasons contains reasons if services at this stop are cancelled
	CancelReasons []string `json:"CancelReasons"`

	// ParkAndRail contains information about park and ride facilities
	ParkAndRail []string `json:"ParkAndRail"`

	// Occupancy indicates how crowded the vehicle is at this stop
	Occupancy string `json:"Occupancy"`
}

// Ticket represents a ticket option available for the journey.
// This includes pricing information for different ticket types.
type Ticket struct {
	// Name is the display name of the ticket type
	Name string `json:"Name"`

	// PriceLevel indicates the fare zone level for this ticket
	PriceLevel int `json:"PriceLevel"`

	// Price is the cost of this ticket type
	Price string `json:"Price"`

	// NumberOfFareZones indicates how many fare zones this ticket covers
	NumberOfFareZones string `json:"NumberOfFareZones"`

	// FareZoneNames lists the names of fare zones covered by this ticket
	FareZoneNames string `json:"FareZoneNames"`
}

// GetRoute plans a journey between two locations using public transport.
// This function provides comprehensive trip planning including multiple route options,
// timing information, transfer details, and pricing for Dresden's public transport network.
//
// The function returns multiple route alternatives with detailed step-by-step instructions,
// real-time information (when available), and complete fare information.
// This is the primary function for journey planning and route discovery.
//
// Parameters:
//   - ctx: Context for the request, allowing for cancellation and timeouts
//   - options: Trip planning parameters including required origin and destination, plus optional timing and routing preferences
//
// Returns:
//   - *GetRouteResponse: Contains multiple route options with detailed journey information
//   - error: Returns an error if origin or destination is empty, or if the API request fails
//
// Example usage:
//
//	params := &GetRouteParams{
//		Origin: "Hauptbahnhof",
//		Destination: "Neustadt Bahnhof",
//		ShortTermChanges: &[]bool{true}[0],
//	}
//	response, err := client.GetRoute(ctx, params)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Found %d route options:\n", len(response.Routes))
//	for i, route := range response.Routes {
//		fmt.Printf("Route %d: %d minutes, %d transfers, Price: %s\n",
//			i+1, route.Duration, route.Interchanges, route.Price)
//	}
func (c *Client) GetRoute(ctx context.Context, options *GetRouteParams) (*GetRouteResponse, error) {
	query := url.Values{}

	if options != nil {
		if options.Origin != "" {
			query.Set("origin", options.Origin)
		} else {
			return nil, errors.New("origin can not be empty")
		}
		if options.Destination != "" {
			query.Set("destination", options.Destination)
		} else {
			return nil, errors.New("destination can not be empty")
		}
		if options.Format != nil && *options.Format != "" {
			query.Set("format", *options.Format)
		}
		if options.IsArrivalTime != nil {
			query.Set("isarrivaltime", strconv.FormatBool(*options.IsArrivalTime))
		}
		if options.ShortTermChanges != nil {
			query.Set("shorttermchanges", strconv.FormatBool(*options.ShortTermChanges))
		}
		if options.Time != nil && *options.Time != "" {
			query.Set("time", *options.Time)
		}
		if options.Via != nil && *options.Via != "" {
			query.Set("via", *options.Via)
		}
	}

	opts := requestOptions{
		Method: http.MethodGet,
		Path:   "/tr/trips",
		Query:  query,
	}

	resp, err := c.doRequest(ctx, opts)
	if err != nil {
		return nil, err
	}

	var resource GetRouteResponse
	if err := c.handleResponse(resp, &resource); err != nil {
		return nil, err
	}

	return &resource, nil
}
