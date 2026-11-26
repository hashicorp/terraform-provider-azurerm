// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/scripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoDatabaseScriptV0ToV1 struct{}

func (s KustoDatabaseScriptV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"database_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"continue_on_errors_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"force_an_update_when_value_changed": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			RequiredWith: []string{"sas_token"},
			ExactlyOneOf: []string{"url", "script_content"},
		},

		"sas_token": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			RequiredWith: []string{"url"},
		},

		"script_content": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ExactlyOneOf: []string{"url", "script_content"},
		},
	}
}

func (s KustoDatabaseScriptV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		oldId := rawState["id"].(string)
		newId, err := scripts.ParseScriptIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
