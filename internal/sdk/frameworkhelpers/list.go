// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frameworkhelpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// WrappedListValidator provides a wrapper for legacy SDKv2 type validations to ease migration to Framework Native
// The provided function is tested against each element in the list, this simulates the SDKv2 behaviour of defining the
// valiation inside the `Elem:` property.
type WrappedListValidator struct {
	Func         func(v interface{}, k string) (warnings []string, errors []error)
	Desc         string
	MarkdownDesc string
}

func (w WrappedListValidator) Description(_ context.Context) string {
	return w.Desc
}

func (w WrappedListValidator) MarkdownDescription(_ context.Context) string {
	return w.MarkdownDesc
}

func (w WrappedListValidator) ValidateList(ctx context.Context, request validator.ListRequest, response *validator.ListResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	path := request.Path.String()

	switch request.ConfigValue.ElementType(ctx) {
	case basetypes.StringType{}, types.StringType:
		items := make([]string, 0)
		request.ConfigValue.ElementsAs(ctx, &items, false)
		for _, v := range items {
			_, errors := w.Func(v, path)
			if errors != nil && len(errors) > 0 {
				response.Diagnostics.AddError(fmt.Sprintf("invalid value for %s", path), fmt.Sprintf("%+v", errors[0]))
				return
			}
		}

	case basetypes.Int64Type{}, types.Int64Type:
		items := make([]int64, 0)
		request.ConfigValue.ElementsAs(ctx, &items, false)
		for _, v := range items {
			_, errors := w.Func(v, path)
			if errors != nil && len(errors) > 0 {
				response.Diagnostics.AddError(fmt.Sprintf("invalid value for %s", path), fmt.Sprintf("%+v", errors[0]))
				return
			}
		}

	case basetypes.Float64Type{}, types.Float64Type:
		items := make([]float64, 0)
		request.ConfigValue.ElementsAs(ctx, &items, false)
		for _, v := range items {
			_, errors := w.Func(v, path)
			if errors != nil && len(errors) > 0 {
				response.Diagnostics.AddError(fmt.Sprintf("invalid value for %s", path), fmt.Sprintf("%+v", errors[0]))
				return
			}
		}

	case basetypes.BoolType{}, types.BoolType:
		items := make([]bool, 0)
		request.ConfigValue.ElementsAs(ctx, &items, false)
		for _, v := range items {
			_, errors := w.Func(v, path)
			if errors != nil && len(errors) > 0 {
				response.Diagnostics.AddError(fmt.Sprintf("invalid value for %s", path), fmt.Sprintf("%+v", errors[0]))
				return
			}
		}
	default:
		response.Diagnostics.AddError(fmt.Sprintf("unsupported list validation wrapper type for %s", path), fmt.Sprintf("%+v", request.ConfigValue))
	}

}

var _ validator.List = &WrappedListValidator{}
