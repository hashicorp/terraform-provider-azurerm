// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_built_in

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
)

type AlertRuleAnomalyBuiltInModel struct {
	Name                         string                                          `tfschema:"name"`
	DisplayName                  string                                          `tfschema:"display_name"`
	WorkspaceId                  string                                          `tfschema:"log_analytics_workspace_id"`
	Enabled                      bool                                            `tfschema:"enabled"`
	Mode                         string                                          `tfschema:"mode"`
	AnomalyVersion               string                                          `tfschema:"anomaly_version"`
	AnomalySettingsVersion       int64                                           `tfschema:"anomaly_settings_version"`
	Description                  string                                          `tfschema:"description"`
	Frequency                    string                                          `tfschema:"frequency"`
	RequiredDataConnectors       []helpers.AnomalyRuleRequiredDataConnectorModel `tfschema:"required_data_connector"`
	SettingsDefinitionId         string                                          `tfschema:"settings_definition_id"`
	Tactics                      []string                                        `tfschema:"tactics"`
	Techniques                   []string                                        `tfschema:"techniques"`
	ThresholdObservation         []helpers.AnomalyRuleThresholdModel             `tfschema:"threshold_observation"`
	MultiSelectObservation       []helpers.AnomalyRuleMultiSelectModel           `tfschema:"multi_select_observation"`
	SingleSelectObservation      []helpers.AnomalyRuleSingleSelectModel          `tfschema:"single_select_observation"`
	PrioritizeExcludeObservation []helpers.AnomalyRulePriorityModel              `tfschema:"prioritized_exclude_observation"`
}

func (r Resource) ModelObject() any {
	return &AlertRuleAnomalyBuiltInModel{}
}
