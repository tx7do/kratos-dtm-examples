package server

import (
	"github.com/dtm-labs/client/workflow"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"kratos-dtm-examples/app/shop/service/internal/service"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"

	serviceName "kratos-dtm-examples/pkg/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	cfg *conf.Bootstrap, logger log.Logger,
	stockService *service.StockService,
	orderService *service.OrderService,
	paymentService *service.PaymentService,
) *grpc.Server {
	if cfg == nil || cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil
	}

	srv := rpc.CreateGrpcServer(
		cfg,
		logging.Server(logger),
	)

	shopV1.RegisterStockServiceServer(srv, stockService)
	shopV1.RegisterOrderServiceServer(srv, orderService)
	shopV1.RegisterPaymentServiceServer(srv, paymentService)

	// 注册操作需要在业务服务启动之后执行，因为当进程crash，dtm会回调业务服务器，继续未完成的任务
	workflow.InitGrpc(serviceName.DtmServerAddress, serviceName.ShopServerAddress, srv.Server)

	return srv
}
