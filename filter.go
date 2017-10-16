package jsonmask

import (
	"reflect"

	"gopkg.in/oleiade/reflections.v1"
)

func filter(jsonObj interface{}, compiledMask nodeMap) (interface{}, error) {
	switch reflect.TypeOf(jsonObj).Kind() {
	case reflect.Slice, reflect.Array:
		return filterProps(arrayWrapper{K: jsonObj}, nodeMap{"k": node{typ: typeArray, props: compiledMask}})
	default:
		return filterProps(jsonObj, compiledMask)
	}
}

func filterProps(obj interface{}, mask nodeMap) (interface{}, error) {
	objFields, err := reflections.Fields(obj)
	if err != nil {
		return nil, err
	}

	for _, fieldName := range objFields {
		if _, ok := mask[fieldName]; !ok {
			obj, err = clearJSONTag(obj, fieldName)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

type arrayWrapper struct {
	K interface{} `json:"k"`
}
