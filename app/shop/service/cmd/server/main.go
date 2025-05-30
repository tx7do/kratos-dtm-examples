package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"kratos-dtm-examples/pkg/service"
)

var version = "1.0.0"

// go build -ldflags "-X main.version=x.y.z"

func newApp(ll log.Logger, rr registry.Registrar, gs *grpc.Server, hs *http.Server) *kratos.App {
	return bootstrap.NewApp(ll, rr, gs, hs)
}

func main() {
	bootstrap.Bootstrap(initApp, trans.Ptr(service.ShopService), trans.Ptr(version))
}
