package model

import (
	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Code         string `gorm:"type:varchar(64);not null"`
	Name         string `gorm:"type:varchar(32);not null"`
	CouponType   int32
	DiscountType int32 `gorm:"default:1"`
	Discount     float32
	DiscountFrom float32 //minimal amount to use coupon
	Added        bool    `gorm:"comment' can be added with different coupon'"`
	Ratio        float32
	Used         bool
	EnableAt     string `gorm:"type:date;not null"`
	ExpiredAt    string `gorm:"type:date;not null"`
	AccountId    uint
	Account      *Account
}
