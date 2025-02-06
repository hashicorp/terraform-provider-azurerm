// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure UnderlyingAttributes satisfies the expected interfaces.
var _ tftypes.AttributePathStepper = UnderlyingAttributes{}

// UnderlyingAttributes represents attributes under a nested attribute.
type UnderlyingAttributes map[string]Attribute

// ApplyTerraform5AttributePathStep performs an AttributeName step on the
// underlying attributes or returns an error.
func (u UnderlyingAttributes) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (any, error) {
	name, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("can't apply %T to Attributes", step)
	}

	attribute, ok := u[string(name)]

	if !ok {
		return nil, fmt.Errorf("no attribute %q on Attributes", name)
	}

	return attribute, nil
}

// Equal returns true if all underlying attributes are equal.
func (u UnderlyingAttributes) Equal(o UnderlyingAttributes) bool {
	if len(u) != len(o) {
		return false
	}

	for name, uAttribute := range u {
		oAttribute, ok := o[name]

		if !ok {
			return false
		}

		if !uAttribute.Equal(oAttribute) {
			return false
		}
	}

	return true
}

// Type returns the framework type of the underlying attributes.
func (u UnderlyingAttributes) Type() basetypes.ObjectTypable {
	attrTypes := make(map[string]attr.Type, len(u))

	for name, attr := range u {
		attrTypes[name] = attr.GetType()
	}

	return basetypes.ObjectType{
		AttrTypes: attrTypes,
	}
}
