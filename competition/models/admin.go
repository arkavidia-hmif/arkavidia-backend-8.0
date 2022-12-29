package models

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Username       string          `json:"username" gorm:"not null;unique"`
	HashedPassword EncryptedString `json:"password" gorm:"not null" visibility:"false"`
	ApprovesPhoto  []Photo         `json:"photos"`
	ApprovesTeam   []Team          `json:"teams"`
}
