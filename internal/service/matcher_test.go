package service

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestMatchLimitAskSamePrice(t *testing.T) {
	matcher := NewMatcher(nil, nil)

	mkt := &Market{
		Name:  "BTCUSDT",
		Stock: "BTC",
		Money: "USDT",
	}

	bids := NewOrderBook()
	mkt.Bids = bids

	makerBid := &Order{
		ID:     1,
		Side:   SideBid,
		Price:  decimal.RequireFromString("100.00"),
		Amount: decimal.RequireFromString("1.0"),
		Left:   decimal.RequireFromString("1.0"),
	}
	bids.Insert(makerBid)

	takerAsk := &Order{
		ID:     2,
		Side:   SideAsk,
		Price:  decimal.RequireFromString("100.00"),
		Amount: decimal.RequireFromString("0.5"),
		Left:   decimal.RequireFromString("0.5"),
		Market: "BTCUSDT",
	}

	results, err := matcher.executeLimitAsk(mkt, takerAsk)
	if err != nil {
		t.Fatalf("match error: %v", err)
	}

	if len(results) == 0 {
		t.Errorf("expected match at same price")
	}
}

func TestMatchLimitBidSamePrice(t *testing.T) {
	matcher := NewMatcher(nil, nil)

	mkt := &Market{Name: "BTCUSDT"}
	mkt.Asks = NewOrderBook()
	mkt.Bids = NewOrderBook()

	makerAsk := &Order{
		ID:    1,
		Side:  SideAsk,
		Price: decimal.RequireFromString("100.00"),
		Left:  decimal.RequireFromString("1.0"),
	}
	mkt.Asks.Insert(makerAsk)

	takerBid := &Order{
		ID:     2,
		Side:   SideBid,
		Price:  decimal.RequireFromString("100.00"),
		Amount: decimal.RequireFromString("0.5"),
		Left:   decimal.RequireFromString("0.5"),
		Market: "BTCUSDT",
	}

	results, _ := matcher.executeLimitBid(mkt, takerBid)
	if len(results) != 1 {
		t.Errorf("expected 1 match at same price, got %d", len(results))
	}
}

func TestMatchPartialFill(t *testing.T) {
	matcher := NewMatcher(nil, nil)

	mkt := &Market{Name: "BTCUSDT"}
	mkt.Asks = NewOrderBook()
	mkt.Bids = NewOrderBook()

	makerAsk := &Order{
		ID:    1,
		Side:  SideAsk,
		Price: decimal.RequireFromString("100.00"),
		Left:  decimal.RequireFromString("1.0"),
	}
	mkt.Asks.Insert(makerAsk)

	takerBid := &Order{
		ID:     2,
		Side:   SideBid,
		Price:  decimal.RequireFromString("100.00"),
		Amount: decimal.RequireFromString("0.5"),
		Left:   decimal.RequireFromString("0.5"),
		Market: "BTCUSDT",
	}

	results, _ := matcher.executeLimitBid(mkt, takerBid)
	if len(results) != 1 {
		t.Fatalf("expected 1 match, got %d", len(results))
	}

	if !results[0].DealAmount.Equal(decimal.RequireFromString("0.5")) {
		t.Errorf("deal amount mismatch: expected 0.5, got %s", results[0].DealAmount)
	}
}