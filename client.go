package hue

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

const (
	defaultBasePath = "api/"
	userAgent       = "go-hue"
	discoveryUrl    = "https://discovery.meethue.com/"
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

type discoverResponse struct {
	ID   string `json:"id"`
	Host string `json:"internalipaddress"`
}

type createUserRequest struct {
	DeviceType string `json:"devicetype"`
}

type Client struct {
	client  *http.Client
	baseURL *url.URL

	userAgent string
	clientId  string // username for hue bridge
	logger    logr.Logger
	common    service

	Lights *LightService
	Groups *GroupService
}

type service struct {
	client *Client
}

type ClientOptions struct {
	HttpClient *http.Client
	LogLevel   logrus.Level
}

// Discover gets hue bridge host address
func Discover() (string, error) {
	req, err := http.NewRequest(http.MethodGet, discoveryUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	discoveryResponses := new([]discoverResponse)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("invalid status code returned")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bytes, discoveryResponses)
	if err != nil {
		return "", err
	}

	if discoveryResponses == nil || len(*discoveryResponses) == 0 {
		return "", errors.New("no bridge found on your network")
	}

	// Use the first bridge on the network
	return (*discoveryResponses)[0].Host, nil
}

func newClient(host string, opts *ClientOptions) (*Client, error) {
	var httpClient *http.Client
	if opts == nil || opts.HttpClient == nil {
		httpClient = &http.Client{} // Create default client
	} else {
		httpClient = opts.HttpClient
	}

	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	c := &Client{client: httpClient, baseURL: u, userAgent: userAgent}
	c.logger = logrusr.NewLogger(logrus.New())
	c.common.client = c

	c.Lights = (*LightService)(&c.common)
	c.Groups = (*GroupService)(&c.common)

	return c, nil
}

func NewClient(host, clientId string, opts *ClientOptions) *Client {
	c, err := newClient(fmt.Sprintf("http://%v/api/", host), opts)
	if err != nil {
		c.logger.Error(err, "Couldn't create client")
		return nil
	}
	c.clientId = clientId

	return c
}

// CreateUser creates local user on the bridge and returns authenticated client instance
// Don't forget to press bridge button otherwise it will fail
func CreateUser(host, deviceType string, opts *ClientOptions) (*Client, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(&createUserRequest{DeviceType: deviceType})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/api", host), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	apiResponses := new([]ApiResponse)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status code returned")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, apiResponses)
	if err != nil {
		return nil, err
	}

	if apiResponses == nil || len(*apiResponses) == 0 || (*apiResponses)[0].Error != nil {
		return nil, errors.New((*apiResponses)[0].Error.Description)
	}

	clientId := (*apiResponses)[0].Success["username"].(string)

	return NewClient(host, clientId, opts), nil
}

// GetHost returns ip address of hue bridge
func (c *Client) GetHost() string {
	return c.baseURL.Host
}

// GetClientID returns clientID of current client
func (c *Client) GetClientID() string {
	return c.clientId
}

func (c *Client) newRequest(method, url string, payload interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.baseURL)
	}
	u, err := c.baseURL.Parse(url)
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
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return &Response{Response: resp}, err
	}
	defer resp.Body.Close()

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

func (c *Client) path(service string, params ...string) string {
	if c.clientId == "" {
		c.logger.Info("clientId is missing")
	}
	if len(params) == 0 {
		return fmt.Sprintf("%v/%v", c.clientId, service)
	} else if len(params) == 1 {
		return fmt.Sprintf("%v/%v/%v", c.clientId, service, params[0])
	}

	return fmt.Sprintf("%v/%v/%v/%v", c.clientId, service, params[0], params[1])
}
