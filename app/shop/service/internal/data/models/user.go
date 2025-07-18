package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `gorm:"column:username; comment:用户名; uniqueIndex"`
	Password string `gorm:"column:password; comment:密码"`
	Nickname string `gorm:"column:nickname; comment:昵称"`
	Avatar   string `gorm:"column:avatar; comment:头像"`
}

func (u User) TableName() string {
	return "users"
}
