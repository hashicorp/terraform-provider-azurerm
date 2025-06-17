// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ReadResourceRequest is the framework server request for the
// ReadResource RPC.
type ReadResourceRequest struct {
	ClientCapabilities resource.ReadClientCapabilities
	CurrentState       *tfsdk.State
	Resource           resource.Resource
	Private            *privatestate.Data
	ProviderMeta       *tfsdk.Config
}

// ReadResourceResponse is the framework server response for the
// ReadResource RPC.
type ReadResourceResponse struct {
	Deferred    *resource.Deferred
	Diagnostics diag.Diagnostics
	NewState    *tfsdk.State
	Private     *privatestate.Data
}

// ReadResource implements the framework server ReadResource RPC.
func (s *Server) ReadResource(ctx context.Context, req *ReadResourceRequest, resp *ReadResourceResponse) {
	if req == nil {
		return
	}

	if req.CurrentState == nil {
		resp.Diagnostics.AddError(
			"Unexpected Read Request",
			"An unexpected error was encountered when reading the resource. The current state was missing.\n\n"+
				"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.",
		)

		return
	}

	if s.deferred != nil {
		logging.FrameworkDebug(ctx, "Provider has deferred response configured, automatically returning deferred response.",
			map[string]interface{}{
				logging.KeyDeferredReason: s.deferred.Reason.String(),
			},
		)
		resp.NewState = req.CurrentState
		resp.Deferred = &resource.Deferred{
			Reason: resource.DeferredReason(s.deferred.Reason),
		}
		return
	}

	if resourceWithConfigure, ok := req.Resource.(resource.ResourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "Resource implements ResourceWithConfigure")

		configureReq := resource.ConfigureRequest{
			ProviderData: s.ResourceConfigureData,
		}
		configureResp := resource.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Resource Configure")
		resourceWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined Resource Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	readReq := resource.ReadRequest{
		ClientCapabilities: req.ClientCapabilities,
		State: tfsdk.State{
			Schema: req.CurrentState.Schema,
			Raw:    req.CurrentState.Raw.Copy(),
		},
	}
	readResp := resource.ReadResponse{
		State: tfsdk.State{
			Schema: req.CurrentState.Schema,
			Raw:    req.CurrentState.Raw.Copy(),
		},
	}

	if req.ProviderMeta != nil {
		readReq.ProviderMeta = *req.ProviderMeta
	}

	privateProviderData := privatestate.EmptyProviderData(ctx)

	readReq.Private = privateProviderData
	readResp.Private = privateProviderData

	if req.Private != nil {
		if req.Private.Provider != nil {
			readReq.Private = req.Private.Provider
			readResp.Private = req.Private.Provider
		}

		resp.Private = req.Private
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Resource Read")
	req.Resource.Read(ctx, readReq, &readResp)
	logging.FrameworkTrace(ctx, "Called provider defined Resource Read")

	resp.Diagnostics = readResp.Diagnostics
	resp.NewState = &readResp.State
	resp.Deferred = readResp.Deferred

	if readResp.Private != nil {
		if resp.Private == nil {
			resp.Private = &privatestate.Data{}
		}

		resp.Private.Provider = readResp.Private
	}

	if resp.Diagnostics.HasError() {
		return
	}

	semanticEqualityReq := SchemaSemanticEqualityRequest{
		PriorData: fwschemadata.Data{
			Description:    fwschemadata.DataDescriptionState,
			Schema:         req.CurrentState.Schema,
			TerraformValue: req.CurrentState.Raw.Copy(),
		},
		ProposedNewData: fwschemadata.Data{
			Description:    fwschemadata.DataDescriptionState,
			Schema:         resp.NewState.Schema,
			TerraformValue: resp.NewState.Raw.Copy(),
		},
	}
	semanticEqualityResp := &SchemaSemanticEqualityResponse{
		NewData: semanticEqualityReq.ProposedNewData,
	}

	SchemaSemanticEquality(ctx, semanticEqualityReq, semanticEqualityResp)

	resp.Diagnostics.Append(semanticEqualityResp.Diagnostics...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !semanticEqualityResp.NewData.TerraformValue.Equal(resp.NewState.Raw) {
		logging.FrameworkDebug(ctx, "State updated due to semantic equality")

		resp.NewState.Raw = semanticEqualityResp.NewData.TerraformValue
	}

	// Set any write-only attributes in the state to null
	modifiedState, err := tftypes.Transform(resp.NewState.Raw, NullifyWriteOnlyAttributes(ctx, resp.NewState.Schema))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Modifying State",
			"There was an unexpected error modifying the NewState. This is always a problem with the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return
	}

	resp.NewState.Raw = modifiedState
}
