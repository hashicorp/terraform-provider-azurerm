package migration

import (
	"context"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-05-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = CosmosDBAccountV0toV1{}

type CosmosDBAccountV0toV1 struct{}

func (c CosmosDBAccountV0toV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"offer_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"analytical_storage": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"schema_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"capacity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_throughput_limit": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		// per Microsoft's documentation, as of April 1 2023 the default minimal TLS version for all new accounts is 1.2
		"minimal_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		// Per Documentation: "The default identity needs to be explicitly set by the users." This should not be optional without a default anymore.
		// DOC: https://learn.microsoft.com/en-us/java/api/com.azure.resourcemanager.cosmos.models.databaseaccountupdateparameters?view=azure-java-stable#method-details
		"default_identity_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "FirstPartyIdentity",
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(cosmosdb.DatabaseAccountKindGlobalDocumentDB),
		},

		"ip_range_filter": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"free_tier_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"analytical_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"automatic_failover_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: !features.FourPointOhBeta(),
		},

		"key_vault_key_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"consistency_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"consistency_level": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					// This value can only change if the 'consistency_level' is set to 'BoundedStaleness'
					"max_interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  5,
					},

					// This value can only change if the 'consistency_level' is set to 'BoundedStaleness'
					"max_staleness_prefix": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  100,
					},
				},
			},
		},

		"geo_location": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"location": commonschema.LocationWithoutForceNew(),

					"failover_priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"zone_redundant": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"capabilities": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"is_virtual_network_filter_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"virtual_network_rule": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"ignore_missing_vnet_service_endpoint": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"access_key_metadata_writes_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"local_authentication_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"mongo_server_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"multiple_write_locations_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"network_acl_bypass_for_azure_services": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"network_acl_bypass_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true, Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"partition_merge_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"burst_capacity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"backup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					// Though `tier` has the default value `Continuous30Days` but `tier` is only for the backup type `Continuous`. So the default value isn't added in the property schema.
					"tier": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"interval_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"retention_in_hours": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"storage_redundancy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"cors_rule": common.SchemaCorsRule(),

		"connection_strings": {
			Type:      pluginsdk.TypeList,
			Computed:  true,
			Sensitive: true,
			Elem: &pluginsdk.Schema{
				Type:      pluginsdk.TypeString,
				Sensitive: true,
			},
		},

		"enable_multiple_write_locations": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"enable_free_tier": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"enable_automatic_failover": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"restore": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_cosmosdb_account_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"restore_timestamp_in_utc": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"database": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
								},

								"collection_names": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"gremlin_database": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
								},

								"graph_names": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"tables_to_restore": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"read_endpoints": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"write_endpoints": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_readonly_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_readonly_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_sql_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_sql_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_readonly_sql_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_readonly_sql_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_mongodb_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_mongodb_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_readonly_mongodb_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_readonly_mongodb_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (c CosmosDBAccountV0toV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		ipfString := rawState["ip_range_filter"].(string)
		ipfSet := make([]string, 0)
		if ipfString != "" {
			ipfSet = strings.Split(ipfString, ",")
		}
		rawState["ip_range_filter"] = ipfSet

		return rawState, nil
	}
}
