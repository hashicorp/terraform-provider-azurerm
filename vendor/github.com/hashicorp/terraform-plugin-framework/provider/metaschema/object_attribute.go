// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metaschema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ Attribute                                    = ObjectAttribute{}
	_ fwschema.AttributeWithValidateImplementation = ObjectAttribute{}
)

// ObjectAttribute represents a schema attribute that is an object with only
// type information for underlying attributes. When retrieving the value for
// this attribute, use types.Object as the value type unless the CustomType
// field is set. The AttributeTypes field must be set.
//
// Prefer SingleNestedAttribute over ObjectAttribute if the provider is
// using protocol version 6 and full attribute functionality is needed.
//
// Terraform configurations configure this attribute using expressions that
// return an object or directly via curly brace syntax.
//
//	# object with one attribute
//	example_attribute = {
//		underlying_attribute = #...
//	}
//
// Terraform configurations reference this attribute using expressions that
// accept an object or an attribute directly via period syntax:
//
//	# underlying attribute
//	.example_attribute.underlying_attribute
type ObjectAttribute struct {
	// AttributeTypes is the mapping of underlying attribute names to attribute
	// types. This field must be set.
	AttributeTypes map[string]attr.Type

	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.ObjectType. When retrieving data, the basetypes.ObjectValuable
	// associated with this custom type must be used in place of types.Object.
	CustomType basetypes.ObjectTypable

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

// ApplyTerraform5AttributePathStep returns the result of stepping into an
// attribute name or an error.
func (a ObjectAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return a.GetType().ApplyTerraform5AttributePathStep(step)
}

// Equal returns true if the given Attribute is a ObjectAttribute
// and all fields are equal.
func (a ObjectAttribute) Equal(o fwschema.Attribute) bool {
	if _, ok := o.(ObjectAttribute); !ok {
		return false
	}

	return fwschema.AttributesEqual(a, o)
}

// GetDeprecationMessage always returns an empty string as there is no
// deprecation validation support for provider meta schemas.
func (a ObjectAttribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value.
func (a ObjectAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a ObjectAttribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetType returns types.ObjectType or the CustomType field value if defined.
func (a ObjectAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.ObjectType{
		AttrTypes: a.AttributeTypes,
	}
}

// IsComputed always returns false as provider schemas cannot be Computed.
func (a ObjectAttribute) IsComputed() bool {
	return false
}

// IsOptional returns the Optional field value.
func (a ObjectAttribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a ObjectAttribute) IsRequired() bool {
	return a.Required
}

// IsSensitive always returns false as there is no plan for provider meta
// schema data.
func (a ObjectAttribute) IsSensitive() bool {
	return false
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the attribute to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC
// and should never include false positives.
func (a ObjectAttribute) ValidateImplementation(ctx context.Context, req fwschema.ValidateImplementationRequest, resp *fwschema.ValidateImplementationResponse) {
	if a.AttributeTypes == nil && a.CustomType == nil {
		resp.Diagnostics.Append(fwschema.AttributeMissingAttributeTypesDiag(req.Path))
	}
}
