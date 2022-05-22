package service

import (
	"errors"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/cloudzxx/crypto-match-core/internal/config"
	mdecimal "github.com/cloudzxx/crypto-match-core/pkg/decimal"
)

var (
	ErrMarketNotFound = errors.New("market not found")
	ErrMarketExists   = errors.New("market already exists")
)

type Market struct {
	Name       string
	Stock      string
	Money      string
	StockPrec  int
	MoneyPrec  int
	MinAmount  decimal.Decimal
	Asks       *OrderBook
	Bids       *OrderBook
	mu         sync.RWMutex
}

type Manager struct {
	markets map[string]*Market
	mu      sync.RWMutex
}

func NewManager(cfg []config.MarketConfig) *Manager {
	m := &Manager{
		markets: make(map[string]*Market),
	}

	for _, c := range cfg {
		market := &Market{
			Name:      c.Name,
			Stock:     c.Stock,
			Money:     c.Money,
			StockPrec: c.StockPrec,
			MoneyPrec: c.MoneyPrec,
			MinAmount: mdecimal.MustNewFromString(c.MinAmount),
			Asks:      NewOrderBook(),
			Bids:      NewOrderBook(),
		}
		m.markets[c.Name] = market
	}

	return m
}

func (m *Manager) Get(name string) (*Market, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	market, ok := m.markets[name]
	if !ok {
		return nil, ErrMarketNotFound
	}
	return market, nil
}

func (m *Manager) List() []*Market {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Market, 0, len(m.markets))
	for _, market := range m.markets {
		result = append(result, market)
	}
	return result
}

func (m *Manager) Add(cfg config.MarketConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.markets[cfg.Name]; ok {
		return ErrMarketExists
	}

	market := &Market{
		Name:      cfg.Name,
		Stock:     cfg.Stock,
		Money:     cfg.Money,
		StockPrec: cfg.StockPrec,
		MoneyPrec: cfg.MoneyPrec,
		MinAmount: mdecimal.MustNewFromString(cfg.MinAmount),
		Asks:      NewOrderBook(),
		Bids:      NewOrderBook(),
	}
	m.markets[cfg.Name] = market
	return nil
}