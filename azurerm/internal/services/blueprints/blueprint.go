package blueprints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprints/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ManagedIdentitySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					Default:  string(blueprint.ManagedServiceIdentityTypeSystemAssigned),
					ValidateFunc: validation.StringInSlice([]string{
						// ManagedServiceIdentityTypeNone is not valid; a valid and privileged Identity is required for the service to apply the changes.
						string(blueprint.ManagedServiceIdentityTypeUserAssigned),
						string(blueprint.ManagedServiceIdentityTypeSystemAssigned),
					}, false),
				},

				"user_assigned_identities": {
					// The API only seems to care about the "key" portion of this struct, which is the ResourceID of the Identity
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validate.UserAssignedIdentityId,
					},
				},

				"principal_id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func blueprintAssignmentCreateStateRefreshFunc(ctx context.Context, client *blueprint.AssignmentsClient, scope, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
		}

		return resp, string(resp.ProvisioningState), nil
	}
}

func blueprintAssignmentDeleteStateRefreshFunc(ctx context.Context, client *blueprint.AssignmentsClient, scope, name string) resource.StateRefreshFunc {
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
