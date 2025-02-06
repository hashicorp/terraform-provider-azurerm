// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/dataconnections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoDatabaseDataConnectionEventHubV0ToV1 struct{}

func (s KustoDatabaseDataConnectionEventHubV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"cluster_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"compression": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"database_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"eventhub_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"event_system_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"consumer_group": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"table_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"identity_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"mapping_rule_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"data_format": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"database_routing_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
	}
}

func (s KustoDatabaseDataConnectionEventHubV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := dataconnections.ParseDataConnectionIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
