package utils

import (
	"reflect"
)

func GetMissingFields(input interface{}) []string {
	//validate input
	missingFields := []string{}
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag != "" {
				missingFields = append(missingFields, jsonTag)
			} else {
				missingFields = append(missingFields, fieldType.Name)
			}
		}
	}
	
	return missingFields
}