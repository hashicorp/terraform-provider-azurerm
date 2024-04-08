// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiManagementPolicyV2ToV3{}

type ApiManagementPolicyV2ToV3 struct{}

func (ApiManagementPolicyV2ToV3) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
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

func (ApiManagementPolicyV2ToV3) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old id : /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/policies/policy
		// new id : /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1
		oldId := rawState["id"].(string)
		newId := strings.TrimSuffix(oldId, "/policies/policy")

		// ensure the casing on the non-user-specifiable segments is correct
		apiManagementId, err := apimanagementservice.ParseServiceID(newId)
		if err != nil {
			return rawState, fmt.Errorf("parsing updated ID %q as an API Management ID: %+v", newId, err)
		}
		newId = apiManagementId.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
