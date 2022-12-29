package types

import (
	"database/sql/driver"
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
