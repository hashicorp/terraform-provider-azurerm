// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateListResourceConfigRequest returns the *fwserver.ValidateListResourceConfigRequest
// equivalent of a *tfprotov6.ValidateListResourceConfigRequest.
func ValidateListResourceConfigRequest(ctx context.Context, proto6 *tfprotov6.ValidateListResourceConfigRequest, listResource list.ListResource, listResourceSchema fwschema.Schema) (*fwserver.ValidateListResourceConfigRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateListResourceConfigRequest{}

	config, diags := Config(ctx, proto6.Config, listResourceSchema)

	fw.Config = config
	fw.ListResource = listResource

	return fw, diags
}
