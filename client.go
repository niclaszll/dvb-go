// Package dvb provides a Go client for interacting with the Dresden Transport (DVB) API.
// It allows you to fetch real-time public transport information including departures,
// routes, lines, and stop information for Dresden's public transportation network.
//
// Example usage:
//
//	client := dvb.NewClient(dvb.Config{})
//	response, err := client.MonitorStop(ctx, &dvb.MonitorStopParams{
//		StopId: "33000028", // Dresden Hauptbahnhof
//	})
package dvb

import (
	"net/http"
	"time"
)

// Client represents a DVB API client with configuration for making requests.
type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

// Config holds configuration options for creating a new DVB client.
type Config struct {
	BaseURL    string        // Base URL for the DVB API (optional, defaults to official API)
	UserAgent  string        // User agent string for requests (optional)
	Timeout    time.Duration // HTTP timeout for requests (optional, defaults to 30s)
	HTTPClient *http.Client  // Custom HTTP client (optional)
}

// NewClient creates a new DVB API client with the provided configuration.
// If no configuration is provided, sensible defaults will be used.
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = "https://webapi.vvo-online.de"
	}

	if config.UserAgent == "" {
		config.UserAgent = "dvb-go-client/1.0.0"
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return &Client{
		baseURL:    config.BaseURL,
		httpClient: httpClient,
		userAgent:  config.UserAgent,
	}
}
