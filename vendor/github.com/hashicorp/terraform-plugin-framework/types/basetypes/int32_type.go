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
)

// Int32Typable extends attr.Type for int32 types.
// Implement this interface to create a custom Int32Type type.
type Int32Typable interface {
	attr.Type

	// ValueFromInt32 should convert the Int32 to a Int32Valuable type.
	ValueFromInt32(context.Context, Int32Value) (Int32Valuable, diag.Diagnostics)
}

var _ Int32Typable = Int32Type{}

// Int32Type is the base framework type for an integer number.
// Int32Value is the associated value type.
type Int32Type struct{}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// type.
func (t Int32Type) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is equivalent.
func (t Int32Type) Equal(o attr.Type) bool {
	_, ok := o.(Int32Type)

	return ok
}

// String returns a human-readable string of the type name.
func (t Int32Type) String() string {
	return "basetypes.Int32Type"
}

// TerraformType returns the tftypes.Type that should be used to represent this
// framework type.
func (t Int32Type) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.Number
}

// ValueFromInt32 returns a Int32Valuable type given a Int32Value.
func (t Int32Type) ValueFromInt32(_ context.Context, v Int32Value) (Int32Valuable, diag.Diagnostics) {
	return v, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to
// convert the tftypes.Value into a more convenient Go type for the provider to
// consume the data with.
func (t Int32Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return NewInt32Unknown(), nil
	}

	if in.IsNull() {
		return NewInt32Null(), nil
	}

	var bigF *big.Float
	err := in.As(&bigF)

	if err != nil {
		return nil, err
	}

	if !bigF.IsInt() {
		return nil, fmt.Errorf("Value %s is not an integer.", bigF)
	}

	i, accuracy := bigF.Int64()

	if accuracy != 0 {
		return nil, fmt.Errorf("Value %s cannot be represented as a 32-bit integer.", bigF)
	}

	if i < math.MinInt32 || i > math.MaxInt32 {
		return nil, fmt.Errorf("Value %s cannot be represented as a 32-bit integer.", bigF)
	}

	return NewInt32Value(int32(i)), nil
}

// ValueType returns the Value type.
func (t Int32Type) ValueType(_ context.Context) attr.Value {
	// This Value does not need to be valid.
	return Int32Value{}
}
