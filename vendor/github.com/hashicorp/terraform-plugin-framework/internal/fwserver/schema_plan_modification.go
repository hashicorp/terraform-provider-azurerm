// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ModifySchemaPlanRequest represents a request for a schema to run all
// attribute plan modification functions.
type ModifySchemaPlanRequest struct {
	// Config is the configuration the user supplied for the resource.
	Config tfsdk.Config

	// State is the current state of the resource.
	State tfsdk.State

	// Plan is the planned new state for the resource.
	Plan tfsdk.Plan

	// ProviderMeta is metadata from the provider_meta block of the module.
	ProviderMeta tfsdk.Config

	// Private is provider private state data.
	Private *privatestate.ProviderData
}

// ModifySchemaPlanResponse represents a response to a ModifySchemaPlanRequest.
type ModifySchemaPlanResponse struct {
	// Plan is the planned new state for the resource.
	Plan tfsdk.Plan

	// RequiresReplace is a list of attribute paths that require the
	// resource to be replaced. They should point to the specific field
	// that changed that requires the resource to be destroyed and
	// recreated.
	RequiresReplace path.Paths

	// Private is provider private state data following potential modifications.
	Private *privatestate.ProviderData

	// Diagnostics report errors or warnings related to running all attribute
	// plan modifiers. Returning an empty slice indicates a successful
	// plan modification with no warnings or errors generated.
	Diagnostics diag.Diagnostics
}

// SchemaModifyPlan runs all AttributePlanModifiers in all schema attributes
// and blocks.
//
// TODO: Clean up this abstraction back into an internal Schema type method.
// The extra Schema parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func SchemaModifyPlan(ctx context.Context, s fwschema.Schema, req ModifySchemaPlanRequest, resp *ModifySchemaPlanResponse) {
	var diags diag.Diagnostics

	configData := &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionConfiguration,
		Schema:         req.Config.Schema,
		TerraformValue: req.Config.Raw,
	}

	planData := &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionPlan,
		Schema:         req.Plan.Schema,
		TerraformValue: req.Plan.Raw,
	}

	stateData := &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionState,
		Schema:         req.State.Schema,
		TerraformValue: req.State.Raw,
	}

	for name, attribute := range s.GetAttributes() {
		attrReq := ModifyAttributePlanRequest{
			AttributePath: path.Root(name),
			Config:        req.Config,
			State:         req.State,
			Plan:          req.Plan,
			ProviderMeta:  req.ProviderMeta,
			Private:       req.Private,
		}

		attrReq.AttributeConfig, diags = configData.ValueAtPath(ctx, attrReq.AttributePath)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		attrReq.AttributePlan, diags = planData.ValueAtPath(ctx, attrReq.AttributePath)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		attrReq.AttributeState, diags = stateData.ValueAtPath(ctx, attrReq.AttributePath)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		attrResp := ModifyAttributePlanResponse{
			AttributePlan: attrReq.AttributePlan,
			Private:       attrReq.Private,
		}

		AttributeModifyPlan(ctx, attribute, attrReq, &attrResp)

		resp.Diagnostics.Append(attrResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, attrReq.AttributePath, attrResp.AttributePlan)...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.RequiresReplace = append(resp.RequiresReplace, attrResp.RequiresReplace...)
		resp.Private = attrResp.Private
	}

	for name, block := range s.GetBlocks() {
		blockReq := ModifyAttributePlanRequest{
			AttributePath: path.Root(name),
			Config:        req.Config,
			State:         req.State,
			Plan:          req.Plan,
			ProviderMeta:  req.ProviderMeta,
			Private:       req.Private,
		}

		blockReq.AttributeConfig, diags = configData.ValueAtPath(ctx, blockReq.AttributePath)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		blockReq.AttributePlan, diags = planData.ValueAtPath(ctx, blockReq.AttributePath)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		blockReq.AttributeState, diags = stateData.ValueAtPath(ctx, blockReq.AttributePath)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		blockResp := ModifyAttributePlanResponse{
			AttributePlan: blockReq.AttributePlan,
			Private:       blockReq.Private,
		}

		BlockModifyPlan(ctx, block, blockReq, &blockResp)

		resp.Diagnostics.Append(blockResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, blockReq.AttributePath, blockResp.AttributePlan)...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.RequiresReplace = append(resp.RequiresReplace, blockResp.RequiresReplace...)
		resp.Private = blockResp.Private
	}
}
