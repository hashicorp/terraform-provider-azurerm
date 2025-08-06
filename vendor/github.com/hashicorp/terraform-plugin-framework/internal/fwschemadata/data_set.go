// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Set replaces the entire value. The value should be a struct whose fields
// have one of the attr.Value types. Each field must have the tfsdk field tag.
func (d *Data) Set(ctx context.Context, val any) diag.Diagnostics {
	attrValue, diags := reflect.FromValue(ctx, d.Schema.Type(), val, path.Empty())

	if diags.HasError() {
		return diags
	}

	tfValue, err := attrValue.ToTerraformValue(ctx)

	if err != nil {
		diags.AddError(
			d.Description.Title()+" Write Error",
			"An unexpected error was encountered trying to write the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Error: Unable to run ToTerraformValue on new value: %s", err),
		)
		return diags
	}

	d.TerraformValue = tfValue

	return diags
}
