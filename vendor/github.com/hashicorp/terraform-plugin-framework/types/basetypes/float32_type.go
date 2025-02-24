// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
)

// Float32Typable extends attr.Type for float32 types.
// Implement this interface to create a custom Float32Type type.
type Float32Typable interface {
	attr.Type

	// ValueFromFloat32 should convert the Float32 to a Float32Valuable type.
	ValueFromFloat32(context.Context, Float32Value) (Float32Valuable, diag.Diagnostics)
}

var _ Float32Typable = Float32Type{}

// Float32Type is the base framework type for a floating point number.
// Float32Value is the associated value type.
type Float32Type struct{}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// type.
func (t Float32Type) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is equivalent.
func (t Float32Type) Equal(o attr.Type) bool {
	_, ok := o.(Float32Type)

	return ok
}

// String returns a human readable string of the type name.
func (t Float32Type) String() string {
	return "basetypes.Float32Type"
}

// TerraformType returns the tftypes.Type that should be used to represent this
// framework type.
func (t Float32Type) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.Number
}

// ValueFromFloat32 returns a Float32Valuable type given a Float32Value.
func (t Float32Type) ValueFromFloat32(_ context.Context, v Float32Value) (Float32Valuable, diag.Diagnostics) {
	return v, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to
// convert the tftypes.Value into a more convenient Go type for the provider to
// consume the data with.
func (t Float32Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return NewFloat32Unknown(), nil
	}

	if in.IsNull() {
		return NewFloat32Null(), nil
	}

	var bigF *big.Float
	err := in.As(&bigF)

	if err != nil {
		return nil, err
	}

	f, accuracy := bigF.Float32()
	f64, f64accuracy := bigF.Float64()

	if accuracy == big.Exact && f64accuracy == big.Exact {
		logging.FrameworkDebug(ctx, fmt.Sprintf("Float32Type ValueFromTerraform: big.Float value has distinct float32 and float64 representations "+
			"(float32 value: %f, float64 value: %f)", f, f64))
	}

	// Underflow
	// Reference: https://pkg.go.dev/math/big#Float.Float32
	if f == 0 && accuracy != big.Exact {
		return nil, fmt.Errorf("Value %s cannot be represented as a 32-bit floating point.", bigF)
	}

	// Overflow
	// Reference: https://pkg.go.dev/math/big#Float.Float32
	if math.IsInf(float64(f), 0) {
		return nil, fmt.Errorf("Value %s cannot be represented as a 32-bit floating point.", bigF)
	}

	// Underlying *big.Float values are not exposed with helper functions, so creating Float32Value via struct literal
	return Float32Value{
		state: attr.ValueStateKnown,
		value: bigF,
	}, nil
}

// ValueType returns the Value type.
func (t Float32Type) ValueType(_ context.Context) attr.Value {
	// This Value does not need to be valid.
	return Float32Value{}
}
