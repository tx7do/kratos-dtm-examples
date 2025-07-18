package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID    uint `gorm:"not null;index"`
	ProductID uint `gorm:"not null;index"`

	Quantity   int     `gorm:"not null"`                  // 商品数量
	TotalPrice float64 `gorm:"not null"`                  // 总价
	Status     string  `gorm:"type:varchar(50);not null"` // 订单状态
}

func (u Order) TableName() string {
	return "orders"
}
