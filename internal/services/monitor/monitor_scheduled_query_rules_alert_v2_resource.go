// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-08-01/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ScheduledQueryRulesAlertV2Model struct {
	Name                                  string                                    `tfschema:"name"`
	ResourceGroupName                     string                                    `tfschema:"resource_group_name"`
	Actions                               []ScheduledQueryRulesAlertV2ActionsModel  `tfschema:"action"`
	AutoMitigate                          bool                                      `tfschema:"auto_mitigation_enabled"`
	CheckWorkspaceAlertsStorageConfigured bool                                      `tfschema:"workspace_alerts_storage_enabled"`
	Criteria                              []ScheduledQueryRulesAlertV2CriteriaModel `tfschema:"criteria"`
	Description                           string                                    `tfschema:"description"`
	DisplayName                           string                                    `tfschema:"display_name"`
	Enabled                               bool                                      `tfschema:"enabled"`
	EvaluationFrequency                   string                                    `tfschema:"evaluation_frequency"`
	Location                              string                                    `tfschema:"location"`
	MuteActionsDuration                   string                                    `tfschema:"mute_actions_after_alert_duration"`
	OverrideQueryTimeRange                string                                    `tfschema:"query_time_range_override"`
	Scopes                                []string                                  `tfschema:"scopes"`
	Severity                              scheduledqueryrules.AlertSeverity         `tfschema:"severity"`
	SkipQueryValidation                   bool                                      `tfschema:"skip_query_validation"`
	Tags                                  map[string]string                         `tfschema:"tags"`
	TargetResourceTypes                   []string                                  `tfschema:"target_resource_types"`
	WindowSize                            string                                    `tfschema:"window_duration"`
	CreatedWithApiVersion                 string                                    `tfschema:"created_with_api_version"`
	IsLegacyLogAnalyticsRule              bool                                      `tfschema:"is_a_legacy_log_analytics_rule"`
	IsWorkspaceAlertsStorageConfigured    bool                                      `tfschema:"is_workspace_alerts_storage_configured"`
}

type ScheduledQueryRulesAlertV2ActionsModel struct {
	ActionGroups     []string          `tfschema:"action_groups"`
	CustomProperties map[string]string `tfschema:"custom_properties"`
}

type ScheduledQueryRulesAlertV2CriteriaModel struct {
	Dimensions          []ScheduledQueryRulesAlertV2DimensionModel      `tfschema:"dimension"`
	FailingPeriods      []ScheduledQueryRulesAlertV2FailingPeriodsModel `tfschema:"failing_periods"`
	MetricMeasureColumn string                                          `tfschema:"metric_measure_column"`
	Operator            scheduledqueryrules.ConditionOperator           `tfschema:"operator"`
	Query               string                                          `tfschema:"query"`
	ResourceIdColumn    string                                          `tfschema:"resource_id_column"`
	Threshold           float64                                         `tfschema:"threshold"`
	TimeAggregation     scheduledqueryrules.TimeAggregation             `tfschema:"time_aggregation_method"`
}

type ScheduledQueryRulesAlertV2DimensionModel struct {
	Name     string                                `tfschema:"name"`
	Operator scheduledqueryrules.DimensionOperator `tfschema:"operator"`
	Values   []string                              `tfschema:"values"`
}

type ScheduledQueryRulesAlertV2FailingPeriodsModel struct {
	MinFailingPeriodsToAlert  int64 `tfschema:"minimum_failing_periods_to_trigger_alert"`
	NumberOfEvaluationPeriods int64 `tfschema:"number_of_evaluation_periods"`
}

type ScheduledQueryRulesAlertV2Resource struct{}

var _ sdk.ResourceWithUpdate = ScheduledQueryRulesAlertV2Resource{}

func (r ScheduledQueryRulesAlertV2Resource) ResourceType() string {
	return "azurerm_monitor_scheduled_query_rules_alert_v2"
}

func (r ScheduledQueryRulesAlertV2Resource) ModelObject() interface{} {
	return &ScheduledQueryRulesAlertV2Model{}
}

func (r ScheduledQueryRulesAlertV2Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scheduledqueryrules.ValidateScheduledQueryRuleID
}

func (r ScheduledQueryRulesAlertV2Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"criteria": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{

				Schema: map[string]*pluginsdk.Schema{
					"query": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							// see https://github.com/Azure/azure-rest-api-specs/issues/21794
							"Equal",
							string(scheduledqueryrules.ConditionOperatorGreaterThan),
							string(scheduledqueryrules.ConditionOperatorGreaterThanOrEqual),
							string(scheduledqueryrules.ConditionOperatorLessThan),
							string(scheduledqueryrules.ConditionOperatorLessThanOrEqual),
						}, false),
					},

					"time_aggregation_method": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(scheduledqueryrules.TimeAggregationCount),
							string(scheduledqueryrules.TimeAggregationAverage),
							string(scheduledqueryrules.TimeAggregationMinimum),
							string(scheduledqueryrules.TimeAggregationMaximum),
							string(scheduledqueryrules.TimeAggregationTotal),
						}, false),
					},

					"threshold": {
						Type:     pluginsdk.TypeFloat,
						Required: true,
					},

					"dimension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(scheduledqueryrules.DimensionOperatorInclude),
										string(scheduledqueryrules.DimensionOperatorExclude),
									}, false),
								},

								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},

					"failing_periods": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"minimum_failing_periods_to_trigger_alert": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 6),
								},

								"number_of_evaluation_periods": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 6),
								},
							},
						},
					},

					"metric_measure_column": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"resource_id_column": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		// lintignore:S013
		"evaluation_frequency": {
			Type: pluginsdk.TypeString,
			// this field is required, missing this field will get an error from service
			Optional: !features.FourPointOhBeta(),
			Required: features.FourPointOhBeta(),
			ValidateFunc: validation.StringInSlice([]string{
				"PT1M",
				"PT5M",
				"PT10M",
				"PT15M",
				"PT30M",
				"PT45M",
				"PT1H",
				"PT2H",
				"PT3H",
				"PT4H",
				"PT5H",
				"PT6H",
				"P1D",
			}, false),
		},

		"scopes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},

		"severity": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"window_duration": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"PT1M",
				"PT5M",
				"PT10M",
				"PT15M",
				"PT30M",
				"PT45M",
				"PT1H",
				"PT2H",
				"PT3H",
				"PT4H",
				"PT5H",
				"PT6H",
				"P1D",
				"P2D",
			}, false),
		},

		"action": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},

					"custom_properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"auto_mitigation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"workspace_alerts_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"mute_actions_after_alert_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"PT5M",
				"PT10M",
				"PT15M",
				"PT30M",
				"PT45M",
				"PT1H",
				"PT2H",
				"PT3H",
				"PT4H",
				"PT5H",
				"PT6H",
				"P1D",
				"P2D",
			}, false),
		},

		"query_time_range_override": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"PT5M",
				"PT10M",
				"PT15M",
				"PT20M",
				"PT30M",
				"PT45M",
				"PT1H",
				"PT2H",
				"PT3H",
				"PT4H",
				"PT5H",
				"PT6H",
				"P1D",
				"P2D",
			}, false),
		},

		"skip_query_validation": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": commonschema.Tags(),

		"target_resource_types": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r ScheduledQueryRulesAlertV2Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"created_with_api_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"is_a_legacy_log_analytics_rule": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"is_workspace_alerts_storage_configured": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r ScheduledQueryRulesAlertV2Resource) Create() sdk.ResourceFunc {
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
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}
			kind := scheduledqueryrules.KindLogAlert
			properties := &scheduledqueryrules.ScheduledQueryRuleResource{
				Kind:     &kind,
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
			return nil
		},
	}
}

func (r ScheduledQueryRulesAlertV2Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.ScheduledQueryRulesV2Client

			id, err := scheduledqueryrules.ParseScheduledQueryRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var resourceModel ScheduledQueryRulesAlertV2Model
			if err := metadata.Decode(&resourceModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if metadata.ResourceData.HasChange("action") {
				model.Properties.Actions = expandScheduledQueryRulesAlertV2ActionsModel(resourceModel.Actions)
			}

			if metadata.ResourceData.HasChange("auto_mitigation_enabled") {
				model.Properties.AutoMitigate = &resourceModel.AutoMitigate
			}

			if metadata.ResourceData.HasChange("workspace_alerts_storage_enabled") {
				model.Properties.CheckWorkspaceAlertsStorageConfigured = &resourceModel.CheckWorkspaceAlertsStorageConfigured
			}

			if metadata.ResourceData.HasChange("criteria") {
				model.Properties.Criteria = expandScheduledQueryRulesAlertV2CriteriaModel(resourceModel.Criteria)
			}

			if metadata.ResourceData.HasChange("description") {
				model.Properties.Description = &resourceModel.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				if resourceModel.DisplayName != "" {
					model.Properties.DisplayName = &resourceModel.DisplayName
				} else {
					model.Properties.DisplayName = nil
				}
			}

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = &resourceModel.Enabled
			}

			if metadata.ResourceData.HasChange("evaluation_frequency") {
				model.Properties.EvaluationFrequency = &resourceModel.EvaluationFrequency
			}

			if metadata.ResourceData.HasChange("mute_actions_after_alert_duration") {
				if resourceModel.MuteActionsDuration != "" {
					if resourceModel.AutoMitigate {
						return fmt.Errorf("auto mitigation must be disabled when mute action duration is set")
					}
					model.Properties.MuteActionsDuration = &resourceModel.MuteActionsDuration
				} else {
					model.Properties.MuteActionsDuration = nil
				}
			}

			if metadata.ResourceData.HasChange("query_time_range_override") {
				if resourceModel.OverrideQueryTimeRange != "" {
					model.Properties.OverrideQueryTimeRange = &resourceModel.OverrideQueryTimeRange
				} else {
					model.Properties.OverrideQueryTimeRange = nil
				}
			}

			if metadata.ResourceData.HasChange("severity") {
				model.Properties.Severity = &resourceModel.Severity
			}

			if metadata.ResourceData.HasChange("skip_query_validation") {
				model.Properties.SkipQueryValidation = &resourceModel.SkipQueryValidation
			}

			if metadata.ResourceData.HasChange("target_resource_types") {
				model.Properties.TargetResourceTypes = &resourceModel.TargetResourceTypes
			}

			if metadata.ResourceData.HasChange("window_duration") {
				if resourceModel.WindowSize != "" {
					model.Properties.WindowSize = &resourceModel.WindowSize
				} else {
					model.Properties.WindowSize = nil
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &resourceModel.Tags
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ScheduledQueryRulesAlertV2Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.ScheduledQueryRulesV2Client

			id, err := scheduledqueryrules.ParseScheduledQueryRuleID(metadata.ResourceData.Id())
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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := ScheduledQueryRulesAlertV2Model{
				Name:              id.ScheduledQueryRuleName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties
			state.Actions = flattenScheduledQueryRulesAlertV2ActionsModel(properties.Actions)

			if properties.AutoMitigate != nil {
				state.AutoMitigate = *properties.AutoMitigate
			}

			if properties.CheckWorkspaceAlertsStorageConfigured != nil {
				state.CheckWorkspaceAlertsStorageConfigured = *properties.CheckWorkspaceAlertsStorageConfigured
			}

			if properties.CreatedWithApiVersion != nil {
				state.CreatedWithApiVersion = *properties.CreatedWithApiVersion
			}

			state.Criteria = flattenScheduledQueryRulesAlertV2CriteriaModel(properties.Criteria)

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.DisplayName != nil {
				state.DisplayName = *properties.DisplayName
			}

			if properties.Enabled != nil {
				state.Enabled = *properties.Enabled
			}

			if properties.EvaluationFrequency != nil {
				state.EvaluationFrequency = *properties.EvaluationFrequency
			}

			if properties.IsLegacyLogAnalyticsRule != nil {
				state.IsLegacyLogAnalyticsRule = *properties.IsLegacyLogAnalyticsRule
			}

			if properties.IsWorkspaceAlertsStorageConfigured != nil {
				state.IsWorkspaceAlertsStorageConfigured = *properties.IsWorkspaceAlertsStorageConfigured
			}

			if properties.MuteActionsDuration != nil {
				state.MuteActionsDuration = *properties.MuteActionsDuration
			}

			if properties.OverrideQueryTimeRange != nil {
				state.OverrideQueryTimeRange = *properties.OverrideQueryTimeRange
			}

			if properties.Scopes != nil {
				state.Scopes = *properties.Scopes
			}

			if properties.Severity != nil {
				state.Severity = *properties.Severity
			}

			if properties.SkipQueryValidation != nil {
				state.SkipQueryValidation = *properties.SkipQueryValidation
			}

			if properties.TargetResourceTypes != nil {
				state.TargetResourceTypes = *properties.TargetResourceTypes
			}

			if properties.WindowSize != nil {
				state.WindowSize = *properties.WindowSize
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ScheduledQueryRulesAlertV2Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.ScheduledQueryRulesV2Client

			id, err := scheduledqueryrules.ParseScheduledQueryRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
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

	return &output
}

func expandScheduledQueryRulesAlertV2CriteriaModel(inputList []ScheduledQueryRulesAlertV2CriteriaModel) *scheduledqueryrules.ScheduledQueryRuleCriteria {
	output := scheduledqueryrules.ScheduledQueryRuleCriteria{}
	var outputList []scheduledqueryrules.Condition
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
	var outputList []scheduledqueryrules.Dimension
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

func flattenScheduledQueryRulesAlertV2ActionsModel(input *scheduledqueryrules.Actions) []ScheduledQueryRulesAlertV2ActionsModel {
	var outputList []ScheduledQueryRulesAlertV2ActionsModel
	if input == nil {
		return outputList
	}

	output := ScheduledQueryRulesAlertV2ActionsModel{}

	if input.ActionGroups != nil {
		output.ActionGroups = *input.ActionGroups
	}

	if input.CustomProperties != nil {
		output.CustomProperties = *input.CustomProperties
	}

	return append(outputList, output)
}

func flattenScheduledQueryRulesAlertV2CriteriaModel(input *scheduledqueryrules.ScheduledQueryRuleCriteria) []ScheduledQueryRulesAlertV2CriteriaModel {
	var outputList []ScheduledQueryRulesAlertV2CriteriaModel
	if input == nil {
		return outputList
	}

	inputList := input.AllOf
	if inputList == nil {
		return outputList
	}

	for _, v := range *inputList {
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
	var outputList []ScheduledQueryRulesAlertV2DimensionModel
	if inputList == nil {
		return outputList
	}

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
		return outputList
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
