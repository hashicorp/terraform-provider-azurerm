// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// EphemeralResultData returns the *tfprotov6.DynamicValue for a *tfsdk.EphemeralResultData.
func EphemeralResultData(ctx context.Context, fw *tfsdk.EphemeralResultData) (*tfprotov6.DynamicValue, diag.Diagnostics) {
	if fw == nil {
		return nil, nil
	}

	data := &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionEphemeralResultData,
		Schema:         fw.Schema,
		TerraformValue: fw.Raw,
	}

	return DynamicValue(ctx, data)
}
