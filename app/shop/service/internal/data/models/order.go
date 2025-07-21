package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID    uint32 `gorm:"not null;index"`                   // 用户ID
	ProductID uint32 `gorm:"not null;index"`                   // 商品ID
	OrderNo   string `gorm:"type:varchar(50);not null;unique"` // 订单号
	RequestID string `gorm:"type:varchar(50);not null;unique"` // 请求ID

	Quantity   int32   `gorm:"not null"`                  // 商品数量
	TotalPrice float64 `gorm:"not null"`                  // 总价
	Status     string  `gorm:"type:varchar(50);not null"` // 订单状态
}

func (u Order) TableName() string {
	return "orders"
}
