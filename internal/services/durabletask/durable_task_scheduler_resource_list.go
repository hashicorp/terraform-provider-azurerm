// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SchedulerListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(SchedulerListResource)

func (SchedulerListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = SchedulerResource{}.ResourceType()
}

func (SchedulerListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(SchedulerResource{})
}

func (SchedulerListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DurableTask.SchedulersClient

	var data sdk.DefaultListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	var results []schedulers.Scheduler

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", SchedulerResource{}.ResourceType()), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", SchedulerResource{}.ResourceType()), err)
			return
		}
		results = resp.Items
	}

	r := SchedulerResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := schedulers.ParseSchedulerIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Durable Task Scheduler ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			state := SchedulerResourceModel{
				Name:              id.SchedulerName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(item.Location),
			}

			state.Tags = pointer.From(item.Tags)

			if props := item.Properties; props != nil {
				state.SkuName = string(props.Sku.Name)
				state.Capacity = pointer.From(props.Sku.Capacity)
				state.IpAllowList = props.IPAllowlist
				state.Endpoint = pointer.From(props.Endpoint)
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
