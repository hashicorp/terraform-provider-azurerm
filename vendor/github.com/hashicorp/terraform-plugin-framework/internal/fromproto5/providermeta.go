// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ProviderMeta returns the *tfsdk.Config for a *tfprotov5.DynamicValue and
// fwschema.Schema. This data handling is different than Config to simplify
// implementors, in that:
//
//   - Missing Schema will return nil, rather than an error
//   - Missing DynamicValue will return nil typed Value, rather than an error
func ProviderMeta(ctx context.Context, proto5DynamicValue *tfprotov5.DynamicValue, schema fwschema.Schema) (*tfsdk.Config, diag.Diagnostics) {
	if schema == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	fw := &tfsdk.Config{
		Raw:    tftypes.NewValue(schema.Type().TerraformType(ctx), nil),
		Schema: schema,
	}

	if proto5DynamicValue == nil {
		return fw, nil
	}

	proto5Value, err := proto5DynamicValue.Unmarshal(schema.Type().TerraformType(ctx))

	if err != nil {
		diags.AddError(
			"Unable to Convert Provider Meta Configuration",
			"An unexpected error was encountered when converting the provider meta configuration from the protocol type. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+err.Error(),
		)

		return nil, diags
	}

	fw.Raw = proto5Value

	return fw, nil
}
