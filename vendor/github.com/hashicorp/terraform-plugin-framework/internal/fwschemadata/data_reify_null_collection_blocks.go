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

// ReifyNullCollectionBlocks converts list and set block null values to empty
// values. This is the reverse conversion of NullifyCollectionBlocks.
func (d *Data) ReifyNullCollectionBlocks(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	blockPathExpressions := fwschema.SchemaBlockPathExpressions(ctx, d.Schema)

	// Errors are handled as richer diag.Diagnostics instead.
	d.TerraformValue, _ = tftypes.Transform(d.TerraformValue, func(tfTypePath *tftypes.AttributePath, tfTypeValue tftypes.Value) (tftypes.Value, error) {
		// Skip the root of the data
		if len(tfTypePath.Steps()) < 1 {
			return tfTypeValue, nil
		}

		// Only transform null values.
		if !tfTypeValue.IsNull() {
			return tfTypeValue, nil
		}

		_, err := d.Schema.AttributeAtTerraformPath(ctx, tfTypePath)
		if err != nil {
			if errors.Is(err, fwschema.ErrPathInsideDynamicAttribute) {
				// ignore attributes/elements inside schema.DynamicAttribute
				logging.FrameworkTrace(ctx, "attribute is inside of a dynamic attribute, skipping reify null collection blocks")
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

		// Transform to empty value.
		switch tfTypeValue.Type().(type) {
		case tftypes.List, tftypes.Set:
			logging.FrameworkTrace(ctx, "Transforming null block to empty block", map[string]any{
				logging.KeyAttributePath: fwPath.String(),
				logging.KeyDescription:   d.Description.String(),
			})
			return tftypes.NewValue(tfTypeValue.Type(), []tftypes.Value{}), nil
		default:
			return tfTypeValue, nil
		}
	})

	return diags
}
