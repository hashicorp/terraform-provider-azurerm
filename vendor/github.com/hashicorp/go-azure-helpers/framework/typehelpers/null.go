// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package typehelpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func NullValueOf(ctx context.Context, v any) (attr.Value, error) {
	var attrType attr.Type
	var tfType tftypes.Type

	switch t := v.(type) {
	case basetypes.BoolValuable:
		attrType = t.Type(ctx)
		tfType = tftypes.Bool
	case basetypes.Float64Valuable:
		attrType = t.Type(ctx)
		tfType = tftypes.Number
	case basetypes.Int64Valuable:
		attrType = t.Type(ctx)
		tfType = tftypes.Number
	case basetypes.StringValuable:
		attrType = t.Type(ctx)
		tfType = tftypes.String
	case basetypes.ListValuable:
		attrType = t.Type(ctx)
		if v, ok := attrType.(attr.TypeWithElementType); ok {
			tfType = tftypes.List{ElementType: v.ElementType().TerraformType(ctx)}
		} else {
			tfType = tftypes.List{}
		}
	case basetypes.SetValuable:
		attrType = t.Type(ctx)
		if tWithE, ok := attrType.(attr.TypeWithElementType); ok {
			tfType = tftypes.Set{ElementType: tWithE.ElementType().TerraformType(ctx)}
		} else {
			tfType = tftypes.Set{}
		}
	case basetypes.MapValuable:
		attrType = t.Type(ctx)
		if tWithE, ok := attrType.(attr.TypeWithElementType); ok {
			tfType = tftypes.Map{ElementType: tWithE.ElementType().TerraformType(ctx)}
		} else {
			tfType = tftypes.Map{}
		}
	case basetypes.ObjectValuable:
		attrType = t.Type(ctx)
		if tWithE, ok := attrType.(attr.TypeWithAttributeTypes); ok {
			tfType = tftypes.Object{AttributeTypes: translateMapTypes(tWithE.AttributeTypes(), func(attrType attr.Type) tftypes.Type {
				return attrType.TerraformType(ctx)
			})}
		} else {
			tfType = tftypes.Object{}
		}
	default:
		return nil, nil
	}

	return attrType.ValueFromTerraform(ctx, tftypes.NewValue(tfType, nil))
}

func translateMapTypes[M ~map[K]V1, K comparable, V1, V2 any](m M, f func(V1) V2) map[K]V2 {
	n := make(map[K]V2, len(m))

	for k, v := range m {
		n[k] = f(v)
	}

	return n
}
