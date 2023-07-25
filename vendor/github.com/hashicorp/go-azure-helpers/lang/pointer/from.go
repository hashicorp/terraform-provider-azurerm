// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pointer

// FromBool turns a boolean into a pointer to a boolean
func FromBool(input bool) *bool {
	return &input
}

// FromFloat64 turns a float64 into a pointer to a float64
func FromFloat64(input float64) *float64 {
	return &input
}

// FromInt turns a int into a pointer to a int
func FromInt(input int) *int {
	return &input
}

// FromInt64 turns a int64 into a pointer to a int64
func FromInt64(input int64) *int64 {
	return &input
}

// FromMapOfStringInterfaces turns a map[string]interface{} into a pointer to a map[string]interface{}
func FromMapOfStringInterfaces(input map[string]interface{}) *map[string]interface{} {
	return &input
}

// FromMapOfStringStrings turns a map[string]string into a pointer to a map[string]string
func FromMapOfStringStrings(input map[string]string) *map[string]string {
	return &input
}

// FromSliceOfStrings turns a slice of stirngs into a pointer to a slice of strings
func FromSliceOfStrings(input []string) *[]string {
	return &input
}

// FromString turns a string into a pointer to a string
func FromString(input string) *string {
	return &input
}
