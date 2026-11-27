// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_automation

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenSecurityCenterAutomationSources(sources *[]automations.AutomationSource) ([]map[string]any, error) {
	if sources == nil {
		return make([]map[string]any, 0), nil
	}

	resultSlice := make([]map[string]any, 0)
	for _, source := range *sources {
		ruleSetSlice := make([]any, 0)

		// RuleSets is an optional field need check for nil
		if source.RuleSets != nil {
			for _, ruleSet := range *source.RuleSets {
				ruleSlice := make([]map[string]string, 0)

				for _, rule := range *ruleSet.Rules {
					if rule.PropertyJPath == nil {
						return nil, errors.New("'Security Center Automation' API returned a rule with an empty PropertyJPath")
					}
					if rule.ExpectedValue == nil {
						return nil, errors.New("'Security Center Automation' API returned a rule with empty ExpectedValue")
					}
					ruleMap := map[string]string{
						"property_path":  *rule.PropertyJPath,
						"expected_value": *rule.ExpectedValue,
						"operator":       string(*rule.Operator),
						"property_type":  string(*rule.PropertyType),
					}
					ruleSlice = append(ruleSlice, ruleMap)
				}

				ruleSetMap := map[string]any{
					"rule": ruleSlice,
				}
				ruleSetSlice = append(ruleSetSlice, ruleSetMap)
			}
		}

		sourceMap := map[string]any{
			"event_source": source.EventSource,
			"rule_set":     ruleSetSlice,
		}
		resultSlice = append(resultSlice, sourceMap)
	}

	return resultSlice, nil
}

func flattenSecurityCenterAutomationScopes(scopes *[]automations.AutomationScope) ([]string, error) {
	if scopes == nil {
		return []string{}, nil
	}

	resultSlice := make([]string, 0)
	for _, scope := range *scopes {
		if scope.ScopePath == nil {
			return nil, errors.New("'Security Center Automation' API returned a scope with an empty ScopePath")
		}

		resultSlice = append(resultSlice, *scope.ScopePath)
	}

	return resultSlice, nil
}

func flattenSecurityCenterAutomationActions(actions *[]automations.AutomationAction, d *pluginsdk.ResourceData) ([]map[string]string, error) {
	if actions == nil {
		return []map[string]string{}, nil
	}

	resultSlice := make([]map[string]string, 0)

	for i, action := range *actions {
		// Use type assertion to discover the underlying action
		// Trying to use action.(automations.AutomationAction).ActionType results in a panic
		actionLogicApp, isLogicApp := action.(automations.AutomationActionLogicApp)
		if isLogicApp {
			if actionLogicApp.LogicAppResourceId == nil {
				return nil, errors.New("'Security Center Automation' API returned an action with empty logicAppResourceId")
			}
			actionMap := map[string]string{
				"resource_id": *actionLogicApp.LogicAppResourceId,
				"type":        string(automations.ActionTypeLogicApp),
				"trigger_url": "",
			}

			// Need to merge in trigger_url as it's not returned by API Get operation
			// Otherwise don't have consistent state
			if triggerURL, ok := d.GetOk(fmt.Sprintf("action.%d.trigger_url", i)); ok {
				actionMap["trigger_url"] = triggerURL.(string)
			}

			resultSlice = append(resultSlice, actionMap)
		}

		actionEventHub, isEventHub := action.(automations.AutomationActionEventHub)
		if isEventHub {
			if actionEventHub.EventHubResourceId == nil {
				return nil, errors.New("'Security Center Automation' API returned an action with empty eventHubResourceId")
			}
			actionMap := map[string]string{
				"resource_id":       *actionEventHub.EventHubResourceId,
				"type":              string(automations.ActionTypeEventHub),
				"connection_string": "",
			}

			// Need to merge in connection_string as it's not returned by API Get operation
			// Otherwise don't have consistent state
			if triggerURL, ok := d.GetOk(fmt.Sprintf("action.%d.connection_string", i)); ok {
				actionMap["connection_string"] = triggerURL.(string)
			}

			resultSlice = append(resultSlice, actionMap)
		}

		actionLogAnalytics, isLogAnalytics := action.(automations.AutomationActionWorkspace)
		if isLogAnalytics {
			if actionLogAnalytics.WorkspaceResourceId == nil {
				return nil, errors.New("security center automation, API returned an action with empty workspaceResourceId")
			}
			actionMap := map[string]string{
				"resource_id": *actionLogAnalytics.WorkspaceResourceId,
				"type":        string(automations.ActionTypeWorkspace),
			}

			resultSlice = append(resultSlice, actionMap)
		}
	}

	return resultSlice, nil
}
