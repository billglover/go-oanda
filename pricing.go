package oanda

import (
	"context"
	"strings"
	"time"
)

// Price is an Account-specific Price.
type Price struct {
	Type        string           `json:"type"`
	Time        time.Time        `json:"time"`
	Bids        []PriceLiquidity `json:"bids"`
	Asks        []PriceLiquidity `json:"asks"`
	CloseoutBid string           `json:"closeoutBid"`
	CloseoutAsk string           `json:"closeoutAsk"`
	Tradeable   bool             `json:"tradeable"`
	Instrument  string           `json:"instrument"`
}

// PriceLiquidity is a list of prices and liquidity available on an Instrument.
// It is possible for this list to be empty if there is liquidity currently
// available for the Instrument in the Account.
type PriceLiquidity struct {
	Price     string `json:"price"`
	Liquidity int    `json:"liquidity"`
}

// pricingResponse contains a slice of prices and an associated timestamp.
type pricingResponse struct {
	Time   time.Time `json:"time"`
	Prices []Price   `json:"prices"`
}

// Pricing returns pricing information for a specified list of instruments within an Account.
func (c *Client) Pricing(ctx context.Context, actID string, instruments []string) ([]Price, error) {
	req, err := c.NewRequest("GET", "/v3/accounts/"+actID+"/pricing", nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Set("instruments", strings.Join(instruments, ","))
	req.URL.RawQuery = params.Encode()

	pricing := pricingResponse{}
	err = c.Do(ctx, req, &pricing)
	if err != nil {
		return nil, err
	}

	return pricing.Prices, nil
}
