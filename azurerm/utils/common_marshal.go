package utils

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

func FlattenStringSlice(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func ExpandBoolSlice(input []interface{}) *[]bool {
	result := make([]bool, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(bool))
		} else {
			result = append(result, false)
		}
	}
	return &result
}

func FlattenBoolSlice(input *[]bool) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func ExpandUintSlice(input []interface{}) *[]uint {
	result := make([]uint, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, uint(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenUintSlice(input *[]uint) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandUint8Slice(input []interface{}) *[]uint8 {
	result := make([]uint8, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, uint8(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenUint8Slice(input *[]uint8) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandUint16Slice(input []interface{}) *[]uint16 {
	result := make([]uint16, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, uint16(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenUint16Slice(input *[]uint16) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandUint32Slice(input []interface{}) *[]uint32 {
	result := make([]uint32, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, uint32(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenUint32Slice(input *[]uint32) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandUint64Slice(input []interface{}) *[]uint64 {
	result := make([]uint64, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, uint64(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenUint64Slice(input *[]uint64) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandIntSlice(input []interface{}) *[]int {
	result := make([]int, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(int))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenIntSlice(input *[]int) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func ExpandInt8Slice(input []interface{}) *[]int8 {
	result := make([]int8, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, int8(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenInt8Slice(input *[]int8) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandInt16Slice(input []interface{}) *[]int16 {
	result := make([]int16, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, int16(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenInt16Slice(input *[]int16) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandInt32Slice(input []interface{}) *[]int32 {
	result := make([]int32, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, int32(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenInt32Slice(input *[]int32) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandInt64Slice(input []interface{}) *[]int64 {
	result := make([]int64, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, int64(item.(int)))
		} else {
			result = append(result, 0)
		}
	}
	return &result
}

func FlattenInt64Slice(input *[]int64) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, int(item))
		}
	}
	return result
}

func ExpandFloat32Slice(input []interface{}) *[]float32 {
	result := make([]float32, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, float32(item.(float64)))
		} else {
			result = append(result, 0.0)
		}
	}
	return &result
}

func FlattenFloat32Slice(input *[]float32) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, float64(item))
		}
	}
	return result
}

func ExpandFloat64Slice(input []interface{}) *[]float64 {
	result := make([]float64, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(float64))
		} else {
			result = append(result, 0.0)
		}
	}
	return &result
}

func FlattenFloat64Slice(input *[]float64) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}
