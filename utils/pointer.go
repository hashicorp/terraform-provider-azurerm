// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

func Bool(input bool) *bool {
	return &input
}

func Int(input int) *int {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func Int64(input int64) *int64 {
	return &input
}

func Float(input float64) *float64 {
	return &input
}

func String(input string) *string {
	return &input
}

func ToPtr[E any](e E) *E {
	return &e
}

func StringSlice(input []string) *[]string {
	if input == nil {
		return nil
	}
	return &input
}
