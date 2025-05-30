package dvb

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

type MonitorStopOptions struct {
	StopId           string
	Format           *string
	Time             *string
	IsArrival        *bool
	Limit            *int
	ShortTermChanges *bool
	MentzOnly        *bool
}

type MonitorStopResponse struct {
	Name           string      `json:"Name"`
	Status         Status      `json:"Status"`
	Place          string      `json:"Place"`
	ExpirationTime string      `json:"ExpirationTime"`
	Departures     []Departure `json:"Departures"`
}

type Departure struct {
	Id            string   `json:"Id"`
	DlId          string   `json:"DlId"`
	LineName      string   `json:"LineName"`
	Direction     string   `json:"Direction"`
	Platform      Platform `json:"Platform"`
	Mot           string   `json:"Mot"`
	RealTime      string   `json:"RealTime"`
	ScheduledTime string   `json:"ScheduledTime"`
	State         string   `json:"State"`
	RouteChanges  []string `json:"RouteChanges"`
	Diva          Diva     `json:"Diva"`
	CancelReasons []string `json:"CancelReasons"`
	Occupancy     string   `json:"Occupancy"`
}

// Monitor a single stop to see every bus or tram leaving this stop after the specified time offset.
func (c *Client) MonitorStop(ctx context.Context, options *MonitorStopOptions) (*MonitorStopResponse, error) {
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

	opts := RequestOptions{
		Method: GET,
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
