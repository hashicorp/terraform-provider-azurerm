// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_automation_rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAutomationRuleConditionsFromJSON(input string) (*[]automationrules.AutomationRuleCondition, error) {
	if input == "" {
		return nil, nil
	}
	triggerLogic := &automationrules.AutomationRuleTriggeringLogic{}
	if err := triggerLogic.UnmarshalJSON([]byte(fmt.Sprintf(`{ "conditions": %s }`, input))); err != nil {
		return nil, err
	}
	return triggerLogic.Conditions, nil
}

func expandAutomationRuleActions(d *pluginsdk.ResourceData, defaultTenantId string) ([]automationrules.AutomationRuleAction, error) {
	actionIncident, err := expandAutomationRuleActionIncident(d.Get("action_incident").([]any))
	if err != nil {
		return nil, err
	}
	actionPlaybook := expandAutomationRuleActionPlaybook(d.Get("action_playbook").([]any), defaultTenantId)

	actionIncidentTask := expandAutomationRuleActionIncidentTask(d.Get("action_incident_task").([]any))

	if len(actionIncident)+len(actionPlaybook)+len(actionIncidentTask) == 0 {
		return nil, nil
	}

	out := make([]automationrules.AutomationRuleAction, 0, len(actionIncident)+len(actionPlaybook))
	out = append(out, actionIncident...)
	out = append(out, actionPlaybook...)
	out = append(out, actionIncidentTask...)
	return out, nil
}

func expandAutomationRuleActionIncident(input []any) ([]automationrules.AutomationRuleAction, error) {
	if len(input) == 0 {
		return nil, nil
	}

	out := make([]automationrules.AutomationRuleAction, 0, len(input))
	for _, b := range input {
		b := b.(map[string]any)

		status := automationrules.IncidentStatus(b["status"].(string))
		l := strings.Split(b["classification"].(string), "_")
		classification, clr := l[0], ""
		if len(l) == 2 {
			clr = l[1]
		}
		classificationComment := b["classification_comment"].(string)

		// sanity check on classification
		if status == automationrules.IncidentStatusClosed && classification == "" {
			return nil, fmt.Errorf("`classification` is required when `status` is set to `Closed`")
		}
		if status != automationrules.IncidentStatusClosed {
			if classification != "" {
				return nil, fmt.Errorf("`classification` can't be set when `status` is not set to `Closed`")
			}
			if classificationComment != "" {
				return nil, fmt.Errorf("`classification_comment` can't be set when `status` is not set to `Closed`")
			}
		}

		var labelsPtr *[]automationrules.IncidentLabel
		if labelStrsPtr := pluginsdk.ExpandStringSlice(b["labels"].([]any)); labelStrsPtr != nil && len(*labelStrsPtr) > 0 {
			labels := make([]automationrules.IncidentLabel, 0, len(*labelStrsPtr))
			for _, label := range *labelStrsPtr {
				labels = append(labels, automationrules.IncidentLabel{
					LabelName: label,
				})
			}
			labelsPtr = &labels
		}

		var ownerPtr *automationrules.IncidentOwnerInfo
		if ownerIdStr := b["owner_id"].(string); ownerIdStr != "" {
			ownerPtr = &automationrules.IncidentOwnerInfo{
				ObjectId: pointer.To(ownerIdStr),
			}
		}

		severity := b["severity"].(string)

		// sanity check on the whole incident action
		if severity == "" && ownerPtr == nil && labelsPtr == nil && status == "" {
			return nil, fmt.Errorf("at least one of `severity`, `owner_id`, `labels` or `status` should be specified")
		}

		out = append(out, automationrules.AutomationRuleModifyPropertiesAction{
			Order: int64(b["order"].(int)),
			ActionConfiguration: &automationrules.IncidentPropertiesAction{
				Status:                &status,
				Classification:        pointer.ToEnum[automationrules.IncidentClassification](classification),
				ClassificationComment: &classificationComment,
				ClassificationReason:  pointer.ToEnum[automationrules.IncidentClassificationReason](clr),
				Labels:                labelsPtr,
				Owner:                 ownerPtr,
				Severity:              pointer.ToEnum[automationrules.IncidentSeverity](severity),
			},
		})
	}

	return out, nil
}

func expandAutomationRuleActionPlaybook(input []any, defaultTenantId string) []automationrules.AutomationRuleAction {
	out := make([]automationrules.AutomationRuleAction, 0, len(input))
	for _, b := range input {
		b := b.(map[string]any)

		tid := defaultTenantId
		if t := b["tenant_id"].(string); t != "" {
			tid = t
		}

		out = append(out, automationrules.AutomationRuleRunPlaybookAction{
			Order: int64(b["order"].(int)),
			ActionConfiguration: &automationrules.PlaybookActionProperties{
				LogicAppResourceId: b["logic_app_id"].(string),
				TenantId:           &tid,
			},
		})
	}
	return out
}

func expandAutomationRuleActionIncidentTask(input []any) []automationrules.AutomationRuleAction {
	output := make([]automationrules.AutomationRuleAction, 0, len(input))

	for _, task := range input {
		task := task.(map[string]any)
		output = append(output, automationrules.AutomationRuleAddIncidentTaskAction{
			Order: int64(task["order"].(int)),
			ActionConfiguration: &automationrules.AddIncidentTaskActionProperties{
				Title:       task["title"].(string),
				Description: pointer.To(task["description"].(string)),
			},
		})
	}

	return output
}
