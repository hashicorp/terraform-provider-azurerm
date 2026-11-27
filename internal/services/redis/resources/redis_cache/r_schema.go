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
			Description: "The name of the Redis instance.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"capacity": {
			Description: "The size of the Redis Cache to deploy.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"family": {
			Description:  "The SKU family/pricing group to use.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(redisresources.PossibleValuesForSkuFamily(), false),
		},

		"sku_name": {
			Description:  "The SKU of Redis to use.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(redisresources.PossibleValuesForSkuName(), false),
		},

		"minimum_tls_version": {
			Description: "The minimum TLS version.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(redisresources.TlsVersionOnePointTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(redisresources.TlsVersionOnePointTwo),
			}, false),
		},

		"shard_count": {
			Description: "*Only available when using the Premium SKU* The number of Shards to create on the Redis Cluster.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
		},

		"non_ssl_port_enabled": {
			Description: "Enable the non-SSL port (6379) - disabled by default.",
			Type:        pluginsdk.TypeBool,

			Optional: true, Default: false,
		},

		"subnet_id": {
			Description:  "*Only available when using the Premium SKU* The ID of the Subnet within which the Redis Cache should be deployed. This Subnet must only contain Azure Cache for Redis instances without any other type of resources.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"private_static_ip_address": {
			Description: "The Static IP Address to assign to the Redis Cache when hosted inside the Virtual Network. This argument implies the use of `subnet_id`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
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
						Description: "Enable Microsoft Entra (AAD) authentication.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"maxclients": {
						Description: "Returns the max number of connected clients at the same time.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"maxmemory_delta": {
						Description: "The max-memory delta for this Redis instance.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						// Note: O+C this gets a variable default based on the cache's total memory. This can remain O+C as it can be updated without issue
						Computed: true,
					},

					"maxmemory_reserved": {
						Description: "Value in megabytes reserved for non-cache usage e.g. failover.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						// Note: O+C this gets a variable default based on the cache's total memory. This can remain O+C as it can be updated without issue
						Computed: true,
					},

					"maxmemory_policy": {
						Description:  "How Redis will select what to remove when `maxmemory` is reached.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "volatile-lru",
						ValidateFunc: validate.MaxMemoryPolicy,
					},

					"maxfragmentationmemory_reserved": {
						Description: "Value in megabytes reserved to accommodate for memory fragmentation.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						// Note: O+C this gets a variable default based on the cache's total memory. This can remain O+C as it can be updated without issue
						Computed: true,
					},

					"data_persistence_authentication_method": {
						Description: "Preferred auth method to communicate to Storage Account used for data persistence.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"SAS",
							"ManagedIdentity",
						}, false),
					},

					"rdb_backup_enabled": {
						Description: "Is Backup Enabled? Only supported on Premium SKUs.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"rdb_backup_frequency": {
						Description:  "The Backup Frequency in Minutes. Only supported on Premium SKUs.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validate.CacheBackupFrequency,
					},

					"rdb_backup_max_snapshot_count": {
						Description: "The maximum number of snapshots to create as a backup. Only supported for Premium SKUs.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
					},

					"rdb_storage_connection_string": {
						Description: "The Connection String to the Storage Account. Only supported for Premium SKUs. In the format: `DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.example.primary_blob_endpoint};AccountName=${azurerm_storage_account.example.name};AccountKey=${azurerm_storage_account.example.primary_access_key}`.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Sensitive:   true,
					},

					"notify_keyspace_events": {
						Description: "Keyspace notifications allows clients to subscribe to Pub/Sub channels in order to receive events affecting the Redis data set in some way. [Reference](https://redis.io/topics/notifications#configuration).",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"aof_backup_enabled": {
						Description: "Enable or disable AOF persistence for this Redis Cache.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"aof_storage_connection_string_0": {
						Description: "First Storage Account connection string for AOF persistence.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Sensitive:   true,
					},

					"aof_storage_connection_string_1": {
						Description: "Second Storage Account connection string for AOF persistence.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Sensitive:   true,
					},
					"authentication_enabled": {
						Description: "If set to `false`, the Redis instance will be accessible without authentication.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"storage_account_subscription_id": {
						Description:  "The ID of the Subscription containing the Storage Account.",
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
						Description:      "The Weekday name.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.IsDayOfTheWeek(true),
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"maintenance_window": {
						Description:  "The ISO 8601 timespan which specifies the amount of time the Redis Cache can be updated.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "PT5H",
						ValidateFunc: validation.ISO8601Duration,
					},

					"start_hour_utc": {
						Description:  "The Start Hour for maintenance in UTC - possible values range from `0 - 23`.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 23),
					},
				},
			},
		},

		"hostname": {
			Description: "The Hostname of the Redis Instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"port": {
			Description: "The non-SSL Port of the Redis Instance.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"ssl_port": {
			Description: "The SSL Port of the Redis Instance.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"primary_access_key": {
			Description: "The Primary Access Key for the Redis Instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_access_key": {
			Description: "The Secondary Access Key for the Redis Instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_connection_string": {
			Description: "The primary connection string of the Redis Instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_connection_string": {
			Description: "The secondary connection string of the Redis Instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"public_network_access_enabled": {
			Description: "Whether or not public network access is allowed for this Redis Cache. `true` means this resource could be accessed by both public and Private Endpoint. `false` means only Private Endpoint access is allowed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		// todo: investigate the difference between `replicas_per_master` and `replicas_per_primary` - are these
		// the same field that's been renamed ala Redis? https://github.com/Azure/azure-rest-api-specs/pull/13005
		"replicas_per_master": {
			Description: "Amount of replicas to create per master for this Redis Cache.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: O+C returns a value when `replicas_per_primary` is set - this can remain as there is no issue updating it
			Computed: true,
			// Can't make more than 3 replicas in portal, assuming it's a limitation
			ValidateFunc: validation.IntBetween(1, 3),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"replicas_per_primary": {
			Description: "Amount of replicas to create per primary for this Redis Cache. If both `replicas_per_primary` and `replicas_per_master` are set, they need to be equal.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: O+C returns a value when `replicas_per_master` is set - this can remain as there is no issue updating it
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 3),
		},

		"tenant_settings": {
			Description: "A mapping of tenant settings to assign to the resource.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"redis_version": {
			Description:  "Redis version. Only major version needed.",
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
			Description: "Whether access key authentication is enabled? `active_directory_authentication_enabled` must be set to `true` to disable access key authentication.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tags": commonschema.Tags(),
	}
}
