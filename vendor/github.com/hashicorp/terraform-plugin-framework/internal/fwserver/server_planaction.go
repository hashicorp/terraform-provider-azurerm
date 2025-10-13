// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// PlanActionRequest is the framework server request for the PlanAction RPC.
type PlanActionRequest struct {
	ClientCapabilities action.ModifyPlanClientCapabilities
	ActionSchema       fwschema.Schema
	Action             action.Action
	Config             *tfsdk.Config
}

// PlanActionResponse is the framework server response for the PlanAction RPC.
type PlanActionResponse struct {
	Deferred    *action.Deferred
	Diagnostics diag.Diagnostics
}

// PlanAction implements the framework server PlanAction RPC.
func (s *Server) PlanAction(ctx context.Context, req *PlanActionRequest, resp *PlanActionResponse) {
	if req == nil {
		return
	}

	if s.deferred != nil {
		logging.FrameworkDebug(ctx, "Provider has deferred response configured, automatically returning deferred response.",
			map[string]interface{}{
				logging.KeyDeferredReason: s.deferred.Reason.String(),
			},
		)

		resp.Deferred = &action.Deferred{
			Reason: action.DeferredReason(s.deferred.Reason),
		}
		return
	}

	if actionWithConfigure, ok := req.Action.(action.ActionWithConfigure); ok {
		logging.FrameworkTrace(ctx, "Action implements ActionWithConfigure")

		configureReq := action.ConfigureRequest{
			ProviderData: s.ActionConfigureData,
		}
		configureResp := action.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Action Configure")
		actionWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined Action Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	if req.Config == nil {
		req.Config = &tfsdk.Config{
			Raw:    tftypes.NewValue(req.ActionSchema.Type().TerraformType(ctx), nil),
			Schema: req.ActionSchema,
		}
	}

	if actionWithModifyPlan, ok := req.Action.(action.ActionWithModifyPlan); ok {
		logging.FrameworkTrace(ctx, "Action implements ActionWithModifyPlan")

		modifyPlanReq := action.ModifyPlanRequest{
			ClientCapabilities: req.ClientCapabilities,
			Config:             *req.Config,
		}

		modifyPlanResp := action.ModifyPlanResponse{
			Diagnostics: resp.Diagnostics,
		}

		logging.FrameworkTrace(ctx, "Calling provider defined Action ModifyPlan")
		actionWithModifyPlan.ModifyPlan(ctx, modifyPlanReq, &modifyPlanResp)
		logging.FrameworkTrace(ctx, "Called provider defined Action ModifyPlan")

		resp.Diagnostics = modifyPlanResp.Diagnostics
		resp.Deferred = modifyPlanResp.Deferred
	}
}
