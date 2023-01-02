package types

import (
	"database/sql/driver"

	"golang.org/x/crypto/bcrypt"
)

type EncryptedString []byte

func (es *EncryptedString) Scan(value interface{}) error {
	*es = EncryptedString(value.([]byte))
	return nil
}

func (es EncryptedString) Value() (driver.Value, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(es, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func (EncryptedString) GormDataType() string {
	return "bytea"
}
