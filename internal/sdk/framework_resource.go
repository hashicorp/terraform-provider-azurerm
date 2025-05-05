// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
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

// Defaults configures the Resource Metadata for client access, Provider Features, and subscriptionId.
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

// DecodeCreate reads a plan from a resource.CreateRequest into a pointer to a target model and sets resource.CreateResponse diags on error.
// Returns true if there are no error Diagnostics.
func (r *ResourceMetadata) DecodeCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan interface{}) bool {
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	return !resp.Diagnostics.HasError()
}

// EncodeCreate writes the Config passed to create to state.
func (r *ResourceMetadata) EncodeCreate(ctx context.Context, resp *resource.CreateResponse, config interface{}) {
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

// DecodeRead reads a resources State from a resource.ReadRequest into a pointer to a target model and sets resource.ReadResponse diags on error.
// Returns true if there are no error Diagnostics.
func (r *ResourceMetadata) DecodeRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, state interface{}) bool {
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	return !resp.Diagnostics.HasError()
}

// DecodeDataSourceRead reads a Data Sources config from a datasource.ReadRequest into a pointer to a target model and sets datasource.ReadResponse diags on error.
// Returns true if there are no error Diagnostics.
func (r *ResourceMetadata) DecodeDataSourceRead(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse, state interface{}) bool {
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)

	return !resp.Diagnostics.HasError()
}

// EncodeRead writes the state to an ReadResponse.
// The state parameter must be a pointer to a model for the resource. This should have been populated with all possible values read from the API.
func (r *ResourceMetadata) EncodeRead(ctx context.Context, resp *resource.ReadResponse, state interface{}) {
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// EncodeDataSourceRead writes the state to an ReadResponse.
// The state parameter must be a pointer to a model for the resource. This should have been populated with all possible values read from the API.
func (r *ResourceMetadata) EncodeDataSourceRead(ctx context.Context, resp *datasource.ReadResponse, state interface{}) {
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// DecodeUpdate reads a plan and state from a resource.UpdateRequest into pointers to a target models and sets resource.UpdateResponse diags on error.
// Returns true if there are no error Diagnostics.
// The plan and state parameters must be pointer to the model for the resource and should have been populated with the decoded plan and existing state prior to being passed to this function.
func (r *ResourceMetadata) DecodeUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, plan interface{}, state interface{}) bool {
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return false
	}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	return !resp.Diagnostics.HasError()
}

// EncodeUpdate writes the state back to an UpdateResponse.
// The state parameter must be a pointer to a model for the resource.
func (r *ResourceMetadata) EncodeUpdate(ctx context.Context, resp *resource.UpdateResponse, state interface{}) {
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// DecodeDelete reads a resources State from a resource.ReadRequest into a pointer to a target model and sets resource.ReadResponse diags on error.
// Returns true if there are no error Diagnostics.
func (r *ResourceMetadata) DecodeDelete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, state interface{}) {
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
}

// SetResponseErrorDiagnostic is a helper function to write an Error Diagnostic to the appropriate Framework response type
// detail can be specified as an error, from which error.Error() will be used or as a string
func SetResponseErrorDiagnostic(resp interface{}, summary string, detail interface{}) {
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
	}
}

// AppendResponseErrorDiagnostic is a helper function to write an Error Diagnostic to the appropriate Framework response type
// detail can be specified as an error, from which error.Error() will be used or as a string
func AppendResponseErrorDiagnostic(resp interface{}, d diag.Diagnostics) {
	switch v := resp.(type) {
	case *resource.CreateResponse:
		v.Diagnostics.Append(d...)
	case *resource.UpdateResponse:
		v.Diagnostics.Append(d...)
	case *resource.DeleteResponse:
		v.Diagnostics.Append(d...)
	case *resource.ReadResponse:
		v.Diagnostics.Append(d...)
	}
}

type FrameworkWrappedResource interface {
	ModelObject() interface{}

	ResourceType() string

	Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse)

	// New() *FrameworkResourceWrapper

	Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse, metadata ResourceMetadata, plan any)

	Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse, metadata ResourceMetadata, state interface{})

	Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse, metadata ResourceMetadata, state interface{})

	Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse, metadata ResourceMetadata, plan interface{}, state interface{})

	ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse, metadata ResourceMetadata)
}

// FrameworkWrappedResourceWithUpdate provides an interface for resources that need custom configuration beyond the standard wrapped Configure()
type FrameworkWrappedResourceWithUpdate interface {
	FrameworkWrappedResource

	Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse, metadata ResourceMetadata, plan interface{}, state interface{})
}

// FrameworkWrappedResourceWithConfigure provides an interface for resources that need custom configuration beyond the standard wrapped Configure()
type FrameworkWrappedResourceWithConfigure interface {
	FrameworkWrappedResource

	Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse, metadata ResourceMetadata)
}

// FrameworkWrappedResourceWithConfigValidators provides an interface for resources that need custom configuration beyond the standard wrapped Configure()
type FrameworkWrappedResourceWithConfigValidators interface {
	FrameworkWrappedResource

	ConfigValidators(ctx context.Context) []resource.ConfigValidator
}

// FrameworkWrappedResourceWithPlanModifier provides an interface for resources that need custom configuration beyond the standard wrapped Configure()
type FrameworkWrappedResourceWithPlanModifier interface {
	FrameworkWrappedResource

	ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse, metadata ResourceMetadata)
}
