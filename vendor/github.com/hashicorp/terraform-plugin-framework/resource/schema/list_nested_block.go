// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwtype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ Block                                    = ListNestedBlock{}
	_ fwschema.BlockWithValidateImplementation = ListNestedBlock{}
	_ fwxschema.BlockWithListPlanModifiers     = ListNestedBlock{}
	_ fwxschema.BlockWithListValidators        = ListNestedBlock{}
)

// ListNestedBlock represents a block that is a list of objects where
// the object attributes can be fully defined, including further attributes
// or blocks. When retrieving the value for this block, use types.List
// as the value type unless the CustomType field is set. The NestedObject field
// must be set.
//
// Prefer ListNestedAttribute over ListNestedBlock if the provider is
// using protocol version 6. Nested attributes allow practitioners to configure
// values directly with expressions.
//
// Terraform configurations configure this block repeatedly using curly brace
// syntax without an equals (=) sign or [Dynamic Block Expressions].
//
//	# list of blocks with two elements
//	example_block {
//		nested_attribute = #...
//	}
//	example_block {
//		nested_attribute = #...
//	}
//
// Terraform configurations reference this block using expressions that
// accept a list of objects or an element directly via square brace 0-based
// index syntax:
//
//	# first known object
//	.example_block[0]
//	# first known object nested_attribute value
//	.example_block[0].nested_attribute
//
// [Dynamic Block Expressions]: https://developer.hashicorp.com/terraform/language/expressions/dynamic-blocks
type ListNestedBlock struct {
	// NestedObject is the underlying object that contains nested attributes or
	// blocks. This field must be set.
	//
	// Nested attributes that contain a dynamic type (i.e. DynamicAttribute) are not supported.
	// If underlying dynamic values are required, replace this block definition with
	// a DynamicAttribute.
	NestedObject NestedBlockObject

	// CustomType enables the use of a custom attribute type in place of the
	// default types.ListType of types.ObjectType. When retrieving data, the
	// basetypes.ListValuable associated with this custom type must be used in
	// place of types.List.
	CustomType basetypes.ListTypable

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
	Validators []validator.List

	// PlanModifiers defines a sequence of modifiers for this attribute at
	// plan time. Schema-based plan modifications occur before any
	// resource-level plan modifications.
	//
	// Schema-based plan modifications can adjust Terraform's plan by:
	//
	//  - Requiring resource recreation. Typically used for configuration
	//    updates which cannot be done in-place.
	//  - Setting the planned value. Typically used for enhancing the plan
	//    to replace unknown values. Computed must be true or Terraform will
	//    return an error. If the plan value is known due to a known
	//    configuration value, the plan value cannot be changed or Terraform
	//    will return an error.
	//
	// Any errors will prevent further execution of this sequence or modifiers.
	PlanModifiers []planmodifier.List
}

// ApplyTerraform5AttributePathStep returns the NestedObject field value if step
// is ElementKeyInt, otherwise returns an error.
func (b ListNestedBlock) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	_, ok := step.(tftypes.ElementKeyInt)

	if !ok {
		return nil, fmt.Errorf("cannot apply step %T to ListNestedBlock", step)
	}

	return b.NestedObject, nil
}

// Equal returns true if the given Block is ListNestedBlock
// and all fields are equal.
func (b ListNestedBlock) Equal(o fwschema.Block) bool {
	if _, ok := o.(ListNestedBlock); !ok {
		return false
	}

	return fwschema.BlocksEqual(b, o)
}

// GetDeprecationMessage returns the DeprecationMessage field value.
func (b ListNestedBlock) GetDeprecationMessage() string {
	return b.DeprecationMessage
}

// GetDescription returns the Description field value.
func (b ListNestedBlock) GetDescription() string {
	return b.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (b ListNestedBlock) GetMarkdownDescription() string {
	return b.MarkdownDescription
}

// GetNestedObject returns the NestedObject field value.
func (b ListNestedBlock) GetNestedObject() fwschema.NestedBlockObject {
	return b.NestedObject
}

// GetNestingMode always returns BlockNestingModeList.
func (b ListNestedBlock) GetNestingMode() fwschema.BlockNestingMode {
	return fwschema.BlockNestingModeList
}

// ListPlanModifiers returns the PlanModifiers field value.
func (b ListNestedBlock) ListPlanModifiers() []planmodifier.List {
	return b.PlanModifiers
}

// ListValidators returns the Validators field value.
func (b ListNestedBlock) ListValidators() []validator.List {
	return b.Validators
}

// Type returns ListType of ObjectType or CustomType.
func (b ListNestedBlock) Type() attr.Type {
	if b.CustomType != nil {
		return b.CustomType
	}

	return types.ListType{
		ElemType: b.NestedObject.Type(),
	}
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the block to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC and
// should never include false positives.
func (b ListNestedBlock) ValidateImplementation(ctx context.Context, req fwschema.ValidateImplementationRequest, resp *fwschema.ValidateImplementationResponse) {
	if b.CustomType == nil && fwtype.ContainsCollectionWithDynamic(b.Type()) {
		resp.Diagnostics.Append(fwtype.BlockCollectionWithDynamicTypeDiag(req.Path))
	}
}
