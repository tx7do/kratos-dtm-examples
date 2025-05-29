package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"kratos-dtm-examples/app/bank/service/internal/service"

	bankV1 "kratos-dtm-examples/api/gen/go/bank/service/v1"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	cfg *conf.Bootstrap, logger log.Logger,
	bankService *service.BankService,
) *grpc.Server {
	if cfg == nil || cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil
	}

	srv := rpc.CreateGrpcServer(
		cfg,
		logging.Server(logger),
	)

	bankV1.RegisterBankServiceServer(srv, bankService)

	return srv
}
