// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// PlanResourceChangeRequest returns the *fwserver.PlanResourceChangeRequest
// equivalent of a *tfprotov5.PlanResourceChangeRequest.
func PlanResourceChangeRequest(ctx context.Context, proto5 *tfprotov5.PlanResourceChangeRequest, resource resource.Resource, resourceSchema fwschema.Schema, providerMetaSchema fwschema.Schema) (*fwserver.PlanResourceChangeRequest, diag.Diagnostics) {
	if proto5 == nil {
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

	fw := &fwserver.PlanResourceChangeRequest{
		ResourceSchema: resourceSchema,
		Resource:       resource,
	}

	config, configDiags := Config(ctx, proto5.Config, resourceSchema)

	diags.Append(configDiags...)

	fw.Config = config

	priorState, priorStateDiags := State(ctx, proto5.PriorState, resourceSchema)

	diags.Append(priorStateDiags...)

	fw.PriorState = priorState

	proposedNewState, proposedNewStateDiags := Plan(ctx, proto5.ProposedNewState, resourceSchema)

	diags.Append(proposedNewStateDiags...)

	fw.ProposedNewState = proposedNewState

	providerMeta, providerMetaDiags := ProviderMeta(ctx, proto5.ProviderMeta, providerMetaSchema)

	diags.Append(providerMetaDiags...)

	fw.ProviderMeta = providerMeta

	privateData, privateDataDiags := privatestate.NewData(ctx, proto5.PriorPrivate)

	diags.Append(privateDataDiags...)

	fw.PriorPrivate = privateData

	return fw, diags
}
