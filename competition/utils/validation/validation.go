package validation

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"

	"arkavidia-backend-8.0/competition/utils/sanitizer"
)

// TODO: Tambahkan validasi payload request dengan menggunakan validator
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-http-request-payload-validation.html
// ASSIGNED TO: @rayhankinan
// STATUS: IN PROGRESS

type CustomValidator struct {
	validate *validator.Validate
	once     sync.Once
}

// Private
func (customValidator *CustomValidator) lazyInit() {
	customValidator.once.Do(func() {
		customValidator.validate = validator.New()
		customValidator.validate.SetTagName("binding")
	})
}

// Public
func (customValidator *CustomValidator) ValidateStruct(obj interface{}) error {
	if sanitizer.KindOfData(obj) == reflect.Struct {
		customValidator.lazyInit()
		if err := customValidator.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (customValidator *CustomValidator) Engine() interface{} {
	customValidator.lazyInit()
	return customValidator.validate
}

var Validator = &CustomValidator{}
