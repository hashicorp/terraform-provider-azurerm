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

// ImportedResource represents a resource that was imported.
type ImportedResource struct {
	Private  *privatestate.Data
	Identity *tfsdk.ResourceIdentity
	State    tfsdk.State
	TypeName string
}

// ImportResourceStateRequest is the framework server request for the
// ImportResourceState RPC.
//
// Either ID or Identity will be supplied depending on how the resource is being imported.
type ImportResourceStateRequest struct {
	// ID will come from the import CLI command or an import config block with the "id" attribute assigned.
	//
	// This ID field is a special string identifier that can be parsed however the provider deems fit.
	ID string

	// Identity will come from an import config block with the "identity" attribute assigned and will conform
	// to the identity schema defined by the resource. (Terraform v1.12+)
	//
	// All attributes marked as RequiredForImport will be populated (enforced by Terraform core) and OptionalForImport
	// attributes may be null, but could have a config value.
	Identity       *tfsdk.ResourceIdentity
	IdentitySchema fwschema.Schema

	Resource resource.Resource

	// EmptyState is an empty State for the resource schema. This is used to
	// initialize the ImportedResource State of the ImportResourceStateResponse
	// and allow the framework server to verify that the provider updated the
	// state after the provider defined logic.
	EmptyState tfsdk.State

	// TypeName is the resource type name, which is necessary for populating
	// the ImportedResource TypeName of the ImportResourceStateResponse.
	TypeName string

	ClientCapabilities resource.ImportStateClientCapabilities
}

// ImportResourceStateResponse is the framework server response for the
// ImportResourceState RPC.
type ImportResourceStateResponse struct {
	Diagnostics       diag.Diagnostics
	ImportedResources []ImportedResource
	Deferred          *resource.Deferred
}

// ImportResourceState implements the framework server ImportResourceState RPC.
func (s *Server) ImportResourceState(ctx context.Context, req *ImportResourceStateRequest, resp *ImportResourceStateResponse) {
	if req == nil {
		return
	}

	if s.deferred != nil {
		logging.FrameworkDebug(ctx, "Provider has deferred response configured, automatically returning deferred response.",
			map[string]interface{}{
				logging.KeyDeferredReason: s.deferred.Reason.String(),
			},
		)
		// Send an unknown value for the imported object
		resp.ImportedResources = []ImportedResource{
			{
				State: tfsdk.State{
					Raw:    tftypes.NewValue(req.EmptyState.Schema.Type().TerraformType(ctx), tftypes.UnknownValue),
					Schema: req.EmptyState.Schema,
				},
				TypeName: req.TypeName,
				Private:  &privatestate.Data{},
			},
		}
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
		ID:                 req.ID,
		ClientCapabilities: req.ClientCapabilities,
	}

	privateProviderData := privatestate.EmptyProviderData(ctx)

	importResp := resource.ImportStateResponse{
		State: tfsdk.State{
			Raw:    req.EmptyState.Raw.Copy(),
			Schema: req.EmptyState.Schema,
		},
		Private: privateProviderData,
	}

	// If the resource supports identity and we are not importing by identity, pre-populate with a null value.
	if req.Identity == nil && req.IdentitySchema != nil {
		nullTfValue := tftypes.NewValue(req.IdentitySchema.Type().TerraformType(ctx), nil)

		req.Identity = &tfsdk.ResourceIdentity{
			Schema: req.IdentitySchema,
			Raw:    nullTfValue.Copy(),
		}
	}

	if req.Identity != nil {
		importReq.Identity = &tfsdk.ResourceIdentity{
			Schema: req.Identity.Schema,
			Raw:    req.Identity.Raw.Copy(),
		}

		importResp.Identity = &tfsdk.ResourceIdentity{
			Schema: req.Identity.Schema,
			Raw:    req.Identity.Raw.Copy(),
		}
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Resource ImportState")
	resourceWithImportState.ImportState(ctx, importReq, &importResp)
	logging.FrameworkTrace(ctx, "Called provider defined Resource ImportState")

	resp.Diagnostics.Append(importResp.Diagnostics...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set any write-only attributes in the import state to null
	modifiedState, err := tftypes.Transform(importResp.State.Raw, NullifyWriteOnlyAttributes(ctx, importResp.State.Schema))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Modifying Import State",
			"There was an unexpected error modifying the Import State. This is always a problem with the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return
	}

	importResp.State.Raw = modifiedState

	// If we are importing by ID, we should ensure that something in the import stub state has been populated,
	// otherwise the resource doesn't actually support import, which is a provider issue.
	if req.ID != "" && importResp.State.Raw.Equal(req.EmptyState.Raw) {
		resp.Diagnostics.AddError(
			"Missing Resource Import State",
			"An unexpected error was encountered when importing the resource. This is always a problem with the provider. Please give the following information to the provider developer:\n\n"+
				"Resource ImportState method returned no State in response. If import is intentionally not supported, remove the Resource type ImportState method or return an error.",
		)
		return
	}

	private := &privatestate.Data{}

	// Set an internal private field that will get sent alongside the imported resource. This will be cleared by
	// the following ReadResource RPC and is primarily used to control validation of resource identities during refresh.
	private.Framework = map[string][]byte{
		privatestate.ImportBeforeReadKey: []byte(`true`), // The actual data isn't important, we just use the map key to detect it.
	}

	if importResp.Private != nil {
		private.Provider = importResp.Private
	}

	if importResp.Identity != nil && req.IdentitySchema == nil {
		resp.Diagnostics.AddError(
			"Unexpected ImportState Response",
			"An unexpected error was encountered when creating the import response. New identity data was returned by the provider import operation, but the resource does not indicate identity support.\n\n"+
				"This is always a problem with the provider and should be reported to the provider developer.",
		)

		return
	}

	resp.Deferred = importResp.Deferred
	resp.ImportedResources = []ImportedResource{
		{
			State:    importResp.State,
			Identity: importResp.Identity,
			TypeName: req.TypeName,
			Private:  private,
		},
	}
}
