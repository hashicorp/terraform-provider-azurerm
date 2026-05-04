// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package reflect

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FromValue is the inverse of Into, taking a Go value (`val`) and transforming it
// into an attr.Value using the attr.Type supplied. `val` will first be
// transformed into a tftypes.Value, then passed to `typ`'s ValueFromTerraform
// method.
func FromValue(ctx context.Context, typ attr.Type, val interface{}, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := val.(attr.Value); ok {
		return FromAttributeValue(ctx, typ, v, path)
	}
	if v, ok := val.(tftypes.ValueCreator); ok {
		return FromValueCreator(ctx, typ, v, path)
	}
	if v, ok := val.(Unknownable); ok {
		return FromUnknownable(ctx, typ, v, path)
	}
	if v, ok := val.(Nullable); ok {
		return FromNullable(ctx, typ, v, path)
	}
	if bf, ok := val.(*big.Float); ok {
		return FromBigFloat(ctx, typ, bf, path)
	}
	if bi, ok := val.(*big.Int); ok {
		return FromBigInt(ctx, typ, bi, path)
	}
	value := reflect.ValueOf(val)
	kind := value.Kind()
	switch kind {
	case reflect.Struct:
		t, ok := typ.(attr.TypeWithAttributeTypes)
		if !ok {
			err := fmt.Errorf("cannot use type %T as schema type %T; %T must be an attr.TypeWithAttributeTypes to hold %T", val, typ, typ, val)
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert from value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return nil, diags
		}
		return FromStruct(ctx, t, value, path)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return FromInt(ctx, typ, value.Int(), path)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return FromUint(ctx, typ, value.Uint(), path)
	case reflect.Float32, reflect.Float64:
		return FromFloat(ctx, typ, value.Float(), path)
	case reflect.Bool:
		return FromBool(ctx, typ, value.Bool(), path)
	case reflect.String:
		return FromString(ctx, typ, value.String(), path)
	case reflect.Slice:
		return FromSlice(ctx, typ, value, path)
	case reflect.Map:
		t, ok := typ.(attr.TypeWithElementType)
		if !ok {
			err := fmt.Errorf("cannot use type %T as schema type %T; %T must be an attr.TypeWithElementType to hold %T", val, typ, typ, val)
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert from value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return nil, diags
		}
		return FromMap(ctx, t, value, path)
	case reflect.Ptr:
		return FromPointer(ctx, typ, value, path)
	default:
		err := fmt.Errorf("cannot construct attr.Type from %T (%s)", val, kind)
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert from value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}
}
