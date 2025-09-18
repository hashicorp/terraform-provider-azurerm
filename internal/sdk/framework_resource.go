// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

var IDPath = path.Root("id")

func (r *ResourceMetadata) MarkAsGone(idFormatter resourceids.Id, state *tfsdk.State, diags *diag.Diagnostics) {
	diags.Append(diag.NewWarningDiagnostic(fmt.Sprintf("[DEBUG] %s was not found - removing from state", idFormatter), ""))
	state.SetAttribute(context.Background(), IDPath, nil)
}

func (r *ResourceMetadata) ResourceRequiresImport(resourceName string, idFormatter resourceids.Id, resp *resource.CreateResponse) {
	msg := "A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	resp.Diagnostics.AddError("Existing Resource Error", fmt.Sprintf(msg, idFormatter.ID(), resourceName))
}

type ResourceMetadata struct {
	Client *clients.Client

	SubscriptionId string

	TimeoutCreate time.Duration
	TimeoutRead   time.Duration
	TimeoutDelete time.Duration
	TimeoutUpdate *time.Duration

	Features features.UserFeatures
}

// Defaults configures the Resource Metadata for client access, Provider Features, default timeouts, and subscriptionId.
func (r *ResourceMetadata) Defaults(req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError("Client Provider Data Error", fmt.Sprintf("invalid provider data supplied, got %+v", req.ProviderData))
		return
	}

	r.Client = c
	r.SubscriptionId = c.Account.SubscriptionId
	r.Features = c.Features

	r.TimeoutCreate = 30 * time.Minute
	r.TimeoutUpdate = pointer.To(30 * time.Minute)
	r.TimeoutRead = 5 * time.Minute
	r.TimeoutDelete = 30 * time.Minute
}

// DefaultsDataSource configures the Resource Metadata for client access, Provider Features, and subscriptionId.
func (r *ResourceMetadata) DefaultsDataSource(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError("Client Provider Data Error", fmt.Sprintf("invalid provider data supplied, got %+v", req.ProviderData))
		return
	}

	r.Client = c
	r.SubscriptionId = c.Account.SubscriptionId
	r.Features = c.Features

	r.TimeoutRead = 5 * time.Minute
}

// DecodeCreate reads a plan from a resource.CreateRequest into a pointer to a target model and sets
// resource.CreateResponse diags on error.
func (r *ResourceMetadata) DecodeCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, planModel any) {
	resp.Diagnostics.Append(req.Plan.Get(ctx, planModel)...)
}

// DecodeCreateWithConfig reads both the config and plan from the CreateRequest for cases where the raw config is
// required before plan logic has been applied
// This should be used ONLY when absolutely required. Plan Modifiers should be preferred.
func (r *ResourceMetadata) DecodeCreateWithConfig(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, planModel, configModel any) {
	resp.Diagnostics.Append(req.Plan.Get(ctx, planModel)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, configModel)...)
}

// EncodeCreate writes the model populated in the Create method to state.
func (r *ResourceMetadata) EncodeCreate(ctx context.Context, resp *resource.CreateResponse, model any) {
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// DecodeRead reads a resource's Previous State from a resource.ReadRequest into a pointer to a target model and sets
// resource.ReadResponse diags on error.
func (r *ResourceMetadata) DecodeRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, model any) {
	if !req.State.Raw.IsNull() {
		resp.Diagnostics.Append(req.State.Get(ctx, model)...)
	} else {
		// Note to maintainers - If this error is reached, there has been an issue converting the Resource Identity into
		// the Azure Resource ID for the resource
		SetResponseErrorDiagnostic(resp, "Current State Error", "Current State was null for read decode")
	}
}

// DecodeDataSourceRead reads a Data Sources config from a datasource.ReadRequest into a pointer to a target model and
// sets datasource.ReadResponse diags on error.
func (r *ResourceMetadata) DecodeDataSourceRead(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse, model any) {
	resp.Diagnostics.Append(req.Config.Get(ctx, model)...)
}

// EncodeRead writes the model populated in the Resource Read method to state.
func (r *ResourceMetadata) EncodeRead(ctx context.Context, resp *resource.ReadResponse, model any) {
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// EncodeDataSourceRead writes the model populated in the Read method to state.
func (r *ResourceMetadata) EncodeDataSourceRead(ctx context.Context, resp *datasource.ReadResponse, model any) {
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// DecodeUpdate reads a plan and state from a resource.UpdateRequest into pointers to the resource models and sets
// resource.UpdateResponse diags on error.
func (r *ResourceMetadata) DecodeUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, plan any, state any) {
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
}

// EncodeUpdate writes the model populated in the Update method to state.
func (r *ResourceMetadata) EncodeUpdate(ctx context.Context, resp *resource.UpdateResponse, state any) {
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// DecodeDelete reads a resources State from a resource.ReadRequest into a pointer to a target model and sets resource.ReadResponse diags on error.
func (r *ResourceMetadata) DecodeDelete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, state any) {
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
}

// SetResponseErrorDiagnostic is a helper function to write an Error Diagnostic to the appropriate Framework response
// type detail can be specified as an error, from which error.Error() will be used or as a string
// Note: For list resource diagnostics, pass in the stream itself, not the stream.Results for resp.
func SetResponseErrorDiagnostic(resp any, summary string, detail any) {
	var errorMsg string
	switch e := detail.(type) {
	case error:
		errorMsg = e.Error()
	case string:
		errorMsg = e
	}
	switch v := resp.(type) {
	case *resource.CreateResponse:
		v.Diagnostics.AddError(summary, errorMsg)
	case *resource.UpdateResponse:
		v.Diagnostics.AddError(summary, errorMsg)
	case *resource.DeleteResponse:
		v.Diagnostics.AddError(summary, errorMsg)
	case *resource.ReadResponse:
		v.Diagnostics.AddError(summary, errorMsg)
	case *ephemeral.OpenResponse:
		v.Diagnostics.AddError(summary, errorMsg)
	case *ephemeral.RenewResponse:
		v.Diagnostics.AddError(summary, errorMsg)
	case *ephemeral.CloseResponse:
		v.Diagnostics.AddError(summary, errorMsg)
	case *action.InvokeResponse:
		v.Diagnostics.AddError(summary, errorMsg)
	case *list.ListResultsStream:
		diags := diag.Diagnostics{}
		diags.Append(diag.NewErrorDiagnostic(summary, errorMsg))
		v.Results = list.ListResultsStreamDiagnostics(diags)
	}
}

// SetResponseWarningDiagnostic is a helper function to write an Error Diagnostic to the appropriate Framework response
// type detail can be specified as an error, from which error.Error() will be used or as a string
// Note: For list resource diagnostics, pass in the stream itself, not the stream.Results for resp.
func SetResponseWarningDiagnostic(resp any, summary string, detail any) {
	var errorMsg string
	switch e := detail.(type) {
	case error:
		errorMsg = e.Error()
	case string:
		errorMsg = e
	}
	switch v := resp.(type) {
	case *resource.CreateResponse:
		v.Diagnostics.AddWarning(summary, errorMsg)
	case *resource.UpdateResponse:
		v.Diagnostics.AddWarning(summary, errorMsg)
	case *resource.DeleteResponse:
		v.Diagnostics.AddWarning(summary, errorMsg)
	case *resource.ReadResponse:
		v.Diagnostics.AddWarning(summary, errorMsg)
	case *ephemeral.OpenResponse:
		v.Diagnostics.AddWarning(summary, errorMsg)
	case *ephemeral.RenewResponse:
		v.Diagnostics.AddWarning(summary, errorMsg)
	case *ephemeral.CloseResponse:
		v.Diagnostics.AddWarning(summary, errorMsg)
	case *list.ListResultsStream:
		diags := diag.Diagnostics{}
		diags.Append(diag.NewWarningDiagnostic(summary, errorMsg))
		v.Results = list.ListResultsStreamDiagnostics(diags)
	}
}

// AppendResponseErrorDiagnostic is a helper function to write an Error Diagnostic to the appropriate Framework response
// type detail can be specified as an error, from which error.Error() will be used or as a string
func AppendResponseErrorDiagnostic(resp any, d diag.Diagnostics) {
	switch v := resp.(type) {
	case *resource.ConfigureResponse:
		v.Diagnostics.Append(d...)
	case *resource.CreateResponse:
		v.Diagnostics.Append(d...)
	case *resource.UpdateResponse:
		v.Diagnostics.Append(d...)
	case *resource.DeleteResponse:
		v.Diagnostics.Append(d...)
	case *resource.ReadResponse:
		v.Diagnostics.Append(d...)
	case *resource.ValidateConfigResponse:
		v.Diagnostics.Append(d...)
	case *datasource.ConfigureResponse:
		v.Diagnostics.Append(d...)
	case *datasource.ValidateConfigResponse:
		v.Diagnostics.Append(d...)
	case *datasource.ReadResponse:
		v.Diagnostics.Append(d...)
	}
}

type FrameworkWrappedResource interface {
	ModelObject() any

	ResourceType() string

	Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse)

	Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse, metadata ResourceMetadata, plan any)

	Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse, metadata ResourceMetadata, state any)

	Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse, metadata ResourceMetadata, state any)

	ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse, metadata ResourceMetadata)

	Identity() (id resourceids.ResourceId, idType ResourceTypeForIdentity)
}

// FrameworkWrappedResourceWithUpdate provides an extension to the base resource for resources that can be updated.
type FrameworkWrappedResourceWithUpdate interface {
	FrameworkWrappedResource

	Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse, metadata ResourceMetadata, plan any, state any)
}

// FrameworkWrappedResourceWithConfigure provides an interface for resources that need custom configuration beyond the
// standard wrapped Configure() which configures the resource metadata.
type FrameworkWrappedResourceWithConfigure interface {
	FrameworkWrappedResource

	Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse, metadata ResourceMetadata)
}

// FrameworkWrappedResourceWithConfigValidators provides an interface for resources that need custom or complex
// validation logic based on the supplied user config, whole or in part.
type FrameworkWrappedResourceWithConfigValidators interface {
	FrameworkWrappedResource

	ConfigValidators(ctx context.Context) []resource.ConfigValidator
}

// FrameworkWrappedResourceWithPlanModifier provides an interface for resources that require Plan Modification
// Plan modifiers happen after validators and before create.
type FrameworkWrappedResourceWithPlanModifier interface {
	FrameworkWrappedResource

	ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse, metadata ResourceMetadata)
}

type FrameworkWrappedResourceWithList interface {
	FrameworkWrappedResource

	List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata ResourceMetadata)

	ListResourceConfigSchema(ctx context.Context, request list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse, metadata ResourceMetadata)
}
