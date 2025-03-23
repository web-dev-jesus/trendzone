package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/util"
)

// Client represents the API client for SportsData.io
type Client struct {
	httpClient  *http.Client
	baseURL     string
	apiKey      string
	rateLimiter *time.Ticker
	logger      *util.Logger
}

// NewClient creates a new API client with rate limiting
func NewClient(config *config.Config, logger *util.Logger) *Client {
	// Set up HTTP client with security best practices
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	// Create rate limiter ticker to avoid hitting API limits
	rateLimiter := time.NewTicker(config.APICallDelay)

	return &Client{
		httpClient:  httpClient,
		baseURL:     config.SportsDataBaseURL,
		apiKey:      config.SportsDataAPIKey,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

// FetchData fetches data from the API and unmarshals it into the provided interface
func (c *Client) FetchData(endpoint string, result interface{}) error {
	// Wait for rate limiter before making API call
	<-c.rateLimiter.C

	// Construct full URL with API key
	url := fmt.Sprintf("%s%s?key=%s", c.baseURL, endpoint, c.apiKey)

	c.logger.Debug("Fetching data from: " + util.SanitizeURL(url))

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers for better security
	req.Header.Set("User-Agent", "SportsData-NFL-Client/1.0")
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	// Read response body with a size limit for security
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024)) // 10MB limit
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Unmarshal response
	if err := json.Unmarshal(bodyBytes, result); err != nil {
		return fmt.Errorf("error unmarshaling response: %w", err)
	}

	return nil
}

// Close closes the rate limiter
func (c *Client) Close() {
	c.rateLimiter.Stop()
}
