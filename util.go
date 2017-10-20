package jsonmask

import (
	"reflect"
)

func checkIsArray(obj interface{}) bool {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func getFiledByJSONKey(obj interface{}, jsonKey string) (*reflect.StructField, bool) {
	objType := reflect.TypeOf(obj)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		if field.Tag.Get("json") == jsonKey {
			return &field, true
		}
	}

	return nil, false
}

func stringsContains(strings []string, s string) bool {
	for _, str := range strings {
		if s == str {
			return true
		}
	}

	return false
}
