package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Mobile     string `gorm:"index:idx_mobile;unique;varchar(11);not null"`
	Password   string `gorm:"type:varchar(64);not null"`
	Salt       string `gorm:"type:varchar(16)"`
	NickName   string `gorm:"type:varchar(32)"`
	Gender     string `gorm:"varchar(6);default:male"`
	IsEmployee bool   `gorm:"default:false"`
	Role       int    `gorm:"type:int;default:1;comment'1-普通用户,2-管理员'"`
}
