// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metaschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ fwschema.NestedAttributeObject = NestedAttributeObject{}

// NestedAttributeObject is the object containing the underlying attributes
// for a ListNestedAttribute, MapNestedAttribute, SetNestedAttribute, or
// SingleNestedAttribute (automatically generated). When retrieving the value
// for this attribute, use types.Object as the value type unless the CustomType
// field is set. The Attributes field must be set. Nested attributes are only
// compatible with protocol version 6.
//
// This object enables customizing and simplifying details within its parent
// NestedAttribute, therefore it cannot have Terraform schema fields such as
// Required, Description, etc.
type NestedAttributeObject struct {
	// Attributes is the mapping of underlying attribute names to attribute
	// definitions. This field must be set.
	Attributes map[string]Attribute

	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.ObjectType. When retrieving data, the basetypes.ObjectValuable
	// associated with this custom type must be used in place of types.Object.
	CustomType basetypes.ObjectTypable
}

// ApplyTerraform5AttributePathStep performs an AttributeName step on the
// underlying attributes or returns an error.
func (o NestedAttributeObject) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (any, error) {
	return fwschema.NestedAttributeObjectApplyTerraform5AttributePathStep(o, step)
}

// Equal returns true if the given NestedAttributeObject is equivalent.
func (o NestedAttributeObject) Equal(other fwschema.NestedAttributeObject) bool {
	if _, ok := other.(NestedAttributeObject); !ok {
		return false
	}

	return fwschema.NestedAttributeObjectEqual(o, other)
}

// GetAttributes returns the Attributes field value.
func (o NestedAttributeObject) GetAttributes() fwschema.UnderlyingAttributes {
	return schemaAttributes(o.Attributes)
}

// Type returns the framework type of the NestedAttributeObject.
func (o NestedAttributeObject) Type() basetypes.ObjectTypable {
	if o.CustomType != nil {
		return o.CustomType
	}

	return fwschema.NestedAttributeObjectType(o)
}
