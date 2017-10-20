package jsonmask

import (
	"reflect"

	"github.com/fatih/structtag"
)

type arrayWrapper struct {
	K interface{} `json:"k"`
}

func filter(obj interface{}, mask nodeMap) (interface{}, error) {
	var newStructFileds []reflect.StructField
	var err error

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice, reflect.Array:
		newStructFileds, err = filterProps(
			arrayWrapper{K: obj},
			nodeMap{"k": node{typ: typeArray, props: mask}},
		)

		if err != nil {
			return nil, err
		}
	default:
		newStructFileds, err = filterProps(obj, mask)

		if err != nil {
			return nil, err
		}
	}

	return reflect.ValueOf(obj).Convert(reflect.StructOf(newStructFileds)).Interface(), nil
}

func filterProps(obj interface{}, mask nodeMap) ([]reflect.StructField, error) {
	newFields := make([]reflect.StructField, 0)
	objType := reflect.TypeOf(obj)
	newStructNames := make([]string, 0)

	for key, node := range mask {
		field, ok := getFiledByJSONKey(obj, key)
		if !ok {
			continue
		}

		var value interface{}
		if node.props != nil && len(node.props) != 0 {
			switch field.Type.Kind() {
			case reflect.Slice, reflect.Array:
				sliceValue := reflect.ValueOf(obj).FieldByName(field.Name)

				if sliceValue.Len() == 0 {
					continue
				}

				value = sliceValue.Index(0).Interface()
			default:
				value = reflect.ValueOf(obj).FieldByName(field.Name).Interface()
			}

			subFields, err := filterProps(value, node.props)
			if err != nil {
				return nil, err
			}

			switch field.Type.Kind() {
			case reflect.Slice, reflect.Array:
				field.Type = reflect.SliceOf(reflect.StructOf(subFields))
			default:
				field.Type = reflect.StructOf(subFields)
			}
		}

		newFields = append(newFields, *field)
		newStructNames = append(newStructNames, field.Name)
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		if stringsContains(newStructNames, field.Name) {
			continue
		}

		structTags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return nil, err
		}

		jsonTag, err := structTags.Get("json")
		if err != nil {
			continue // "json" tag does not exist
		}
		jsonTag.Name = "-"

		if err := structTags.Set(jsonTag); err != nil {
			return nil, err
		}

		field.Tag = reflect.StructTag(structTags.String())

		newFields = append(newFields, field)
	}

	return newFields, nil
}
