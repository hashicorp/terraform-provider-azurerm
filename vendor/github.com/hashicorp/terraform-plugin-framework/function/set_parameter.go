// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwfunction"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwtype"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ Parameter                                      = SetParameter{}
	_ fwfunction.ParameterWithValidateImplementation = SetParameter{}
	_ ParameterWithSetValidators                     = SetParameter{}
)

// SetParameter represents a function parameter that is an unordered set of a
// single element type. Either the ElementType or CustomType field must be set.
//
// When retrieving the argument value for this parameter:
//
//   - If CustomType is set, use its associated value type.
//   - If AllowUnknownValues is enabled, you must use the [types.Set] value
//     type.
//   - Otherwise, use [types.Set] or any Go slice value types compatible with
//     the element type.
//
// Terraform configurations set this parameter's argument data using expressions
// that return a set or directly via set ("[...]") syntax.
type SetParameter struct {
	// ElementType is the type for all elements of the set. This field must be
	// set.
	//
	// Element types that contain a dynamic type (i.e. types.Dynamic) are not supported.
	// If underlying dynamic values are required, replace this parameter definition with
	// DynamicParameter instead.
	ElementType attr.Type

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
	// default [basetypes.SetType]. When retrieving data, the
	// [basetypes.SetValuable] implementation associated with this custom
	// type must be used in place of [types.Set].
	CustomType basetypes.SetTypable

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

	// Validators is a list of set validators that should be applied to the
	// parameter.
	Validators []SetParameterValidator
}

// GetValidators returns the list of validators for the parameter.
func (p SetParameter) GetValidators() []SetParameterValidator {
	return p.Validators
}

// GetAllowNullValue returns if the parameter accepts a null value.
func (p SetParameter) GetAllowNullValue() bool {
	return p.AllowNullValue
}

// GetAllowUnknownValues returns if the parameter accepts an unknown value.
func (p SetParameter) GetAllowUnknownValues() bool {
	return p.AllowUnknownValues
}

// GetDescription returns the parameter plaintext description.
func (p SetParameter) GetDescription() string {
	return p.Description
}

// GetMarkdownDescription returns the parameter Markdown description.
func (p SetParameter) GetMarkdownDescription() string {
	return p.MarkdownDescription
}

// GetName returns the parameter name.
func (p SetParameter) GetName() string {
	return p.Name
}

// GetType returns the parameter data type.
func (p SetParameter) GetType() attr.Type {
	if p.CustomType != nil {
		return p.CustomType
	}

	return basetypes.SetType{
		ElemType: p.ElementType,
	}
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the parameter to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC and
// should never include false positives.
func (p SetParameter) ValidateImplementation(ctx context.Context, req fwfunction.ValidateParameterImplementationRequest, resp *fwfunction.ValidateParameterImplementationResponse) {
	if p.CustomType == nil {
		if fwtype.ContainsCollectionWithDynamic(p.GetType()) {
			if req.ParameterPosition != nil {
				resp.Diagnostics.Append(fwtype.ParameterCollectionWithDynamicTypeDiag(*req.ParameterPosition, p.GetName()))
			} else {
				resp.Diagnostics.Append(fwtype.VariadicParameterCollectionWithDynamicTypeDiag(p.GetName()))
			}
		}

		if fwtype.ContainsMissingUnderlyingType(p.GetType()) {
			resp.Diagnostics.Append(fwtype.ParameterMissingUnderlyingTypeDiag(p.GetName(), req.ParameterPosition))
		}
	}

	if p.GetName() == "" {
		resp.Diagnostics.Append(fwfunction.MissingParameterNameDiag(req.FunctionName, req.ParameterPosition))
	}
}
