// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func importSentinelAlertRule(expectKind alertrules.AlertRuleKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := alertrules.ParseAlertRuleID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Sentinel.AlertRulesClient
		resp, err := client.AlertRulesGet(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		if err = assertAlertRuleKind(resp.Model, expectKind); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}

func importSentinelAlertRuleForTypedSdk(expectKind alertrules.AlertRuleKind) sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := alertrules.ParseAlertRuleID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.Sentinel.AlertRulesClient
		resp, err := client.AlertRulesGet(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Model, expectKind); err != nil {
			return err
		}
		return nil
	}
}

func assertAlertRuleKind(rule *alertrules.AlertRule, expectKind alertrules.AlertRuleKind) error {
	if rule == nil {
		return fmt.Errorf("model was nil")
	}

	rulePtr := *rule
	var kind alertrules.AlertRuleKind
	switch rulePtr.(type) {
	case alertrules.MLBehaviorAnalyticsAlertRule:
		kind = alertrules.AlertRuleKindMLBehaviorAnalytics
	case alertrules.FusionAlertRule:
		kind = alertrules.AlertRuleKindFusion
	case alertrules.MicrosoftSecurityIncidentCreationAlertRule:
		kind = alertrules.AlertRuleKindMicrosoftSecurityIncidentCreation
	case alertrules.ScheduledAlertRule:
		kind = alertrules.AlertRuleKindScheduled
	case alertrules.NrtAlertRule:
		kind = alertrules.AlertRuleKindNRT
	case alertrules.ThreatIntelligenceAlertRule:
		kind = alertrules.AlertRuleKindThreatIntelligence
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Alert Rule has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}

func expandAlertRuleTactics(input []interface{}) *[]alertrules.AttackTactic {
	result := make([]alertrules.AttackTactic, 0)

	for _, e := range input {
		result = append(result, alertrules.AttackTactic(e.(string)))
	}

	return &result
}

func flattenAlertRuleTactics(input *[]alertrules.AttackTactic) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, string(e))
	}

	return output
}

func expandAlertRuleTechnicals(input []interface{}) *[]string {
	result := make([]string, 0)

	for _, e := range input {
		result = append(result, e.(string))
	}

	return &result
}

func expandAlertRuleIncidentConfiguration(input []interface{}, createIncidentKey string, withGroupByPrefix bool) *alertrules.IncidentConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &alertrules.IncidentConfiguration{
		CreateIncident:        raw[createIncidentKey].(bool),
		GroupingConfiguration: expandAlertRuleGrouping(raw["grouping"].([]interface{}), withGroupByPrefix),
	}

	return output
}

func flattenAlertRuleIncidentConfiguration(input *alertrules.IncidentConfiguration, createIncidentKey string, withGroupByPrefix bool) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			createIncidentKey: input.CreateIncident,
			"grouping":        flattenAlertRuleGrouping(input.GroupingConfiguration, withGroupByPrefix),
		},
	}
}

func expandAlertRuleGrouping(input []interface{}, withGroupPrefix bool) *alertrules.GroupingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &alertrules.GroupingConfiguration{
		Enabled:              raw["enabled"].(bool),
		ReopenClosedIncident: raw["reopen_closed_incidents"].(bool),
		LookbackDuration:     raw["lookback_duration"].(string),
		MatchingMethod:       alertrules.MatchingMethod(raw["entity_matching_method"].(string)),
	}

	key := "by_entities"
	if withGroupPrefix {
		key = "group_" + key
	}
	groupByEntitiesList := raw[key].([]interface{})
	groupByEntities := make([]alertrules.EntityMappingType, len(groupByEntitiesList))
	for idx, t := range groupByEntitiesList {
		groupByEntities[idx] = alertrules.EntityMappingType(t.(string))
	}
	output.GroupByEntities = &groupByEntities

	key = "by_alert_details"
	if withGroupPrefix {
		key = "group_" + key
	}
	groupByAlertDetailsList := raw[key].([]interface{})
	groupByAlertDetails := make([]alertrules.AlertDetail, len(groupByAlertDetailsList))
	for idx, t := range groupByAlertDetailsList {
		groupByAlertDetails[idx] = alertrules.AlertDetail(t.(string))
	}
	output.GroupByAlertDetails = &groupByAlertDetails

	key = "by_custom_details"
	if withGroupPrefix {
		key = "group_" + key
	}
	output.GroupByCustomDetails = utils.ExpandStringSlice(raw[key].([]interface{}))

	return output
}

func flattenAlertRuleGrouping(input *alertrules.GroupingConfiguration, withGroupPrefix bool) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var groupByEntities []interface{}
	if input.GroupByEntities != nil {
		for _, entity := range *input.GroupByEntities {
			groupByEntities = append(groupByEntities, string(entity))
		}
	}

	var groupByAlertDetails []interface{}
	if input.GroupByAlertDetails != nil {
		for _, detail := range *input.GroupByAlertDetails {
			groupByAlertDetails = append(groupByAlertDetails, string(detail))
		}
	}

	var groupByCustomDetails []interface{}
	if input.GroupByCustomDetails != nil {
		for _, detail := range *input.GroupByCustomDetails {
			groupByCustomDetails = append(groupByCustomDetails, detail)
		}
	}

	var (
		k1 = "by_entities"
		k2 = "by_alert_details"
		k3 = "by_custom_details"
	)

	if withGroupPrefix {
		k1 = "group_" + k1
		k2 = "group_" + k2
		k3 = "group_" + k3
	}
	return []interface{}{
		map[string]interface{}{
			"enabled":                 input.Enabled,
			"lookback_duration":       input.LookbackDuration,
			"reopen_closed_incidents": input.ReopenClosedIncident,
			"entity_matching_method":  string(input.MatchingMethod),
			k1:                        groupByEntities,
			k2:                        groupByAlertDetails,
			k3:                        groupByCustomDetails,
		},
	}
}

func expandAlertRuleAlertDetailsOverride(input []interface{}) *alertrules.AlertDetailsOverride {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	b := input[0].(map[string]interface{})
	output := &alertrules.AlertDetailsOverride{}

	if v := b["description_format"]; v != "" {
		output.AlertDescriptionFormat = utils.String(v.(string))
	}
	if v := b["display_name_format"]; v != "" {
		output.AlertDisplayNameFormat = utils.String(v.(string))
	}
	if v := b["severity_column_name"]; v != "" {
		output.AlertSeverityColumnName = utils.String(v.(string))
	}
	if v := b["tactics_column_name"]; v != "" {
		output.AlertTacticsColumnName = utils.String(v.(string))
	}
	if v := b["dynamic_property"]; v != nil && len(v.([]interface{})) > 0 {
		output.AlertDynamicProperties = expandAlertRuleAlertDynamicProperties(v.([]interface{}))
	}

	return output
}

func flattenAlertRuleAlertDetailsOverride(input *alertrules.AlertDetailsOverride) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var descriptionFormat string
	if input.AlertDescriptionFormat != nil {
		descriptionFormat = *input.AlertDescriptionFormat
	}

	var displayNameFormat string
	if input.AlertDisplayNameFormat != nil {
		displayNameFormat = *input.AlertDisplayNameFormat
	}

	var severityColumnName string
	if input.AlertSeverityColumnName != nil {
		severityColumnName = *input.AlertSeverityColumnName
	}

	var tacticsColumnName string
	if input.AlertTacticsColumnName != nil {
		tacticsColumnName = *input.AlertTacticsColumnName
	}

	var dynamicProperties []interface{}
	if input.AlertDynamicProperties != nil {
		dynamicProperties = flattenAlertRuleAlertDynamicProperties(input.AlertDynamicProperties)
	}

	return []interface{}{
		map[string]interface{}{
			"description_format":   descriptionFormat,
			"display_name_format":  displayNameFormat,
			"severity_column_name": severityColumnName,
			"tactics_column_name":  tacticsColumnName,
			"dynamic_property":     dynamicProperties,
		},
	}
}

func expandAlertRuleAlertDynamicProperties(input []interface{}) *[]alertrules.AlertPropertyMapping {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	var output []alertrules.AlertPropertyMapping

	for _, v := range input {
		b := v.(map[string]interface{})
		property := alertrules.AlertProperty(b["name"].(string))
		output = append(output, alertrules.AlertPropertyMapping{
			AlertProperty: &property,
			Value:         utils.String(b["value"].(string)),
		})
	}

	return &output
}

func flattenAlertRuleAlertDynamicProperties(input *[]alertrules.AlertPropertyMapping) []interface{} {
	output := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return output
	}

	for _, i := range *input {
		name := ""
		if i.AlertProperty != nil {
			name = string(*i.AlertProperty)
		}
		output = append(output, map[string]interface{}{
			"name":  name,
			"value": i.Value,
		})
	}

	return output
}

func expandAlertRuleEntityMapping(input []interface{}) *[]alertrules.EntityMapping {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.EntityMapping, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		mappingType := alertrules.EntityMappingType(b["entity_type"].(string))
		result = append(result, alertrules.EntityMapping{
			EntityType:    &mappingType,
			FieldMappings: expandAlertRuleFieldMapping(b["field_mapping"].([]interface{})),
		})
	}

	return &result
}

func flattenAlertRuleEntityMapping(input *[]alertrules.EntityMapping) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		entityType := ""
		if e.EntityType != nil {
			entityType = string(*e.EntityType)
		}
		output = append(output, map[string]interface{}{
			"entity_type":   entityType,
			"field_mapping": flattenAlertRuleFieldMapping(e.FieldMappings),
		})
	}

	return output
}

func expandAlertRuleFieldMapping(input []interface{}) *[]alertrules.FieldMapping {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.FieldMapping, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		result = append(result, alertrules.FieldMapping{
			Identifier: utils.String(b["identifier"].(string)),
			ColumnName: utils.String(b["column_name"].(string)),
		})
	}

	return &result
}

func flattenAlertRuleFieldMapping(input *[]alertrules.FieldMapping) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var identifier string
		if e.Identifier != nil {
			identifier = *e.Identifier
		}

		var columnName string
		if e.ColumnName != nil {
			columnName = *e.ColumnName
		}

		output = append(output, map[string]interface{}{
			"identifier":  identifier,
			"column_name": columnName,
		})
	}

	return output
}

func expandAlertRuleSentinelEntityMapping(input []interface{}) *[]alertrules.SentinelEntityMapping {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.SentinelEntityMapping, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		result = append(result, alertrules.SentinelEntityMapping{
			ColumnName: utils.String(b["column_name"].(string)),
		})
	}

	return &result
}

func flattenAlertRuleSentinelEntityMapping(input *[]alertrules.SentinelEntityMapping) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var columnName string
		if e.ColumnName != nil {
			columnName = *e.ColumnName
		}

		output = append(output, map[string]interface{}{
			"column_name": columnName,
		})
	}

	return output
}
