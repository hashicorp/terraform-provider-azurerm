// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert_v2

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *scheduledqueryrules.ScheduledQueryRuleId, model *scheduledqueryrules.ScheduledQueryRuleResource) error {
	state := ScheduledQueryRulesAlertV2Model{
		Name:              id.ScheduledQueryRuleName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening SystemOrUserAssigned Identity: %+v", err)
		}

		state.Location = location.Normalize(model.Location)
		state.Identity = *flattenedIdentity

		properties := &model.Properties
		state.Actions = flattenScheduledQueryRulesAlertV2ActionsModel(properties.Actions)
		state.AutoMitigate = pointer.From(properties.AutoMitigate)
		state.CheckWorkspaceAlertsStorageConfigured = pointer.From(properties.CheckWorkspaceAlertsStorageConfigured)
		state.CreatedWithApiVersion = pointer.From(properties.CreatedWithApiVersion)
		state.Criteria = flattenScheduledQueryRulesAlertV2CriteriaModel(properties.Criteria)
		state.Description = pointer.From(properties.Description)
		state.DisplayName = pointer.From(properties.DisplayName)
		state.Enabled = pointer.From(properties.Enabled)
		state.EvaluationFrequency = pointer.From(properties.EvaluationFrequency)
		state.IsLegacyLogAnalyticsRule = pointer.From(properties.IsLegacyLogAnalyticsRule)
		state.IsWorkspaceAlertsStorageConfigured = pointer.From(properties.IsWorkspaceAlertsStorageConfigured)
		state.MuteActionsDuration = pointer.From(properties.MuteActionsDuration)
		state.OverrideQueryTimeRange = pointer.From(properties.OverrideQueryTimeRange)
		state.Scopes = pointer.From(properties.Scopes)
		state.Severity = pointer.From(properties.Severity)
		state.SkipQueryValidation = pointer.From(properties.SkipQueryValidation)
		state.TargetResourceTypes = pointer.From(properties.TargetResourceTypes)
		state.WindowSize = pointer.From(properties.WindowSize)
		state.Tags = pointer.From(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenScheduledQueryRulesAlertV2ActionsModel(input *scheduledqueryrules.Actions) []ScheduledQueryRulesAlertV2ActionsModel {
	var outputList []ScheduledQueryRulesAlertV2ActionsModel
	if input == nil {
		return []ScheduledQueryRulesAlertV2ActionsModel{}
	}

	output := ScheduledQueryRulesAlertV2ActionsModel{}

	if input.ActionGroups != nil {
		output.ActionGroups = *input.ActionGroups
	}

	if input.CustomProperties != nil {
		output.CustomProperties = *input.CustomProperties
	}

	if input.ActionProperties != nil {
		if s, ok := (*input.ActionProperties)["Email.Subject"]; ok {
			output.EmailSubject = s
		}
	}

	return append(outputList, output)
}

func flattenScheduledQueryRulesAlertV2CriteriaModel(input *scheduledqueryrules.ScheduledQueryRuleCriteria) []ScheduledQueryRulesAlertV2CriteriaModel {
	if input == nil || input.AllOf == nil {
		return []ScheduledQueryRulesAlertV2CriteriaModel{}
	}

	outputList := make([]ScheduledQueryRulesAlertV2CriteriaModel, 0, len(*input.AllOf))
	for _, v := range *input.AllOf {
		output := ScheduledQueryRulesAlertV2CriteriaModel{}

		output.Dimensions = flattenScheduledQueryRulesAlertV2DimensionModel(v.Dimensions)
		output.FailingPeriods = flattenScheduledQueryRulesAlertV2FailingPeriodsModel(v.FailingPeriods)

		if v.MetricMeasureColumn != nil {
			output.MetricMeasureColumn = *v.MetricMeasureColumn
		}

		if v.Operator != nil {
			output.Operator = *v.Operator
		}

		if v.Query != nil {
			output.Query = *v.Query
		}

		if v.ResourceIdColumn != nil {
			output.ResourceIdColumn = *v.ResourceIdColumn
		}

		if v.Threshold != nil {
			output.Threshold = *v.Threshold
		}

		if v.TimeAggregation != nil {
			output.TimeAggregation = *v.TimeAggregation
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenScheduledQueryRulesAlertV2DimensionModel(inputList *[]scheduledqueryrules.Dimension) []ScheduledQueryRulesAlertV2DimensionModel {
	if inputList == nil {
		return []ScheduledQueryRulesAlertV2DimensionModel{}
	}

	outputList := make([]ScheduledQueryRulesAlertV2DimensionModel, 0, len(*inputList))
	for _, input := range *inputList {
		output := ScheduledQueryRulesAlertV2DimensionModel{
			Name:     input.Name,
			Operator: input.Operator,
			Values:   input.Values,
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenScheduledQueryRulesAlertV2FailingPeriodsModel(input *scheduledqueryrules.ConditionFailingPeriods) []ScheduledQueryRulesAlertV2FailingPeriodsModel {
	var outputList []ScheduledQueryRulesAlertV2FailingPeriodsModel
	if input == nil {
		return []ScheduledQueryRulesAlertV2FailingPeriodsModel{}
	}

	output := ScheduledQueryRulesAlertV2FailingPeriodsModel{}

	if input.MinFailingPeriodsToAlert != nil {
		output.MinFailingPeriodsToAlert = *input.MinFailingPeriodsToAlert
	}

	if input.NumberOfEvaluationPeriods != nil {
		output.NumberOfEvaluationPeriods = *input.NumberOfEvaluationPeriods
	}

	return append(outputList, output)
}
