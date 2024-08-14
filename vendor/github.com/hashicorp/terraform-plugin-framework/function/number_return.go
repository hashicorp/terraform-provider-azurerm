// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Return = NumberReturn{}

// NumberReturn represents a function return that is a 512-bit arbitrary
// precision number.
//
// When setting the value for this return:
//
// - If CustomType is set, use its associated value type.
// - Otherwise, use [types.Number] or *big.Float.
type NumberReturn struct {
	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.NumberType]. When setting data, the
	// [basetypes.NumberValuable] implementation associated with this custom
	// type must be used in place of [types.Number].
	CustomType basetypes.NumberTypable
}

// GetType returns the return data type.
func (r NumberReturn) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.NumberType{}
}

// NewResultData returns a new result data based on the type.
func (r NumberReturn) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewNumberUnknown()

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromNumber(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}
