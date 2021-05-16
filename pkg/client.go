package hue

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

const (
	defaultBasePath = "api/"
	userAgent       = "gohue"
)

type Response struct {
	*http.Response
}

// ApiResponse that Hue returns
type ApiResponse struct {
	Success map[string]interface{} `json:"success,omitempty"`
	Error   *ApiError              `json:"error,omitempty"`
}

// ErrorResponse that Hue returns
type ApiError struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

type DiscoverResponse struct {
	ID   string `json:"id"`
	Host string `json:"internalipaddress"`
}

type Client struct {
	client *http.Client

	BaseURL *url.URL

	UserAgent string

	Username string

	Verbose bool

	common service

	logger logr.Logger

	User  *UserService
	Light *LightService
}

type service struct {
	client *Client
}

type ClientOptions struct {
	HttpClient *http.Client
}

// Discover gets unauthorized client instance
// You need to login to make authorized service call
func Discover() (*Client, error) {
	c := &Client{client: &http.Client{}, UserAgent: userAgent}

	req, err := c.NewRequest(http.MethodGet, "https://discovery.meethue.com/", nil)
	if err != nil {
		return nil, err
	}

	discoveryResponses := new([]*DiscoverResponse)
	resp, err := c.Do(context.Background(), req, discoveryResponses)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid status code returned")
	}

	if discoveryResponses == nil || len(*discoveryResponses) == 0 {
		return nil, errors.New("No bridge found on your network")
	}

	// Use the first bridge on the network
	baseURL, _ := url.Parse(fmt.Sprintf("http://%s/%s", (*discoveryResponses)[0].Host, defaultBasePath))
	c.BaseURL = baseURL
	c.logger = logrusr.NewLogger(logrus.New()) // TODO: Make it configurable
	c.common.client = c
	c.Verbose = true // TODO: Make it configurable
	c.User = (*UserService)(&c.common)
	c.Light = (*LightService)(&c.common)

	return c, nil
}

func NewClient(host, username string, opts *ClientOptions) *Client {
	var httpClient *http.Client
	if opts == nil || opts.HttpClient == nil {
		httpClient = &http.Client{} // Create default client
	} else {
		httpClient = opts.HttpClient
	}

	baseURL, _ := url.Parse(fmt.Sprintf("http://%s/%s", host, defaultBasePath))

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent, Username: username}
	c.logger = logrusr.NewLogger(logrus.New()) // TODO: Make it configurable
	c.common.client = c
	c.Verbose = true // TODO: Make it configurable
	c.User = (*UserService)(&c.common)
	c.Light = (*LightService)(&c.common)

	return c
}

func (c *Client) NewRequest(method, url string, payload interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if payload != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)

	if c.Verbose {
		body, _ := httputil.DumpRequest(req, true)
		c.logger.Info(string(body))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return &Response{Response: resp}, err
	}
	defer resp.Body.Close()

	if c.Verbose {
		body, _ := httputil.DumpResponse(resp, true)
		c.logger.Info(string(body))
	}

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return &Response{Response: resp}, err
}

func (c *Client) Login(username string) *Client {
	c.Username = username

	// TODO: Make a request to verify user login

	return c
}

func path(service string, params ...string) string {
	if len(params) == 1 {
		return fmt.Sprintf("%v/%v", params[0], service)
	} else if len(params) == 2 {
		return fmt.Sprintf("%v/%v/%v", params[0], service, params[1])
	}

	return fmt.Sprintf("%v/%v/%v/%v", params[0], service, params[1], params[2])
}
