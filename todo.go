package main

import (
	"reflect"
	"strings"
)

type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type UpdateTodo struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Done        *bool   `json:"done,omitempty"`
}

func (t *Todo) Toggle() bool {
	t.Done = !t.Done
	return t.Done
}

type Todos []Todo

func UpdateFields(target interface{}, source interface{}) []string {
	targetVal := reflect.ValueOf(target).Elem()
	sourceVal := reflect.ValueOf(source).Elem()
	sourceType := reflect.TypeOf(source).Elem()

	var updatedFields []string

	for i := 0; i < sourceVal.NumField(); i++ {
		sourceField := sourceVal.Field(i)
		fieldName := sourceType.Field(i).Name

		// Check if it's a pointer and not nil
		if sourceField.Kind() == reflect.Ptr && !sourceField.IsNil() {
			// Find corresponding field in target
			targetField := targetVal.FieldByName(fieldName)
			if targetField.IsValid() && targetField.CanSet() {
				// Dereference the pointer and set the value
				targetField.Set(sourceField.Elem())
				updatedFields = append(updatedFields, strings.ToLower(fieldName))
			}
		}
	}

	return updatedFields
}
