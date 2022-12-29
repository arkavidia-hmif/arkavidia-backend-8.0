package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Submission struct {
	gorm.Model
	FileName      uuid.UUID             `json:"file_name" gorm:"type:uuid;unique"`
	FileExtension string                `json:"file_extension" gorm:"not null"`
	TeamID        uint                  `json:"team_id" gorm:"not null"`
	Stage         types.SubmissionStage `json:"stage" gorm:"not null"`
	Team          Team                  `json:"team" gorm:"foreignKey:TeamID;references:ID"`
}
