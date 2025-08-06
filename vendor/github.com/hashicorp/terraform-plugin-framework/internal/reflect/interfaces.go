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

// Unknownable is an interface for types that can be explicitly set to known or
// unknown.
type Unknownable interface {
	SetUnknown(context.Context, bool) error
	SetValue(context.Context, interface{}) error
	GetUnknown(context.Context) bool
	GetValue(context.Context) interface{}
}

// NewUnknownable creates a zero value of `target` (or the concrete type it's
// referencing, if it's a pointer) and calls its SetUnknown method.
//
// It is meant to be called through Into, not directly.
func NewUnknownable(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	receiver := pointerSafeZeroValue(ctx, target)
	method := receiver.MethodByName("SetUnknown")
	if !method.IsValid() {
		err := fmt.Errorf("cannot find SetUnknown method on type %s", receiver.Type().String())
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return target, diags
	}
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(!val.IsKnown()),
	})
	err := results[0].Interface()
	if err != nil {
		var underlyingErr error
		switch e := err.(type) {
		case error:
			underlyingErr = e
		default:
			underlyingErr = fmt.Errorf("unknown error type %T: %v", e, e)
		}
		underlyingErr = fmt.Errorf("reflection error: %w", underlyingErr)
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert into a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+underlyingErr.Error(),
		)
		return target, diags
	}
	return receiver, diags
}

// FromUnknownable creates an attr.Value from the data in an Unknownable.
//
// It is meant to be called through FromValue, not directly.
func FromUnknownable(ctx context.Context, typ attr.Type, val Unknownable, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if val.GetUnknown(ctx) {
		tfVal := tftypes.NewValue(typ.TerraformType(ctx), tftypes.UnknownValue)

		res, err := typ.ValueFromTerraform(ctx, tfVal)
		if err != nil {
			return nil, append(diags, valueFromTerraformErrorDiag(err, path))
		}

		switch t := res.(type) {
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

		return res, nil
	}
	err := tftypes.ValidateValue(typ.TerraformType(ctx), val.GetValue(ctx))
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}

	tfVal := tftypes.NewValue(typ.TerraformType(ctx), val.GetValue(ctx))

	res, err := typ.ValueFromTerraform(ctx, tfVal)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	switch t := res.(type) {
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

	return res, nil
}

// Nullable is an interface for types that can be explicitly set to null.
type Nullable interface {
	SetNull(context.Context, bool) error
	SetValue(context.Context, interface{}) error
	GetNull(context.Context) bool
	GetValue(context.Context) interface{}
}

// NewNullable creates a zero value of `target` (or the concrete type it's
// referencing, if it's a pointer) and calls its SetNull method.
//
// It is meant to be called through Into, not directly.
func NewNullable(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	receiver := pointerSafeZeroValue(ctx, target)
	method := receiver.MethodByName("SetNull")
	if !method.IsValid() {
		err := fmt.Errorf("cannot find SetNull method on type %s", receiver.Type().String())
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return target, diags
	}
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(val.IsNull()),
	})
	err := results[0].Interface()
	if err != nil {
		var underlyingErr error
		switch e := err.(type) {
		case error:
			underlyingErr = e
		default:
			underlyingErr = fmt.Errorf("unknown error type: %T", e)
		}
		underlyingErr = fmt.Errorf("reflection error: %w", underlyingErr)
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert into a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+underlyingErr.Error(),
		)
		return target, diags
	}
	return receiver, diags
}

// FromNullable creates an attr.Value from the data in a Nullable.
//
// It is meant to be called through FromValue, not directly.
func FromNullable(ctx context.Context, typ attr.Type, val Nullable, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if val.GetNull(ctx) {
		tfVal := tftypes.NewValue(typ.TerraformType(ctx), nil)

		res, err := typ.ValueFromTerraform(ctx, tfVal)
		if err != nil {
			return nil, append(diags, valueFromTerraformErrorDiag(err, path))
		}

		switch t := res.(type) {
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

		return res, nil
	}
	err := tftypes.ValidateValue(typ.TerraformType(ctx), val.GetValue(ctx))
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}

	tfVal := tftypes.NewValue(typ.TerraformType(ctx), val.GetValue(ctx))

	res, err := typ.ValueFromTerraform(ctx, tfVal)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	switch t := res.(type) {
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

	return res, diags
}

// NewValueConverter creates a zero value of `target` (or the concrete type
// it's referencing, if it's a pointer) and calls its FromTerraform5Value
// method.
//
// It is meant to be called through Into, not directly.
func NewValueConverter(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	receiver := pointerSafeZeroValue(ctx, target)
	method := receiver.MethodByName("FromTerraform5Value")
	if !method.IsValid() {
		err := fmt.Errorf("could not find FromTerraform5Type method on type %s", receiver.Type().String())
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert into a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return target, diags
	}
	results := method.Call([]reflect.Value{reflect.ValueOf(val)})
	err := results[0].Interface()
	if err != nil {
		var underlyingErr error
		switch e := err.(type) {
		case error:
			underlyingErr = e
		default:
			underlyingErr = fmt.Errorf("unknown error type: %T", e)
		}
		underlyingErr = fmt.Errorf("reflection error: %w", underlyingErr)
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert into a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+underlyingErr.Error(),
		)
		return target, diags
	}
	return receiver, diags
}

// FromValueCreator creates an attr.Value from the data in a
// tftypes.ValueCreator, calling its ToTerraform5Value method and converting
// the result to an attr.Value using `typ`.
//
// It is meant to be called from FromValue, not directly.
func FromValueCreator(ctx context.Context, typ attr.Type, val tftypes.ValueCreator, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	raw, err := val.ToTerraform5Value()
	if err != nil {
		return nil, append(diags, toTerraform5ValueErrorDiag(err, path))
	}
	err = tftypes.ValidateValue(typ.TerraformType(ctx), raw)
	if err != nil {
		return nil, append(diags, validateValueErrorDiag(err, path))
	}
	tfVal := tftypes.NewValue(typ.TerraformType(ctx), raw)

	res, err := typ.ValueFromTerraform(ctx, tfVal)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	switch t := res.(type) {
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

	return res, diags
}

// NewAttributeValue creates a new reflect.Value by calling the
// ValueFromTerraform method on `typ`. It will return an error if the returned
// `attr.Value` is not the same type as `target`.
//
// It is meant to be called through Into, not directly.
func NewAttributeValue(ctx context.Context, typ attr.Type, val tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	res, err := typ.ValueFromTerraform(ctx, val)
	if err != nil {
		return target, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	switch t := res.(type) {
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
			return target, diags
		}
	default:
		//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
		if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
			diags.Append(typeWithValidate.Validate(ctx, val, path)...)

			if diags.HasError() {
				return target, diags
			}
		}
	}

	if reflect.TypeOf(res) != target.Type() {
		diags.Append(diag.WithPath(path, DiagNewAttributeValueIntoWrongType{
			ValType:    reflect.TypeOf(res),
			TargetType: target.Type(),
			SchemaType: typ,
		}))
		return target, diags
	}
	return reflect.ValueOf(res), diags
}

// FromAttributeValue creates an attr.Value from an attr.Value. It just returns
// the attr.Value it is passed or an error if there is an unexpected mismatch
// between the attr.Type and attr.Value.
//
// It is meant to be called through FromValue, not directly.
func FromAttributeValue(ctx context.Context, typ attr.Type, val attr.Value, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Since the reflection logic is a generic Go type implementation with
	// user input, it is possible to get into awkward situations where
	// the logic is expecting a certain type while a value may not be
	// compatible. This check will ensure the framework raises its own
	// error is there is a mismatch, rather than a terraform-plugin-go
	// error or worse a panic.
	if !typ.TerraformType(ctx).Equal(val.Type(ctx).TerraformType(ctx)) {
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered while verifying an attribute value matched its expected type to prevent unexpected behavior or panics. "+
				"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Expected framework type from provider logic: %s / underlying type: %s\n", typ, typ.TerraformType(ctx))+
				fmt.Sprintf("Received framework type from provider logic: %s / underlying type: %s\n", val.Type(ctx), val.Type(ctx).TerraformType(ctx))+
				fmt.Sprintf("Path: %s", path),
		)

		return nil, diags
	}

	switch t := val.(type) {
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
			return val, diags
		}
	default:
		//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
		if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
			tfVal, err := val.ToTerraformValue(ctx)
			if err != nil {
				return val, append(diags, toTerraformValueErrorDiag(err, path))
			}

			diags.Append(typeWithValidate.Validate(ctx, tfVal, path)...)

			if diags.HasError() {
				return val, diags
			}
		}
	}

	return val, diags
}
