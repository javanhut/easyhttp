// Package easyhttp provides a simple, powerful HTTP client library for Go
// that combines the performance of Go's native HTTP client with the ease of use
// of Python's requests library and the flexibility of easyjson for JSON handling.
package easyhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/javanhut/easyjson"
)

// Response wraps http.Response with convenient methods
type Response struct {
	*http.Response
	body []byte
}

// Text returns the response body as a string
func (r *Response) Text() string {
	if r.body == nil {
		r.readBody()
	}
	return string(r.body)
}

// JSON unmarshals the response body into the provided interface using easyjson
func (r *Response) JSON(v interface{}) error {
	if r.body == nil {
		r.readBody()
	}

	// Parse JSON using easyjson and then extract to the target interface
	jsonValue, err := easyjson.Load(r.body)
	if err != nil {
		return err
	}

	// Convert the easyjson value back to bytes and unmarshal into target
	jsonBytes, err := jsonValue.Dump()
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, v)
}

// JSONValue returns the response body as an easyjson JSONValue for fluent access
func (r *Response) JSONValue() (*easyjson.JSONValue, error) {
	if r.body == nil {
		r.readBody()
	}
	return easyjson.Load(r.body)
}

// Bytes returns the response body as bytes
func (r *Response) Bytes() []byte {
	if r.body == nil {
		r.readBody()
	}
	return r.body
}

// OK returns true if status code is 2xx
func (r *Response) OK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func (r *Response) readBody() {
	if r.Body != nil {
		defer r.Body.Close()
		r.body, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(r.body))
	}
}

// Client provides an easy-to-use HTTP client
type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
	auth    *Auth
}

// Auth holds authentication information
type Auth struct {
	Username string
	Password string
	Token    string
}

// RequestOptions holds options for HTTP requests
type RequestOptions struct {
	Headers       map[string]string
	Params        map[string]string
	JSON          interface{}
	Data          interface{}
	Auth          *Auth
	Timeout       time.Duration
	AllowRedirect bool
}

// New creates a new HTTP client
func New() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

// SetBaseURL sets the base URL for all requests
func (c *Client) SetBaseURL(baseURL string) *Client {
	c.baseURL = strings.TrimSuffix(baseURL, "/")
	return c
}

// SetTimeout sets the default timeout for requests
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.client.Timeout = timeout
	return c
}

// SetHeaders sets default headers for all requests
func (c *Client) SetHeaders(headers map[string]string) *Client {
	for k, v := range headers {
		c.headers[k] = v
	}
	return c
}

// SetAuth sets default authentication
func (c *Client) SetAuth(auth *Auth) *Client {
	c.auth = auth
	return c
}

// GET performs a GET request
func (c *Client) GET(url string, opts ...*RequestOptions) (*Response, error) {
	return c.request("GET", url, opts...)
}

// POST performs a POST request
func (c *Client) POST(url string, opts ...*RequestOptions) (*Response, error) {
	return c.request("POST", url, opts...)
}

// PUT performs a PUT request
func (c *Client) PUT(url string, opts ...*RequestOptions) (*Response, error) {
	return c.request("PUT", url, opts...)
}

// DELETE performs a DELETE request
func (c *Client) DELETE(url string, opts ...*RequestOptions) (*Response, error) {
	return c.request("DELETE", url, opts...)
}

// PATCH performs a PATCH request
func (c *Client) PATCH(url string, opts ...*RequestOptions) (*Response, error) {
	return c.request("PATCH", url, opts...)
}

// HEAD performs a HEAD request
func (c *Client) HEAD(url string, opts ...*RequestOptions) (*Response, error) {
	return c.request("HEAD", url, opts...)
}

func (c *Client) request(method, reqURL string, opts ...*RequestOptions) (*Response, error) {
	// Merge options
	opt := &RequestOptions{}
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	// Build full URL
	fullURL := c.buildURL(reqURL, opt.Params)

	// Prepare body
	var body io.Reader
	var contentType string

	if opt.JSON != nil {
		// Handle different JSON input types
		var jsonData []byte
		var err error

		switch v := opt.JSON.(type) {
		case *easyjson.JSONValue:
			// If it's already an easyjson JSONValue, use its Dump method
			jsonData, err = v.Dump()
		case string:
			// If it's a JSON string, parse it first to validate
			jsonValue, parseErr := easyjson.Loads(v)
			if parseErr != nil {
				return nil, fmt.Errorf("invalid JSON string: %w", parseErr)
			}
			jsonData, err = jsonValue.Dump()
		default:
			// For other types, use standard JSON marshaling
			jsonData, err = json.Marshal(opt.JSON)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		body = bytes.NewReader(jsonData)
		contentType = "application/json"
	} else if opt.Data != nil {
		switch data := opt.Data.(type) {
		case string:
			body = strings.NewReader(data)
		case []byte:
			body = bytes.NewReader(data)
		case io.Reader:
			body = data
		default:
			return nil, fmt.Errorf("unsupported data type: %T", data)
		}
	}

	// Create request
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	c.setHeaders(req, opt, contentType)

	// Set auth
	c.setAuth(req, opt)

	// Set timeout if specified
	if opt.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	// Configure redirect policy
	client := c.client
	if opt.AllowRedirect == false {
		client = &http.Client{
			Timeout: c.client.Timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return &Response{Response: resp}, nil
}

func (c *Client) buildURL(reqURL string, params map[string]string) string {
	fullURL := reqURL
	if c.baseURL != "" && !strings.HasPrefix(reqURL, "http") {
		fullURL = c.baseURL + "/" + strings.TrimPrefix(reqURL, "/")
	}

	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return fullURL
		}
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	return fullURL
}

func (c *Client) setHeaders(req *http.Request, opt *RequestOptions, contentType string) {
	// Set default headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Set content type if provided
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Set request-specific headers
	if opt.Headers != nil {
		for k, v := range opt.Headers {
			req.Header.Set(k, v)
		}
	}
}

func (c *Client) setAuth(req *http.Request, opt *RequestOptions) {
	auth := c.auth
	if opt.Auth != nil {
		auth = opt.Auth
	}

	if auth != nil {
		if auth.Token != "" {
			req.Header.Set("Authorization", "Bearer "+auth.Token)
		} else if auth.Username != "" || auth.Password != "" {
			req.SetBasicAuth(auth.Username, auth.Password)
		}
	}
}

// Convenience functions for quick requests

// GET performs a simple GET request
func GET(url string, opts ...*RequestOptions) (*Response, error) {
	return New().GET(url, opts...)
}

// POST performs a simple POST request
func POST(url string, opts ...*RequestOptions) (*Response, error) {
	return New().POST(url, opts...)
}

// PUT performs a simple PUT request
func PUT(url string, opts ...*RequestOptions) (*Response, error) {
	return New().PUT(url, opts...)
}

// DELETE performs a simple DELETE request
func DELETE(url string, opts ...*RequestOptions) (*Response, error) {
	return New().DELETE(url, opts...)
}

// PATCH performs a simple PATCH request
func PATCH(url string, opts ...*RequestOptions) (*Response, error) {
	return New().PATCH(url, opts...)
}

// HEAD performs a simple HEAD request
func HEAD(url string, opts ...*RequestOptions) (*Response, error) {
	return New().HEAD(url, opts...)
}
