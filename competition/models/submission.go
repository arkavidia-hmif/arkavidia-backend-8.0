package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Submission struct {
	gorm.Model
	FileName      uuid.UUID             `gorm:"type:uuid;unique"`
	FileExtension string                `gorm:"not null"`
	TeamID        uint                  `gorm:"not null"`
	Stage         types.SubmissionStage `gorm:"not null"`
	Team          Team                  `gorm:"foreignKey:TeamID;references:ID"`
}

type DisplaySubmission struct {
	ID            uint                  `json:"id,omitempty"`
	CreatedAt     time.Time             `json:"created_at,omitempty"`
	UpdatedAt     time.Time             `json:"updated_at,omitempty"`
	FileName      uuid.UUID             `json:"file_name,omitempty" gorm:"type:uuid;unique"`
	FileExtension string                `json:"file_extension,omitempty" gorm:"not null"`
	TeamID        uint                  `json:"team_id,omitempty" gorm:"not null"`
	Stage         types.SubmissionStage `json:"stage,omitempty" gorm:"not null"`
}

func (submission Submission) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DisplaySubmission{
		ID:            submission.ID,
		CreatedAt:     submission.CreatedAt,
		UpdatedAt:     submission.UpdatedAt,
		FileName:      submission.FileName,
		FileExtension: submission.FileExtension,
		TeamID:        submission.TeamID,
		Stage:         submission.Stage,
	})
}
