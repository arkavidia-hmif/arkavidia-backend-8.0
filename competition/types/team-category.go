package types

import (
	"database/sql/driver"
)

type TeamCategory string

const (
	CP         TeamCategory = "competitive-programming"
	Datavidia  TeamCategory = "datavidia"
	UXVidia    TeamCategory = "uxvidia"
	Arkalogica TeamCategory = "arkalogica"
)

func (teamCategory *TeamCategory) Scan(value interface{}) error {
	*teamCategory = TeamCategory(value.(string))
	return nil
}

func (teamCategory TeamCategory) Value() (driver.Value, error) {
	return string(teamCategory), nil
}

func (TeamCategory) GormDataType() string {
	return "team_category"
}
