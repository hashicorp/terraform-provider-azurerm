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

// build a slice of elements, matching the type of `target`, and fill it with
// the data in `val`.
func reflectSlice(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	// this only works with slices, so check that out first
	if target.Kind() != reflect.Slice {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        fmt.Errorf("expected a slice type, got %s", target.Type()),
		}))
		return target, diags
	}

	// we need our value to become a list of values so we can iterate over
	// them and handle them individually
	var values []tftypes.Value
	err := val.As(&values)
	if err != nil {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        err,
		}))
		return target, diags
	}

	switch t := typ.(type) {
	// List or Set
	case attr.TypeWithElementType:
		// we need to know the type the slice is wrapping
		elemType := target.Type().Elem()
		elemAttrType := t.ElementType()

		// we want an empty version of the slice
		slice := reflect.MakeSlice(target.Type(), 0, len(values))

		// go over each of the values passed in, create a Go value of the right
		// type for them, and add it to our new slice
		for pos, value := range values {
			// create a new Go value of the type that can go in the slice
			targetValue := reflect.Zero(elemType)

			// update our path so we can have nice errors
			valPath := path.AtListIndex(pos)

			if typ.TerraformType(ctx).Is(tftypes.Set{}) {
				attrVal, err := elemAttrType.ValueFromTerraform(ctx, value)

				if err != nil {
					diags.AddAttributeError(
						path,
						"Value Conversion Error",
						"An unexpected error was encountered trying to convert to slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
					)
					return target, diags
				}

				valPath = path.AtSetValue(attrVal)
			}

			// reflect the value into our new target
			val, valDiags := BuildValue(ctx, elemAttrType, value, targetValue, opts, valPath)
			diags.Append(valDiags...)

			if diags.HasError() {
				return target, diags
			}

			// add the new target to our slice
			slice = reflect.Append(slice, val)
		}

		return slice, diags

	// Tuple reflection into slices is currently limited to use-cases where all tuple element types are the same.
	//
	// Overall, Tuple support is limited in the framework, but the main path that executes tuple reflection is the provider-defined function variadic
	// parameter. All tuple elements in this variadic parameter will have the same element type. For use-cases where the variadic parameter is a dynamic type,
	// all elements will have the same type of `DynamicType` and value of `DynamicValue`, with an underlying value that may be different.
	case attr.TypeWithElementTypes:
		// we need to know the type the slice is wrapping
		elemType := target.Type().Elem()

		// we want an empty version of the slice
		slice := reflect.MakeSlice(target.Type(), 0, len(values))

		if len(t.ElementTypes()) <= 0 {
			// If the tuple values are empty as well, we can just pass back an empty slice of the type we received.
			if len(values) == 0 {
				return slice, diags
			}

			diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
				Val:        val,
				TargetType: target.Type(),
				Err:        fmt.Errorf("cannot reflect %s using type information provided by %T, tuple type contained no element types but received values", val.Type(), t),
			}))
			return target, diags
		}

		// Ensure that all tuple element types are the same by comparing each element type to the first
		multipleTypes := false
		allElemTypes := t.ElementTypes()
		elemAttrType := allElemTypes[0]
		for _, elemType := range allElemTypes[1:] {
			if !elemAttrType.Equal(elemType) {
				multipleTypes = true
				break
			}
		}

		if multipleTypes {
			diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
				Val:        val,
				TargetType: target.Type(),
				Err:        fmt.Errorf("cannot reflect %s using type information provided by %T, reflection support for tuples is limited to multiple elements of the same element type. Expected all element types to be %T", val.Type(), t, elemAttrType),
			}))
			return target, diags
		}

		// go over each of the values passed in, create a Go value of the right
		// type for them, and add it to our new slice
		for pos, value := range values {
			// create a new Go value of the type that can go in the slice
			targetValue := reflect.Zero(elemType)

			// update our path so we can have nice errors
			valPath := path.AtTupleIndex(pos)

			// reflect the value into our new target
			val, valDiags := BuildValue(ctx, elemAttrType, value, targetValue, opts, valPath)
			diags.Append(valDiags...)

			if diags.HasError() {
				return target, diags
			}

			// add the new target to our slice
			slice = reflect.Append(slice, val)
		}

		return slice, diags
	default:
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        fmt.Errorf("cannot reflect %s using type information provided by %T, %T must be an attr.TypeWithElementType or attr.TypeWithElementTypes", val.Type(), typ, typ),
		}))
		return target, diags
	}
}

// FromSlice returns an attr.Value as produced by `typ` using the data in
// `val`. `val` must be a slice. `typ` must be an attr.TypeWithElementType or
// attr.TypeWithElementTypes. If the slice is nil, the representation of null
// for `typ` will be returned. Otherwise, FromSlice will recurse into FromValue
// for each element in the slice, using the element type or types defined on
// `typ` to construct values for them.
//
// It is meant to be called through FromValue, not directly.
func FromSlice(ctx context.Context, typ attr.Type, val reflect.Value, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfType := typ.TerraformType(ctx)

	if val.IsNil() {
		tfVal := tftypes.NewValue(tfType, nil)

		attrVal, err := typ.ValueFromTerraform(ctx, tfVal)

		if err != nil {
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert from slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
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

	tfElems := make([]tftypes.Value, 0, val.Len())
	switch t := typ.(type) {
	// List or Set
	case attr.TypeWithElementType:
		elemType := t.ElementType()
		for i := 0; i < val.Len(); i++ {
			// The underlying reflect.Slice is fetched by Index(). For set types,
			// the path is value-based instead of index-based. Since there is only
			// the index until the value is retrieved, this will pass the
			// technically incorrect index-based path at first for framework
			// debugging purposes, then correct the path afterwards.
			valPath := path.AtListIndex(i)

			// If the element implements xattr.ValidateableAttribute, or xattr.TypeWithValidate,
			// and the element does not validate then diagnostics will be added here and returned
			// before reaching the switch statement below.
			val, valDiags := FromValue(ctx, elemType, val.Index(i).Interface(), valPath)
			diags.Append(valDiags...)

			if diags.HasError() {
				return nil, diags
			}

			tfVal, err := val.ToTerraformValue(ctx)
			if err != nil {
				return nil, append(diags, toTerraformValueErrorDiag(err, path))
			}

			if tfType.Is(tftypes.Set{}) {
				valPath = path.AtSetValue(val)
			}

			switch t := val.(type) {
			case xattr.ValidateableAttribute:
				resp := xattr.ValidateAttributeResponse{}

				t.ValidateAttribute(ctx,
					xattr.ValidateAttributeRequest{
						Path: valPath,
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
					diags.Append(typeWithValidate.Validate(ctx, tfVal, valPath)...)

					if diags.HasError() {
						return nil, diags
					}
				}
			}

			tfElems = append(tfElems, tfVal)
		}

	// Tuple reflection from slices is currently limited to use-cases where all tuple element types are the same.
	//
	// Overall, Tuple support is limited in the framework, but the main path that executes tuple reflection is the provider-defined function variadic
	// parameter. All tuple elements in this variadic parameter will have the same element type. For use-cases where the variadic parameter is a dynamic type,
	// all elements will have the same type of `DynamicType` and value of `DynamicValue`, with an underlying value that may be different.
	case attr.TypeWithElementTypes:
		if len(t.ElementTypes()) <= 0 {
			// If the tuple values are empty as well, we can just pass back an empty slice of the type we received.
			if val.Len() == 0 {
				break
			}

			err := fmt.Errorf("cannot use type %s as schema type %T; tuple type contained no element types but received values", val.Type(), t)
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert from slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return nil, diags
		}

		// Ensure that all tuple element types are the same by comparing each element type to the first
		multipleTypes := false
		allElemTypes := t.ElementTypes()
		elemAttrType := allElemTypes[0]
		for _, elemType := range allElemTypes[1:] {
			if !elemAttrType.Equal(elemType) {
				multipleTypes = true
				break
			}
		}

		if multipleTypes {
			err := fmt.Errorf("cannot use type %s as schema type %T; reflection support for tuples is limited to multiple elements of the same element type. Expected all element types to be %T", val.Type(), t, elemAttrType)
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert from slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return nil, diags
		}

		for i := 0; i < val.Len(); i++ {
			valPath := path.AtTupleIndex(i)

			// If the element implements xattr.ValidateableAttribute, or xattr.TypeWithValidate,
			// and the element does not validate then diagnostics will be added here and returned
			// before reaching the switch statement below.
			val, valDiags := FromValue(ctx, elemAttrType, val.Index(i).Interface(), valPath)
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
						Path: valPath,
					},
					&resp,
				)

				diags.Append(resp.Diagnostics...)

				if diags.HasError() {
					return nil, diags
				}
			default:
				//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
				if typeWithValidate, ok := elemAttrType.(xattr.TypeWithValidate); ok {
					diags.Append(typeWithValidate.Validate(ctx, tfVal, valPath)...)

					if diags.HasError() {
						return nil, diags
					}
				}
			}

			tfElems = append(tfElems, tfVal)
		}
	default:
		err := fmt.Errorf("cannot use type %s as schema type %T; %T must be an attr.TypeWithElementType or attr.TypeWithElementTypes", val.Type(), t, t)
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert from slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
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
			"An unexpected error was encountered trying to convert from slice value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
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
