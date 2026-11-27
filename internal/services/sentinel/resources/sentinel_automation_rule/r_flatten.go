// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_automation_rule

import (
	"encoding/json"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
)

func flattenAutomationRuleConditionsToJSON(input *[]automationrules.AutomationRuleCondition) (string, error) {
	if input == nil || len(*input) == 0 {
		return "", nil
	}
	result, err := json.Marshal(input)
	return string(result), err
}

func flattenAutomationRuleActions(input []automationrules.AutomationRuleAction) (actionIncident []any, actionPlaybook []any, actionIncidentTask []any) {
	actionIncident = make([]any, 0)
	actionPlaybook = make([]any, 0)
	actionIncidentTask = make([]any, 0)

	for _, action := range input {
		switch action := action.(type) {
		case automationrules.AutomationRuleModifyPropertiesAction:
			actionIncident = append(actionIncident, flattenAutomationRuleActionIncident(action))
		case automationrules.AutomationRuleRunPlaybookAction:
			actionPlaybook = append(actionPlaybook, flattenAutomationRuleActionPlaybook(action))
		case automationrules.AutomationRuleAddIncidentTaskAction:
			actionIncidentTask = append(actionIncidentTask, flattenAutomationRuleACtionIncidentTask(action))
		}
	}

	return
}

func flattenAutomationRuleActionIncident(input automationrules.AutomationRuleModifyPropertiesAction) map[string]any {
	var (
		status      string
		clsf        string
		clsfComment string
		clsfReason  string
		labels      []any
		owner       string
		severity    string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		if cfg.Status != nil {
			status = string(*cfg.Status)
		}
		if cfg.Classification != nil {
			clsf = string(*cfg.Classification)
		}
		if cfg.ClassificationComment != nil {
			clsfComment = *cfg.ClassificationComment
		}
		if cfg.ClassificationReason != nil {
			clsfReason = string(*cfg.ClassificationReason)
		}

		if cfg.Labels != nil {
			for _, label := range *cfg.Labels {
				labels = append(labels, label.LabelName)
			}
		}

		if cfg.Owner != nil && cfg.Owner.ObjectId != nil {
			owner = *cfg.Owner.ObjectId
		}

		if cfg.Severity != nil {
			severity = string(*cfg.Severity)
		}
	}

	classification := clsf
	if clsfReason != "" {
		classification = classification + "_" + clsfReason
	}

	return map[string]any{
		"order":                  input.Order,
		"status":                 status,
		"classification":         classification,
		"classification_comment": clsfComment,
		"labels":                 labels,
		"owner_id":               owner,
		"severity":               severity,
	}
}

func flattenAutomationRuleActionPlaybook(input automationrules.AutomationRuleRunPlaybookAction) map[string]any {
	var (
		logicAppId string
		tenantId   string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		logicAppId = cfg.LogicAppResourceId

		if cfg.TenantId != nil {
			tenantId = *cfg.TenantId
		}
	}

	return map[string]any{
		"order":        input.Order,
		"logic_app_id": logicAppId,
		"tenant_id":    tenantId,
	}
}

func flattenAutomationRuleACtionIncidentTask(input automationrules.AutomationRuleAddIncidentTaskAction) map[string]any {
	var (
		title       string
		description string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		title = cfg.Title
		description = pointer.From(cfg.Description)
	}

	return map[string]any{
		"order":       input.Order,
		"title":       title,
		"description": description,
	}
}
