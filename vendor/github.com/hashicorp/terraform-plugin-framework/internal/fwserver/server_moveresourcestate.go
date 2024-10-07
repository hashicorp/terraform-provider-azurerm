// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// MoveResourceStateRequest is the framework server request for the
// MoveResourceState RPC.
type MoveResourceStateRequest struct {
	// SourcePrivate is the private state of the source resource as given by
	// Terraform across the protocol.
	SourcePrivate *privatestate.Data

	// SourceProviderAddress is the address of the source provider as given by
	// Terraform across the protocol.
	SourceProviderAddress string

	// SourceSchemaVersion is the version of the source resource schema as given
	// by Terraform across the protocol.
	SourceSchemaVersion int64

	// SourceRawState is the raw state of the source resource as given by
	// Terraform across the protocol.
	//
	// Using the tfprotov6 type here was a pragmatic effort decision around when
	// the framework introduced compatibility promises. This type was chosen as
	// it was readily available and trivial to convert between tfprotov5.
	//
	// Using a terraform-plugin-go type is not ideal for the framework as almost
	// all terraform-plugin-go types have framework abstractions, but if there
	// is ever a time where it makes sense to re-evaluate this decision, such as
	// a major version bump, it could be changed then.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/340
	SourceRawState *tfprotov6.RawState

	// SourceTypeName is the type name of the source resource as given by
	// Terraform across the protocol.
	SourceTypeName string

	// TargetResource is the provider-defined resource implementation as
	// determined by the framework looking up the resource name from the
	// provider.Provider implementation Resources method defined by the
	// provider developer.
	TargetResource resource.Resource

	// TargetResourceSchema is the evaluated schema definition of the target
	// resource as determined by the framework calling the resource.Resource
	// implementation Schema method defined by the provider developer.
	TargetResourceSchema fwschema.Schema

	// TargetTypeName is the type name of the target resource as given by
	// Terraform across the protocol.
	TargetTypeName string
}

// MoveResourceStateResponse is the framework server response for the
// MoveResourceState RPC.
type MoveResourceStateResponse struct {
	Diagnostics   diag.Diagnostics
	TargetPrivate *privatestate.Data
	TargetState   *tfsdk.State
}

// MoveResourceState implements the framework server MoveResourceState RPC.
func (s *Server) MoveResourceState(ctx context.Context, req *MoveResourceStateRequest, resp *MoveResourceStateResponse) {
	if req == nil {
		return
	}

	if req.SourceRawState == nil {
		resp.Diagnostics.AddError(
			"Missing Source Resource State",
			"The source resource state was not provided to the provider for the MoveResourceState operation. "+
				"This is always an issue in Terraform and should be reported to the Terraform maintainers.",
		)

		return
	}

	resourceWithMoveState, ok := req.TargetResource.(resource.ResourceWithMoveState)

	if !ok {
		resp.Diagnostics.AddError(
			"Unable to Move Resource State",
			"The target resource implementation does not include move resource state support. "+
				"The resource implementation can be updated by the provider developers to include this support with the ResourceWithMoveState interface.\n\n"+
				"Source Provider Address: "+req.SourceProviderAddress+"\n"+
				"Source Resource Type: "+req.SourceTypeName+"\n"+
				"Source Resource Schema Version: "+strconv.FormatInt(req.SourceSchemaVersion, 10)+"\n"+
				"Target Resource Type: "+req.TargetTypeName,
		)

		return
	}

	logging.FrameworkTrace(ctx, "Resource implements ResourceWithMoveState")

	logging.FrameworkTrace(ctx, "Calling provider defined Resource MoveState")
	resourceStateMovers := resourceWithMoveState.MoveState(ctx)
	logging.FrameworkTrace(ctx, "Called provider defined Resource MoveState")

	sourcePrivate := privatestate.EmptyProviderData(ctx)

	if req.SourcePrivate != nil && req.SourcePrivate.Provider != nil {
		sourcePrivate = req.SourcePrivate.Provider
	}

	if resp.TargetPrivate == nil {
		resp.TargetPrivate = privatestate.EmptyData(ctx)
	}

	for _, resourceStateMover := range resourceStateMovers {
		moveStateReq := resource.MoveStateRequest{
			SourcePrivate:         sourcePrivate,
			SourceProviderAddress: req.SourceProviderAddress,
			SourceRawState:        req.SourceRawState,
			SourceSchemaVersion:   req.SourceSchemaVersion,
			SourceTypeName:        req.SourceTypeName,
		}
		moveStateResp := resource.MoveStateResponse{
			TargetPrivate: privatestate.EmptyProviderData(ctx),
			TargetState: tfsdk.State{
				Schema: req.TargetResourceSchema,
				Raw:    tftypes.NewValue(req.TargetResourceSchema.Type().TerraformType(ctx), nil),
			},
		}

		if resourceStateMover.SourceSchema != nil {
			logging.FrameworkTrace(ctx, "Attempting to populate MoveResourceStateRequest SourceState from provider defined SourceSchema")

			sourceSchemaType := resourceStateMover.SourceSchema.Type().TerraformType(ctx)
			unmarshalOpts := tfprotov6.UnmarshalOpts{
				ValueFromJSONOpts: tftypes.ValueFromJSONOpts{
					// IgnoreUndefinedAttributes will silently skip over fields
					// in the JSON that do not have a matching definition in the
					// given schema. The purpose of this is to allow for
					// additive changes to the source resource schema without
					// breaking target resource state moves. It also enables
					// simplified implementations, if certain source data is not
					// needed anyways.
					IgnoreUndefinedAttributes: true,
				},
			}

			rawStateValue, err := req.SourceRawState.UnmarshalWithOpts(sourceSchemaType, unmarshalOpts)

			// Resources may support multiple source resources, so returning the
			// error here as an error or warning diagnostic is not appropriate
			// since both the developer and calling practitioner cannot avoid
			// the situation. Instead, developers will still have a nil
			// SourceState and they can investigate any error as logged here.
			//
			// It is also important to note that the error generally only occurs
			// if the source schema declared incompatible types. The
			// IgnoreUndefinedAttributes option above can cause the error to be
			// nil and the SourceState to be populated with null values. It is
			// always recommended for StateMover implementations to check the
			// other request fields (SourceTypeName, SourceProviderAddress,
			// etc.) instead of relying on SourceState to be populated or not.
			if err != nil {
				logging.FrameworkDebug(
					ctx,
					"Error unmarshalling SourceRawState using the provided SourceSchema for source "+
						req.SourceProviderAddress+" resource type "+
						req.SourceTypeName+" with schema version "+
						strconv.FormatInt(req.SourceSchemaVersion, 10)+". "+
						"This is not a fatal error since resources can support multiple source resources which cause this type of error to be unavoidable, "+
						"but due to this error the SourceState will not be populated for the implementation.",
					map[string]any{
						logging.KeyError: err,
					},
				)
			} else {
				moveStateReq.SourceState = &tfsdk.State{
					Raw:    rawStateValue,
					Schema: *resourceStateMover.SourceSchema,
				}
			}
		}

		logging.FrameworkTrace(ctx, "Calling provider defined Resource StateMover")
		resourceStateMover.StateMover(ctx, moveStateReq, &moveStateResp)
		logging.FrameworkTrace(ctx, "Called provider defined Resource StateMover")

		resp.Diagnostics.Append(moveStateResp.Diagnostics...)

		// If the implementation has error diagnostics, return the diagnostics.
		if moveStateResp.Diagnostics.HasError() {
			resp.Diagnostics = moveStateResp.Diagnostics

			return
		}

		// If the implement has set the state in any way, return the response.
		if !moveStateResp.TargetState.Raw.Equal(tftypes.NewValue(req.TargetResourceSchema.Type().TerraformType(ctx), nil)) {
			resp.Diagnostics = moveStateResp.Diagnostics
			resp.TargetState = &moveStateResp.TargetState

			if moveStateResp.TargetPrivate != nil {
				resp.TargetPrivate.Provider = moveStateResp.TargetPrivate
			}

			return
		}
	}

	resp.Diagnostics.AddError(
		"Unable to Move Resource State",
		"The target resource implementation does not include support for the given source resource. "+
			"The resource implementation can be updated by the provider developers to include this support by returning the moved state when the request matches this source.\n\n"+
			"Source Provider Address: "+req.SourceProviderAddress+"\n"+
			"Source Resource Type: "+req.SourceTypeName+"\n"+
			"Source Resource Schema Version: "+strconv.FormatInt(req.SourceSchemaVersion, 10)+"\n"+
			"Target Resource Type: "+req.TargetTypeName,
	)
}
