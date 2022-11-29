package models

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
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

type Membership struct {
	gorm.Model
	TeamID        uint           `json:"team_id" gorm:"uniqueIndex:membership_index"`
	ParticipantID uint           `json:"participant_id" gorm:"uniqueIndex:membership_index"`
	Role          MembershipRole `json:"role" gorm:"type:membership_role;not null"`
	Team          Team           `json:"team" gorm:"foreignKey:TeamID;references:ID"`
	Participant   Participant    `json:"participant" gorm:"foreignKey:ParticipantID;references:ID"`
}

// Menambahkan constraint untuk mengecek apakah terdapat participant yang mengikuti dua team atau lebih
// dengan jenis lomba yang sama atau memiliki role leader lebih dari satu kali
func (membership *Membership) BeforeSave(tx *gorm.DB) error {
	condition1 := Membership{ParticipantID: membership.ParticipantID}
	oldMemberships1 := []Membership{}
	if err := tx.Where(&condition1).Find(&oldMemberships1).Error; err != nil {
		return err
	}

	for _, oldMembership := range oldMemberships1 {
		if oldMembership.Team.TeamCategory != "" && oldMembership.Team.TeamCategory == membership.Team.TeamCategory {
			return fmt.Errorf("ERROR: CANNOT PARTICIPATE MORE THAN ONCE")
		}
		if oldMembership.Role == Leader && membership.Role == Leader {
			return fmt.Errorf("ERROR: INELIGIBLE LEADER")
		}
	}

	condition2 := Membership{ParticipantID: membership.TeamID}
	oldMemberships2 := []Membership{}
	if err := tx.Where(&condition2).Find(&oldMemberships2).Error; err != nil {
		return err
	}

	for _, oldMembership := range oldMemberships2 {
		if oldMembership.Role == Leader && membership.Role == Leader {
			return fmt.Errorf("ERROR: INELIGIBLE LEADER")
		}
	}

	return nil
}
