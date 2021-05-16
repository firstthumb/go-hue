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

func NewClient(httpClient *http.Client, host, username string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(fmt.Sprintf("http://%s/%s", host, defaultBasePath))

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent, Username: username}
	c.logger = logrusr.NewLogger(logrus.New())
	c.common.client = c
	c.Verbose = true
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

	body, _ := httputil.DumpRequest(req, true)
	c.logger.Info(fmt.Sprintf("%s", string(body)))

	resp, err := c.client.Do(req)
	if err != nil {
		return &Response{Response: resp}, err
	}
	defer resp.Body.Close()

	body, _ = httputil.DumpResponse(resp, true)
	c.logger.Info(fmt.Sprintf("%s", string(body)))

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

type ApiResponse struct {
	Success map[string]interface{} `json:"success,omitempty"`
	Error   *ApiError              `json:"error,omitempty"`
}

type ApiError struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

func Bool(v bool) *bool { return &v }

func Int(v int) *int { return &v }

func UInt8(v uint8) *uint8 { return &v }

func Int64(v int64) *int64 { return &v }

func String(v string) *string { return &v }
