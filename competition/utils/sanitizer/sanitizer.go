package sanitizer

// TODO: Tambahkan fungsi dengan menggunakan reflection untuk mengecek apakah relation dibutuhkan atau tidak
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/A-reflect.html
// ASSIGNED TO: @graceclaudia19
// STATUS: IN PROGRESS

import (
	"reflect"
)

func KindOfData(obj interface{}) reflect.Kind {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}

	return valueType
}
