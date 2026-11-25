// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_tag_rules

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.TagRulesClient
			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if model := resp.Model; model != nil {
				props := model.Properties
				monitorId := monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName)

				state := TagRulesResourceModel{
					Name:        id.TagRuleName,
					Monitor:     monitorId.ID(),
					LogRules:    FlattenLogRules(props.LogRules),
					MetricRules: FlattenMetricRules(props.MetricRules),
				}

				return metadata.Encode(&state)
			}

			return nil
		},
	}
}

type MetricRule struct {
	FilteringTags  []FilteringTag `tfschema:"filtering_tag"`
	SendingMetrics bool           `tfschema:"sending_metrics_enabled"`
}

type LogRule struct {
	FilteringTags        []FilteringTag `tfschema:"filtering_tag"`
	SendAadLogs          bool           `tfschema:"send_azure_active_directory_logs_enabled"`
	SendActivityLogs     bool           `tfschema:"send_activity_logs_enabled"`
	SendSubscriptionLogs bool           `tfschema:"send_subscription_logs_enabled"`
}

type FilteringTag struct {
	Name   string `tfschema:"name"`
	Value  string `tfschema:"value"`
	Action string `tfschema:"action"`
}
