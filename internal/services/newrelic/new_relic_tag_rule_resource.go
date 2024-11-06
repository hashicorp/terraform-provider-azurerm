// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package newrelic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NewRelicTagRuleModel struct {
	NewRelicMonitorId      string              `tfschema:"monitor_id"`
	AadLogEnabled          bool                `tfschema:"azure_active_directory_log_enabled"`
	ActivityLogEnabled     bool                `tfschema:"activity_log_enabled"`
	LogTagFilter           []FilteringTagModel `tfschema:"log_tag_filter"`
	MetricEnabled          bool                `tfschema:"metric_enabled"`
	MetricTagFilter        []FilteringTagModel `tfschema:"metric_tag_filter"`
	SubscriptionLogEnabled bool                `tfschema:"subscription_log_enabled"`
}

type FilteringTagModel struct {
	Action tagrules.TagAction `tfschema:"action"`
	Name   string             `tfschema:"name"`
	Value  string             `tfschema:"value"`
}

type NewRelicTagRuleResource struct{}

var _ sdk.ResourceWithUpdate = NewRelicTagRuleResource{}

func (r NewRelicTagRuleResource) ResourceType() string {
	return "azurerm_new_relic_tag_rule"
}

func (r NewRelicTagRuleResource) ModelObject() interface{} {
	return &NewRelicTagRuleModel{}
}

func (r NewRelicTagRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tagrules.ValidateTagRuleID
}

func (r NewRelicTagRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"monitor_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: monitors.ValidateMonitorID,
		},

		"azure_active_directory_log_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"activity_log_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"log_tag_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(tagrules.TagActionExclude),
							string(tagrules.TagActionInclude),
						}, false),
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"metric_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"metric_tag_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(tagrules.TagActionExclude),
							string(tagrules.TagActionInclude),
						}, false),
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"subscription_log_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r NewRelicTagRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NewRelicTagRuleResource) Create() sdk.ResourceFunc {
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
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NewRelicTagRuleResource) Update() sdk.ResourceFunc {
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

func (r NewRelicTagRuleResource) Read() sdk.ResourceFunc {
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

func (r NewRelicTagRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.TagRulesClient

			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r NewRelicTagRuleResource) getEmail(ctx context.Context, monitorClient *monitors.MonitorsClient, monitorId *monitors.MonitorId) (string, error) {
	monitor, err := monitorClient.Get(ctx, *monitorId)
	if err != nil {
		return "", fmt.Errorf("getting monitor: %+v", err)
	}

	if monitor.Model == nil || monitor.Model.Properties.UserInfo == nil || monitor.Model.Properties.UserInfo.EmailAddress == nil || *monitor.Model.Properties.UserInfo.EmailAddress == "" {
		return "", fmt.Errorf("failed to get user email address from monitor")
	}

	return *monitor.Model.Properties.UserInfo.EmailAddress, nil
}

func expandFilteringTagModelArray(inputList []FilteringTagModel) *[]tagrules.FilteringTag {
	var outputList []tagrules.FilteringTag
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
