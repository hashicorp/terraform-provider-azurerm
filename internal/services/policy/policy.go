// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" // nolint: staticcheck
	assignments "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func getPolicyDefinitionByDisplayName(ctx context.Context, client *policy.DefinitionsClient, displayName, managementGroupName string,
	builtInOnly bool) (policy.Definition, error) {
	var policyDefinitions policy.DefinitionListResultIterator
	var err error

	if managementGroupName != "" {
		policyDefinitions, err = client.ListByManagementGroupComplete(ctx, managementGroupName, "", nil)
	} else {
		if builtInOnly {
			policyDefinitions, err = client.ListBuiltInComplete(ctx, "", nil)
		} else {
			policyDefinitions, err = client.ListComplete(ctx, "", nil)
		}
	}
	if err != nil {
		return policy.Definition{}, fmt.Errorf("loading Policy Definition List: %+v", err)
	}

	var results []policy.Definition
	for policyDefinitions.NotDone() {
		def := policyDefinitions.Value()
		if def.DisplayName != nil && *def.DisplayName == displayName && def.ID != nil {
			results = append(results, def)
		}

		if err := policyDefinitions.NextWithContext(ctx); err != nil {
			return policy.Definition{}, fmt.Errorf("loading Policy Definition List: %s", err)
		}
	}

	// we found none
	if len(results) == 0 {
		return policy.Definition{}, fmt.Errorf("loading Policy Definition List: could not find policy '%s'. has the policies name changed? list available with `az policy definition list`", displayName)
	}

	// we found more than one
	if len(results) > 1 {
		return policy.Definition{}, fmt.Errorf("loading Policy Definition List: found more than one (%d) policy '%s'", len(results), displayName)
	}

	return results[0], nil
}

func getPolicyDefinitionByName(ctx context.Context, client *policy.DefinitionsClient, name, managementGroupName string) (res policy.Definition, err error) {
	if managementGroupName == "" {
		res, err = client.GetBuiltIn(ctx, name)
		if utils.ResponseWasNotFound(res.Response) {
			res, err = client.Get(ctx, name)
		}
	} else {
		res, err = client.GetAtManagementGroup(ctx, name, managementGroupName)
	}

	return res, err
}

func getPolicySetDefinitionByName(ctx context.Context, client *policy.SetDefinitionsClient, name, managementGroupID string) (res policy.SetDefinition, err error) {
	if managementGroupID == "" {
		res, err = client.GetBuiltIn(ctx, name)
		if utils.ResponseWasNotFound(res.Response) {
			res, err = client.Get(ctx, name)
		}
	} else {
		res, err = client.GetAtManagementGroup(ctx, name, managementGroupID)
	}

	return res, err
}

func getPolicySetDefinitionByDisplayName(ctx context.Context, client *policy.SetDefinitionsClient, displayName, managementGroupID string) (policy.SetDefinition, error) {
	var setDefinitions policy.SetDefinitionListResultIterator
	var err error

	if managementGroupID != "" {
		setDefinitions, err = client.ListByManagementGroupComplete(ctx, managementGroupID, "", nil)
	} else {
		setDefinitions, err = client.ListComplete(ctx, "", nil)
	}
	if err != nil {
		return policy.SetDefinition{}, fmt.Errorf("loading Policy Set Definition List: %+v", err)
	}

	var results []policy.SetDefinition
	for setDefinitions.NotDone() {
		def := setDefinitions.Value()
		if def.DisplayName != nil && *def.DisplayName == displayName && def.ID != nil {
			results = append(results, def)
		}

		if err := setDefinitions.NextWithContext(ctx); err != nil {
			return policy.SetDefinition{}, fmt.Errorf("loading Policy Set Definition List: %s", err)
		}
	}

	// throw error when we found none
	if len(results) == 0 {
		return policy.SetDefinition{}, fmt.Errorf("loading Policy Set Definition List: could not find policy '%s'", displayName)
	}

	// throw error when we found more than one
	if len(results) > 1 {
		return policy.SetDefinition{}, fmt.Errorf("loading Policy Set Definition List: found more than one policy set definition '%s'", displayName)
	}

	return results[0], nil
}

func expandParameterDefinitionsValueFromString(jsonString string) (map[string]*policy.ParameterDefinitionsValue, error) {
	var result map[string]*policy.ParameterDefinitionsValue

	err := json.Unmarshal([]byte(jsonString), &result)

	return result, err
}

func flattenParameterDefinitionsValueToString(input map[string]*policy.ParameterDefinitionsValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	result, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	compactJson := bytes.Buffer{}
	if err := json.Compact(&compactJson, result); err != nil {
		return "", err
	}

	return compactJson.String(), nil
}

func expandParameterValuesValueFromString(jsonString string) (map[string]assignments.ParameterValuesValue, error) {
	var result map[string]assignments.ParameterValuesValue

	err := json.Unmarshal([]byte(jsonString), &result)

	return result, err
}

func flattenParameterValuesValueToString(input map[string]*policy.ParameterValuesValue) (string, error) {
	if input == nil {
		return "", nil
	}

	// no need to call `json.Compact` for the result of `json.Marshal`, it's compacted already
	result, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(result), err
}

func flattenParameterValuesValueToStringV2(input *map[string]assignments.ParameterValuesValue) (string, error) {
	if input == nil || *input == nil {
		return "", nil
	}
	bs, err := json.Marshal(input)
	return string(bs), err
}

func getPolicyRoleDefinitionIDs(ruleStr string) (res []string, err error) {
	type policyRule struct {
		Then struct {
			Details struct {
				RoleDefinitionIds []string `json:"roleDefinitionIds"`
			} `json:"details"`
		} `json:"then"`
	}
	var ins policyRule
	err = json.Unmarshal([]byte(ruleStr), &ins)
	res = ins.Then.Details.RoleDefinitionIds
	return
}
