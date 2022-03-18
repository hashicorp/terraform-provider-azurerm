package blueprints

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func blueprintAssignmentCreateStateRefreshFunc(ctx context.Context, client *blueprint.AssignmentsClient, scope, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
		}
		if resp.ProvisioningState == blueprint.Failed {
			return resp, string(resp.ProvisioningState), fmt.Errorf("Blueprint Assignment provisioning entered a Failed state.")
		}

		return resp, string(resp.ProvisioningState), nil
	}
}

func blueprintAssignmentDeleteStateRefreshFunc(ctx context.Context, client *blueprint.AssignmentsClient, scope, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			} else {
				return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
			}
		}

		return resp, string(resp.ProvisioningState), nil
	}
}

func normalizeAssignmentParameterValuesJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	var values map[string]*blueprint.ParameterValue
	if err := json.Unmarshal([]byte(jsonString.(string)), &values); err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}

	b, _ := json.Marshal(values)
	return string(b)
}

func normalizeAssignmentResourceGroupValuesJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	var values map[string]*blueprint.ResourceGroupValue
	if err := json.Unmarshal([]byte(jsonString.(string)), &values); err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}

	b, _ := json.Marshal(values)
	return string(b)
}

func expandArmBlueprintAssignmentParameters(input string) map[string]*blueprint.ParameterValue {
	var result map[string]*blueprint.ParameterValue
	// the string has been validated by the schema, therefore the error is ignored here, since it will never happen.
	_ = json.Unmarshal([]byte(input), &result)
	return result
}

func expandArmBlueprintAssignmentResourceGroups(input string) map[string]*blueprint.ResourceGroupValue {
	var result map[string]*blueprint.ResourceGroupValue
	// the string has been validated by the schema, therefore the error is ignored here, since it will never happen.
	_ = json.Unmarshal([]byte(input), &result)
	return result
}

func expandArmBlueprintAssignmentIdentity(input []interface{}) (*blueprint.ManagedServiceIdentity, error) {
	expanded, err := identity.ExpandUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := blueprint.ManagedServiceIdentity{
		Type:                   blueprint.ManagedServiceIdentityType(string(expanded.Type)),
		UserAssignedIdentities: make(map[string]*blueprint.UserAssignedIdentity),
	}
	for k := range expanded.IdentityIds {
		out.UserAssignedIdentities[k] = &blueprint.UserAssignedIdentity{
			// intentionally empty
		}
	}
	return &out, nil
}

func flattenArmBlueprintAssignmentIdentity(input *blueprint.ManagedServiceIdentity) (*[]interface{}, error) {
	var transform *identity.UserAssignedMap

	if input != nil {
		transform = &identity.UserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenUserAssignedMap(transform)
}

func flattenArmBlueprintAssignmentParameters(input map[string]*blueprint.ParameterValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func flattenArmBlueprintAssignmentResourceGroups(input map[string]*blueprint.ResourceGroupValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
