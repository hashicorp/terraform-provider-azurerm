package sdk

import (
	"fmt"
	"reflect"
)

// HasChanges will determine if the Terraform Attribute has changed between Terraform runs.
// NOTE: this attribute must be passed by value - and must contain `tfschema`
// struct tags
//
// Example Usage:
//
// type Person struct {
//	 Name string `tfschema:"name"
// }
// var person Person
// if err := metadata.Decode(&person); err != nil { .. }
// if metadata.HasChanges(person, "Name") {
//   ...
// }
func (rmd ResourceMetaData) HasChanges(input interface{}, fieldName string) (bool, error) {
	return hasChangesReflectedType(input, fieldName, rmd.ResourceData, rmd.serializationDebugLogger)
}

// stateDiffRetriever is a convenience wrapper around the Plugin SDK to be able to test it more accurately
type stateDiffRetriever interface {
	HasChanges(keys ...string) bool
}

func hasChangesReflectedType(input interface{}, fieldName string, stateRetriever stateDiffRetriever, debugLogger Logger) (bool, error) {
	field, ok := reflect.TypeOf(input).Elem().FieldByName(fieldName)
	if !ok {
		return false, fmt.Errorf("unable to find field %q", fieldName)
	}

	if val, exists := field.Tag.Lookup("tfschema"); exists {
		return stateRetriever.HasChanges(val), nil
	}
	return false, nil
}
