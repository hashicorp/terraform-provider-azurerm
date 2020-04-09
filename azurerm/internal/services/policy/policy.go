package policy

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
)

func getPolicyDefinitionByDisplayName(ctx context.Context, client *policy.DefinitionsClient, displayName, managementGroupName string) (policy.Definition, error) {
	var policyDefinitions policy.DefinitionListResultIterator
	var err error

	if managementGroupName != "" {
		policyDefinitions, err = client.ListByManagementGroupComplete(ctx, managementGroupName)
	} else {
		policyDefinitions, err = client.ListComplete(ctx)
	}
	if err != nil {
		return policy.Definition{}, fmt.Errorf("failed to load Policy Definition List: %+v", err)
	}

	var results []policy.Definition
	for policyDefinitions.NotDone() {
		def := policyDefinitions.Value()
		if def.DisplayName != nil && *def.DisplayName == displayName && def.ID != nil {
			results = append(results, def)
		}

		if err := policyDefinitions.NextWithContext(ctx); err != nil {
			return policy.Definition{}, fmt.Errorf("failed to load Policy Definition List: %s", err)
		}
	}

	// we found none
	if len(results) == 0 {
		return policy.Definition{}, fmt.Errorf("failed to load Policy Definition List: could not find policy '%s'", displayName)
	}

	// we found more than one
	if len(results) > 1 {
		return policy.Definition{}, fmt.Errorf("failed to load Policy Definition List: found more than one policy '%s'", displayName)
	}

	return results[0], nil
}

func getPolicyDefinitionByName(ctx context.Context, client *policy.DefinitionsClient, name string, managementGroupName string) (res policy.Definition, err error) {
	if managementGroupName == "" {
		res, err = client.Get(ctx, name)
	} else {
		res, err = client.GetAtManagementGroup(ctx, name, managementGroupName)
	}

	return res, err
}
