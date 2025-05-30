package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        *string        `gorm:"column:name"`
	Description sql.NullString `gorm:"column:description"`
	Stock       int32          `gorm:"column:stock"`
}

func (u Product) TableName() string {
	return "products"
}
