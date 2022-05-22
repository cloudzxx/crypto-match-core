package service

import (
	"sync"

	"github.com/shopspring/decimal"
)

// 吃单方费率 0.1%，挂单方返佣 0.05%
var (
	takerFeeRate    = decimal.NewFromFloat(0.001)
	makerRebateRate = decimal.NewFromFloat(0.0005)
)

// 一次撮合的结果：吃单与挂单在某个价格上成交了一定数量
type MatchResult struct {
	TakerOrder   *Order
	MakerOrder   *Order
	DealPrice    decimal.Decimal
	DealAmount   decimal.Decimal
	DealMoney    decimal.Decimal
	TakerFee     decimal.Decimal
	MakerRebate  decimal.Decimal
}

// 撮合引擎：价格优先 → 时间优先
// 全市场共享一把锁，保证撮合串行化
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

// Match 是撮合入口：根据买卖方向委托到不同的执行路径
func (m *Matcher) Match(order *Order) ([]MatchResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	marketData, err := m.mm.Get(order.Market)
	if err != nil {
		return nil, err
	}

	var results []MatchResult

	if order.Side == SideBid {
		// 买单：从卖一价开始逐级向上吃
		results, err = m.executeLimitBid(marketData, order)
	} else {
		// 卖单：从买一价开始逐级向下吃
		results, err = m.executeLimitAsk(marketData, order)
	}

	return results, err
}

// executeLimitAsk 卖单限价撮合
// 流程：遍历买单队列（从最高价开始），价格 >= taker.Price 即可成交
func (m *Matcher) executeLimitAsk(mkt *Market, taker *Order) ([]MatchResult, error) {
	var results []MatchResult

	takerLeft := taker.Left

	bestPrice := decimal.Zero

	mkt.Bids.Range(func(price decimal.Decimal, orders []*Order) bool {
		if bestPrice.IsZero() {
			bestPrice = price
		}

		// 买方价格低于卖方心理价，停止撮合
		if taker.Price.Cmp(price) < 0 {
			return false
		}

		// 同价位按时间顺序逐个吃单
		for _, maker := range orders {
			if takerLeft.IsZero() {
				return false
			}

			// 取吃单剩余量与挂单剩余量的较小值
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
			maker.Left = maker.Left.Sub(fillAmount) // 挂单剩余减少
		}
		return true
	})

	if len(results) > 0 {
		taker.Left = takerLeft
	}

	return results, nil
}

// executeLimitBid 买单限价撮合
// 流程：遍历卖单队列（从最低价开始），价格 <= taker.Price 即可成交
func (m *Matcher) executeLimitBid(mkt *Market, taker *Order) ([]MatchResult, error) {
	var results []MatchResult

	takerLeft := taker.Left

	mkt.Asks.Range(func(price decimal.Decimal, orders []*Order) bool {
		// 卖方价格高于买方心理价，停止撮合
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
