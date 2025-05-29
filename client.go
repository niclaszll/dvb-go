package dvb

import (
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

type Config struct {
	BaseURL    string
	UserAgent  string
	Timeout    time.Duration
	HTTPClient *http.Client
}

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
