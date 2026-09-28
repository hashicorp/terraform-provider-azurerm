// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SynapseIntegrationRuntimeAzureV1ToV2{}

type SynapseIntegrationRuntimeAzureV1ToV2 struct{}

func (SynapseIntegrationRuntimeAzureV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"synapse_workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
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

func (SynapseIntegrationRuntimeAzureV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// the V0 -> V1 upgrader normalised the states that existed at the time, but IDs imported since
		// were still parsed with the legacy resourceids.ParseAzureResourceID and can contain
		// non-canonically cased static segments (e.g. `resourcegroups`, `microsoft.synapse`), which the
		// case-sensitive SDK parser rejects - normalise them to the canonical casing
		oldId := rawState["id"].(string)
		id, err := integrationruntimes.ParseIntegrationRuntimeIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
