package sdk

import (
	"fmt"
	"reflect"
)

// Encode will encode the specified object into the Terraform State
// NOTE: this requires that the object passed in is a pointer and
// all fields contain `hcl` struct tags
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
			return fmt.Errorf("settting %q: %+v", k, err)
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

	output = make(map[string]interface{}, 0)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldVal := objVal.Field(i)
		if hclTag, exists := field.Tag.Lookup("hcl"); exists {
			switch field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv := fieldVal.Int()
				debugLogger.Infof("Setting %q to %d", hclTag, iv)

				output[hclTag] = iv

			case reflect.Float32, reflect.Float64:
				fv := fieldVal.Float()
				debugLogger.Infof("Setting %q to %f", hclTag, fv)

				output[hclTag] = fv

			case reflect.String:
				sv := fieldVal.String()
				debugLogger.Infof("Setting %q to %q", hclTag, sv)
				output[hclTag] = sv

			case reflect.Bool:
				bv := fieldVal.Bool()
				debugLogger.Infof("Setting %q to %t", hclTag, bv)
				output[hclTag] = bv

			case reflect.Map:
				iter := fieldVal.MapRange()
				attr := make(map[string]interface{})
				for iter.Next() {
					attr[iter.Key().String()] = iter.Value().Interface()
				}
				output[hclTag] = attr

			case reflect.Slice:
				sv := fieldVal.Slice(0, fieldVal.Len())
				attr := make([]interface{}, sv.Len())
				switch sv.Type() {
				case reflect.TypeOf([]string{}):
					debugLogger.Infof("Setting %q to []string", hclTag)
					if sv.Len() > 0 {
						output[hclTag] = sv.Interface()
					} else {
						output[hclTag] = make([]string, 0)
					}

				case reflect.TypeOf([]int{}):
					debugLogger.Infof("Setting %q to []int", hclTag)
					if sv.Len() > 0 {
						output[hclTag] = sv.Interface()
					} else {
						output[hclTag] = make([]int, 0)
					}

				case reflect.TypeOf([]float64{}):
					debugLogger.Infof("Setting %q to []float64", hclTag)
					if sv.Len() > 0 {
						output[hclTag] = sv.Interface()
					} else {
						output[hclTag] = make([]float64, 0)
					}

				case reflect.TypeOf([]bool{}):
					debugLogger.Infof("Setting %q to []bool", hclTag)
					if sv.Len() > 0 {
						output[hclTag] = sv.Interface()
					} else {
						output[hclTag] = make([]bool, 0)
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
					debugLogger.Infof("[SLICE] Setting %q to %+v", hclTag, attr)
					output[hclTag] = attr
				}
			default:
				return output, fmt.Errorf("unknown type %+v for key %q", field.Type.Kind(), hclTag)
			}
		}
	}

	return output, nil
}
