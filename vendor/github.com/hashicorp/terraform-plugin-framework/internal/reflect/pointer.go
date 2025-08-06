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

// Pointer builds a new zero value of the concrete type that `target`
// references, populates it with BuildValue, and takes a pointer to it.
//
// It is meant to be called through Into, not directly.
func Pointer(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if target.Kind() != reflect.Ptr {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        val,
			TargetType: target.Type(),
			Err:        fmt.Errorf("cannot dereference pointer, not a pointer, is a %s (%s)", target.Type(), target.Kind()),
		}))
		return target, diags
	}
	// we may have gotten a nil pointer, so we need to create our own that
	// we can set
	pointer := reflect.New(target.Type().Elem())
	// build out whatever the pointer is pointing to
	pointed, pointedDiags := BuildValue(ctx, typ, val, pointer.Elem(), opts, path)
	diags.Append(pointedDiags...)

	if diags.HasError() {
		return target, diags
	}
	// to be able to set the pointer to our new pointer, we need to create
	// a pointer to the pointer
	pointerPointer := reflect.New(pointer.Type())
	// we set the pointer we created on the pointer to the pointer
	pointerPointer.Elem().Set(pointer)
	// then it's settable, so we can now set the concrete value we created
	// on the pointer
	pointerPointer.Elem().Elem().Set(pointed)
	// return the pointer we created
	return pointerPointer.Elem(), diags
}

// create a zero value of concrete type underlying any number of pointers, then
// wrap it in that number of pointers again. The end result is to wind up with
// the same exact type, except now you can be sure it's pointing to actual data
// and will not give you a nil pointer dereference panic unexpectedly.
func pointerSafeZeroValue(_ context.Context, target reflect.Value) reflect.Value {
	pointer := target.Type()
	var pointers int
	for pointer.Kind() == reflect.Ptr {
		pointer = pointer.Elem()
		pointers++
	}
	receiver := reflect.Zero(pointer)
	for i := 0; i < pointers; i++ {
		newReceiver := reflect.New(receiver.Type())
		newReceiver.Elem().Set(receiver)
		receiver = newReceiver
	}
	return receiver
}

// FromPointer turns a pointer into an attr.Value using `typ`. If the pointer
// is nil, the attr.Value will use its null representation. If it is not nil,
// it will recurse into FromValue to find the attr.Value of the type the value
// the pointer is referencing.
//
// It is meant to be called through FromValue, not directly.
func FromPointer(ctx context.Context, typ attr.Type, value reflect.Value, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if value.Kind() != reflect.Ptr {
		err := fmt.Errorf("cannot use type %s as a pointer", value.Type())
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert from pointer value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}
	if value.IsNil() {
		tfVal := tftypes.NewValue(typ.TerraformType(ctx), nil)

		attrVal, err := typ.ValueFromTerraform(ctx, tfVal)

		if err != nil {
			diags.AddAttributeError(
				path,
				"Value Conversion Error",
				"An unexpected error was encountered trying to convert from pointer value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
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

	attrVal, attrValDiags := FromValue(ctx, typ, value.Elem().Interface(), path)
	diags.Append(attrValDiags...)

	return attrVal, diags
}
