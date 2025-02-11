// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

type Float64 = basetypes.Float64Value

// Float64Null creates a Float64 with a null value. Determine whether the value is
// null via the Float64 type IsNull method.
func Float64Null() basetypes.Float64Value {
	return basetypes.NewFloat64Null()
}

// Float64Unknown creates a Float64 with an unknown value. Determine whether the
// value is unknown via the Float64 type IsUnknown method.
func Float64Unknown() basetypes.Float64Value {
	return basetypes.NewFloat64Unknown()
}

// Float64Value creates a Float64 with a known value. Access the value via the Float64
// type ValueFloat64 method.
func Float64Value(value float64) basetypes.Float64Value {
	return basetypes.NewFloat64Value(value)
}

// Float64PointerValue creates a Float64 with a null value if nil or a known value.
func Float64PointerValue(value *float64) basetypes.Float64Value {
	return basetypes.NewFloat64PointerValue(value)
}
