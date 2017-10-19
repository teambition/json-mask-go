package jsonmask

import "reflect"

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

	for key, node := range mask {
		field, ok := getFiledByJSONKey(obj, key)
		if !ok {
			continue
		}

		var value interface{}
		if node.props != nil {
			switch field.Type.Kind() {
			case reflect.Slice, reflect.Array:
				if field.Type.Len() == 0 {
					continue
				}

				value = reflect.ValueOf(obj).FieldByName(field.Name).Index(0).Interface()
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
	}

	return newFields, nil
}
