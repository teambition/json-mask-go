package jsonmask

import (
	"reflect"

	"gopkg.in/oleiade/reflections.v1"
)

type arrayWrapper struct {
	K interface{} `json:"k"`
}

func filter(obj interface{}, mask nodeMap) (interface{}, error) {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice, reflect.Array:
		filtered := make([]map[string]interface{}, 0)

		for i := 0; i < reflect.ValueOf(obj).Len(); i++ {
			val, err := filterProps(reflect.ValueOf(obj).Index(i), mask)
			if err != nil {
				return nil, err
			}
			filtered = append(filtered, val)
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
		field, _, ok := getFiledByJSONKey(obj, key)
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
