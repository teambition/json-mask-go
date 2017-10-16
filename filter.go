package jsonmask

import (
	"reflect"
)

func filter(jsonObj interface{}, compiledMask nodeMap) interface{} {
	switch reflect.TypeOf(jsonObj).Kind() {
	case reflect.Slice, reflect.Array:
		return filterProps(arrayWrapper{
			K: jsonObj,
		}, nodeMap{
			"k": node{typ: typeArray, props: compiledMask},
		})
	default:
		return filterProps(jsonObj, compiledMask)
	}
}

func filterProps(obj interface{}, mask nodeMap) interface{} {
	return nil
}

type arrayWrapper struct {
	K interface{} `json:"k"`
}
