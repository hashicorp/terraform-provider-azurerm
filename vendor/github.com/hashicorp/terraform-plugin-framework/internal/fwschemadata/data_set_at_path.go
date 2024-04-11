// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// SetAtPath sets the attribute at `path` using the supplied Go value.
//
// The attribute path and value must be valid with the current schema. If the
// attribute path already has a value, it will be overwritten. If the attribute
// path does not have a value, it will be added, including any parent attribute
// paths as necessary.
//
// Lists can only have the next element added according to the current length.
func (d *Data) SetAtPath(ctx context.Context, path path.Path, val interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	ctx = logging.FrameworkWithAttributePath(ctx, path.String())

	tftypesPath, tftypesPathDiags := totftypes.AttributePath(ctx, path)

	diags.Append(tftypesPathDiags...)

	if diags.HasError() {
		return diags
	}

	attrType, err := d.Schema.TypeAtTerraformPath(ctx, tftypesPath)

	if err != nil {
		diags.AddAttributeError(
			path,
			d.Description.Title()+" Write Error",
			"An unexpected error was encountered trying to retrieve type information at a given path. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Error: "+err.Error(),
		)
		return diags
	}

	newVal, newValDiags := reflect.FromValue(ctx, attrType, val, path)
	diags.Append(newValDiags...)

	if diags.HasError() {
		return diags
	}

	tfVal, err := newVal.ToTerraformValue(ctx)

	if err != nil {
		diags.AddAttributeError(
			path,
			d.Description.Title()+" Write Error",
			"An unexpected error was encountered trying to write an attribute to the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Error: Cannot run ToTerraformValue on new data value: "+err.Error(),
		)
		return diags
	}

	if attrTypeWithValidate, ok := attrType.(xattr.TypeWithValidate); ok {
		logging.FrameworkTrace(ctx, "Type implements TypeWithValidate")
		logging.FrameworkTrace(ctx, "Calling provider defined Type Validate")
		diags.Append(attrTypeWithValidate.Validate(ctx, tfVal, path)...)
		logging.FrameworkTrace(ctx, "Called provider defined Type Validate")

		if diags.HasError() {
			return diags
		}
	}

	transformFunc, transformFuncDiags := d.SetAtPathTransformFunc(ctx, path, tfVal, nil)
	diags.Append(transformFuncDiags...)

	if diags.HasError() {
		return diags
	}

	d.TerraformValue, err = tftypes.Transform(d.TerraformValue, transformFunc)

	if err != nil {
		diags.AddAttributeError(
			path,
			d.Description.Title()+" Write Error",
			"An unexpected error was encountered trying to write an attribute to the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Error: Cannot transform data: "+err.Error(),
		)
		return diags
	}

	return diags
}

// SetAttributeTransformFunc recursively creates a value based on the current
// Plan values along the path. If the value at the path does not yet exist,
// this will perform recursion to add the child value to a parent value,
// creating the parent value if necessary.
func (d Data) SetAtPathTransformFunc(ctx context.Context, path path.Path, tfVal tftypes.Value, diags diag.Diagnostics) (func(*tftypes.AttributePath, tftypes.Value) (tftypes.Value, error), diag.Diagnostics) {
	exists, pathExistsDiags := d.PathExists(ctx, path)
	diags.Append(pathExistsDiags...)

	if diags.HasError() {
		return nil, diags
	}

	tftypesPath, tftypesPathDiags := totftypes.AttributePath(ctx, path)

	diags.Append(tftypesPathDiags...)

	if diags.HasError() {
		return nil, diags
	}

	if exists {
		// Overwrite existing value
		return func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
			if p.Equal(tftypesPath) {
				return tfVal, nil
			}
			return v, nil
		}, diags
	}

	parentPath := path.ParentPath()
	parentTftypesPath := tftypesPath.WithoutLastStep()
	parentAttrType, err := d.Schema.TypeAtTerraformPath(ctx, parentTftypesPath)

	if err != nil {
		err = fmt.Errorf("error getting parent attribute type in schema: %w", err)
		diags.AddAttributeError(
			parentPath,
			d.Description.Title()+" Write Error",
			"An unexpected error was encountered trying to write an attribute to the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	parentValue, err := d.TerraformValueAtTerraformPath(ctx, parentTftypesPath)

	if err != nil && !errors.Is(err, tftypes.ErrInvalidStep) {
		diags.AddAttributeError(
			parentPath,
			d.Description.Title()+" Read Error",
			"An unexpected error was encountered trying to read an attribute from the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	if parentValue.IsNull() || !parentValue.IsKnown() {
		parentType := parentAttrType.TerraformType(ctx)
		var childValue interface{}

		if !parentValue.IsKnown() {
			childValue = tftypes.UnknownValue
		}

		var parentValueDiags diag.Diagnostics
		parentValue, parentValueDiags = CreateParentTerraformValue(ctx, parentPath, parentType, childValue)
		diags.Append(parentValueDiags...)

		if diags.HasError() {
			return nil, diags
		}
	}

	var childValueDiags diag.Diagnostics
	childStep, _ := path.Steps().LastStep()
	parentValue, childValueDiags = UpsertChildTerraformValue(ctx, parentPath, parentValue, childStep, tfVal)
	diags.Append(childValueDiags...)

	if diags.HasError() {
		return nil, diags
	}

	if attrTypeWithValidate, ok := parentAttrType.(xattr.TypeWithValidate); ok {
		logging.FrameworkTrace(ctx, "Type implements TypeWithValidate")
		logging.FrameworkTrace(ctx, "Calling provider defined Type Validate")
		diags.Append(attrTypeWithValidate.Validate(ctx, parentValue, parentPath)...)
		logging.FrameworkTrace(ctx, "Called provider defined Type Validate")

		if diags.HasError() {
			return nil, diags
		}
	}

	return d.SetAtPathTransformFunc(ctx, parentPath, parentValue, diags)
}
