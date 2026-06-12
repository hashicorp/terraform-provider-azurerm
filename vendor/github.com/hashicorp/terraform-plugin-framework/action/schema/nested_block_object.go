// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ fwxschema.NestedBlockObjectWithValidators = NestedBlockObject{}

// NestedBlockObject is the object containing the underlying attributes and
// blocks for a ListNestedBlock or SetNestedBlock. When retrieving the value
// for this attribute, use types.Object as the value type unless the CustomType
// field is set.
//
// This object enables customizing and simplifying details within its parent
// Block, therefore it cannot have Terraform schema fields such as Description,
// etc.
type NestedBlockObject struct {
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

// ApplyTerraform5AttributePathStep performs an AttributeName step on the
// underlying attributes or returns an error.
func (o NestedBlockObject) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (any, error) {
	return fwschema.NestedBlockObjectApplyTerraform5AttributePathStep(o, step)
}

// Equal returns true if the given NestedBlockObject is equivalent.
func (o NestedBlockObject) Equal(other fwschema.NestedBlockObject) bool {
	if _, ok := other.(NestedBlockObject); !ok {
		return false
	}

	return fwschema.NestedBlockObjectEqual(o, other)
}

// GetAttributes returns the Attributes field value.
func (o NestedBlockObject) GetAttributes() fwschema.UnderlyingAttributes {
	return schemaAttributes(o.Attributes)
}

// GetAttributes returns the Blocks field value.
func (o NestedBlockObject) GetBlocks() map[string]fwschema.Block {
	return schemaBlocks(o.Blocks)
}

// ObjectValidators returns the Validators field value.
func (o NestedBlockObject) ObjectValidators() []validator.Object {
	return o.Validators
}

// Type returns the framework type of the NestedBlockObject.
func (o NestedBlockObject) Type() basetypes.ObjectTypable {
	if o.CustomType != nil {
		return o.CustomType
	}

	return fwschema.NestedBlockObjectType(o)
}
