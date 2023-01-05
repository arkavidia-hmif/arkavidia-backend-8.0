package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Photo struct {
	gorm.Model
	FileName      uuid.UUID         `gorm:"type:uuid;unique"`
	FileExtension string            `gorm:"not null"`
	ParticipantID uint              `gorm:"not null"`
	AdminID       uint              `gorm:"default:null"`
	Status        types.PhotoStatus `gorm:"not null"`
	Type          types.PhotoType   `gorm:"not null"`
	Participant   Participant       `gorm:"foreignKey:ParticipantID;references:ID"`
	ApprovedBy    Admin             `gorm:"foreignKey:AdminID;references:ID"`
}

type DisplayPhoto struct {
	ID            uint              `json:"id,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	UpdatedAt     time.Time         `json:"updated_at,omitempty"`
	FileName      uuid.UUID         `json:"file_name,omitempty"`
	FileExtension string            `json:"file_extension,omitempty" `
	ParticipantID uint              `json:"participant_id,omitempty"`
	AdminID       uint              `json:"admin_id,omitempty"`
	Type          types.PhotoType   `json:"type,omitempty"`
	Status        types.PhotoStatus `json:"status,omitempty"`
}

func (photo Photo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DisplayPhoto{
		ID:            photo.ID,
		CreatedAt:     photo.CreatedAt,
		UpdatedAt:     photo.UpdatedAt,
		FileName:      photo.FileName,
		FileExtension: photo.FileExtension,
		ParticipantID: photo.ParticipantID,
		AdminID:       photo.AdminID,
		Type:          photo.Type,
		Status:        photo.Status,
	})
}

// Menambahkan constraint untuk mengecek apakah terdapat photos yang telah diapprove namun admin tidak tercatat
// atau photos yang belum diapprove namun admin tercatat
func (photo *Photo) BeforeSave(tx *gorm.DB) error {
	if photo.Status != "" {
		if photo.AdminID == 0 {
			return fmt.Errorf("ERROR: ADMIN MUST BE RECORDED")
		}
		if photo.Status == types.WaitingForApproval {
			return fmt.Errorf("ERROR: STATUS MUST BE RECORDED")
		}
	}

	return nil
}
