// Package client implements a Go client for the SMS Manager APIs.
//
// Sending goes through the JSON API v2 (https://api.smsmngr.com/v2), account,
// status and pricing lookups go through the REST API v1
// (https://rest-api.smsmngr.com/v1). Both are authenticated with the same API
// key, passed in the x-api-key header.
//
// https://smsmanager.cz/docs/api-reference/json-v2/overview
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// JSONAPIBaseURL is the base URL of the JSON API v2 (message sending).
	JSONAPIBaseURL = "https://api.smsmngr.com/v2"

	// RESTAPIBaseURL is the base URL of the REST API v1 (account, message status, pricing).
	RESTAPIBaseURL = "https://rest-api.smsmngr.com/v1"

	// defaultUserAgent is sent with every request unless overridden by WithUserAgent.
	defaultUserAgent = "sms-manager-go"

	// defaultTimeout is used by the HTTP client created by New.
	defaultTimeout = 30 * time.Second
)

var (
	// ErrMissingAPIKey is returned by New when no API key is given.
	ErrMissingAPIKey = errors.New("api key is empty")

	// ErrNotFound is returned when the API has no record for the given identifier.
	ErrNotFound = errors.New("not found")
)

// APIError is a non-2xx response, or a REST API response reporting success: false.
type APIError struct {
	// StatusCode is the HTTP status code returned by the server.
	StatusCode int

	// Message is the human-readable error message returned by the API, if any.
	Message string

	// RequestID is the request identifier returned by the API, if any.
	RequestID string
}

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = http.StatusText(e.StatusCode)
	}

	if e.RequestID != "" {
		return fmt.Sprintf("sms manager api: %d: %s (request_id %s)", e.StatusCode, msg, e.RequestID)
	}

	return fmt.Sprintf("sms manager api: %d: %s", e.StatusCode, msg)
}

// Client talks to the SMS Manager APIs. Use New to create one; it is safe for
// concurrent use by multiple goroutines.
type Client struct {
	apiKey      string
	jsonBaseURL string
	restBaseURL string
	userAgent   string
	httpClient  *http.Client
}

// Option customizes a Client created by New.
type Option func(*Client)

// WithHTTPClient sets the HTTP client used for all requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithJSONAPIBaseURL overrides the JSON API v2 base URL (useful for testing).
func WithJSONAPIBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.jsonBaseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithRESTAPIBaseURL overrides the REST API v1 base URL (useful for testing).
func WithRESTAPIBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.restBaseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// New returns a Client authenticated with the given API key. The key is issued
// in the account settings at https://app.smsmanager.com/api-cloud.
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	c := &Client{
		apiKey:      apiKey,
		jsonBaseURL: JSONAPIBaseURL,
		restBaseURL: RESTAPIBaseURL,
		userAgent:   defaultUserAgent,
		httpClient:  &http.Client{Timeout: defaultTimeout},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// postJSON sends payload as a JSON body and decodes the response into out.
func (c *Client) postJSON(ctx context.Context, baseURL, path string, payload, out interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, out)
}

// postForm sends params as an urlencoded body and decodes the response into out.
func (c *Client) postForm(ctx context.Context, baseURL, path string, params url.Values, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req, out)
}

// get issues a GET request with the given query params and decodes the response into out.
func (c *Client) get(ctx context.Context, baseURL, path string, params url.Values, out interface{}) error {
	fullURL := baseURL + path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	return c.do(req, out)
}

// do authenticates and executes req, then decodes a successful response body into out.
func (c *Client) do(req *http.Request, out interface{}) error {
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return newAPIError(resp.StatusCode, body)
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

// newAPIError builds an APIError out of an error response body. Both APIs return
// a JSON object with a message field, the JSON API v2 sometimes capitalizes it
// (decoding is case-insensitive, so a single field covers both).
func newAPIError(statusCode int, body []byte) error {
	apiErr := &APIError{StatusCode: statusCode}

	var payload struct {
		Message   string `json:"message"`
		RequestID string `json:"request_id"`
	}

	if err := json.Unmarshal(body, &payload); err == nil {
		apiErr.Message = payload.Message
		apiErr.RequestID = payload.RequestID
	}

	if apiErr.Message == "" {
		apiErr.Message = strings.TrimSpace(string(body))
	}

	return apiErr
}
