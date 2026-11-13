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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NewRelicTagRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NewRelic.TagRulesClient
			monitorId, err := monitors.ParseMonitorID(model.NewRelicMonitorId)
			if err != nil {
				return err
			}

			id := tagrules.NewTagRuleID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, "default")

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			logRules := tagrules.LogRules{
				FilteringTags:        expandFilteringTagModelArray(model.LogTagFilter),
				SendAadLogs:          pointer.To(tagrules.SendAadLogsStatusDisabled),
				SendActivityLogs:     pointer.To(tagrules.SendActivityLogsStatusDisabled),
				SendSubscriptionLogs: pointer.To(tagrules.SendSubscriptionLogsStatusDisabled),
			}

			email, err := r.getEmail(ctx, metadata.Client.NewRelic.MonitorsClient, monitorId)
			if err != nil {
				return err
			}

			metricRules := tagrules.MetricRules{
				FilteringTags: expandFilteringTagModelArray(model.MetricTagFilter),
				SendMetrics:   pointer.To(tagrules.SendMetricsStatusDisabled),
				UserEmail:     &email,
			}

			if model.AadLogEnabled {
				logRules.SendAadLogs = pointer.To(tagrules.SendAadLogsStatusEnabled)
			}

			if model.ActivityLogEnabled {
				logRules.SendActivityLogs = pointer.To(tagrules.SendActivityLogsStatusEnabled)
			}

			if model.SubscriptionLogEnabled {
				logRules.SendSubscriptionLogs = pointer.To(tagrules.SendSubscriptionLogsStatusEnabled)
			}

			if model.MetricEnabled {
				metricRules.SendMetrics = pointer.To(tagrules.SendMetricsStatusEnabled)
			}

			properties := &tagrules.TagRule{
				Properties: tagrules.MonitoringTagRulesProperties{
					LogRules:    &logRules,
					MetricRules: &metricRules,
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandFilteringTagModelArray(inputList []FilteringTagModel) *[]tagrules.FilteringTag {
	outputList := make([]tagrules.FilteringTag, 0, len(inputList))
	for _, v := range inputList {
		input := v
		output := tagrules.FilteringTag{
			Action: &input.Action,
		}

		if input.Name != "" {
			output.Name = &input.Name
		}

		if input.Value != "" {
			output.Value = &input.Value
		}
		outputList = append(outputList, output)
	}
	return &outputList
}
