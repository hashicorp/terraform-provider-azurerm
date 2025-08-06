// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type List = basetypes.ListValue

// ListNull creates a List with a null value. Determine whether the value is
// null via the List type IsNull method.
func ListNull(elementType attr.Type) basetypes.ListValue {
	return basetypes.NewListNull(elementType)
}

// ListUnknown creates a List with an unknown value. Determine whether the
// value is unknown via the List type IsUnknown method.
func ListUnknown(elementType attr.Type) basetypes.ListValue {
	return basetypes.NewListUnknown(elementType)
}

// ListValue creates a List with a known value. Access the value via the List
// type Elements or ElementsAs methods.
func ListValue(elementType attr.Type, elements []attr.Value) (basetypes.ListValue, diag.Diagnostics) {
	return basetypes.NewListValue(elementType, elements)
}

// ListValueFrom creates a List with a known value, using reflection rules.
// The elements must be a slice which can convert into the given element type.
// Access the value via the List type Elements or ElementsAs methods.
func ListValueFrom(ctx context.Context, elementType attr.Type, elements any) (basetypes.ListValue, diag.Diagnostics) {
	return basetypes.NewListValueFrom(ctx, elementType, elements)
}

// ListValueMust creates a List with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the List
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create List values which will
// not potentially affect practitioners, such as testing, or exhaustively
// tested provider logic.
func ListValueMust(elementType attr.Type, elements []attr.Value) basetypes.ListValue {
	return basetypes.NewListValueMust(elementType, elements)
}
