package sdk

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidateModelObject validates that the object contains the specified `tfschema` tags
// required to be used with the Encode and Decode functions
func ValidateModelObject(input interface{}) error {
	if reflect.TypeOf(input).Kind() != reflect.Ptr {
		return fmt.Errorf("need a pointer")
	}

	// TODO: could we also validate that each `tfschema` tag exists in the schema?

	objType := reflect.TypeOf(input).Elem()
	objVal := reflect.ValueOf(input).Elem()
	return validateModelObjectRecursively("", objType, objVal)
}

func validateModelObjectRecursively(prefix string, objType reflect.Type, objVal reflect.Value) (errOut error) {
	defer func() {
		if r := recover(); r != nil {
			out, ok := r.(error)
			if !ok {
				return
			}

			errOut = out
		}
	}()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldVal := objVal.Field(i)

		if field.Type.Kind() == reflect.Slice {
			sv := fieldVal.Slice(0, fieldVal.Len())
			innerType := sv.Type().Elem()
			innerVal := reflect.Indirect(reflect.New(innerType))
			fieldName := strings.TrimPrefix(fmt.Sprintf("%s.%s", prefix, field.Name), ".")
			if err := validateModelObjectRecursively(fieldName, innerType, innerVal); err != nil {
				return err
			}
		}

		if _, exists := field.Tag.Lookup("tfschema"); !exists {
			fieldName := strings.TrimPrefix(fmt.Sprintf("%s.%s", prefix, field.Name), ".")
			return fmt.Errorf("field %q is missing an `tfschema` label", fieldName)
		}
	}

	return nil
}
