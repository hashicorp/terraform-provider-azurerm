// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_automation

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
)

func expandSecurityCenterAutomationSources(sourcesRaw []any) (*[]automations.AutomationSource, error) {
	if len(sourcesRaw) == 0 {
		return &[]automations.AutomationSource{}, nil
	}

	// Output is an array of AutomationSource
	output := make([]automations.AutomationSource, 0)

	// Top level loop over sources array
	for _, sourceRaw := range sourcesRaw {
		sourceMap, ok := sourceRaw.(map[string]any)
		if !ok {
			return nil, errors.New("'Security Center Automation' unable to decode sources")
		}

		// Build and parse array of RuleSets
		ruleSets := make([]automations.AutomationRuleSet, 0)
		ruleSetsRaw := sourceMap["rule_set"].([]any)
		for _, ruleSetRaw := range ruleSetsRaw {
			ruleSetMap := ruleSetRaw.(map[string]any)
			rulesRaw := ruleSetMap["rule"].([]any)

			// Build and parse array of Rules in each RuleSet
			rules := make([]automations.AutomationTriggeringRule, 0)
			for _, ruleRaw := range rulesRaw {
				// Parse the rule fields
				ruleMap := ruleRaw.(map[string]any)

				// Create AutomationTriggeringRule struct and push into array
				rule := automations.AutomationTriggeringRule{
					PropertyJPath: pointer.To(ruleMap["property_path"].(string)),
					PropertyType:  pointer.ToEnum[automations.PropertyType](ruleMap["property_type"].(string)),
					ExpectedValue: pointer.To(ruleMap["expected_value"].(string)),
					Operator:      pointer.ToEnum[automations.Operator](ruleMap["operator"].(string)),
				}
				rules = append(rules, rule)
			}

			// Create AutomationRuleSet struct and push into array
			ruleSet := automations.AutomationRuleSet{
				Rules: &rules,
			}
			ruleSets = append(ruleSets, ruleSet)
		}

		// Finally create AutomationSource struct holding our list of RuleSets
		source := automations.AutomationSource{
			EventSource: pointer.ToEnum[automations.EventSource](sourceMap["event_source"].(string)),
			RuleSets:    &ruleSets,
		}

		// Finally (no really this time), push AutomationSource into output
		output = append(output, source)
	}

	return &output, nil
}

func expandSecurityCenterAutomationScopes(scopePathsRaw []any) *[]automations.AutomationScope {
	scopes := make([]automations.AutomationScope, 0)

	for _, scopePathRaw := range scopePathsRaw {
		if path, ok := scopePathRaw.(string); ok {
			desc := fmt.Sprintf("scope for %s", path)
			scope := automations.AutomationScope{
				ScopePath:   &path,
				Description: &desc,
			}
			scopes = append(scopes, scope)
		}
	}

	return &scopes
}

func expandSecurityCenterAutomationActions(actionsRaw []any) (*[]automations.AutomationAction, error) {
	if len(actionsRaw) == 0 {
		return &[]automations.AutomationAction{}, nil
	}

	output := make([]automations.AutomationAction, 0)

	for _, actionRaw := range actionsRaw {
		actionMap := actionRaw.(map[string]any)

		var autoAction automations.AutomationAction
		var resourceID string
		var ok bool

		// No checking, as fields are enforced by resource schema
		resourceID = actionMap["resource_id"].(string)
		actionType := automations.ActionType(actionMap["type"].(string))

		switch actionType {
		// Handle LogicApp action type
		case automations.ActionTypeLogicApp:
			var triggerURL string
			if triggerURL, ok = actionMap["trigger_url"].(string); !ok || triggerURL == "" {
				return nil, errors.New("'Security Center Automation' trigger_url is required for LogicApp action")
			}
			autoAction = automations.AutomationActionLogicApp{
				LogicAppResourceId: &resourceID,
				Uri:                &triggerURL,
			}

			// Handle LogAnalytics action type
		case automations.ActionTypeWorkspace:
			autoAction = automations.AutomationActionWorkspace{
				WorkspaceResourceId: &resourceID,
			}

			// Handle EventHub action type
		case automations.ActionTypeEventHub:
			var connString string
			if connString, ok = actionMap["connection_string"].(string); !ok || connString == "" {
				return nil, errors.New("'Security Center Automation' connection_string is required for EventHub action")
			}
			autoAction = automations.AutomationActionEventHub{
				EventHubResourceId: &resourceID,
				ConnectionString:   &connString,
			}
		default:
			return nil, fmt.Errorf("'Security Center automation' expected action type to be one of: %s, %s or %s", automations.ActionTypeEventHub, automations.ActionTypeWorkspace, automations.ActionTypeLogicApp)
		}
		output = append(output, autoAction)
	}

	return &output, nil
}
