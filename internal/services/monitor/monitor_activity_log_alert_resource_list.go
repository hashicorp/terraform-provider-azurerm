// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const monitorActivityLogAlertResourceName = "azurerm_monitor_activity_log_alert"

type MonitorActivityLogAlertListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(MonitorActivityLogAlertListResource)

func (MonitorActivityLogAlertListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMonitorActivityLogAlert()
}

func (MonitorActivityLogAlertListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = monitorActivityLogAlertResourceName
}

func (MonitorActivityLogAlertListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Monitor.ActivityLogAlertsClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	results := make([]activitylogalertsapis.ActivityLogAlertResource, 0)

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ActivityLogAlertsListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s` by resource group", monitorActivityLogAlertResourceName), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.ActivityLogAlertsListBySubscriptionIdComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s` by subscription", monitorActivityLogAlertResourceName), err)
			return
		}
		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := activitylogalertsapis.ParseActivityLogAlertIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Activity Log Alert ID", err)
				return
			}

			rd := resourceMonitorActivityLogAlert().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceMonitorActivityLogAlertFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data for %s", monitorActivityLogAlertResourceName, pointer.From(item.Name)), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
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
