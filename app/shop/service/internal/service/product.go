package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

type ProductService struct {
	shopV1.UnimplementedProductServiceServer

	log *log.Helper
}

func NewProductService(logger log.Logger) *ProductService {
	l := log.NewHelper(log.With(logger, "module", "product/service/shop-service"))
	return &ProductService{
		log: l,
	}
}

func (s *ProductService) DecreaseStock(ctx context.Context, req *shopV1.DecreaseStockRequest) (*shopV1.DecreaseStockResponse, error) {
	s.log.Infof("decrease stock for product %d, quantity %d", req.ProductId, req.Quantity)
	return &shopV1.DecreaseStockResponse{}, nil
}
