package dvb

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type GetRouteOptions struct {
	Origin           string
	Destination      string
	Format           *string
	IsArrivalTime    *bool
	ShortTermChanges *bool
	Time             *string
	Via              *string
}

type GetRouteResponse struct {
	SessionId string  `json:"SessionId"`
	Status    Status  `json:"Status"`
	Routes    []Route `json:"Routes"`
}

type Route struct {
	PriceLevel                 int            `json:"PriceLevel"`
	Price                      string         `json:"Price"`
	PriceDayTicket             string         `json:"PriceDayTicket"`
	Net                        string         `json:"Net"`
	Duration                   int            `json:"Duration"`
	Interchanges               int            `json:"Interchanges"`
	MotChain                   []MotChain     `json:"MotChain"`
	NumberOfFareZones          string         `json:"NumberOfFareZones"`
	NumberOfFareZonesDayTicket string         `json:"NumberOfFareZonesDayTicket"`
	FareZoneNames              string         `json:"FareZoneNames"`
	FareZoneNamesDayTicket     string         `json:"FareZoneNamesDayTicket"`
	FareZoneOrigin             int            `json:"FareZoneOrigin"`
	FareZoneDestination        int            `json:"FareZoneDestination"`
	RouteId                    int            `json:"RouteId"`
	PartialRoutes              []PartialRoute `json:"PartialRoutes"`
	MapData                    []string       `json:"MapData"`
	Tickets                    []Ticket       `json:"Tickets"`
}

type MotChain struct {
	DlId                  string   `json:"DlId"`
	StatelessId           string   `json:"StatelessId"`
	Type                  string   `json:"Type"`
	Name                  string   `json:"Name"`
	Direction             string   `json:"Direction"`
	Changes               []string `json:"Changes"`
	Diva                  Diva     `json:"Diva"`
	TransportationCompany string   `json:"TransportationCompany"`
	OperatorCode          string   `json:"OperatorCode"`
	ProductName           string   `json:"ProductName"`
	TrainNumber           string   `json:"TrainNumber"`
}

type PartialRoute struct {
	PartialRouteId         *int          `json:"PartialRouteId,omitempty"`
	Duration               int           `json:"Duration"`
	Mot                    Mot           `json:"Mot"`
	MapDataIndex           *int          `json:"MapDataIndex,omitempty"`
	Shift                  string        `json:"Shift"`
	RegularStops           []RegularStop `json:"RegularStops,omitempty"`
	ChangeoverEndangered   *bool         `json:"ChangeoverEndangered,omitempty"`
	NextDepartureTimes     []string      `json:"NextDepartureTimes,omitempty"`
	PreviousDepartureTimes []string      `json:"PreviousDepartureTimes,omitempty"`
}

type Mot struct {
	DlId                  *string  `json:"DlId,omitempty"`
	StatelessId           *string  `json:"StatelessId,omitempty"`
	Type                  string   `json:"Type"`
	Name                  *string  `json:"Name,omitempty"`
	Direction             *string  `json:"Direction,omitempty"`
	Changes               []string `json:"Changes,omitempty"`
	Diva                  *Diva    `json:"Diva,omitempty"`
	TransportationCompany *string  `json:"TransportationCompany,omitempty"`
	OperatorCode          *string  `json:"OperatorCode,omitempty"`
	ProductName           *string  `json:"ProductName,omitempty"`
	TrainNumber           *string  `json:"TrainNumber,omitempty"`
}

type RegularStop struct {
	ArrivalTime       string   `json:"ArrivalTime"`
	DepartureTime     string   `json:"DepartureTime"`
	ArrivalRealTime   *string  `json:"ArrivalRealTime,omitempty"`
	DepartureRealTime *string  `json:"DepartureRealTime,omitempty"`
	Place             string   `json:"Place"`
	Name              string   `json:"Name"`
	Type              string   `json:"Type"`
	DataId            string   `json:"DataId"`
	DhId              string   `json:"DhId"`
	Platform          Platform `json:"Platform"`
	Latitude          int      `json:"Latitude"`
	Longitude         int      `json:"Longitude"`
	DepartureState    *string  `json:"DepartureState,omitempty"`
	ArrivalState      *string  `json:"ArrivalState,omitempty"`
	CancelReasons     []string `json:"CancelReasons"`
	ParkAndRail       []string `json:"ParkAndRail"`
	Occupancy         string   `json:"Occupancy"`
}

type Ticket struct {
	Name              string `json:"Name"`
	PriceLevel        int    `json:"PriceLevel"`
	Price             string `json:"Price"`
	NumberOfFareZones string `json:"NumberOfFareZones"`
	FareZoneNames     string `json:"FareZoneNames"`
}

// Get a route between two stops
func (c *Client) GetRoute(ctx context.Context, options *GetRouteOptions) (*GetRouteResponse, error) {
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
