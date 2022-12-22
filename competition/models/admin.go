package models

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Username       string  `json:"username" gorm:"not null;unique"`
	HashedPassword []byte  `json:"password" gorm:"not null"`
	Approves       []Photo `json:"photos"`
}
