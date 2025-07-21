package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"kratos-dtm-examples/app/shop/service/internal/data/models"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"

	"kratos-dtm-examples/pkg/dtmgorm"
)

type StockRepo struct {
	data *Data
	log  *log.Helper
}

func NewStockRepo(logger log.Logger, db *Data) *StockRepo {
	return &StockRepo{
		data: db,
		log:  log.NewHelper(log.With(logger, "module", "shop/service/data/stock")),
	}
}

func (r *StockRepo) DeductStock(productID uint32, quantity int32) error {
	return r.DeductStockWithTx(r.data.db, productID, quantity)
}

func (r *StockRepo) DeductStockWithTx(tx *gorm.DB, productID uint32, quantity int32) error {
	// 使用事务来确保数据一致性
	return tx.Transaction(func(tx *gorm.DB) error {
		// 查询当前库存
		var stock models.Stock
		if err := tx.Model(&models.Stock{}).
			Select("quantity", "locked").
			Where("product_id = ?", productID).
			First(&stock).Error; err != nil {
			return err
		}

		// 检查库存是否足够
		if stock.Quantity-stock.Locked < quantity {
			return gorm.ErrRecordNotFound
		}

		// 扣减库存
		result := tx.Model(&models.Stock{}).
			Where("product_id = ? AND quantity >= ?", productID, quantity).
			UpdateColumn("quantity", gorm.Expr("quantity - ?", quantity))

		if result.Error != nil {
			return result.Error
		}

		// 如果没有更新任何记录，说明库存不足
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func (r *StockRepo) TryDeductStock(ctx context.Context, req *shopV1.TryDeductStockRequest) (*shopV1.StockResponse, error) {
	var err error

	err = dtmgorm.BarrierGorm(ctx, r.data.db, func(tx *gorm.DB) error {
		// 查询当前库存记录
		var stock models.Stock
		if err := tx.Model(&models.Stock{}).
			Where("product_id = ?", req.GetProductId()).
			First(&stock).Error; err != nil {
			return err
		}

		// 检查库存是否足够
		if stock.Quantity-stock.Locked < req.GetQuantity() {
			return gorm.ErrRecordNotFound
		}

		// 锁定库存
		result := tx.Model(&models.Stock{}).
			Where("product_id = ? AND quantity - locked >= ?", req.GetProductId(), req.GetQuantity()).
			UpdateColumn("locked", gorm.Expr("locked + ?", req.GetQuantity()))
		if result.Error != nil {
			return result.Error
		}

		// 检查是否更新成功
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
	if err != nil {
		return &shopV1.StockResponse{Success: false, Message: err.Error()}, nil
	}

	r.log.Infof("Attempting to deduct stock for product_id: %d, quantity: %d", req.GetProductId(), req.GetQuantity())

	return &shopV1.StockResponse{
		Success: true,
		Message: "Stock deduction initiated successfully",
	}, nil
}

func (r *StockRepo) ConfirmDeductStock(ctx context.Context, req *shopV1.ConfirmDeductStockRequest) (*shopV1.StockResponse, error) {
	var err error

	err = dtmgorm.BarrierGorm(ctx, r.data.db, func(tx *gorm.DB) error {
		// 查询当前库存记录
		var stock models.Stock
		if err = tx.Model(&models.Stock{}).
			Where("product_id = ?", req.GetProductId()).
			First(&stock).Error; err != nil {
			return err
		}

		// 确认扣减库存，将锁定的库存正式扣减
		result := tx.Model(&models.Stock{}).
			Where("product_id = ? AND locked >= ?", req.GetProductId(), req.GetQuantity()).
			Updates(map[string]interface{}{
				"locked":   gorm.Expr("locked - ?", req.GetQuantity()),
				"quantity": gorm.Expr("quantity - ?", req.GetQuantity()),
			})
		if result.Error != nil {
			return result.Error
		}

		// 检查是否更新成功
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
	if err != nil {
		return &shopV1.StockResponse{Success: false, Message: err.Error()}, nil
	}

	r.log.Infof("Confirming stock deduction for product_id: %d, quantity: %d", req.GetProductId(), req.GetQuantity())

	return &shopV1.StockResponse{
		Success: true,
		Message: "Stock deduction confirmed successfully",
	}, nil
}

func (r *StockRepo) CancelDeductStock(ctx context.Context, req *shopV1.CancelDeductStockRequest) (*shopV1.StockResponse, error) {
	var err error

	err = dtmgorm.BarrierGorm(ctx, r.data.db, func(tx *gorm.DB) error {
		// 查询当前库存记录
		var stock models.Stock
		if err = tx.Model(&models.Stock{}).
			Where("product_id = ?", req.GetProductId()).
			First(&stock).Error; err != nil {
			return err
		}

		// 恢复库存
		result := tx.Model(&models.Stock{}).
			Where("product_id = ?", req.GetProductId()).
			UpdateColumn("quantity", gorm.Expr("quantity + ?", req.GetQuantity()))
		if result.Error != nil {
			return result.Error
		}

		// 检查是否更新成功
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
	if err != nil {
		return &shopV1.StockResponse{Success: false, Message: err.Error()}, nil
	}

	r.log.Infof("Cancelling stock deduction for product_id: %d, quantity: %d", req.GetProductId(), req.GetQuantity())

	return &shopV1.StockResponse{
		Success: true,
		Message: "Stock deduction canceled successfully",
	}, nil
}

func (r *StockRepo) RefundStock(ctx context.Context, req *shopV1.RefundStockRequest) (*shopV1.StockResponse, error) {
	var err error

	err = dtmgorm.BarrierGorm(ctx, r.data.db, func(tx *gorm.DB) error {
		// 查询当前库存记录
		var stock models.Stock
		if err = tx.Model(&models.Stock{}).
			Where("product_id = ?", req.GetProductId()).
			First(&stock).Error; err != nil {
			return err
		}

		// 增加库存
		result := tx.Model(&models.Stock{}).
			Where("product_id = ?", req.GetProductId()).
			UpdateColumn("quantity", gorm.Expr("quantity + ?", req.GetQuantity()))
		if result.Error != nil {
			return result.Error
		}

		// 检查是否更新成功
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
	if err != nil {
		return &shopV1.StockResponse{Success: false, Message: err.Error()}, nil
	}

	r.log.Infof("Refunding stock for product_id: %d, quantity: %d", req.GetProductId(), req.GetQuantity())

	return &shopV1.StockResponse{
		Success: true,
		Message: "Stock refunded successfully",
	}, nil
}
