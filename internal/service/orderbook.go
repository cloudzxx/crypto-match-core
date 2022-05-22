package service

import (
	"sync"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID        uint64
	Type      OrderType
	Side      Side
	UserID    uint32
	Market    string
	Price     decimal.Decimal
	Amount    decimal.Decimal
	Left      decimal.Decimal
	Freeze    decimal.Decimal
	Source    string
	CreatedAt int64
	UpdatedAt int64
}

type OrderType uint8

const (
	OrderTypeLimit OrderType = iota
	OrderTypeMarket
)

type Side uint8

const (
	SideAsk Side = iota
	SideBid
)

type PriceLevel struct {
	Price  decimal.Decimal
	Orders []*Order
}

type OrderBook struct {
	tree       *redblacktree.Tree
	orderIndex map[uint64]*Order
	mu         sync.RWMutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		tree: redblacktree.NewWith(func(a, b interface{}) int {
			return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
		}),
		orderIndex: make(map[uint64]*Order),
	}
}

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

func (ob *OrderBook) Range(callback func(price decimal.Decimal, orders []*Order) bool) {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

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
