package models

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PhotoStatus string

const (
	WaitingForVerification PhotoStatus = "waiting-for-verification"
	Verified               PhotoStatus = "verified"
	Declined               PhotoStatus = "declined"
)

func (photoStatus *PhotoStatus) Scan(value interface{}) error {
	*photoStatus = PhotoStatus(value.(string))
	return nil
}

func (photoStatus PhotoStatus) Value() (driver.Value, error) {
	return string(photoStatus), nil
}

type PhotoType string

const (
	Pribadi             PhotoType = "pribadi"
	KartuPelajar        PhotoType = "kartu-pelajar"
	BuktiMahasiswaAktif PhotoType = "bukti-mahasiswa-aktif"
	BuktiPembayaran 		PhotoType = "bukti-pembayaran"
)

func (photoType *PhotoType) Scan(value interface{}) error {
	*photoType = PhotoType(value.(string))
	return nil
}

func (photoType PhotoType) Value() (driver.Value, error) {
	return string(photoType), nil
}

type Photo struct {
	gorm.Model
	FileName      uuid.UUID   `json:"file_name" gorm:"type:uuid;unique"`
	FileExtension string      `json:"file_extension" gorm:"not null"`
	ParticipantID uint        `json:"participant_id" gorm:"not null"`
	Status        PhotoStatus `json:"status" gorm:"type:photo_status;not null"`
	Type          PhotoType   `json:"type" gorm:"type:photo_type;not null"`
	Participant   Participant `json:"participant" gorm:"foreignKey:ParticipantID;references:ID"`
}
