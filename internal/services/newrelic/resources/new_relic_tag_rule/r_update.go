// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package new_relic_tag_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.TagRulesClient

			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if properties.Properties.LogRules == nil {
				return fmt.Errorf("retrieving %s: log rules was nil", id)
			}

			if properties.Properties.MetricRules == nil {
				return fmt.Errorf("retrieving %s: metric rules was nil", id)
			}

			var model NewRelicTagRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			monitorId, err := monitors.ParseMonitorID(model.NewRelicMonitorId)
			if err != nil {
				return err
			}

			email, err := r.getEmail(ctx, metadata.Client.NewRelic.MonitorsClient, monitorId)
			if err != nil {
				return err
			}

			properties.Properties.MetricRules.UserEmail = &email

			if metadata.ResourceData.HasChange("azure_active_directory_log_enabled") {
				if model.AadLogEnabled {
					properties.Properties.LogRules.SendAadLogs = pointer.To(tagrules.SendAadLogsStatusEnabled)
				} else {
					properties.Properties.LogRules.SendAadLogs = pointer.To(tagrules.SendAadLogsStatusDisabled)
				}
			}

			if metadata.ResourceData.HasChange("activity_log_enabled") {
				if model.ActivityLogEnabled {
					properties.Properties.LogRules.SendActivityLogs = pointer.To(tagrules.SendActivityLogsStatusEnabled)
				} else {
					properties.Properties.LogRules.SendActivityLogs = pointer.To(tagrules.SendActivityLogsStatusDisabled)
				}
			}

			if metadata.ResourceData.HasChange("log_tag_filter") {
				properties.Properties.LogRules.FilteringTags = expandFilteringTagModelArray(model.LogTagFilter)
			}

			if metadata.ResourceData.HasChange("metric_enabled") {
				if model.MetricEnabled {
					properties.Properties.MetricRules.SendMetrics = pointer.To(tagrules.SendMetricsStatusEnabled)
				} else {
					properties.Properties.MetricRules.SendMetrics = pointer.To(tagrules.SendMetricsStatusDisabled)
				}
			}

			if metadata.ResourceData.HasChange("metric_tag_filter") {
				properties.Properties.MetricRules.FilteringTags = expandFilteringTagModelArray(model.MetricTagFilter)
			}

			if metadata.ResourceData.HasChange("subscription_log_enabled") {
				if model.SubscriptionLogEnabled {
					properties.Properties.LogRules.SendSubscriptionLogs = pointer.To(tagrules.SendSubscriptionLogsStatusEnabled)
				} else {
					properties.Properties.LogRules.SendSubscriptionLogs = pointer.To(tagrules.SendSubscriptionLogsStatusDisabled)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
