// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromtftypes"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// TransformDefaults walks the schema and applies schema defined default values
// when configRaw contains a null value at the same path.
func (d *Data) TransformDefaults(ctx context.Context, configRaw tftypes.Value) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	configData := Data{
		Description:    DataDescriptionConfiguration,
		Schema:         d.Schema,
		TerraformValue: configRaw,
	}

	d.TerraformValue, err = tftypes.Transform(d.TerraformValue, func(tfTypePath *tftypes.AttributePath, tfTypeValue tftypes.Value) (tftypes.Value, error) {
		// Skip the root of the data, only applying defaults to attributes
		if len(tfTypePath.Steps()) < 1 {
			return tfTypeValue, nil
		}

		attrAtPath, err := d.Schema.AttributeAtTerraformPath(ctx, tfTypePath)

		if err != nil {
			if errors.Is(err, fwschema.ErrPathInsideAtomicAttribute) {
				// ignore attributes/elements inside schema.Attributes, they have no schema of their own
				logging.FrameworkTrace(ctx, "attribute is a non-schema attribute, not setting default")
				return tfTypeValue, nil
			}

			if errors.Is(err, fwschema.ErrPathIsBlock) {
				// ignore blocks, they do not have a computed field
				logging.FrameworkTrace(ctx, "attribute is a block, not setting default")
				return tfTypeValue, nil
			}

			if errors.Is(err, fwschema.ErrPathInsideDynamicAttribute) {
				// ignore attributes/elements inside schema.DynamicAttribute, they have no schema of their own
				logging.FrameworkTrace(ctx, "attribute is inside of a dynamic attribute, not setting default")
				return tfTypeValue, nil
			}

			return tftypes.Value{}, fmt.Errorf("couldn't find attribute in resource schema: %w", err)
		}

		fwPath, fwPathDiags := fromtftypes.AttributePath(ctx, tfTypePath, d.Schema)

		diags.Append(fwPathDiags...)

		// Do not transform if path cannot be converted.
		// Checking against fwPathDiags will capture all errors.
		if fwPathDiags.HasError() {
			return tfTypeValue, nil
		}

		configValue, configValueDiags := configData.ValueAtPath(ctx, fwPath)

		diags.Append(configValueDiags...)

		// Do not transform if rawConfig value cannot be retrieved.
		if configValueDiags.HasError() {
			return tfTypeValue, nil
		}

		// Do not transform if rawConfig value is not null.
		if !configValue.IsNull() {
			// Dynamic values need to perform more logic to check the config value for null-ness
			dynValuable, ok := configValue.(basetypes.DynamicValuable)
			if !ok {
				return tfTypeValue, nil
			}

			dynConfigVal, dynDiags := dynValuable.ToDynamicValue(ctx)
			if dynDiags.HasError() {
				return tfTypeValue, nil
			}

			// For dynamic values, it's possible to be known when only the type is known.
			// The underlying value can still be null, so check for that here
			if !dynConfigVal.IsUnderlyingValueNull() {
				return tfTypeValue, nil
			}
		}

		switch a := attrAtPath.(type) {
		case fwschema.AttributeWithBoolDefaultValue:
			defaultValue := a.BoolDefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.BoolRequest{
				Path: fwPath,
			}
			resp := defaults.BoolResponse{}

			defaultValue.DefaultBool(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithFloat64DefaultValue:
			defaultValue := a.Float64DefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.Float64Request{
				Path: fwPath,
			}
			resp := defaults.Float64Response{}

			defaultValue.DefaultFloat64(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithInt64DefaultValue:
			defaultValue := a.Int64DefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.Int64Request{
				Path: fwPath,
			}
			resp := defaults.Int64Response{}

			defaultValue.DefaultInt64(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithListDefaultValue:
			defaultValue := a.ListDefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.ListRequest{
				Path: fwPath,
			}
			resp := defaults.ListResponse{}

			defaultValue.DefaultList(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			if resp.PlanValue.ElementType(ctx) == nil {
				logging.FrameworkWarn(ctx, "attribute default declared, but returned no value")

				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithMapDefaultValue:
			defaultValue := a.MapDefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}
			req := defaults.MapRequest{
				Path: fwPath,
			}
			resp := defaults.MapResponse{}

			defaultValue.DefaultMap(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			if resp.PlanValue.ElementType(ctx) == nil {
				logging.FrameworkWarn(ctx, "attribute default declared, but returned no value")

				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithNumberDefaultValue:
			defaultValue := a.NumberDefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.NumberRequest{
				Path: fwPath,
			}
			resp := defaults.NumberResponse{}

			defaultValue.DefaultNumber(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithObjectDefaultValue:
			defaultValue := a.ObjectDefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.ObjectRequest{
				Path: fwPath,
			}
			resp := defaults.ObjectResponse{}

			defaultValue.DefaultObject(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithSetDefaultValue:
			defaultValue := a.SetDefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.SetRequest{
				Path: fwPath,
			}
			resp := defaults.SetResponse{}

			defaultValue.DefaultSet(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			if resp.PlanValue.ElementType(ctx) == nil {
				logging.FrameworkWarn(ctx, "attribute default declared, but returned no value")

				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithStringDefaultValue:
			defaultValue := a.StringDefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.StringRequest{
				Path: fwPath,
			}
			resp := defaults.StringResponse{}

			defaultValue.DefaultString(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		case fwschema.AttributeWithDynamicDefaultValue:
			defaultValue := a.DynamicDefaultValue()

			if defaultValue == nil {
				return tfTypeValue, nil
			}

			req := defaults.DynamicRequest{
				Path: fwPath,
			}
			resp := defaults.DynamicResponse{}

			defaultValue.DefaultDynamic(ctx, req, &resp)

			diags.Append(resp.Diagnostics...)

			if resp.Diagnostics.HasError() {
				return tfTypeValue, nil
			}

			logging.FrameworkTrace(ctx, fmt.Sprintf("setting attribute %s to default value: %s", fwPath, resp.PlanValue))

			return resp.PlanValue.ToTerraformValue(ctx)
		}

		return tfTypeValue, nil
	})

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/930
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic(
			"Error Handling Schema Defaults",
			"An unexpected error occurred while handling schema default values. "+
				"Please report the following to the provider developer:\n\n"+
				"Error: "+err.Error(),
		))
	}

	return diags
}
