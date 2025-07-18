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

	repo *data.ProductRepo
}

func NewStockService(logger log.Logger, repo *data.ProductRepo) *StockService {
	l := log.NewHelper(log.With(logger, "module", "stock/service/shop-service"))
	return &StockService{
		log:  l,
		repo: repo,
	}
}

func (s *StockService) DeductStock(ctx context.Context, req *shopV1.DeductStockRequest) (*shopV1.DeductStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeductStock not implemented")
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
