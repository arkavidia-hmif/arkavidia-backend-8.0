package types

import (
	"database/sql/driver"
)

type MembershipRole string

const (
	Leader MembershipRole = "leader"
	Member MembershipRole = "member"
)

func (membershipRole *MembershipRole) Scan(value interface{}) error {
	*membershipRole = MembershipRole(value.(string))
	return nil
}

func (membershipRole MembershipRole) Value() (driver.Value, error) {
	return string(membershipRole), nil
}

func (MembershipRole) GormDataType() string {
	return "membership_role"
}
