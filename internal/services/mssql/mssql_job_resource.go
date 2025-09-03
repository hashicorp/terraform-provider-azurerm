package mssql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/commonschema"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobs"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MsSqlJobResource struct{}

var _ sdk.FrameworkWrappedResourceWithUpdate = &MsSqlJobResource{}

func (r MsSqlJobResource) ModelObject() any {
	return new(MsSqlJobResourceModel)
}

func (r MsSqlJobResource) ResourceType() string {
	return "azurerm_mssql_job"
}

func (r MsSqlJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, metadata sdk.ResourceMetadata) {
	if req.ID == "" {
		resourceIdentity := &MsSqlJobResourceIdentityModel{}
		req.Identity.Get(ctx, resourceIdentity)
		id := pointer.To(jobs.NewJobID(resourceIdentity.SubscriptionId, resourceIdentity.ResourceGroupName, resourceIdentity.ServerName, resourceIdentity.JobAgentName, resourceIdentity.Name))
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id.ID())...)
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r MsSqlJobResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			commonschema.Name: schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringIsNotEmpty,
					},
				},
			},

			"job_agent_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: jobs.ValidateJobAgentID,
					},
				},
			},

			"description": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r MsSqlJobResource) Create(ctx context.Context, _ resource.CreateRequest, resp *resource.CreateResponse, metadata sdk.ResourceMetadata, decodedPlan any) {
	client := metadata.Client.MSSQL.JobsClient
	subscriptionID := metadata.Client.Account.SubscriptionId

	data := sdk.AssertResourceModelType[MsSqlJobResourceModel](decodedPlan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	jobAgentID, err := jobs.ParseJobAgentID(data.JobAgentID.ValueString())
	if err != nil {
		// TODO: what do we want to use as the summary for resource ID parsing errors? I think this should be standardised for consistency while we're migrating to FW
		// - "ID Parsing Error" -- can be misleading if it's parsing a parent ID rather than resource's ID
		// - "ID Parsing Error For `property`" -- meh
		// - "`job_agent_id`" -- just the property that's being parsed? `err` already contains the relevant info I think
		sdk.SetResponseErrorDiagnostic(resp, "`job_agent_id`", err)
	}

	id := jobs.NewJobID(subscriptionID, jobAgentID.ResourceGroupName, jobAgentID.ServerName, jobAgentID.JobAgentName, data.Name.ValueString())

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("checking for presence of existing %s: %+v", id, err), err)
			return
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		metadata.ResourceRequiresImport(r.ResourceType(), id, resp)
		return
	}

	payload := r.buildPayload(data, nil)

	if _, err = client.CreateOrUpdate(ctx, id, payload); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("creating %s:", id), err)
		return
	}

	data.ID = types.StringValue(id.ID())
}

func (r MsSqlJobResource) Read(ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse, metadata sdk.ResourceMetadata, decodedState any) {
	client := metadata.Client.MSSQL.JobsClient

	state := sdk.AssertResourceModelType[MsSqlJobResourceModel](decodedState, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := jobs.ParseJobID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "ID parsing error", err)
		return
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			metadata.MarkAsGone(id, &resp.State, &resp.Diagnostics)
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id), err)
		return
	}

	state.Name = types.StringValue(id.JobName)
	state.JobAgentID = types.StringValue(jobs.NewJobAgentID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.JobAgentName).ID())

	if model := existing.Model; model != nil {
		if props := model.Properties; props != nil {
			// API returns an empty string regardless
			// and Terraform will show a non-empty plan with no visible changes if we store an empty string?
			if pointer.From(props.Description) != "" {
				state.Description = types.StringPointerValue(props.Description)
			}
		}
	}
}

func (r MsSqlJobResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse, metadata sdk.ResourceMetadata, decodedPlan any, decodedState any) {
	client := metadata.Client.MSSQL.JobsClient

	plan := sdk.AssertResourceModelType[MsSqlJobResourceModel](decodedPlan, resp)
	state := sdk.AssertResourceModelType[MsSqlJobResourceModel](decodedState, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := jobs.ParseJobID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "ID parsing error", err)
		return
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			metadata.MarkAsGone(id, &resp.State, &resp.Diagnostics)
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id), err)
		return
	}

	// TODO: is there a `HasChange` equivalent for FW / do we even need it?
	payload := r.buildPayload(plan, existing.Model)

	if _, err = client.CreateOrUpdate(ctx, *id, payload); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("updating %s:", id), err)
		return
	}
}

func (r MsSqlJobResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse, metadata sdk.ResourceMetadata, decodedState any) {
	client := metadata.Client.MSSQL.JobsClient

	state := sdk.AssertResourceModelType[MsSqlJobResourceModel](decodedState, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := jobs.ParseJobID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "ID parsing error", err)
		return
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("deleting %s:", *id), err)
	}
}

func (r MsSqlJobResource) Identity() (id resourceids.ResourceId, idType sdk.ResourceTypeForIdentity) {
	return &jobs.JobId{}, sdk.ResourceTypeForIdentityDefault
}

func (r MsSqlJobResource) buildPayload(data *MsSqlJobResourceModel, existing *jobs.Job) jobs.Job {
	// TODO: should create/update be separate? / is using existing even required anymore?
	// current implementation (as in sdkv2) we use the config to update the existing model
	// do we just default to pulling everything from `data` regardless of what's in existing?
	// Or always pass `existing` (rename to `payload`?)
	if existing == nil {
		return jobs.Job{
			Name: data.Name.ValueStringPointer(),
			Properties: &jobs.JobProperties{
				Description: data.Description.ValueStringPointer(),
			},
		}
	}

	// Only updatable property is `description`
	if existing.Properties == nil {
		existing.Properties = &jobs.JobProperties{}
	}

	existing.Properties.Description = data.Description.ValueStringPointer()

	return *existing
}
