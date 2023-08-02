// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/sourcecontrol"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationSourceControlV0ToV1 struct{}

func (s AutomationSourceControlV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"repository_url": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"branch": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"folder_path": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"automatic_sync": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"publish_runbook_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"source_control_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"security": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"token": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"refresh_token": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"token_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func (s AutomationSourceControlV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		newId, err := sourcecontrol.ParseSourceControlIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
