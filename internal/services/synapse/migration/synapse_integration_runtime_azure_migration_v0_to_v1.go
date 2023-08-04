// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SynapseIntegrationRuntimeAzureV0ToV1 struct{}

func (s SynapseIntegrationRuntimeAzureV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"synapse_workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"compute_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"core_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"time_to_live_min": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
	}
}

func (s SynapseIntegrationRuntimeAzureV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.IntegrationRuntimeIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
