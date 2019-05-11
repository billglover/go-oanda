package oanda

import (
	"net/http"
	"net/url"
)

const (
	demoAPI       string = "https://api-fxpractice.oanda.com"
	liveAPI       string = "https://api-fxtrade.oanda.com"
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
