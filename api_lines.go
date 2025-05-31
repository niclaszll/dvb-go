package dvb

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

type GetLinesParams struct {
	StopId string
	Format *string
}

type GetLinesResponse struct {
	Lines          []Line `json:"Lines"`
	Status         Status `json:"Status"`
	ExpirationTime string `json:"ExpirationTime"`
}

type Line struct {
	Name       string      `json:"Name"`
	Mot        string      `json:"Mot"`
	Changes    []string    `json:"Changes,omitempty"`
	Directions []Direction `json:"Directions"`
	Diva       Diva        `json:"Diva"`
}

type Direction struct {
	Name       string      `json:"Name"`
	TimeTables []TimeTable `json:"TimeTables"`
}

type TimeTable struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

// Get a list of available tram/bus lines for a stop
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
