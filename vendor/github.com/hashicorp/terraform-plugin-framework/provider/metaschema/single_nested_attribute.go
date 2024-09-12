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
	_ NestedAttribute = SingleNestedAttribute{}
)

// SingleNestedAttribute represents an attribute that is a single object where
// the object attributes can be fully defined, including further nested
// attributes. When retrieving the value for this attribute, use types.Object
// as the value type unless the CustomType field is set. The Attributes field
// must be set. Nested attributes are only compatible with protocol version 6.
//
// Use ObjectAttribute if the underlying attributes do not require definition
// beyond type information.
//
// Terraform configurations configure this attribute using expressions that
// return an object or directly via curly brace syntax.
//
//	# single object
//	example_attribute = {
//		nested_attribute = #...
//	}
//
// Terraform configurations reference this attribute using expressions that
// accept an object or an attribute name directly via period syntax:
//
//	# object nested_attribute value
//	.example_attribute.nested_attribute
type SingleNestedAttribute struct {
	// Attributes is the mapping of underlying attribute names to attribute
	// definitions. This field must be set.
	Attributes map[string]Attribute

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

// ApplyTerraform5AttributePathStep returns the Attributes field value if step
// is AttributeName, otherwise returns an error.
func (a SingleNestedAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	name, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("cannot apply step %T to SingleNestedAttribute", step)
	}

	attribute, ok := a.Attributes[string(name)]

	if !ok {
		return nil, fmt.Errorf("no attribute %q on SingleNestedAttribute", name)
	}

	return attribute, nil
}

// Equal returns true if the given Attribute is a SingleNestedAttribute
// and all fields are equal.
func (a SingleNestedAttribute) Equal(o fwschema.Attribute) bool {
	other, ok := o.(SingleNestedAttribute)

	if !ok {
		return false
	}

	return fwschema.NestedAttributesEqual(a, other)
}

// GetAttributes returns the Attributes field value.
func (a SingleNestedAttribute) GetAttributes() fwschema.UnderlyingAttributes {
	return schemaAttributes(a.Attributes)
}

// GetDeprecationMessage always returns an empty string as there is no
// deprecation validation support for provider meta schemas.
func (a SingleNestedAttribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value.
func (a SingleNestedAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a SingleNestedAttribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetNestedObject returns a generated NestedAttributeObject from the
// Attributes and CustomType field values.
func (a SingleNestedAttribute) GetNestedObject() fwschema.NestedAttributeObject {
	return NestedAttributeObject{
		Attributes: a.Attributes,
		CustomType: a.CustomType,
	}
}

// GetNestingMode always returns NestingModeList.
func (a SingleNestedAttribute) GetNestingMode() fwschema.NestingMode {
	return fwschema.NestingModeSingle
}

// GetType returns ListType of ObjectType or CustomType.
func (a SingleNestedAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	attrTypes := make(map[string]attr.Type, len(a.Attributes))

	for name, attribute := range a.Attributes {
		attrTypes[name] = attribute.GetType()
	}

	return types.ObjectType{
		AttrTypes: attrTypes,
	}
}

// IsComputed always returns false as provider schemas cannot be Computed.
func (a SingleNestedAttribute) IsComputed() bool {
	return false
}

// IsOptional returns the Optional field value.
func (a SingleNestedAttribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a SingleNestedAttribute) IsRequired() bool {
	return a.Required
}

// IsSensitive always returns false as there is no plan for provider meta
// schema data.
func (a SingleNestedAttribute) IsSensitive() bool {
	return false
}
