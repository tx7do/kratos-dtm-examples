package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"kratos-dtm-examples/app/shop/service/internal/data"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

type OrderService struct {
	shopV1.UnimplementedOrderServiceServer

	log *log.Helper

	repo *data.OrderRepo
}

func NewOrderService(
	logger log.Logger,
	repo *data.OrderRepo,
) *OrderService {
	l := log.NewHelper(log.With(logger, "module", "order/service/shop-service"))
	return &OrderService{
		log:  l,
		repo: repo,
	}
}

func (s *OrderService) CreateOrder(_ context.Context, req *shopV1.CreateOrderRequest) (*shopV1.OrderResponse, error) {
	s.log.Infof("Creating order for user %d with product %d, quantity %d", req.UserId, req.ProductId, req.Quantity)

	// 检查订单是否已存在
	exists, err := s.repo.OrderExistsByRequestID(req.RequestId)
	if err != nil {
		s.log.Errorf("Failed to check if order exists: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to check order existence")
	}
	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "order already exists")
	}

	// 创建订单
	err = s.repo.CreateOrder(&shopV1.Order{
		UserId:    req.UserId,
		ProductId: req.ProductId,
		RequestId: req.RequestId,
		OrderNo:   req.OrderNo,
		Quantity:  req.Quantity,
		Status:    shopV1.OrderStatus_PENDING, // 设置订单状态为 PENDING
	})
	if err != nil {
		s.log.Errorf("Failed to create order: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create order")
	}

	s.log.Infof("Order created successfully for user %d", req.UserId)
	return &shopV1.OrderResponse{
		Success: true,
		Message: "Order created successfully",
	}, nil
}

func (s *OrderService) CreateOrderXA(ctx context.Context, req *shopV1.CreateOrderRequest) (*shopV1.OrderResponse, error) {
	// 从上下文获取XA事务
	//xa, err := dtmgrpc.XaGrpcFromRequest(ctx)
	//if err != nil {
	//	s.log.Errorf("Failed to get XA transaction from context: %v", err)
	//	return nil, shopV1.ErrorInternalServerError("failed to get XA transaction from context")
	//}

	//dtmgrpc.XaLocalTransaction()
	//
	//xa.CallBranch()

	return nil, nil
}

func (s *OrderService) TryCreateOrder(ctx context.Context, req *shopV1.TryCreateOrderRequest) (*shopV1.OrderResponse, error) {
	s.log.Infof("尝试创建订单： %s", req.OrderNo)

	return s.repo.TryCreateOrder(ctx, req)
}

func (s *OrderService) ConfirmCreateOrder(ctx context.Context, req *shopV1.ConfirmCreateOrderRequest) (*shopV1.OrderResponse, error) {
	s.log.Infof("确认创建订单： %s", req.OrderNo)

	return s.repo.ConfirmCreateOrder(ctx, req)
}

func (s *OrderService) CancelCreateOrder(ctx context.Context, req *shopV1.CancelCreateOrderRequest) (*shopV1.OrderResponse, error) {
	s.log.Infof("取消创建订单： %s", req.OrderNo)

	return s.repo.CancelCreateOrder(ctx, req)
}

func (s *OrderService) RefundOrder(ctx context.Context, req *shopV1.RefundOrderRequest) (*shopV1.OrderResponse, error) {
	s.log.Infof("Processing refund for order ID %s", req.OrderNo)

	return s.repo.RefundOrder(ctx, req)
}
