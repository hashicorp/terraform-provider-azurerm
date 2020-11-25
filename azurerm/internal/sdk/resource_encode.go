package sdk

import (
	"fmt"
	"reflect"
)

// Encode will encode the specified object into the Terraform State
// NOTE: this requires that the object passed in is a pointer and
// all fields contain `tfschema` struct tags
func (rmd ResourceMetaData) Encode(input interface{}) error {
	if reflect.TypeOf(input).Kind() != reflect.Ptr {
		return fmt.Errorf("need a pointer")
	}

	objType := reflect.TypeOf(input).Elem()
	objVal := reflect.ValueOf(input).Elem()

	fieldName := reflect.ValueOf(input).Elem().String()
	serialized, err := recurse(objType, objVal, fieldName, rmd.serializationDebugLogger)
	if err != nil {
		return err
	}

	for k, v := range serialized {
		if err := rmd.ResourceData.Set(k, v); err != nil {
			return fmt.Errorf("setting %q: %+v", k, err)
		}
	}
	return nil
}

func recurse(objType reflect.Type, objVal reflect.Value, fieldName string, debugLogger Logger) (output map[string]interface{}, errOut error) {
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

	output = make(map[string]interface{})
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldVal := objVal.Field(i)
		if tfschemaTag, exists := field.Tag.Lookup("tfschema"); exists {
			switch field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv := fieldVal.Int()
				debugLogger.Infof("Setting %q to %d", tfschemaTag, iv)
				output[tfschemaTag] = iv

			case reflect.Float32, reflect.Float64:
				fv := fieldVal.Float()
				debugLogger.Infof("Setting %q to %f", tfschemaTag, fv)
				output[tfschemaTag] = fv

			case reflect.String:
				sv := fieldVal.String()
				debugLogger.Infof("Setting %q to %q", tfschemaTag, sv)
				output[tfschemaTag] = sv

			case reflect.Bool:
				bv := fieldVal.Bool()
				debugLogger.Infof("Setting %q to %t", tfschemaTag, bv)
				output[tfschemaTag] = bv

			case reflect.Map:
				iter := fieldVal.MapRange()
				attr := make(map[string]interface{})
				for iter.Next() {
					attr[iter.Key().String()] = iter.Value().Interface()
				}
				output[tfschemaTag] = attr

			case reflect.Slice:
				sv := fieldVal.Slice(0, fieldVal.Len())
				attr := make([]interface{}, sv.Len())
				switch sv.Type() {
				case reflect.TypeOf([]string{}):
					debugLogger.Infof("Setting %q to []string", tfschemaTag)
					if sv.Len() > 0 {
						output[tfschemaTag] = sv.Interface()
					} else {
						output[tfschemaTag] = make([]string, 0)
					}

				case reflect.TypeOf([]int{}):
					debugLogger.Infof("Setting %q to []int", tfschemaTag)
					if sv.Len() > 0 {
						output[tfschemaTag] = sv.Interface()
					} else {
						output[tfschemaTag] = make([]int, 0)
					}

				case reflect.TypeOf([]float64{}):
					debugLogger.Infof("Setting %q to []float64", tfschemaTag)
					if sv.Len() > 0 {
						output[tfschemaTag] = sv.Interface()
					} else {
						output[tfschemaTag] = make([]float64, 0)
					}

				case reflect.TypeOf([]bool{}):
					debugLogger.Infof("Setting %q to []bool", tfschemaTag)
					if sv.Len() > 0 {
						output[tfschemaTag] = sv.Interface()
					} else {
						output[tfschemaTag] = make([]bool, 0)
					}

				default:
					for i := 0; i < sv.Len(); i++ {
						debugLogger.Infof("[SLICE] Index %d is %q", i, sv.Index(i).Interface())
						debugLogger.Infof("[SLICE] Type %+v", sv.Type())
						nestedType := sv.Index(i).Type()
						nestedValue := sv.Index(i)

						fieldName := field.Name
						serialized, err := recurse(nestedType, nestedValue, fieldName, debugLogger)
						if err != nil {
							return nil, fmt.Errorf("serializing nested object %q: %+v", sv.Type(), exists)
						}
						attr[i] = serialized
					}
					debugLogger.Infof("[SLICE] Setting %q to %+v", tfschemaTag, attr)
					output[tfschemaTag] = attr
				}
			default:
				return output, fmt.Errorf("unknown type %+v for key %q", field.Type.Kind(), tfschemaTag)
			}
		}
	}

	return output, nil
}
