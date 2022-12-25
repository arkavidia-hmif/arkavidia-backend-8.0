package validation

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		cv.lazyInit()
		if err := cv.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (cv *CustomValidator) Engine() interface{} {
	cv.lazyInit()
	return cv.validate
}

func (cv *CustomValidator) lazyInit() {
	cv.once.Do(func() {
		cv.validate = validator.New()
		cv.validate.SetTagName("binding")
	})
}

func kindOfData(obj interface{}) reflect.Kind {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}

	return valueType
}

var customValidator *CustomValidator = nil

func Init() *CustomValidator {
	return &CustomValidator{}
}

func GetValidator() *CustomValidator {
	if customValidator == nil {
		customValidator = Init()
	}
	return customValidator
}
