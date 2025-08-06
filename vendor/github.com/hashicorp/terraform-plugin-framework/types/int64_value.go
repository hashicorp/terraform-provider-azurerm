// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

type Int64 = basetypes.Int64Value

// Int64Null creates a Int64 with a null value. Determine whether the value is
// null via the Int64 type IsNull method.
func Int64Null() basetypes.Int64Value {
	return basetypes.NewInt64Null()
}

// Int64Unknown creates a Int64 with an unknown value. Determine whether the
// value is unknown via the Int64 type IsUnknown method.
func Int64Unknown() basetypes.Int64Value {
	return basetypes.NewInt64Unknown()
}

// Int64Value creates a Int64 with a known value. Access the value via the
// Int64 type ValueInt64 method.
func Int64Value(value int64) basetypes.Int64Value {
	return basetypes.NewInt64Value(value)
}

// Int64PointerValue creates a Int64 with a null value if nil or a known value.
func Int64PointerValue(value *int64) basetypes.Int64Value {
	return basetypes.NewInt64PointerValue(value)
}
