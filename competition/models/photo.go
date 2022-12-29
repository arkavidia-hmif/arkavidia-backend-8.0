package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Photo struct {
	gorm.Model
	FileName      uuid.UUID         `json:"file_name" gorm:"type:uuid;unique"`
	FileExtension string            `json:"file_extension" gorm:"not null"`
	ParticipantID uint              `json:"participant_id" gorm:"not null"`
	AdminID       uint              `json:"admin_id" gorm:"default:null"`
	Status        types.PhotoStatus `json:"status" gorm:"not null"`
	Type          types.PhotoType   `json:"type" gorm:"not null"`
	Participant   Participant       `json:"participant" gorm:"foreignKey:ParticipantID;references:ID"`
	ApprovedBy    Admin             `json:"admin" gorm:"foreignKey:AdminID;references:ID"`
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
