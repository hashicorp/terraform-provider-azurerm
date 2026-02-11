// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package typehelpers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FlattenMapPointer[T any](input *map[string]T) (result types.Map, diags diag.Diagnostics) {
	outType := reflect.TypeOf(input).Elem().Elem().Kind()
	diags = make(diag.Diagnostics, 0)
	switch outType {
	case reflect.String:
		if len(*input) > 0 {
			result, diags = types.MapValueFrom(context.Background(), types.StringType, input)
		} else {
			result = types.MapNull(types.StringType)
		}
	case reflect.Int64:
		if len(*input) > 0 {
			result, diags = types.MapValueFrom(context.Background(), types.Int64Type, input)
		} else {
			result = types.MapNull(types.Int64Type)
		}
	case reflect.Float64:
		if len(*input) > 0 {
			result, diags = types.MapValueFrom(context.Background(), types.Float64Type, input)
		} else {
			result = types.MapNull(types.Float64Type)
		}
	case reflect.Bool:
		if len(*input) > 0 {
			result, diags = types.MapValueFrom(context.Background(), types.BoolType, input)
		} else {
			result = types.MapNull(types.BoolType)
		}
	default:
		diags.AddError("unsupported map element type", fmt.Sprintf("got %s", outType.String()))
	}

	return result, diags
}

func FlattenMap[T any](input map[string]T) (result types.Map, diags diag.Diagnostics) {
	outType := reflect.TypeOf(input).Elem().Kind()
	diags = make(diag.Diagnostics, 0)
	switch outType {
	case reflect.String:
		if len(input) > 0 {
			result, diags = types.MapValueFrom(context.Background(), types.StringType, input)
		} else {
			result = types.MapNull(types.StringType)
		}
	case reflect.Int64:
		if len(input) > 0 {
			result, diags = types.MapValueFrom(context.Background(), types.Int64Type, input)
		} else {
			result = types.MapNull(types.Int64Type)
		}
	case reflect.Float64:
		if len(input) > 0 {
			result, diags = types.MapValueFrom(context.Background(), types.Float64Type, input)
		} else {
			result = types.MapNull(types.Float64Type)
		}
	case reflect.Bool:
		if len(input) > 0 {
			result, diags = types.MapValueFrom(context.Background(), types.BoolType, input)
		} else {
			result = types.MapNull(types.BoolType)
		}
	default:
		diags.AddError("unsupported map element type", fmt.Sprintf("got %s", outType.String()))
	}

	return result, diags
}

func ExpandMap[T any](input types.Map) (result map[string]T, diags diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, diags
	}

	diags = input.ElementsAs(context.Background(), &result, false)

	return
}

func ExpandMapPointer[T any](input types.Map) (*map[string]T, diag.Diagnostics) {
	r, d := ExpandMap[T](input)

	if r == nil {
		return nil, d
	}

	return &r, d
}
