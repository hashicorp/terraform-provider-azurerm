// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert_v2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ScheduledQueryRulesAlertV2Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Monitor.ScheduledQueryRulesV2Client
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := scheduledqueryrules.NewScheduledQueryRuleID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &scheduledqueryrules.ScheduledQueryRuleResource{
				Kind:     pointer.To(scheduledqueryrules.KindLogAlert),
				Location: location.Normalize(model.Location),
				Properties: scheduledqueryrules.ScheduledQueryRuleProperties{
					AutoMitigate:                          &model.AutoMitigate,
					CheckWorkspaceAlertsStorageConfigured: &model.CheckWorkspaceAlertsStorageConfigured,
					Enabled:                               &model.Enabled,
					Scopes:                                &model.Scopes,
					Severity:                              &model.Severity,
					SkipQueryValidation:                   &model.SkipQueryValidation,
					TargetResourceTypes:                   &model.TargetResourceTypes,
				},
				Tags: &model.Tags,
			}

			properties.Properties.Actions = expandScheduledQueryRulesAlertV2ActionsModel(model.Actions)

			properties.Properties.Criteria = expandScheduledQueryRulesAlertV2CriteriaModel(model.Criteria)

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.DisplayName != "" {
				properties.Properties.DisplayName = &model.DisplayName
			}

			if model.EvaluationFrequency != "" {
				properties.Properties.EvaluationFrequency = &model.EvaluationFrequency
			}

			if len(model.Identity) != 0 {
				ExpandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding SystemOrUserAssigned Identity: %+v", err)
				}
				properties.Identity = ExpandedIdentity
			}

			if model.MuteActionsDuration != "" {
				if model.AutoMitigate {
					return fmt.Errorf("auto mitigation must be disabled when mute action duration is set")
				}
				properties.Properties.MuteActionsDuration = &model.MuteActionsDuration
			}

			if model.OverrideQueryTimeRange != "" {
				properties.Properties.OverrideQueryTimeRange = &model.OverrideQueryTimeRange
			}

			if model.WindowSize != "" {
				properties.Properties.WindowSize = &model.WindowSize
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandScheduledQueryRulesAlertV2ActionsModel(inputList []ScheduledQueryRulesAlertV2ActionsModel) *scheduledqueryrules.Actions {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := scheduledqueryrules.Actions{
		ActionGroups:     &input.ActionGroups,
		CustomProperties: &input.CustomProperties,
	}

	if input.EmailSubject != "" {
		m := map[string]string{
			"Email.Subject": input.EmailSubject,
		}
		output.ActionProperties = &m
	} else {
		output.ActionProperties = pointer.To(map[string]string{})
	}

	return &output
}

func expandScheduledQueryRulesAlertV2CriteriaModel(inputList []ScheduledQueryRulesAlertV2CriteriaModel) *scheduledqueryrules.ScheduledQueryRuleCriteria {
	output := scheduledqueryrules.ScheduledQueryRuleCriteria{}
	outputList := make([]scheduledqueryrules.Condition, 0, len(inputList))
	for _, v := range inputList {
		input := v
		condition := scheduledqueryrules.Condition{
			Operator:        &input.Operator,
			Threshold:       &input.Threshold,
			TimeAggregation: &input.TimeAggregation,
		}

		condition.Dimensions = expandScheduledQueryRulesAlertV2DimensionModel(input.Dimensions)
		condition.FailingPeriods = expandScheduledQueryRulesAlertV2FailingPeriodsModel(input.FailingPeriods)

		if input.MetricMeasureColumn != "" {
			condition.MetricMeasureColumn = &input.MetricMeasureColumn
		}

		if input.Query != "" {
			condition.Query = &input.Query
		}

		if input.ResourceIdColumn != "" {
			condition.ResourceIdColumn = &input.ResourceIdColumn
		}

		outputList = append(outputList, condition)
	}
	output.AllOf = &outputList
	return &output
}

func expandScheduledQueryRulesAlertV2DimensionModel(inputList []ScheduledQueryRulesAlertV2DimensionModel) *[]scheduledqueryrules.Dimension {
	outputList := make([]scheduledqueryrules.Dimension, 0, len(inputList))
	for _, v := range inputList {
		input := v
		output := scheduledqueryrules.Dimension{
			Name:     input.Name,
			Operator: input.Operator,
			Values:   input.Values,
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func expandScheduledQueryRulesAlertV2FailingPeriodsModel(inputList []ScheduledQueryRulesAlertV2FailingPeriodsModel) *scheduledqueryrules.ConditionFailingPeriods {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := scheduledqueryrules.ConditionFailingPeriods{
		MinFailingPeriodsToAlert:  &input.MinFailingPeriodsToAlert,
		NumberOfEvaluationPeriods: &input.NumberOfEvaluationPeriods,
	}

	return &output
}
