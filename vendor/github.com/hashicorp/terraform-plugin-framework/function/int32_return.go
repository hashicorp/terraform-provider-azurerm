// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Return = Int32Return{}

// Int32Return represents a function return that is a 32-bit integer number.
//
// When setting the value for this return:
//
// - If CustomType is set, use its associated value type.
// - Otherwise, use [types.Int32], *int32, or int32.
//
// Return documentation is expected in the function [Definition] documentation.
type Int32Return struct {
	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.Int32Type]. When setting data, the
	// [basetypes.Int32Valuable] implementation associated with this custom
	// type must be used in place of [types.Int32].
	CustomType basetypes.Int32Typable
}

// GetType returns the return data type.
func (r Int32Return) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.Int32Type{}
}

// NewResultData returns a new result data based on the type.
func (r Int32Return) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewInt32Unknown()

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromInt32(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}
