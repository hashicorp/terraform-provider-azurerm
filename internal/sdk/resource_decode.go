// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Decode will decode the Terraform Schema into the specified object consisting of Supported Go native Types
// These are: int64, float64, string, bool, as well as lists and maps of these base types.
// NOTE: this object must be passed by value - and must contain `tfschema`
// struct tags for all fields
//
// Example Usage:
//
//	type Person struct {
//		 Name string `tfschema:"name"
//	}
//
// var person Person
// if err := metadata.Decode(&person); err != nil { .. }
func (rmd ResourceMetaData) Decode(input interface{}) error {
	if rmd.ResourceData == nil {
		return fmt.Errorf("ResourceData was nil")
	}
	return decodeReflectedType(input, rmd.ResourceData, rmd.serializationDebugLogger)
}

// DecodeDiff decodes the Terraform Schema into the specified object in the
// same manner as Decode, but using the ResourceDiff as a source. Intended
// for use in CustomizeDiff functions.
func (rmd ResourceMetaData) DecodeDiff(input interface{}) error {
	if rmd.ResourceDiff == nil {
		return fmt.Errorf("ResourceDiff was nil")
	}
	return decodeReflectedType(input, rmd.ResourceDiff, rmd.serializationDebugLogger)
}

// stateRetriever is a convenience wrapper around the Plugin SDK to be able to test it more accurately
type stateRetriever interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
	GetOkExists(key string) (interface{}, bool)
}

func decodeReflectedType(input interface{}, stateRetriever stateRetriever, debugLogger Logger) error {
	if reflect.TypeOf(input).Kind() != reflect.Ptr {
		return fmt.Errorf("need a pointer")
	}

	objType := reflect.TypeOf(input).Elem()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		debugLogger.Infof("Field: %#v", field.Name)

		structTags, err := parseStructTags(field.Tag)
		if err != nil {
			return fmt.Errorf("parsing struct tags for %q: %+v", field.Name, err)
		}

		if structTags != nil {
			tfschemaValue, valExists := stateRetriever.GetOkExists(structTags.hclPath)
			if !valExists {
				continue
			}

			debugLogger.Infof("TFSchemaValue: %+v", tfschemaValue)
			debugLogger.Infof("Input Type: %+v", reflect.ValueOf(input).Elem().Field(i).Type())

			if err := setValue(input, tfschemaValue, i, field.Name, debugLogger); err != nil {
				return fmt.Errorf("while setting value %+v of model field %q: %+v", tfschemaValue, field.Name, err)
			}
		}
	}
	return nil
}

func setValue(input, tfschemaValue interface{}, index int, fieldName string, debugLogger Logger) (errOut error) {
	debugLogger.Infof("setting value for %q..", fieldName)
	defer func() {
		if r := recover(); r != nil {
			debugLogger.Warnf("error setting value for %q: %+v", fieldName, r)
			out, ok := r.(error)
			if !ok {
				return
			}

			errOut = out
		}
	}()

	if v, ok := tfschemaValue.(string); ok {
		n := reflect.ValueOf(input).Elem().Field(index)
		if n.Kind() == reflect.Pointer {
			debugLogger.Infof("*[String] Decode %+v", v)
			tmp := reflect.New(n.Type().Elem())
			tmp.Elem().SetString(v)
			n.Set(tmp)
		} else {
			debugLogger.Infof("[String] Decode %+v", v)
			debugLogger.Infof("Input %+v", reflect.ValueOf(input))
			debugLogger.Infof("Input Elem %+v", reflect.ValueOf(input).Elem())
			n.SetString(v)
		}
		return nil
	}

	if v, ok := tfschemaValue.(int); ok {
		n := reflect.ValueOf(input).Elem().Field(index)
		if n.Kind() == reflect.Pointer {
			debugLogger.Infof("*[INT] Decode %+v", v)
			tmp := reflect.New(n.Type().Elem())
			tmp.Elem().Set(reflect.ValueOf(int64(v)))
			n.Set(tmp)
		} else {
			debugLogger.Infof("[INT] Decode %+v", v)
			n.SetInt(int64(v))
		}
		return nil
	}

	if v, ok := tfschemaValue.(float64); ok {
		n := reflect.ValueOf(input).Elem().Field(index)
		if n.Kind() == reflect.Pointer {
			debugLogger.Infof("*[Float] Decode %+v", v)
			tmp := reflect.New(n.Type().Elem())
			tmp.Elem().SetFloat(v)
			n.Set(tmp)
		} else {
			debugLogger.Infof("[Float] Decode %+v", v)
			n.SetFloat(v)
		}
		return nil
	}

	if v, ok := tfschemaValue.(bool); ok {
		n := reflect.ValueOf(input).Elem().Field(index)
		if n.Kind() == reflect.Pointer {
			debugLogger.Infof("*[Bool] Decode %+v", v)
			tmp := reflect.New(n.Type().Elem())
			tmp.Elem().Set(reflect.ValueOf(v))
			n.Set(tmp)
		} else {
			debugLogger.Infof("[BOOL] Decode %+v", v)
			n.Set(reflect.ValueOf(v))
		}
		return nil
	}

	if v, ok := tfschemaValue.(*schema.Set); ok {
		return setListValue(input, index, fieldName, v.List(), debugLogger)
	}

	if mapConfig, ok := tfschemaValue.(map[string]interface{}); ok {
		n := reflect.ValueOf(input).Elem().Field(index)
		if n.Kind() == reflect.Pointer {
			tmp := reflect.New(n.Type().Elem())
			ty := reflect.Indirect(tmp)

			mapOutput := reflect.MakeMap(ty.Type())
			for key, val := range mapConfig {
				switch t := val.(type) {
				case int:
					mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(int64(t)))

				default:
					mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				}
			}

			tmp.Elem().Set(mapOutput)

			n.Set(tmp)
		} else {
			mapOutput := reflect.MakeMap(n.Type())
			for key, val := range mapConfig {
				switch t := val.(type) {
				case int:
					mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(int64(t)))
				default:
					mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				}
			}

			reflect.ValueOf(input).Elem().Field(index).Set(mapOutput)
		}
		return nil
	}

	if mapConfig, ok := tfschemaValue.(*map[string]interface{}); ok {
		n := reflect.ValueOf(input).Elem().Field(index).Type()

		tmp := reflect.New(n.Elem())
		ty := reflect.Indirect(tmp)

		mapOutput := reflect.MakeMap(ty.Type())
		for key, val := range *mapConfig {
			// GetOkExists always returns map[string]int rather than map[string]int64 - so we need to ensure we don't try to set the wrong type into the map
			switch mapOutput.Type().Elem().Kind() {
			case reflect.Int64:
				// Guard against GetOkExists being updated to return int64
				if v, ok := val.(int); ok {
					mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(int64(v)))
				} else {
					mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				}
			default:
				mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
			}
		}

		tmp.Elem().Set(mapOutput)
		loc := reflect.ValueOf(input).Elem().Field(index)

		loc.Set(tmp)
		return nil
	}

	if v, ok := tfschemaValue.([]interface{}); ok {
		return setListValue(input, index, fieldName, v, debugLogger)
	}

	return nil
}

func setListValue(input interface{}, index int, fieldName string, v []interface{}, debugLogger Logger) error {
	fieldType := reflect.ValueOf(input).Elem().Field(index).Type()
	var slice reflect.Value
	if reflect.TypeOf(input).Elem().Field(index).Type.Kind() != reflect.Ptr {
		slice = reflect.MakeSlice(reflect.TypeOf(input).Elem().Field(index).Type, len(v), len(v))
	} else {
		slice = reflect.MakeSlice(fieldType.Elem(), len(v), len(v))
	}
	isPtr := fieldType.Kind() == reflect.Ptr
	var dereferenceFieldType reflect.Kind
	if isPtr {
		dereferenceFieldType = fieldType.Elem().Elem().Kind()
	} else {
		dereferenceFieldType = fieldType.Elem().Kind()
	}
	switch dereferenceFieldType {
	case reflect.String:
		if !isPtr {
			for i, stringVal := range v {
				slice.Index(i).SetString(stringVal.(string))
			}
			reflect.ValueOf(input).Elem().Field(index).Set(slice)
		} else {
			tmp := reflect.New(fieldType.Elem())
			for i, stringVal := range v {
				slice.Index(i).SetString(stringVal.(string))
			}
			tmp.Elem().Set(slice)
			reflect.ValueOf(input).Elem().Field(index).Set(tmp)
		}

	case reflect.Int64:
		if !isPtr {
			for i, iVal := range v {
				slice.Index(i).SetInt(int64(iVal.(int)))
			}
			reflect.ValueOf(input).Elem().Field(index).Set(slice)
		} else {
			tmp := reflect.New(fieldType.Elem())
			for i, iVal := range v {
				slice.Index(i).SetInt(int64(iVal.(int)))
			}
			tmp.Elem().Set(slice)
			reflect.ValueOf(input).Elem().Field(index).Set(tmp)
		}

	case reflect.Float64:
		if !isPtr {
			for i, fVal := range v {
				slice.Index(i).SetFloat(fVal.(float64))
			}
			reflect.ValueOf(input).Elem().Field(index).Set(slice)
		} else {
			tmp := reflect.New(fieldType.Elem())
			for i, fVal := range v {
				slice.Index(i).SetFloat(fVal.(float64))
			}
			tmp.Elem().Set(slice)
			reflect.ValueOf(input).Elem().Field(index).Set(tmp)
		}

	case reflect.Bool:
		if !isPtr {
			for i, bVal := range v {
				slice.Index(i).SetBool(bVal.(bool))
			}
			reflect.ValueOf(input).Elem().Field(index).Set(slice)
		} else {
			tmp := reflect.New(fieldType.Elem())
			for i, bVal := range v {
				slice.Index(i).SetBool(bVal.(bool))
			}
			tmp.Elem().Set(slice)
			reflect.ValueOf(input).Elem().Field(index).Set(tmp)

		}

	default:
		n := reflect.ValueOf(input).Elem().Field(index)
		if n.Kind() == reflect.Pointer {
			tmp := reflect.New(fieldType.Elem())
			valueToSet := reflect.MakeSlice(tmp.Elem().Type(), 0, 0)
			for _, mapVal := range v {
				if test, ok := mapVal.(map[string]interface{}); ok && test != nil {
					elem := reflect.New(fieldType.Elem().Elem())
					debugLogger.Infof("element %s", elem.String())
					for j := 0; j < elem.Type().Elem().NumField(); j++ {
						nestedField := elem.Type().Elem().Field(j)
						debugLogger.Infof("nestedField Name: '%s', Tags: '%+v'", nestedField.Name, nestedField.Tag)

						structTags, err := parseStructTags(nestedField.Tag)
						if err != nil {
							return fmt.Errorf("parsing struct tags for nested field `%s`: %+v", nestedField.Name, err)
						}

						if structTags != nil {
							nestedTFSchemaValue := test[structTags.hclPath]
							if err := setValue(elem.Interface(), nestedTFSchemaValue, j, fieldName, debugLogger); err != nil {
								return err
							}
						}
					}

					if !elem.CanSet() {
						elem = elem.Elem()
					}

					if valueToSet.Kind() == reflect.Ptr {
						valueToSet.Elem().Set(reflect.Append(valueToSet.Elem(), elem))
					} else {
						valueToSet = reflect.Append(valueToSet, elem)
					}

					debugLogger.Infof("value to set type after changes %s", valueToSet.Type().String())
				}
			}

			tmp.Elem().Set(valueToSet)
			n.Set(tmp)
		} else {
			valueToSet := reflect.MakeSlice(n.Type(), 0, 0)
			debugLogger.Infof("List Type '%s'", valueToSet.Type().String())

			for _, mapVal := range v {
				if test, ok := mapVal.(map[string]interface{}); ok && test != nil {
					elem := reflect.New(fieldType.Elem())
					debugLogger.Infof("element '%s'", elem.String())
					for j := 0; j < elem.Type().Elem().NumField(); j++ {
						nestedField := elem.Type().Elem().Field(j)
						debugLogger.Infof("nestedField Name '%s', Tags: '%+v'", nestedField.Name, nestedField.Tag)

						structTags, err := parseStructTags(nestedField.Tag)
						if err != nil {
							return fmt.Errorf("parsing struct tags for nested field '%s': %+v", nestedField.Name, err)
						}

						if structTags != nil {
							nestedTFSchemaValue := test[structTags.hclPath]
							if err := setValue(elem.Interface(), nestedTFSchemaValue, j, nestedField.Name, debugLogger); err != nil {
								return err
							}
						}
					}

					if !elem.CanSet() {
						elem = elem.Elem()
					}

					if valueToSet.Kind() == reflect.Ptr {
						valueToSet.Elem().Set(reflect.Append(valueToSet.Elem(), elem))
					} else {
						valueToSet = reflect.Append(valueToSet, elem)
					}

					debugLogger.Infof("value to set type after changes '%s'", valueToSet.Type().String())
				}
			}

			valueToSet = reflect.Indirect(valueToSet)
			fieldToSet := reflect.ValueOf(input).Elem().Field(index)
			fieldToSet.Set(reflect.Indirect(valueToSet))
		}
	}

	return nil
}
