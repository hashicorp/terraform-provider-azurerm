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

type KustoDatabaseDataConnectionEventGridV0ToV1 struct{}

func (s KustoDatabaseDataConnectionEventGridV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"eventhub_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"eventhub_consumer_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"blob_storage_event_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"skip_first_record": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"table_name": {
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

		"eventgrid_resource_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"managed_identity_resource_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s KustoDatabaseDataConnectionEventGridV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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
