// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Return = Float64Return{}

// Float64Return represents a function return that is a 64-bit floating point
// number.
//
// When setting the value for this return:
//
// - If CustomType is set, use its associated value type.
// - Otherwise, use [types.Float64], *float64, or float64.
type Float64Return struct {
	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.Float64Type]. When setting data, the
	// [basetypes.Float64Valuable] implementation associated with this custom
	// type must be used in place of [types.Float64].
	CustomType basetypes.Float64Typable
}

// GetType returns the return data type.
func (r Float64Return) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.Float64Type{}
}

// NewResultData returns a new result data based on the type.
func (r Float64Return) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewFloat64Unknown()

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromFloat64(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}
