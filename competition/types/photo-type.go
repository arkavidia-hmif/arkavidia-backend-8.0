package types

import (
	"database/sql/driver"
)

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
