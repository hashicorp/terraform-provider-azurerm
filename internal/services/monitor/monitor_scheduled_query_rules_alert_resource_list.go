// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const monitorScheduledQueryRulesAlertResourceName = "azurerm_monitor_scheduled_query_rules_alert"

type MonitorScheduledQueryRulesAlertListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(MonitorScheduledQueryRulesAlertListResource)

func (r MonitorScheduledQueryRulesAlertListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMonitorScheduledQueryRulesAlert()
}

func (r MonitorScheduledQueryRulesAlertListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = monitorScheduledQueryRulesAlertResourceName
}

func (r MonitorScheduledQueryRulesAlertListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Monitor.ScheduledQueryRulesClient

	var data sdk.DefaultListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	results := make([]scheduledqueryrules.LogSearchRuleResource, 0)

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroup(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()), scheduledqueryrules.DefaultListByResourceGroupOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", monitorScheduledQueryRulesAlertResourceName), err)
			return
		}

		if resp.Model != nil && resp.Model.Value != nil {
			results = *resp.Model.Value
		}
	default:
		resp, err := client.ListBySubscription(ctx, commonids.NewSubscriptionID(subscriptionID), scheduledqueryrules.DefaultListBySubscriptionOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", monitorScheduledQueryRulesAlertResourceName), err)
			return
		}

		if resp.Model != nil && resp.Model.Value != nil {
			results = *resp.Model.Value
		}
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			// The scheduledQueryRules API returns both `AlertingAction` rules (azurerm_monitor_scheduled_query_rules_alert)
			// and `LogToMetricAction` rules (azurerm_monitor_scheduled_query_rules_log) - only the former belong to this resource.
			if _, ok := item.Properties.Action.(scheduledqueryrules.AlertingAction); !ok {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := scheduledqueryrules.ParseScheduledQueryRuleIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Scheduled Query Rule ID", err)
				return
			}

			rd := resourceMonitorScheduledQueryRulesAlert().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceMonitorScheduledQueryRulesAlertFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", monitorScheduledQueryRulesAlertResourceName), err)
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
