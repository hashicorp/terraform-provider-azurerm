// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SnapshotV0ToV1{}

type SnapshotV0ToV1 struct{}

func (SnapshotV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := snapshots.ParseSnapshotIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating the ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}

func (SnapshotV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"create_option": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"source_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"source_resource_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"disk_size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"encryption_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"disk_encryption_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"secret_url": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"source_vault_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"key_encryption_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_url": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"source_vault_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"trusted_launch_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
