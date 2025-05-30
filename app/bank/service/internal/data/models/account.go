package models

import (
	"gorm.io/gorm"

	bankV1 "kratos-dtm-examples/api/gen/go/bank/service/v1"
)

type Account struct {
	gorm.Model

	AccountId *string `gorm:"column:account_id; comment:账户ID"`
	UserId    *int64  `gorm:"column:user_id; comment:用户ID"`

	Balance  int64                `gorm:"column:balance; comment:账户余额，以分为单位"`
	Currency *bankV1.CurrencyType `gorm:"column:currency; comment:货币类型"`

	Status *bankV1.AccountStatus `gorm:"column:status; comment:账户状态"`
}

func (u Account) TableName() string {
	return "bank_accounts"
}
