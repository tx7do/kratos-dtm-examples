package data

import (
	"gorm.io/gorm"

	"kratos-dtm-examples/app/shop/service/internal/data/models"
)

type StockRepo struct {
	db *gorm.DB
}

func NewStockRepo(db *gorm.DB) *StockRepo {
	return &StockRepo{
		db: db,
	}
}

func (r *StockRepo) DeductStock(productID uint, quantity int) error {
	// 使用事务来确保数据一致性
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查库存是否足够
		var stock int
		if err := tx.Model(&models.Product{}).Where("id = ?", productID).Select("stock").Scan(&stock).Error; err != nil {
			return err
		}
		if stock < quantity {
			return gorm.ErrRecordNotFound // 库存不足
		}

		// 扣减库存
		if err := tx.Model(&models.Product{}).Where("id = ?", productID).UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
			return err
		}

		return nil
	})
}
