package types

import (
	"database/sql/driver"
)

type ParticipantStatus string

const (
	WaitingForVerification ParticipantStatus = "waiting-for-verification"
	Verified               ParticipantStatus = "verified"
	Declined               ParticipantStatus = "declined"
)

func (participantStatus *ParticipantStatus) Scan(value interface{}) error {
	*participantStatus = ParticipantStatus(value.(string))
	return nil
}

func (participantStatus ParticipantStatus) Value() (driver.Value, error) {
	return string(participantStatus), nil
}

func (ParticipantStatus) GormDataType() string {
	return "participant_status"
}
