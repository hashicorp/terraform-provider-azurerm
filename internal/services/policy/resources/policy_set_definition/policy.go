// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_set_definition

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
)

func getPolicySetDefinitionByName(ctx context.Context, client *policy.SetDefinitionsClient, name, managementGroupID string) (res policy.SetDefinition, err error) {
	if managementGroupID == "" {
		res, err = client.GetBuiltIn(ctx, name)
		if response.WasNotFound(res.Response.Response) {
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

func flattenParameterValuesValueToStringTrack1(input map[string]*policy.ParameterValuesValue) (string, error) {
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
