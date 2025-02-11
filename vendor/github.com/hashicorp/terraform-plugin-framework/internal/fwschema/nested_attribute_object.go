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

// NestedAttributeObject represents the Object inside a NestedAttribute.
// Refer to the fwxschema package for validation and plan modification
// extensions to this interface.
type NestedAttributeObject interface {
	tftypes.AttributePathStepper

	// Equal should return true if given NestedAttributeObject is equivalent.
	Equal(NestedAttributeObject) bool

	// GetAttributes should return the nested attributes of an attribute.
	GetAttributes() UnderlyingAttributes

	// Type should return the framework type of the object.
	Type() basetypes.ObjectTypable
}

// NestedAttributeObjectApplyTerraform5AttributePathStep is a helper function
// to perform base tftypes.AttributePathStepper handling using the
// GetAttributes method. NestedAttributeObject implementations should still
// include custom type functionality in addition to using this helper.
func NestedAttributeObjectApplyTerraform5AttributePathStep(o NestedAttributeObject, step tftypes.AttributePathStep) (any, error) {
	name, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("cannot apply AttributePathStep %T to NestedAttributeObject", step)
	}

	attribute, ok := o.GetAttributes()[string(name)]

	if ok {
		return attribute, nil
	}

	return nil, fmt.Errorf("no attribute %q on NestedAttributeObject", name)
}

// NestedAttributeObjectEqual is a helper function to perform base equality testing
// on two NestedAttributeObject. NestedAttributeObject implementations should still
// compare the concrete types and other custom functionality in addition to
// using this helper.
func NestedAttributeObjectEqual(a, b NestedAttributeObject) bool {
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

	return true
}

// NestedAttributeObjectType is a helper function to perform base type handling
// using the GetAttributes and GetBlocks methods. NestedAttributeObject
// implementations should still include custom type functionality in addition
// to using this helper.
func NestedAttributeObjectType(o NestedAttributeObject) basetypes.ObjectTypable {
	attrTypes := make(map[string]attr.Type, len(o.GetAttributes()))

	for name, attribute := range o.GetAttributes() {
		attrTypes[name] = attribute.GetType()
	}

	return types.ObjectType{
		AttrTypes: attrTypes,
	}
}
