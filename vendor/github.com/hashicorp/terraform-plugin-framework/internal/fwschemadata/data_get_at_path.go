// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// GetAtPath retrieves the attribute found at `path` and populates the
// `target` with the value.
func (d Data) GetAtPath(ctx context.Context, schemaPath path.Path, target any) diag.Diagnostics {
	ctx = logging.FrameworkWithAttributePath(ctx, schemaPath.String())

	attrValue, diags := d.ValueAtPath(ctx, schemaPath)

	if diags.HasError() {
		return diags
	}

	if attrValue == nil {
		diags.AddAttributeError(
			schemaPath,
			d.Description.Title()+" Read Error",
			"An unexpected error was encountered trying to read an attribute from the "+d.Description.String()+". This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Missing attribute value, however no error was returned. Preventing the panic from this situation.",
		)
		return diags
	}

	if reflect.IsGenericAttrValue(ctx, target) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `reflect.IsGenericAttrValue` function
		*(target.(*attr.Value)) = attrValue
		return nil
	}

	raw, err := attrValue.ToTerraformValue(ctx)

	if err != nil {
		diags.AddAttributeError(
			schemaPath,
			d.Description.Title()+" Value Conversion Error",
			fmt.Sprintf("An unexpected error was encountered converting a %T to its equivalent Terraform representation. This is always a bug in the provider.\n\n"+
				"Error: %s", attrValue, err),
		)
		return diags
	}

	reflectDiags := reflect.Into(ctx, attrValue.Type(ctx), raw, target, reflect.Options{}, schemaPath)

	diags.Append(reflectDiags...)

	return diags
}
