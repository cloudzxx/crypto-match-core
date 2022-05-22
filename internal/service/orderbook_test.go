package service

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestOrderBookInsert(t *testing.T) {
	ob := NewOrderBook()

	order := &Order{
		ID:    1,
		Price: decimal.RequireFromString("100.00"),
		Left:  decimal.RequireFromString("1.0"),
	}

	ob.Insert(order)

	if ob.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", ob.Len())
	}
}

func TestOrderBookRemove(t *testing.T) {
	ob := NewOrderBook()

	order := &Order{
		ID:    1,
		Price: decimal.RequireFromString("100.00"),
		Left:  decimal.RequireFromString("1.0"),
	}

	ob.Insert(order)
	ob.Remove(order)

	if ob.Len() != 0 {
		t.Errorf("expected 0 entries after remove")
	}
}

func TestOrderBookBestAsk(t *testing.T) {
	ob := NewOrderBook()

	ob.Insert(&Order{
		ID:    1,
		Price: decimal.RequireFromString("100.00"),
		Left:  decimal.RequireFromString("1.0"),
	})

	ob.Insert(&Order{
		ID:    2,
		Price: decimal.RequireFromString("99.00"),
		Left:  decimal.RequireFromString("0.5"),
	})

	best := ob.BestAsk()
	if best == nil {
		t.Fatal("expected best ask order")
	}
	if best.Price.Cmp(decimal.RequireFromString("99.00")) != 0 {
		t.Errorf("expected price 99.00, got %s", best.Price)
	}
}

func TestOrderBookBestBid(t *testing.T) {
	ob := NewOrderBook()

	ob.Insert(&Order{
		ID:    1,
		Price: decimal.RequireFromString("100.00"),
		Left:  decimal.RequireFromString("1.0"),
	})

	ob.Insert(&Order{
		ID:    2,
		Price: decimal.RequireFromString("99.00"),
		Left:  decimal.RequireFromString("0.5"),
	})

	best := ob.BestBid()
	if best == nil {
		t.Fatal("expected best bid order")
	}
	if best.Price.Cmp(decimal.RequireFromString("100.00")) != 0 {
		t.Errorf("expected price 100.00, got %s", best.Price)
	}
}