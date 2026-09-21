// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = BotChannelDirectlineV0ToV1{}

type BotChannelDirectlineV0ToV1 struct{}

func (BotChannelDirectlineV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"bot_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"site": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"user_upload_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"endpoint_parameters_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"storage_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"v1_allowed": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"v3_allowed": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"enhanced_authentication_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  false,
						Optional: true,
					},

					"trusted_origins": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"key": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"key2": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"extension_key_1": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"extension_key_2": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (BotChannelDirectlineV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// IDs imported while this resource parsed them with the legacy resourceids.ParseAzureResourceID
		// can contain non-canonically cased static segments (e.g. `resourcegroups`, `microsoft.botservice`),
		// which the case-sensitive SDK parser rejects - normalise them to the canonical casing
		oldId := rawState["id"].(string)
		id, err := commonids.ParseBotServiceChannelIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
