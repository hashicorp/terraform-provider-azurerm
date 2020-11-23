package sdk

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Decode will decode the Terraform Schema into the specified object
// NOTE: this object must be passed by value - and must contain `tfschema`
// struct tags for all fields
//
// Example Usage:
//
// type Person struct {
//	 Name string `tfschema:"name"
// }
// var person Person
// if err := metadata.Decode(&person); err != nil { .. }
func (rmd ResourceMetaData) Decode(input interface{}) error {
	return decodeReflectedType(input, rmd.ResourceData, rmd.serializationDebugLogger)
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
		debugLogger.Infof("Field", field)

		if val, exists := field.Tag.Lookup("tfschema"); exists {
			tfschemaValue, valExists := stateRetriever.GetOkExists(val)
			if !valExists {
				continue
			}

			debugLogger.Infof("TFSchemaValue: ", tfschemaValue)
			debugLogger.Infof("Input Type: ", reflect.ValueOf(input).Elem().Field(i).Type())

			fieldName := reflect.ValueOf(input).Elem().Field(i).String()
			if err := setValue(input, tfschemaValue, i, fieldName, debugLogger); err != nil {
				return err
			}
		}
	}
	return nil
}

func setValue(input, tfschemaValue interface{}, index int, fieldName string, debugLogger Logger) (errOut error) {
	debugLogger.Infof("setting list value for %q..", fieldName)
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
		debugLogger.Infof("[String] Decode %+v", v)
		debugLogger.Infof("Input %+v", reflect.ValueOf(input))
		debugLogger.Infof("Input Elem %+v", reflect.ValueOf(input).Elem())
		reflect.ValueOf(input).Elem().Field(index).SetString(v)
		return nil
	}

	if v, ok := tfschemaValue.(int); ok {
		debugLogger.Infof("[INT] Decode %+v", v)
		reflect.ValueOf(input).Elem().Field(index).SetInt(int64(v))
		return nil
	}

	if v, ok := tfschemaValue.(int32); ok {
		debugLogger.Infof("[INT] Decode %+v", v)
		reflect.ValueOf(input).Elem().Field(index).SetInt(int64(v))
		return nil
	}

	if v, ok := tfschemaValue.(int64); ok {
		debugLogger.Infof("[INT] Decode %+v", v)
		reflect.ValueOf(input).Elem().Field(index).SetInt(v)
		return nil
	}

	if v, ok := tfschemaValue.(float64); ok {
		debugLogger.Infof("[Float] Decode %+v", v)
		reflect.ValueOf(input).Elem().Field(index).SetFloat(v)
		return nil
	}

	// Doesn't work for empty bools?
	if v, ok := tfschemaValue.(bool); ok {
		debugLogger.Infof("[BOOL] Decode %+v", v)

		reflect.ValueOf(input).Elem().Field(index).SetBool(v)
		return nil
	}

	if v, ok := tfschemaValue.(*schema.Set); ok {
		return setListValue(input, index, fieldName, v.List(), debugLogger)
	}

	if mapConfig, ok := tfschemaValue.(map[string]interface{}); ok {
		mapOutput := reflect.MakeMap(reflect.ValueOf(input).Elem().Field(index).Type())
		for key, val := range mapConfig {
			mapOutput.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
		}

		reflect.ValueOf(input).Elem().Field(index).Set(mapOutput)
		return nil
	}

	if v, ok := tfschemaValue.([]interface{}); ok {
		return setListValue(input, index, fieldName, v, debugLogger)
	}

	return nil
}

func setListValue(input interface{}, index int, fieldName string, v []interface{}, debugLogger Logger) error {
	switch fieldType := reflect.ValueOf(input).Elem().Field(index).Type(); fieldType {
	case reflect.TypeOf([]string{}):
		stringSlice := reflect.MakeSlice(reflect.TypeOf([]string{}), len(v), len(v))
		for i, stringVal := range v {
			stringSlice.Index(i).SetString(stringVal.(string))
		}
		reflect.ValueOf(input).Elem().Field(index).Set(stringSlice)

	case reflect.TypeOf([]int{}):
		iSlice := reflect.MakeSlice(reflect.TypeOf([]int{}), len(v), len(v))
		for i, iVal := range v {
			iSlice.Index(i).SetInt(int64(iVal.(int)))
		}
		reflect.ValueOf(input).Elem().Field(index).Set(iSlice)

	case reflect.TypeOf([]float64{}):
		fSlice := reflect.MakeSlice(reflect.TypeOf([]float64{}), len(v), len(v))
		for i, fVal := range v {
			fSlice.Index(i).SetFloat(fVal.(float64))
		}
		reflect.ValueOf(input).Elem().Field(index).Set(fSlice)

	case reflect.TypeOf([]bool{}):
		bSlice := reflect.MakeSlice(reflect.TypeOf([]bool{}), len(v), len(v))
		for i, bVal := range v {
			bSlice.Index(i).SetBool(bVal.(bool))
		}
		reflect.ValueOf(input).Elem().Field(index).Set(bSlice)

	default:
		valueToSet := reflect.MakeSlice(reflect.ValueOf(input).Elem().Field(index).Type(), 0, 0)
		debugLogger.Infof("List Type", valueToSet.Type())

		for _, mapVal := range v {
			if test, ok := mapVal.(map[string]interface{}); ok && test != nil {
				elem := reflect.New(fieldType.Elem())
				debugLogger.Infof("element ", elem)
				for j := 0; j < elem.Type().Elem().NumField(); j++ {
					nestedField := elem.Type().Elem().Field(j)
					debugLogger.Infof("nestedField ", nestedField)

					if val, exists := nestedField.Tag.Lookup("tfschema"); exists {
						nestedTFSchemaValue := test[val]
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

				debugLogger.Infof("value to set type after changes", valueToSet.Type())
			}
		}

		valueToSet = reflect.Indirect(valueToSet)
		fieldToSet := reflect.ValueOf(input).Elem().Field(index)
		fieldToSet.Set(reflect.Indirect(valueToSet))
	}

	return nil
}
