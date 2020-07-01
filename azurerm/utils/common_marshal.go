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
	for _, item := range input {
		if item != nil {
			var vitem reflect.Value
			if convert == nil {
				vitem = reflect.New(targetType).Elem()
				vitem.Set(reflect.ValueOf(item))
			} else {
				vitem = reflect.ValueOf(convert(item))
			}
			result = reflect.Append(result, vitem)
		} else {
			result = reflect.Append(result, reflect.Zero(targetType))
		}
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
		if convert == nil {
			vv = reflect.New(targetType).Elem()
			vv.Set(reflect.ValueOf(v))
		} else {
			vv = reflect.ValueOf(convert(v))
		}
		result.SetMapIndex(reflect.ValueOf(k), vv)
	}
	return result.Interface()
}

// FlattenSlice flattens the input pointer to slice of element whose type is specified by "t", into slice.
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
		ev := ve.Index(i)
		if convert != nil {
			ev = reflect.ValueOf(convert(ev.Interface()))
		}
		result = append(result, ev.Interface())
	}
	return result
}
