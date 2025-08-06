// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// RenewEphemeralResourceRequest returns the *fwserver.RenewEphemeralResourceRequest
// equivalent of a *tfprotov6.RenewEphemeralResourceRequest.
func RenewEphemeralResourceRequest(ctx context.Context, proto6 *tfprotov6.RenewEphemeralResourceRequest, ephemeralResource ephemeral.EphemeralResource, ephemeralResourceSchema fwschema.Schema) (*fwserver.RenewEphemeralResourceRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Panic prevention here to simplify the calling implementations.
	// This should not happen, but just in case.
	if ephemeralResourceSchema == nil {
		diags.AddError(
			"Missing EphemeralResource Schema",
			"An unexpected error was encountered when handling the request. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+
				"Missing schema.",
		)

		return nil, diags
	}

	fw := &fwserver.RenewEphemeralResourceRequest{
		EphemeralResource:       ephemeralResource,
		EphemeralResourceSchema: ephemeralResourceSchema,
	}

	privateData, privateDataDiags := privatestate.NewData(ctx, proto6.Private)

	diags.Append(privateDataDiags...)

	fw.Private = privateData

	return fw, diags
}
