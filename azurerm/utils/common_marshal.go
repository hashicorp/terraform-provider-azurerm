package utils

import "reflect"

func ExpandStringSlice(input []interface{}) *[]string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return &result
}

func ExpandFloatSlice(input []interface{}) *[]float64 {
	result := make([]float64, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(float64))
		}
	}
	return &result
}

func ExpandMapStringPtrString(input map[string]interface{}) map[string]*string {
	result := make(map[string]*string)
	for k, v := range input {
		result[k] = String(v.(string))
	}
	return result
}

func ExpandInt32Slice(input []interface{}) *[]int32 {
	result := make([]int32, len(input))
	for i, item := range input {
		result[i] = int32(item.(int))
	}

	return &result
}

func FlattenStringSlice(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func FlattenFloatSlice(input *[]float64) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func FlattenMapStringPtrString(input map[string]*string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range input {
		if v == nil {
			result[k] = ""
		} else {
			result[k] = *v
		}
	}
	return result
}

func FlattenInt32Slice(input *[]int32) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

// ExpandSlice expands the input slice into pointer to slice of element whose type is specified by "t".
// If "t" is different from the element type of input, then user has to specify a customized converter via
// "convert", which guides the conversion from element of the input slice to the element of the output slice.
// Otherwise, user can pass a nil "convert".
func ExpandSlice(input []interface{}, t interface{}, convert func(interface{}) interface{}) interface{} {
	targetType := reflect.TypeOf(t)
	result := reflect.MakeSlice(reflect.SliceOf(targetType), 0, 0)
	for _, v := range input {
		var vv reflect.Value
		if v != nil {
			if convert == nil {
				vv = reflect.New(targetType).Elem()
				vv.Set(reflect.ValueOf(v))
			} else {
				vv = reflect.ValueOf(convert(v))
			}
		} else {
			vv = underlyingZeroValue(targetType)
		}
		result = reflect.Append(result, vv)
	}

	resultp := reflect.New(result.Type())
	resultp.Elem().Set(result)
	return resultp.Interface()
}

// ExpandMap expands the input map (key is of type string) into map (key is of type string) where the type of value is specified by "t".
// If "t" is different from the element type of input, then user has to specify a customized converter via
// "convert", which guides the conversion from value of the input map to the value of the output map.
// Otherwise, user can pass a nil "convert".
func ExpandMap(input map[string]interface{}, t interface{}, convert func(interface{}) interface{}) interface{} {
	targetType := reflect.TypeOf(t)
	result := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(""), targetType))
	for k, v := range input {
		var vv reflect.Value
		if v != nil {
			if convert == nil {
				vv = reflect.New(targetType).Elem()
				vv.Set(reflect.ValueOf(v))
			} else {
				vv = reflect.ValueOf(convert(v))
			}
		} else {
			vv = underlyingZeroValue(targetType)
		}
		result.SetMapIndex(reflect.ValueOf(k), vv)
	}
	return result.Interface()
}

// FlattenSlice flattens the input pointer to slice of a certain type ("t"), into slice of type of `interface{}`
// If "t" is different from the element type of output, then user has to specify a customized converter via
// "convert", which guides the conversion from element of the input slice to the element of the output slice.
// Otherwise, user can pass a nil "convert".
func FlattenSlicePtr(input interface{}, convert func(interface{}) interface{}) []interface{} {
	v := reflect.ValueOf(input)
	// safe guard
	if v.Type().Kind() != reflect.Ptr {
		panic("Invalid input: input is not a pointer")
	}

	ve := v.Elem()
	if ve.Type().Kind() != reflect.Slice {
		panic("Invalid input: value of input is not a slice")
	}

	result := make([]interface{}, 0)
	if v.IsNil() {
		return result
	}

	for i := 0; i < ve.Len(); i++ {
		v := ve.Index(i)
		var ov reflect.Value
		if isNilable(v) && v.IsNil() {
			ov = underlyingZeroValue(ve.Type().Elem())
		} else {
			if convert == nil {
				ov = v
			} else {
				ov = reflect.ValueOf(convert(v.Interface()))
			}
		}
		result = append(result, ov.Interface())
	}
	return result
}

// FlattenStringMap flattens the input map (key is of type string), whose value is of a certain type "t", into map (key is of type string), whose value is of type of `interface{}`.
// If "t" is different from the value type of output map, then user has to specify a customized converter via
// "convert", which guides the conversion from value of the input map to the value of the output map.
// Otherwise, user can pass a nil "convert".
func FlattenStringMap(input interface{}, convert func(interface{}) interface{}) map[string]interface{} {
	v := reflect.ValueOf(input)
	// safe guard
	if v.Type().Kind() != reflect.Map {
		panic("Invalid input: input is not a map")
	}
	if v.Type().Key().Kind() != reflect.String {
		panic("Invalid input: key of input map is not a string")
	}

	result := make(map[string]interface{}, 0)
	iter := v.MapRange()
	for iter.Next() {
		k, v := iter.Key(), iter.Value()
		var ov reflect.Value

		if isNilable(v) && v.IsNil() {
			ov = underlyingZeroValue(v.Type().Elem()) // Elem() of a map type will give you the value's type!
		} else {
			if convert == nil {
				ov = v
			} else {
				ov = reflect.ValueOf(convert(v.Interface()))
			}
		}
		result[k.Interface().(string)] = ov.Interface()
	}
	return result
}

// isNilable check whether the passed in value is fine to call the `Nil()` method.
// The supported kinds are got from the source code of `Nil()` itself.
func isNilable(v reflect.Value) bool {
	switch v.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

// underlyingZeroValue returns the zero value of the passed in type. Especially, if the "t" is a pointer
// then it will return the zero value of the de-referenced value, recursively.
func underlyingZeroValue(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		return underlyingZeroValue(t.Elem())
	}
	return reflect.Zero(t)
}
