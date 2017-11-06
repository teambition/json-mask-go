package jsonmask

import (
	"reflect"

	"github.com/fatih/structtag"
)

type arrayWrapper struct {
	K interface{} `json:"k"`
}

func filter(obj interface{}, mask nodeMap) (interface{}, error) {
	var (
		newStructFileds []reflect.StructField
		err             error
	)

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
	objType := reflect.TypeOf(obj)
	newFields := make([]reflect.StructField, objType.NumField())

	for key, node := range mask {
		field, pos, ok := getFiledByJSONKey(obj, key)
		if !ok {
			continue
		}

		if node.props != nil && len(node.props) != 0 {
			switch field.Type.Kind() {
			case reflect.Slice, reflect.Array:
				sliceValue := reflect.ValueOf(obj).FieldByName(field.Name)

				if sliceValue.Len() == 0 {
					continue
				}

				subFields, err := filterProps(sliceValue.Index(0).Interface(), node.props)
				if err != nil {
					return nil, err
				}

				subStruct := reflect.StructOf(subFields)

				subSliceStruct := reflect.SliceOf(subStruct)

				field.Type = reflect.SliceOf(subSliceStruct)
			default:
				subFields, err := filterProps(reflect.ValueOf(obj).FieldByName(field.Name).Interface(), node.props)
				if err != nil {
					return nil, err
				}

				field.Type = reflect.StructOf(subFields)
			}
		}

		newFields[pos] = *field
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		if newFields[i].Name != "" {
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

		newFields[i] = field
	}

	return newFields, nil
}
