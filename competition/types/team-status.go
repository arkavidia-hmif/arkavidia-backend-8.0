package types

import (
	"database/sql/driver"
)

type TeamStatus string

const (
	WaitingForEvaluation TeamStatus = "waiting-for-evaluation"
	Passed               TeamStatus = "passed"
	Eliminated           TeamStatus = "eliminated"
)

func (teamStatus *TeamStatus) Scan(value interface{}) error {
	*teamStatus = TeamStatus(value.(string))
	return nil
}

func (teamStatus TeamStatus) Value() (driver.Value, error) {
	return string(teamStatus), nil
}

func (TeamStatus) GormDataType() string {
	return "team_status"
}
