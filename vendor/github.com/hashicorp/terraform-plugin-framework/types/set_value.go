// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Set = basetypes.SetValue

// SetNull creates a Set with a null value. Determine whether the value is
// null via the Set type IsNull method.
func SetNull(elementType attr.Type) basetypes.SetValue {
	return basetypes.NewSetNull(elementType)
}

// SetUnknown creates a Set with an unknown value. Determine whether the
// value is unknown via the Set type IsUnknown method.
func SetUnknown(elementType attr.Type) basetypes.SetValue {
	return basetypes.NewSetUnknown(elementType)
}

// SetValue creates a Set with a known value. Access the value via the Set
// type Elements or ElementsAs methods.
func SetValue(elementType attr.Type, elements []attr.Value) (basetypes.SetValue, diag.Diagnostics) {
	return basetypes.NewSetValue(elementType, elements)
}

// SetValueFrom creates a Set with a known value, using reflection rules.
// The elements must be a slice which can convert into the given element type.
// Access the value via the Set type Elements or ElementsAs methods.
func SetValueFrom(ctx context.Context, elementType attr.Type, elements any) (basetypes.SetValue, diag.Diagnostics) {
	return basetypes.NewSetValueFrom(ctx, elementType, elements)
}

// SetValueMust creates a Set with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Set
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create Set values which will
// not potentially affect practitioners, such as testing, or exhaustively
// tested provider logic.
func SetValueMust(elementType attr.Type, elements []attr.Value) basetypes.SetValue {
	return basetypes.NewSetValueMust(elementType, elements)
}
