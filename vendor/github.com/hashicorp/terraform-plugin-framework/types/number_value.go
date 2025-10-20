// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Number = basetypes.NumberValue

// NumberNull creates a Number with a null value. Determine whether the value is
// null via the Number type IsNull method.
func NumberNull() basetypes.NumberValue {
	return basetypes.NewNumberNull()
}

// NumberUnknown creates a Number with an unknown value. Determine whether the
// value is unknown via the Number type IsUnknown method.
func NumberUnknown() basetypes.NumberValue {
	return basetypes.NewNumberUnknown()
}

// NumberValue creates a Number with a known value. Access the value via the Number
// type ValueBigFloat method. If the given value is nil, a null Number is created.
func NumberValue(value *big.Float) basetypes.NumberValue {
	return basetypes.NewNumberValue(value)
}
