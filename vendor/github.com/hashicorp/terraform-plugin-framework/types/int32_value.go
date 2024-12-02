// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

type Int32 = basetypes.Int32Value

// Int32Null creates a Int32 with a null value. Determine whether the value is
// null via the Int32 type IsNull method.
func Int32Null() basetypes.Int32Value {
	return basetypes.NewInt32Null()
}

// Int32Unknown creates a Int32 with an unknown value. Determine whether the
// value is unknown via the Int32 type IsUnknown method.
func Int32Unknown() basetypes.Int32Value {
	return basetypes.NewInt32Unknown()
}

// Int32Value creates a Int32 with a known value. Access the value via the
// Int32 type ValueInt32 method.
func Int32Value(value int32) basetypes.Int32Value {
	return basetypes.NewInt32Value(value)
}

// Int32PointerValue creates a Int32 with a null value if nil or a known value.
func Int32PointerValue(value *int32) basetypes.Int32Value {
	return basetypes.NewInt32PointerValue(value)
}
