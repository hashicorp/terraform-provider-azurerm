// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ResourceIdentity returns the *tfprotov5.ResourceIdentityData for a *tfsdk.ResourceIdentity.
func ResourceIdentity(ctx context.Context, fw *tfsdk.ResourceIdentity) (*tfprotov5.ResourceIdentityData, diag.Diagnostics) {
	if fw == nil {
		return nil, nil
	}

	identitySchemaData := &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionResourceIdentity,
		Schema:         fw.Schema,
		TerraformValue: fw.Raw,
	}

	identityData, diags := DynamicValue(ctx, identitySchemaData)
	if diags.HasError() {
		return nil, diags
	}

	return &tfprotov5.ResourceIdentityData{
		IdentityData: identityData,
	}, nil
}
