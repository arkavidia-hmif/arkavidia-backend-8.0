package models

import (
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Admin struct {
	gorm.Model
	Username       string                `json:"username" gorm:"not null;unique"`
	HashedPassword types.EncryptedString `json:"password" gorm:"not null" visibility:"false"`
	ApprovesPhoto  []Photo               `json:"photos"`
	ApprovesTeam   []Team                `json:"teams"`
}
