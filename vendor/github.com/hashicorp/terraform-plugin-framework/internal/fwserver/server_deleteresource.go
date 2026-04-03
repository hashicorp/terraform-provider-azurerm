// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// DeleteResourceRequest is the framework server request for a delete request
// with the ApplyResourceChange RPC.
type DeleteResourceRequest struct {
	PlannedPrivate *privatestate.Data
	PriorState     *tfsdk.State
	PriorIdentity  *tfsdk.ResourceIdentity
	ProviderMeta   *tfsdk.Config
	ResourceSchema fwschema.Schema
	IdentitySchema fwschema.Schema
	Resource       resource.Resource
}

// DeleteResourceResponse is the framework server response for a delete request
// with the ApplyResourceChange RPC.
type DeleteResourceResponse struct {
	Diagnostics diag.Diagnostics
	NewState    *tfsdk.State
	NewIdentity *tfsdk.ResourceIdentity
	Private     *privatestate.Data
}

// DeleteResource implements the framework server delete request logic for the
// ApplyResourceChange RPC.
func (s *Server) DeleteResource(ctx context.Context, req *DeleteResourceRequest, resp *DeleteResourceResponse) {
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

	deleteReq := resource.DeleteRequest{
		State: tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    tftypes.NewValue(req.ResourceSchema.Type().TerraformType(ctx), nil),
		},
	}
	deleteResp := resource.DeleteResponse{
		State: tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    tftypes.NewValue(req.ResourceSchema.Type().TerraformType(ctx), nil),
		},
	}

	if req.PriorState != nil {
		deleteReq.State = *req.PriorState
		deleteResp.State = *req.PriorState
	}

	if req.ProviderMeta != nil {
		deleteReq.ProviderMeta = *req.ProviderMeta
	}

	privateProviderData := privatestate.EmptyProviderData(ctx)

	deleteReq.Private = privateProviderData
	deleteResp.Private = privateProviderData

	if req.PlannedPrivate != nil {
		if req.PlannedPrivate.Provider != nil {
			deleteReq.Private = req.PlannedPrivate.Provider
			deleteResp.Private = req.PlannedPrivate.Provider
		}

		resp.Private = req.PlannedPrivate
	}

	if req.PriorIdentity == nil && req.IdentitySchema != nil {
		nullIdentityTfValue := tftypes.NewValue(req.IdentitySchema.Type().TerraformType(ctx), nil)

		req.PriorIdentity = &tfsdk.ResourceIdentity{
			Schema: req.IdentitySchema,
			Raw:    nullIdentityTfValue.Copy(),
		}
	}

	if req.PriorIdentity != nil {
		deleteReq.Identity = &tfsdk.ResourceIdentity{
			Schema: req.PriorIdentity.Schema,
			Raw:    req.PriorIdentity.Raw.Copy(),
		}

		deleteResp.Identity = &tfsdk.ResourceIdentity{
			Schema: req.PriorIdentity.Schema,
			Raw:    req.PriorIdentity.Raw.Copy(),
		}
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Resource Delete")
	req.Resource.Delete(ctx, deleteReq, &deleteResp)
	logging.FrameworkTrace(ctx, "Called provider defined Resource Delete")

	if !deleteResp.Diagnostics.HasError() {
		logging.FrameworkTrace(ctx, "No provider defined Delete errors detected, ensuring State and Private are cleared")
		deleteResp.State.RemoveResource(ctx)

		// Preserve prior behavior of always returning nil.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/863
		deleteResp.Private = nil
		resp.Private = nil

		// If the resource supports identity send a null value.
		if req.IdentitySchema != nil {
			nullIdentityTfValue := tftypes.NewValue(req.IdentitySchema.Type().TerraformType(ctx), nil)

			deleteResp.Identity = &tfsdk.ResourceIdentity{
				Schema: req.IdentitySchema,
				Raw:    nullIdentityTfValue.Copy(),
			}
		}
	}

	resp.Diagnostics = deleteResp.Diagnostics
	resp.NewState = &deleteResp.State
	resp.NewIdentity = deleteResp.Identity

	if deleteResp.Private != nil {
		if resp.Private == nil {
			resp.Private = &privatestate.Data{}
		}

		resp.Private.Provider = deleteResp.Private
	}

	if resp.NewIdentity != nil && req.IdentitySchema == nil {
		resp.Diagnostics.AddError(
			"Unexpected Delete Response",
			"An unexpected error was encountered when creating the apply response. New identity data was returned by the provider delete operation, but the resource does not indicate identity support.\n\n"+
				"This is always a problem with the provider and should be reported to the provider developer.",
		)

		return
	}
}
