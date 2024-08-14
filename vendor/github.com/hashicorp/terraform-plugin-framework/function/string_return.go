// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Return = StringReturn{}

// StringReturn represents a function return that is a string.
//
// When setting the value for this return:
//
// - If CustomType is set, use its associated value type.
// - Otherwise, use [types.String], *string, or string.
type StringReturn struct {
	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.StringType]. When setting data, the
	// [basetypes.StringValuable] implementation associated with this custom
	// type must be used in place of [types.String].
	CustomType basetypes.StringTypable
}

// GetType returns the return data type.
func (r StringReturn) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.StringType{}
}

// NewResultData returns a new result data based on the type.
func (r StringReturn) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewStringUnknown()

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromString(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}
