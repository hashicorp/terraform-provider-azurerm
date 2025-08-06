// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Return = Float32Return{}

// Float32Return represents a function return that is a 32-bit floating point
// number.
//
// When setting the value for this return:
//
// - If CustomType is set, use its associated value type.
// - Otherwise, use [types.Float32], *float32, or float32.
//
// Return documentation is expected in the function [Definition] documentation.
type Float32Return struct {
	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.Float32Type]. When setting data, the
	// [basetypes.Float32Valuable] implementation associated with this custom
	// type must be used in place of [types.Float32].
	CustomType basetypes.Float32Typable
}

// GetType returns the return data type.
func (r Float32Return) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.Float32Type{}
}

// NewResultData returns a new result data based on the type.
func (r Float32Return) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewFloat32Unknown()

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromFloat32(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}
