// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// NestedBlockObject represents the Object inside a Block.
// Refer to the fwxschema package for validation and plan modification
// extensions to this interface.
type NestedBlockObject interface {
	tftypes.AttributePathStepper

	// Equal should return true if given NestedBlockObject is equivalent.
	Equal(NestedBlockObject) bool

	// GetAttributes should return the nested attributes of the object.
	GetAttributes() UnderlyingAttributes

	// GetBlocks should return the nested attributes of the object.
	GetBlocks() map[string]Block

	// Type should return the framework type of the object.
	Type() basetypes.ObjectTypable
}

// NestedBlockObjectApplyTerraform5AttributePathStep is a helper function to
// perform base tftypes.AttributePathStepper handling using the GetAttributes
// and GetBlocks methods. NestedBlockObject implementations should still
// include custom type functionality in addition to using this helper.
func NestedBlockObjectApplyTerraform5AttributePathStep(o NestedBlockObject, step tftypes.AttributePathStep) (any, error) {
	name, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("cannot apply AttributePathStep %T to NestedBlockObject", step)
	}

	attribute, ok := o.GetAttributes()[string(name)]

	if ok {
		return attribute, nil
	}

	block, ok := o.GetBlocks()[string(name)]

	if ok {
		return block, nil
	}

	return nil, fmt.Errorf("no attribute or block %q on NestedBlockObject", name)
}

// NestedBlockObjectEqual is a helper function to perform base equality testing
// on two NestedBlockObject. NestedBlockObject implementations should still
// compare the concrete types and other custom functionality in addition to
// using this helper.
func NestedBlockObjectEqual(a, b NestedBlockObject) bool {
	if !a.Type().Equal(b.Type()) {
		return false
	}

	if len(a.GetAttributes()) != len(b.GetAttributes()) {
		return false
	}

	for name, aAttribute := range a.GetAttributes() {
		bAttribute, ok := b.GetAttributes()[name]

		if !ok {
			return false
		}

		if !aAttribute.Equal(bAttribute) {
			return false
		}
	}

	if len(a.GetBlocks()) != len(b.GetBlocks()) {
		return false
	}

	for name, aBlock := range a.GetBlocks() {
		bBlock, ok := b.GetBlocks()[name]

		if !ok {
			return false
		}

		if !aBlock.Equal(bBlock) {
			return false
		}
	}

	return true
}

// NestedBlockObjectType is a helper function to perform base type handling
// using the GetAttributes and GetBlocks methods. NestedBlockObject
// implementations should still include custom type functionality in addition
// to using this helper.
func NestedBlockObjectType(o NestedBlockObject) basetypes.ObjectTypable {
	attrTypes := make(map[string]attr.Type, len(o.GetAttributes())+len(o.GetBlocks()))

	for name, attribute := range o.GetAttributes() {
		attrTypes[name] = attribute.GetType()
	}

	for name, block := range o.GetBlocks() {
		attrTypes[name] = block.Type()
	}

	return types.ObjectType{
		AttrTypes: attrTypes,
	}
}
