// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// OpenEphemeralResourceRequest is the framework server request for the
// OpenEphemeralResource RPC.
type OpenEphemeralResourceRequest struct {
	ClientCapabilities      ephemeral.OpenClientCapabilities
	Config                  *tfsdk.Config
	EphemeralResourceSchema fwschema.Schema
	EphemeralResource       ephemeral.EphemeralResource
}

// OpenEphemeralResourceResponse is the framework server response for the
// OpenEphemeralResource RPC.
type OpenEphemeralResourceResponse struct {
	Result      *tfsdk.EphemeralResultData
	Private     *privatestate.Data
	Diagnostics diag.Diagnostics
	RenewAt     time.Time
	Deferred    *ephemeral.Deferred
}

// OpenEphemeralResource implements the framework server OpenEphemeralResource RPC.
func (s *Server) OpenEphemeralResource(ctx context.Context, req *OpenEphemeralResourceRequest, resp *OpenEphemeralResourceResponse) {
	if req == nil {
		return
	}

	if s.deferred != nil {
		logging.FrameworkDebug(ctx, "Provider has deferred response configured, automatically returning deferred response.",
			map[string]interface{}{
				logging.KeyDeferredReason: s.deferred.Reason.String(),
			},
		)
		// Send an unknown value for the ephemeral resource. This will replace any configured values
		// for ease of implementation as Terraform Core currently does not use these values for
		// deferred actions, but this design could change in the future.
		resp.Result = &tfsdk.EphemeralResultData{
			Raw:    tftypes.NewValue(req.EphemeralResourceSchema.Type().TerraformType(ctx), tftypes.UnknownValue),
			Schema: req.EphemeralResourceSchema,
		}
		resp.Deferred = &ephemeral.Deferred{
			Reason: ephemeral.DeferredReason(s.deferred.Reason),
		}
		return
	}

	if ephemeralResourceWithConfigure, ok := req.EphemeralResource.(ephemeral.EphemeralResourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "EphemeralResource implements EphemeralResourceWithConfigure")

		configureReq := ephemeral.ConfigureRequest{
			ProviderData: s.EphemeralResourceConfigureData,
		}
		configureResp := ephemeral.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined EphemeralResource Configure")
		ephemeralResourceWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined EphemeralResource Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	openReq := ephemeral.OpenRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config: tfsdk.Config{
			Schema: req.EphemeralResourceSchema,
		},
	}
	openResp := ephemeral.OpenResponse{
		Result: tfsdk.EphemeralResultData{
			Schema: req.EphemeralResourceSchema,
		},
		Private: privatestate.EmptyProviderData(ctx),
	}

	if req.Config != nil {
		openReq.Config = *req.Config
		openResp.Result.Raw = req.Config.Raw.Copy()
	}

	logging.FrameworkTrace(ctx, "Calling provider defined EphemeralResource Open")
	req.EphemeralResource.Open(ctx, openReq, &openResp)
	logging.FrameworkTrace(ctx, "Called provider defined EphemeralResource Open")

	resp.Diagnostics = openResp.Diagnostics
	resp.Result = &openResp.Result
	resp.RenewAt = openResp.RenewAt
	resp.Deferred = openResp.Deferred

	resp.Private = privatestate.EmptyData(ctx)
	if openResp.Private != nil {
		resp.Private.Provider = openResp.Private
	}
}
