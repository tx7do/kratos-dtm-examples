package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"kratos-dtm-examples/app/shop/service/internal/data"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

type OrderService struct {
	shopV1.UnimplementedOrderServiceServer

	log *log.Helper

	repo *data.OrderRepo
}

func NewOrderService(logger log.Logger, repo *data.OrderRepo) *OrderService {
	l := log.NewHelper(log.With(logger, "module", "order/service/shop-service"))
	return &OrderService{
		log:  l,
		repo: repo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *shopV1.CreateOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}

func (s *OrderService) TryCreateOrder(ctx context.Context, req *shopV1.TryCreateOrderRequest) (*shopV1.TryCreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryCreateOrder not implemented")
}

func (s *OrderService) ConfirmCreateOrder(ctx context.Context, req *shopV1.ConfirmCreateOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmCreateOrder not implemented")
}

func (s *OrderService) CancelCreateOrder(ctx context.Context, req *shopV1.CancelCreateOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelCreateOrder not implemented")
}

func (s *OrderService) RefundOrder(ctx context.Context, req *shopV1.RefundOrderRequest) (*shopV1.RefundOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundOrder not implemented")
}
