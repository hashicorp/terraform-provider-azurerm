// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Parameter = BoolParameter{}
var _ ParameterWithBoolValidators = BoolParameter{}

// BoolParameter represents a function parameter that is a boolean.
//
// When retrieving the argument value for this parameter:
//
//   - If CustomType is set, use its associated value type.
//   - If AllowUnknownValues is enabled, you must use the [types.Bool] value
//     type.
//   - If AllowNullValue is enabled, you must use [types.Bool] or *bool
//     value types.
//   - Otherwise, use [types.Bool] or *bool, or bool value types.
//
// Terraform configurations set this parameter's argument data using expressions
// that return a bool or directly via true/false keywords.
type BoolParameter struct {
	// AllowNullValue when enabled denotes that a null argument value can be
	// passed to the function. When disabled, Terraform returns an error if the
	// argument value is null.
	//
	// Enabling this requires reading argument values as *bool or [types.Bool].
	AllowNullValue bool

	// AllowUnknownValues when enabled denotes that an unknown argument value
	// can be passed to the function. When disabled, Terraform skips the
	// function call entirely and assumes an unknown value result from the
	// function.
	//
	// Enabling this requires reading argument values as [types.Bool].
	AllowUnknownValues bool

	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.BoolType]. When retrieving data, the
	// [basetypes.BoolValuable] implementation associated with this custom
	// type must be used in place of [types.Bool].
	CustomType basetypes.BoolTypable

	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this parameter is,
	// what it is for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this parameter is, what it is for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription string

	// Name is a short usage name for the parameter, such as "data". This name
	// is used in documentation, such as generating a function signature,
	// however its usage may be extended in the future.
	//
	// If no name is provided, this will default to "param" with a suffix of the
	// position the parameter is in the function definition. ("param1", "param2", etc.)
	// If the parameter is variadic, the default name will be "varparam".
	//
	// This must be a valid Terraform identifier, such as starting with an
	// alphabetical character and followed by alphanumeric or underscore
	// characters.
	Name string

	// Validators is a list of bool validators that should be applied to the
	// parameter.
	Validators []BoolParameterValidator
}

// GetValidators returns the list of validators for the parameter.
func (p BoolParameter) GetValidators() []BoolParameterValidator {
	return p.Validators
}

// GetAllowNullValue returns if the parameter accepts a null value.
func (p BoolParameter) GetAllowNullValue() bool {
	return p.AllowNullValue
}

// GetAllowUnknownValues returns if the parameter accepts an unknown value.
func (p BoolParameter) GetAllowUnknownValues() bool {
	return p.AllowUnknownValues
}

// GetDescription returns the parameter plaintext description.
func (p BoolParameter) GetDescription() string {
	return p.Description
}

// GetMarkdownDescription returns the parameter Markdown description.
func (p BoolParameter) GetMarkdownDescription() string {
	return p.MarkdownDescription
}

// GetName returns the parameter name.
func (p BoolParameter) GetName() string {
	return p.Name
}

// GetType returns the parameter data type.
func (p BoolParameter) GetType() attr.Type {
	if p.CustomType != nil {
		return p.CustomType
	}

	return basetypes.BoolType{}
}
