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

// StringTypable extends attr.Type for string types.
// Implement this interface to create a custom StringType type.
type StringTypable interface {
	attr.Type

	// ValueFromString should convert the String to a StringValuable type.
	ValueFromString(context.Context, StringValue) (StringValuable, diag.Diagnostics)
}

var _ StringTypable = StringType{}

// StringType is the base framework type for a string. StringValue is the
// associated value type.
type StringType struct{}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// type.
func (t StringType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is equivalent.
func (t StringType) Equal(o attr.Type) bool {
	_, ok := o.(StringType)

	return ok
}

// String returns a human readable string of the type name.
func (t StringType) String() string {
	return "basetypes.StringType"
}

// TerraformType returns the tftypes.Type that should be used to represent this
// framework type.
func (t StringType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t StringType) ValueFromString(_ context.Context, v StringValue) (StringValuable, diag.Diagnostics) {
	return v, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to
// convert the tftypes.Value into a more convenient Go type for the provider to
// consume the data with.
func (t StringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return NewStringUnknown(), nil
	}

	if in.IsNull() {
		return NewStringNull(), nil
	}

	var s string

	err := in.As(&s)

	if err != nil {
		return nil, err
	}

	return NewStringValue(s), nil
}

// ValueType returns the Value type.
func (t StringType) ValueType(_ context.Context) attr.Value {
	// This Value does not need to be valid.
	return StringValue{}
}
