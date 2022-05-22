package handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/cloudzxx/crypto-match-core/internal/publisher"
	"github.com/cloudzxx/crypto-match-core/internal/repository"
	"github.com/cloudzxx/crypto-match-core/internal/service"
	"google.golang.org/grpc"
	"github.com/shopspring/decimal"
)

// gRPC 服务端：负责协议层解析 + 调用撮合核心 + 写操作日志 + 推送事件
type Server struct {
	UnimplementedMatchEngineServer
	addr        string
	grpcServer  *grpc.Server
	mm          *service.Manager
	bm          *service.BalanceManager
	matcher     *service.Matcher
	operlog     *repository.Writer
	publisher   *publisher.Publisher
	nextOrderID uint64
	stopped     int32
}

func NewServer(addr string, mm *service.Manager, bm *service.BalanceManager, operlogWriter *repository.Writer, pub *publisher.Publisher) *Server {
	return &Server{
		addr:      addr,
		mm:        mm,
		bm:        bm,
		matcher:   service.NewMatcher(mm, bm),
		operlog:   operlogWriter,
		publisher: pub,
		nextOrderID: 1,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	RegisterMatchEngineServer(s.grpcServer, s)

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Printf("grpc serve error: %v", err)
		}
	}()

	log.Printf("grpc server started at %s", s.addr)
	return nil
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// 限价单入口：分配 ID → 冻结资产 → 入订单簿 → 写操作日志
func (s *Server) PutLimitOrder(ctx context.Context, req *PutLimitRequest) (*OrderReply, error) {
	orderID := atomic.AddUint64(&s.nextOrderID, 1)

	mkt, err := s.mm.Get(req.Market)
	if err != nil {
		return nil, err
	}

	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return nil, err
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, err
	}

	order := &service.Order{
		ID:        orderID,
		Type:      service.OrderTypeLimit,
		Side:      service.SideAsk,
		UserID:    req.UserId,
		Market:    req.Market,
		Price:     price,
		Amount:    amount,
		Left:      amount,
		Source:    req.Source,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if req.Side == "bid" {
		// 买单冻结 quote（Money）市值
		order.Side = service.SideBid
		freezeAmount := price.Mul(amount)
		if err := s.bm.Freeze(req.UserId, mkt.Money, freezeAmount); err != nil {
			return nil, err
		}
		order.Freeze = freezeAmount
	} else {
		// 卖单冻结 base（Stock）数量
		if err := s.bm.Freeze(req.UserId, mkt.Stock, amount); err != nil {
			return nil, err
		}
		order.Freeze = amount
	}

	// 入订单簿等待撮合
	if order.Side == service.SideAsk {
		mkt.Asks.Insert(order)
	} else {
		mkt.Bids.Insert(order)
	}

	s.operlog.Write(&repository.OperLog{
		Type:  repository.OperTypeOrderPut,
		Data:  order,
		SeqID: s.operlog.NextSeqID(),
	})

	return &OrderReply{
		OrderId: order.ID,
		UserId:  order.UserID,
		Market:  order.Market,
		Side:    req.Side,
		Price:   req.Price,
		Amount:  req.Amount,
		Left:    req.Amount,
		Status:  "PUT",
	}, nil
}

func (s *Server) PutMarketOrder(ctx context.Context, req *PutMarketRequest) (*OrderReply, error) {
	return nil, fmt.Errorf("market order not implemented")
}

// 撤单流程：遍历所有市场检索订单 → 验证用户所有权 → 从订单簿移除 → 解冻资产
func (s *Server) CancelOrder(ctx context.Context, req *CancelRequest) (*CancelReply, error) {
	mktName := ""
	var order *service.Order

	markets := s.mm.List()
	for _, m := range markets {
		if o := m.Bids.Get(req.OrderId); o != nil {
			order = o
			mktName = m.Name
			break
		}
		if o := m.Asks.Get(req.OrderId); o != nil {
			order = o
			mktName = m.Name
			break
		}
	}

	if order == nil {
		return &CancelReply{Success: false, Message: "order not found"}, nil
	}

	if order.UserID != req.UserId {
		return &CancelReply{Success: false, Message: "unauthorized"}, nil
	}

	mkt, _ := s.mm.Get(mktName)
	if order.Side == service.SideBid {
		mkt.Bids.Remove(order)
		if err := s.bm.Unfreeze(order.UserID, mkt.Money, order.Freeze); err != nil {
			return &CancelReply{Success: false, Message: err.Error()}, nil
		}
	} else {
		mkt.Asks.Remove(order)
		if err := s.bm.Unfreeze(order.UserID, mkt.Stock, order.Freeze); err != nil {
			return &CancelReply{Success: false, Message: err.Error()}, nil
		}
	}

	return &CancelReply{Success: true, Message: ""}, nil
}

func (s *Server) QueryOrder(ctx context.Context, req *QueryRequest) (*OrderReply, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Server) QueryBalance(ctx context.Context, req *BalanceRequest) (*BalanceReply, error) {
	balance := s.bm.GetBalance(req.UserId, req.Asset)
	return &BalanceReply{
		UserId:    req.UserId,
		Asset:     req.Asset,
		Available: balance.Available.String(),
		Freeze:    balance.Freeze.String(),
	}, nil
}

func (s *Server) QueryOrderBook(ctx context.Context, req *OrderBookRequest) (*OrderBookReply, error) {
	mkt, err := s.mm.Get(req.Market)
	if err != nil {
		return nil, err
	}

	reply := &OrderBookReply{
		Market: req.Market,
		Asks:   make([]*OrderBookEntry, 0),
		Bids:   make([]*OrderBookEntry, 0),
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}

	askCount := 0
	mkt.Asks.Range(func(price decimal.Decimal, orders []*service.Order) bool {
		var total decimal.Decimal
		for _, o := range orders {
			total = total.Add(o.Left)
		}
		reply.Asks = append(reply.Asks, &OrderBookEntry{
			Price:  price.String(),
			Amount: total.String(),
		})
		askCount++
		return askCount < limit
	})

	bidCount := 0
	mkt.Bids.Range(func(price decimal.Decimal, orders []*service.Order) bool {
		var total decimal.Decimal
		for _, o := range orders {
			total = total.Add(o.Left)
		}
		reply.Bids = append(reply.Bids, &OrderBookEntry{
			Price:  price.String(),
			Amount: total.String(),
		})
		bidCount++
		return bidCount < limit
	})

	return reply, nil
}

func (s *Server) QueryDepth(ctx context.Context, req *DepthRequest) (*DepthReply, error) {
	return &DepthReply{}, nil
}

func (s *Server) MarketList(ctx context.Context, req *Empty) (*MarketListReply, error) {
	markets := s.mm.List()
	reply := &MarketListReply{
		Markets: make([]*MarketInfo, 0, len(markets)),
	}

	for _, mkt := range markets {
		reply.Markets = append(reply.Markets, &MarketInfo{
			Name:  mkt.Name,
			Stock: mkt.Stock,
			Money: mkt.Money,
		})
	}

	return reply, nil
}

func (s *Server) MarketSummary(ctx context.Context, req *MarketRequest) (*MarketSummaryReply, error) {
	return &MarketSummaryReply{}, nil
}