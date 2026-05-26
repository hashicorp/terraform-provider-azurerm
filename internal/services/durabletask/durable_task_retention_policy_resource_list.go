// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/retentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	RetentionPolicyListResource struct{}
	RetentionPolicyListModel    struct {
		DurableTaskSchedulerId types.String `tfsdk:"durable_task_scheduler_id"`
	}
)

var _ sdk.FrameworkListWrappedResourceWithConfig = new(RetentionPolicyListResource)

func (RetentionPolicyListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = RetentionPolicyResource{}.ResourceType()
}

func (RetentionPolicyListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(RetentionPolicyResource{})
}

func (RetentionPolicyListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"durable_task_scheduler_id": schema.StringAttribute{
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

func (RetentionPolicyListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DurableTask.RetentionPoliciesClient

	var data RetentionPolicyListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	schedulerId, err := schedulers.ParseSchedulerID(data.DurableTaskSchedulerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Scheduler ID for `%s`", RetentionPolicyResource{}.ResourceType()), err)
		return
	}

	retentionPolicySchedulerId := retentionpolicies.NewSchedulerID(schedulerId.SubscriptionId, schedulerId.ResourceGroupName, schedulerId.SchedulerName)

	resp, err := client.ListBySchedulerComplete(ctx, retentionPolicySchedulerId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", RetentionPolicyResource{}.ResourceType()), err)
		return
	}

	r := RetentionPolicyResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)

			id := NewRetentionPolicyID(schedulerId.SubscriptionId, schedulerId.ResourceGroupName, schedulerId.SchedulerName)

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			state := RetentionPolicyResourceModel{
				DurableTaskSchedulerId: schedulerId.ID(),
				RetentionPolicy:        make([]RetentionPolicyModel, 0),
			}

			if props := item.Properties; props != nil && props.RetentionPolicies != nil {
				state.RetentionPolicy = flattenRetentionPolicyDetails(props.RetentionPolicies)
			}

			if err := meta.Encode(&state); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
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
