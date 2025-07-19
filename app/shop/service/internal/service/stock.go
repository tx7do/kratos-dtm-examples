package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/log"

	"kratos-dtm-examples/app/shop/service/internal/data"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

type StockService struct {
	shopV1.UnimplementedStockServiceServer

	log *log.Helper

	productRepo           *data.ProductRepo
	stockDeductionLogRepo *data.StockDeductionLogRepo
	stockRepo             *data.StockRepo
}

func NewStockService(
	logger log.Logger,
	productRepo *data.ProductRepo,
	stockDeductionLogRepo *data.StockDeductionLogRepo,
	stockRepo *data.StockRepo,
) *StockService {
	l := log.NewHelper(log.With(logger, "module", "stock/service/shop-service"))
	return &StockService{
		log:                   l,
		productRepo:           productRepo,
		stockDeductionLogRepo: stockDeductionLogRepo,
		stockRepo:             stockRepo,
	}
}

func (s *StockService) DeductStock(ctx context.Context, req *shopV1.DeductStockRequest) (*shopV1.DeductStockResponse, error) {
	exist, err := s.stockDeductionLogRepo.ExistLogByRequestID(req.GetRequestId())
	if err != nil {
		s.log.Errorf("failed to check stock deduction log existence for request_id: %s, error: %v", req.GetRequestId(), err)
		return nil, shopV1.ErrorInternalServerError("failed to check stock deduction log existence: %v", err)
	}
	if exist {
		s.log.Infof("stock deduction log already exists for request_id: %s", req.GetRequestId())
		return &shopV1.DeductStockResponse{Success: true}, nil
	}

	if err = s.stockRepo.DeductStock(uint(req.GetProductId()), int(req.GetQuantity())); err != nil {
		s.log.Errorf("failed to deduct stock for product_id: %d, quantity: %d, error: %v", req.GetProductId(), req.GetQuantity(), err)
		return nil, shopV1.ErrorInternalServerError("failed to deduct stock: %v", err)
	}

	if err = s.stockDeductionLogRepo.CreateLog(&shopV1.StockDeductionLog{
		ProductId: req.GetProductId(),
		//UserId:,
		RequestId: req.GetRequestId(),
		Quantity:  req.GetQuantity(),
	}); err != nil {
		s.log.Errorf("failed to create stock deduction log for request_id: %s, error: %v", req.GetRequestId(), err)
		return nil, shopV1.ErrorInternalServerError("failed to create stock deduction log: %v", err)
	}

	return &shopV1.DeductStockResponse{}, nil
}

func (s *StockService) TryDeductStock(ctx context.Context, req *shopV1.TryDeductStockRequest) (*shopV1.TryDeductStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryDeductStock not implemented")
}

func (s *StockService) ConfirmDeductStock(ctx context.Context, req *shopV1.ConfirmDeductStockRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmDeductStock not implemented")
}

func (s *StockService) CancelDeductStock(ctx context.Context, req *shopV1.CancelDeductStockRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelDeductStock not implemented")
}

func (s *StockService) RefundStock(ctx context.Context, req *shopV1.RefundStockRequest) (*shopV1.RefundStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundStock not implemented")
}
