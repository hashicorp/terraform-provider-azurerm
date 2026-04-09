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

// BoolTypable extends attr.Type for bool types.
// Implement this interface to create a custom BoolType type.
type BoolTypable interface {
	attr.Type

	// ValueFromBool should convert the Bool to a BoolValuable type.
	ValueFromBool(context.Context, BoolValue) (BoolValuable, diag.Diagnostics)
}

var _ BoolTypable = BoolType{}

// BoolType is the base framework type for a boolean. BoolValue is the
// associated value type.
type BoolType struct{}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// type.
func (t BoolType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is equivalent.
func (t BoolType) Equal(o attr.Type) bool {
	_, ok := o.(BoolType)

	return ok
}

// String returns a human readable string of the type name.
func (t BoolType) String() string {
	return "basetypes.BoolType"
}

// TerraformType returns the tftypes.Type that should be used to represent this
// framework type.
func (t BoolType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.Bool
}

// ValueFromBool returns a BoolValuable type given a BoolValue.
func (t BoolType) ValueFromBool(_ context.Context, v BoolValue) (BoolValuable, diag.Diagnostics) {
	return v, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to
// convert the tftypes.Value into a more convenient Go type for the provider to
// consume the data with.
func (t BoolType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return NewBoolUnknown(), nil
	}

	if in.IsNull() {
		return NewBoolNull(), nil
	}

	var v bool

	err := in.As(&v)

	if err != nil {
		return nil, err
	}

	return NewBoolValue(v), nil
}

// ValueType returns the Value type.
func (t BoolType) ValueType(_ context.Context) attr.Value {
	// This Value does not need to be valid.
	return BoolValue{}
}
