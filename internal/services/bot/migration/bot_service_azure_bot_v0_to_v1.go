// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = BotServiceAzureBotV0ToV1{}

type BotServiceAzureBotV0ToV1 struct{}

func (BotServiceAzureBotV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"microsoft_app_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"developer_app_insights_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"developer_app_insights_api_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"developer_app_insights_application_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"cmk_key_vault_key_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"microsoft_app_msi_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"microsoft_app_tenant_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"microsoft_app_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"luis_app_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"luis_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"streaming_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"icon_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "https://docs.botframework.com/static/devportal/client/images/bot-framework-default.png",
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (BotServiceAzureBotV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// IDs imported while this resource parsed them with the legacy resourceids.ParseAzureResourceID
		// can contain non-canonically cased static segments (e.g. `resourcegroups`, `microsoft.botservice`),
		// which the case-sensitive SDK parser rejects - normalise them to the canonical casing
		oldId := rawState["id"].(string)
		id, err := commonids.ParseBotServiceIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
