package jsonmask

import (
	"reflect"
	"sync"

	"gopkg.in/oleiade/reflections.v1"
)

var (
	structTagCache = sync.Map{}
)

func getFiledNamesByJSONKeys(obj interface{}, jsonKeys []string) (map[string]string, error) {
	var (
		fieldNames  = make(map[string]string)
		structName  = reflect.TypeOf(obj).Name()
		fieldTagMap = map[string]string{}
		err         error
	)

	cache, ok := structTagCache.Load(structName)
	if !ok {
		fieldTagMap, err = reflections.TagsDeep(obj, "json")
		if err != nil {
			return nil, err
		}

		structTagCache.Store(structName, fieldTagMap)
	} else {
		fieldTagMap = cache.(map[string]string)
	}

	for fieldName, key := range fieldTagMap {
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
