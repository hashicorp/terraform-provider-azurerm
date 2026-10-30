// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package new_relic_tag_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.TagRulesClient

			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := NewRelicTagRuleModel{
				NewRelicMonitorId: monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName).ID(),
			}

			if model := resp.Model; model != nil {
				properties := &model.Properties
				if properties.LogRules != nil {
					state.AadLogEnabled = properties.LogRules.SendAadLogs != nil && *properties.LogRules.SendAadLogs == tagrules.SendAadLogsStatusEnabled
					state.ActivityLogEnabled = properties.LogRules.SendActivityLogs != nil && *properties.LogRules.SendActivityLogs == tagrules.SendActivityLogsStatusEnabled
					state.LogTagFilter = flattenFilteringTagModelArray(properties.LogRules.FilteringTags)
					state.SubscriptionLogEnabled = properties.LogRules.SendSubscriptionLogs != nil && *properties.LogRules.SendSubscriptionLogs == tagrules.SendSubscriptionLogsStatusEnabled
				}

				if properties.MetricRules != nil {
					state.MetricEnabled = properties.MetricRules.SendMetrics != nil && *properties.MetricRules.SendMetrics == tagrules.SendMetricsStatusEnabled
					state.MetricTagFilter = flattenFilteringTagModelArray(properties.MetricRules.FilteringTags)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenFilteringTagModelArray(inputList *[]tagrules.FilteringTag) []FilteringTagModel {
	outputList := make([]FilteringTagModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		outputList = append(outputList, FilteringTagModel{
			Action: pointer.From(input.Action),
			Name:   pointer.From(input.Name),
			Value:  pointer.From(input.Value),
		})
	}

	return outputList
}

func (r Resource) getEmail(ctx context.Context, monitorClient *monitors.MonitorsClient, monitorId *monitors.MonitorId) (string, error) {
	monitor, err := monitorClient.Get(ctx, *monitorId)
	if err != nil {
		return "", fmt.Errorf("getting monitor: %+v", err)
	}

	if monitor.Model == nil || monitor.Model.Properties.UserInfo == nil || monitor.Model.Properties.UserInfo.EmailAddress == nil || *monitor.Model.Properties.UserInfo.EmailAddress == "" {
		return "", fmt.Errorf("failed to get user email address from monitor")
	}

	return *monitor.Model.Properties.UserInfo.EmailAddress, nil
}
