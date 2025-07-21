package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/grpc/resolver"
	"gorm.io/gorm"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	_ "kratos-dtm-examples/pkg/dtmdriver-kratos"
)

// Data .
type Data struct {
	log *log.Helper

	db *gorm.DB
}

// NewData .
func NewData(logger log.Logger, rr registry.Discovery, db *gorm.DB) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data/shop-service"))

	d := &Data{
		log: l,
		db:  db,
	}

	//l.Info("message", "initializing data resources")

	// 注册Kratos的gRPC解析器的用于动态解析服务地址，用于与Dtm服务通信
	resolver.Register(discovery.NewBuilder(rr, discovery.WithInsecure(true)))

	return d, func() {
		l.Info("message", "closing the data resources")
	}, nil
}

func NewDiscovery(cfg *conf.Bootstrap) registry.Discovery {
	return bootstrap.NewDiscovery(cfg.Registry)
}
