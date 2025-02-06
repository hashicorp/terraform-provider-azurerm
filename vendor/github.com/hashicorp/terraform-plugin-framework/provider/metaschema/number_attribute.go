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
	_ Attribute = NumberAttribute{}
)

// NumberAttribute represents a schema attribute that is a generic number with
// up to 512 bits of floating point or integer precision. When retrieving the
// value for this attribute, use types.Number as the value type unless the
// CustomType field is set.
//
// Use Float64Attribute for 64-bit floating point number attributes or
// Int64Attribute for 64-bit integer number attributes.
//
// Terraform configurations configure this attribute using expressions that
// return a number or directly via a floating point or integer value.
//
//	example_attribute = 123
//
// Terraform configurations reference this attribute using the attribute name.
//
//	.example_attribute
type NumberAttribute struct {
	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.NumberType. When retrieving data, the basetypes.NumberValuable
	// associated with this custom type must be used in place of types.Number.
	CustomType basetypes.NumberTypable

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
// possible to step further into a NumberAttribute.
func (a NumberAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return a.GetType().ApplyTerraform5AttributePathStep(step)
}

// Equal returns true if the given Attribute is a NumberAttribute
// and all fields are equal.
func (a NumberAttribute) Equal(o fwschema.Attribute) bool {
	if _, ok := o.(NumberAttribute); !ok {
		return false
	}

	return fwschema.AttributesEqual(a, o)
}

// GetDeprecationMessage always returns an empty string as there is no
// deprecation validation support for provider meta schemas.
func (a NumberAttribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value.
func (a NumberAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a NumberAttribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetType returns types.NumberType or the CustomType field value if defined.
func (a NumberAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.NumberType
}

// IsComputed always returns false as provider schemas cannot be Computed.
func (a NumberAttribute) IsComputed() bool {
	return false
}

// IsOptional returns the Optional field value.
func (a NumberAttribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a NumberAttribute) IsRequired() bool {
	return a.Required
}

// IsSensitive always returns false as there is no plan for provider meta
// schema data.
func (a NumberAttribute) IsSensitive() bool {
	return false
}
