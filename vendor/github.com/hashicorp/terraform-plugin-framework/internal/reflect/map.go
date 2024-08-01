// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package reflect

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Map creates a map value that matches the type of `target`, and populates it
// with the contents of `val`.
func Map(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	underlyingValue := trueReflectValue(target)

	// this only works with maps, so check that out first
	if underlyingValue.Kind() != reflect.Map {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        fmt.Errorf("expected a map type, got %s", target.Type()),
		}))
		return target, diags
	}
	if !val.Type().Is(tftypes.Map{}) {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        fmt.Errorf("cannot reflect %s into a map, must be a map", val.Type().String()),
		}))
		return target, diags
	}
	elemTyper, ok := typ.(attr.TypeWithElementType)
	if !ok {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        fmt.Errorf("cannot reflect map using type information provided by %T, %T must be an attr.TypeWithElementType", typ, typ),
		}))
		return target, diags
	}

	// we need our value to become a map of values so we can iterate over
	// them and handle them individually
	values := map[string]tftypes.Value{}
	err := val.As(&values)
	if err != nil {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        err,
		}))
		return target, diags
	}

	// we need to know the type the slice is wrapping
	elemType := underlyingValue.Type().Elem()
	elemAttrType := elemTyper.ElementType()

	// we want an empty version of the map
	m := reflect.MakeMapWithSize(underlyingValue.Type(), len(values))

	// go over each of the values passed in, create a Go value of the right
	// type for them, and add it to our new map
	for key, value := range values {
		// create a new Go value of the type that can go in the map
		targetValue := reflect.Zero(elemType)

		// update our path so we can have nice errors
		path := path.AtMapKey(key)

		// reflect the value into our new target
		result, elemDiags := BuildValue(ctx, elemAttrType, value, targetValue, opts, path)
		diags.Append(elemDiags...)

		if diags.HasError() {
			return target, diags
		}

		m.SetMapIndex(reflect.ValueOf(key), result)
	}

	return m, diags
}

// FromMap returns an attr.Value representing the data contained in `val`.
// `val` must be a map type with keys that are a string type. The attr.Value
// will be of the type produced by `typ`.
//
// It is meant to be called through FromValue, not directly.
func FromMap(ctx context.Context, typ attr.TypeWithElementType, val reflect.Value, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfType := typ.TerraformType(ctx)

	if val.IsNil() {
		tfVal := tftypes.NewValue(tfType, nil)

		attrVal, err := typ.ValueFromTerraform(ctx, tfVal)

		if err != nil {
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert from map value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return nil, diags
		}

		switch t := attrVal.(type) {
		case xattr.ValidateableAttribute:
			resp := xattr.ValidateAttributeResponse{}

			t.ValidateAttribute(ctx,
				xattr.ValidateAttributeRequest{
					Path: path,
				},
				&resp,
			)

			diags.Append(resp.Diagnostics...)

			if diags.HasError() {
				return nil, diags
			}
		default:
			//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
			if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
				diags.Append(typeWithValidate.Validate(ctx, tfVal, path)...)

				if diags.HasError() {
					return nil, diags
				}
			}
		}

		return attrVal, diags
	}

	elemType := typ.ElementType()
	tfElems := map[string]tftypes.Value{}
	for _, key := range val.MapKeys() {
		if key.Kind() != reflect.String {
			err := fmt.Errorf("map keys must be strings, got %s", key.Type())
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert into a Terraform value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return nil, diags
		}

		mapKeyPath := path.AtMapKey(key.String())

		// If the element implements xattr.ValidateableAttribute, or xattr.TypeWithValidate,
		// and the element does not validate then diagnostics will be added here and returned
		// before reaching the switch statement below.
		val, valDiags := FromValue(ctx, elemType, val.MapIndex(key).Interface(), mapKeyPath)
		diags.Append(valDiags...)

		if diags.HasError() {
			return nil, diags
		}

		tfVal, err := val.ToTerraformValue(ctx)
		if err != nil {
			return nil, append(diags, toTerraformValueErrorDiag(err, path))
		}

		switch t := val.(type) {
		case xattr.ValidateableAttribute:
			resp := xattr.ValidateAttributeResponse{}

			t.ValidateAttribute(ctx,
				xattr.ValidateAttributeRequest{
					Path: mapKeyPath,
				},
				&resp,
			)

			diags.Append(resp.Diagnostics...)

			if diags.HasError() {
				return nil, diags
			}
		default:
			//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
			if typeWithValidate, ok := elemType.(xattr.TypeWithValidate); ok {
				diags.Append(typeWithValidate.Validate(ctx, tfVal, mapKeyPath)...)

				if diags.HasError() {
					return nil, diags
				}
			}
		}

		tfElems[key.String()] = tfVal
	}

	err := tftypes.ValidateValue(tfType, tfElems)
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}

	tfVal := tftypes.NewValue(tfType, tfElems)

	attrVal, err := typ.ValueFromTerraform(ctx, tfVal)

	if err != nil {
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert to map value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	switch t := attrVal.(type) {
	case xattr.ValidateableAttribute:
		resp := xattr.ValidateAttributeResponse{}

		t.ValidateAttribute(ctx,
			xattr.ValidateAttributeRequest{
				Path: path,
			},
			&resp,
		)

		diags.Append(resp.Diagnostics...)

		if diags.HasError() {
			return nil, diags
		}
	default:
		//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
		if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
			diags.Append(typeWithValidate.Validate(ctx, tfVal, path)...)

			if diags.HasError() {
				return nil, diags
			}
		}
	}

	return attrVal, diags
}
