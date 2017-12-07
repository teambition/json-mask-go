package jsonmask

import (
	"reflect"
)

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
