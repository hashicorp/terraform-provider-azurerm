// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package pointer

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

// FromEnumFW is a helper function to return a string from a pointer to an Enum without having to cast it
// example code simplification:
// myStruct.SomeStringValue = string(pointer.From(model.EnumValue))
// becomes
// myStruct.SomeStringValue = pointer.FromEnumFW(model.EnumValue)
// if input is nil, returns a Null String
func FromEnumFW[T ~string](input *T) (output types.String) {
	if input == nil {
		return types.StringNull()
	}

	return types.StringValue(string(*input))
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

// ToEnumFW is a helper function to cast strings as an Enum type where API objects expect a pointer to the Enum value
// example code simplification:
// APIModel.SomeValue = pointer.To(someservice.SomeEnumType(model.SomeVariable.ValueString))
// becomes
// APIModel.SomeValue = pointer.FWToEnum[someservice.SomeEnumType](model.SomeVariable)
// This function returns nil if the provided `types.String` value is Null or Unknown.
func ToEnumFW[T ~string](input types.String) *T {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	return To(T(input.ValueString()))
}