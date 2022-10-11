package model

import "gorm.io/gorm"

type Coffee struct {
	gorm.Model
	Name        string  `gorm:"varchar(32);not null"`
	Price       float32 `gorm:"not null"`
	RealPrice   float32 `gorm:"not null"`
	SoldNum     int32   `gorm:"default:0"`
	Sku         string  `gorm:"varchar(32);not null"`
	Description string  `gorm:"varchar(512)"`
	Image       string  `gorm:"varchar(1024)"`
}
