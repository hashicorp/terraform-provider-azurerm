// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/taskhubs"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	TaskHubListResource struct{}
	TaskHubListModel    struct {
		SchedulerId types.String `tfsdk:"scheduler_id"`
	}
)

var _ sdk.FrameworkListWrappedResourceWithConfig = new(TaskHubListResource)

func (TaskHubListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = TaskHubResource{}.ResourceType()
}

func (TaskHubListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(TaskHubResource{})
}

func (TaskHubListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"scheduler_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: schedulers.ValidateSchedulerID,
					},
				},
			},
		},
	}
}

func (TaskHubListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DurableTask.TaskHubsClient

	var data TaskHubListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	schedulerId, err := schedulers.ParseSchedulerID(data.SchedulerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Scheduler ID for `%s`", TaskHubResource{}.ResourceType()), err)
		return
	}

	taskHubSchedulerId := taskhubs.NewSchedulerID(schedulerId.SubscriptionId, schedulerId.ResourceGroupName, schedulerId.SchedulerName)

	resp, err := client.ListBySchedulerComplete(ctx, taskHubSchedulerId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", TaskHubResource{}.ResourceType()), err)
		return
	}

	r := TaskHubResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := taskhubs.ParseTaskHubIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Durable Task Hub ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			state := TaskHubResourceModel{
				Name:        id.TaskHubName,
				SchedulerId: schedulerId.ID(),
			}

			if props := item.Properties; props != nil {
				state.DashboardUrl = pointer.From(props.DashboardURL)
			}

			if err := meta.Encode(&state); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", pointer.From(item.Name)), err)
				return
			}

			sdk.EncodeListResult(ctx, meta.ResourceData, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
