// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwtype"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ NestedAttribute                              = ListNestedAttribute{}
	_ fwschema.AttributeWithValidateImplementation = ListNestedAttribute{}
	_ fwxschema.AttributeWithListValidators        = ListNestedAttribute{}
)

// ListNestedAttribute represents an attribute that is a list of objects where
// the object attributes can be fully defined, including further nested
// attributes. When retrieving the value for this attribute, use types.List
// as the value type unless the CustomType field is set. The NestedObject field
// must be set. Nested attributes are only compatible with protocol version 6.
//
// Use ListAttribute if the underlying elements are of a single type and do
// not require definition beyond type information.
//
// Terraform configurations configure this attribute using expressions that
// return a list of objects or directly via square and curly brace syntax.
//
//	# list of objects
//	example_attribute = [
//		{
//			nested_attribute = #...
//		},
//	]
//
// Terraform configurations reference this attribute using expressions that
// accept a list of objects or an element directly via square brace 0-based
// index syntax:
//
//	# first known object
//	.example_attribute[0]
//	# first known object nested_attribute value
//	.example_attribute[0].nested_attribute
type ListNestedAttribute struct {
	// NestedObject is the underlying object that contains nested attributes.
	// This field must be set.
	//
	// Nested attributes that contain a dynamic type (i.e. DynamicAttribute) are not supported.
	// If underlying dynamic values are required, replace this attribute definition with
	// DynamicAttribute instead.
	NestedObject NestedAttributeObject

	// CustomType enables the use of a custom attribute type in place of the
	// default types.ListType of types.ObjectType. When retrieving data, the
	// basetypes.ListValuable associated with this custom type must be used in
	// place of types.List.
	CustomType basetypes.ListTypable

	// Required indicates whether the practitioner must enter a value for
	// this attribute or not. Required and Optional cannot both be true,
	// and Required and Computed cannot both be true.
	Required bool

	// Optional indicates whether the practitioner can choose to enter a value
	// for this attribute or not. Optional and Required cannot both be true.
	Optional bool

	// Sensitive indicates whether the value of this attribute should be
	// considered sensitive data. Setting it to true will obscure the value
	// in CLI output. Sensitive does not impact how values are stored, and
	// practitioners are encouraged to store their state as if the entire
	// file is sensitive.
	Sensitive bool

	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this attribute is,
	// what it's for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this attribute is, what it's for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription string

	// DeprecationMessage defines warning diagnostic details to display when
	// practitioner configurations use this Attribute. The warning diagnostic
	// summary is automatically set to "Attribute Deprecated" along with
	// configuration source file and line information.
	//
	// Set this field to a practitioner actionable message such as:
	//
	//  - "Configure other_attribute instead. This attribute will be removed
	//    in the next major version of the provider."
	//  - "Remove this attribute's configuration as it no longer is used and
	//    the attribute will be removed in the next major version of the
	//    provider."
	//
	// In Terraform 1.2.7 and later, this warning diagnostic is displayed any
	// time a practitioner attempts to configure a value for this attribute and
	// certain scenarios where this attribute is referenced.
	//
	// In Terraform 1.2.6 and earlier, this warning diagnostic is only
	// displayed when the Attribute is Required or Optional, and if the
	// practitioner configuration sets the value to a known or unknown value
	// (which may eventually be null). It has no effect when the Attribute is
	// Computed-only (read-only; not Required or Optional).
	//
	// Across any Terraform version, there are no warnings raised for
	// practitioner configuration values set directly to null, as there is no
	// way for the framework to differentiate between an unset and null
	// configuration due to how Terraform sends configuration information
	// across the protocol.
	//
	// Additional information about deprecation enhancements for read-only
	// attributes can be found in:
	//
	//  - https://github.com/hashicorp/terraform/issues/7569
	//
	DeprecationMessage string

	// Validators define value validation functionality for the attribute. All
	// elements of the slice of AttributeValidator are run, regardless of any
	// previous error diagnostics.
	//
	// Many common use case validators can be found in the
	// github.com/hashicorp/terraform-plugin-framework-validators Go module.
	//
	// If the Type field points to a custom type that implements the
	// xattr.TypeWithValidate interface, the validators defined in this field
	// are run in addition to the validation defined by the type.
	Validators []validator.List
}

// ApplyTerraform5AttributePathStep returns the Attributes field value if step
// is ElementKeyInt, otherwise returns an error.
func (a ListNestedAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	_, ok := step.(tftypes.ElementKeyInt)

	if !ok {
		return nil, fmt.Errorf("cannot apply step %T to ListNestedAttribute", step)
	}

	return a.NestedObject, nil
}

// Equal returns true if the given Attribute is a ListNestedAttribute
// and all fields are equal.
func (a ListNestedAttribute) Equal(o fwschema.Attribute) bool {
	other, ok := o.(ListNestedAttribute)

	if !ok {
		return false
	}

	return fwschema.NestedAttributesEqual(a, other)
}

// GetDeprecationMessage returns the DeprecationMessage field value.
func (a ListNestedAttribute) GetDeprecationMessage() string {
	return a.DeprecationMessage
}

// GetDescription returns the Description field value.
func (a ListNestedAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a ListNestedAttribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetNestedObject returns the NestedObject field value.
func (a ListNestedAttribute) GetNestedObject() fwschema.NestedAttributeObject {
	return a.NestedObject
}

// GetNestingMode always returns NestingModeList.
func (a ListNestedAttribute) GetNestingMode() fwschema.NestingMode {
	return fwschema.NestingModeList
}

// GetType returns ListType of ObjectType or CustomType.
func (a ListNestedAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.ListType{
		ElemType: a.NestedObject.Type(),
	}
}

// IsComputed always returns false as provider schemas cannot be Computed.
func (a ListNestedAttribute) IsComputed() bool {
	return false
}

// IsOptional returns the Optional field value.
func (a ListNestedAttribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a ListNestedAttribute) IsRequired() bool {
	return a.Required
}

// IsSensitive returns the Sensitive field value.
func (a ListNestedAttribute) IsSensitive() bool {
	return a.Sensitive
}

// ListValidators returns the Validators field value.
func (a ListNestedAttribute) ListValidators() []validator.List {
	return a.Validators
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the attribute to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC and
// should never include false positives.
func (a ListNestedAttribute) ValidateImplementation(ctx context.Context, req fwschema.ValidateImplementationRequest, resp *fwschema.ValidateImplementationResponse) {
	if a.CustomType == nil && fwtype.ContainsCollectionWithDynamic(a.GetType()) {
		resp.Diagnostics.Append(fwtype.AttributeCollectionWithDynamicTypeDiag(req.Path))
	}
}
