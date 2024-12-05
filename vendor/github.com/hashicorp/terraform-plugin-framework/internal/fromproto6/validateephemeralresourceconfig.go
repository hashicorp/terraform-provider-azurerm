// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateEphemeralResourceConfigRequest returns the *fwserver.ValidateEphemeralResourceConfigRequest
// equivalent of a *tfprotov6.ValidateEphemeralResourceConfigRequest.
func ValidateEphemeralResourceConfigRequest(ctx context.Context, proto6 *tfprotov6.ValidateEphemeralResourceConfigRequest, ephemeralResource ephemeral.EphemeralResource, ephemeralResourceSchema fwschema.Schema) (*fwserver.ValidateEphemeralResourceConfigRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateEphemeralResourceConfigRequest{}

	config, diags := Config(ctx, proto6.Config, ephemeralResourceSchema)

	fw.Config = config
	fw.EphemeralResource = ephemeralResource

	return fw, diags
}
