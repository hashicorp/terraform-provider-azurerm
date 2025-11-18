// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pointer

// From is a generic function that returns the value of a pointer
// If the pointer is nil, a zero value for the underlying type of the pointer is returned.
func From[T any](input *T) (output T) {
	var v T
	if input != nil {
		return *input
	}
	return v
}

// FromEnum is a helper function to return a string from a pointer to an Enum without having to cast it
// example code simplification:
// myStruct.SomeStringValue = string(pointer.From(model.EnumValue))
// becomes
// myStruct.SomeStringValue = pointer.FromEnum(model.EnumValue)
// if input is nil, returns an empty string
func FromEnum[T ~string](input *T) (output string) {
	if input == nil {
		return ""
	}

	return string(*input)
}

// To is a generic function that returns a pointer to the value provided.
func To[T any](input T) *T {
	return &input
}

// ToEnum is a helper function to cast strings as an Enum type where API objects expect a pointer to the Enum value
// example code simplification:
// APIModel.SomeValue = pointer.To(someservice.SomeEnumType(model.SomeVariable))
// becomes
// APIModel.SomeValue = pointer.ToEnum[someservice.SomeEnumType](model.SomeVariable)
func ToEnum[T ~string](input string) *T {
	result := T(input)
	return &result
}
