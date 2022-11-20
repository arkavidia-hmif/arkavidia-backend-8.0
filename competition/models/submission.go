package models

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionStage string

const (
	FirstStage  SubmissionStage = "first-stage"
	SecondStage SubmissionStage = "second-stage"
	FinalStage  SubmissionStage = "final-stage"
)

func (submissionStage *SubmissionStage) Scan(value interface{}) error {
	*submissionStage = SubmissionStage(value.([]byte))
	return nil
}

func (submissionStage SubmissionStage) Value() (driver.Value, error) {
	return string(submissionStage), nil
}

type Submission struct {
	gorm.Model
	FileName      uuid.UUID       `json:"file_name" gorm:"type:uuid;unique"`
	FileExtension string          `json:"file_extension" gorm:"not null"`
	TeamID        uint            `json:"team_id" gorm:"not null;uniqueIndex:submission_index"`
	Stage         SubmissionStage `json:"stage" gorm:"type:submission_stage;not null;uniqueIndex:submission_index"`
	Team          Team            `json:"team" gorm:"foreignKey:TeamID;references:ID"`
}
