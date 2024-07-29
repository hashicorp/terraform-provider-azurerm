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
	_ Attribute                                    = ListAttribute{}
	_ fwschema.AttributeWithValidateImplementation = ListAttribute{}
)

// ListAttribute represents a schema attribute that is a list with a single
// element type. When retrieving the value for this attribute, use types.List
// as the value type unless the CustomType field is set. The ElementType field
// must be set.
//
// Use ListNestedAttribute if the underlying elements should be objects and
// require definition beyond type information.
//
// Terraform configurations configure this attribute using expressions that
// return a list or directly via square brace syntax.
//
//	# list of strings
//	example_attribute = ["first", "second"]
//
// Terraform configurations reference this attribute using expressions that
// accept a list or an element directly via square brace 0-based index syntax:
//
//	# first known element
//	.example_attribute[0]
type ListAttribute struct {
	// ElementType is the type for all elements of the list. This field must be
	// set.
	ElementType attr.Type

	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.ListType. When retrieving data, the basetypes.ListValuable
	// associated with this custom type must be used in place of types.List.
	CustomType basetypes.ListTypable

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

// ApplyTerraform5AttributePathStep returns the result of stepping into a list
// index or an error.
func (a ListAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return a.GetType().ApplyTerraform5AttributePathStep(step)
}

// Equal returns true if the given Attribute is a ListAttribute
// and all fields are equal.
func (a ListAttribute) Equal(o fwschema.Attribute) bool {
	if _, ok := o.(ListAttribute); !ok {
		return false
	}

	return fwschema.AttributesEqual(a, o)
}

// GetDeprecationMessage always returns an empty string as there is no
// deprecation validation support for provider meta schemas.
func (a ListAttribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value.
func (a ListAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a ListAttribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetType returns types.ListType or the CustomType field value if defined.
func (a ListAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.ListType{
		ElemType: a.ElementType,
	}
}

// IsComputed always returns false as provider schemas cannot be Computed.
func (a ListAttribute) IsComputed() bool {
	return false
}

// IsOptional returns the Optional field value.
func (a ListAttribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a ListAttribute) IsRequired() bool {
	return a.Required
}

// IsSensitive always returns false as there is no plan for provider meta
// schema data.
func (a ListAttribute) IsSensitive() bool {
	return false
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the attribute to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC
// and should never include false positives.
func (a ListAttribute) ValidateImplementation(ctx context.Context, req fwschema.ValidateImplementationRequest, resp *fwschema.ValidateImplementationResponse) {
	if a.CustomType == nil && a.ElementType == nil {
		resp.Diagnostics.Append(fwschema.AttributeMissingElementTypeDiag(req.Path))
	}
}
