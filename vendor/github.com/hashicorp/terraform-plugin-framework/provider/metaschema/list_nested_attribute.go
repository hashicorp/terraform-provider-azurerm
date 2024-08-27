// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metaschema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ NestedAttribute = ListNestedAttribute{}
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

// GetDeprecationMessage always returns an empty string as there is no
// deprecation validation support for provider meta schemas.
func (a ListNestedAttribute) GetDeprecationMessage() string {
	return ""
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

// IsSensitive always returns false as there is no plan for provider meta
// schema data.
func (a ListNestedAttribute) IsSensitive() bool {
	return false
}
