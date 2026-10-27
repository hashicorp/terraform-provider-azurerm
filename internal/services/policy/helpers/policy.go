// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
)

func GetPolicyDefinitionByDisplayName(ctx context.Context, client *policy.DefinitionsClient, displayName, managementGroupName string,
	builtInOnly bool,
) (policy.Definition, error) {
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

func GetPolicyDefinitionByName(ctx context.Context, client *policy.DefinitionsClient, name, managementGroupName string) (res policy.Definition, err error) {
	if managementGroupName == "" {
		res, err = client.GetBuiltIn(ctx, name)
		if response.WasNotFound(res.Response.Response) {
			res, err = client.Get(ctx, name)
		}
	} else {
		res, err = client.GetAtManagementGroup(ctx, name, managementGroupName)
	}

	return res, err
}

func FlattenParameterDefinitionsValueToStringTrack1(input map[string]*policy.ParameterDefinitionsValue) (string, error) {
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

func ExpandParameterValuesValueFromString(jsonString string) (map[string]policyassignments.ParameterValuesValue, error) {
	var result map[string]policyassignments.ParameterValuesValue

	err := json.Unmarshal([]byte(jsonString), &result)

	return result, err
}

func FlattenParameterValuesValueToStringV2(input *map[string]policyassignments.ParameterValuesValue) (string, error) {
	if input == nil || *input == nil {
		return "", nil
	}
	bs, err := json.Marshal(input)
	return string(bs), err
}

func GetPolicyRoleDefinitionIDs(ruleStr string) (res []string, err error) {
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
