// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ Block                               = SingleNestedBlock{}
	_ fwxschema.BlockWithObjectValidators = SingleNestedBlock{}
)

// SingleNestedBlock represents a block that is a single object where
// the object attributes can be fully defined, including further attributes
// or blocks. When retrieving the value for this block, use types.Object
// as the value type unless the CustomType field is set.
//
// Prefer SingleNestedAttribute over SingleNestedBlock if the provider is
// using protocol version 6. Nested attributes allow practitioners to configure
// values directly with expressions.
//
// Terraform configurations configure this block only once using curly brace
// syntax without an equals (=) sign or [Dynamic Block Expressions].
//
//	# single block
//	example_block {
//		nested_attribute = #...
//	}
//
// Terraform configurations reference this block using expressions that
// accept an object or an attribute name directly via period syntax:
//
//	# object nested_attribute value
//	.example_block.nested_attribute
//
// [Dynamic Block Expressions]: https://developer.hashicorp.com/terraform/language/expressions/dynamic-blocks
type SingleNestedBlock struct {
	// Attributes is the mapping of underlying attribute names to attribute
	// definitions.
	//
	// Names must only contain lowercase letters, numbers, and underscores.
	// Names must not collide with any Blocks names.
	Attributes map[string]Attribute

	// Blocks is the mapping of underlying block names to block definitions.
	//
	// Names must only contain lowercase letters, numbers, and underscores.
	// Names must not collide with any Attributes names.
	Blocks map[string]Block

	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.ObjectType. When retrieving data, the basetypes.ObjectValuable
	// associated with this custom type must be used in place of types.Object.
	CustomType basetypes.ObjectTypable

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
	Validators []validator.Object
}

// ApplyTerraform5AttributePathStep returns the Attributes field value if step
// is AttributeName, otherwise returns an error.
func (b SingleNestedBlock) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	name, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("cannot apply step %T to SingleNestedBlock", step)
	}

	if attribute, ok := b.Attributes[string(name)]; ok {
		return attribute, nil
	}

	if block, ok := b.Blocks[string(name)]; ok {
		return block, nil
	}

	return nil, fmt.Errorf("no attribute or block %q on SingleNestedBlock", name)
}

// Equal returns true if the given Attribute is b SingleNestedBlock
// and all fields are equal.
func (b SingleNestedBlock) Equal(o fwschema.Block) bool {
	if _, ok := o.(SingleNestedBlock); !ok {
		return false
	}

	return fwschema.BlocksEqual(b, o)
}

// GetDeprecationMessage returns the DeprecationMessage field value.
func (b SingleNestedBlock) GetDeprecationMessage() string {
	return b.DeprecationMessage
}

// GetDescription returns the Description field value.
func (b SingleNestedBlock) GetDescription() string {
	return b.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (b SingleNestedBlock) GetMarkdownDescription() string {
	return b.MarkdownDescription
}

// GetNestedObject returns a generated NestedBlockObject from the
// Attributes, CustomType, and Validators field values.
func (b SingleNestedBlock) GetNestedObject() fwschema.NestedBlockObject {
	return NestedBlockObject{
		Attributes: b.Attributes,
		Blocks:     b.Blocks,
		CustomType: b.CustomType,
		Validators: b.Validators,
	}
}

// GetNestingMode always returns BlockNestingModeSingle.
func (b SingleNestedBlock) GetNestingMode() fwschema.BlockNestingMode {
	return fwschema.BlockNestingModeSingle
}

// ObjectValidators returns the Validators field value.
func (b SingleNestedBlock) ObjectValidators() []validator.Object {
	return b.Validators
}

// Type returns ObjectType or CustomType.
func (b SingleNestedBlock) Type() attr.Type {
	if b.CustomType != nil {
		return b.CustomType
	}

	attrTypes := make(map[string]attr.Type, len(b.Attributes)+len(b.Blocks))

	for name, attribute := range b.Attributes {
		attrTypes[name] = attribute.GetType()
	}

	for name, block := range b.Blocks {
		attrTypes[name] = block.Type()
	}

	return types.ObjectType{
		AttrTypes: attrTypes,
	}
}
