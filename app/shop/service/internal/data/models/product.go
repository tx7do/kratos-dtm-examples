package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        *string        `gorm:"column:name; comment:商品名称"`
	Description sql.NullString `gorm:"column:description; comment:商品描述"`
	Price       int64          `gorm:"column:price; comment:商品价格(分)"`
}

func (u Product) TableName() string {
	return "products"
}
