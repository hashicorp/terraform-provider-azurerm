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

// ResourceIdentity returns the *tfprotov6.ResourceIdentityData for a *tfsdk.ResourceIdentity.
func ResourceIdentity(ctx context.Context, fw *tfsdk.ResourceIdentity) (*tfprotov6.ResourceIdentityData, diag.Diagnostics) {
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

	return &tfprotov6.ResourceIdentityData{
		IdentityData: identityData,
	}, nil
}
