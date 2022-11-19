package models

import (
	"github.com/google/uuid"
)

type PhotoStatus string

const (
	NotUploaded            PhotoStatus = "Not Uploaded"
	WaitingForVerification PhotoStatus = "Waiting for Verification"
	Verified               PhotoStatus = "Verified"
	Declined               PhotoStatus = "Declined"
)

type PhotoType string

const (
	Pribadi             PhotoType = "Pribadi"
	KartuPelajar        PhotoType = "Kartu Pelajar"
	BuktiMahasiswaAktif PhotoType = "Bukti Mahasiswa Aktif"
)

type Photo struct {
	ID            uuid.UUID   `json:"photo_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ParticipantID uuid.UUID   `json:"participant_id" gorm:"type:uuid;not null"`
	LinkToFile    uuid.UUID   `json:"link_to_file" gorm:"type:uuid;not null"`
	FileExtension string      `json:"file_extension" gorm:"not null"`
	Status        PhotoStatus `json:"status" gorm:"not null"`
	Type          PhotoType   `json:"type" gorm:"not null"`
	Participant   Participant `json:"-" gorm:"foreignKey:ParticipantID;references:ID"`
}
