// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

type Float32 = basetypes.Float32Value

// Float32Null creates a Float32 with a null value. Determine whether the value is
// null via the Float32 type IsNull method.
func Float32Null() basetypes.Float32Value {
	return basetypes.NewFloat32Null()
}

// Float32Unknown creates a Float32 with an unknown value. Determine whether the
// value is unknown via the Float32 type IsUnknown method.
func Float32Unknown() basetypes.Float32Value {
	return basetypes.NewFloat32Unknown()
}

// Float32Value creates a Float32 with a known value. Access the value via the Float32
// type ValueFloat32 method.
func Float32Value(value float32) basetypes.Float32Value {
	return basetypes.NewFloat32Value(value)
}

// Float32PointerValue creates a Float32 with a null value if nil or a known value.
func Float32PointerValue(value *float32) basetypes.Float32Value {
	return basetypes.NewFloat32PointerValue(value)
}
