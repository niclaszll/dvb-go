package dvb

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

type GetPointOptions struct {
	Query         string
	Format        *string
	StopsOnly     *bool
	AssignedStops *bool
	Limit         *int
	Dvb           *bool
}

type GetPointResponse struct {
	PointStatus    string   `json:"PointStatus"`
	Status         Status   `json:"Status"`
	Points         []string `json:"Points"`
	ExpirationTime string   `json:"ExpirationTime"`
}

// Get a point by query
func (c *Client) GetPoint(ctx context.Context, options *GetPointOptions) (*GetPointResponse, error) {
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

	opts := RequestOptions{
		Method: GET,
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
