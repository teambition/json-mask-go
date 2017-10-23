package jsonmask

import (
	"reflect"
)

func getFiledByJSONKey(obj interface{}, jsonKey string) (*reflect.StructField, int, bool) {
	objType := reflect.TypeOf(obj)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		if field.Tag.Get("json") == jsonKey {
			return &field, i, true
		}
	}

	return nil, 0, false
}

func stringsContains(strings []string, s string) bool {
	for _, str := range strings {
		if s == str {
			return true
		}
	}

	return false
}
