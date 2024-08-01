// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Parameter = StringParameter{}
var _ ParameterWithStringValidators = StringParameter{}

// StringParameter represents a function parameter that is a string.
//
// When retrieving the argument value for this parameter:
//
//   - If CustomType is set, use its associated value type.
//   - If AllowUnknownValues is enabled, you must use the [types.String] value
//     type.
//   - If AllowNullValue is enabled, you must use [types.String] or *string
//     value types.
//   - Otherwise, use [types.String] or *string, or string value types.
//
// Terraform configurations set this parameter's argument data using expressions
// that return a string or directly via double quote ("value") syntax.
type StringParameter struct {
	// AllowNullValue when enabled denotes that a null argument value can be
	// passed to the function. When disabled, Terraform returns an error if the
	// argument value is null.
	AllowNullValue bool

	// AllowUnknownValues when enabled denotes that an unknown argument value
	// can be passed to the function. When disabled, Terraform skips the
	// function call entirely and assumes an unknown value result from the
	// function.
	AllowUnknownValues bool

	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.StringType]. When retrieving data, the
	// [basetypes.StringValuable] implementation associated with this custom
	// type must be used in place of [types.String].
	CustomType basetypes.StringTypable

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

	// Validators is a list of string validators that should be applied to the
	// parameter.
	Validators []StringParameterValidator
}

// GetValidators returns the string validators for the parameter.
func (p StringParameter) GetValidators() []StringParameterValidator {
	return p.Validators
}

// GetAllowNullValue returns if the parameter accepts a null value.
func (p StringParameter) GetAllowNullValue() bool {
	return p.AllowNullValue
}

// GetAllowUnknownValues returns if the parameter accepts an unknown value.
func (p StringParameter) GetAllowUnknownValues() bool {
	return p.AllowUnknownValues
}

// GetDescription returns the parameter plaintext description.
func (p StringParameter) GetDescription() string {
	return p.Description
}

// GetMarkdownDescription returns the parameter Markdown description.
func (p StringParameter) GetMarkdownDescription() string {
	return p.MarkdownDescription
}

// GetName returns the parameter name.
func (p StringParameter) GetName() string {
	return p.Name
}

// GetType returns the parameter data type.
func (p StringParameter) GetType() attr.Type {
	if p.CustomType != nil {
		return p.CustomType
	}

	return basetypes.StringType{}
}
