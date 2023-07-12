// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01/redis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

var _ pluginsdk.StateUpgrade = RedisCacheV0ToV1{}

type RedisCacheV0ToV1 struct{}

func (RedisCacheV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"zones": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"capacity": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"family": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"minimum_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "1.2",
		},

		"shard_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		// TODO 4.0: change this from enable_* to *_enabled
		"enable_non_ssl_port": {
			Type:     pluginsdk.TypeBool,
			Default:  false,
			Optional: true,
		},

		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"private_static_ip_address": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"redis_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"maxclients": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"maxmemory_delta": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"maxmemory_reserved": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"maxmemory_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "volatile-lru",
					},

					"maxfragmentationmemory_reserved": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"rdb_backup_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"rdb_backup_frequency": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"rdb_backup_max_snapshot_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"rdb_storage_connection_string": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"notify_keyspace_events": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"aof_backup_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"aof_storage_connection_string_0": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"aof_storage_connection_string_1": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},
					"enable_authentication": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"patch_schedule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day_of_week": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"maintenance_window": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "PT5H",
					},

					"start_hour_utc": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"ssl_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"replicas_per_master": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"identity": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"replicas_per_primary": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"tenant_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"redis_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
				n := strings.Split(old, ".")
				if len(n) >= 1 {
					newMajor := n[0]
					return new == newMajor
				}
				return false
			},
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (RedisCacheV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := redis.ParseRediIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
