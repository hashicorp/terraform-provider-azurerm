// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromtftypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// PathMatches returns all matching path.Paths from the given path.Expression.
//
// If a parent path is null or unknown, which would prevent a full expression
// from matching, the parent path is returned rather than no match to prevent
// false positives.
func (d Data) PathMatches(ctx context.Context, pathExpr path.Expression) (path.Paths, diag.Diagnostics) {
	var diags diag.Diagnostics
	var paths path.Paths

	if !d.ValidPathExpression(ctx, pathExpr) {
		diags.AddError(
			"Invalid Path Expression for Schema",
			"The Terraform Provider unexpectedly provided a path expression that does not match the current schema. "+
				"This can happen if the path expression does not correctly follow the schema in structure or types. "+
				"Please report this to the provider developers.\n\n"+
				"Path Expression: "+pathExpr.String(),
		)

		return paths, diags
	}

	_ = tftypes.Walk(d.TerraformValue, func(tfTypePath *tftypes.AttributePath, tfTypeValue tftypes.Value) (bool, error) {
		fwPath, fwPathDiags := fromtftypes.AttributePath(ctx, tfTypePath, d.Schema)

		diags.Append(fwPathDiags...)

		if diags.HasError() {
			// If there was an error with conversion of the path at this level,
			// no need to traverse further since a deeper path will error.
			return false, nil
		}

		if pathExpr.Matches(fwPath) {
			paths.Append(fwPath)

			// If we matched, there is no need to traverse further since a
			// deeper path will never match.
			return false, nil
		}

		// If current path cannot be parent path, there is no need to traverse
		// further since a deeper path will never match.
		if !pathExpr.MatchesParent(fwPath) {
			return false, nil
		}

		// If value at current path (now known to be a parent path of the
		// expression) is null or unknown, return it as a valid path match
		// since Walk will stop traversing deeper anyways and we want
		// consumers to know about the path with the null or unknown value.
		//
		// This behavior may be confusing for consumers as fetching the value
		// at this parent path will return a potentially unexpected type,
		// however this is an implementation tradeoff to prevent false
		// positives of missing null or unknown values.
		if tfTypeValue.IsNull() || !tfTypeValue.IsKnown() {
			paths.Append(fwPath)

			return false, nil
		}

		return true, nil
	})

	return paths, diags
}
