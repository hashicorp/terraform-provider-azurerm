// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// PathExists returns true if the path can be reached. The value at the path
// may be null or unknown.
func (d Data) PathExists(ctx context.Context, path path.Path) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	tftypesPath, tftypesPathDiags := totftypes.AttributePath(ctx, path)

	diags.Append(tftypesPathDiags...)

	if diags.HasError() {
		return false, diags
	}

	_, remaining, err := tftypes.WalkAttributePath(d.TerraformValue, tftypesPath)

	if err != nil {
		if errors.Is(err, tftypes.ErrInvalidStep) {
			return false, diags
		}

		diags.AddAttributeError(
			path,
			d.Description.Title()+" Read Error",
			"An unexpected error was encountered trying to read an attribute from the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Cannot walk attribute path in %s: %s", d.Description, err),
		)
		return false, diags
	}

	return len(remaining.Steps()) == 0, diags
}
