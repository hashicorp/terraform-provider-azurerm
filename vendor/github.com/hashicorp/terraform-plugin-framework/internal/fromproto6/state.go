// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// State returns the *tfsdk.State for a *tfprotov6.DynamicValue and
// fwschema.Schema.
func State(ctx context.Context, proto6DynamicValue *tfprotov6.DynamicValue, schema fwschema.Schema) (*tfsdk.State, diag.Diagnostics) {
	if proto6DynamicValue == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Panic prevention here to simplify the calling implementations.
	// This should not happen, but just in case.
	if schema == nil {
		diags.AddError(
			"Unable to Convert State",
			"An unexpected error was encountered when converting the state from the protocol type. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+
				"Missing schema.",
		)

		return nil, diags
	}

	data, dynamicValueDiags := DynamicValue(ctx, proto6DynamicValue, schema, fwschemadata.DataDescriptionState)

	diags.Append(dynamicValueDiags...)

	if diags.HasError() {
		return nil, diags
	}

	fw := &tfsdk.State{
		Raw:    data.TerraformValue,
		Schema: schema,
	}

	return fw, diags
}
