// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// UpgradeResourceStateRequest returns the *fwserver.UpgradeResourceStateRequest
// equivalent of a *tfprotov6.UpgradeResourceStateRequest.
func UpgradeResourceStateRequest(ctx context.Context, proto6 *tfprotov6.UpgradeResourceStateRequest, resource resource.Resource, resourceSchema fwschema.Schema) (*fwserver.UpgradeResourceStateRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Panic prevention here to simplify the calling implementations.
	// This should not happen, but just in case.
	if resourceSchema == nil {
		diags.AddError(
			"Unable to Create Empty State",
			"An unexpected error was encountered when creating the empty state. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+
				"Missing schema.",
		)

		return nil, diags
	}

	fw := &fwserver.UpgradeResourceStateRequest{
		RawState:       proto6.RawState,
		ResourceSchema: resourceSchema,
		Resource:       resource,
		Version:        proto6.Version,
	}

	return fw, diags
}
