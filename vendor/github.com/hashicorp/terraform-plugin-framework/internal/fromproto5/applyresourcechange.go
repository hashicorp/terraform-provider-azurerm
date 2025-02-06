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

// ApplyResourceChangeRequest returns the *fwserver.ApplyResourceChangeRequest
// equivalent of a *tfprotov5.ApplyResourceChangeRequest.
func ApplyResourceChangeRequest(ctx context.Context, proto5 *tfprotov5.ApplyResourceChangeRequest, resource resource.Resource, resourceSchema fwschema.Schema, providerMetaSchema fwschema.Schema) (*fwserver.ApplyResourceChangeRequest, diag.Diagnostics) {
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

	fw := &fwserver.ApplyResourceChangeRequest{
		ResourceSchema: resourceSchema,
		Resource:       resource,
	}

	config, configDiags := Config(ctx, proto5.Config, resourceSchema)

	diags.Append(configDiags...)

	fw.Config = config

	plannedState, plannedStateDiags := Plan(ctx, proto5.PlannedState, resourceSchema)

	diags.Append(plannedStateDiags...)

	fw.PlannedState = plannedState

	priorState, priorStateDiags := State(ctx, proto5.PriorState, resourceSchema)

	diags.Append(priorStateDiags...)

	fw.PriorState = priorState

	providerMeta, providerMetaDiags := ProviderMeta(ctx, proto5.ProviderMeta, providerMetaSchema)

	diags.Append(providerMetaDiags...)

	fw.ProviderMeta = providerMeta

	privateData, privateDataDiags := privatestate.NewData(ctx, proto5.PlannedPrivate)

	diags.Append(privateDataDiags...)

	fw.PlannedPrivate = privateData

	return fw, diags
}
