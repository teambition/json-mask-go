package jsonmask

import (
	"reflect"

	"gopkg.in/oleiade/reflections.v1"
)

func getFiledNamesByJSONKeys(obj interface{}, jsonKeys []string) (map[string]string, error) {
	fieldNames := make(map[string]string)

	m, err := reflections.TagsDeep(obj, "json")
	if err != nil {
		return nil, err
	}

	for fieldName, key := range m {
		for _, jsonKey := range jsonKeys {
			if key == jsonKey {
				fieldNames[jsonKey] = fieldName
				break
			}
		}
	}

	return fieldNames, nil
}

func getMaskKeys(mask nodeMap) []string {
	keys := make([]string, 0)

	for key := range mask {
		keys = append(keys, key)
	}

	return keys
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

func containsString(sli []string, s string) bool {
	for _, item := range sli {
		if item == s {
			return true
		}
	}

	return false
}
