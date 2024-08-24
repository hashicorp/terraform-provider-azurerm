// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ImportedResource represents a resource that was imported.
type ImportedResource struct {
	Private  *privatestate.Data
	State    tfsdk.State
	TypeName string
}

// ImportResourceStateRequest is the framework server request for the
// ImportResourceState RPC.
type ImportResourceStateRequest struct {
	ID       string
	Resource resource.Resource

	// EmptyState is an empty State for the resource schema. This is used to
	// initialize the ImportedResource State of the ImportResourceStateResponse
	// and allow the framework server to verify that the provider updated the
	// state after the provider defined logic.
	EmptyState tfsdk.State

	// TypeName is the resource type name, which is necessary for populating
	// the ImportedResource TypeName of the ImportResourceStateResponse.
	TypeName string
}

// ImportResourceStateResponse is the framework server response for the
// ImportResourceState RPC.
type ImportResourceStateResponse struct {
	Diagnostics       diag.Diagnostics
	ImportedResources []ImportedResource
}

// ImportResourceState implements the framework server ImportResourceState RPC.
func (s *Server) ImportResourceState(ctx context.Context, req *ImportResourceStateRequest, resp *ImportResourceStateResponse) {
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

	resourceWithImportState, ok := req.Resource.(resource.ResourceWithImportState)

	if !ok {
		// If there is a feature request for customizing this messaging,
		// provider developers can implement a ImportState method that
		// immediately returns a custom error diagnostic.
		//
		// However, implementing the ImportState method could cause issues
		// with automated documentation generation, which likely would check
		// if the resource implements the ResourceWithImportState interface.
		// Instead, a separate "ResourceWithoutImportState" interface could be
		// created with a method such as:
		//    ImportNotImplementedMessage(context.Context) string.
		resp.Diagnostics.AddError(
			"Resource Import Not Implemented",
			"This resource does not support import. Please contact the provider developer for additional information.",
		)
		return
	}

	importReq := resource.ImportStateRequest{
		ID: req.ID,
	}

	privateProviderData := privatestate.EmptyProviderData(ctx)

	importResp := resource.ImportStateResponse{
		State: tfsdk.State{
			Raw:    req.EmptyState.Raw.Copy(),
			Schema: req.EmptyState.Schema,
		},
		Private: privateProviderData,
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Resource ImportState")
	resourceWithImportState.ImportState(ctx, importReq, &importResp)
	logging.FrameworkTrace(ctx, "Called provider defined Resource ImportState")

	resp.Diagnostics.Append(importResp.Diagnostics...)

	if resp.Diagnostics.HasError() {
		return
	}

	if importResp.State.Raw.Equal(req.EmptyState.Raw) {
		resp.Diagnostics.AddError(
			"Missing Resource Import State",
			"An unexpected error was encountered when importing the resource. This is always a problem with the provider. Please give the following information to the provider developer:\n\n"+
				"Resource ImportState method returned no State in response. If import is intentionally not supported, remove the Resource type ImportState method or return an error.",
		)
		return
	}

	private := &privatestate.Data{}

	if importResp.Private != nil {
		private.Provider = importResp.Private
	}

	resp.ImportedResources = []ImportedResource{
		{
			State:    importResp.State,
			TypeName: req.TypeName,
			Private:  private,
		},
	}
}
