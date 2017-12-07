package jsonmask

import (
	"reflect"

	"gopkg.in/oleiade/reflections.v1"
)

func filter(obj interface{}, mask nodeMap) (interface{}, error) {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice, reflect.Array:
		len := reflect.ValueOf(obj).Len()
		filtered := make([]map[string]interface{}, len)

		for i := 0; i < len; i++ {
			val, err := filterProps(reflect.ValueOf(obj).Index(i).Interface(), mask)
			if err != nil {
				return nil, err
			}

			filtered[i] = val
		}

		return filtered, nil
	default:
		filtered, err := filterProps(obj, mask)
		if err != nil {
			return nil, err
		}

		return filtered, nil
	}
}

func filterProps(obj interface{}, mask nodeMap) (map[string]interface{}, error) {
	filteredMap := make(map[string]interface{})

	for key := range mask {
		field, ok := getFiledByJSONKey(obj, key)
		if !ok {
			continue
		}

		val, err := reflections.GetField(obj, field.Name)
		if err != nil {
			return nil, err
		}

		filteredMap[key] = val
	}

	return filteredMap, nil
}
