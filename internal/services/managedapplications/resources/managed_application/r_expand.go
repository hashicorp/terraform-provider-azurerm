// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_application

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandManagedApplicationPlan(input []any) *applications.Plan {
	if len(input) == 0 {
		return nil
	}
	plan := input[0].(map[string]any)

	return &applications.Plan{
		Name:          plan["name"].(string),
		Product:       plan["product"].(string),
		Publisher:     plan["publisher"].(string),
		Version:       plan["version"].(string),
		PromotionCode: pointer.To(plan["promotion_code"].(string)),
	}
}

func expandManagedApplicationParameters(d *pluginsdk.ResourceData) (*map[string]any, error) {
	newParams := make(map[string]any)

	if v, ok := d.GetOk("parameter_values"); ok {
		if err := json.Unmarshal([]byte(v.(string)), &newParams); err != nil {
			return nil, fmt.Errorf("unmarshalling `parameter_values`: %+v", err)
		}
	}

	return &newParams, nil
}

func expandManagedApplicationIdentity(input []any) (*applications.Identity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := &applications.Identity{
		Type: pointer.ToEnum[applications.ResourceIdentityType](string(expanded.Type)),
	}

	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		userAssignedIdentities := make(map[string]applications.UserAssignedResourceIdentity)
		for k := range expanded.IdentityIds {
			userAssignedIdentities[k] = applications.UserAssignedResourceIdentity{}
		}

		out.UserAssignedIdentities = &userAssignedIdentities
	}

	return out, nil
}
