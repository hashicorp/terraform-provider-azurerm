// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var _ Return = DynamicReturn{}

// DynamicReturn represents a function return that is a dynamic, rather
// than a static type. Static types are always preferable over dynamic
// types in Terraform as practitioners will receive less helpful configuration
// assistance from validation error diagnostics and editor integrations.
//
// When setting the value for this return:
//
// - If CustomType is set, use its associated value type.
// - Otherwise, use the [types.Dynamic] value type.
type DynamicReturn struct {
	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.DynamicType]. When setting data, the
	// [basetypes.DynamicValuable] implementation associated with this custom
	// type must be used in place of [types.Dynamic].
	CustomType basetypes.DynamicTypable
}

// GetType returns the return data type.
func (r DynamicReturn) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.DynamicType{}
}

// NewResultData returns a new result data based on the type.
func (r DynamicReturn) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewDynamicUnknown()

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromDynamic(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}
