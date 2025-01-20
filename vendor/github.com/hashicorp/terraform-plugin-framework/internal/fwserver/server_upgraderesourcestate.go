// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// UpgradeResourceStateRequest is the framework server request for the
// UpgradeResourceState RPC.
type UpgradeResourceStateRequest struct {
	// Using the tfprotov6 type here was a pragmatic effort decision around when
	// the framework introduced compatibility promises. This type was chosen as
	// it was readily available and trivial to convert between tfprotov5.
	//
	// Using a terraform-plugin-go type is not ideal for the framework as almost
	// all terraform-plugin-go types have framework abstractions, but if there
	// is ever a time where it makes sense to re-evaluate this decision, such as
	// a major version bump, it could be changed then.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/340
	RawState *tfprotov6.RawState

	ResourceSchema fwschema.Schema
	Resource       resource.Resource
	Version        int64
}

// UpgradeResourceStateResponse is the framework server response for the
// UpgradeResourceState RPC.
type UpgradeResourceStateResponse struct {
	Diagnostics   diag.Diagnostics
	UpgradedState *tfsdk.State
}

// UpgradeResourceState implements the framework server UpgradeResourceState RPC.
func (s *Server) UpgradeResourceState(ctx context.Context, req *UpgradeResourceStateRequest, resp *UpgradeResourceStateResponse) {
	if req == nil {
		return
	}

	// No UpgradedState to return. This could return an error diagnostic about
	// the odd scenario, but seems best to allow Terraform CLI to handle the
	// situation itself in case it might be expected behavior.
	if req.RawState == nil {
		return
	}

	// Define options to be used when unmarshalling raw state.
	// IgnoreUndefinedAttributes will silently skip over fields in the JSON
	// that do not have a matching entry in the schema.
	unmarshalOpts := tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{
			IgnoreUndefinedAttributes: true,
		},
	}

	// Terraform CLI can call UpgradeResourceState even if the stored state
	// version matches the current schema. Presumably this is to account for
	// the previous terraform-plugin-sdk implementation, which handled some
	// state fixups on behalf of Terraform CLI. When this happens, we do not
	// want to return errors for a missing ResourceWithUpgradeState
	// implementation or an undefined version within an existing
	// ResourceWithUpgradeState implementation as that would be confusing
	// detail for provider developers. Instead, the framework will attempt to
	// roundtrip the prior RawState to a State matching the current Schema.
	//
	// TODO: To prevent provider developers from accidentally implementing
	// ResourceWithUpgradeState with a version matching the current schema
	// version which would never get called, the framework can introduce a
	// unit test helper.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/113
	//
	// UnmarshalWithOpts allows optionally ignoring instances in which elements being
	// do not have a corresponding attribute within the schema.
	if req.Version == req.ResourceSchema.GetVersion() {
		logging.FrameworkTrace(ctx, "UpgradeResourceState request version matches current Schema version, using framework defined passthrough implementation")

		resourceSchemaType := req.ResourceSchema.Type().TerraformType(ctx)

		rawStateValue, err := req.RawState.UnmarshalWithOpts(resourceSchemaType, unmarshalOpts)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Previously Saved State for UpgradeResourceState",
				"There was an error reading the saved resource state using the current resource schema.\n\n"+
					"If this resource state was last refreshed with Terraform CLI 0.11 and earlier, it must be refreshed or applied with an older provider version first. "+
					"If you manually modified the resource state, you will need to manually modify it to match the current resource schema. "+
					"Otherwise, please report this to the provider developer:\n\n"+err.Error(),
			)
			return
		}

		resp.UpgradedState = &tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    rawStateValue,
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

	resourceWithUpgradeState, ok := req.Resource.(resource.ResourceWithUpgradeState)

	if !ok {
		resp.Diagnostics.AddError(
			"Unable to Upgrade Resource State",
			"This resource was implemented without an UpgradeState() method, "+
				fmt.Sprintf("however Terraform was expecting an implementation for version %d upgrade.\n\n", req.Version)+
				"This is always an issue with the Terraform Provider and should be reported to the provider developer.",
		)
		return
	}

	logging.FrameworkTrace(ctx, "Resource implements ResourceWithUpgradeState")

	logging.FrameworkTrace(ctx, "Calling provider defined Resource UpgradeState")
	resourceStateUpgraders := resourceWithUpgradeState.UpgradeState(ctx)
	logging.FrameworkTrace(ctx, "Called provider defined Resource UpgradeState")

	// Panic prevention
	if resourceStateUpgraders == nil {
		resourceStateUpgraders = make(map[int64]resource.StateUpgrader, 0)
	}

	resourceStateUpgrader, ok := resourceStateUpgraders[req.Version]

	if !ok {
		resp.Diagnostics.AddError(
			"Unable to Upgrade Resource State",
			"This resource was implemented with an UpgradeState() method, "+
				fmt.Sprintf("however Terraform was expecting an implementation for version %d upgrade.\n\n", req.Version)+
				"This is always an issue with the Terraform Provider and should be reported to the provider developer.",
		)
		return
	}

	upgradeResourceStateRequest := resource.UpgradeStateRequest{
		RawState: req.RawState,
	}

	if resourceStateUpgrader.PriorSchema != nil {
		logging.FrameworkTrace(ctx, "Initializing populated UpgradeResourceStateRequest state from provider defined prior schema and request RawState")

		priorSchemaType := resourceStateUpgrader.PriorSchema.Type().TerraformType(ctx)

		rawStateValue, err := req.RawState.UnmarshalWithOpts(priorSchemaType, unmarshalOpts)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Previously Saved State for UpgradeResourceState",
				fmt.Sprintf("There was an error reading the saved resource state using the prior resource schema defined for version %d upgrade.\n\n", req.Version)+
					"Please report this to the provider developer:\n\n"+err.Error(),
			)
			return
		}

		upgradeResourceStateRequest.State = &tfsdk.State{
			Raw:    rawStateValue,
			Schema: *resourceStateUpgrader.PriorSchema,
		}
	}

	upgradeResourceStateResponse := resource.UpgradeStateResponse{
		State: tfsdk.State{
			Schema: req.ResourceSchema,
			// Raw is intentionally not set.
		},
	}

	// To simplify provider logic, this could perform a best effort attempt
	// to populate the response State by looping through all Attribute/Block
	// by calling the equivalent of SetAttribute(GetAttribute()) and skipping
	// any errors.

	logging.FrameworkTrace(ctx, "Calling provider defined StateUpgrader")
	resourceStateUpgrader.StateUpgrader(ctx, upgradeResourceStateRequest, &upgradeResourceStateResponse)
	logging.FrameworkTrace(ctx, "Called provider defined StateUpgrader")

	resp.Diagnostics.Append(upgradeResourceStateResponse.Diagnostics...)

	if resp.Diagnostics.HasError() {
		return
	}

	if upgradeResourceStateResponse.DynamicValue != nil {
		logging.FrameworkTrace(ctx, "UpgradeResourceStateResponse DynamicValue set, overriding State")

		upgradedStateValue, err := upgradeResourceStateResponse.DynamicValue.Unmarshal(req.ResourceSchema.Type().TerraformType(ctx))

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Upgrade Resource State",
				fmt.Sprintf("After attempting a resource state upgrade to version %d, the provider returned state data that was not compatible with the current schema.\n\n", req.Version)+
					"This is always an issue with the Terraform Provider and should be reported to the provider developer:\n\n"+err.Error(),
			)
			return
		}

		resp.UpgradedState = &tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    upgradedStateValue,
		}

		return
	}

	if upgradeResourceStateResponse.State.Raw.Type() == nil || upgradeResourceStateResponse.State.Raw.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Upgraded Resource State",
			fmt.Sprintf("After attempting a resource state upgrade to version %d, the provider did not return any state data. ", req.Version)+
				"Preventing the unexpected loss of resource state data. "+
				"This is always an issue with the Terraform Provider and should be reported to the provider developer.",
		)
		return
	}

	resp.UpgradedState = &upgradeResourceStateResponse.State
}
