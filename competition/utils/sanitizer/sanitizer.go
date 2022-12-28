package sanitizer

// TODO: Tambahkan fungsi dengan menggunakan reflection untuk mengecek apakah relation dibutuhkan atau tidak
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/A-reflect.html
// ASSIGNED TO: @graceclaudia19
// STATUS: DONE

import (
	"reflect"
)

type Response[T interface{}] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
	URL     string `json:"URL"`
}

func SanitizeArray(obj interface{}) []interface{} {
	objValue := reflect.ValueOf(obj)
	responses := []interface{}{}

	for i := 0; i < objValue.Len(); i++ {
		childObjValue := objValue.Index(i)

		if childObjValue.CanInterface() && !reflect.DeepEqual(childObjValue.Interface(), reflect.Zero(childObjValue.Type()).Interface()) {
			switch childObjValue.Kind() {
			case reflect.Array:
				{
					responses = append(responses, SanitizeArray(childObjValue.Interface())...)
				}
			case reflect.Struct:
				{
					responses = append(responses, SanitizeStruct(childObjValue.Interface()))
				}
			default:
				{
					responses = append(responses, childObjValue.Interface())
				}
			}
		}
	}

	return responses
}

func SanitizeStruct(obj interface{}) map[string]interface{} {
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)
	response := map[string]interface{}{}

	for i := 0; i < objValue.NumField(); i++ {
		childObjValue := objValue.Field(i)
		childObjType := objType.Field(i)

		if childObjType.Tag.Get("json") != "" && childObjType.Tag.Get("visibility") != "false" && childObjValue.CanInterface() && !reflect.DeepEqual(childObjValue.Interface(), reflect.Zero(childObjValue.Type()).Interface()) {
			switch childObjValue.Kind() {
			case reflect.Slice:
				{
					response[childObjType.Tag.Get("json")] = SanitizeArray(childObjValue.Interface())
				}
			case reflect.Struct:
				{
					response[childObjType.Tag.Get("json")] = SanitizeStruct(childObjValue.Interface())
				}
			default:
				{
					response[childObjType.Tag.Get("json")] = childObjValue.Interface()
				}
			}
		}

	}

	return response
}

func KindOfData(obj interface{}) reflect.Kind {
	objValue := reflect.ValueOf(obj)
	objValueKind := objValue.Kind()

	if objValueKind == reflect.Ptr {
		objValueKind = objValue.Elem().Kind()
	}

	return objValueKind
}
