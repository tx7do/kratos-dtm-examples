package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

// Data .
type Data struct {
	log *log.Helper
}

// NewData .
func NewData(logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data/bank-service"))

	d := &Data{
		log: l,
	}

	return d, func() {
		l.Info("message", "closing the data resources")
	}, nil
}

// NewDiscovery 创建服務发现客户端
func NewDiscovery(cfg *conf.Bootstrap) registry.Discovery {
	return bRegistry.NewDiscovery(cfg.Registry)
}
