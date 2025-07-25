// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"kratos-dtm-examples/app/shop/service/internal/data"
	"kratos-dtm-examples/app/shop/service/internal/server"
	"kratos-dtm-examples/app/shop/service/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(logger log.Logger, registrar registry.Registrar, bootstrap *v1.Bootstrap) (*kratos.App, func(), error) {
	discovery := data.NewDiscovery(bootstrap)
	db := data.NewGormClient(bootstrap, logger)
	dataData, cleanup, err := data.NewData(logger, discovery, db)
	if err != nil {
		return nil, nil, err
	}
	productRepo := data.NewProductRepo(logger, dataData)
	stockDeductionLogRepo := data.NewStockDeductionLogRepo(logger, dataData)
	stockRepo := data.NewStockRepo(logger, dataData)
	stockService := service.NewStockService(logger, productRepo, stockDeductionLogRepo, stockRepo)
	orderRepo := data.NewOrderRepo(logger, dataData)
	orderService := service.NewOrderService(logger, orderRepo)
	paymentService := service.NewPaymentService(logger)
	grpcServer := server.NewGRPCServer(bootstrap, logger, stockService, orderService, paymentService)
	shopService := service.NewShopService(logger, stockService, orderService, paymentService)
	httpServer := server.NewRestServer(bootstrap, shopService)
	app := newApp(logger, registrar, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
