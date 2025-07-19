package models

import "gorm.io/gorm"

type Stock struct {
	gorm.Model

	ProductID uint  `gorm:"not null;uniqueIndex;comment:商品ID"` // 商品ID
	Quantity  int32 `gorm:"not null;comment:当前库存"`             // 当前库存数量
}

func (s Stock) TableName() string {
	return "stocks"
}
