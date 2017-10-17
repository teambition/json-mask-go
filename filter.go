package jsonmask

import (
	"reflect"
)

type arrayWrapper struct {
	K interface{} `json:"k"`
}

func filter(obj interface{}, mask nodeMap) (interface{}, error) {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice, reflect.Array:
		filterdObj, err := filterProps(
			arrayWrapper{K: obj},
			nodeMap{"k": node{typ: typeArray, props: mask}},
		)
		if err != nil {
			return nil, err
		}

		return reflect.ValueOf(filterdObj).Field(0).Interface(), nil
	default:
		return filterProps(obj, mask)
	}
}

func filterProps(obj interface{}, mask nodeMap) (interface{}, error) {
	newFields := make([]reflect.StructField, 0)

	for key, _ := range mask {
		field, ok := getFiledByJSONKey(obj, key)
		if !ok {
			continue
		}

		newFields = append(newFields, *field)
	}

	return reflect.ValueOf(obj).Convert(reflect.StructOf(newFields)).Interface(), nil
}
