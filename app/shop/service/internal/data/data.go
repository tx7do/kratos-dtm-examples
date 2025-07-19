package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"google.golang.org/grpc/resolver"

	_ "github.com/dtm-labs/driver-kratos"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
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
	return bootstrap.NewDiscovery(cfg.Registry)
}
