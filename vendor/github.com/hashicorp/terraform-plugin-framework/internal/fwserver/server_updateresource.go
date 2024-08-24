// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// UpdateResourceRequest is the framework server request for an update request
// with the ApplyResourceChange RPC.
type UpdateResourceRequest struct {
	Config         *tfsdk.Config
	PlannedPrivate *privatestate.Data
	PlannedState   *tfsdk.Plan
	PriorState     *tfsdk.State
	ProviderMeta   *tfsdk.Config
	ResourceSchema fwschema.Schema
	Resource       resource.Resource
}

// UpdateResourceResponse is the framework server response for an update request
// with the ApplyResourceChange RPC.
type UpdateResourceResponse struct {
	Diagnostics diag.Diagnostics
	NewState    *tfsdk.State
	Private     *privatestate.Data
}

// UpdateResource implements the framework server update request logic for the
// ApplyResourceChange RPC.
func (s *Server) UpdateResource(ctx context.Context, req *UpdateResourceRequest, resp *UpdateResourceResponse) {
	if req == nil {
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

	nullSchemaData := tftypes.NewValue(req.ResourceSchema.Type().TerraformType(ctx), nil)

	updateReq := resource.UpdateRequest{
		Config: tfsdk.Config{
			Schema: req.ResourceSchema,
			Raw:    nullSchemaData,
		},
		Plan: tfsdk.Plan{
			Schema: req.ResourceSchema,
			Raw:    nullSchemaData,
		},
		State: tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    nullSchemaData,
		},
	}
	updateResp := resource.UpdateResponse{
		State: tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    nullSchemaData,
		},
	}

	if req.Config != nil {
		updateReq.Config = *req.Config
	}

	if req.PlannedState != nil {
		updateReq.Plan = *req.PlannedState
	}

	if req.PriorState != nil {
		updateReq.State = *req.PriorState
		// Require explicit provider updates for tracking successful updates.
		updateResp.State = *req.PriorState
	}

	if req.ProviderMeta != nil {
		updateReq.ProviderMeta = *req.ProviderMeta
	}

	privateProviderData := privatestate.EmptyProviderData(ctx)

	updateReq.Private = privateProviderData
	updateResp.Private = privateProviderData

	if req.PlannedPrivate != nil {
		if req.PlannedPrivate.Provider != nil {
			updateReq.Private = req.PlannedPrivate.Provider
			updateResp.Private = req.PlannedPrivate.Provider
		}

		resp.Private = req.PlannedPrivate
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Resource Update")
	req.Resource.Update(ctx, updateReq, &updateResp)
	logging.FrameworkTrace(ctx, "Called provider defined Resource Update")

	resp.Diagnostics = updateResp.Diagnostics
	resp.NewState = &updateResp.State

	if !resp.Diagnostics.HasError() && updateResp.State.Raw.Equal(nullSchemaData) {
		resp.Diagnostics.AddError(
			"Missing Resource State After Update",
			"The Terraform Provider unexpectedly returned no resource state after having no errors in the resource update. "+
				"This is always an issue in the Terraform Provider and should be reported to the provider developers.",
		)
	}

	if updateResp.Private != nil {
		if resp.Private == nil {
			resp.Private = &privatestate.Data{}
		}

		resp.Private.Provider = updateResp.Private
	}

	if resp.Diagnostics.HasError() {
		return
	}

	semanticEqualityReq := SchemaSemanticEqualityRequest{
		PriorData: fwschemadata.Data{
			Description:    fwschemadata.DataDescriptionPlan,
			Schema:         req.PlannedState.Schema,
			TerraformValue: req.PlannedState.Raw.Copy(),
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

	if semanticEqualityResp.NewData.TerraformValue.Equal(resp.NewState.Raw) {
		return
	}

	logging.FrameworkDebug(ctx, "State updated due to semantic equality")

	resp.NewState.Raw = semanticEqualityResp.NewData.TerraformValue
}
