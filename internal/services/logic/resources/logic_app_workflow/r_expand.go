// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_workflow

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandLogicAppWorkflowParameters(input map[string]any, paramDefs map[string]any) (*map[string]workflows.WorkflowParameter, error) {
	output := make(map[string]workflows.WorkflowParameter)

	for k, v := range input {
		defRaw, ok := paramDefs[k]
		if !ok {
			return nil, fmt.Errorf("no parameter definition for %s", k)
		}
		def := defRaw.(map[string]any)
		t := workflows.ParameterType(def["type"].(string))

		v := v.(string)

		var value any
		switch t {
		case workflows.ParameterTypeBool:
			var uv bool
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to bool: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeFloat:
			var uv float64
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to float64: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeInt:
			var uv int
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to int: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeArray:
			var uv []any
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to []interface{}: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeObject,
			workflows.ParameterTypeSecureObject:
			var uv map[string]any
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to map[string]interface{}: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeString,
			workflows.ParameterTypeSecureString:
			value = v
		}

		if k == "$connections" {
			output[k] = workflows.WorkflowParameter{
				Value: &value,
			}
		} else {
			output[k] = workflows.WorkflowParameter{
				Type:  &t,
				Value: &value,
			}
		}
	}

	return &output, nil
}

func expandLogicAppWorkflowWorkflowParameters(input map[string]any) (map[string]any, error) {
	if len(input) == 0 {
		return nil, nil
	}

	output := make(map[string]any)
	for k, v := range input {
		obj, err := pluginsdk.ExpandJsonFromString(v.(string))
		if err != nil {
			return nil, err
		}
		output[k] = obj
	}
	return output, nil
}

func expandLogicAppWorkflowAccessControl(input []any) *workflows.FlowAccessControlConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)

	result := workflows.FlowAccessControlConfiguration{}

	if contents := v["content"].([]any); len(contents) != 0 {
		result.Contents = expandLogicAppWorkflowAccessControlConfigurationPolicy(contents)
	}

	if actions := v["action"].([]any); len(actions) != 0 {
		result.Actions = expandLogicAppWorkflowAccessControlConfigurationPolicy(actions)
	}

	if triggers := v["trigger"].([]any); len(triggers) != 0 {
		result.Triggers = expandLogicAppWorkflowAccessControlTriggerConfigurationPolicy(triggers)
	}

	if workflowManagement := v["workflow_management"].([]any); len(workflowManagement) != 0 {
		result.WorkflowManagement = expandLogicAppWorkflowAccessControlConfigurationPolicy(workflowManagement)
	}

	return &result
}

func expandLogicAppWorkflowAccessControlConfigurationPolicy(input []any) *workflows.FlowAccessControlConfigurationPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)

	return &workflows.FlowAccessControlConfigurationPolicy{
		AllowedCallerIPAddresses: expandLogicAppWorkflowIPAddressRanges(v["allowed_caller_ip_address_range"].(*pluginsdk.Set).List()),
	}
}

func expandLogicAppWorkflowAccessControlTriggerConfigurationPolicy(input []any) *workflows.FlowAccessControlConfigurationPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)

	result := workflows.FlowAccessControlConfigurationPolicy{
		AllowedCallerIPAddresses: expandLogicAppWorkflowIPAddressRanges(v["allowed_caller_ip_address_range"].(*pluginsdk.Set).List()),
	}

	if openAuthenticationPolicy, ok := v["open_authentication_policy"]; ok {
		openAuthenticationPolicies := openAuthenticationPolicy.(*pluginsdk.Set).List()
		if len(openAuthenticationPolicies) != 0 {
			result.OpenAuthenticationPolicies = &workflows.OpenAuthenticationAccessPolicies{
				Policies: expandLogicAppWorkflowOpenAuthenticationPolicy(openAuthenticationPolicies),
			}
		}
	}

	return &result
}

func expandLogicAppWorkflowIPAddressRanges(input []any) *[]workflows.IPAddressRange {
	results := make([]workflows.IPAddressRange, 0)

	for _, item := range input {
		results = append(results, workflows.IPAddressRange{
			AddressRange: pointer.To(item.(string)),
		})
	}

	return &results
}

func expandLogicAppWorkflowOpenAuthenticationPolicy(input []any) *map[string]workflows.OpenAuthenticationAccessPolicy {
	if len(input) == 0 {
		return nil
	}
	results := make(map[string]workflows.OpenAuthenticationAccessPolicy)

	for _, item := range input {
		if item == nil {
			continue
		}
		v := item.(map[string]any)
		policyName := v["name"].(string)

		results[policyName] = workflows.OpenAuthenticationAccessPolicy{
			Type:   pointer.To(workflows.OpenAuthenticationProviderTypeAAD),
			Claims: expandLogicAppWorkflowOpenAuthenticationPolicyClaim(v["claim"].(*pluginsdk.Set).List()),
		}
	}

	return &results
}

func expandLogicAppWorkflowOpenAuthenticationPolicyClaim(input []any) *[]workflows.OpenAuthenticationPolicyClaim {
	results := make([]workflows.OpenAuthenticationPolicyClaim, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, workflows.OpenAuthenticationPolicyClaim{
			Name:  pointer.To(v["name"].(string)),
			Value: pointer.To(v["value"].(string)),
		})
	}
	return &results
}
