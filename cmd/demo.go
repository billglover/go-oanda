package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/billglover/go-oanda"
	"golang.org/x/oauth2"
)

func main() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("OANDA_TOKEN")})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	api := oanda.NewClient(tc)

	acts, err := api.Accounts(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	actID := acts[0].ID

	prices, err := api.Pricing(context.Background(), actID, []string{"GBP_USD", "EUR_GBP", "EUR_USD"})
	if err != nil {
		log.Fatal(err)
	}

	for i := range prices {
		fmt.Println(prices[i])
	}
}
