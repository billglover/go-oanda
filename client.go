package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	demoAPI       string = "https://api-fxpractice.oanda.com/"
	liveAPI       string = "https://api-fxtrade.oanda.com/"
	demoStreamAPI string = "https://stream-fxpractice.oanda.com/"
	liveStreamAPI string = "https://stream-fxtrade.oanda.com/"
	userAgent     string = "github.com/billglover/go-oanda"
)

// Client holds configuration and state for the OANDA API client.
type Client struct {
	baseAPIURL    *url.URL
	baseStreamURL *url.URL

	userAgent string
	client    *http.Client
}

// NewClient returns a new OANDA client. If a nil httpClient is provided,
// http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication.
// For example: golang.org/x/oauth2
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = http.DefaultClient
	}

	baseAPIURL, err := url.Parse(demoAPI)
	if err != nil {
		panic(err)
	}

	baseStreamURL, err := url.Parse(demoStreamAPI)
	if err != nil {
		panic(err)
	}

	client := &Client{
		baseAPIURL:    baseAPIURL,
		baseStreamURL: baseStreamURL,
		userAgent:     userAgent,
		client:        c,
	}
	return client
}

// NewRequest creates an HTTP Request. The client baseURL is checked to confirm that it has a trailing
// slash. A relative URL should be provided without the leading slash. If a non-nil body is provided
// it will be JSON encoded and included in the request.
// Inspiration: https://github.com/google/go-github/blob/master/github/github.go
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	return c.newRequest(method, c.baseAPIURL, path, body)
}

// NewStreamRequest creates an HTTP Request. The client baseURL is checked to confirm that it has a trailing
// slash. A relative URL should be provided without the leading slash. If a non-nil body is provided
// it will be JSON encoded and included in the request.
// Inspiration: https://github.com/google/go-github/blob/master/github/github.g
func (c *Client) newStreamRequest(method, path string, body interface{}) (*http.Request, error) {
	return c.newRequest(method, c.baseStreamURL, path, body)
}

func (c *Client) newRequest(method string, baseURL *url.URL, path string, body interface{}) (*http.Request, error) {
	if strings.HasSuffix(baseURL.Path, "/") == false {
		fmt.Println(baseURL.Path)
		return nil, fmt.Errorf("client baseURL does not have a trailing slash: %q", baseURL)
	}

	u, err := baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}

// Do sends a request and returns the response. An error is returned if the request cannot
// be sent or if the API returns an error. If a response is received, the body response body
// is decoded and stored in the value pointed to by v.
// Inspiration: https://github.com/google/go-github/blob/master/github/github.go
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	// TODO: remove deb
	fmt.Println(req)
	fmt.Println(string(data))

	// Anything other than a HTTP 2xx response code is treated as an error.
	if c := resp.StatusCode; c >= 300 {
		return fmt.Errorf(resp.Status)
	}

	if v != nil && len(data) != 0 {
		err = json.Unmarshal(data, v)

		switch err {
		case nil:
		case io.EOF:
			return nil
		default:
			return fmt.Errorf("unable to parse response body")
		}
	}

	return nil
}
