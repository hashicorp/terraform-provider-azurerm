// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// DynamicValue returns the *tfprotov6.DynamicValue for a given
// fwschemadata.Data.
//
// If necessary, the underlying data is modified to convert list and set block
// values from a null collection to an empty collection. This is to prevent
// developers from needing to understand Terraform's differences between
// block and attribute values where blocks are technically never null, but from
// a developer perspective this distinction introduces unnecessary complexity.
func DynamicValue(ctx context.Context, data *fwschemadata.Data) (*tfprotov6.DynamicValue, diag.Diagnostics) {
	if data == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Prevent Terraform core errors for null list/set blocks.
	diags.Append(data.ReifyNullCollectionBlocks(ctx)...)

	proto6, err := tfprotov6.NewDynamicValue(data.Schema.Type().TerraformType(ctx), data.TerraformValue)

	if err != nil {
		diags.AddError(
			"Unable to Convert "+data.Description.Title(),
			"An unexpected error was encountered when converting the "+data.Description.String()+" to the protocol type. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+
				"Unable to create DynamicValue: "+err.Error(),
		)

		return nil, diags
	}

	return &proto6, nil
}
