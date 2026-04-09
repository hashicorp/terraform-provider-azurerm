// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var (
	_ NumberValuable = NumberValue{}
)

// NumberValuable extends attr.Value for number value types.
// Implement this interface to create a custom Number value type.
type NumberValuable interface {
	attr.Value

	// ToNumberValue should convert the value type to a Number.
	ToNumberValue(ctx context.Context) (NumberValue, diag.Diagnostics)
}

// NumberValuableWithSemanticEquals extends NumberValuable with semantic
// equality logic.
type NumberValuableWithSemanticEquals interface {
	NumberValuable

	// NumberSemanticEquals should return true if the given value is
	// semantically equal to the current value. This logic is used to prevent
	// Terraform data consistency errors and resource drift where a value change
	// may have inconsequential differences, such as rounding.
	//
	// Only known values are compared with this method as changing a value's
	// state implicitly represents a different value.
	NumberSemanticEquals(context.Context, NumberValuable) (bool, diag.Diagnostics)
}

// NewNumberNull creates a Number with a null value. Determine whether the value is
// null via the Number type IsNull method.
func NewNumberNull() NumberValue {
	return NumberValue{
		state: attr.ValueStateNull,
	}
}

// NewNumberUnknown creates a Number with an unknown value. Determine whether the
// value is unknown via the Number type IsUnknown method.
func NewNumberUnknown() NumberValue {
	return NumberValue{
		state: attr.ValueStateUnknown,
	}
}

// NewNumberValue creates a Number with a known value. Access the value via the Number
// type ValueBigFloat method. If the given value is nil, a null Number is created.
func NewNumberValue(value *big.Float) NumberValue {
	if value == nil {
		return NewNumberNull()
	}

	return NumberValue{
		state: attr.ValueStateKnown,
		value: value,
	}
}

// NumberValue represents a number value, exposed as a *big.Float. Numbers can be
// floats or integers.
type NumberValue struct {
	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState

	// value contains the known value, if not null or unknown.
	value *big.Float
}

// Type returns a NumberType.
func (n NumberValue) Type(_ context.Context) attr.Type {
	return NumberType{}
}

// ToTerraformValue returns the data contained in the Number as a tftypes.Value.
func (n NumberValue) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	switch n.state {
	case attr.ValueStateKnown:
		if n.value == nil {
			return tftypes.NewValue(tftypes.Number, nil), nil
		}

		if err := tftypes.ValidateValue(tftypes.Number, n.value); err != nil {
			return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(tftypes.Number, n.value), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(tftypes.Number, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Number state in ToTerraformValue: %s", n.state))
	}
}

// Equal returns true if `other` is a Number and has the same value as `n`.
func (n NumberValue) Equal(other attr.Value) bool {
	o, ok := other.(NumberValue)

	if !ok {
		return false
	}

	if n.state != o.state {
		return false
	}

	if n.state != attr.ValueStateKnown {
		return true
	}

	return n.value.Cmp(o.value) == 0
}

// IsNull returns true if the Number represents a null value.
func (n NumberValue) IsNull() bool {
	return n.state == attr.ValueStateNull
}

// IsUnknown returns true if the Number represents a currently unknown value.
func (n NumberValue) IsUnknown() bool {
	return n.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Number value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (n NumberValue) String() string {
	if n.IsUnknown() {
		return attr.UnknownValueString
	}

	if n.IsNull() {
		return attr.NullValueString
	}

	return n.value.String()
}

// ValueBigFloat returns the known *big.Float value. If Number is null or unknown, returns
// 0.0.
func (n NumberValue) ValueBigFloat() *big.Float {
	return n.value
}

// ToNumberValue returns Number.
func (n NumberValue) ToNumberValue(context.Context) (NumberValue, diag.Diagnostics) {
	return n, nil
}
