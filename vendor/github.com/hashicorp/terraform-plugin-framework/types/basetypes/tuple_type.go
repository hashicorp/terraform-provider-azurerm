// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Type                 = TupleType{}
	_ attr.TypeWithElementTypes = TupleType{}
)

// TupleType implements a tuple type definition. This type intentionally includes less functionality
// than other types in the type system as it has limited real world application and therefore
// is not exposed to provider developers.
type TupleType struct {
	// ElemTypes is an ordered list of element types for the tuple.
	ElemTypes []attr.Type
}

// ElementTypes returns the ordered attr.Type slice for the tuple.
func (t TupleType) ElementTypes() []attr.Type {
	return t.ElemTypes
}

// WithElementTypes returns a TupleType that is identical to `t`, but with the element types set to `types`.
func (t TupleType) WithElementTypes(types []attr.Type) attr.TypeWithElementTypes {
	return TupleType{ElemTypes: types}
}

// Equal returns true if `o` is also a TupleType and has the same ElemTypes in the same order.
func (t TupleType) Equal(o attr.Type) bool {
	other, ok := o.(TupleType)

	if !ok {
		return false
	}

	if len(t.ElemTypes) != len(other.ElemTypes) {
		return false
	}

	for i, elemType := range t.ElemTypes {
		if !elemType.Equal(other.ElemTypes[i]) {
			return false
		}
	}

	return true
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the tuple.
func (t TupleType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	indexStep, ok := step.(tftypes.ElementKeyInt)
	if !ok {
		return nil, fmt.Errorf("cannot apply step %T to TupleType", step)
	}

	index := int(indexStep)
	if index < 0 || index >= len(t.ElemTypes) {
		return nil, fmt.Errorf("no element defined at index %d in TupleType", index)
	}

	return t.ElemTypes[index], nil
}

// String returns a human-friendly description of the TupleType.
func (t TupleType) String() string {
	typeStrings := make([]string, len(t.ElemTypes))

	for i, elemType := range t.ElemTypes {
		typeStrings[i] = elemType.String()
	}

	return "types.TupleType[" + strings.Join(typeStrings, ", ") + "]"
}

// TerraformType returns the tftypes.Type that should be used to represent this type.
func (t TupleType) TerraformType(ctx context.Context) tftypes.Type {
	tfTypes := make([]tftypes.Type, len(t.ElemTypes))

	for i, elemType := range t.ElemTypes {
		tfTypes[i] = elemType.TerraformType(ctx)
	}

	return tftypes.Tuple{
		ElementTypes: tfTypes,
	}
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value. This is meant to convert
// the tftypes.Value into a more convenient Go type for the provider to consume the data with.
func (t TupleType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewTupleNull(t.ElementTypes()), nil
	}
	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}
	if !in.IsKnown() {
		return NewTupleUnknown(t.ElementTypes()), nil
	}
	if in.IsNull() {
		return NewTupleNull(t.ElementTypes()), nil
	}
	val := []tftypes.Value{}
	err := in.As(&val)
	if err != nil {
		return nil, err
	}
	elems := make([]attr.Value, 0, len(val))
	for i, elem := range val {
		// Accessing this index is safe because of the type comparison above
		av, err := t.ElemTypes[i].ValueFromTerraform(ctx, elem)
		if err != nil {
			return nil, err
		}
		elems = append(elems, av)
	}

	// ValueFromTerraform above on each element should make this safe.
	// Otherwise, this will need to do some Diagnostics to error conversion.
	return NewTupleValueMust(t.ElementTypes(), elems), nil
}

// ValueType returns the Value type.
func (t TupleType) ValueType(_ context.Context) attr.Value {
	return TupleValue{
		elementTypes: t.ElementTypes(),
	}
}
