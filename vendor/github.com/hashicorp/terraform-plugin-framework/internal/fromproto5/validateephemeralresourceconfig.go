// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateEphemeralResourceConfigRequest returns the *fwserver.ValidateEphemeralResourceConfigRequest
// equivalent of a *tfprotov5.ValidateEphemeralResourceConfigRequest.
func ValidateEphemeralResourceConfigRequest(ctx context.Context, proto5 *tfprotov5.ValidateEphemeralResourceConfigRequest, ephemeralResource ephemeral.EphemeralResource, ephemeralResourceSchema fwschema.Schema) (*fwserver.ValidateEphemeralResourceConfigRequest, diag.Diagnostics) {
	if proto5 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateEphemeralResourceConfigRequest{}

	config, diags := Config(ctx, proto5.Config, ephemeralResourceSchema)

	fw.Config = config
	fw.EphemeralResource = ephemeralResource

	return fw, diags
}
