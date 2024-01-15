// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiManagementApiOperationPolicyV0ToV1{}

type ApiManagementApiOperationPolicyV0ToV1 struct{}

func (ApiManagementApiOperationPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"operation_id": {
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

func (ApiManagementApiOperationPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old id : /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/apis/api1/operations/operation1/policies/xml
		// new id : /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/apis/api1/operations/operation1
		oldId := rawState["id"].(string)
		newId := strings.TrimSuffix(oldId, "/policies/xml")

		log.Printf("[DEBUG] Updating ID from %q tmakeo %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
