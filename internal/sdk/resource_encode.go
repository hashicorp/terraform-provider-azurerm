// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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

	serialized, err := recurse(objType, objVal, rmd.serializationDebugLogger)
	if err != nil {
		return err
	}

	for k, v := range serialized {
		// lintignore:R001
		if err := rmd.ResourceData.Set(k, v); err != nil {
			return fmt.Errorf("setting %q: %+v", k, err)
		}
	}
	return nil
}

func recurse(objType reflect.Type, objVal reflect.Value, debugLogger Logger) (output map[string]interface{}, errOut error) {
	var fieldName string
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
		fieldName = field.Name
		fieldVal := objVal.Field(i)
		structTags, err := parseStructTags(field.Tag)
		if err != nil {
			return nil, fmt.Errorf("parsing struct tags for %q: %+v", field.Name, err)
		}

		if structTags != nil {
			if structTags.removedInNextMajorVersion && features.FivePointOh() {
				debugLogger.Infof("The HCL Path %q is marked as removed - skipping", structTags.hclPath)
				continue
			}

			if structTags.addedInNextMajorVersion && !features.FivePointOh() {
				debugLogger.Infof("The HCL Path %q is marked as not yet present - skipping", structTags.hclPath)
				continue
			}

			switch field.Type.Kind() {
			case reflect.Int64:
				iv := fieldVal.Int()
				debugLogger.Infof("Setting %q to %d", structTags.hclPath, iv)
				output[structTags.hclPath] = iv

			case reflect.Float64:
				fv := fieldVal.Float()
				debugLogger.Infof("Setting %q to %f", structTags.hclPath, fv)
				output[structTags.hclPath] = fv

			case reflect.String:
				sv := fieldVal.String()
				debugLogger.Infof("Setting %q to %q", structTags.hclPath, sv)
				output[structTags.hclPath] = sv

			case reflect.Bool:
				bv := fieldVal.Bool()
				debugLogger.Infof("Setting %q to %t", structTags.hclPath, bv)
				output[structTags.hclPath] = bv

			case reflect.Map:
				iter := fieldVal.MapRange()
				attr := make(map[string]interface{})
				for iter.Next() {
					attr[iter.Key().String()] = iter.Value().Interface()
				}
				output[structTags.hclPath] = attr

			case reflect.Slice:
				sv := fieldVal.Slice(0, fieldVal.Len())
				attr := make([]interface{}, sv.Len())
				switch sv.Type().Elem().Kind() {
				case reflect.String:
					debugLogger.Infof("Setting %q to []string", structTags.hclPath)
					if sv.Len() > 0 {
						output[structTags.hclPath] = sv.Interface()
					} else {
						output[structTags.hclPath] = make([]string, 0)
					}

				case reflect.Int64:
					debugLogger.Infof("Setting %q to []int", structTags.hclPath)
					if sv.Len() > 0 {
						output[structTags.hclPath] = sv.Interface()
					} else {
						output[structTags.hclPath] = make([]int64, 0)
					}

				case reflect.Float64:
					debugLogger.Infof("Setting %q to []float64", structTags.hclPath)
					if sv.Len() > 0 {
						output[structTags.hclPath] = sv.Interface()
					} else {
						output[structTags.hclPath] = make([]float64, 0)
					}

				case reflect.Bool:
					debugLogger.Infof("Setting %q to []bool", structTags.hclPath)
					if sv.Len() > 0 {
						output[structTags.hclPath] = sv.Interface()
					} else {
						output[structTags.hclPath] = make([]bool, 0)
					}

				default:
					for i := 0; i < sv.Len(); i++ {
						debugLogger.Infof("[SLICE] Index %d is %q", i, sv.Index(i).Interface())
						debugLogger.Infof("[SLICE] Type %+v", sv.Type())
						nestedType := sv.Index(i).Type()
						nestedValue := sv.Index(i)

						serialized, err := recurse(nestedType, nestedValue, debugLogger)
						if err != nil {
							return nil, fmt.Errorf("serializing nested object %q: %+v", sv.Type(), err)
						}
						attr[i] = serialized
					}
					debugLogger.Infof("[SLICE] Setting %q to %+v", structTags.hclPath, attr)
					output[structTags.hclPath] = attr
				}

			case reflect.Pointer:
				if !fieldVal.IsNil() {
					pv := fieldVal.Elem()
					switch pv.Kind() {
					case reflect.Int, reflect.Int64:
						iv := pv.Int()
						debugLogger.Infof("Setting %q to %d", structTags.hclPath, iv)
						output[structTags.hclPath] = iv

					case reflect.Float64:
						fv := pv.Float()
						debugLogger.Infof("Setting %q to %f", structTags.hclPath, fv)
						output[structTags.hclPath] = fv

					case reflect.String:
						sv := pv.String()
						debugLogger.Infof("Setting %q to %q", structTags.hclPath, sv)
						output[structTags.hclPath] = sv

					case reflect.Bool:
						bv := pv.Bool()
						debugLogger.Infof("Setting %q to %t", structTags.hclPath, bv)
						output[structTags.hclPath] = bv

					case reflect.Map:
						iter := pv.MapRange()
						attr := make(map[string]interface{})
						for iter.Next() {
							attr[iter.Key().String()] = iter.Value().Interface()
						}
						output[structTags.hclPath] = attr

					case reflect.Slice:
						sv := pv.Slice(0, pv.Len())
						attr := make([]interface{}, sv.Len())
						switch sv.Type().Elem().Kind() {
						case reflect.String:
							debugLogger.Infof("Setting %q to []string", structTags.hclPath)
							if sv.Len() > 0 {
								output[structTags.hclPath] = sv.Interface()
							} else {
								output[structTags.hclPath] = make([]string, 0)
							}

						case reflect.Int64:
							debugLogger.Infof("Setting %q to []int", structTags.hclPath)
							if sv.Len() > 0 {
								output[structTags.hclPath] = sv.Interface()
							} else {
								output[structTags.hclPath] = make([]int64, 0)
							}

						case reflect.Float64:
							debugLogger.Infof("Setting %q to []float64", structTags.hclPath)
							if sv.Len() > 0 {
								output[structTags.hclPath] = sv.Interface()
							} else {
								output[structTags.hclPath] = make([]float64, 0)
							}

						case reflect.Bool:
							debugLogger.Infof("Setting %q to []bool", structTags.hclPath)
							if sv.Len() > 0 {
								output[structTags.hclPath] = sv.Interface()
							} else {
								output[structTags.hclPath] = make([]bool, 0)
							}

						default:
							for i := 0; i < sv.Len(); i++ {
								debugLogger.Infof("[SLICE] Index %d is %q", i, sv.Index(i).Interface())
								debugLogger.Infof("[SLICE] Type %+v", sv.Type())
								nestedType := sv.Index(i).Type()
								nestedValue := sv.Index(i)

								serialized, err := recurse(nestedType, nestedValue, debugLogger)
								if err != nil {
									return nil, fmt.Errorf("serializing nested object %q: %+v", sv.Type(), err)
								}
								attr[i] = serialized
							}
							debugLogger.Infof("[SLICE] Setting %q to %+v", structTags.hclPath, attr)
							output[structTags.hclPath] = attr
						}
					}
				} else {
					debugLogger.Infof("Setting %q to nil", structTags.hclPath)
					output[structTags.hclPath] = nil
				}

			default:
				return output, fmt.Errorf("unknown type %+v for key %q", field.Type.Kind(), structTags.hclPath)
			}
		}
	}

	return output, nil
}
