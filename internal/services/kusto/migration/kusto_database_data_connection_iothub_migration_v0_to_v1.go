// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoDatabaseDataConnectionIotHub0ToV1 struct{}

func (s KustoDatabaseDataConnectionIotHub0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"database_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"iothub_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"consumer_group": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"shared_access_policy_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"table_name": {
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			Optional: true,
		},

		"mapping_rule_name": {
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			Optional: true,
		},

		"data_format": {
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			Optional: true,
		},

		"database_routing_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"event_system_properties": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (s KustoDatabaseDataConnectionIotHub0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.DataConnectionIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
