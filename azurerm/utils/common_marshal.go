package utils

import "strings"

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

func ExpandStringSliceWithDelimiter(input []interface{}, delimiter string) *string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return String(strings.Join(result, delimiter))
}

func FlattenStringSliceWithDelimiter(input *string, delimiter string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		inputStrings := strings.Split(*input, delimiter)
		for _, item := range inputStrings {
			result = append(result, item)
		}
	}
	return result
}
