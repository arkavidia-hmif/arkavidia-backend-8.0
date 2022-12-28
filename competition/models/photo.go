package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

type PhotoType string

const (
	Pribadi             PhotoType = "pribadi"
	KartuPelajar        PhotoType = "kartu-pelajar"
	BuktiMahasiswaAktif PhotoType = "bukti-mahasiswa-aktif"
	BuktiPembayaran     PhotoType = "bukti-pembayaran"
)

func (photoType *PhotoType) Scan(value interface{}) error {
	*photoType = PhotoType(value.(string))
	return nil
}

func (photoType PhotoType) Value() (driver.Value, error) {
	return string(photoType), nil
}

func (PhotoType) GormDataType() string {
	return "photo_type"
}

type Photo struct {
	gorm.Model
	FileName      uuid.UUID   `json:"file_name" gorm:"type:uuid;unique"`
	FileExtension string      `json:"file_extension" gorm:"not null"`
	ParticipantID uint        `json:"participant_id" gorm:"not null"`
	AdminID       uint        `json:"admin_id" gorm:"default:null"`
	Status        PhotoStatus `json:"status" gorm:"not null"`
	Type          PhotoType   `json:"type" gorm:"not null"`
	Participant   Participant `json:"participant" gorm:"foreignKey:ParticipantID;references:ID"`
	ApprovedBy    Admin       `json:"admin" gorm:"foreignKey:AdminID;references:ID"`
}

// Menambahkan constraint untuk mengecek apakah terdapat photos yang telah diapprove namun admin tidak tercatat
// atau photos yang belum diapprove namun admin tercatat
func (photo *Photo) BeforeSave(tx *gorm.DB) error {
	if photo.Status != "" {
		if photo.AdminID == 0 {
			return fmt.Errorf("ERROR: ADMIN MUST BE RECORDED")
		}
		if photo.Status == WaitingForApproval {
			return fmt.Errorf("ERROR: STATUS MUST BE RECORDED")
		}
	}

	return nil
}
