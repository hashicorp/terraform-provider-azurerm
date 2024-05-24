// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var _ MapTypable = MapType{}

// MapTypable extends attr.Type for map types.
// Implement this interface to create a custom MapType type.
type MapTypable interface {
	attr.Type

	// ValueFromMap should convert the Map to a MapValuable type.
	ValueFromMap(context.Context, MapValue) (MapValuable, diag.Diagnostics)
}

// MapType is an AttributeType representing a map of values. All values must
// be of the same type, which the provider must specify as the ElemType
// property. Keys will always be strings.
type MapType struct {
	ElemType attr.Type
}

// WithElementType returns a new copy of the type with its element type set.
func (m MapType) WithElementType(typ attr.Type) attr.TypeWithElementType {
	return MapType{
		ElemType: typ,
	}
}

// ElementType returns the type's element type.
func (m MapType) ElementType() attr.Type {
	if m.ElemType == nil {
		return missingType{}
	}

	return m.ElemType
}

// TerraformType returns the tftypes.Type that should be used to represent this
// type. This constrains what user input will be accepted and what kind of data
// can be set in state. The framework will use this to translate the
// AttributeType to something Terraform can understand.
func (m MapType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.Map{
		ElementType: m.ElementType().TerraformType(ctx),
	}
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value. This is
// meant to convert the tftypes.Value into a more convenient Go type for the
// provider to consume the data with.
func (m MapType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewMapNull(m.ElementType()), nil
	}
	if !in.Type().Is(tftypes.Map{}) {
		return nil, fmt.Errorf("can't use %s as value of MapValue, can only use tftypes.Map values", in.String())
	}

	// MAINTAINER NOTE:
	// MapType does not support DynamicType as an element type. It is not explicitly prevented from being created with the
	// Framework type system, but the Framework-supported MapAttribute and MapNestedAttribute prevent DynamicType
	// from being used as an element type. An attempt to use DynamicType as the element type will eventually lead you to an error on this line :)
	//
	// In the future, if we ever need to support a map of dynamic element types, this type equality check will need to be modified to allow
	// dynamic types to not return an error, as the tftypes.Value coming in (if known) will be a concrete value, for example:
	//
	// - m.TerraformType(ctx): tftypes.Map[tftypes.DynamicPseudoType]
	// - in.Type(): tftypes.Map[tftypes.String]
	//
	// The `ValueFromTerraform` function for a dynamic type will be able create the correct concrete dynamic value with this modification in place.
	//
	if !in.Type().Equal(tftypes.Map{ElementType: m.ElementType().TerraformType(ctx)}) {
		return nil, fmt.Errorf("can't use %s as value of Map with ElementType %T, can only use %s values", in.String(), m.ElementType(), m.ElementType().TerraformType(ctx).String())
	}
	if !in.IsKnown() {
		return NewMapUnknown(m.ElementType()), nil
	}
	if in.IsNull() {
		return NewMapNull(m.ElementType()), nil
	}
	val := map[string]tftypes.Value{}
	err := in.As(&val)
	if err != nil {
		return nil, err
	}
	elems := make(map[string]attr.Value, len(val))
	for key, elem := range val {
		av, err := m.ElementType().ValueFromTerraform(ctx, elem)
		if err != nil {
			return nil, err
		}
		elems[key] = av
	}
	// ValueFromTerraform above on each element should make this safe.
	// Otherwise, this will need to do some Diagnostics to error conversion.
	return NewMapValueMust(m.ElementType(), elems), nil
}

// Equal returns true if `o` is also a MapType and has the same ElemType.
func (m MapType) Equal(o attr.Type) bool {
	// Preserve prior ElemType nil check behavior
	if m.ElementType().Equal(missingType{}) {
		return false
	}

	other, ok := o.(MapType)

	if !ok {
		return false
	}

	return m.ElementType().Equal(other.ElementType())
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// map.
func (m MapType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	if _, ok := step.(tftypes.ElementKeyString); !ok {
		return nil, fmt.Errorf("cannot apply step %T to MapType", step)
	}

	return m.ElementType(), nil
}

// String returns a human-friendly description of the MapType.
func (m MapType) String() string {
	return "types.MapType[" + m.ElementType().String() + "]"
}

// Validate validates all elements of the map that are of type
// xattr.TypeWithValidate.
func (m MapType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.Map{}) {
		err := fmt.Errorf("expected Map value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Map Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var elems map[string]tftypes.Value

	if err := in.As(&elems); err != nil {
		diags.AddAttributeError(
			path,
			"Map Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
	validatableType, isValidatable := m.ElementType().(xattr.TypeWithValidate)
	if !isValidatable {
		return diags
	}

	for index, elem := range elems {
		if !elem.IsFullyKnown() {
			continue
		}
		diags = append(diags, validatableType.Validate(ctx, elem, path.AtMapKey(index))...)
	}

	return diags
}

// ValueType returns the Value type.
func (m MapType) ValueType(_ context.Context) attr.Value {
	return MapValue{
		elementType: m.ElementType(),
	}
}

// ValueFromMap returns a MapValuable type given a Map.
func (m MapType) ValueFromMap(_ context.Context, ma MapValue) (MapValuable, diag.Diagnostics) {
	return ma, nil
}
