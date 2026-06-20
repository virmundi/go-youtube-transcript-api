package lib

import (
	"net/http"
	"time"
)

// Client is a wrapper around http.Client to handle requests to YouTube.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get performs a GET request to the specified URL.
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// YouTube expects a user-agent.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")

	return c.httpClient.Do(req)
}
