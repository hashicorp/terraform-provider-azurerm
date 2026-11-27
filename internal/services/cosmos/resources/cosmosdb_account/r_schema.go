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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbAccountSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-z0-9]{3,50}$"),
				"Cosmos DB Account name must be 3 - 50 characters long, contain only lowercase letters, numbers and hyphens.",
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		// resource fields
		"offer_type": {
			Description:  "The Offer Type to use for this CosmosDB Account; currently, this can only be set to `Standard`.",
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
						Description:  "The schema type of the Analytical Storage for this Cosmos DB account.",
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
						Description:  "The total throughput limit imposed on this Cosmos DB account (RU/s). `-1` means no limit.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(-1),
					},
				},
			},
		},

		// per Microsoft's documentation, as of April 1 2023 the default minimal TLS version for all new accounts is 1.2
		"minimal_tls_version": {
			Description: "The minimal TLS version for the CosmosDB account.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(cosmosdb.MinimalTlsVersionTlsOneTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(cosmosdb.MinimalTlsVersionTlsOneTwo),
			}, false),
		},

		"create_mode": {
			Description:  "The creation mode for the CosmosDB Account.",
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
			Description:  "The Kind of CosmosDB to create.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(cosmosdb.DatabaseAccountKindGlobalDocumentDB),
			ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForDatabaseAccountKind(), false),
		},

		"ip_range_filter": {
			Description: "A set of IP addresses or IP address ranges in CIDR form to be included as the allowed list of client IPs for a given database account. For example `[\"55.0.1.0/24\", \"55.0.2.0/24\"]`.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.Any(validation.IsCIDR, validation.IsIPv4Address),
			},
		},

		"free_tier_enabled": {
			Description: "Enable the Free Tier pricing option for this Cosmos DB account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"analytical_storage_enabled": {
			Description: "Enable Analytical Storage option for this Cosmos DB account. Enabling and then disabling analytical storage forces a new resource to be created.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for this CosmosDB account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"automatic_failover_enabled": {
			Description: "Enable automatic failover for this Cosmos DB account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"key_vault_key_id": {
			Description:  "A versionless Key Vault Key ID for CMK encryption.",
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
						Description:  "The Consistency Level to use for this CosmosDB Account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForDefaultConsistencyLevel(), false),
					},

					// This value can only change if the 'consistency_level' is set to 'BoundedStaleness'
					"max_interval_in_seconds": {
						Description:      "When used with the Bounded Staleness consistency level, this value represents the time amount of staleness (in seconds) tolerated. The accepted range for this value is `5` - `86400` (1 day).",
						Type:             pluginsdk.TypeInt,
						Optional:         true,
						Default:          5,
						ValidateFunc:     validation.IntBetween(5, 86400), // single region values
						DiffSuppressFunc: suppressConsistencyPolicyStalenessConfiguration,
					},

					// This value can only change if the 'consistency_level' is set to 'BoundedStaleness'
					"max_staleness_prefix": {
						Description:      "When used with the Bounded Staleness consistency level, this value represents the number of stale requests tolerated. The accepted range for this value is `10` - `2147483647`.",
						Type:             pluginsdk.TypeInt,
						Optional:         true,
						Default:          100,
						ValidateFunc:     validation.IntBetween(10, math.MaxInt32), // single region values
						DiffSuppressFunc: suppressConsistencyPolicyStalenessConfiguration,
					},
				},
			},
		},

		"geo_location": {
			Description: "A `geo_location` resource, used to define where data should be replicated with the `failover_priority` 0 specifying the primary location.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Set:         resourceAzureRMCosmosDBAccountGeoLocationHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"location": commonschema.LocationWithoutForceNew(),

					"failover_priority": {
						Description:  "The failover priority of the region. A failover priority of `0` indicates a write region. The maximum value for a failover priority = (total number of regions - 1). Failover priority values must be unique for each of the regions in which the database account exists.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"zone_redundant": {
						Description: "Whether zone redundancy should be enabled for this region.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"capabilities": {
			Description: "The capabilities which should be enabled for this Cosmos DB account.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Set:         resourceAzureRMCosmosDBAccountCapabilitiesHash,
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
		},

		"virtual_network_filter_enabled": {
			Description: "Enables Virtual Network filtering for this Cosmos DB account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_network_rule": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Set:      resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description:  "The ID of the Virtual Network subnet.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"ignore_missing_vnet_service_endpoint": {
						Description: "If set to true, the specified subnet will be added as a Virtual Network rule even if its CosmosDB service endpoint is not active.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"access_key_metadata_writes_enabled": {
			Description: "Whether write operations on metadata resources (databases, containers, throughput) via account keys are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"local_authentication_enabled": {
			Description: "Whether local authentication is enabled, when disabled only MSI and AAD can be used exclusively for authentication.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"mongo_server_version": {
			Description:  "The Server Version of a MongoDB account.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForServerVersion(), false),
		},

		"multiple_write_locations_enabled": {
			Description: "Enable multiple write locations for this Cosmos DB account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"network_acl_bypass_for_azure_services": {
			Description: "If Azure services can bypass ACLs.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"network_acl_bypass_ids": {
			Description: "The list of resource IDs for Network ACL Bypass for this Cosmos DB account.",
			Type:        pluginsdk.TypeList,
			Optional:    true, Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},

		"partition_merge_enabled": {
			Description: "Whether partition merge on the Cosmos DB account is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"burst_capacity_enabled": {
			Description: "Enable burst capacity for this Cosmos DB account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"backup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Description:  "The type of the `backup`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForBackupPolicyType(), false),
					},

					// Though `tier` has the default value `Continuous30Days` but `tier` is only for the backup type `Continuous`. So the default value isn't added in the property schema.
					"tier": {
						Description:  "The continuous backup tier.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForContinuousTier(), false),
					},

					"interval_in_minutes": {
						Description:  "The interval in minutes between two backups.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(60, 1440),
					},

					"retention_in_hours": {
						Description:  "The time in hours that each backup is retained.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(8, 720),
					},

					"storage_redundancy": {
						Description:  "The storage redundancy is used to indicate the type of backup residency.",
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
						Description:  "The resource ID of the restorable database account from which the restore has to be initiated. The example is `/subscriptions/{subscriptionId}/providers/Microsoft.DocumentDB/locations/{location}/restorableDatabaseAccounts/{restorableDatabaseAccountName}`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: restorables.ValidateRestorableDatabaseAccountID,
					},

					"restore_timestamp_in_utc": {
						Description:  "The creation time of the database or the collection (Datetime Format `RFC 3339`).",
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
									Description:  "The database name for the restore request.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"collection_names": {
									Description: "A list of the collection names for the restore request.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									ForceNew:    true,
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
									Description:  "The Gremlin Database name for the restore request.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.CosmosEntityName,
								},

								"graph_names": {
									Description: "A list of the Graph names for the restore request.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									ForceNew:    true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.CosmosEntityName,
									},
								},
							},
						},
					},

					"tables_to_restore": {
						Description: "A list of specific tables available for restore.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
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
			Description: "The endpoint used to connect to the CosmosDB account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"read_endpoints": {
			Description: "A list of read endpoints available for this CosmosDB account.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"write_endpoints": {
			Description: "A list of write endpoints available for this CosmosDB account.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"primary_key": {
			Description: "The Primary key for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_key": {
			Description: "The Secondary key for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_readonly_key": {
			Description: "The Primary read-only Key for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_readonly_key": {
			Description: "The Secondary read-only key for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_sql_connection_string": {
			Description: "Primary SQL connection string for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_sql_connection_string": {
			Description: "Secondary SQL connection string for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_readonly_sql_connection_string": {
			Description: "Primary readonly SQL connection string for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_readonly_sql_connection_string": {
			Description: "Secondary readonly SQL connection string for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_mongodb_connection_string": {
			Description: "Primary Mongodb connection string for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_mongodb_connection_string": {
			Description: "Secondary Mongodb connection string for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_readonly_mongodb_connection_string": {
			Description: "Primary readonly Mongodb connection string for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_readonly_mongodb_connection_string": {
			Description: "Secondary readonly Mongodb connection string for the CosmosDB Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		s["is_virtual_network_filter_enabled"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Description:   "Enables Virtual Network filtering for this Cosmos DB account.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"virtual_network_filter_enabled"},
			Deprecated:    "`is_virtual_network_filter_enabled` has been deprecated in favour of `virtual_network_filter_enabled` and will be removed in v6.0 of the AzureRM Provider",
		}

		s["virtual_network_filter_enabled"] = &pluginsdk.Schema{
			Description:   "Enables Virtual Network filtering for this Cosmos DB account.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"is_virtual_network_filter_enabled"},
		}
	}

	return s
}
