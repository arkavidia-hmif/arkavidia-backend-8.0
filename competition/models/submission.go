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
	*submissionStage = SubmissionStage(value.(string))
	return nil
}

func (submissionStage SubmissionStage) Value() (driver.Value, error) {
	return string(submissionStage), nil
}

func (SubmissionStage) GormDataType() string {
	return "submission_stage"
}

type Submission struct {
	gorm.Model
	FileName      uuid.UUID       `json:"file_name" gorm:"type:uuid;unique"`
	FileExtension string          `json:"file_extension" gorm:"not null"`
	TeamID        uint            `json:"team_id" gorm:"not null"`
	Stage         SubmissionStage `json:"stage" gorm:"not null"`
	Team          Team            `json:"team" gorm:"foreignKey:TeamID;references:ID"`
}
