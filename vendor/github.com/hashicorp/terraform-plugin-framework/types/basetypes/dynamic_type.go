// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// DynamicTypable extends attr.Type for dynamic types. Implement this interface to create a custom DynamicType type.
type DynamicTypable interface {
	attr.Type

	// ValueFromDynamic should convert the DynamicValue to a DynamicValuable type.
	ValueFromDynamic(context.Context, DynamicValue) (DynamicValuable, diag.Diagnostics)
}

var _ DynamicTypable = DynamicType{}

// DynamicType is the base framework type for a dynamic. Static types are always
// preferable over dynamic types in Terraform as practitioners will receive less
// helpful configuration assistance from validation error diagnostics and editor
// integrations.
//
// DynamicValue is the associated value type and, when known, contains a concrete
// value of another framework type. (StringValue, ListValue, ObjectValue, MapValue, etc.)
type DynamicType struct{}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the type.
func (t DynamicType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	// MAINTAINER NOTE: Based on dynamic type alone, there is no alternative type information to return related to a path step.
	// When working with dynamics, we should always use DynamicValue to determine underlying type information.
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is equivalent.
//
// Dynamic types do not contain a reference to the underlying `attr.Value` type, so this equality check
// only asserts that both types are DynamicType.
func (t DynamicType) Equal(o attr.Type) bool {
	_, ok := o.(DynamicType)

	return ok
}

// String returns a human-friendly description of the DynamicType.
func (t DynamicType) String() string {
	return "basetypes.DynamicType"
}

// TerraformType returns the tftypes.Type that should be used to represent this type.
func (t DynamicType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.DynamicPseudoType
}

// ValueFromDynamic returns a DynamicValuable type given a DynamicValue.
func (t DynamicType) ValueFromDynamic(ctx context.Context, v DynamicValue) (DynamicValuable, diag.Diagnostics) {
	return v, nil
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value. This is meant to convert
// the tftypes.Value into a more convenient Go type for the provider to consume the data with.
func (t DynamicType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewDynamicNull(), nil
	}

	// For dynamic values, it's possible the incoming value is unknown but the concrete type itself is known. In this
	// situation, we can't return a dynamic unknown as we will lose that concrete type information. If the type here
	// is not dynamic, then we use the concrete `(attr.Type).ValueFromTerraform` below to produce the unknown value.
	if !in.IsKnown() && in.Type().Is(tftypes.DynamicPseudoType) {
		return NewDynamicUnknown(), nil
	}

	// For dynamic values, it's possible the incoming value is null but the concrete type itself is known. In this
	// situation, we can't return a dynamic null as we will lose that concrete type information. If the type here
	// is not dynamic, then we use the concrete `(attr.Type).ValueFromTerraform` below to produce the null value.
	if in.IsNull() && in.Type().Is(tftypes.DynamicPseudoType) {
		return NewDynamicNull(), nil
	}

	// MAINTAINER NOTE: It should not be possible for Terraform core to send a known value of `tftypes.DynamicPseudoType`.
	// This check prevents an infinite recursion that would result if this scenario occurs when attempting to create a dynamic value.
	if in.Type().Is(tftypes.DynamicPseudoType) {
		return nil, errors.New("ambiguous known value for `tftypes.DynamicPseudoType` detected")
	}

	attrType, err := tftypeToFrameworkType(in.Type())
	if err != nil {
		return nil, err
	}

	val, err := attrType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	return NewDynamicValue(val), nil
}

// ValueType returns the Value type.
func (t DynamicType) ValueType(_ context.Context) attr.Value {
	return DynamicValue{}
}

// tftypeToFrameworkType is a helper function that returns the framework type equivalent for a given Terraform type.
//
// Custom dynamic type implementations shouldn't need to override this method, but if needed, they can implement similar logic
// in their `ValueFromTerraform` implementation.
func tftypeToFrameworkType(in tftypes.Type) (attr.Type, error) {
	// Primitive types
	if in.Is(tftypes.Bool) {
		return BoolType{}, nil
	}
	if in.Is(tftypes.Number) {
		return NumberType{}, nil
	}
	if in.Is(tftypes.String) {
		return StringType{}, nil
	}

	if in.Is(tftypes.DynamicPseudoType) {
		// Null and Unknown values that do not have a type determined will have a type of DynamicPseudoType
		return DynamicType{}, nil
	}

	// Collection types
	if in.Is(tftypes.List{}) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
		l := in.(tftypes.List)

		elemType, err := tftypeToFrameworkType(l.ElementType)
		if err != nil {
			return nil, err
		}
		return ListType{ElemType: elemType}, nil
	}
	if in.Is(tftypes.Map{}) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
		m := in.(tftypes.Map)

		elemType, err := tftypeToFrameworkType(m.ElementType)
		if err != nil {
			return nil, err
		}

		return MapType{ElemType: elemType}, nil
	}
	if in.Is(tftypes.Set{}) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
		s := in.(tftypes.Set)

		elemType, err := tftypeToFrameworkType(s.ElementType)
		if err != nil {
			return nil, err
		}

		return SetType{ElemType: elemType}, nil
	}

	// Structural types
	if in.Is(tftypes.Object{}) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
		o := in.(tftypes.Object)

		attrTypes := make(map[string]attr.Type, len(o.AttributeTypes))
		for name, tfType := range o.AttributeTypes {
			t, err := tftypeToFrameworkType(tfType)
			if err != nil {
				return nil, err
			}
			attrTypes[name] = t
		}
		return ObjectType{AttrTypes: attrTypes}, nil
	}
	if in.Is(tftypes.Tuple{}) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
		tup := in.(tftypes.Tuple)

		elemTypes := make([]attr.Type, len(tup.ElementTypes))
		for i, tfType := range tup.ElementTypes {
			t, err := tftypeToFrameworkType(tfType)
			if err != nil {
				return nil, err
			}
			elemTypes[i] = t
		}
		return TupleType{ElemTypes: elemTypes}, nil
	}

	return nil, fmt.Errorf("unsupported tftypes.Type detected: %T", in)
}
