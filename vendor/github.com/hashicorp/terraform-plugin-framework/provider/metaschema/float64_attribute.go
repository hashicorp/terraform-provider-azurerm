// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metaschema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ Attribute = Float64Attribute{}
)

// Float64Attribute represents a schema attribute that is a 64-bit floating
// point number. When retrieving the value for this attribute, use
// types.Float64 as the value type unless the CustomType field is set.
//
// Use Int64Attribute for 64-bit integer attributes or NumberAttribute for
// 512-bit generic number attributes.
//
// Terraform configurations configure this attribute using expressions that
// return a number or directly via a floating point value.
//
//	example_attribute = 123.45
//
// Terraform configurations reference this attribute using the attribute name.
//
//	.example_attribute
type Float64Attribute struct {
	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.Float64Type. When retrieving data, the basetypes.Float64Valuable
	// associated with this custom type must be used in place of types.Float64.
	CustomType basetypes.Float64Typable

	// Required indicates whether the practitioner must enter a value for
	// this attribute or not. Required and Optional cannot both be true,
	// and Required and Computed cannot both be true.
	Required bool

	// Optional indicates whether the practitioner can choose to enter a value
	// for this attribute or not. Optional and Required cannot both be true.
	Optional bool

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
}

// ApplyTerraform5AttributePathStep always returns an error as it is not
// possible to step further into a Float64Attribute.
func (a Float64Attribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return a.GetType().ApplyTerraform5AttributePathStep(step)
}

// Equal returns true if the given Attribute is a Float64Attribute
// and all fields are equal.
func (a Float64Attribute) Equal(o fwschema.Attribute) bool {
	if _, ok := o.(Float64Attribute); !ok {
		return false
	}

	return fwschema.AttributesEqual(a, o)
}

// GetDeprecationMessage always returns an empty string as there is no
// deprecation validation support for provider meta schemas.
func (a Float64Attribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value.
func (a Float64Attribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a Float64Attribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetType returns types.Float64Type or the CustomType field value if defined.
func (a Float64Attribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.Float64Type
}

// IsComputed always returns false as provider schemas cannot be Computed.
func (a Float64Attribute) IsComputed() bool {
	return false
}

// IsOptional returns the Optional field value.
func (a Float64Attribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a Float64Attribute) IsRequired() bool {
	return a.Required
}

// IsSensitive always returns false as there is no plan for provider meta
// schema data.
func (a Float64Attribute) IsSensitive() bool {
	return false
}
