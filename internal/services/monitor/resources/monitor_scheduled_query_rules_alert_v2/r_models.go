// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert_v2

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
)

type ScheduledQueryRulesAlertV2Model struct {
	Name                                  string                                     `tfschema:"name"`
	ResourceGroupName                     string                                     `tfschema:"resource_group_name"`
	Actions                               []ScheduledQueryRulesAlertV2ActionsModel   `tfschema:"action"`
	AutoMitigate                          bool                                       `tfschema:"auto_mitigation_enabled"`
	CheckWorkspaceAlertsStorageConfigured bool                                       `tfschema:"workspace_alerts_storage_enabled"`
	Criteria                              []ScheduledQueryRulesAlertV2CriteriaModel  `tfschema:"criteria"`
	Description                           string                                     `tfschema:"description"`
	DisplayName                           string                                     `tfschema:"display_name"`
	Enabled                               bool                                       `tfschema:"enabled"`
	EvaluationFrequency                   string                                     `tfschema:"evaluation_frequency"`
	Location                              string                                     `tfschema:"location"`
	MuteActionsDuration                   string                                     `tfschema:"mute_actions_after_alert_duration"`
	OverrideQueryTimeRange                string                                     `tfschema:"query_time_range_override"`
	Scopes                                []string                                   `tfschema:"scopes"`
	Severity                              scheduledqueryrules.AlertSeverity          `tfschema:"severity"`
	SkipQueryValidation                   bool                                       `tfschema:"skip_query_validation"`
	Tags                                  map[string]string                          `tfschema:"tags"`
	TargetResourceTypes                   []string                                   `tfschema:"target_resource_types"`
	WindowSize                            string                                     `tfschema:"window_duration"`
	CreatedWithApiVersion                 string                                     `tfschema:"created_with_api_version"`
	IsLegacyLogAnalyticsRule              bool                                       `tfschema:"is_a_legacy_log_analytics_rule"`
	IsWorkspaceAlertsStorageConfigured    bool                                       `tfschema:"is_workspace_alerts_storage_configured"`
	Identity                              []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
}

type ScheduledQueryRulesAlertV2ActionsModel struct {
	ActionGroups     []string          `tfschema:"action_groups"`
	CustomProperties map[string]string `tfschema:"custom_properties"`
	EmailSubject     string            `tfschema:"email_subject"`
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

func (r Resource) ModelObject() any {
	return &ScheduledQueryRulesAlertV2Model{}
}
