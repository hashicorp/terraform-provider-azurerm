// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ManagedDiskV0ToV1{}

type ManagedDiskV0ToV1 struct{}

func (ManagedDiskV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := disks.ParseDiskIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating the ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}

func (ManagedDiskV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"create_option": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"edge_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"logical_sector_size": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"source_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"source_resource_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"image_reference_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"gallery_image_reference_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"disk_size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"disk_iops_read_write": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"disk_mbps_read_write": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"disk_iops_read_only": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"disk_mbps_read_only": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"disk_encryption_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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

		"network_access_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"disk_access_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"max_shares": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},

		"trusted_launch_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"secure_vm_disk_encryption_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"security_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"hyper_v_generation": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"on_demand_bursting_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"zone": {
			Type:     schema.TypeString,
			Optional: true,
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
