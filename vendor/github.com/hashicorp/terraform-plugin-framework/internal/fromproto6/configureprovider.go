// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ConfigureProviderRequest returns the *fwserver.ConfigureProviderRequest
// equivalent of a *tfprotov6.ConfigureProviderRequest.
func ConfigureProviderRequest(ctx context.Context, proto6 *tfprotov6.ConfigureProviderRequest, providerSchema fwschema.Schema) (*provider.ConfigureRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	fw := &provider.ConfigureRequest{
		TerraformVersion: proto6.TerraformVersion,
	}

	config, diags := Config(ctx, proto6.Config, providerSchema)

	if config != nil {
		fw.Config = *config
	}

	return fw, diags
}
