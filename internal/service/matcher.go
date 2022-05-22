package service

import (
	"fmt"
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
// 所有订单操作（下单、撤单）均应通过此锁串行化
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

// Lock 暴露锁用于外部协调订单生命周期
func (m *Matcher) Lock()   { m.mu.Lock() }
func (m *Matcher) Unlock() { m.mu.Unlock() }

// ProcessOrder 在撮合锁内完成完整流程：入簿 → 撮合 → 结算 → 清理
func (m *Matcher) ProcessOrder(mkt *Market, order *Order) ([]MatchResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if order.Side == SideAsk {
		mkt.Asks.Insert(order)
	} else {
		mkt.Bids.Insert(order)
	}

	var results []MatchResult
	var err error
	if order.Side == SideBid {
		results, err = m.executeLimitBid(mkt, order)
	} else {
		results, err = m.executeLimitAsk(mkt, order)
	}
	if err != nil {
		return nil, err
	}

	for i := range results {
		if err := settleMatchResult(&results[i], mkt, m.bm); err != nil {
			return nil, err
		}
	}
	cleanupFilledOrders(mkt, results, order)
	return results, nil
}

// CancelOrder 在撮合锁内完成撤单：校验 → 移除 → 解冻
func (m *Matcher) CancelOrder(mkt *Market, orderID uint64, userID uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	order := mkt.Asks.Get(orderID)
	if order == nil {
		order = mkt.Bids.Get(orderID)
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}
	if order.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	if order.Side == SideAsk {
		mkt.Asks.Remove(order)
	} else {
		mkt.Bids.Remove(order)
	}
	if order.Side == SideBid {
		return m.bm.Unfreeze(userID, mkt.Money, order.Freeze)
	}
	return m.bm.Unfreeze(userID, mkt.Stock, order.Freeze)
}

// settleMatchResult 成交结算 + 费用扣收
func settleMatchResult(r *MatchResult, mkt *Market, bm *BalanceManager) error {
	if r.TakerOrder.Side == SideAsk {
		// 卖单吃单：释出 base（Stock），收入 quote（Money）扣手续费
		if err := bm.SettleTrade(r.TakerOrder.UserID, mkt.Stock, decimal.Zero, r.DealAmount.Neg()); err != nil {
			return err
		}
		if err := bm.SettleTrade(r.TakerOrder.UserID, mkt.Money, r.DealMoney.Sub(r.TakerFee), decimal.Zero); err != nil {
			return err
		}
	} else {
		// 买单吃单：释出 quote（Money）扣手续费，收入 base（Stock）
		if err := bm.SettleTrade(r.TakerOrder.UserID, mkt.Money, r.TakerFee.Neg(), r.DealMoney.Neg()); err != nil {
			return err
		}
		if err := bm.SettleTrade(r.TakerOrder.UserID, mkt.Stock, r.DealAmount, decimal.Zero); err != nil {
			return err
		}
	}
	if r.MakerOrder.Side == SideAsk {
		if err := bm.SettleTrade(r.MakerOrder.UserID, mkt.Stock, decimal.Zero, r.DealAmount.Neg()); err != nil {
			return err
		}
		if err := bm.SettleTrade(r.MakerOrder.UserID, mkt.Money, r.DealMoney.Add(r.MakerRebate), decimal.Zero); err != nil {
			return err
		}
	} else {
		if err := bm.SettleTrade(r.MakerOrder.UserID, mkt.Money, r.MakerRebate, r.DealMoney.Neg()); err != nil {
			return err
		}
		if err := bm.SettleTrade(r.MakerOrder.UserID, mkt.Stock, r.DealAmount, decimal.Zero); err != nil {
			return err
		}
	}
	return nil
}

// cleanupFilledOrders 移除已完全成交的订单
func cleanupFilledOrders(mkt *Market, results []MatchResult, taker *Order) {
	for _, r := range results {
		if r.MakerOrder.Left.IsZero() {
			if r.MakerOrder.Side == SideAsk {
				mkt.Asks.Remove(r.MakerOrder)
			} else {
				mkt.Bids.Remove(r.MakerOrder)
			}
		}
	}
	if taker.Left.IsZero() {
		if taker.Side == SideAsk {
			mkt.Asks.Remove(taker)
		} else {
			mkt.Bids.Remove(taker)
		}
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
		results, err = m.executeLimitBid(marketData, order)
	} else {
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
