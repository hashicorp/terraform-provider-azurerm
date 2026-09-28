// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = BotWebAppV0ToV1{}

type BotWebAppV0ToV1 struct{}

func (BotWebAppV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"microsoft_app_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"microsoft_app_tenant_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"microsoft_app_user_assigned_identity_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (BotWebAppV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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
