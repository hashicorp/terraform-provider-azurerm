// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// GenerateResourceConfigRequest returns the *fwserver.GenerateResourceConfigRequest
// equivalent of a *tfprotov6.GenerateResourceConfigRequest.
func GenerateResourceConfigRequest(ctx context.Context, proto6 *tfprotov6.GenerateResourceConfigRequest, resourceSchema fwschema.Schema) (*fwserver.GenerateResourceConfigRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Panic prevention here to simplify the calling implementations.
	// This should not happen, but just in case.
	if resourceSchema == nil {
		diags.AddError(
			"Missing Resource Schema",
			"An unexpected error was encountered when handling the request. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+
				"Missing schema.",
		)

		return nil, diags
	}

	state, stateDiags := State(ctx, proto6.State, resourceSchema)

	diags.Append(stateDiags...)

	fw := &fwserver.GenerateResourceConfigRequest{
		ResourceSchema: resourceSchema,
		State:          state,
	}

	return fw, diags
}
