// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Dynamic = basetypes.DynamicValue

// DynamicNull creates a Dynamic with a null value. Determine whether the value is
// null via the Dynamic type IsNull method.
func DynamicNull() basetypes.DynamicValue {
	return basetypes.NewDynamicNull()
}

// DynamicUnknown creates a Dynamic with an unknown value. Determine whether the
// value is unknown via the Dynamic type IsUnknown method.
func DynamicUnknown() basetypes.DynamicValue {
	return basetypes.NewDynamicUnknown()
}

// DynamicValue creates a Dynamic with a known value. Access the value via the Dynamic
// type UnderlyingValue method.
func DynamicValue(value attr.Value) basetypes.DynamicValue {
	return basetypes.NewDynamicValue(value)
}
