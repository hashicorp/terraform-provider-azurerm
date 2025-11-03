// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataFactoryIntegrationRuntimeAzureV0ToV1 struct{}

var _ pluginsdk.StateUpgrade = DataFactoryIntegrationRuntimeAzureV0ToV1{}

func (DataFactoryIntegrationRuntimeAzureV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"data_factory_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cleanup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"compute_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"core_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"time_to_live_min": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"virtual_network_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

func (v DataFactoryIntegrationRuntimeAzureV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// Migration to update ID segment from lowercase to camelCase (integrationruntimename to integrationRuntimeName)

		oldId := rawState["id"].(string)
		parsedId, err := integrationruntimes.ParseIntegrationRuntimeIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		newId := parsedId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
