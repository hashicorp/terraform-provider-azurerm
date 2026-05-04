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

// CreateResourceRequest is the framework server request for a create request
// with the ApplyResourceChange RPC.
type CreateResourceRequest struct {
	Config          *tfsdk.Config
	PlannedPrivate  *privatestate.Data
	PlannedState    *tfsdk.Plan
	PlannedIdentity *tfsdk.ResourceIdentity
	ProviderMeta    *tfsdk.Config
	ResourceSchema  fwschema.Schema
	IdentitySchema  fwschema.Schema
	Resource        resource.Resource
}

// CreateResourceResponse is the framework server response for a create request
// with the ApplyResourceChange RPC.
type CreateResourceResponse struct {
	Diagnostics diag.Diagnostics
	NewState    *tfsdk.State
	NewIdentity *tfsdk.ResourceIdentity
	Private     *privatestate.Data
}

// CreateResource implements the framework server create request logic for the
// ApplyResourceChange RPC.
func (s *Server) CreateResource(ctx context.Context, req *CreateResourceRequest, resp *CreateResourceResponse) {
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

	createReq := resource.CreateRequest{
		Config: tfsdk.Config{
			Schema: req.ResourceSchema,
			Raw:    nullSchemaData,
		},
		Plan: tfsdk.Plan{
			Schema: req.ResourceSchema,
			Raw:    nullSchemaData,
		},
	}

	privateProviderData := privatestate.EmptyProviderData(ctx)

	createResp := resource.CreateResponse{
		State: tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    nullSchemaData,
		},
		Private: privateProviderData,
	}

	if req.Config != nil {
		createReq.Config = *req.Config
	}

	if req.PlannedState != nil {
		createReq.Plan = *req.PlannedState
	}

	if req.ProviderMeta != nil {
		createReq.ProviderMeta = *req.ProviderMeta
	}

	// If the resource supports identity and there is no planned identity data, pre-populate with a null value.
	if req.PlannedIdentity == nil && req.IdentitySchema != nil {
		nullIdentityTfValue := tftypes.NewValue(req.IdentitySchema.Type().TerraformType(ctx), nil)

		req.PlannedIdentity = &tfsdk.ResourceIdentity{
			Schema: req.IdentitySchema,
			Raw:    nullIdentityTfValue.Copy(),
		}
	}

	// Pre-populate the new identity with the planned identity.
	if req.PlannedIdentity != nil {
		createReq.Identity = &tfsdk.ResourceIdentity{
			Schema: req.PlannedIdentity.Schema,
			Raw:    req.PlannedIdentity.Raw.Copy(),
		}

		createResp.Identity = &tfsdk.ResourceIdentity{
			Schema: req.PlannedIdentity.Schema,
			Raw:    req.PlannedIdentity.Raw.Copy(),
		}
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Resource Create")
	req.Resource.Create(ctx, createReq, &createResp)
	logging.FrameworkTrace(ctx, "Called provider defined Resource Create")

	resp.Diagnostics = createResp.Diagnostics
	resp.NewState = &createResp.State
	resp.NewIdentity = createResp.Identity

	if !resp.Diagnostics.HasError() && createResp.State.Raw.Equal(nullSchemaData) {
		detail := "The Terraform Provider unexpectedly returned no resource state after having no errors in the resource creation. " +
			"This is always an issue in the Terraform Provider and should be reported to the provider developers.\n\n" +
			"The resource may have been successfully created, but Terraform is not tracking it. " +
			"Applying the configuration again with no other action may result in duplicate resource errors."

		if _, ok := req.Resource.(resource.ResourceWithImportState); ok {
			detail += " Import the resource if the resource was actually created and Terraform should be tracking it."
		}

		resp.Diagnostics.AddError(
			"Missing Resource State After Create",
			detail,
		)
	}

	if createResp.Private != nil {
		if resp.Private == nil {
			resp.Private = &privatestate.Data{}
		}

		resp.Private.Provider = createResp.Private
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if resp.NewIdentity != nil && req.IdentitySchema == nil {
		resp.Diagnostics.AddError(
			"Unexpected Create Response",
			"An unexpected error was encountered when creating the apply response. New identity data was returned by the provider create operation, but the resource does not indicate identity support.\n\n"+
				"This is always a problem with the provider and should be reported to the provider developer.",
		)

		return
	}

	if req.IdentitySchema != nil {
		if resp.NewIdentity.Raw.IsFullyNull() {
			resp.Diagnostics.AddError(
				"Missing Resource Identity After Create",
				"The Terraform Provider unexpectedly returned no resource identity data after having no errors in the resource create. "+
					"This is always an issue in the Terraform Provider and should be reported to the provider developers.",
			)
			return
		}
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
