// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// PrepareProviderConfigRequest returns the *fwserver.ValidateProviderConfigRequest
// equivalent of a *tfprotov5.PrepareProviderConfigRequest.
func PrepareProviderConfigRequest(ctx context.Context, proto5 *tfprotov5.PrepareProviderConfigRequest, providerSchema fwschema.Schema) (*fwserver.ValidateProviderConfigRequest, diag.Diagnostics) {
	if proto5 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateProviderConfigRequest{}

	config, diags := Config(ctx, proto5.Config, providerSchema)

	fw.Config = config

	return fw, diags
}
