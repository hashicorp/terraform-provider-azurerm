// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Set replaces the entire value. The value can be a tftypes.Value or a struct whose fields
// have one of the attr.Value types. Each field must have the tfsdk field tag.
func (d *Data) Set(ctx context.Context, val any) diag.Diagnostics {
	var diags diag.Diagnostics

	if v, ok := val.(tftypes.Value); ok {
		objType := d.Schema.Type().TerraformType(ctx)

		if !objType.Equal(v.Type()) {
			diags.AddError(
				d.Description.Title()+" Write Error",
				"An unexpected error was encountered trying to write the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Error: Type mismatch between provided value and type of %s, expected %+v, got %+v", d.Description.String(), objType.String(), v.Type().String()),
			)
			return diags

		}
		d.TerraformValue = v

		return diags
	}

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
