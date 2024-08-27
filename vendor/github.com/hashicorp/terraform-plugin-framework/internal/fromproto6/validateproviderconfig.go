// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateProviderConfigRequest returns the *fwserver.ValidateProviderConfigRequest
// equivalent of a *tfprotov6.ValidateProviderConfigRequest.
func ValidateProviderConfigRequest(ctx context.Context, proto6 *tfprotov6.ValidateProviderConfigRequest, providerSchema fwschema.Schema) (*fwserver.ValidateProviderConfigRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateProviderConfigRequest{}

	config, diags := Config(ctx, proto6.Config, providerSchema)

	fw.Config = config

	return fw, diags
}
