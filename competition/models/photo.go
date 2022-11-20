package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PhotoStatus string

const (
	WaitingForVerification PhotoStatus = "waiting-for-verification"
	Verified               PhotoStatus = "verified"
	Declined               PhotoStatus = "declined"
)

type PhotoType string

const (
	Pribadi             PhotoType = "pribadi"
	KartuPelajar        PhotoType = "kartu-pelajar"
	BuktiMahasiswaAktif PhotoType = "bukti-mahasiswa-aktif"
)

type Photo struct {
	gorm.Model
	FileName      uuid.UUID   `json:"file_name" gorm:"type:uuid;unique"`
	FileExtension string      `json:"file_extension" gorm:"not null"`
	ParticipantID uuid.UUID   `json:"participant_id" gorm:"type:uuid;not null;uniqueIndex:photo_index"`
	Status        PhotoStatus `json:"status" gorm:"not null"`
	Type          PhotoType   `json:"type" gorm:"not null;uniqueIndex:photo_index"`
	Participant   Participant `json:"-" gorm:"foreignKey:ParticipantID;references:ID"`
}
