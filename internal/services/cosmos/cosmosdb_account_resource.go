// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/customermanagedkeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	keyVaultSuppress "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/suppress"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	managedHsmValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var CosmosDbAccountResourceName = "azurerm_cosmosdb_account"

var connStringPropertyMap = map[string]string{
	"Primary SQL Connection String":                 "primary_sql_connection_string",
	"Secondary SQL Connection String":               "secondary_sql_connection_string",
	"Primary Read-Only SQL Connection String":       "primary_readonly_sql_connection_string",
	"Secondary Read-Only SQL Connection String":     "secondary_readonly_sql_connection_string",
	"Primary MongoDB Connection String":             "primary_mongodb_connection_string",
	"Secondary MongoDB Connection String":           "secondary_mongodb_connection_string",
	"Primary Read-Only MongoDB Connection String":   "primary_readonly_mongodb_connection_string",
	"Secondary Read-Only MongoDB Connection String": "secondary_readonly_mongodb_connection_string",
}

type databaseAccountCapabilities string

const (
	databaseAccountCapabilitiesEnableAggregationPipeline         databaseAccountCapabilities = "EnableAggregationPipeline"
	databaseAccountCapabilitiesEnableCassandra                   databaseAccountCapabilities = "EnableCassandra"
	databaseAccountCapabilitiesEnableGremlin                     databaseAccountCapabilities = "EnableGremlin"
	databaseAccountCapabilitiesEnableTable                       databaseAccountCapabilities = "EnableTable"
	databaseAccountCapabilitiesEnableServerless                  databaseAccountCapabilities = "EnableServerless"
	databaseAccountCapabilitiesEnableMongo                       databaseAccountCapabilities = "EnableMongo"
	databaseAccountCapabilitiesEnableMongo16MBDocumentSupport    databaseAccountCapabilities = "EnableMongo16MBDocumentSupport"
	databaseAccountCapabilitiesMongoDBv34                        databaseAccountCapabilities = "MongoDBv3.4"
	databaseAccountCapabilitiesMongoEnableDocLevelTTL            databaseAccountCapabilities = "mongoEnableDocLevelTTL"
	databaseAccountCapabilitiesDeleteAllItemsByPartitionKey      databaseAccountCapabilities = "DeleteAllItemsByPartitionKey"
	databaseAccountCapabilitiesDisableRateLimitingResponses      databaseAccountCapabilities = "DisableRateLimitingResponses"
	databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36    databaseAccountCapabilities = "AllowSelfServeUpgradeToMongo36"
	databaseAccountCapabilitiesEnableMongoRetryableWrites        databaseAccountCapabilities = "EnableMongoRetryableWrites"
	databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl databaseAccountCapabilities = "EnableMongoRoleBasedAccessControl"
	databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs    databaseAccountCapabilities = "EnableUniqueCompoundNestedDocs"
	databaseAccountCapabilitiesEnableNoSqlVectorSearch           databaseAccountCapabilities = "EnableNoSQLVectorSearch"
	databaseAccountCapabilitiesEnableNoSqlFullTextSearch         databaseAccountCapabilities = "EnableNoSQLFullTextSearch"
	databaseAccountCapabilitiesEnableTtlOnCustomPath             databaseAccountCapabilities = "EnableTtlOnCustomPath"
	databaseAccountCapabilitiesEnablePartialUniqueIndex          databaseAccountCapabilities = "EnablePartialUniqueIndex"
)

/*
	The mapping of capabilities and kinds of cosmosdb account confirmed by service team is as follows:

EnableMongo :                    	MongoDB
EnableCassandra :                	GlobalDocumentDB, Parse
EnableGremlin :                  	GlobalDocumentDB, Parse
EnableTable :                    	GlobalDocumentDB, Parse
EnableAggregationPipeline :      	GlobalDocumentDB, MongoDB, Parse
EnableServerless :               	GlobalDocumentDB, MongoDB, Parse
MongoDBv3.4 :                    	GlobalDocumentDB, MongoDB, Parse
mongoEnableDocLevelTTL :         	GlobalDocumentDB, MongoDB, Parse
DeleteAllItemsByPartitionKey :   	GlobalDocumentDB, MongoDB, Parse
DisableRateLimitingResponses :   	GlobalDocumentDB, MongoDB, Parse
AllowSelfServeUpgradeToMongo36 : 	GlobalDocumentDB, MongoDB, Parse
EnableMongoRetryableWrites :		MongoDB
EnableMongoRoleBasedAccessControl : MongoDB
EnableUniqueCompoundNestedDocs : 	MongoDB
EnableTtlOnCustomPath:              MongoDB
EnablePartialUniqueIndex:           MongoDB
*/
var capabilitiesToKindMap = map[string]interface{}{
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongo)):                       []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongo16MBDocumentSupport)):    []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl)): []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)):        []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs)):    []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableTtlOnCustomPath)):             []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnablePartialUniqueIndex)):          []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableNoSqlVectorSearch)):           []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableNoSqlFullTextSearch)):         []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableCassandra)):                   []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableGremlin)):                     []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableTable)):                       []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableServerless)):                  []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableAggregationPipeline)):         []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesMongoDBv34)):                        []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesMongoEnableDocLevelTTL)):            []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesDeleteAllItemsByPartitionKey)):      []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesDisableRateLimitingResponses)):      []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36)):    []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
}

// If the consistency policy of the Cosmos DB Database Account is not bounded staleness,
// any changes to the configuration for bounded staleness should be suppressed.
func suppressConsistencyPolicyStalenessConfiguration(_, _, _ string, d *pluginsdk.ResourceData) bool {
	consistencyPolicyList := d.Get("consistency_policy").([]interface{})
	if len(consistencyPolicyList) == 0 || consistencyPolicyList[0] == nil {
		return false
	}

	consistencyPolicy := consistencyPolicyList[0].(map[string]interface{})

	return consistencyPolicy["consistency_level"].(string) != string(cosmosdb.DefaultConsistencyLevelBoundedStaleness)
}

func resourceCosmosDbAccount() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceCosmosDbAccountCreate,
		Read:   resourceCosmosDbAccountRead,
		Update: resourceCosmosDbAccountUpdate,
		Delete: resourceCosmosDbAccountDelete,
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("backup.0.type", func(ctx context.Context, old, new, _ interface{}) bool {
				// backup type can only change from Periodic to Continuous
				return old.(string) == string(cosmosdb.BackupPolicyTypeContinuous) && new.(string) == string(cosmosdb.BackupPolicyTypePeriodic)
			}),

			pluginsdk.ForceNewIfChange("analytical_storage_enabled", func(ctx context.Context, old, new, _ interface{}) bool {
				// analytical_storage_enabled can not be changed after being set to true
				return old.(bool) && !new.(bool)
			}),

			pluginsdk.ForceNewIf("capabilities", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				kind := d.Get("kind").(string)
				old, new := d.GetChange("capabilities")

				return !checkCapabilitiesCanBeUpdated(kind, prepareCapabilities(old), prepareCapabilities(new))
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
				caps := diff.Get("capabilities")
				mongo34found := false
				enableMongo := false
				isMongo := strings.EqualFold(diff.Get("kind").(string), string(cosmosdb.DatabaseAccountKindMongoDB))

				for _, cap := range caps.(*pluginsdk.Set).List() {
					m := cap.(map[string]interface{})
					if v, ok := m["name"].(string); ok {
						if v == "MongoDBv3.4" {
							mongo34found = true
						} else if v == "EnableMongo" {
							enableMongo = true
						}
					}
				}

				if isMongo && (mongo34found && !enableMongo) {
					return fmt.Errorf("capability EnableMongo must be enabled if MongoDBv3.4 is also enabled")
				}
				return nil
			}),
		),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatabaseAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(180 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(180 * time.Minute),
		},

		SchemaVersion: 1,

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CosmosDBAccountV0toV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cosmosdb.DatabaseAccountOfferTypeStandard),
				}, false),
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
							ValidateFunc: validation.StringInSlice([]string{
								string(cosmosdb.AnalyticalStorageSchemaTypeWellDefined),
								string(cosmosdb.AnalyticalStorageSchemaTypeFullFidelity),
							}, false),
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cosmosdb.CreateModeDefault),
					string(cosmosdb.CreateModeRestore),
				}, false),
			},

			// Per Documentation: "The default identity needs to be explicitly set by the users." This should not be optional without a default anymore.
			// DOC: https://learn.microsoft.com/en-us/java/api/com.azure.resourcemanager.cosmos.models.databaseaccountupdateparameters?view=azure-java-stable#method-details
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(cosmosdb.DatabaseAccountKindGlobalDocumentDB),
				ValidateFunc: validation.StringInSlice([]string{
					string(cosmosdb.DatabaseAccountKindGlobalDocumentDB),
					string(cosmosdb.DatabaseAccountKindMongoDB),
					string(cosmosdb.DatabaseAccountKindParse),
				}, false),
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
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: keyVaultSuppress.DiffSuppressIgnoreKeyVaultKeyVersion,
				ValidateFunc:     keyVaultValidate.VersionlessNestedItemId,
				ConflictsWith:    []string{"managed_hsm_key_id"},
			},

			"managed_hsm_key_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  managedHsmValidate.ManagedHSMDataPlaneVersionlessKeyID,
				ConflictsWith: []string{"key_vault_key_id"},
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
							ValidateFunc: validation.StringInSlice([]string{
								string(cosmosdb.DefaultConsistencyLevelBoundedStaleness),
								string(cosmosdb.DefaultConsistencyLevelConsistentPrefix),
								string(cosmosdb.DefaultConsistencyLevelEventual),
								string(cosmosdb.DefaultConsistencyLevelSession),
								string(cosmosdb.DefaultConsistencyLevelStrong),
							}, false),
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
				Computed: true,
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

			"local_authentication_disabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"mongo_server_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
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
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cosmosdb.BackupPolicyTypeContinuous),
								string(cosmosdb.BackupPolicyTypePeriodic),
							}, false),
						},

						// Though `tier` has the default value `Continuous30Days` but `tier` is only for the backup type `Continuous`. So the default value isn't added in the property schema.
						"tier": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForContinuousTier(), false),
						},

						"interval_in_minutes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(60, 1440),
						},

						"retention_in_hours": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(8, 720),
						},

						"storage_redundancy": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cosmosdb.BackupStorageRedundancyGeo),
								string(cosmosdb.BackupStorageRedundancyLocal),
								string(cosmosdb.BackupStorageRedundancyZone),
							}, false),
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
							ValidateFunc: validate.RestorableDatabaseAccountID,
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
		},
	}

	if !features.FivePointOh() {
		resource.Schema["minimal_tls_version"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(cosmosdb.MinimalTlsVersionTlsOneTwo),
			ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForMinimalTlsVersion(), false),
		}
	}

	return resource
}

func resourceCosmosDbAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	databaseClient := meta.(*clients.Client).Cosmos.DatabaseClient
	accountClient := meta.(*clients.Client).Account
	subscriptionId := accountClient.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] Preparing arguments for AzureRM Cosmos DB Account creation")

	id := cosmosdb.NewDatabaseAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.DatabaseAccountsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_account", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	kind := d.Get("kind").(string)
	offerType := d.Get("offer_type").(string)

	ipRangeFilter := common.CosmosDBIpRangeFilterToIpRules(*utils.ExpandStringSlice(d.Get("ip_range_filter").(*pluginsdk.Set).List()))
	isVirtualNetworkFilterEnabled := d.Get("is_virtual_network_filter_enabled").(bool)

	enableFreeTier := d.Get("free_tier_enabled").(bool)
	enableAutomaticFailover := d.Get("automatic_failover_enabled").(bool)
	enableMultipleWriteLocations := d.Get("multiple_write_locations_enabled").(bool)

	partitionMergeEnabled := d.Get("partition_merge_enabled").(bool)
	burstCapacityEnabled := d.Get("burst_capacity_enabled").(bool)
	enableAnalyticalStorage := d.Get("analytical_storage_enabled").(bool)
	disableLocalAuthentication := d.Get("local_authentication_disabled").(bool)

	r, err := databaseClient.CheckNameExists(ctx, id.DatabaseAccountName)
	if err != nil {
		// TODO: remove when https://github.com/Azure/azure-sdk-for-go/issues/9891 is fixed
		if !utils.ResponseWasStatusCode(r, http.StatusInternalServerError) {
			return fmt.Errorf("checking if CosmosDB Account %s: %+v", id, err)
		}
	} else {
		if !utils.ResponseWasNotFound(r) {
			return fmt.Errorf("CosmosDB Account %s already exists, please import the resource via terraform import", id.DatabaseAccountName)
		}
	}
	geoLocations, err := expandAzureRmCosmosDBAccountGeoLocations(d)
	if err != nil {
		return fmt.Errorf("expanding %s geo locations: %+v", id, err)
	}

	publicNetworkAccess := cosmosdb.PublicNetworkAccessEnabled
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		publicNetworkAccess = cosmosdb.PublicNetworkAccessDisabled
	}

	networkByPass := cosmosdb.NetworkAclBypassNone
	if d.Get("network_acl_bypass_for_azure_services").(bool) {
		networkByPass = cosmosdb.NetworkAclBypassAzureServices
	}

	expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	capabilities := expandAzureRmCosmosDBAccountCapabilities(d)

	account := cosmosdb.DatabaseAccountCreateUpdateParameters{
		Location: pointer.To(location),
		Kind:     pointer.To(cosmosdb.DatabaseAccountKind(kind)),
		Identity: expandedIdentity,
		Properties: cosmosdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:           cosmosdb.DatabaseAccountOfferType(offerType),
			IPRules:                            ipRangeFilter,
			IsVirtualNetworkFilterEnabled:      utils.Bool(isVirtualNetworkFilterEnabled),
			EnableFreeTier:                     utils.Bool(enableFreeTier),
			EnableAutomaticFailover:            utils.Bool(enableAutomaticFailover),
			ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                          geoLocations,
			Capabilities:                       capabilities,
			MinimalTlsVersion:                  pointer.To(cosmosdb.MinimalTlsVersion(d.Get("minimal_tls_version").(string))),
			VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:       utils.Bool(enableMultipleWriteLocations),
			EnablePartitionMerge:               pointer.To(partitionMergeEnabled),
			EnableBurstCapacity:                pointer.To(burstCapacityEnabled),
			PublicNetworkAccess:                pointer.To(publicNetworkAccess),
			EnableAnalyticalStorage:            utils.Bool(enableAnalyticalStorage),
			Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]interface{})),
			DisableKeyBasedMetadataWriteAccess: utils.Bool(!d.Get("access_key_metadata_writes_enabled").(bool)),
			NetworkAclBypass:                   pointer.To(networkByPass),
			NetworkAclBypassResourceIds:        utils.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]interface{})),
			DisableLocalAuth:                   utils.Bool(disableLocalAuthentication),
		},
		Tags: tags.Expand(t),
	}

	// These values may not have changed but they need to be in the update params...
	if v, ok := d.GetOk("default_identity_type"); ok {
		account.Properties.DefaultIdentity = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("analytical_storage"); ok {
		account.Properties.AnalyticalStorageConfiguration = expandCosmosDBAccountAnalyticalStorageConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("capacity"); ok {
		account.Properties.Capacity = expandCosmosDBAccountCapacity(v.([]interface{}))
	}

	var createMode string
	if v, ok := d.GetOk("create_mode"); ok {
		createMode = v.(string)
		account.Properties.CreateMode = pointer.To(cosmosdb.CreateMode(createMode))
	}

	if v, ok := d.GetOk("restore"); ok {
		account.Properties.RestoreParameters = expandCosmosdbAccountRestoreParameters(v.([]interface{}))
	}

	if v, ok := d.GetOk("mongo_server_version"); ok {
		account.Properties.ApiProperties = &cosmosdb.ApiProperties{
			ServerVersion: pointer.To(cosmosdb.ServerVersion(v.(string))),
		}
	}

	if v, ok := d.GetOk("backup"); ok {
		policy, err := expandCosmosdbAccountBackup(v.([]interface{}), false, createMode)
		if err != nil {
			return fmt.Errorf("expanding `backup`: %+v", err)
		}
		account.Properties.BackupPolicy = policy
	} else if createMode != "" {
		return fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
	}

	if key, err := customermanagedkeys.ExpandKeyVaultOrManagedHSMKey(d, customermanagedkeys.VersionTypeAny, accountClient.Environment.KeyVault, accountClient.Environment.ManagedHSM); err != nil {
		return fmt.Errorf("parse key vault key id: %+v", err)
	} else if key != nil {
		account.Properties.KeyVaultKeyUri = pointer.To(key.ID())
	}

	// additional validation on MaxStalenessPrefix as it varies depending on if the DB is multi region or not
	consistencyPolicy := account.Properties.ConsistencyPolicy
	if len(geoLocations) > 1 && consistencyPolicy != nil && consistencyPolicy.DefaultConsistencyLevel == cosmosdb.DefaultConsistencyLevelBoundedStaleness {
		if msp := consistencyPolicy.MaxStalenessPrefix; msp != nil && pointer.From(msp) < 100000 {
			return fmt.Errorf("max_staleness_prefix (%d) must be greater then 100000 when more then one geo_location is used", *msp)
		}
		if mis := consistencyPolicy.MaxIntervalInSeconds; mis != nil && pointer.From(mis) < 300 {
			return fmt.Errorf("max_interval_in_seconds (%d) must be greater then 300 (5min) when more then one geo_location is used", *mis)
		}
	}

	err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id, account, d)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// NOTE: this is to work around the issue here: https://github.com/Azure/azure-rest-api-specs/issues/27596
	// Once the above issue is resolved we shouldn't need this check and update anymore
	if d.Get("create_mode").(string) == string(cosmosdb.CreateModeRestore) {
		err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id, account, d)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceCosmosDbAccountRead(d, meta)
}

func resourceCosmosDbAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	apiEnvs := meta.(*clients.Client).Account.Environment
	// subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] Preparing arguments for AzureRM Cosmos DB Account update")

	id, err := cosmosdb.ParseDatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	// get existing locations (if exists)
	existing, err := client.DatabaseAccountsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: properties were nil", id)
	}

	configLocations, err := expandAzureRmCosmosDBAccountGeoLocations(d)
	if err != nil {
		return fmt.Errorf("expanding %s geo locations: %+v", id, err)
	}

	// Normalize Locations...
	cosmosLocations := make([]cosmosdb.Location, 0)
	cosmosLocationsMap := map[string]cosmosdb.Location{}

	if existing.Model.Properties.Locations != nil {
		for _, l := range *existing.Model.Properties.Locations {
			location := cosmosdb.Location{
				Id:               l.Id,
				LocationName:     l.LocationName,
				FailoverPriority: l.FailoverPriority,
				IsZoneRedundant:  l.IsZoneRedundant,
			}

			cosmosLocations = append(cosmosLocations, location)
			cosmosLocationsMap[azure.NormalizeLocation(*location.LocationName)] = location
		}
	}

	var capabilities *[]cosmosdb.Capability
	if existing.Model.Properties.Capabilities != nil {
		capabilities = existing.Model.Properties.Capabilities
	}

	// backup must be updated independently
	var backup cosmosdb.BackupPolicy
	if existing.Model.Properties.BackupPolicy != nil {
		backup = existing.Model.Properties.BackupPolicy
		if d.HasChange("backup") {
			if v, ok := d.GetOk("backup"); ok {
				newBackup, err := expandCosmosdbAccountBackup(v.([]interface{}), d.HasChange("backup.0.type"), string(pointer.From(existing.Model.Properties.CreateMode)))
				if err != nil {
					return fmt.Errorf("expanding `backup`: %+v", err)
				}
				updateParameters := cosmosdb.DatabaseAccountUpdateParameters{
					Properties: &cosmosdb.DatabaseAccountUpdateProperties{
						BackupPolicy: newBackup,
					},
				}

				// Update Database 'backup'...
				if err := client.DatabaseAccountsUpdateThenPoll(ctx, *id, updateParameters); err != nil {
					return fmt.Errorf("updating CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
				}

				backup = newBackup
			} else if string(pointer.From(existing.Model.Properties.CreateMode)) != "" {
				return fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
			}
		}
	}

	updateRequired := false
	if props := existing.Model.Properties; props != nil {
		location := location.Normalize(pointer.From(existing.Model.Location))
		offerType := d.Get("offer_type").(string)
		t := tags.Expand(d.Get("tags").(map[string]interface{}))
		kind := cosmosdb.DatabaseAccountKind(d.Get("kind").(string))
		isVirtualNetworkFilterEnabled := pointer.To(d.Get("is_virtual_network_filter_enabled").(bool))
		enableAnalyticalStorage := pointer.To(d.Get("analytical_storage_enabled").(bool))
		disableLocalAuthentication := pointer.To(d.Get("local_authentication_disabled").(bool))
		enableAutomaticFailover := pointer.To(d.Get("automatic_failover_enabled").(bool))

		networkByPass := cosmosdb.NetworkAclBypassNone
		if d.Get("network_acl_bypass_for_azure_services").(bool) {
			networkByPass = cosmosdb.NetworkAclBypassAzureServices
		}

		ipRangeFilter := common.CosmosDBIpRangeFilterToIpRules(*utils.ExpandStringSlice(d.Get("ip_range_filter").(*pluginsdk.Set).List()))

		publicNetworkAccess := cosmosdb.PublicNetworkAccessEnabled
		if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
			publicNetworkAccess = cosmosdb.PublicNetworkAccessDisabled
		}

		// NOTE: these fields are expanded directly into the
		// 'DatabaseAccountCreateUpdateParameters' below or
		// are included in the 'DatabaseAccountCreateUpdateParameters'
		// later, however we need to know if they changed or not...
		// TODO Post 4.0 remove `enable_automatic_failover` from this list
		if d.HasChanges("consistency_policy", "virtual_network_rule", "cors_rule", "access_key_metadata_writes_enabled",
			"network_acl_bypass_for_azure_services", "network_acl_bypass_ids", "analytical_storage",
			"capacity", "create_mode", "restore", "key_vault_key_id", "managed_hsm_key_id", "mongo_server_version",
			"public_network_access_enabled", "ip_range_filter", "offer_type", "is_virtual_network_filter_enabled",
			"kind", "tags", "enable_automatic_failover", "automatic_failover_enabled", "analytical_storage_enabled",
			"local_authentication_disabled", "partition_merge_enabled", "minimal_tls_version", "burst_capacity_enabled") {
			updateRequired = true
		}

		// Incident : #383341730
		// Azure Bug: #2209567 'Updating identities and default identity at the same time fails silently'
		//
		// The 'Identity' field should only ever be sent once to the endpoint, except for updates and removal. If the
		// 'Identity' field is included in the update call with the 'DefaultIdentity' it will silently fail
		// per the bug noted above (e.g. Azure Bug #2209567).
		//
		// In the update scenario where the end-user would like to update their 'Identity' and their 'DefaultIdentity'
		// fields at the same time both of these operations need to happen atomically in separate PUT/PATCH calls
		// to the service else you will hit the bug mentioned above. You need to update the 'Identity' field
		// first then update the 'DefaultIdentity' in totally different PUT/PATCH calls where you have to drop
		// the 'Identity' field on the floor when updating the 'DefaultIdentity' field.
		//
		// NOTE      : If the 'Identity' field has not changed in the resource, do not send it in the payload.
		//             this workaround can be removed once the service team fixes the above mentioned bug.
		//
		// ADDITIONAL: You cannot update properties and add/remove replication locations or update the enabling of
		//             multiple write locations at the same time. So you must update any changed properties
		//             first, then address the replication locations and/or updating/enabling of
		//             multiple write locations.

		account := cosmosdb.DatabaseAccountCreateUpdateParameters{
			Location: pointer.To(location),
			Kind:     pointer.To(kind),
			Properties: cosmosdb.DatabaseAccountCreateUpdateProperties{
				DatabaseAccountOfferType:           cosmosdb.DatabaseAccountOfferType(offerType),
				IPRules:                            ipRangeFilter,
				IsVirtualNetworkFilterEnabled:      isVirtualNetworkFilterEnabled,
				EnableFreeTier:                     existing.Model.Properties.EnableFreeTier,
				EnableAutomaticFailover:            enableAutomaticFailover,
				MinimalTlsVersion:                  pointer.To(cosmosdb.MinimalTlsVersion(d.Get("minimal_tls_version").(string))),
				Capabilities:                       capabilities,
				ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
				Locations:                          cosmosLocations,
				VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
				EnableMultipleWriteLocations:       props.EnableMultipleWriteLocations,
				PublicNetworkAccess:                pointer.To(publicNetworkAccess),
				EnableAnalyticalStorage:            enableAnalyticalStorage,
				Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]interface{})),
				DisableKeyBasedMetadataWriteAccess: pointer.To(!d.Get("access_key_metadata_writes_enabled").(bool)),
				NetworkAclBypass:                   pointer.To(networkByPass),
				NetworkAclBypassResourceIds:        utils.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]interface{})),
				DisableLocalAuth:                   disableLocalAuthentication,
				BackupPolicy:                       backup,
				EnablePartitionMerge:               pointer.To(d.Get("partition_merge_enabled").(bool)),
				EnableBurstCapacity:                pointer.To(d.Get("burst_capacity_enabled").(bool)),
			},
			Tags: t,
		}

		if key, err := customermanagedkeys.ExpandKeyVaultOrManagedHSMKey(d, customermanagedkeys.VersionTypeAny, apiEnvs.KeyVault, apiEnvs.ManagedHSM); err != nil {
			return err
		} else if key != nil {
			account.Properties.KeyVaultKeyUri = pointer.To(key.ID())
		}

		// 'default_identity_type' will always have a value since it now has a default value of "FirstPartyIdentity" per the API documentation.
		// I do not include 'DefaultIdentity' and 'Identity' in the 'accountProps' intentionally, these operations need to be
		// performed mutually exclusive from each other in an atomic fashion, else you will hit the service teams bug...
		updateDefaultIdentity := false
		if d.HasChange("default_identity_type") {
			updateDefaultIdentity = true
		}

		// adding 'DefaultIdentity' to avoid causing it to fallback
		// to "FirstPartyIdentity" on update(s), issue #22466
		if v, ok := d.GetOk("default_identity_type"); ok {
			account.Properties.DefaultIdentity = pointer.To(v.(string))
		}

		// we need the following in the accountProps even if they have not changed...
		if v, ok := d.GetOk("analytical_storage"); ok {
			account.Properties.AnalyticalStorageConfiguration = expandCosmosDBAccountAnalyticalStorageConfiguration(v.([]interface{}))
		}

		if v, ok := d.GetOk("capacity"); ok {
			account.Properties.Capacity = expandCosmosDBAccountCapacity(v.([]interface{}))
		}

		var createMode string
		if v, ok := d.GetOk("create_mode"); ok {
			createMode = v.(string)
			account.Properties.CreateMode = pointer.To(cosmosdb.CreateMode(createMode))
		}

		if v, ok := d.GetOk("restore"); ok {
			account.Properties.RestoreParameters = expandCosmosdbAccountRestoreParameters(v.([]interface{}))
		}

		if !pluginsdk.IsExplicitlyNullInConfig(d, "mongo_server_version") {
			account.Properties.ApiProperties = &cosmosdb.ApiProperties{
				ServerVersion: pointer.To(cosmosdb.ServerVersion(d.Get("mongo_server_version").(string))),
			}
		}

		// Only do this update if a value has changed above...
		if updateRequired {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'DatabaseAccountCreateUpdateParameters'")

			// Update the database...
			if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, *id, account, d); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Update 'DatabaseAccountCreateUpdateParameters' [NO CHANGE]")
		}

		// Update the following properties independently after the initial CreateOrUpdate...
		if d.HasChange("multiple_write_locations_enabled") {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'EnableMultipleWriteLocations'")

			enableMultipleWriteLocations := pointer.To(d.Get("multiple_write_locations_enabled").(bool))
			if props.EnableMultipleWriteLocations != enableMultipleWriteLocations {
				account.Properties.EnableMultipleWriteLocations = enableMultipleWriteLocations

				// Update the database...
				if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, *id, account, d); err != nil {
					return fmt.Errorf("updating %q EnableMultipleWriteLocations: %+v", id, err)
				}
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Updating 'EnableMultipleWriteLocations' [NO CHANGE]")
		}

		// determine if any locations have been renamed/priority reordered and remove them
		updateLocations := false
		for _, configLoc := range configLocations {
			if cosmosLoc, ok := cosmosLocationsMap[pointer.From(configLoc.LocationName)]; ok {
				// is the location in the config also in the database with the same 'FailoverPriority'?
				if pointer.From(configLoc.FailoverPriority) != pointer.From(cosmosLoc.FailoverPriority) {
					// The Failover Priority has been changed in the config...
					if pointer.From(configLoc.FailoverPriority) == 0 {
						return fmt.Errorf("cannot change the failover priority of %q location %q to %d", id, pointer.From(configLoc.LocationName), pointer.From(configLoc.FailoverPriority))
					}

					// since the Locations FailoverPriority changed remove it from the map because
					// we have to update the Location in the database. The Locations
					// left in the map after this loop are the Locations that are
					// the same in the database and in the config file...
					delete(cosmosLocationsMap, pointer.From(configLoc.LocationName))
					updateLocations = true
				}
			}
		}

		if updateLocations {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Removing renamed 'Locations'")
			locationsUnchanged := make([]cosmosdb.Location, 0, len(cosmosLocationsMap))
			for _, value := range cosmosLocationsMap {
				locationsUnchanged = append(locationsUnchanged, value)
			}

			account.Properties.Locations = locationsUnchanged

			// Update the database...
			if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, *id, account, d); err != nil {
				return fmt.Errorf("removing %q renamed `locations`: %+v", id, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Removing renamed 'Locations' [NO CHANGE]")
		}

		if d.HasChanges("geo_location") {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'Locations'")
			// add any new/renamed locations
			account.Properties.Locations = configLocations

			// Update the database locations...
			err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, *id, account, d)
			if err != nil {
				return fmt.Errorf("updating %q `locations`: %+v", id, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Updating 'Locations' [NO CHANGE]")
		}

		// Update Identity and Default Identity...
		identityChanged := false
		expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		if d.HasChange("identity") {
			identityChanged = true

			// Looks like you have to always remove all the identities first before you can
			// reassign/modify them, else it will append any new/changed identities
			// resulting in a diff...
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Setting 'Identity' to 'None'")

			// can't set this back to account, because that will hit the bug...
			identityVal := cosmosdb.DatabaseAccountUpdateParameters{
				Identity: pointer.To(identity.LegacySystemAndUserAssignedMap{
					Type: identity.TypeNone,
				}),
			}

			// Update the database 'Identity' to 'None'...
			err = resourceCosmosDbAccountApiUpdate(client, ctx, *id, identityVal, d)
			if err != nil {
				return fmt.Errorf("updating 'identity' %q: %+v", id, err)
			}

			// If the Identity was removed from the configuration file it will be set as type None
			// so we can skip setting the Identity if it is going to be set to None...
			if expandedIdentity.Type != identity.TypeNone {
				log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'Identity' to %q", expandedIdentity.Type)

				identityVal := cosmosdb.DatabaseAccountUpdateParameters{
					Identity: expandedIdentity,
				}

				// Update the database...
				err = resourceCosmosDbAccountApiUpdate(client, ctx, *id, identityVal, d)
				if err != nil {
					return fmt.Errorf("updating 'identity' %q: %+v", id, err)
				}
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Updating 'Identity' [NO CHANGE]")
		}

		// NOTE: updateDefaultIdentity now has a default value of 'FirstPartyIdentity'... This value now gets
		//       triggered if the default value does not match the value in Azure...
		//
		// NOTE: When you change the 'Identity', the 'DefaultIdentity' will be set to 'undefined', so if you change
		//       the identity you must also update the 'DefaultIdentity' as well...
		if updateDefaultIdentity || identityChanged {
			// This will now return the default of 'FirstPartyIdentity' if it
			// is not set in the config, which is correct.
			configDefaultIdentity := d.Get("default_identity_type").(string)
			if identityChanged {
				log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'DefaultIdentity' to %q because the 'Identity' was changed to %q", configDefaultIdentity, expandedIdentity.Type)
			} else {
				log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'DefaultIdentity' to %q because 'default_identity_type' was changed", configDefaultIdentity)
			}

			// PATCH instead of PUT...
			defaultIdentity := cosmosdb.DatabaseAccountUpdateParameters{
				Properties: &cosmosdb.DatabaseAccountUpdateProperties{
					DefaultIdentity: pointer.To(configDefaultIdentity),
				},
			}

			// Update the database...
			err = resourceCosmosDbAccountApiUpdate(client, ctx, *id, defaultIdentity, d)
			if err != nil {
				return fmt.Errorf("updating 'default_identity_type' %q: %+v", id, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Updating 'DefaultIdentity' [NO CHANGE]")
		}
	}

	if existing.Model.Properties.Capabilities != nil {
		if d.HasChange("capabilities") {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'Capabilities'")

			newCapabilities := expandAzureRmCosmosDBAccountCapabilities(d)
			updateParameters := cosmosdb.DatabaseAccountUpdateParameters{
				Properties: &cosmosdb.DatabaseAccountUpdateProperties{
					Capabilities: newCapabilities,
				},
			}

			// Update Database 'capabilities'...
			if err := client.DatabaseAccountsUpdateThenPoll(ctx, *id, updateParameters); err != nil {
				return fmt.Errorf("updating CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Updating 'Capabilities' [NO CHANGE]")
		}
	}

	return resourceCosmosDbAccountRead(d, meta)
}

func resourceCosmosDbAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseDatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.DatabaseAccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
	}

	d.Set("name", id.DatabaseAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(existing.Model.Location))
	d.Set("kind", pointer.From(existing.Model.Kind))

	identity, err := identity.FlattenLegacySystemAndUserAssignedMap(existing.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := existing.Model.Properties; props != nil {
		d.Set("offer_type", pointer.From(props.DatabaseAccountOfferType))

		d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilter(props.IPRules))

		d.Set("endpoint", props.DocumentEndpoint)

		d.Set("free_tier_enabled", props.EnableFreeTier)
		d.Set("analytical_storage_enabled", props.EnableAnalyticalStorage)
		d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == cosmosdb.PublicNetworkAccessEnabled)
		if props.DefaultIdentity == nil || *props.DefaultIdentity != "" {
			d.Set("default_identity_type", props.DefaultIdentity)
		} else {
			d.Set("default_identity_type", "FirstPartyIdentity")
		}
		d.Set("minimal_tls_version", pointer.From(props.MinimalTlsVersion))
		d.Set("create_mode", pointer.From(props.CreateMode))
		d.Set("partition_merge_enabled", pointer.From(props.EnablePartitionMerge))
		d.Set("burst_capacity_enabled", pointer.From(props.EnableBurstCapacity))

		if v := props.IsVirtualNetworkFilterEnabled; v != nil {
			d.Set("is_virtual_network_filter_enabled", props.IsVirtualNetworkFilterEnabled)
		}

		if v := props.EnableAutomaticFailover; v != nil {
			d.Set("automatic_failover_enabled", props.EnableAutomaticFailover)
		}

		if v := props.KeyVaultKeyUri; v != nil {
			envs := meta.(*clients.Client).Account.Environment
			if key, err := customermanagedkeys.FlattenKeyVaultOrManagedHSMID(*v, envs.ManagedHSM); err != nil {
				return fmt.Errorf("flatten key vault uri: %+v", err)
			} else if key.IsSet() {
				if key.KeyVaultKeyId != nil {
					d.Set("key_vault_key_id", key.KeyVaultKeyId.ID())
				} else {
					d.Set("managed_hsm_key_id", key.ManagedHSMKeyID())
				}
			}
		}

		if v := existing.Model.Properties.EnableMultipleWriteLocations; v != nil {
			d.Set("multiple_write_locations_enabled", props.EnableMultipleWriteLocations)
		}

		if err := d.Set("analytical_storage", flattenCosmosDBAccountAnalyticalStorageConfiguration(props.AnalyticalStorageConfiguration)); err != nil {
			return fmt.Errorf("setting `analytical_storage`: %+v", err)
		}

		if err := d.Set("capacity", flattenCosmosDBAccountCapacity(props.Capacity)); err != nil {
			return fmt.Errorf("setting `capacity`: %+v", err)
		}

		if err := d.Set("restore", flattenCosmosdbAccountRestoreParameters(props.RestoreParameters)); err != nil {
			return fmt.Errorf("setting `restore`: %+v", err)
		}

		if err = d.Set("consistency_policy", flattenAzureRmCosmosDBAccountConsistencyPolicy(props.ConsistencyPolicy)); err != nil {
			return fmt.Errorf("setting CosmosDB Account %q `consistency_policy` (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
		}

		if err = d.Set("geo_location", flattenAzureRmCosmosDBAccountGeoLocations(props)); err != nil {
			return fmt.Errorf("setting `geo_location`: %+v", err)
		}

		if err = d.Set("capabilities", flattenAzureRmCosmosDBAccountCapabilities(props.Capabilities)); err != nil {
			return fmt.Errorf("setting `capabilities`: %+v", err)
		}

		if err = d.Set("virtual_network_rule", flattenAzureRmCosmosDBAccountVirtualNetworkRules(props.VirtualNetworkRules)); err != nil {
			return fmt.Errorf("setting `virtual_network_rule`: %+v", err)
		}

		d.Set("access_key_metadata_writes_enabled", !*props.DisableKeyBasedMetadataWriteAccess)
		if apiProps := props.ApiProperties; apiProps != nil {
			d.Set("mongo_server_version", pointer.From(apiProps.ServerVersion))
		}
		d.Set("network_acl_bypass_for_azure_services", pointer.From(props.NetworkAclBypass) == cosmosdb.NetworkAclBypassAzureServices)
		d.Set("network_acl_bypass_ids", utils.FlattenStringSlice(props.NetworkAclBypassResourceIds))

		if v := existing.Model.Properties.DisableLocalAuth; v != nil {
			d.Set("local_authentication_disabled", props.DisableLocalAuth)
		}

		policy, err := flattenCosmosdbAccountBackup(props.BackupPolicy)
		if err != nil {
			return err
		}

		if err = d.Set("backup", policy); err != nil {
			return fmt.Errorf("setting `backup`: %+v", err)
		}

		d.Set("cors_rule", common.FlattenCosmosCorsRule(props.Cors))
	}

	readEndpoints := make([]string, 0)
	if p := existing.Model.Properties.ReadLocations; p != nil {
		for _, l := range *p {
			if l.DocumentEndpoint == nil {
				continue
			}

			readEndpoints = append(readEndpoints, *l.DocumentEndpoint)
		}
	}
	if err := d.Set("read_endpoints", readEndpoints); err != nil {
		return fmt.Errorf("setting `read_endpoints`: %s", err)
	}

	writeEndpoints := make([]string, 0)
	if p := existing.Model.Properties.WriteLocations; p != nil {
		for _, l := range *p {
			if l.DocumentEndpoint == nil {
				continue
			}

			writeEndpoints = append(writeEndpoints, *l.DocumentEndpoint)
		}
	}
	if err := d.Set("write_endpoints", writeEndpoints); err != nil {
		return fmt.Errorf("setting `write_endpoints`: %s", err)
	}

	// ListKeys returns a data structure containing a DatabaseAccountListReadOnlyKeysResult pointer
	// implying that it also returns the read only keys, however this appears to not be the case
	keys, err := client.DatabaseAccountsListKeys(ctx, *id)
	if err != nil {
		if response.WasNotFound(keys.HttpResponse) {
			log.Printf("[DEBUG] Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.DatabaseAccountName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List Write keys for CosmosDB Account %s: %s", id.DatabaseAccountName, err)
	}
	d.Set("primary_key", keys.Model.PrimaryMasterKey)
	d.Set("secondary_key", keys.Model.SecondaryMasterKey)

	readonlyKeys, err := client.DatabaseAccountsListReadOnlyKeys(ctx, *id)
	if err != nil {
		if response.WasNotFound(keys.HttpResponse) {
			log.Printf("[DEBUG] Read Only Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.DatabaseAccountName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List read-only keys for CosmosDB Account %s: %s", id.DatabaseAccountName, err)
	}
	d.Set("primary_readonly_key", readonlyKeys.Model.PrimaryReadonlyMasterKey)
	d.Set("secondary_readonly_key", readonlyKeys.Model.SecondaryReadonlyMasterKey)

	connStringResp, err := client.DatabaseAccountsListConnectionStrings(ctx, *id)
	if err != nil {
		if response.WasNotFound(keys.HttpResponse) {
			log.Printf("[DEBUG] Connection Strings were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.ResourceGroupName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List connection strings for CosmosDB Account %s: %s", id.DatabaseAccountName, err)
	}

	var connStrings []string
	if connStringResp.Model.ConnectionStrings != nil {
		connStrings = make([]string, len(*connStringResp.Model.ConnectionStrings))
		for i, v := range *connStringResp.Model.ConnectionStrings {
			connStrings[i] = *v.ConnectionString
			if propertyName, propertyExists := connStringPropertyMap[*v.Description]; propertyExists {
				d.Set(propertyName, v.ConnectionString) // lintignore:R001
			}
		}
	}

	return tags.FlattenAndSet(d, existing.Model.Tags)
}

func resourceCosmosDbAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseDatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DatabaseAccountsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
	}

	// the SDK now will return a `WasNotFound` response even when still deleting
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"NotFound"},
		MinTimeout: 30 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.DatabaseAccountsGet(ctx, *id)
			if err2 != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return resp, "NotFound", nil
				}
				return nil, "", err2
			}

			return resp, "Deleting", nil
		},
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for CosmosDB Account %q (Resource Group %q) to be deleted: %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
	}

	return nil
}

func resourceCosmosDbAccountApiUpdate(client *cosmosdb.CosmosDBClient, ctx context.Context, id cosmosdb.DatabaseAccountId, account cosmosdb.DatabaseAccountUpdateParameters, d *pluginsdk.ResourceData) error {
	if err := client.DatabaseAccountsUpdateThenPoll(ctx, id, account); err != nil {
		return fmt.Errorf("updating CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Creating", "Updating", "Deleting", "Initializing", "Dequeued", "Enqueued"},
		Target:                    []string{"Succeeded"},
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 2,
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.DatabaseAccountsGet(ctx, id)
			if err2 != nil || resp.HttpResponse == nil || resp.HttpResponse.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("reading CosmosDB Account %q after update (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err2)
			}
			status := "Succeeded"

			return resp, status, nil
		},
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to update: %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
	}

	return nil
}

func resourceCosmosDbAccountApiCreateOrUpdate(client *cosmosdb.CosmosDBClient, ctx context.Context, id cosmosdb.DatabaseAccountId, account cosmosdb.DatabaseAccountCreateUpdateParameters, d *pluginsdk.ResourceData) error {
	if err := client.DatabaseAccountsCreateOrUpdateThenPoll(ctx, id, account); err != nil {
		return fmt.Errorf("creating/updating CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
	}

	// if a replication location is added or removed it can take some time to provision
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating", "Updating", "Deleting", "Initializing", "Dequeued", "Enqueued"},
		Target:     []string{"Succeeded"},
		MinTimeout: 15 * time.Second,
		Delay:      30 * time.Second, // required because it takes some time before the 'creating' location shows up
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.DatabaseAccountsGet(ctx, id)
			if err2 != nil || resp.HttpResponse == nil || resp.HttpResponse.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("reading CosmosDB Account %q after create/update (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err2)
			}
			status := "Succeeded"
			if props := resp.Model.Properties; props != nil {
				var locations []cosmosdb.Location

				if props.ReadLocations != nil {
					locations = append(locations, *props.ReadLocations...)
				}
				if props.WriteLocations != nil {
					locations = append(locations, *props.WriteLocations...)
				}

				for _, l := range locations {
					if status = *l.ProvisioningState; status == "Creating" || status == "Updating" || status == "Deleting" {
						break // return the first non successful status.
					}
				}

				for _, desiredLocation := range account.Properties.Locations {
					for index, l := range locations {
						if azure.NormalizeLocation(*desiredLocation.LocationName) == azure.NormalizeLocation(*l.LocationName) {
							break
						}

						if (index + 1) == len(locations) {
							return resp, "Updating", nil
						}
					}
				}
			}

			return resp, status, nil
		},
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to provision: %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
	}

	return nil
}

func expandAzureRmCosmosDBAccountConsistencyPolicy(d *pluginsdk.ResourceData) *cosmosdb.ConsistencyPolicy {
	i := d.Get("consistency_policy").([]interface{})
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	consistencyLevel := input["consistency_level"].(string)
	policy := cosmosdb.ConsistencyPolicy{
		DefaultConsistencyLevel: cosmosdb.DefaultConsistencyLevel(consistencyLevel),
	}

	if stalenessPrefix, ok := input["max_staleness_prefix"].(int); ok {
		if stalenessPrefix == 0 {
			stalenessPrefix = 100
		}
		policy.MaxStalenessPrefix = pointer.FromInt64(int64(stalenessPrefix))
	}
	if maxInterval, ok := input["max_interval_in_seconds"].(int); ok {
		if maxInterval == 0 {
			maxInterval = 5
		}
		policy.MaxIntervalInSeconds = utils.Int64(int64(maxInterval))
	}

	return &policy
}

func expandAzureRmCosmosDBAccountGeoLocations(d *pluginsdk.ResourceData) ([]cosmosdb.Location, error) {
	locations := make([]cosmosdb.Location, 0)
	for _, l := range d.Get("geo_location").(*pluginsdk.Set).List() {
		data := l.(map[string]interface{})

		location := cosmosdb.Location{
			LocationName:     pointer.To(azure.NormalizeLocation(data["location"].(string))),
			FailoverPriority: utils.Int64(int64(data["failover_priority"].(int))),
			IsZoneRedundant:  pointer.FromBool(data["zone_redundant"].(bool)),
		}

		locations = append(locations, location)
	}

	// TODO: maybe this should be in a CustomizeDiff
	// all priorities & locations must be unique
	byPriorities := make(map[int]interface{}, len(locations))
	byName := make(map[string]interface{}, len(locations))
	locationsCount := len(locations)
	for _, location := range locations {
		priority := int(*location.FailoverPriority)
		name := *location.LocationName

		if _, ok := byPriorities[priority]; ok {
			return nil, fmt.Errorf("each `geo_location` needs to have a unique failover_prioroty. Multiple instances of '%d' found", priority)
		}

		if _, ok := byName[name]; ok {
			return nil, fmt.Errorf("each `geo_location` needs to be in unique location. Multiple instances of %q found", name)
		}

		if priority > locationsCount-1 {
			return nil, fmt.Errorf("the maximum value for a failover priority = (total number of regions - 1). '%d' was found", priority)
		}

		byPriorities[priority] = location
		byName[name] = location
	}

	// and one of them must have a priority of 0...
	if _, ok := byPriorities[0]; !ok {
		return nil, fmt.Errorf("there needs to be a `geo_location` with a `failover_priority` of `0`")
	}

	return locations, nil
}

func expandAzureRmCosmosDBAccountCapabilities(d *pluginsdk.ResourceData) *[]cosmosdb.Capability {
	capabilities := d.Get("capabilities").(*pluginsdk.Set).List()
	s := make([]cosmosdb.Capability, 0)

	for _, c := range capabilities {
		m := c.(map[string]interface{})
		s = append(s, cosmosdb.Capability{Name: pointer.To(m["name"].(string))})
	}

	return &s
}

func expandAzureRmCosmosDBAccountVirtualNetworkRules(d *pluginsdk.ResourceData) *[]cosmosdb.VirtualNetworkRule {
	virtualNetworkRules := d.Get("virtual_network_rule").(*pluginsdk.Set).List()

	s := make([]cosmosdb.VirtualNetworkRule, len(virtualNetworkRules))
	for i, r := range virtualNetworkRules {
		m := r.(map[string]interface{})
		s[i] = cosmosdb.VirtualNetworkRule{
			Id:                               pointer.To(m["id"].(string)),
			IgnoreMissingVNetServiceEndpoint: pointer.FromBool(m["ignore_missing_vnet_service_endpoint"].(bool)),
		}
	}

	return &s
}

func flattenAzureRmCosmosDBAccountConsistencyPolicy(policy *cosmosdb.ConsistencyPolicy) []interface{} {
	result := map[string]interface{}{}
	result["consistency_level"] = string(policy.DefaultConsistencyLevel)
	if policy.MaxIntervalInSeconds != nil {
		result["max_interval_in_seconds"] = int(*policy.MaxIntervalInSeconds)
	}
	if policy.MaxStalenessPrefix != nil {
		result["max_staleness_prefix"] = int(*policy.MaxStalenessPrefix)
	}

	return []interface{}{result}
}

func flattenAzureRmCosmosDBAccountGeoLocations(account *cosmosdb.DatabaseAccountGetProperties) *pluginsdk.Set {
	locationSet := pluginsdk.Set{
		F: resourceAzureRMCosmosDBAccountGeoLocationHash,
	}
	if account == nil || account.FailoverPolicies == nil {
		return &locationSet
	}

	for _, l := range *account.FailoverPolicies {
		if l.Id == nil {
			continue
		}

		id := *l.Id
		lb := map[string]interface{}{
			"id":                id,
			"location":          location.NormalizeNilable(l.LocationName),
			"failover_priority": int(pointer.From(l.FailoverPriority)),
			// there is no zone redundancy information in FailoverPolicies currently, we have to search it by `id` in the Locations property.
			"zone_redundant": findZoneRedundant(account.Locations, id),
		}

		locationSet.Add(lb)
	}

	return &locationSet
}

func findZoneRedundant(locations *[]cosmosdb.Location, id string) bool {
	if locations == nil {
		return false
	}

	for _, location := range *locations {
		if location.Id != nil && *location.Id == id {
			if location.IsZoneRedundant != nil {
				return *location.IsZoneRedundant
			}
		}
	}

	return false
}

func isServerlessCapacityMode(accResp documentdb.DatabaseAccountGetResults) bool {
	if props := accResp.DatabaseAccountGetProperties; props != nil && props.Capabilities != nil {
		for _, v := range *props.Capabilities {
			if v.Name != nil && *v.Name == "EnableServerless" {
				return true
			}
		}
	}

	return false
}

func flattenAzureRmCosmosDBAccountCapabilities(capabilities *[]cosmosdb.Capability) *pluginsdk.Set {
	s := pluginsdk.Set{
		F: resourceAzureRMCosmosDBAccountCapabilitiesHash,
	}

	for _, c := range *capabilities {
		if v := c.Name; v != nil {
			e := map[string]interface{}{
				"name": *v,
			}
			s.Add(e)
		}
	}

	return &s
}

func flattenAzureRmCosmosDBAccountVirtualNetworkRules(rules *[]cosmosdb.VirtualNetworkRule) *pluginsdk.Set {
	results := pluginsdk.Set{
		F: resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash,
	}

	if rules != nil {
		for _, r := range *rules {
			rule := map[string]interface{}{
				"id":                                   *r.Id,
				"ignore_missing_vnet_service_endpoint": *r.IgnoreMissingVNetServiceEndpoint,
			}
			results.Add(rule)
		}
	}

	return &results
}

func resourceAzureRMCosmosDBAccountGeoLocationHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		location := location.Normalize(m["location"].(string))
		priority := int32(m["failover_priority"].(int))

		buf.WriteString(fmt.Sprintf("%s-%d", location, priority))
	}

	return pluginsdk.HashString(buf.String())
}

func resourceAzureRMCosmosDBAccountCapabilitiesHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}

func resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(strings.ToLower(m["id"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}

func expandCosmosdbAccountBackup(input []interface{}, backupHasChange bool, createMode string) (cosmosdb.BackupPolicy, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	attr := input[0].(map[string]interface{})

	switch attr["type"].(string) {
	case string(cosmosdb.BackupPolicyTypeContinuous):
		if v := attr["interval_in_minutes"].(int); v != 0 && !backupHasChange {
			return nil, fmt.Errorf("`interval_in_minutes` cannot be defined when the `backup.type` is set to %q", cosmosdb.BackupPolicyTypeContinuous)
		}

		if v := attr["retention_in_hours"].(int); v != 0 && !backupHasChange {
			return nil, fmt.Errorf("`retention_in_hours` cannot be defined when the `backup.type` is set to %q", cosmosdb.BackupPolicyTypeContinuous)
		}

		if v := attr["storage_redundancy"].(string); v != "" && !backupHasChange {
			return nil, fmt.Errorf("`storage_redundancy` cannot be defined when the `backup.type` is set to %q", cosmosdb.BackupPolicyTypeContinuous)
		}

		result := cosmosdb.ContinuousModeBackupPolicy{}

		if v := attr["tier"].(string); v != "" {
			result.ContinuousModeProperties = &cosmosdb.ContinuousModeProperties{
				Tier: pointer.To(cosmosdb.ContinuousTier(v)),
			}
		}

		return result, nil

	case string(cosmosdb.BackupPolicyTypePeriodic):
		if createMode != "" {
			return nil, fmt.Errorf("`create_mode` can only be defined when the `backup.type` is set to %q, got %q", cosmosdb.BackupPolicyTypeContinuous, cosmosdb.BackupPolicyTypePeriodic)
		}

		if v := attr["tier"].(string); v != "" && !backupHasChange {
			return nil, fmt.Errorf("`tier` can not be set when `type` in `backup` is `Periodic`")
		}

		// Mirror the behavior of the old SDK...
		periodicModeBackupPolicy := cosmosdb.PeriodicModeBackupPolicy{
			PeriodicModeProperties: &cosmosdb.PeriodicModeProperties{
				BackupIntervalInMinutes:        utils.Int64(int64(attr["interval_in_minutes"].(int))),
				BackupRetentionIntervalInHours: utils.Int64(int64(attr["retention_in_hours"].(int))),
			},
		}

		if v := attr["storage_redundancy"].(string); v != "" {
			periodicModeBackupPolicy.PeriodicModeProperties.BackupStorageRedundancy = pointer.To(cosmosdb.BackupStorageRedundancy(attr["storage_redundancy"].(string)))
		}

		return periodicModeBackupPolicy, nil

	default:
		return nil, fmt.Errorf("unknown `type` in `backup`:%+v", attr["type"].(string))
	}
}

func flattenCosmosdbAccountBackup(input cosmosdb.BackupPolicy) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	switch backupPolicy := input.(type) {
	case cosmosdb.ContinuousModeBackupPolicy:
		var tier cosmosdb.ContinuousTier
		if v := backupPolicy.ContinuousModeProperties; v != nil {
			tier = pointer.From(v.Tier)
		}
		return []interface{}{
			map[string]interface{}{
				"type": string(cosmosdb.BackupPolicyTypeContinuous),
				"tier": string(tier),
			},
		}, nil

	case cosmosdb.PeriodicModeBackupPolicy:
		var interval, retention int
		if v := backupPolicy.PeriodicModeProperties.BackupIntervalInMinutes; v != nil {
			interval = int(*v)
		}

		if v := backupPolicy.PeriodicModeProperties.BackupRetentionIntervalInHours; v != nil {
			retention = int(*v)
		}

		var storageRedundancy cosmosdb.BackupStorageRedundancy
		if backupPolicy.PeriodicModeProperties.BackupStorageRedundancy != nil {
			storageRedundancy = pointer.From(backupPolicy.PeriodicModeProperties.BackupStorageRedundancy)
		}

		return []interface{}{
			map[string]interface{}{
				"type":                string(cosmosdb.BackupPolicyTypePeriodic),
				"interval_in_minutes": interval,
				"retention_in_hours":  retention,
				"storage_redundancy":  storageRedundancy,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown `type` in `backup`: %+v", input)
	}
}

func expandCosmosDBAccountAnalyticalStorageConfiguration(input []interface{}) *cosmosdb.AnalyticalStorageConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &cosmosdb.AnalyticalStorageConfiguration{
		SchemaType: pointer.To(cosmosdb.AnalyticalStorageSchemaType(v["schema_type"].(string))),
	}
}

func expandCosmosDBAccountCapacity(input []interface{}) *cosmosdb.Capacity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &cosmosdb.Capacity{
		TotalThroughputLimit: utils.Int64(int64(v["total_throughput_limit"].(int))),
	}
}

func flattenCosmosDBAccountAnalyticalStorageConfiguration(input *cosmosdb.AnalyticalStorageConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var schemaType cosmosdb.AnalyticalStorageSchemaType
	if input.SchemaType != nil {
		schemaType = pointer.From(input.SchemaType)
	}

	return []interface{}{
		map[string]interface{}{
			"schema_type": schemaType,
		},
	}
}

func flattenCosmosDBAccountCapacity(input *cosmosdb.Capacity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var totalThroughputLimit int64
	if input.TotalThroughputLimit != nil {
		totalThroughputLimit = *input.TotalThroughputLimit
	}

	return []interface{}{
		map[string]interface{}{
			"total_throughput_limit": totalThroughputLimit,
		},
	}
}

func expandCosmosdbAccountRestoreParameters(input []interface{}) *cosmosdb.RestoreParameters {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	restoreParameters := cosmosdb.RestoreParameters{
		RestoreMode:               pointer.To(cosmosdb.RestoreModePointInTime),
		RestoreSource:             pointer.To(v["source_cosmosdb_account_id"].(string)),
		DatabasesToRestore:        expandCosmosdbAccountDatabasesToRestore(v["database"].(*pluginsdk.Set).List()),
		GremlinDatabasesToRestore: expandCosmosdbAccountGremlinDatabasesToRestore(v["gremlin_database"].([]interface{})),
	}

	restoreTimestampInUtc, _ := time.Parse(time.RFC3339, v["restore_timestamp_in_utc"].(string))
	restoreParameters.SetRestoreTimestampInUtcAsTime(restoreTimestampInUtc)

	if tablesToRestore := v["tables_to_restore"].([]interface{}); len(tablesToRestore) > 0 {
		restoreParameters.TablesToRestore = utils.ExpandStringSlice(tablesToRestore)
	}

	return &restoreParameters
}

func expandCosmosdbAccountDatabasesToRestore(input []interface{}) *[]cosmosdb.DatabaseRestoreResource {
	results := make([]cosmosdb.DatabaseRestoreResource, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, cosmosdb.DatabaseRestoreResource{
			DatabaseName:    pointer.To(v["name"].(string)),
			CollectionNames: utils.ExpandStringSlice(v["collection_names"].(*pluginsdk.Set).List()),
		})
	}
	return &results
}

func expandCosmosdbAccountGremlinDatabasesToRestore(input []interface{}) *[]cosmosdb.GremlinDatabaseRestoreResource {
	results := make([]cosmosdb.GremlinDatabaseRestoreResource, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, cosmosdb.GremlinDatabaseRestoreResource{
			DatabaseName: pointer.To(v["name"].(string)),
			GraphNames:   utils.ExpandStringSlice(v["graph_names"].([]interface{})),
		})
	}

	return &results
}

func flattenCosmosdbAccountRestoreParameters(input *cosmosdb.RestoreParameters) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	var restoreSource string
	if input.RestoreSource != nil {
		restoreSource = *input.RestoreSource
	}

	var restoreTimestampInUtc string
	if input.RestoreTimestampInUtc != nil {
		restoreTimestampInUtc = pointer.From(input.RestoreTimestampInUtc)
	}

	return []interface{}{
		map[string]interface{}{
			"database":                   flattenCosmosdbAccountDatabasesToRestore(input.DatabasesToRestore),
			"gremlin_database":           flattenCosmosdbAccountGremlinDatabasesToRestore(input.GremlinDatabasesToRestore),
			"source_cosmosdb_account_id": restoreSource,
			"restore_timestamp_in_utc":   restoreTimestampInUtc,
			"tables_to_restore":          pointer.From(input.TablesToRestore),
		},
	}
}

func flattenCosmosdbAccountDatabasesToRestore(input *[]cosmosdb.DatabaseRestoreResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var databaseName string
		if item.DatabaseName != nil {
			databaseName = *item.DatabaseName
		}

		results = append(results, map[string]interface{}{
			"collection_names": utils.FlattenStringSlice(item.CollectionNames),
			"name":             databaseName,
		})
	}
	return results
}

func flattenCosmosdbAccountGremlinDatabasesToRestore(input *[]cosmosdb.GremlinDatabaseRestoreResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"graph_names": utils.FlattenStringSlice(item.GraphNames),
			"name":        pointer.From(item.DatabaseName),
		})
	}

	return results
}

func checkCapabilitiesCanBeUpdated(kind string, oldCapabilities *[]cosmosdb.Capability, newCapabilities *[]cosmosdb.Capability) bool {
	// The feedback from service team : capabilities that can be added to an existing account
	canBeAddedCaps := []string{
		strings.ToLower(string(databaseAccountCapabilitiesDeleteAllItemsByPartitionKey)),
		strings.ToLower(string(databaseAccountCapabilitiesDisableRateLimitingResponses)),
		strings.ToLower(string(databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableAggregationPipeline)),
		strings.ToLower(string(databaseAccountCapabilitiesMongoDBv34)),
		strings.ToLower(string(databaseAccountCapabilitiesMongoEnableDocLevelTTL)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableMongo16MBDocumentSupport)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableTtlOnCustomPath)),
		strings.ToLower(string(databaseAccountCapabilitiesEnablePartialUniqueIndex)),
	}

	// The feedback from service team: capabilities that can be removed from an existing account
	canBeRemovedCaps := []string{
		strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)),
		strings.ToLower(string(databaseAccountCapabilitiesDisableRateLimitingResponses)),
	}

	// first check the new capabilities can be added
	for _, capability := range *newCapabilities {
		if capability.Name == nil {
			continue
		}
		existedPreviously := false
		for _, existing := range *oldCapabilities {
			if existing.Name != nil && strings.EqualFold(*existing.Name, *capability.Name) {
				existedPreviously = true
				break
			}
		}
		if existedPreviously {
			continue
		}

		// retrieve a list of capabilities for this DB type
		supportedKindsForCapability, ok := capabilitiesToKindMap[strings.ToLower(*capability.Name)]
		if !ok {
			return false
		}

		// first check if this is supported
		if isSupported := utils.SliceContainsValue(supportedKindsForCapability.([]string), strings.ToLower(kind)); !isSupported {
			return false
		}

		// then check if it can be added via an update
		if !utils.SliceContainsValue(canBeAddedCaps, strings.ToLower(*capability.Name)) {
			return false
		}
	}
	// then check if we're removing any that they can be removed
	for _, capability := range *oldCapabilities {
		existsNow := false
		for _, new := range *newCapabilities {
			if new.Name != nil && strings.EqualFold(*new.Name, *capability.Name) {
				existsNow = true
				break
			}
		}
		if existsNow {
			continue
		}

		if !utils.SliceContainsValue(canBeRemovedCaps, strings.ToLower(*capability.Name)) {
			return false
		}
	}

	return true
}

func prepareCapabilities(capabilities interface{}) *[]cosmosdb.Capability {
	output := make([]cosmosdb.Capability, 0)
	for _, v := range capabilities.(*pluginsdk.Set).List() {
		m := v.(map[string]interface{})
		if c, ok := m["name"].(string); ok {
			capability := cosmosdb.Capability{
				Name: pointer.To(c),
			}
			output = append(output, capability)
		}
	}
	return &output
}
