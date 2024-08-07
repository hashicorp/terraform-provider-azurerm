// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ ObjectTypable = ObjectType{}

// ObjectTypable extends attr.Type for object types.
// Implement this interface to create a custom ObjectType type.
type ObjectTypable interface {
	attr.Type

	// ValueFromObject should convert the Object to an ObjectValuable type.
	ValueFromObject(context.Context, ObjectValue) (ObjectValuable, diag.Diagnostics)
}

// ObjectType is an AttributeType representing an object.
type ObjectType struct {
	AttrTypes map[string]attr.Type
}

// WithAttributeTypes returns a new copy of the type with its attribute types
// set.
func (o ObjectType) WithAttributeTypes(typs map[string]attr.Type) attr.TypeWithAttributeTypes {
	return ObjectType{
		AttrTypes: typs,
	}
}

// AttributeTypes returns a copy of the type's attribute types.
func (o ObjectType) AttributeTypes() map[string]attr.Type {
	// Ensure callers cannot mutate the value
	result := make(map[string]attr.Type, len(o.AttrTypes))

	for key, value := range o.AttrTypes {
		result[key] = value
	}

	return result
}

// TerraformType returns the tftypes.Type that should be used to
// represent this type. This constrains what user input will be
// accepted and what kind of data can be set in state. The framework
// will use this to translate the AttributeType to something Terraform
// can understand.
func (o ObjectType) TerraformType(ctx context.Context) tftypes.Type {
	attributeTypes := map[string]tftypes.Type{}
	for k, v := range o.AttrTypes {
		attributeTypes[k] = v.TerraformType(ctx)
	}
	return tftypes.Object{
		AttributeTypes: attributeTypes,
	}
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value.
// This is meant to convert the tftypes.Value into a more convenient Go
// type for the provider to consume the data with.
func (o ObjectType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewObjectNull(o.AttrTypes), nil
	}
	if !in.Type().Equal(o.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", o.TerraformType(ctx), in.Type())
	}
	if !in.IsKnown() {
		return NewObjectUnknown(o.AttrTypes), nil
	}
	if in.IsNull() {
		return NewObjectNull(o.AttrTypes), nil
	}
	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}
	err := in.As(&val)
	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := o.AttrTypes[k].ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, err
		}
		attributes[k] = a
	}
	// ValueFromTerraform above on each attribute should make this safe.
	// Otherwise, this will need to do some Diagnostics to error conversion.
	return NewObjectValueMust(o.AttrTypes, attributes), nil
}

// Equal returns true if `candidate` is also an ObjectType and has the same
// AttributeTypes.
func (o ObjectType) Equal(candidate attr.Type) bool {
	other, ok := candidate.(ObjectType)
	if !ok {
		return false
	}
	if len(other.AttrTypes) != len(o.AttrTypes) {
		return false
	}
	for k, v := range o.AttrTypes {
		attr, ok := other.AttrTypes[k]
		if !ok {
			return false
		}
		if !v.Equal(attr) {
			return false
		}
	}
	return true
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// object.
func (o ObjectType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	attrName, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("cannot apply step %T to ObjectType", step)
	}

	attrType, ok := o.AttrTypes[string(attrName)]

	if !ok {
		return nil, fmt.Errorf("undefined attribute name %s in ObjectType", attrName)
	}

	return attrType, nil
}

// String returns a human-friendly description of the ObjectType.
func (o ObjectType) String() string {
	var res strings.Builder
	res.WriteString("types.ObjectType[")
	keys := make([]string, 0, len(o.AttrTypes))
	for k := range o.AttrTypes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for pos, key := range keys {
		if pos != 0 {
			res.WriteString(", ")
		}
		res.WriteString(`"` + key + `":`)
		res.WriteString(o.AttrTypes[key].String())
	}
	res.WriteString("]")
	return res.String()
}

// ValueType returns the Value type.
func (o ObjectType) ValueType(_ context.Context) attr.Value {
	return ObjectValue{
		attributeTypes: o.AttrTypes,
	}
}

// ValueFromObject returns an ObjectValuable type given an Object.
func (o ObjectType) ValueFromObject(_ context.Context, obj ObjectValue) (ObjectValuable, diag.Diagnostics) {
	return obj, nil
}
