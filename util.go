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

func containsMapKey(sli []reflect.Value, s string) *reflect.Value {
	for _, sliItem := range sli {
		if sliItem.Kind() != reflect.String {
			continue
		}

		if sliItem.Interface().(string) == s {
			return &sliItem
		}
	}

	return nil
}
