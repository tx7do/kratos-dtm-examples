package service

import (
	"context"

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

func (s *StockService) DeductStock(_ context.Context, req *shopV1.DeductStockRequest) (*shopV1.StockResponse, error) {
	exist, err := s.stockDeductionLogRepo.ExistLogByRequestID(req.GetRequestId())
	if err != nil {
		s.log.Errorf("failed to check stock deduction log existence for request_id: %s, error: %v", req.GetRequestId(), err)
		return nil, shopV1.ErrorInternalServerError("failed to check stock deduction log existence: %v", err)
	}
	if exist {
		s.log.Infof("stock deduction log already exists for request_id: %s", req.GetRequestId())
		return &shopV1.StockResponse{Success: true}, nil
	}

	if err = s.stockRepo.DeductStock(req.GetProductId(), req.GetQuantity()); err != nil {
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

	return &shopV1.StockResponse{
		Success: true,
		Message: "Stock deducted successfully",
	}, nil
}

func (s *StockService) DeductStockXA(ctx context.Context, req *shopV1.DeductStockRequest) (*shopV1.StockResponse, error) {
	//// 从上下文获取XA事务
	//xa, err := dtmgrpc.XaGrpcFromRequest(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//xa.XaLocalTransaction()

	return nil, nil
}

func (s *StockService) TryDeductStock(ctx context.Context, req *shopV1.TryDeductStockRequest) (*shopV1.StockResponse, error) {
	s.log.Infof("尝试扣除库存: %+v", req.RequestId)

	return s.stockRepo.TryDeductStock(ctx, req)
}

func (s *StockService) ConfirmDeductStock(ctx context.Context, req *shopV1.ConfirmDeductStockRequest) (*shopV1.StockResponse, error) {
	s.log.Infof("确认扣除库存: %+v", req.RequestId)

	return s.stockRepo.ConfirmDeductStock(ctx, req)
}

func (s *StockService) CancelDeductStock(ctx context.Context, req *shopV1.CancelDeductStockRequest) (*shopV1.StockResponse, error) {
	s.log.Infof("取消扣除库存: %+v", req.RequestId)

	return s.stockRepo.CancelDeductStock(ctx, req)
}

func (s *StockService) RefundStock(ctx context.Context, req *shopV1.RefundStockRequest) (*shopV1.StockResponse, error) {
	s.log.Infof("RefundStock called with request: %+v", req)

	return s.stockRepo.RefundStock(ctx, req)
}
