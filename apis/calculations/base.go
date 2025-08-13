package calculations

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the public API host.
	DefaultBaseURL = "https://apis.rentier.io"
	// DefaultCalculateAVMPath is the path as referenced by the docs UI.
	DefaultCalculateAVMPath = "/calculation/calculateAVM"
)

// Client is a reusable HTTP client for the Rentier APIs.
// It is safe for concurrent use.
type Client struct {
	baseURL *url.URL
	hc      *http.Client
	// auth header to use (e.g., "Authorization", "X-API-Key") and its value
	authHeader string
	authValue  string

	// default headers applied to every request
	defaultHeaders http.Header

	// retry configuration
	retryMax     int
	retryWaitMin time.Duration
	retryWaitMax time.Duration
}

// Option configures the Client.

type Option func(*Client) error

// WithBaseURL overrides the API base URL. Accepts values like
// "https://apis.rentier.io" (no trailing slash required).
func WithBaseURL(raw string) Option {
	return func(c *Client) error {
		u, err := url.Parse(strings.TrimRight(raw, "/"))
		if err != nil {
			return fmt.Errorf("invalid base URL: %w", err)
		}
		c.baseURL = u
		return nil
	}
}

// WithHTTPClient injects a custom http.Client (timeouts, proxies, etc.).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) error {
		if h == nil {
			return errors.New("nil http.Client")
		}
		c.hc = h
		return nil
	}
}

// WithAPIKey sets a custom API key header (e.g., "X-API-Key").
func WithAPIKey(headerName, key string) Option {
	return func(c *Client) error {
		if headerName == "" {
			return errors.New("headerName cannot be empty")
		}
		c.authHeader = headerName
		c.authValue = key
		return nil
	}
}

// WithDefaultHeader attaches a default header to every request.
func WithDefaultHeader(k, v string) Option {
	return func(c *Client) error {
		if c.defaultHeaders == nil {
			c.defaultHeaders = make(http.Header)
		}
		c.defaultHeaders.Set(k, v)
		return nil
	}
}

// WithRetry configures simple exponential backoff for transient errors.
// retryMax is the total number of attempts (including the first).
func WithRetry(retryMax int, waitMin, waitMax time.Duration) Option {
	return func(c *Client) error {
		if retryMax < 1 {
			return errors.New("retryMax must be >= 1")
		}
		c.retryMax = retryMax
		c.retryWaitMin = waitMin
		c.retryWaitMax = waitMax
		return nil
	}
}

func NewClient(opts ...Option) (*Client, error) {
	base, _ := url.Parse(DefaultBaseURL)
	c := &Client{
		baseURL:      base,
		hc:           &http.Client{Timeout: 30 * time.Second},
		retryMax:     3,
		retryWaitMin: 250 * time.Millisecond,
		retryWaitMax: 2 * time.Second,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// RequestOption customizes an individual request.

type RequestOption func(r *http.Request)

// WithHeader adds/overrides a single header for the request.
func WithHeader(k, v string) RequestOption { return func(r *http.Request) { r.Header.Set(k, v) } }

// WithPath lets you override the path (useful if the server uses a different route).
func WithPath(path string) RequestOption { return func(r *http.Request) { r.URL.Path = path } }

// APIError represents a non-2xx response.

type APIError struct {
	StatusCode int
	Body       []byte
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("api error: status=%d msg=%s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("api error: status=%d body=%q", e.StatusCode, string(e.Body))
}
