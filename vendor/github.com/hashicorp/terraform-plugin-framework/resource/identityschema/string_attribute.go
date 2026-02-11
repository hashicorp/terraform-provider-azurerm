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
type StringAttribute struct {
	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.StringType. When retrieving data, the basetypes.StringValuable
	// associated with this custom type must be used in place of types.String.
	CustomType basetypes.StringTypable

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

// GetDeprecationMessage returns an empty string as identity attributes cannot
// surface deprecation messages.
func (a StringAttribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value. For identity attributes,
// there is only a single description field that is permitted to contain plaintext or Markdown.
func (a StringAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the Description field value. For identity attributes,
// there is only a single description field that is permitted to contain Markdown or plaintext.
func (a StringAttribute) GetMarkdownDescription() string {
	return a.Description
}

// GetType returns types.StringType or the CustomType field value if defined.
func (a StringAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.StringType
}

// IsComputed returns false as it's not relevant for identity schemas.
func (a StringAttribute) IsComputed() bool {
	return false
}

// IsOptional returns false as it's not relevant for identity schemas.
func (a StringAttribute) IsOptional() bool {
	return false
}

// IsRequired returns false as it's not relevant for identity schemas.
func (a StringAttribute) IsRequired() bool {
	return false
}

// IsSensitive returns false as it's not relevant for identity schemas.
func (a StringAttribute) IsSensitive() bool {
	return false
}

// IsWriteOnly returns false as it's not relevant for identity schemas.
func (a StringAttribute) IsWriteOnly() bool {
	return false
}

// IsRequiredForImport returns the RequiredForImport field value.
func (a StringAttribute) IsRequiredForImport() bool {
	return a.RequiredForImport
}

// IsOptionalForImport returns the OptionalForImport field value.
func (a StringAttribute) IsOptionalForImport() bool {
	return a.OptionalForImport
}
