package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc"
	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"

	"kratos-dtm-examples/app/shop/service/internal/service"
)

// NewRestServer new an HTTP server.
func NewRestServer(cfg *conf.Bootstrap, shopService *service.ShopService) *http.Server {
	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil
	}

	srv := rpc.CreateRestServer(cfg)

	shopV1.RegisterShopServiceHTTPServer(srv, shopService)

	return srv
}
