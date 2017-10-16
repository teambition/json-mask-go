package jsonmask

import (
	"reflect"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
)

func clearJSONTag(obj interface{}, key string) (interface{}, error) {
	field, err := getField(obj, key)
	if err != nil {
		return nil, err
	}

	structTags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return nil, err
	}

	jsonTag, err := structTags.Get("json")
	if err != nil {
		return obj, nil // "json" tag does not exist
	}
	jsonTag.Name = "-"

	if err := structTags.Set(jsonTag); err != nil {
		return nil, err
	}

	objType := reflect.TypeOf(obj)
	newObjFields := make([]reflect.StructField, 0)

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		if field.Name == key {
			field.Tag = reflect.StructTag(structTags.String())
		}

		newObjFields = append(newObjFields, field)
	}

	return reflect.ValueOf(obj).Convert(reflect.StructOf(newObjFields)).Interface(), nil
}

func getField(obj interface{}, fieldName string) (*reflect.StructField, error) {
	field, ok := reflectValue(obj).Type().FieldByName(fieldName)
	if !ok {
		return nil, errors.Errorf("No such field: %s in obj", fieldName)
	}

	return &field, nil
}

func reflectValue(obj interface{}) reflect.Value {
	var val reflect.Value

	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		val = reflect.ValueOf(obj).Elem()
	} else {
		val = reflect.ValueOf(obj)
	}

	return val
}
