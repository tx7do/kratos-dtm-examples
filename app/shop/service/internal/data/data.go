package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"

	"google.golang.org/grpc/resolver"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
	"github.com/tx7do/kratos-bootstrap/rpc"

	_ "github.com/dtm-labs/driver-kratos"

	bankV1 "kratos-dtm-examples/api/gen/go/bank/service/v1"

	"kratos-dtm-examples/pkg/service"
)

// Data .
type Data struct {
	log *log.Helper
}

// NewData .
func NewData(logger log.Logger, rr registry.Discovery) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data/shop-service"))

	d := &Data{
		log: l,
	}

	resolver.Register(discovery.NewBuilder(rr, discovery.WithInsecure(true)))

	return d, func() {
		l.Info("message", "closing the data resources")
	}, nil
}

func NewDiscovery(cfg *conf.Bootstrap) registry.Discovery {
	return bRegistry.NewDiscovery(cfg.Registry)
}

func NewBankServiceClient(r registry.Discovery, cfg *conf.Bootstrap) bankV1.BankServiceClient {
	return bankV1.NewBankServiceClient(rpc.CreateGrpcClient(context.Background(), r, service.BankService, cfg))
}
