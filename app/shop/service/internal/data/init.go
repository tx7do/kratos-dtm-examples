package data

import "github.com/google/wire"

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDiscovery,
	NewGormClient,

	NewOrderRepo,
	NewProductRepo,
	NewUserRepo,
	NewStockDeductionLogRepo,
	NewStockRepo,
)
