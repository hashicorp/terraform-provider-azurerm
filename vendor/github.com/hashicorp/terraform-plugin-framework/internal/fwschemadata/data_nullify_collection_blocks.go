// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromtftypes"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// NullifyCollectionBlocks converts list and set block empty values to null
// values. The reverse conversion is ReifyNullCollectionBlocks.
func (d *Data) NullifyCollectionBlocks(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	blockPathExpressions := fwschema.SchemaBlockPathExpressions(ctx, d.Schema)

	// Errors are handled as richer diag.Diagnostics instead.
	d.TerraformValue, _ = tftypes.Transform(d.TerraformValue, func(tfTypePath *tftypes.AttributePath, tfTypeValue tftypes.Value) (tftypes.Value, error) {
		// Skip the root of the data
		if len(tfTypePath.Steps()) < 1 {
			return tfTypeValue, nil
		}

		// Do not transform if value is already null or is not fully known.
		if tfTypeValue.IsNull() || !tfTypeValue.IsFullyKnown() {
			return tfTypeValue, nil
		}

		_, err := d.Schema.AttributeAtTerraformPath(ctx, tfTypePath)
		if err != nil {
			if errors.Is(err, fwschema.ErrPathInsideDynamicAttribute) {
				// ignore attributes/elements inside schema.DynamicAttribute
				logging.FrameworkTrace(ctx, "attribute is inside of a dynamic attribute, skipping nullify collection blocks")
				return tfTypeValue, nil
			}
		}

		fwPath, fwPathDiags := fromtftypes.AttributePath(ctx, tfTypePath, d.Schema)

		diags.Append(fwPathDiags...)

		// Do not transform if path cannot be converted.
		// Checking against fwPathDiags will capture all errors.
		if fwPathDiags.HasError() {
			return tfTypeValue, nil
		}

		// Do not transform if path is not a block.
		if !blockPathExpressions.Matches(fwPath) {
			return tfTypeValue, nil
		}

		var elements []tftypes.Value

		switch tfTypeValue.Type().(type) {
		case tftypes.List, tftypes.Set:
			err := tfTypeValue.As(&elements)

			// If this occurs, it likely is an upstream issue in Terraform
			// or terraform-plugin-go.
			if err != nil {
				diags.AddAttributeError(
					fwPath,
					d.Description.Title()+" Data Transformation Error",
					"An unexpected error occurred while transforming "+d.Description.String()+" data. "+
						"This is always an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
						"Path: "+fwPath.String()+"\n"+
						"Error: (tftypes.Value).As() error: "+err.Error(),
				)

				return tfTypeValue, nil //nolint:nilerr // Using richer diag.Diagnostics instead.
			}
		default:
			return tfTypeValue, nil
		}

		// Do not transform if there are any elements.
		if len(elements) > 0 {
			return tfTypeValue, nil
		}

		// Transform to null value.
		logging.FrameworkTrace(ctx, "Transforming empty block to null block", map[string]any{
			logging.KeyAttributePath: fwPath.String(),
			logging.KeyDescription:   d.Description.String(),
		})
		return tftypes.NewValue(tfTypeValue.Type(), nil), nil
	})

	return diags
}
