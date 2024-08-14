// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Return = Int64Return{}

// Int64Return represents a function return that is a 64-bit integer number.
//
// When setting the value for this return:
//
// - If CustomType is set, use its associated value type.
// - Otherwise, use [types.Int64], *int64, or int64.
type Int64Return struct {
	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.Int64Type]. When setting data, the
	// [basetypes.Int64Valuable] implementation associated with this custom
	// type must be used in place of [types.Int64].
	CustomType basetypes.Int64Typable
}

// GetType returns the return data type.
func (r Int64Return) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.Int64Type{}
}

// NewResultData returns a new result data based on the type.
func (r Int64Return) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewInt64Unknown()

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromInt64(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}
