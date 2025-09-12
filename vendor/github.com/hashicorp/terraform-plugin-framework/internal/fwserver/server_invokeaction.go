// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// InvokeActionRequest is the framework server request for the InvokeAction RPC.
type InvokeActionRequest struct {
	Action       action.Action
	ActionSchema fwschema.Schema
	Config       *tfsdk.Config
}

// InvokeActionEventsStream is the framework server stream for the InvokeAction RPC.
type InvokeActionResponse struct {
	// ProgressEvents is a channel provided by the consuming proto{5/6}server implementation
	// that allows the provider developers to return progress events while the action is being invoked.
	ProgressEvents chan InvokeProgressEvent
	Diagnostics    diag.Diagnostics
}

type InvokeProgressEvent struct {
	Message string
}

// SendProgress is injected into the action.InvokeResponse for use by the provider developer
func (r *InvokeActionResponse) SendProgress(event action.InvokeProgressEvent) {
	r.ProgressEvents <- InvokeProgressEvent{
		Message: event.Message,
	}
}

// InvokeAction implements the framework server InvokeAction RPC.
func (s *Server) InvokeAction(ctx context.Context, req *InvokeActionRequest, resp *InvokeActionResponse) {
	if req == nil {
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

	invokeReq := action.InvokeRequest{
		Config: *req.Config,
	}
	invokeResp := action.InvokeResponse{
		SendProgress: resp.SendProgress,
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Action Invoke")
	req.Action.Invoke(ctx, invokeReq, &invokeResp)
	logging.FrameworkTrace(ctx, "Called provider defined Action Invoke")

	resp.Diagnostics = invokeResp.Diagnostics
}
