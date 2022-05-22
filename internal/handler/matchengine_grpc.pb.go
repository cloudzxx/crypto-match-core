package handler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MatchEngineClient interface {
	PutLimitOrder(context.Context, *PutLimitRequest, ...grpc.CallOption) (*OrderReply, error)
	PutMarketOrder(context.Context, *PutMarketRequest, ...grpc.CallOption) (*OrderReply, error)
	CancelOrder(context.Context, *CancelRequest, ...grpc.CallOption) (*CancelReply, error)
	QueryOrder(context.Context, *QueryRequest, ...grpc.CallOption) (*OrderReply, error)
	QueryBalance(context.Context, *BalanceRequest, ...grpc.CallOption) (*BalanceReply, error)
	QueryOrderBook(context.Context, *OrderBookRequest, ...grpc.CallOption) (*OrderBookReply, error)
	QueryDepth(context.Context, *DepthRequest, ...grpc.CallOption) (*DepthReply, error)
	MarketList(context.Context, *Empty, ...grpc.CallOption) (*MarketListReply, error)
	MarketSummary(context.Context, *MarketRequest, ...grpc.CallOption) (*MarketSummaryReply, error)
}

type MatchEngineServer interface {
	PutLimitOrder(context.Context, *PutLimitRequest) (*OrderReply, error)
	PutMarketOrder(context.Context, *PutMarketRequest) (*OrderReply, error)
	CancelOrder(context.Context, *CancelRequest) (*CancelReply, error)
	QueryOrder(context.Context, *QueryRequest) (*OrderReply, error)
	QueryBalance(context.Context, *BalanceRequest) (*BalanceReply, error)
	QueryOrderBook(context.Context, *OrderBookRequest) (*OrderBookReply, error)
	QueryDepth(context.Context, *DepthRequest) (*DepthReply, error)
	MarketList(context.Context, *Empty) (*MarketListReply, error)
	MarketSummary(context.Context, *MarketRequest) (*MarketSummaryReply, error)
}

func RegisterMatchEngineServer(s grpc.ServiceRegistrar, srv MatchEngineServer) {
	s.RegisterService(&_MatchEngine_serviceDesc, srv)
}

var _MatchEngine_serviceDesc = grpc.ServiceDesc{
	ServiceName: "matchengine.MatchEngine",
	HandlerType: (*MatchEngineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PutLimitOrder",
			Handler:    _MatchEngine_PutLimitOrder_Handler,
		},
		{
			MethodName: "PutMarketOrder",
			Handler:    _MatchEngine_PutMarketOrder_Handler,
		},
		{
			MethodName: "CancelOrder",
			Handler:    _MatchEngine_CancelOrder_Handler,
		},
		{
			MethodName: "QueryOrder",
			Handler:    _MatchEngine_QueryOrder_Handler,
		},
		{
			MethodName: "QueryBalance",
			Handler:    _MatchEngine_QueryBalance_Handler,
		},
		{
			MethodName: "QueryOrderBook",
			Handler:    _MatchEngine_QueryOrderBook_Handler,
		},
		{
			MethodName: "QueryDepth",
			Handler:    _MatchEngine_QueryDepth_Handler,
		},
		{
			MethodName: "MarketList",
			Handler:    _MatchEngine_MarketList_Handler,
		},
		{
			MethodName: "MarketSummary",
			Handler:    _MatchEngine_MarketSummary_Handler,
		},
	},
	Metadata: "matchengine.proto",
}

func _MatchEngine_PutLimitOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).PutLimitOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/PutLimitOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).PutLimitOrder(ctx, req.(*PutLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchEngine_PutMarketOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutMarketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).PutMarketOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/PutMarketOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).PutMarketOrder(ctx, req.(*PutMarketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchEngine_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/CancelOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).CancelOrder(ctx, req.(*CancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchEngine_QueryOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).QueryOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/QueryOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).QueryOrder(ctx, req.(*QueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchEngine_QueryBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).QueryBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/QueryBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).QueryBalance(ctx, req.(*BalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchEngine_QueryOrderBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).QueryOrderBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/QueryOrderBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).QueryOrderBook(ctx, req.(*OrderBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchEngine_QueryDepth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).QueryDepth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/QueryDepth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).QueryDepth(ctx, req.(*DepthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchEngine_MarketList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).MarketList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/MarketList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).MarketList(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchEngine_MarketSummary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchEngineServer).MarketSummary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/matchengine.MatchEngine/MarketSummary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchEngineServer).MarketSummary(ctx, req.(*MarketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

type UnimplementedMatchEngineServer struct{}

func (UnimplementedMatchEngineServer) PutLimitOrder(context.Context, *PutLimitRequest) (*OrderReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutLimitOrder not implemented")
}
func (UnimplementedMatchEngineServer) PutMarketOrder(context.Context, *PutMarketRequest) (*OrderReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutMarketOrder not implemented")
}
func (UnimplementedMatchEngineServer) CancelOrder(context.Context, *CancelRequest) (*CancelReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
func (UnimplementedMatchEngineServer) QueryOrder(context.Context, *QueryRequest) (*OrderReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryOrder not implemented")
}
func (UnimplementedMatchEngineServer) QueryBalance(context.Context, *BalanceRequest) (*BalanceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryBalance not implemented")
}
func (UnimplementedMatchEngineServer) QueryOrderBook(context.Context, *OrderBookRequest) (*OrderBookReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryOrderBook not implemented")
}
func (UnimplementedMatchEngineServer) QueryDepth(context.Context, *DepthRequest) (*DepthReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryDepth not implemented")
}
func (UnimplementedMatchEngineServer) MarketList(context.Context, *Empty) (*MarketListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarketList not implemented")
}
func (UnimplementedMatchEngineServer) MarketSummary(context.Context, *MarketRequest) (*MarketSummaryReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarketSummary not implemented")
}