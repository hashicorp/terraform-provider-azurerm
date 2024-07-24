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
	_ Attribute = StringAttribute{}
)

// StringAttribute represents a schema attribute that is a string. When
// retrieving the value for this attribute, use types.String as the value type
// unless the CustomType field is set.
//
// Terraform configurations configure this attribute using expressions that
// return a string or directly via double quote syntax.
//
//	example_attribute = "value"
//
// Terraform configurations reference this attribute using the attribute name.
//
//	.example_attribute
type StringAttribute struct {
	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.StringType. When retrieving data, the basetypes.StringValuable
	// associated with this custom type must be used in place of types.String.
	CustomType basetypes.StringTypable

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
// possible to step further into a StringAttribute.
func (a StringAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return a.GetType().ApplyTerraform5AttributePathStep(step)
}

// Equal returns true if the given Attribute is a StringAttribute
// and all fields are equal.
func (a StringAttribute) Equal(o fwschema.Attribute) bool {
	if _, ok := o.(StringAttribute); !ok {
		return false
	}

	return fwschema.AttributesEqual(a, o)
}

// GetDeprecationMessage always returns an empty string as there is no
// deprecation validation support for provider meta schemas.
func (a StringAttribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value.
func (a StringAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a StringAttribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetType returns types.StringType or the CustomType field value if defined.
func (a StringAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.StringType
}

// IsComputed always returns false as provider schemas cannot be Computed.
func (a StringAttribute) IsComputed() bool {
	return false
}

// IsOptional returns the Optional field value.
func (a StringAttribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a StringAttribute) IsRequired() bool {
	return a.Required
}

// IsSensitive always returns false as there is no plan for provider meta
// schema data.
func (a StringAttribute) IsSensitive() bool {
	return false
}
