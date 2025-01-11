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
	_ Float32Valuable                   = Float32Value{}
	_ Float32ValuableWithSemanticEquals = Float32Value{}
)

// Float32Valuable extends attr.Value for float32 value types.
// Implement this interface to create a custom Float32 value type.
type Float32Valuable interface {
	attr.Value

	// ToFloat32Value should convert the value type to a Float32.
	ToFloat32Value(ctx context.Context) (Float32Value, diag.Diagnostics)
}

// Float32ValuableWithSemanticEquals extends Float32Valuable with semantic
// equality logic.
type Float32ValuableWithSemanticEquals interface {
	Float32Valuable

	// Float32SemanticEquals should return true if the given value is
	// semantically equal to the current value. This logic is used to prevent
	// Terraform data consistency errors and resource drift where a value change
	// may have inconsequential differences, such as rounding.
	//
	// Only known values are compared with this method as changing a value's
	// state implicitly represents a different value.
	Float32SemanticEquals(context.Context, Float32Valuable) (bool, diag.Diagnostics)
}

// NewFloat32Null creates a Float32 with a null value. Determine whether the value is
// null via the Float32 type IsNull method.
func NewFloat32Null() Float32Value {
	return Float32Value{
		state: attr.ValueStateNull,
	}
}

// NewFloat32Unknown creates a Float32 with an unknown value. Determine whether the
// value is unknown via the Float32 type IsUnknown method.
func NewFloat32Unknown() Float32Value {
	return Float32Value{
		state: attr.ValueStateUnknown,
	}
}

// NewFloat32Value creates a Float32 with a known value. Access the value via the Float32
// type ValueFloat32 method. Passing a value of `NaN` will result in a panic.
func NewFloat32Value(value float32) Float32Value {
	return Float32Value{
		state: attr.ValueStateKnown,
		value: big.NewFloat(float64(value)),
	}
}

// NewFloat32PointerValue creates a Float32 with a null value if nil or a known
// value. Access the value via the Float32 type ValueFloat32Pointer method.
// Passing a value of `NaN` will result in a panic.
func NewFloat32PointerValue(value *float32) Float32Value {
	if value == nil {
		return NewFloat32Null()
	}

	return NewFloat32Value(*value)
}

// Float32Value represents a 32-bit floating point value, exposed as a float32.
type Float32Value struct {
	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState

	// value contains the known value, if not null or unknown.
	value *big.Float
}

// Float32SemanticEquals returns true if the given Float32Value is semantically equal to the current Float32Value.
// The underlying value *big.Float can store more precise float values then the Go built-in float32 type. This only
// compares the precision of the value that can be represented as the Go built-in float32, which is 53 bits of precision.
func (f Float32Value) Float32SemanticEquals(ctx context.Context, newValuable Float32Valuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(Float32Value)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", f)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return f.ValueFloat32() == newValue.ValueFloat32(), diags
}

// Equal returns true if `other` is a Float32 and has the same value as `f`.
func (f Float32Value) Equal(other attr.Value) bool {
	o, ok := other.(Float32Value)

	if !ok {
		return false
	}

	if f.state != o.state {
		return false
	}

	if f.state != attr.ValueStateKnown {
		return true
	}

	// Not possible to create a known Float32Value with a nil value, but check anyways
	if f.value == nil || o.value == nil {
		return f.value == o.value
	}

	return f.value.Cmp(o.value) == 0
}

// ToTerraformValue returns the data contained in the Float32 as a tftypes.Value.
func (f Float32Value) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	switch f.state {
	case attr.ValueStateKnown:
		if err := tftypes.ValidateValue(tftypes.Number, f.value); err != nil {
			return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(tftypes.Number, f.value), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(tftypes.Number, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Float32 state in ToTerraformValue: %s", f.state))
	}
}

// Type returns a Float32Type.
func (f Float32Value) Type(ctx context.Context) attr.Type {
	return Float32Type{}
}

// IsNull returns true if the Float32 represents a null value.
func (f Float32Value) IsNull() bool {
	return f.state == attr.ValueStateNull
}

// IsUnknown returns true if the Float32 represents a currently unknown value.
func (f Float32Value) IsUnknown() bool {
	return f.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Float32 value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (f Float32Value) String() string {
	if f.IsUnknown() {
		return attr.UnknownValueString
	}

	if f.IsNull() {
		return attr.NullValueString
	}

	f32 := f.ValueFloat32()

	return fmt.Sprintf("%f", f32)
}

// ValueFloat32 returns the known float32 value. If Float32 is null or unknown, returns
// 0.0.
func (f Float32Value) ValueFloat32() float32 {
	if f.IsNull() || f.IsUnknown() {
		return float32(0.0)
	}

	f32, _ := f.value.Float32()
	return f32
}

// ValueFloat32Pointer returns a pointer to the known float32 value, nil for a
// null value, or a pointer to 0.0 for an unknown value.
func (f Float32Value) ValueFloat32Pointer() *float32 {
	if f.IsNull() {
		return nil
	}

	if f.IsUnknown() {
		f32 := float32(0.0)
		return &f32
	}

	f32, _ := f.value.Float32()
	return &f32
}

// ToFloat32Value returns Float32.
func (f Float32Value) ToFloat32Value(context.Context) (Float32Value, diag.Diagnostics) {
	return f, nil
}
