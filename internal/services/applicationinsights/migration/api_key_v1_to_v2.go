// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	apikeys "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiKeyUpgradeV1ToV2{}

type ApiKeyUpgradeV1ToV2 struct{}

func (ApiKeyUpgradeV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return apiKeySchemaForV1AndV2()
}

func (ApiKeyUpgradeV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/microsoft.insights/components/component1/apikeys/key1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/component1/apiKeys/key1
		oldIdRaw := rawState["id"].(string)
		oldId, err := apikeys.ParseApiKeyIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		id := apikeys.NewApiKeyID(oldId.SubscriptionId, oldId.ResourceGroupName, oldId.ComponentName, oldId.KeyId)
		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func apiKeySchemaForV1AndV2() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"application_insights_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"read_permissions": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"write_permissions": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"api_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}
