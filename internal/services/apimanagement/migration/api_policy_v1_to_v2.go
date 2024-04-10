// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiManagementApiPolicyV1ToV2{}

type ApiManagementApiPolicyV1ToV2 struct{}

func (ApiManagementApiPolicyV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"api_management_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"api_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"xml_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"xml_link": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (ApiManagementApiPolicyV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old id : /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/apis/exampleId/policies/policy
		// new id : /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/apis/exampleId
		oldId := rawState["id"].(string)

		// Prior to v3.70.0 of Terraform Provider, after importing resource, the id in state file ends with "/policies/policy", the id in state file ends with "/policies/xml" for creating resource by Terraform.
		// So after migrating pandora SDK (starting from v3.70.0), these two cases need to be migrated.
		// In ApiManagementApiPolicyV0ToV1, only the case where the ID ends with "/policies/xml" is processed, so the case where the ID ends with "/policies/policy" is processed here to solve the parse id error.
		newId := strings.TrimSuffix(oldId, "/policies/policy")
		parsed, err := apipolicy.ParseApiID(newId)
		if err != nil {
			return rawState, err
		}
		newId = parsed.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
