// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = BotChannelDirectLineSpeechV0ToV1{}

type BotChannelDirectLineSpeechV0ToV1 struct{}

func (BotChannelDirectLineSpeechV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"cognitive_service_access_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"cognitive_service_location": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"cognitive_account_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"custom_speech_model_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"custom_voice_deployment_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (BotChannelDirectLineSpeechV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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
