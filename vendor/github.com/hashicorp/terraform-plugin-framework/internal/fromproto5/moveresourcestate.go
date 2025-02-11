// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// MoveResourceStateRequest returns the *fwserver.MoveResourceStateRequest
// equivalent of a *tfprotov5.MoveResourceStateRequest.
func MoveResourceStateRequest(ctx context.Context, proto5 *tfprotov5.MoveResourceStateRequest, resource resource.Resource, resourceSchema fwschema.Schema) (*fwserver.MoveResourceStateRequest, diag.Diagnostics) {
	if proto5 == nil {
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
		SourceProviderAddress: proto5.SourceProviderAddress,
		SourceRawState:        (*tfprotov6.RawState)(proto5.SourceState),
		SourceSchemaVersion:   proto5.SourceSchemaVersion,
		SourceTypeName:        proto5.SourceTypeName,
		TargetResource:        resource,
		TargetResourceSchema:  resourceSchema,
		TargetTypeName:        proto5.TargetTypeName,
	}

	sourcePrivate, sourcePrivateDiags := privatestate.NewData(ctx, proto5.SourcePrivate)

	diags.Append(sourcePrivateDiags...)

	fw.SourcePrivate = sourcePrivate

	return fw, diags
}
