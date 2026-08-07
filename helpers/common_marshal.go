package helpers

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

// ExpandSlice safely iterates over a slice of interfaces and converts each non-nil element
// to the target type O using the provided convert function. If appendOnNil is true, it inserts
// the zero value of O for nil elements; otherwise, nil elements are skipped.
//
// Example:
//
//	input := []interface{}{"a", "b", nil, "c"}
//	result := ExpandSlice(input, func(i string) string { return i }, true)
//	// *result is []string{"a", "b", "", "c"}
func ExpandSlice[I, O any](input []interface{}, convert func(I) O, appendOnNil bool) *[]O {
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

// FlattenSlice takes a pointer to a slice of any type T and flattens it into a slice of interface{}.
//
// Example:
//
//	input := []string{"a", "b", "c"}
//	result := FlattenSlice(&input)
//	// result is []interface{}{"a", "b", "c"}
func FlattenSlice[T any](input *[]T) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

// ExpandMap converts a map of string to interface{} into a map of string to the target type O.
//
// Example:
//
//	input := map[string]interface{}{"key": "value"}
//	result := ExpandMap(input, func(i string) string { return i })
//	// result is map[string]string{"key": "value"}
func ExpandMap[I, O any](input map[string]interface{}, convert func(I) O) map[string]O {
	result := make(map[string]O, len(input))
	for k, v := range input {
		result[k] = convert(v.(I))
	}
	return result
}

// ExpandSliceWithDelimiter joins the elements of the interface slice using the provided delimiter.
// Nil elements are treated as empty strings.
//
// Example:
//
//	input := []interface{}{int(1), int(2), nil, int(3)}
//	result := ExpandSliceWithDelimiter(input, strconv.Itoa, ",")
//	// *result is "1,2,,3"
func ExpandSliceWithDelimiter[I any](input []interface{}, convert func(I) string, delimiter string) *string {
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

// ExpandStringSlice converts a slice of interface{} to a pointer to a slice of strings.
// Nil elements are converted to empty strings.
//
// Deprecated: Use ExpandSlice instead.
func ExpandStringSlice(input []interface{}) *[]string {
	return ExpandSlice(input, func(i string) string { return i }, true)
}

// ExpandFloatSlice converts a slice of interface{} to a pointer to a slice of float64s.
// Nil elements are ignored.
//
// Deprecated: Use ExpandSlice instead.
func ExpandFloatSlice(input []interface{}) *[]float64 {
	return ExpandSlice(input, func(i float64) float64 { return i }, false)
}

// ExpandFloatRangeSlice converts a slice of interface{} to a pointer to a slice of float64 slices.
// Nil elements are ignored.
//
// Deprecated: Use ExpandSlice instead.
func ExpandFloatRangeSlice(input []interface{}) *[][]float64 {
	return ExpandSlice(input, func(i []interface{}) []float64 { return *ExpandFloatSlice(i) }, false)
}

// ExpandPtrMapStringString converts a map of interface{} to a pointer to a map of strings.
//
// Deprecated: Use ExpandMap instead.
func ExpandPtrMapStringString(input map[string]interface{}) *map[string]string {
	res := ExpandMap(input, func(i string) string { return i })
	return &res
}

// ExpandMapStringPtrString converts a map of interface{} to a map of string pointers.
//
// Deprecated: Use ExpandMap instead.
func ExpandMapStringPtrString(input map[string]interface{}) map[string]*string {
	return ExpandMap(input, func(i string) *string { return pointer.To(i) })
}

// ExpandInt32Slice converts a slice of interface{} to a pointer to a slice of int32s.
// Nil elements are converted to 0.
//
// Deprecated: Use ExpandSlice instead.
func ExpandInt32Slice(input []interface{}) *[]int32 {
	return ExpandSlice(input, func(i int) int32 { return int32(i) }, true)
}

// ExpandInt64Slice converts a slice of interface{} to a pointer to a slice of int64s.
// Nil elements are converted to 0.
//
// Deprecated: Use ExpandSlice instead.
func ExpandInt64Slice(input []interface{}) *[]int64 {
	return ExpandSlice(input, func(i int) int64 { return int64(i) }, true)
}

// FlattenStringSlice converts a pointer to a slice of strings into a slice of interface{}.
//
// Deprecated: Use FlattenSlice instead.
func FlattenStringSlice(input *[]string) []interface{} {
	return FlattenSlice(input)
}

// FlattenFloatSlice converts a pointer to a slice of float64s into a slice of interface{}.
//
// Deprecated: Use FlattenSlice instead.
func FlattenFloatSlice(input *[]float64) []interface{} {
	return FlattenSlice(input)
}

// FlattenFloatRangeSlice converts a pointer to a slice of float64 slices into a slice of interface{} slices.
//
// Deprecated: Use FlattenSlice instead.
func FlattenFloatRangeSlice(input *[][]float64) [][]interface{} {
	result := make([][]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, FlattenFloatSlice(&item))
		}
	}
	return result
}

// FlattenMapStringPtrString converts a map of string pointers into a map of interface{}.
// Nil pointers are converted to empty strings.
//
// Deprecated: Use FlattenSlice or custom flatten map logic instead.
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

// FlattenPtrMapStringString converts a pointer to a map of strings into a map of interface{}.
//
// Deprecated: Use FlattenSlice or custom flatten map logic instead.
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

// FlattenInt32Slice converts a pointer to a slice of int32s into a slice of interface{}.
//
// Deprecated: Use FlattenSlice instead.
func FlattenInt32Slice(input *[]int32) []interface{} {
	return FlattenSlice(input)
}

// FlattenInt64Slice converts a pointer to a slice of int64s into a slice of interface{}.
//
// Deprecated: Use FlattenSlice instead.
func FlattenInt64Slice(input *[]int64) []interface{} {
	return FlattenSlice(input)
}

// ExpandStringSliceWithDelimiter converts an interface slice to a delimited string pointer.
//
// Deprecated: Use ExpandSliceWithDelimiter instead.
func ExpandStringSliceWithDelimiter(input []interface{}, delimiter string) *string {
	return ExpandSliceWithDelimiter(input, func(i string) string { return i }, delimiter)
}

// ExpandIntSliceWithDelimiter converts an interface slice of ints to a delimited string pointer.
//
// Deprecated: Use ExpandSliceWithDelimiter instead.
func ExpandIntSliceWithDelimiter(input []interface{}, delimiter string) *string {
	return ExpandSliceWithDelimiter(input, strconv.Itoa, delimiter)
}

// FlattenStringSliceWithDelimiter splits a delimited string pointer into a slice of interface{}.
//
// Deprecated: Custom string splitting should be done inline.
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
