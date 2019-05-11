package oanda

import (
	"context"
)

// accountsResponse contains a slice of AccountIdentifiers.
type accountsResponse struct {
	Accounts []AccountIdentifier `json:"accounts"`
}

// AccountIdentifier contains an Account ID and associated tags.
type AccountIdentifier struct {
	ID   string        `json:"id"`
	Tags []interface{} `json:"tags"`
}

// Accounts returns a list of all Accounts.
func (c *Client) Accounts(ctx context.Context) ([]AccountIdentifier, error) {
	req, err := c.NewRequest("GET", "/v3/accounts", nil)
	if err != nil {
		return nil, err
	}

	acts := accountsResponse{}
	err = c.Do(ctx, req, &acts)
	if err != nil {
		return nil, err
	}

	return acts.Accounts, nil
}
