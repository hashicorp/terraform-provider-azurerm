// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Map = basetypes.MapValue

// MapNull creates a Map with a null value. Determine whether the value is
// null via the Map type IsNull method.
func MapNull(elementType attr.Type) basetypes.MapValue {
	return basetypes.NewMapNull(elementType)
}

// MapUnknown creates a Map with an unknown value. Determine whether the
// value is unknown via the Map type IsUnknown method.
func MapUnknown(elementType attr.Type) basetypes.MapValue {
	return basetypes.NewMapUnknown(elementType)
}

// MapValue creates a Map with a known value. Access the value via the Map
// type Elements or ElementsAs methods.
func MapValue(elementType attr.Type, elements map[string]attr.Value) (basetypes.MapValue, diag.Diagnostics) {
	return basetypes.NewMapValue(elementType, elements)
}

// MapValueFrom creates a Map with a known value, using reflection rules.
// The elements must be a map which can convert into the given element type.
// Access the value via the Map type Elements or ElementsAs methods.
func MapValueFrom(ctx context.Context, elementType attr.Type, elements any) (basetypes.MapValue, diag.Diagnostics) {
	return basetypes.NewMapValueFrom(ctx, elementType, elements)
}

// MapValueMust creates a Map with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Map
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create Map values which will
// not potentially affect practitioners, such as testing, or exhaustively
// tested provider logic.
func MapValueMust(elementType attr.Type, elements map[string]attr.Value) basetypes.MapValue {
	return basetypes.NewMapValueMust(elementType, elements)
}
