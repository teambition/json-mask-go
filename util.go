package jsonmask

import (
	"reflect"
	"strings"
	"sync"

	"gopkg.in/oleiade/reflections.v1"
)

var (
	structTagCache = sync.Map{}
)

// getFiledNamesByJSONKeys gets the map of key and JSON tag of the given struct.
func getFiledNamesByJSONKeys(obj interface{}, jsonKeys []string) (map[string]string, error) {
	var (
		fieldNames  = make(map[string]string)
		fieldTagMap = make(map[string]string)
		structName  = reflect.TypeOf(obj).String()
		err         error
	)

	// try to load from cache at first
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
		// ignore the content after first ','
		// example: `json: "_id,omitempty"`
		keys := strings.Split(key, ",")
		if len(keys) == 0 {
			continue
		}
		key := keys[0]
		for _, jsonKey := range jsonKeys {
			if key == jsonKey {
				fieldNames[jsonKey] = fieldName
				break
			}
		}
	}

	return fieldNames, nil
}

// getMaskKeys returns all keys in the nodeMap.
func getMaskKeys(mask nodeMap) []string {
	keys := make([]string, 0)

	for key := range mask {
		keys = append(keys, key)
	}

	return keys
}

// containsMapKey checks whether the given key exists in the given map.
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
