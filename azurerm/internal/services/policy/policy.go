package policy

import ("context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
)

func getPolicySetDefinitionByName(ctx context.Context, client *policy.SetDefinitionsClient, name string, managementGroupID string) (res policy.SetDefinition, err error) {
	if managementGroupID == "" {
		res, err = client.Get(ctx, name)
	} else {
		res, err = client.GetAtManagementGroup(ctx, name, managementGroupID)
	}

	return res, err
}

func getPolicySetDefinitionByDisplayName(ctx context.Context, client *policy.SetDefinitionsClient, displayName, managementGroupID string) (policy.SetDefinition, error) {
	var setDefinitions policy.SetDefinitionListResultIterator
	var err error

	if managementGroupID != "" {
		setDefinitions, err = client.ListByManagementGroupComplete(ctx, managementGroupID)
	} else {
		setDefinitions, err = client.ListComplete(ctx)
	}
	if err != nil {
		return policy.SetDefinition{}, fmt.Errorf("failed to load Policy Set Definition List: %+v", err)
	}

	var results []policy.SetDefinition
	for setDefinitions.NotDone() {
		def := setDefinitions.Value()
		if def.DisplayName != nil && *def.DisplayName == displayName && def.ID != nil {
			results = append(results, def)
		}

		if err := setDefinitions.NextWithContext(ctx); err != nil {
			return policy.SetDefinition{}, fmt.Errorf("failed to load Policy Set Definition List: %s", err)
		}
	}

	// throw error when we found none
	if len(results) == 0 {
		return policy.SetDefinition{}, fmt.Errorf("failed to load Policy Set Definition List: could not find policy '%s'", displayName)
	}

	// throw error when we found more than one
	if len(results) > 1 {
		return policy.SetDefinition{}, fmt.Errorf("failed to load Policy Set Definition List: found more than one policy set definition '%s'", displayName)
	}
	
	return results[0], nil
}
