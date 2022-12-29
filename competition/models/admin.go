package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Admin struct {
	gorm.Model
	Username       string                `gorm:"not null;unique"`
	HashedPassword types.EncryptedString `gorm:"not null"`
	ApprovesPhoto  []Photo
	ApprovesTeam   []Team
}

type DisplayAdmin struct {
	ID             uint                  `json:"id,omitempty"`
	CreatedAt      time.Time             `json:"created_at,omitempty"`
	UpdatedAt      time.Time             `json:"updated_at,omitempty"`
	Username       string                `json:"username,omitempty"`
	HashedPassword types.EncryptedString `json:"-"`
	ApprovesPhoto  []Photo               `json:"photos,omitempty"`
	ApprovesTeam   []Team                `json:"teams,omitempty"`
}

func (admin Admin) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DisplayAdmin{
		ID:            admin.ID,
		CreatedAt:     admin.CreatedAt,
		UpdatedAt:     admin.UpdatedAt,
		Username:      admin.Username,
		ApprovesPhoto: admin.ApprovesPhoto,
		ApprovesTeam:  admin.ApprovesTeam,
	})
}
