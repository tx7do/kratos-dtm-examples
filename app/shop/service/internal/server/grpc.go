package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"kratos-dtm-examples/app/shop/service/internal/service"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	cfg *conf.Bootstrap, logger log.Logger,
	shopService *service.ShopService,
	productService *service.ProductService,
) *grpc.Server {
	if cfg == nil || cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil
	}

	srv := rpc.CreateGrpcServer(
		cfg,
		logging.Server(logger),
	)

	shopV1.RegisterShopServiceServer(srv, shopService)
	shopV1.RegisterProductServiceServer(srv, productService)

	return srv
}
