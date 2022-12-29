package types

import (
	"database/sql/driver"
)

type PhotoStatus string

const (
	WaitingForApproval PhotoStatus = "waiting-for-approval"
	Approved           PhotoStatus = "approved"
	Denied             PhotoStatus = "denied"
)

func (photoStatus *PhotoStatus) Scan(value interface{}) error {
	*photoStatus = PhotoStatus(value.(string))
	return nil
}

func (photoStatus PhotoStatus) Value() (driver.Value, error) {
	return string(photoStatus), nil
}

func (PhotoStatus) GormDataType() string {
	return "photo_status"
}
