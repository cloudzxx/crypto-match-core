package service

import (
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
)

type BalanceType uint8

const (
	BalanceTypeAvailable BalanceType = iota
	BalanceTypeFreeze
)

type Balance struct {
	Available decimal.Decimal
	Freeze    decimal.Decimal
}

type BalanceManager struct {
	balances map[uint32]map[string]*Balance
	mu       sync.RWMutex
}

func NewBalanceManager() *BalanceManager {
	return &BalanceManager{
		balances: make(map[uint32]map[string]*Balance),
	}
}

func (bm *BalanceManager) GetBalance(userID uint32, asset string) *Balance {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if userBalances, ok := bm.balances[userID]; ok {
		if balance, ok := userBalances[asset]; ok {
			return &Balance{
				Available: balance.Available,
				Freeze:    balance.Freeze,
			}
		}
	}

	return &Balance{
		Available: decimal.Zero,
		Freeze:    decimal.Zero,
	}
}

func (bm *BalanceManager) Freeze(userID uint32, asset string, amount decimal.Decimal) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if _, ok := bm.balances[userID]; !ok {
		bm.balances[userID] = make(map[string]*Balance)
	}

	if _, ok := bm.balances[userID][asset]; !ok {
		bm.balances[userID][asset] = &Balance{
			Available: decimal.Zero,
			Freeze:    decimal.Zero,
		}
	}

	balance := bm.balances[userID][asset]
	if balance.Available.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance")
	}

	balance.Available = balance.Available.Sub(amount)
	balance.Freeze = balance.Freeze.Add(amount)
	return nil
}

func (bm *BalanceManager) Unfreeze(userID uint32, asset string, amount decimal.Decimal) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if _, ok := bm.balances[userID]; !ok {
		return fmt.Errorf("balance not found")
	}

	if _, ok := bm.balances[userID][asset]; !ok {
		return fmt.Errorf("balance not found")
	}

	balance := bm.balances[userID][asset]
	if balance.Freeze.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient frozen balance")
	}

	balance.Freeze = balance.Freeze.Sub(amount)
	balance.Available = balance.Available.Add(amount)
	return nil
}

func (bm *BalanceManager) SettleTrade(userID uint32, asset string, availableDelta, freezeDelta decimal.Decimal) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if _, ok := bm.balances[userID]; !ok {
		bm.balances[userID] = make(map[string]*Balance)
	}

	if _, ok := bm.balances[userID][asset]; !ok {
		bm.balances[userID][asset] = &Balance{
			Available: decimal.Zero,
			Freeze:    decimal.Zero,
		}
	}

	balance := bm.balances[userID][asset]
	newAvailable := balance.Available.Add(availableDelta)
	newFreeze := balance.Freeze.Add(freezeDelta)
	if newAvailable.LessThan(decimal.Zero) || newFreeze.LessThan(decimal.Zero) {
		return fmt.Errorf("balance would go negative")
	}
	balance.Available = newAvailable
	balance.Freeze = newFreeze
	return nil
}

func (bm *BalanceManager) SetBalance(userID uint32, asset string, available, freeze decimal.Decimal) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if _, ok := bm.balances[userID]; !ok {
		bm.balances[userID] = make(map[string]*Balance)
	}

	bm.balances[userID][asset] = &Balance{
		Available: available,
		Freeze:    freeze,
	}
}

func (bm *BalanceManager) GetAllBalances(userID uint32) map[string]*Balance {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	result := make(map[string]*Balance)
	if userBalances, ok := bm.balances[userID]; ok {
		for asset, balance := range userBalances {
			result[asset] = &Balance{
				Available: balance.Available,
				Freeze:    balance.Freeze,
			}
		}
	}
	return result
}
