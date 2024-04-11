// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Int64Typable extends attr.Type for int64 types.
// Implement this interface to create a custom Int64Type type.
type Int64Typable interface {
	xattr.TypeWithValidate

	// ValueFromInt64 should convert the Int64 to a Int64Valuable type.
	ValueFromInt64(context.Context, Int64Value) (Int64Valuable, diag.Diagnostics)
}

var _ Int64Typable = Int64Type{}

// Int64Type is the base framework type for an integer number.
// Int64Value is the associated value type.
type Int64Type struct{}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// type.
func (t Int64Type) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is equivalent.
func (t Int64Type) Equal(o attr.Type) bool {
	_, ok := o.(Int64Type)

	return ok
}

// String returns a human readable string of the type name.
func (t Int64Type) String() string {
	return "basetypes.Int64Type"
}

// TerraformType returns the tftypes.Type that should be used to represent this
// framework type.
func (t Int64Type) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.Number
}

// Validate implements type validation.
func (t Int64Type) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Equal(tftypes.Number) {
		diags.AddAttributeError(
			path,
			"Int64 Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Expected Number value, received %T with value: %v", in, in),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var value *big.Float
	err := in.As(&value)

	if err != nil {
		diags.AddAttributeError(
			path,
			"Int64 Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Cannot convert value to big.Float: %s", err),
		)
		return diags
	}

	if !value.IsInt() {
		diags.AddAttributeError(
			path,
			"Int64 Type Validation Error",
			fmt.Sprintf("Value %s is not an integer.", value),
		)
		return diags
	}

	_, accuracy := value.Int64()

	if accuracy != 0 {
		diags.AddAttributeError(
			path,
			"Int64 Type Validation Error",
			fmt.Sprintf("Value %s cannot be represented as a 64-bit integer.", value),
		)
		return diags
	}

	return diags
}

// ValueFromInt64 returns a Int64Valuable type given a Int64Value.
func (t Int64Type) ValueFromInt64(_ context.Context, v Int64Value) (Int64Valuable, diag.Diagnostics) {
	return v, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to
// convert the tftypes.Value into a more convenient Go type for the provider to
// consume the data with.
func (t Int64Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return NewInt64Unknown(), nil
	}

	if in.IsNull() {
		return NewInt64Null(), nil
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
		return nil, fmt.Errorf("Value %s cannot be represented as a 64-bit integer.", bigF)
	}

	return NewInt64Value(i), nil
}

// ValueType returns the Value type.
func (t Int64Type) ValueType(_ context.Context) attr.Value {
	// This Value does not need to be valid.
	return Int64Value{}
}
