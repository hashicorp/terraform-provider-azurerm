package helpers

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func expandSlice[I, O any](input []interface{}, convert func(I) O, appendOnNil bool) *[]O {
	result := make([]O, 0, len(input))
	for _, item := range input {
		if item != nil {
			result = append(result, convert(item.(I)))
		} else if appendOnNil {
			var zero O
			result = append(result, zero)
		}
	}
	return &result
}

func flattenSlice[T any](input *[]T) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func expandMap[I, O any](input map[string]interface{}, convert func(I) O) map[string]O {
	result := make(map[string]O, len(input))
	for k, v := range input {
		result[k] = convert(v.(I))
	}
	return result
}

func expandSliceWithDelimiter[I any](input []interface{}, convert func(I) string, delimiter string) *string {
	result := make([]string, 0, len(input))
	for _, item := range input {
		if item != nil {
			result = append(result, convert(item.(I)))
		} else {
			result = append(result, "")
		}
	}
	return pointer.To(strings.Join(result, delimiter))
}

func ExpandStringSlice(input []interface{}) *[]string {
	return expandSlice(input, func(i string) string { return i }, true)
}

func ExpandFloatSlice(input []interface{}) *[]float64 {
	return expandSlice(input, func(i float64) float64 { return i }, false)
}

func ExpandFloatRangeSlice(input []interface{}) *[][]float64 {
	return expandSlice(input, func(i []interface{}) []float64 { return *ExpandFloatSlice(i) }, false)
}

func ExpandPtrMapStringString(input map[string]interface{}) *map[string]string {
	res := expandMap(input, func(i string) string { return i })
	return &res
}

func ExpandMapStringPtrString(input map[string]interface{}) map[string]*string {
	return expandMap(input, func(i string) *string { return pointer.To(i) })
}

func ExpandInt32Slice(input []interface{}) *[]int32 {
	return expandSlice(input, func(i int) int32 { return int32(i) }, true)
}

func ExpandInt64Slice(input []interface{}) *[]int64 {
	return expandSlice(input, func(i int) int64 { return int64(i) }, true)
}

func FlattenStringSlice(input *[]string) []interface{} {
	return flattenSlice(input)
}

func FlattenFloatSlice(input *[]float64) []interface{} {
	return flattenSlice(input)
}

func FlattenFloatRangeSlice(input *[][]float64) [][]interface{} {
	result := make([][]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, FlattenFloatSlice(&item))
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

func FlattenPtrMapStringString(input *map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	if input == nil {
		return result
	}
	for k, v := range *input {
		result[k] = v
	}
	return result
}

func FlattenInt32Slice(input *[]int32) []interface{} {
	return flattenSlice(input)
}

func FlattenInt64Slice(input *[]int64) []interface{} {
	return flattenSlice(input)
}

func ExpandStringSliceWithDelimiter(input []interface{}, delimiter string) *string {
	return expandSliceWithDelimiter(input, func(i string) string { return i }, delimiter)
}

func ExpandIntSliceWithDelimiter(input []interface{}, delimiter string) *string {
	return expandSliceWithDelimiter(input, strconv.Itoa, delimiter)
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
