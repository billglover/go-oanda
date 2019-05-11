package oanda_test

import (
	"testing"

	"github.com/billglover/go-oanda"
)

func TestNewClient(t *testing.T) {
	c := oanda.NewClient(nil)

	if c == nil {
		t.Fatal("should return a client: got nil")
	}
}
