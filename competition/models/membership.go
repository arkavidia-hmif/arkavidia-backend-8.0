package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipRole string

const (
	Leader    MembershipRole = "leader"
	MemberOne MembershipRole = "member-1"
	MemberTwo MembershipRole = "member-2"
)

func (membershipRole *MembershipRole) Scan(value interface{}) error {
	*membershipRole = MembershipRole(value.([]byte))
	return nil
}

func (membershipRole MembershipRole) Value() (driver.Value, error) {
	return string(membershipRole), nil
}

type Membership struct {
	gorm.Model
	TeamID        uuid.UUID      `json:"team_id" gorm:"type:uuid;uniqueIndex:membership_index"`
	ParticipantID uuid.UUID      `json:"participant_id" gorm:"type:uuid;uniqueIndex:membership_index"`
	Role          MembershipRole `json:"role" gorm:"type:membership_role;not null"`
	Team          Team           `json:"-" gorm:"foreignKey:TeamID;references:UUID"`
	Participant   Participant    `json:"-" gorm:"foreignKey:ParticipantID;references:UUID"`
}

// Menambahkan constraint untuk mengecek apakah terdapat participant yang mengikuti dua team atau lebih
// dengan jenis lomba yang sama atau memiliki role leader lebhi dari satu kali
func (membership *Membership) BeforeSave(tx *gorm.DB) error {
	condition := Membership{ParticipantID: membership.ParticipantID}
	newMemberships := []Membership{}
	if err := tx.Where(&condition).Find(&newMemberships).Error; err != nil {
		return err
	}

	for _, membershipA := range newMemberships {
		for _, membershipB := range newMemberships {
			if membershipA.TeamID != membershipB.TeamID {
				if membershipA.Team.TeamCategory == membershipB.Team.TeamCategory {
					return fmt.Errorf("Error: Invalid Database Operation!")
				}
				if membershipA.Role == Leader && membershipB.Role == Leader {
					return fmt.Errorf("Error: Invalid Database Operation!")
				}
			}
		}
	}

	return nil
}
