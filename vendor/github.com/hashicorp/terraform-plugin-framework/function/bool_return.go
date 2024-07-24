// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Return = BoolReturn{}

// BoolReturn represents a function return that is a boolean.
//
// When setting the value for this return:
//
// - If CustomType is set, use its associated value type.
// - Otherwise, use [types.Bool], *bool, or bool.
type BoolReturn struct {
	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.BoolType]. When setting data, the
	// [basetypes.BoolValuable] implementation associated with this custom
	// type must be used in place of [types.Bool].
	CustomType basetypes.BoolTypable
}

// GetType returns the return data type.
func (r BoolReturn) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.BoolType{}
}

// NewResultData returns a new result data based on the type.
func (r BoolReturn) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewBoolUnknown()

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromBool(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}
