package models

import (
	"gorm.io/gorm"
)

// StockDeductionLog 库存扣减日志
type StockDeductionLog struct {
	gorm.Model

	UserID    uint   `gorm:"not null;index"`                   // 用户ID
	ProductID uint   `gorm:"not null;index"`                   // 商品ID
	Quantity  int    `gorm:"not null"`                         // 扣减数量
	RequestID string `gorm:"type:varchar(50);not null;unique"` // 请求ID
}

func (s StockDeductionLog) TableName() string {
	return "stock_deduction_logs"
}
