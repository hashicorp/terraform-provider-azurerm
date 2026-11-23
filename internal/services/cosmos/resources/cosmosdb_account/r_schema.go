// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_account

import (
	"math"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/restorables"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-z0-9]{3,50}$"),
				"Cosmos DB Account name must be 3 - 50 characters long, contain only lowercase letters, numbers and hyphens.",
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		// resource fields
		"offer_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForDatabaseAccountOfferType(), false),
		},

		"analytical_storage": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"schema_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForAnalyticalStorageSchemaType(), false),
					},
				},
			},
		},

		"capacity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_throughput_limit": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(-1),
					},
				},
			},
		},

		// per Microsoft's documentation, as of April 1 2023 the default minimal TLS version for all new accounts is 1.2
		"minimal_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(cosmosdb.MinimalTlsVersionTlsOneTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(cosmosdb.MinimalTlsVersionTlsOneTwo),
			}, false),
		},

		"create_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForCreateMode(), false),
		},

		// Per Documentation: "The default identity needs to be explicitly set by the users." This should not be optional without a default anymore.
		// DOC: https://learn.microsoft.com/java/api/com.azure.resourcemanager.cosmos.models.databaseaccountupdateparameters?view=azure-java-stable#method-details
		"default_identity_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "FirstPartyIdentity",
			ValidateFunc: validation.Any(
				validation.StringMatch(regexp.MustCompile(`^UserAssignedIdentity(.)+$`), "user assigned identity must be in the format of: 'UserAssignedIdentity=/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{userAssignedIdentityName}'"),
				validation.StringInSlice([]string{
					"FirstPartyIdentity",
					"SystemAssignedIdentity",
				}, false),
			),
		},

		"kind": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(cosmosdb.DatabaseAccountKindGlobalDocumentDB),
			ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForDatabaseAccountKind(), false),
		},

		"ip_range_filter": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.Any(validation.IsCIDR, validation.IsIPv4Address),
			},
		},

		"free_tier_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
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
			Default:  false,
		},

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeKey),
		},

		"consistency_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"consistency_level": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForDefaultConsistencyLevel(), false),
					},

					// This value can only change if the 'consistency_level' is set to 'BoundedStaleness'
					"max_interval_in_seconds": {
						Type:             pluginsdk.TypeInt,
						Optional:         true,
						Default:          5,
						DiffSuppressFunc: suppressConsistencyPolicyStalenessConfiguration,
						ValidateFunc:     validation.IntBetween(5, 86400), // single region values
					},

					// This value can only change if the 'consistency_level' is set to 'BoundedStaleness'
					"max_staleness_prefix": {
						Type:             pluginsdk.TypeInt,
						Optional:         true,
						Default:          100,
						DiffSuppressFunc: suppressConsistencyPolicyStalenessConfiguration,
						ValidateFunc:     validation.IntBetween(10, math.MaxInt32), // single region values
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
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"zone_redundant": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
			Set: resourceAzureRMCosmosDBAccountGeoLocationHash,
		},

		"capabilities": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(databaseAccountCapabilitiesEnableAggregationPipeline),
							string(databaseAccountCapabilitiesEnableCassandra),
							string(databaseAccountCapabilitiesEnableGremlin),
							string(databaseAccountCapabilitiesEnableTable),
							string(databaseAccountCapabilitiesEnableServerless),
							string(databaseAccountCapabilitiesEnableMongo),
							string(databaseAccountCapabilitiesEnableMongo16MBDocumentSupport),
							string(databaseAccountCapabilitiesMongoDBv34),
							string(databaseAccountCapabilitiesMongoEnableDocLevelTTL),
							string(databaseAccountCapabilitiesDeleteAllItemsByPartitionKey),
							string(databaseAccountCapabilitiesDisableRateLimitingResponses),
							string(databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36),
							string(databaseAccountCapabilitiesEnableMongoRetryableWrites),
							string(databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl),
							string(databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs),
							string(databaseAccountCapabilitiesEnableNoSqlVectorSearch),
							string(databaseAccountCapabilitiesEnableNoSqlFullTextSearch),
							string(databaseAccountCapabilitiesEnableTtlOnCustomPath),
							string(databaseAccountCapabilitiesEnablePartialUniqueIndex),
							string(databaseAccountCapabilitiesEnableFabricNetworkAclBypass),
						}, false),
					},
				},
			},
			Set: resourceAzureRMCosmosDBAccountCapabilitiesHash,
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
			Set: resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash,
		},

		"access_key_metadata_writes_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"mongo_server_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForServerVersion(), false),
		},

		"multiple_write_locations_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"network_acl_bypass_for_azure_services": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"network_acl_bypass_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true, Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
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
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForBackupPolicyType(), false),
					},

					// Though `tier` has the default value `Continuous30Days` but `tier` is only for the backup type `Continuous`. So the default value isn't added in the property schema.
					"tier": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForContinuousTier(), false),
					},

					"interval_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(60, 1440),
					},

					"retention_in_hours": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(8, 720),
					},

					"storage_redundancy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForBackupStorageRedundancy(), false),
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"cors_rule": common.SchemaCorsRule(),

		"restore": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_cosmosdb_account_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: restorables.ValidateRestorableDatabaseAccountID,
					},

					"restore_timestamp_in_utc": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},

					"database": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"collection_names": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.CosmosEntityName,
								},

								"graph_names": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.CosmosEntityName,
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
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.CosmosEntityName,
						},
					},
				},
			},
		},

		// computed
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
