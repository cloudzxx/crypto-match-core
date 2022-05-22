package service

import (
	"sync"

	"github.com/shopspring/decimal"
)

var (
	takerFeeRate    = decimal.NewFromFloat(0.001)
	makerRebateRate = decimal.NewFromFloat(0.0005)
)

type MatchResult struct {
	TakerOrder   *Order
	MakerOrder   *Order
	DealPrice    decimal.Decimal
	DealAmount   decimal.Decimal
	DealMoney    decimal.Decimal
	TakerFee     decimal.Decimal
	MakerRebate  decimal.Decimal
}

type Matcher struct {
	mm       *Manager
	bm       *BalanceManager
	mu       sync.Mutex
}

func NewMatcher(mm *Manager, bm *BalanceManager) *Matcher {
	return &Matcher{
		mm: mm,
		bm: bm,
	}
}

func (m *Matcher) Match(order *Order) ([]MatchResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	marketData, err := m.mm.Get(order.Market)
	if err != nil {
		return nil, err
	}

	var results []MatchResult

	if order.Side == SideBid {
		results, err = m.executeLimitBid(marketData, order)
	} else {
		results, err = m.executeLimitAsk(marketData, order)
	}

	return results, err
}

func (m *Matcher) executeLimitAsk(mkt *Market, taker *Order) ([]MatchResult, error) {
	var results []MatchResult

	takerLeft := taker.Left

	bestPrice := decimal.Zero

	mkt.Bids.Range(func(price decimal.Decimal, orders []*Order) bool {
		if bestPrice.IsZero() {
			bestPrice = price
		}

		if taker.Price.Cmp(price) < 0 {
			return false
		}

		for _, maker := range orders {
			if takerLeft.IsZero() {
				return false
			}

			fillAmount := decimal.Min(takerLeft, maker.Left)

			dealMoney := fillAmount.Mul(maker.Price)
			takerFee := dealMoney.Mul(takerFeeRate)
			makerRebate := dealMoney.Mul(makerRebateRate)

			result := MatchResult{
				TakerOrder:  taker,
				MakerOrder:  maker,
				DealPrice:   maker.Price,
				DealAmount:  fillAmount,
				DealMoney:   dealMoney,
				TakerFee:    takerFee,
				MakerRebate: makerRebate,
			}
			results = append(results, result)

			takerLeft = takerLeft.Sub(fillAmount)
			maker.Left = maker.Left.Sub(fillAmount)
		}
		return true
	})

	if len(results) > 0 {
		taker.Left = takerLeft
	}

	return results, nil
}

func (m *Matcher) executeLimitBid(mkt *Market, taker *Order) ([]MatchResult, error) {
	var results []MatchResult

	takerLeft := taker.Left

	mkt.Asks.Range(func(price decimal.Decimal, orders []*Order) bool {
		if taker.Price.Cmp(price) > 0 {
			return false
		}

		for _, maker := range orders {
			if takerLeft.IsZero() {
				return false
			}

			fillAmount := decimal.Min(takerLeft, maker.Left)

			dealMoney := fillAmount.Mul(maker.Price)
			takerFee := dealMoney.Mul(takerFeeRate)
			makerRebate := dealMoney.Mul(makerRebateRate)

			result := MatchResult{
				TakerOrder:  taker,
				MakerOrder:  maker,
				DealPrice:   maker.Price,
				DealAmount:  fillAmount,
				DealMoney:   dealMoney,
				TakerFee:    takerFee,
				MakerRebate: makerRebate,
			}
			results = append(results, result)

			takerLeft = takerLeft.Sub(fillAmount)
			maker.Left = maker.Left.Sub(fillAmount)
		}
		return true
	})

	if len(results) > 0 {
		taker.Left = takerLeft
	}

	return results, nil
}
