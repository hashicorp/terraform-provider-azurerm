package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func alertRuleID(rule securityinsight.BasicAlertRule) *string {
	if rule == nil {
		return nil
	}
	switch rule := rule.(type) {
	case securityinsight.FusionAlertRule:
		return rule.ID
	case securityinsight.MicrosoftSecurityIncidentCreationAlertRule:
		return rule.ID
	case securityinsight.ScheduledAlertRule:
		return rule.ID
	case securityinsight.MLBehaviorAnalyticsAlertRule:
		return rule.ID
	case securityinsight.NrtAlertRule:
		return rule.ID
	default:
		return nil
	}
}

func importSentinelAlertRule(expectKind securityinsight.AlertRuleKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.AlertRuleID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Sentinel.AlertRulesClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, expectKind); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}

func assertAlertRuleKind(rule securityinsight.BasicAlertRule, expectKind securityinsight.AlertRuleKind) error {
	var kind securityinsight.AlertRuleKind
	switch rule.(type) {
	case securityinsight.MLBehaviorAnalyticsAlertRule:
		kind = securityinsight.AlertRuleKindMLBehaviorAnalytics
	case securityinsight.FusionAlertRule:
		kind = securityinsight.AlertRuleKindFusion
	case securityinsight.MicrosoftSecurityIncidentCreationAlertRule:
		kind = securityinsight.AlertRuleKindMicrosoftSecurityIncidentCreation
	case securityinsight.ScheduledAlertRule:
		kind = securityinsight.AlertRuleKindScheduled
	case securityinsight.NrtAlertRule:
		kind = securityinsight.AlertRuleKindNRT
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Alert Rule has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}

func expandAlertRuleTactics(input []interface{}) *[]securityinsight.AttackTactic {
	result := make([]securityinsight.AttackTactic, 0)

	for _, e := range input {
		result = append(result, securityinsight.AttackTactic(e.(string)))
	}

	return &result
}

func flattenAlertRuleTactics(input *[]securityinsight.AttackTactic) []interface{} {
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

func expandAlertRuleIncidentConfiguration(input []interface{}, createIncidentKey string, withGroupByPrefix bool) *securityinsight.IncidentConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &securityinsight.IncidentConfiguration{
		CreateIncident:        utils.Bool(raw[createIncidentKey].(bool)),
		GroupingConfiguration: expandAlertRuleGrouping(raw["grouping"].([]interface{}), withGroupByPrefix),
	}

	return output
}

func flattenAlertRuleIncidentConfiguration(input *securityinsight.IncidentConfiguration, createIncidentKey string, withGroupByPrefix bool) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	createIncident := false
	if input.CreateIncident != nil {
		createIncident = *input.CreateIncident
	}

	return []interface{}{
		map[string]interface{}{
			createIncidentKey: createIncident,
			"grouping":        flattenAlertRuleGrouping(input.GroupingConfiguration, withGroupByPrefix),
		},
	}
}

func expandAlertRuleGrouping(input []interface{}, withGroupPrefix bool) *securityinsight.GroupingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &securityinsight.GroupingConfiguration{
		Enabled:              utils.Bool(raw["enabled"].(bool)),
		ReopenClosedIncident: utils.Bool(raw["reopen_closed_incidents"].(bool)),
		LookbackDuration:     utils.String(raw["lookback_duration"].(string)),
		MatchingMethod:       securityinsight.MatchingMethod(raw["entity_matching_method"].(string)),
	}

	key := "by_entities"
	if withGroupPrefix {
		key = "group_" + key
	}
	groupByEntitiesList := raw[key].([]interface{})
	groupByEntities := make([]securityinsight.EntityMappingType, len(groupByEntitiesList))
	for idx, t := range groupByEntitiesList {
		groupByEntities[idx] = securityinsight.EntityMappingType(t.(string))
	}
	output.GroupByEntities = &groupByEntities

	key = "by_alert_details"
	if withGroupPrefix {
		key = "group_" + key
	}
	groupByAlertDetailsList := raw[key].([]interface{})
	groupByAlertDetails := make([]securityinsight.AlertDetail, len(groupByAlertDetailsList))
	for idx, t := range groupByAlertDetailsList {
		groupByAlertDetails[idx] = securityinsight.AlertDetail(t.(string))
	}
	output.GroupByAlertDetails = &groupByAlertDetails

	key = "by_custom_details"
	if withGroupPrefix {
		key = "group_" + key
	}
	output.GroupByCustomDetails = utils.ExpandStringSlice(raw[key].([]interface{}))

	return output
}

func flattenAlertRuleGrouping(input *securityinsight.GroupingConfiguration, withGroupPrefix bool) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	enabled := false
	if input.Enabled != nil {
		enabled = *input.Enabled
	}

	lookbackDuration := ""
	if input.LookbackDuration != nil {
		lookbackDuration = *input.LookbackDuration
	}

	reopenClosedIncidents := false
	if input.ReopenClosedIncident != nil {
		reopenClosedIncidents = *input.ReopenClosedIncident
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
			"enabled":                 enabled,
			"lookback_duration":       lookbackDuration,
			"reopen_closed_incidents": reopenClosedIncidents,
			"entity_matching_method":  string(input.MatchingMethod),
			k1:                        groupByEntities,
			k2:                        groupByAlertDetails,
			k3:                        groupByCustomDetails,
		},
	}
}

func expandAlertRuleAlertDetailsOverride(input []interface{}) *securityinsight.AlertDetailsOverride {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	b := input[0].(map[string]interface{})
	output := &securityinsight.AlertDetailsOverride{}

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

func flattenAlertRuleAlertDetailsOverride(input *securityinsight.AlertDetailsOverride) []interface{} {
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

func expandAlertRuleAlertDynamicProperties(input []interface{}) *[]securityinsight.AlertPropertyMapping {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	var output []securityinsight.AlertPropertyMapping

	for _, v := range input {
		b := v.(map[string]interface{})
		output = append(output, securityinsight.AlertPropertyMapping{
			AlertProperty: securityinsight.AlertProperty(b["name"].(string)),
			Value:         utils.String(b["value"].(string)),
		})
	}

	return &output
}

func flattenAlertRuleAlertDynamicProperties(input *[]securityinsight.AlertPropertyMapping) []interface{} {
	output := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return output
	}

	for _, i := range *input {
		output = append(output, map[string]interface{}{
			"name":  string(i.AlertProperty),
			"value": i.Value,
		})
	}

	return output
}

func expandAlertRuleEntityMapping(input []interface{}) *[]securityinsight.EntityMapping {
	if len(input) == 0 {
		return nil
	}

	result := make([]securityinsight.EntityMapping, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		result = append(result, securityinsight.EntityMapping{
			EntityType:    securityinsight.EntityMappingType(b["entity_type"].(string)),
			FieldMappings: expandAlertRuleFieldMapping(b["field_mapping"].([]interface{})),
		})
	}

	return &result
}

func flattenAlertRuleEntityMapping(input *[]securityinsight.EntityMapping) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, map[string]interface{}{
			"entity_type":   string(e.EntityType),
			"field_mapping": flattenAlertRuleFieldMapping(e.FieldMappings),
		})
	}

	return output
}

func expandAlertRuleFieldMapping(input []interface{}) *[]securityinsight.FieldMapping {
	if len(input) == 0 {
		return nil
	}

	result := make([]securityinsight.FieldMapping, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		result = append(result, securityinsight.FieldMapping{
			Identifier: utils.String(b["identifier"].(string)),
			ColumnName: utils.String(b["column_name"].(string)),
		})
	}

	return &result
}

func flattenAlertRuleFieldMapping(input *[]securityinsight.FieldMapping) []interface{} {
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

func expandAlertRuleSentinelEntityMapping(input []interface{}) *[]securityinsight.SentinelEntityMapping {
	if len(input) == 0 {
		return nil
	}

	result := make([]securityinsight.SentinelEntityMapping, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		result = append(result, securityinsight.SentinelEntityMapping{
			ColumnName: utils.String(b["column_name"].(string)),
		})
	}

	return &result
}

func flattenAlertRuleSentinelEntityMapping(input *[]securityinsight.SentinelEntityMapping) []interface{} {
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
