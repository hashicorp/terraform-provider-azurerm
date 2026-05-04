// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ BoolValuable = BoolValue{}
)

// BoolValuable extends attr.Value for boolean value types.
// Implement this interface to create a custom Bool value type.
type BoolValuable interface {
	attr.Value

	// ToBoolValue should convert the value type to a Bool.
	ToBoolValue(ctx context.Context) (BoolValue, diag.Diagnostics)
}

// BoolValuableWithSemanticEquals extends BoolValuable with semantic
// equality logic.
type BoolValuableWithSemanticEquals interface {
	BoolValuable

	// BoolSemanticEquals should return true if the given value is
	// semantically equal to the current value. This logic is used to prevent
	// Terraform data consistency errors and resource drift where a value change
	// may have inconsequential differences.
	//
	// Only known values are compared with this method as changing a value's
	// state implicitly represents a different value.
	BoolSemanticEquals(context.Context, BoolValuable) (bool, diag.Diagnostics)
}

// NewBoolNull creates a Bool with a null value. Determine whether the value is
// null via the Bool type IsNull method.
func NewBoolNull() BoolValue {
	return BoolValue{
		state: attr.ValueStateNull,
	}
}

// NewBoolUnknown creates a Bool with an unknown value. Determine whether the
// value is unknown via the Bool type IsUnknown method.
func NewBoolUnknown() BoolValue {
	return BoolValue{
		state: attr.ValueStateUnknown,
	}
}

// NewBoolValue creates a Bool with a known value. Access the value via the Bool
// type ValueBool method.
func NewBoolValue(value bool) BoolValue {
	return BoolValue{
		state: attr.ValueStateKnown,
		value: value,
	}
}

// NewBoolPointerValue creates a Bool with a null value if nil or a known
// value. Access the value via the Bool type ValueBoolPointer method.
func NewBoolPointerValue(value *bool) BoolValue {
	if value == nil {
		return NewBoolNull()
	}

	return NewBoolValue(*value)
}

// BoolValue represents a boolean value.
type BoolValue struct {
	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState

	// value contains the known value, if not null or unknown.
	value bool
}

// Type returns a BoolType.
func (b BoolValue) Type(_ context.Context) attr.Type {
	return BoolType{}
}

// ToTerraformValue returns the data contained in the Bool as a tftypes.Value.
func (b BoolValue) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	switch b.state {
	case attr.ValueStateKnown:
		if err := tftypes.ValidateValue(tftypes.Bool, b.value); err != nil {
			return tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(tftypes.Bool, b.value), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(tftypes.Bool, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Bool state in ToTerraformValue: %s", b.state))
	}
}

// Equal returns true if `other` is a *Bool and has the same value as `b`.
func (b BoolValue) Equal(other attr.Value) bool {
	o, ok := other.(BoolValue)

	if !ok {
		return false
	}

	if b.state != o.state {
		return false
	}

	if b.state != attr.ValueStateKnown {
		return true
	}

	return b.value == o.value
}

// IsNull returns true if the Bool represents a null value.
func (b BoolValue) IsNull() bool {
	return b.state == attr.ValueStateNull
}

// IsUnknown returns true if the Bool represents a currently unknown value.
func (b BoolValue) IsUnknown() bool {
	return b.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Bool value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (b BoolValue) String() string {
	if b.IsUnknown() {
		return attr.UnknownValueString
	}

	if b.IsNull() {
		return attr.NullValueString
	}

	return fmt.Sprintf("%t", b.value)
}

// ValueBool returns the known bool value. If Bool is null or unknown, returns
// false.
func (b BoolValue) ValueBool() bool {
	return b.value
}

// ValueBoolPointer returns a pointer to the known bool value, nil for a null
// value, or a pointer to false for an unknown value.
func (b BoolValue) ValueBoolPointer() *bool {
	if b.IsNull() {
		return nil
	}

	return &b.value
}

// ToBoolValue returns Bool.
func (b BoolValue) ToBoolValue(context.Context) (BoolValue, diag.Diagnostics) {
	return b, nil
}
