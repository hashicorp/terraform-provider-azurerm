// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateListResourceConfigRequest returns the *fwserver.ValidateListResourceConfigRequest
// equivalent of a *tfprotov5.ValidateListResourceConfigRequest.
func ValidateListResourceConfigRequest(ctx context.Context, proto5 *tfprotov5.ValidateListResourceConfigRequest, listResource list.ListResource, listResourceSchema fwschema.Schema) (*fwserver.ValidateListResourceConfigRequest, diag.Diagnostics) {
	if proto5 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateListResourceConfigRequest{}

	config, diags := Config(ctx, proto5.Config, listResourceSchema)

	fw.Config = config
	fw.ListResource = listResource

	return fw, diags
}
