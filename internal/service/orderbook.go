package service

import (
	"sync"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

// 订单：价格优先 → 时间优先排序的基础单位
type Order struct {
	ID        uint64
	Type      OrderType
	Side      Side
	UserID    uint32
	Market    string
	Price     decimal.Decimal
	Amount    decimal.Decimal
	Left      decimal.Decimal          // 未成交数量（撮合时递减）
	Freeze    decimal.Decimal          // 已冻结的资产数量
	Source    string
	CreatedAt int64
	UpdatedAt int64
}

type OrderType uint8

const (
	OrderTypeLimit  OrderType = iota // 限价单
	OrderTypeMarket                  // 市价单（预留）
)

type Side uint8

const (
	SideAsk Side = iota // 卖单
	SideBid             // 买单
)

// 同价格档位的订单列表（按到达时间排序）
type PriceLevel struct {
	Price  decimal.Decimal
	Orders []*Order
}

// 基于红黑树的订单簿
// - Asks 升序排列（最低卖价在最前）
// - Bids 降序排列（最高买价在最前）
// - orderIndex 提供 O(1) 的订单查询
type OrderBook struct {
	tree       *redblacktree.Tree
	orderIndex map[uint64]*Order
	mu         sync.RWMutex
}

// 红黑树按价格升序排序；买单遍历时从右向左（降序），卖单从左向右（升序）
func NewOrderBook() *OrderBook {
	return &OrderBook{
		tree: redblacktree.NewWith(func(a, b interface{}) int {
			return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
		}),
		orderIndex: make(map[uint64]*Order),
	}
}

// 插入订单：写入 orderIndex + 追加到对应价格档位
func (ob *OrderBook) Insert(order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.orderIndex[order.ID] = order

	price := order.Price
	node, ok := ob.tree.Get(price)
	if !ok {
		priceKey := price
		node = &PriceLevel{
			Price:  priceKey,
			Orders: []*Order{},
		}
		ob.tree.Put(priceKey, node)
	}

	pl := node.(*PriceLevel)
	pl.Orders = append(pl.Orders, order)
}

// 撤销订单：从 orderIndex 和价格档位中移除，空档位自动清理
func (ob *OrderBook) Remove(order *Order) bool {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	delete(ob.orderIndex, order.ID)

	price := order.Price
	node, ok := ob.tree.Get(price)
	if !ok {
		return false
	}

	pl := node.(*PriceLevel)
	for i, o := range pl.Orders {
		if o.ID == order.ID {
			pl.Orders = append(pl.Orders[:i], pl.Orders[i+1:]...)
			if len(pl.Orders) == 0 {
				ob.tree.Remove(price)
			}
			return true
		}
	}
	return false
}

func (ob *OrderBook) Get(id uint64) *Order {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	return ob.orderIndex[id]
}

func (ob *OrderBook) Len() int {
	ob.mu.RLock()
	defer ob.mu.RUnlock()
	return ob.tree.Size()
}

// 按价格升序遍历所有档位，callback 返回 false 时提前终止
// 注：callback 内可能修改订单状态，故使用写锁
func (ob *OrderBook) Range(callback func(price decimal.Decimal, orders []*Order) bool) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	it := ob.tree.Iterator()
	for it.Next() {
		pl := it.Value().(*PriceLevel)
		if !callback(pl.Price, pl.Orders) {
			return
		}
	}
}

func (ob *OrderBook) Best() *Order {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.tree.Empty() {
		return nil
	}

	node := ob.tree.Left()
	if node == nil {
		return nil
	}

	pl := node.Value.(*PriceLevel)
	if len(pl.Orders) == 0 {
		return nil
	}
	return pl.Orders[0]
}

// 最高买入价（红黑树最右节点）
func (ob *OrderBook) BestBid() *Order {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.tree.Size() == 0 {
		return nil
	}

	node := ob.tree.Right()
	if node == nil {
		return nil
	}

	pl := node.Value.(*PriceLevel)
	if len(pl.Orders) == 0 {
		return nil
	}
	return pl.Orders[0]
}

// 最低卖出价（红黑树最左节点）
func (ob *OrderBook) BestAsk() *Order {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.tree.Size() == 0 {
		return nil
	}

	node := ob.tree.Left()
	if node == nil {
		return nil
	}

	pl := node.Value.(*PriceLevel)
	if len(pl.Orders) == 0 {
		return nil
	}
	return pl.Orders[0]
}

func (ob *OrderBook) BestPrice() (decimal.Decimal, bool) {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.tree.Empty() {
		return decimal.Zero, false
	}

	node := ob.tree.Left()
	if node == nil {
		return decimal.Zero, false
	}

	pl := node.Value.(*PriceLevel)
	return pl.Price, true
}

func (ob *OrderBook) Clear() {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.tree.Clear()
	ob.orderIndex = make(map[uint64]*Order)
}
