package jsonmask

import (
	"container/list"
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

func filterProps(obj interface{}, mask nodeMap) ([]reflect.StructField, error) {
	newFields := make([]reflect.StructField, 0)

	fieldStack := list.New()

	for key, node := range mask {
		field, ok := getFiledByJSONKey(obj, key)
		if !ok {
			continue
		}

		listElement := fieldStack.PushBack(field)

		if node.props != nil {
			subFields, err := filterProps(reflect.ValueOf(obj).FieldByName(field.Name).Interface(), node.props)
			if err != nil {
				return nil, err
			}

			field.Type = reflect.StructOf(subFields)
		}

		newFields = append(newFields, *field)

		fieldStack.Remove(listElement)
	}

	return newFields, nil
}
