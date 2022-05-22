package service

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestBalanceFreeze(t *testing.T) {
	bm := NewBalanceManager()

	bm.SetBalance(1, "BTC", decimal.NewFromFloat(10.0), decimal.Zero)

	err := bm.Freeze(1, "BTC", decimal.NewFromFloat(5.0))
	if err != nil {
		t.Fatalf("freeze failed: %v", err)
	}

	balance := bm.GetBalance(1, "BTC")
	if balance.Available.Cmp(decimal.NewFromFloat(5.0)) != 0 {
		t.Errorf("expected available 5.0, got %s", balance.Available)
	}
	if balance.Freeze.Cmp(decimal.NewFromFloat(5.0)) != 0 {
		t.Errorf("expected freeze 5.0, got %s", balance.Freeze)
	}
}

func TestBalanceUnfreeze(t *testing.T) {
	bm := NewBalanceManager()

	bm.SetBalance(1, "BTC", decimal.NewFromFloat(5.0), decimal.NewFromFloat(5.0))

	err := bm.Unfreeze(1, "BTC", decimal.NewFromFloat(3.0))
	if err != nil {
		t.Fatalf("unfreeze failed: %v", err)
	}

	balance := bm.GetBalance(1, "BTC")
	if balance.Available.Cmp(decimal.NewFromFloat(8.0)) != 0 {
		t.Errorf("expected available 8.0, got %s", balance.Available)
	}
}

func TestInsufficientBalance(t *testing.T) {
	bm := NewBalanceManager()

	bm.SetBalance(1, "BTC", decimal.NewFromFloat(3.0), decimal.Zero)

	err := bm.Freeze(1, "BTC", decimal.NewFromFloat(5.0))
	if err == nil {
		t.Errorf("expected insufficient balance error")
	}
}