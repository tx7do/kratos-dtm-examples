// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: shop/service/v1/order.proto

package servicev1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	OrderService_CreateOrder_FullMethodName        = "/shop.service.v1.OrderService/CreateOrder"
	OrderService_CreateOrderXA_FullMethodName      = "/shop.service.v1.OrderService/CreateOrderXA"
	OrderService_TryCreateOrder_FullMethodName     = "/shop.service.v1.OrderService/TryCreateOrder"
	OrderService_ConfirmCreateOrder_FullMethodName = "/shop.service.v1.OrderService/ConfirmCreateOrder"
	OrderService_CancelCreateOrder_FullMethodName  = "/shop.service.v1.OrderService/CancelCreateOrder"
	OrderService_RefundOrder_FullMethodName        = "/shop.service.v1.OrderService/RefundOrder"
)

// OrderServiceClient is the client API for OrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 订单服务
type OrderServiceClient interface {
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	CreateOrderXA(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	TryCreateOrder(ctx context.Context, in *TryCreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	ConfirmCreateOrder(ctx context.Context, in *ConfirmCreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	CancelCreateOrder(ctx context.Context, in *CancelCreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	RefundOrder(ctx context.Context, in *RefundOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, OrderService_CreateOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) CreateOrderXA(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, OrderService_CreateOrderXA_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) TryCreateOrder(ctx context.Context, in *TryCreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, OrderService_TryCreateOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) ConfirmCreateOrder(ctx context.Context, in *ConfirmCreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, OrderService_ConfirmCreateOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) CancelCreateOrder(ctx context.Context, in *CancelCreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, OrderService_CancelCreateOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) RefundOrder(ctx context.Context, in *RefundOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, OrderService_RefundOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServiceServer is the server API for OrderService service.
// All implementations must embed UnimplementedOrderServiceServer
// for forward compatibility.
//
// 订单服务
type OrderServiceServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*OrderResponse, error)
	CreateOrderXA(context.Context, *CreateOrderRequest) (*OrderResponse, error)
	TryCreateOrder(context.Context, *TryCreateOrderRequest) (*OrderResponse, error)
	ConfirmCreateOrder(context.Context, *ConfirmCreateOrderRequest) (*OrderResponse, error)
	CancelCreateOrder(context.Context, *CancelCreateOrderRequest) (*OrderResponse, error)
	RefundOrder(context.Context, *RefundOrderRequest) (*OrderResponse, error)
	mustEmbedUnimplementedOrderServiceServer()
}

// UnimplementedOrderServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOrderServiceServer struct{}

func (UnimplementedOrderServiceServer) CreateOrder(context.Context, *CreateOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedOrderServiceServer) CreateOrderXA(context.Context, *CreateOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrderXA not implemented")
}
func (UnimplementedOrderServiceServer) TryCreateOrder(context.Context, *TryCreateOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryCreateOrder not implemented")
}
func (UnimplementedOrderServiceServer) ConfirmCreateOrder(context.Context, *ConfirmCreateOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmCreateOrder not implemented")
}
func (UnimplementedOrderServiceServer) CancelCreateOrder(context.Context, *CancelCreateOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelCreateOrder not implemented")
}
func (UnimplementedOrderServiceServer) RefundOrder(context.Context, *RefundOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundOrder not implemented")
}
func (UnimplementedOrderServiceServer) mustEmbedUnimplementedOrderServiceServer() {}
func (UnimplementedOrderServiceServer) testEmbeddedByValue()                      {}

// UnsafeOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServiceServer will
// result in compilation errors.
type UnsafeOrderServiceServer interface {
	mustEmbedUnimplementedOrderServiceServer()
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	// If the following call pancis, it indicates UnimplementedOrderServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

func _OrderService_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_CreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_CreateOrderXA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CreateOrderXA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_CreateOrderXA_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CreateOrderXA(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_TryCreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TryCreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).TryCreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_TryCreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).TryCreateOrder(ctx, req.(*TryCreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_ConfirmCreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmCreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).ConfirmCreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_ConfirmCreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).ConfirmCreateOrder(ctx, req.(*ConfirmCreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_CancelCreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelCreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CancelCreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_CancelCreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CancelCreateOrder(ctx, req.(*CancelCreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_RefundOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefundOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).RefundOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_RefundOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).RefundOrder(ctx, req.(*RefundOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderService_ServiceDesc is the grpc.ServiceDesc for OrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shop.service.v1.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _OrderService_CreateOrder_Handler,
		},
		{
			MethodName: "CreateOrderXA",
			Handler:    _OrderService_CreateOrderXA_Handler,
		},
		{
			MethodName: "TryCreateOrder",
			Handler:    _OrderService_TryCreateOrder_Handler,
		},
		{
			MethodName: "ConfirmCreateOrder",
			Handler:    _OrderService_ConfirmCreateOrder_Handler,
		},
		{
			MethodName: "CancelCreateOrder",
			Handler:    _OrderService_CancelCreateOrder_Handler,
		},
		{
			MethodName: "RefundOrder",
			Handler:    _OrderService_RefundOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shop/service/v1/order.proto",
}
