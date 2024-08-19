// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Object = basetypes.ObjectValue

// ObjectNull creates a Object with a null value. Determine whether the value is
// null via the Object type IsNull method.
func ObjectNull(attributeTypes map[string]attr.Type) basetypes.ObjectValue {
	return basetypes.NewObjectNull(attributeTypes)
}

// ObjectUnknown creates a Object with an unknown value. Determine whether the
// value is unknown via the Object type IsUnknown method.
func ObjectUnknown(attributeTypes map[string]attr.Type) basetypes.ObjectValue {
	return basetypes.NewObjectUnknown(attributeTypes)
}

// ObjectValue creates a Object with a known value. Access the value via the Object
// type Attributes or As methods.
func ObjectValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (basetypes.ObjectValue, diag.Diagnostics) {
	return basetypes.NewObjectValue(attributeTypes, attributes)
}

// ObjectValueFrom creates a Object with a known value, using reflection rules.
// The attributes must be a struct which can convert into the given attribute types.
// Access the value via the Object type Attributes or As methods.
func ObjectValueFrom(ctx context.Context, attributeTypes map[string]attr.Type, attributes any) (basetypes.ObjectValue, diag.Diagnostics) {
	return basetypes.NewObjectValueFrom(ctx, attributeTypes, attributes)
}

// ObjectValueMust creates a Object with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Object
// type Attributes or As methods.
//
// This creation function is only recommended to create Object values which will
// not potentially affect practitioners, such as testing, or exhaustively
// tested provider logic.
func ObjectValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) basetypes.ObjectValue {
	return basetypes.NewObjectValueMust(attributeTypes, attributes)
}
