// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identityschema

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the desired interfaces.
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
type Float64Attribute struct {
	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.Float64Type. When retrieving data, the basetypes.Float64Valuable
	// associated with this custom type must be used in place of types.Float64.
	CustomType basetypes.Float64Typable

	// RequiredForImport indicates whether the practitioner must enter a value for
	// this attribute when importing a managed resource by this identity.
	// RequiredForImport and OptionalForImport cannot both be true.
	RequiredForImport bool

	// OptionalForImport indicates whether the practitioner can choose to enter a value
	// for this attribute when importing a managed resource by this identity.
	// OptionalForImport and RequiredForImport cannot both be true.
	OptionalForImport bool

	// Description is used in various tooling, like the language server or the documentation
	// generator, to give practitioners more information about what this attribute is,
	// what it's for, and how it should be used. It can be written as plain text with no
	// special formatting, or formatted as Markdown.
	Description string
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

// GetDeprecationMessage returns an empty string as identity attributes cannot
// surface deprecation messages.
func (a Float64Attribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value. For identity attributes,
// there is only a single description field that is permitted to contain plaintext or Markdown.
func (a Float64Attribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the Description field value. For identity attributes,
// there is only a single description field that is permitted to contain Markdown or plaintext.
func (a Float64Attribute) GetMarkdownDescription() string {
	return a.Description
}

// GetType returns types.Float64Type or the CustomType field value if defined.
func (a Float64Attribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.Float64Type
}

// IsComputed returns false as it's not relevant for identity schemas.
func (a Float64Attribute) IsComputed() bool {
	return false
}

// IsOptional returns false as it's not relevant for identity schemas.
func (a Float64Attribute) IsOptional() bool {
	return false
}

// IsRequired returns false as it's not relevant for identity schemas.
func (a Float64Attribute) IsRequired() bool {
	return false
}

// IsSensitive returns false as it's not relevant for identity schemas.
func (a Float64Attribute) IsSensitive() bool {
	return false
}

// IsWriteOnly returns false as it's not relevant for identity schemas.
func (a Float64Attribute) IsWriteOnly() bool {
	return false
}

// IsRequiredForImport returns the RequiredForImport field value.
func (a Float64Attribute) IsRequiredForImport() bool {
	return a.RequiredForImport
}

// IsOptionalForImport returns the OptionalForImport field value.
func (a Float64Attribute) IsOptionalForImport() bool {
	return a.OptionalForImport
}
