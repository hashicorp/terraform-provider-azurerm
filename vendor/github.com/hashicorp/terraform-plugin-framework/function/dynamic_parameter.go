// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Parameter = DynamicParameter{}
var _ ParameterWithDynamicValidators = DynamicParameter{}

// DynamicParameter represents a function parameter that is a dynamic, rather
// than a static type. Static types are always preferable over dynamic
// types in Terraform as practitioners will receive less helpful configuration
// assistance from validation error diagnostics and editor integrations.
//
// When retrieving the argument value for this parameter:
//
//   - If CustomType is set, use its associated value type.
//   - Otherwise, use the [types.Dynamic] value type.
//
// The concrete value type for a dynamic is determined at runtime by Terraform,
// if defined in the configuration.
type DynamicParameter struct {
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
	// default [basetypes.DynamicType]. When retrieving data, the
	// [basetypes.DynamicValuable] implementation associated with this custom
	// type must be used in place of [types.Dynamic].
	CustomType basetypes.DynamicTypable

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

	// Validators is a list of dynamic validators that should be applied to the
	// parameter.
	Validators []DynamicParameterValidator
}

// GetValidators returns the list of validators for the parameter.
func (p DynamicParameter) GetValidators() []DynamicParameterValidator {
	return p.Validators
}

// GetAllowNullValue returns if the parameter accepts a null value.
func (p DynamicParameter) GetAllowNullValue() bool {
	return p.AllowNullValue
}

// GetAllowUnknownValues returns if the parameter accepts an unknown value.
func (p DynamicParameter) GetAllowUnknownValues() bool {
	return p.AllowUnknownValues
}

// GetDescription returns the parameter plaintext description.
func (p DynamicParameter) GetDescription() string {
	return p.Description
}

// GetMarkdownDescription returns the parameter Markdown description.
func (p DynamicParameter) GetMarkdownDescription() string {
	return p.MarkdownDescription
}

// GetName returns the parameter name.
func (p DynamicParameter) GetName() string {
	return p.Name
}

// GetType returns the parameter data type.
func (p DynamicParameter) GetType() attr.Type {
	if p.CustomType != nil {
		return p.CustomType
	}

	return basetypes.DynamicType{}
}
