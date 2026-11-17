// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_workflow

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenLogicAppWorkflowParameters(d *pluginsdk.ResourceData, input *map[string]workflows.WorkflowParameter, paramDefs map[string]any) (map[string]any, error) {
	output := make(map[string]any)
	if input == nil {
		return output, nil
	}

	// Read the "parameters" from state, which is used to fill in the "sensitive" properties.
	paramInState := make(map[string]any)
	paramsRaw := d.Get("parameters")
	if params, ok := paramsRaw.(map[string]any); ok {
		paramInState = params
	}

	for k, v := range *input {
		defRaw, ok := paramDefs[k]
		if !ok {
			// This should never happen.
			log.Printf("[WARN] The parameter %s is not defined in the Logic App Workflow", k)
			continue
		}

		def := defRaw.(map[string]any)
		t := workflows.ParameterType(def["type"].(string))

		var value string
		switch t {
		case workflows.ParameterTypeBool:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got nil", k)
			}
			valueRaw := *v.Value
			tv, ok := valueRaw.(bool)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got %T", k, v.Value)
			}
			value = "true"
			if !tv {
				value = "false"
			}
		case workflows.ParameterTypeFloat:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be float64, but got nil", k)
			}
			valueRaw := *v.Value
			// Note that the json unmarshalled response doesn't differ between float and int, as json has only type number.
			tv, ok := valueRaw.(float64)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be float64, but got %T", k, v.Value)
			}
			value = strconv.FormatFloat(tv, 'f', -1, 64)
		case workflows.ParameterTypeInt:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be float64, but got nil", k)
			}
			valueRaw := *v.Value
			// Note that the json unmarshalled response doesn't differ between float and int, as json has only type number.
			tv, ok := valueRaw.(float64)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be float64, but got %T", k, v.Value)
			}
			value = strconv.Itoa(int(tv))

		case workflows.ParameterTypeArray:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be []interface{}, but got nil", k)
			}
			valueRaw := *v.Value
			tv, ok := valueRaw.([]any)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be []interface{}, but got %T", k, v.Value)
			}
			obj, err := json.Marshal(tv)
			if err != nil {
				return nil, fmt.Errorf("converting %+v from json: %v", tv, err)
			}
			value = string(obj)

		case workflows.ParameterTypeObject:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be map[string]interface{}, but got nil", k)
			}
			valueRaw := *v.Value
			tv, ok := valueRaw.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be map[string]interface{}, but got %T", k, v.Value)
			}
			obj, err := json.Marshal(tv)
			if err != nil {
				return nil, fmt.Errorf("converting %+v from json: %v", tv, err)
			}
			value = string(obj)

		case workflows.ParameterTypeString:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be string, but got nil", k)
			}
			valueRaw := *v.Value
			tv, ok := valueRaw.(string)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be string, but got %T", k, v.Value)
			}
			value = tv

		case workflows.ParameterTypeSecureString,
			workflows.ParameterTypeSecureObject:
			// This is not returned from API, we will try to read them from the state instead.
			if v, ok := paramInState[k]; ok {
				value = v.(string) // The value in state here is guaranteed to be a string, so directly cast the type.
			}
		}

		output[k] = value
	}

	return output, nil
}

func flattenLogicAppWorkflowWorkflowParameters(input map[string]any) (map[string]any, error) {
	if input == nil {
		return map[string]any{}, nil
	}
	output := make(map[string]any)
	for k, v := range input {
		objstr, err := pluginsdk.FlattenJsonToString(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		output[k] = objstr
	}
	return output, nil
}

func flattenIPAddresses(input *[]workflows.IPAddress) []any {
	if input == nil {
		return []any{}
	}

	addresses := make([]any, 0, len(*input))
	for _, addr := range *input {
		addresses = append(addresses, *addr.Address)
	}
	return addresses
}

func flattenLogicAppWorkflowFlowAccessControl(input *workflows.FlowAccessControlConfiguration) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"action":              flattenLogicAppWorkflowAccessControlConfigurationPolicy(input.Actions),
			"content":             flattenLogicAppWorkflowAccessControlConfigurationPolicy(input.Contents),
			"trigger":             flattenLogicAppWorkflowAccessControlTriggerConfigurationPolicy(input.Triggers),
			"workflow_management": flattenLogicAppWorkflowAccessControlConfigurationPolicy(input.WorkflowManagement),
		},
	}
}

func flattenLogicAppWorkflowAccessControlConfigurationPolicy(input *workflows.FlowAccessControlConfigurationPolicy) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"allowed_caller_ip_address_range": flattenLogicAppWorkflowIPAddressRanges(input.AllowedCallerIPAddresses),
		},
	}
}

func flattenLogicAppWorkflowAccessControlTriggerConfigurationPolicy(input *workflows.FlowAccessControlConfigurationPolicy) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"allowed_caller_ip_address_range": flattenLogicAppWorkflowIPAddressRanges(input.AllowedCallerIPAddresses),
			"open_authentication_policy":      flattenLogicAppWorkflowOpenAuthenticationPolicy(input.OpenAuthenticationPolicies),
		},
	}
}

func flattenLogicAppWorkflowIPAddressRanges(input *[]workflows.IPAddressRange) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, pointer.From(item.AddressRange))
	}

	return results
}

func flattenLogicAppWorkflowOpenAuthenticationPolicy(input *workflows.OpenAuthenticationAccessPolicies) []any {
	results := make([]any, 0)
	if input == nil || input.Policies == nil {
		return results
	}

	for k, v := range *input.Policies {
		results = append(results, map[string]any{
			"name":  k,
			"claim": flattenLogicAppWorkflowOpenAuthenticationPolicyClaim(v.Claims),
		})
	}

	return results
}

func flattenLogicAppWorkflowOpenAuthenticationPolicyClaim(input *[]workflows.OpenAuthenticationPolicyClaim) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"name":  pointer.From(item.Name),
			"value": pointer.From(item.Value),
		})
	}

	return results
}
