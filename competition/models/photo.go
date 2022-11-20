package models

import (
	"github.com/google/uuid"
)

type PhotoStatus string

const (
	NotUploaded            PhotoStatus = "not-uploaded"
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
	FileName      uuid.UUID   `json:"file_name" gorm:"type:uuid;primaryKey"`
	FileExtension string      `json:"file_extension" gorm:"not null"`
	ParticipantID uuid.UUID   `json:"participant_id" gorm:"type:uuid;not null"`
	Status        PhotoStatus `json:"status" gorm:"not null"`
	Type          PhotoType   `json:"type" gorm:"not null"`
	Participant   Participant `json:"-" gorm:"foreignKey:ParticipantID;references:ID"`
}
