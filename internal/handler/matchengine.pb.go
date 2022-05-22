package handler

import (
	"fmt"
)

type Empty struct{}

func (x *Empty) Reset()         { *x = Empty{} }
func (x *Empty) String() string { return fmt.Sprintf("Empty{}") }
func (Empty) ProtoMessage()    {}

type PutLimitRequest struct {
	UserId uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3"`
	Market string `protobuf:"bytes,2,opt,name=market,proto3"`
	Side   string `protobuf:"bytes,3,opt,name=side,proto3"`
	Price  string `protobuf:"bytes,4,opt,name=price,proto3"`
	Amount string `protobuf:"bytes,5,opt,name=amount,proto3"`
	Source string `protobuf:"bytes,6,opt,name=source,proto3"`
}

func (m *PutLimitRequest) Reset()         { *m = PutLimitRequest{} }
func (m *PutLimitRequest) String() string  { return fmt.Sprintf("PutLimitRequest%+v", m) }
func (PutLimitRequest) ProtoMessage()      {}

func (m *PutLimitRequest) GetUserId() uint32 { return m.UserId }
func (m *PutLimitRequest) GetMarket() string  { return m.Market }
func (m *PutLimitRequest) GetSide() string    { return m.Side }
func (m *PutLimitRequest) GetPrice() string   { return m.Price }
func (m *PutLimitRequest) GetAmount() string  { return m.Amount }
func (m *PutLimitRequest) GetSource() string  { return m.Source }

type PutMarketRequest struct {
	UserId uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3"`
	Market string `protobuf:"bytes,2,opt,name=market,proto3"`
	Side   string `protobuf:"bytes,3,opt,name=side,proto3"`
	Amount string `protobuf:"bytes,4,opt,name=amount,proto3"`
	Source string `protobuf:"bytes,5,opt,name=source,proto3"`
}

func (m *PutMarketRequest) Reset()         { *m = PutMarketRequest{} }
func (m *PutMarketRequest) String() string  { return fmt.Sprintf("PutMarketRequest%+v", m) }
func (PutMarketRequest) ProtoMessage()      {}

func (m *PutMarketRequest) GetUserId() uint32 { return m.UserId }
func (m *PutMarketRequest) GetMarket() string  { return m.Market }
func (m *PutMarketRequest) GetSide() string    { return m.Side }
func (m *PutMarketRequest) GetAmount() string  { return m.Amount }
func (m *PutMarketRequest) GetSource() string  { return m.Source }

type CancelRequest struct {
	UserId  uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3"`
	OrderId uint64 `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3"`
}

func (m *CancelRequest) Reset()         { *m = CancelRequest{} }
func (m *CancelRequest) String() string  { return fmt.Sprintf("CancelRequest%+v", m) }
func (CancelRequest) ProtoMessage()      {}

func (m *CancelRequest) GetUserId() uint32  { return m.UserId }
func (m *CancelRequest) GetOrderId() uint64  { return m.OrderId }

type QueryRequest struct {
	UserId  uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3"`
	OrderId uint64 `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3"`
}

func (m *QueryRequest) Reset()         { *m = QueryRequest{} }
func (m *QueryRequest) String() string  { return fmt.Sprintf("QueryRequest%+v", m) }
func (QueryRequest) ProtoMessage()      {}

func (m *QueryRequest) GetUserId() uint32  { return m.UserId }
func (m *QueryRequest) GetOrderId() uint64  { return m.OrderId }

type BalanceRequest struct {
	UserId uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3"`
	Asset  string `protobuf:"bytes,2,opt,name=asset,proto3"`
}

func (m *BalanceRequest) Reset()         { *m = BalanceRequest{} }
func (m *BalanceRequest) String() string  { return fmt.Sprintf("BalanceRequest%+v", m) }
func (BalanceRequest) ProtoMessage()      {}

func (m *BalanceRequest) GetUserId() uint32 { return m.UserId }
func (m *BalanceRequest) GetAsset() string   { return m.Asset }

type OrderBookRequest struct {
	Market string `protobuf:"bytes,1,opt,name=market,proto3"`
	Limit  int32  `protobuf:"varint,2,opt,name=limit,proto3"`
}

func (m *OrderBookRequest) Reset()         { *m = OrderBookRequest{} }
func (m *OrderBookRequest) String() string  { return fmt.Sprintf("OrderBookRequest%+v", m) }
func (OrderBookRequest) ProtoMessage()      {}

func (m *OrderBookRequest) GetMarket() string { return m.Market }
func (m *OrderBookRequest) GetLimit() int32   { return m.Limit }

type DepthRequest struct {
	Market string `protobuf:"bytes,1,opt,name=market,proto3"`
	Limit  int32  `protobuf:"varint,2,opt,name=limit,proto3"`
}

func (m *DepthRequest) Reset()         { *m = DepthRequest{} }
func (m *DepthRequest) String() string { return fmt.Sprintf("DepthRequest%+v", m) }
func (DepthRequest) ProtoMessage()     {}

func (m *DepthRequest) GetMarket() string { return m.Market }
func (m *DepthRequest) GetLimit() int32   { return m.Limit }

type MarketRequest struct {
	Market string `protobuf:"bytes,1,opt,name=market,proto3"`
}

func (m *MarketRequest) Reset()         { *m = MarketRequest{} }
func (m *MarketRequest) String() string  { return fmt.Sprintf("MarketRequest%+v", m) }
func (MarketRequest) ProtoMessage()      {}

func (m *MarketRequest) GetMarket() string { return m.Market }

type OrderReply struct {
	OrderId   uint64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3"`
	UserId    uint32 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3"`
	Market    string `protobuf:"bytes,3,opt,name=market,proto3"`
	Side      string `protobuf:"bytes,4,opt,name=side,proto3"`
	Price     string `protobuf:"bytes,5,opt,name=price,proto3"`
	Amount    string `protobuf:"bytes,6,opt,name=amount,proto3"`
	Left      string `protobuf:"bytes,7,opt,name=left,proto3"`
	DealStock string `protobuf:"bytes,8,opt,name=deal_stock,json=dealStock,proto3"`
	DealMoney string `protobuf:"bytes,9,opt,name=deal_money,json=dealMoney,proto3"`
	DealFee   string `protobuf:"bytes,10,opt,name=deal_fee,json=dealFee,proto3"`
	Status    string `protobuf:"bytes,11,opt,name=status,proto3"`
	CreatedAt int64  `protobuf:"varint,12,opt,name=created_at,json=createdAt,proto3"`
}

func (m *OrderReply) Reset()         { *m = OrderReply{} }
func (m *OrderReply) String() string  { return fmt.Sprintf("OrderReply%+v", m) }
func (OrderReply) ProtoMessage()      {}

func (m *OrderReply) GetOrderId() uint64   { return m.OrderId }
func (m *OrderReply) GetUserId() uint32    { return m.UserId }
func (m *OrderReply) GetMarket() string    { return m.Market }
func (m *OrderReply) GetSide() string      { return m.Side }
func (m *OrderReply) GetPrice() string     { return m.Price }
func (m *OrderReply) GetAmount() string    { return m.Amount }
func (m *OrderReply) GetLeft() string      { return m.Left }
func (m *OrderReply) GetDealStock() string { return m.DealStock }
func (m *OrderReply) GetDealMoney() string { return m.DealMoney }
func (m *OrderReply) GetDealFee() string   { return m.DealFee }
func (m *OrderReply) GetStatus() string    { return m.Status }
func (m *OrderReply) GetCreatedAt() int64  { return m.CreatedAt }

type CancelReply struct {
	Success bool   `protobuf:"varint,1,opt,name=success,proto3"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3"`
}

func (m *CancelReply) Reset()         { *m = CancelReply{} }
func (m *CancelReply) String() string { return fmt.Sprintf("CancelReply%+v", m) }
func (CancelReply) ProtoMessage()     {}

func (m *CancelReply) GetSuccess() bool   { return m.Success }
func (m *CancelReply) GetMessage() string { return m.Message }

type BalanceReply struct {
	UserId    uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3"`
	Asset     string `protobuf:"bytes,2,opt,name=asset,proto3"`
	Available string `protobuf:"bytes,3,opt,name=available,proto3"`
	Freeze    string `protobuf:"bytes,4,opt,name=freeze,proto3"`
}

func (m *BalanceReply) Reset()         { *m = BalanceReply{} }
func (m *BalanceReply) String() string  { return fmt.Sprintf("BalanceReply%+v", m) }
func (BalanceReply) ProtoMessage()      {}

func (m *BalanceReply) GetUserId() uint32    { return m.UserId }
func (m *BalanceReply) GetAsset() string      { return m.Asset }
func (m *BalanceReply) GetAvailable() string  { return m.Available }
func (m *BalanceReply) GetFreeze() string     { return m.Freeze }

type OrderBookEntry struct {
	Price  string `protobuf:"bytes,1,opt,name=price,proto3"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3"`
}

func (m *OrderBookEntry) Reset()         { *m = OrderBookEntry{} }
func (m *OrderBookEntry) String() string { return fmt.Sprintf("OrderBookEntry%+v", m) }
func (OrderBookEntry) ProtoMessage()    {}

func (m *OrderBookEntry) GetPrice() string  { return m.Price }
func (m *OrderBookEntry) GetAmount() string { return m.Amount }

type OrderBookReply struct {
	Market string           `protobuf:"bytes,1,opt,name=market,proto3"`
	Asks   []*OrderBookEntry `protobuf:"bytes,2,rep,name=asks,proto3"`
	Bids   []*OrderBookEntry `protobuf:"bytes,3,rep,name=bids,proto3"`
}

func (m *OrderBookReply) Reset()         { *m = OrderBookReply{} }
func (m *OrderBookReply) String() string { return fmt.Sprintf("OrderBookReply%+v", m) }
func (OrderBookReply) ProtoMessage()     {}

func (m *OrderBookReply) GetMarket() string          { return m.Market }
func (m *OrderBookReply) GetAsks() []*OrderBookEntry { return m.Asks }
func (m *OrderBookReply) GetBids() []*OrderBookEntry { return m.Bids }

type DepthEntry struct {
	Price  string `protobuf:"bytes,1,opt,name=price,proto3"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3"`
	Total  string `protobuf:"bytes,3,opt,name=total,proto3"`
}

func (m *DepthEntry) Reset()         { *m = DepthEntry{} }
func (m *DepthEntry) String() string { return fmt.Sprintf("DepthEntry%+v", m) }
func (DepthEntry) ProtoMessage()    {}

func (m *DepthEntry) GetPrice() string  { return m.Price }
func (m *DepthEntry) GetAmount() string { return m.Amount }
func (m *DepthEntry) GetTotal() string  { return m.Total }

type DepthReply struct {
	Market  string       `protobuf:"bytes,1,opt,name=market,proto3"`
	Asks    []*DepthEntry `protobuf:"bytes,2,rep,name=asks,proto3"`
	Bids    []*DepthEntry `protobuf:"bytes,3,rep,name=bids,proto3"`
}

func (m *DepthReply) Reset()         { *m = DepthReply{} }
func (m *DepthReply) String() string { return fmt.Sprintf("DepthReply%+v", m) }
func (DepthReply) ProtoMessage()     {}

func (m *DepthReply) GetMarket() string     { return m.Market }
func (m *DepthReply) GetAsks() []*DepthEntry  { return m.Asks }
func (m *DepthReply) GetBids() []*DepthEntry  { return m.Bids }

type MarketInfo struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3"`
	Stock string `protobuf:"bytes,2,opt,name=stock,proto3"`
	Money string `protobuf:"bytes,3,opt,name=money,proto3"`
}

func (m *MarketInfo) Reset()         { *m = MarketInfo{} }
func (m *MarketInfo) String() string  { return fmt.Sprintf("MarketInfo%+v", m) }
func (MarketInfo) ProtoMessage()     {}

func (m *MarketInfo) GetName() string  { return m.Name }
func (m *MarketInfo) GetStock() string  { return m.Stock }
func (m *MarketInfo) GetMoney() string  { return m.Money }

type MarketListReply struct {
	Markets []*MarketInfo `protobuf:"bytes,1,rep,name=markets,proto3"`
}

func (m *MarketListReply) Reset()         { *m = MarketListReply{} }
func (m *MarketListReply) String() string { return fmt.Sprintf("MarketListReply%+v", m) }
func (MarketListReply) ProtoMessage()     {}

func (m *MarketListReply) GetMarkets() []*MarketInfo { return m.Markets }

type MarketSummaryReply struct {
	Name   string `protobuf:"bytes,1,opt,name=name,proto3"`
	Last   string `protobuf:"bytes,2,opt,name=last,proto3"`
	High   string `protobuf:"bytes,3,opt,name=high,proto3"`
	Low    string `protobuf:"bytes,4,opt,name=low,proto3"`
	Volume string `protobuf:"bytes,5,opt,name=volume,proto3"`
	Deal   string `protobuf:"bytes,6,opt,name=deal,proto3"`
}

func (m *MarketSummaryReply) Reset()         { *m = MarketSummaryReply{} }
func (m *MarketSummaryReply) String() string  { return fmt.Sprintf("MarketSummaryReply%+v", m) }
func (MarketSummaryReply) ProtoMessage()      {}

func (m *MarketSummaryReply) GetName() string   { return m.Name }
func (m *MarketSummaryReply) GetLast() string   { return m.Last }
func (m *MarketSummaryReply) GetHigh() string   { return m.High }
func (m *MarketSummaryReply) GetLow() string    { return m.Low }
func (m *MarketSummaryReply) GetVolume() string { return m.Volume }
func (m *MarketSummaryReply) GetDeal() string    { return m.Deal }

