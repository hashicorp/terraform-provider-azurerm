// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// MoveResourceStateRequest returns the *fwserver.MoveResourceStateRequest
// equivalent of a *tfprotov6.MoveResourceStateRequest.
func MoveResourceStateRequest(ctx context.Context, proto6 *tfprotov6.MoveResourceStateRequest, resource resource.Resource, resourceSchema fwschema.Schema) (*fwserver.MoveResourceStateRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Panic prevention here to simplify the calling implementations.
	// This should not happen, but just in case.
	if resourceSchema == nil {
		diags.AddError(
			"Framework Implementation Error",
			"An unexpected issue was encountered when converting the MoveResourceState RPC request information from the protocol type to the framework type. "+
				"The resource schema was missing. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.",
		)

		return nil, diags
	}

	fw := &fwserver.MoveResourceStateRequest{
		SourceProviderAddress: proto6.SourceProviderAddress,
		SourceRawState:        proto6.SourceState,
		SourceSchemaVersion:   proto6.SourceSchemaVersion,
		SourceTypeName:        proto6.SourceTypeName,
		TargetResource:        resource,
		TargetResourceSchema:  resourceSchema,
		TargetTypeName:        proto6.TargetTypeName,
	}

	sourcePrivate, sourcePrivateDiags := privatestate.NewData(ctx, proto6.SourcePrivate)

	diags.Append(sourcePrivateDiags...)

	fw.SourcePrivate = sourcePrivate

	return fw, diags
}
