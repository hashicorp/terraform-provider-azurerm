// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TerraformTypeToFrameworkType is a helper function that returns the framework type equivalent for a given Terraform type.
func TerraformTypeToFrameworkType(in tftypes.Type) (attr.Type, error) {
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
		return DynamicType{}, nil
	}

	// Collection types
	if in.Is(tftypes.List{}) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
		l := in.(tftypes.List)

		elemType, err := TerraformTypeToFrameworkType(l.ElementType)
		if err != nil {
			return nil, err
		}
		return ListType{ElemType: elemType}, nil
	}
	if in.Is(tftypes.Map{}) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
		m := in.(tftypes.Map)

		elemType, err := TerraformTypeToFrameworkType(m.ElementType)
		if err != nil {
			return nil, err
		}

		return MapType{ElemType: elemType}, nil
	}
	if in.Is(tftypes.Set{}) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
		s := in.(tftypes.Set)

		elemType, err := TerraformTypeToFrameworkType(s.ElementType)
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
			t, err := TerraformTypeToFrameworkType(tfType)
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
			t, err := TerraformTypeToFrameworkType(tfType)
			if err != nil {
				return nil, err
			}
			elemTypes[i] = t
		}
		return TupleType{ElemTypes: elemTypes}, nil
	}

	return nil, fmt.Errorf("unsupported tftypes.Type detected: %T", in)
}
