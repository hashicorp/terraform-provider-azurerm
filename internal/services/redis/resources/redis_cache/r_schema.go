// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func redisCacheSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"capacity": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"family": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(redisresources.PossibleValuesForSkuFamily(), false),
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(redisresources.PossibleValuesForSkuName(), false),
		},

		"minimum_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(redisresources.TlsVersionOnePointTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(redisresources.TlsVersionOnePointTwo),
			}, false),
		},

		"shard_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"non_ssl_port_enabled": {
			Type:     pluginsdk.TypeBool,
			Default:  false,
			Optional: true,
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"private_static_ip_address": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C in some cases this gets a default value if omitted. This can remain o+c as it is ForceNew and cannot be updated
			Computed: true,
			ForceNew: true,
		},

		"redis_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_directory_authentication_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"maxclients": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"maxmemory_delta": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// Note: O+C this gets a variable default based on the cache's total memory. This can remain O+C as it can be updated without issue
						Computed: true,
					},

					"maxmemory_reserved": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// Note: O+C this gets a variable default based on the cache's total memory. This can remain O+C as it can be updated without issue
						Computed: true,
					},

					"maxmemory_policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "volatile-lru",
						ValidateFunc: validate.MaxMemoryPolicy,
					},

					"maxfragmentationmemory_reserved": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// Note: O+C this gets a variable default based on the cache's total memory. This can remain O+C as it can be updated without issue
						Computed: true,
					},

					"data_persistence_authentication_method": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"SAS",
							"ManagedIdentity",
						}, false),
					},

					"rdb_backup_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"rdb_backup_frequency": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validate.CacheBackupFrequency,
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
					"authentication_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"storage_account_subscription_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
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
						ValidateFunc:     validation.IsDayOfTheWeek(true),
					},

					"maintenance_window": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "PT5H",
						ValidateFunc: validation.ISO8601Duration,
					},

					"start_hour_utc": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 23),
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

		// todo: investigate the difference between `replicas_per_master` and `replicas_per_primary` - are these
		// the same field that's been renamed ala Redis? https://github.com/Azure/azure-rest-api-specs/pull/13005
		"replicas_per_master": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C returns a value when `replicas_per_primary` is set - this can remain as there is no issue updating it
			Computed: true,
			// Can't make more than 3 replicas in portal, assuming it's a limitation
			ValidateFunc: validation.IntBetween(1, 3),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"replicas_per_primary": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C returns a value when `replicas_per_master` is set - this can remain as there is no issue updating it
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 3),
		},

		"tenant_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"redis_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "6",
			ValidateFunc: validation.StringInSlice([]string{"4", "6"}, false),
			DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
				n := strings.Split(old, ".")
				if len(n) >= 1 {
					newMajor := n[0]
					return new == newMajor
				}
				return false
			},
		},

		"access_keys_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}
