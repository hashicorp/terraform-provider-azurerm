// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultSuppress "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/suppress"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var CosmosDbAccountResourceName = "azurerm_cosmosdb_account"

var connStringPropertyMap = map[string]string{
	"Primary SQL Connection String":             "primary_sql_connection_string",
	"Secondary SQL Connection String":           "secondary_sql_connection_string",
	"Primary Read-Only SQL Connection String":   "primary_readonly_sql_connection_string",
	"Secondary Read-Only SQL Connection String": "secondary_readonly_sql_connection_string",
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
	databaseAccountCapabilitiesDisableRateLimitingResponses      databaseAccountCapabilities = "DisableRateLimitingResponses"
	databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36    databaseAccountCapabilities = "AllowSelfServeUpgradeToMongo36"
	databaseAccountCapabilitiesEnableMongoRetryableWrites        databaseAccountCapabilities = "EnableMongoRetryableWrites"
	databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl databaseAccountCapabilities = "EnableMongoRoleBasedAccessControl"
	databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs    databaseAccountCapabilities = "EnableUniqueCompoundNestedDocs"
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
DisableRateLimitingResponses :   	GlobalDocumentDB, MongoDB, Parse
AllowSelfServeUpgradeToMongo36 : 	GlobalDocumentDB, MongoDB, Parse
EnableMongoRetryableWrites :		MongoDB
EnableMongoRoleBasedAccessControl : MongoDB
EnableUniqueCompoundNestedDocs : 	MongoDB
EnableTtlOnCustomPath:              MongoDB
EnablePartialUniqueIndex:           MongoDB
*/
var capabilitiesToKindMap = map[string]interface{}{
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongo)):                       []string{strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongo16MBDocumentSupport)):    []string{strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl)): []string{strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)):        []string{strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs)):    []string{strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableTtlOnCustomPath)):             []string{strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnablePartialUniqueIndex)):          []string{strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableCassandra)):                   []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableGremlin)):                     []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableTable)):                       []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableServerless)):                  []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableAggregationPipeline)):         []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesMongoDBv34)):                        []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesMongoEnableDocLevelTTL)):            []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesDisableRateLimitingResponses)):      []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36)):    []string{strings.ToLower(string(documentdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(documentdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(documentdb.DatabaseAccountKindParse))},
}

// If the consistency policy of the Cosmos DB Database Account is not bounded staleness,
// any changes to the configuration for bounded staleness should be suppressed.
func suppressConsistencyPolicyStalenessConfiguration(_, _, _ string, d *pluginsdk.ResourceData) bool {
	consistencyPolicyList := d.Get("consistency_policy").([]interface{})
	if len(consistencyPolicyList) == 0 || consistencyPolicyList[0] == nil {
		return false
	}

	consistencyPolicy := consistencyPolicyList[0].(map[string]interface{})

	return consistencyPolicy["consistency_level"].(string) != string(documentdb.DefaultConsistencyLevelBoundedStaleness)
}

func resourceCosmosDbAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbAccountCreate,
		Read:   resourceCosmosDbAccountRead,
		Update: resourceCosmosDbAccountUpdate,
		Delete: resourceCosmosDbAccountDelete,
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("backup.0.type", func(ctx context.Context, old, new, _ interface{}) bool {
				// backup type can only change from Periodic to Continuous
				return old.(string) == string(documentdb.TypeContinuous) && new.(string) == string(documentdb.TypePeriodic)
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
				isMongo := strings.EqualFold(diff.Get("kind").(string), string(documentdb.DatabaseAccountKindMongoDB))

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
					string(documentdb.DatabaseAccountOfferTypeStandard),
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
								string(documentdb.AnalyticalStorageSchemaTypeWellDefined),
								string(documentdb.AnalyticalStorageSchemaTypeFullFidelity),
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

			"create_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.CreateModeDefault),
					string(documentdb.CreateModeRestore),
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
				Default:  string(documentdb.DatabaseAccountKindGlobalDocumentDB),
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.DatabaseAccountKindGlobalDocumentDB),
					string(documentdb.DatabaseAccountKindMongoDB),
					string(documentdb.DatabaseAccountKindParse),
				}, false),
			},

			"ip_range_filter": func() *schema.Schema {
				if features.FourPointOhBeta() {
					return &schema.Schema{
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDR,
						},
					}
				}
				return &schema.Schema{
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^(\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(/([1-2][0-9]|3[0-2]))?\b[,]?)*$`),
						"Cosmos DB ip_range_filter must be a set of CIDR IP addresses separated by commas with no spaces: '10.0.0.1,10.0.0.2,10.20.0.0/16'",
					),
				}
			}(),

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_free_tier": {
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

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_automatic_failover": {
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
								string(documentdb.DefaultConsistencyLevelBoundedStaleness),
								string(documentdb.DefaultConsistencyLevelConsistentPrefix),
								string(documentdb.DefaultConsistencyLevelEventual),
								string(documentdb.DefaultConsistencyLevelSession),
								string(documentdb.DefaultConsistencyLevelStrong),
							}, false),
						},

						"max_interval_in_seconds": {
							Type:             pluginsdk.TypeInt,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressConsistencyPolicyStalenessConfiguration,
							ValidateFunc:     validation.IntBetween(5, 86400), // single region values
						},

						"max_staleness_prefix": {
							Type:             pluginsdk.TypeInt,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressConsistencyPolicyStalenessConfiguration,
							ValidateFunc:     validation.IntBetween(10, 2147483647), // single region values
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
								string(databaseAccountCapabilitiesDisableRateLimitingResponses),
								string(databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36),
								string(databaseAccountCapabilitiesEnableMongoRetryableWrites),
								string(databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl),
								string(databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs),
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

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_multiple_write_locations": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
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
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.ServerVersionThreeFullStopTwo),
					string(documentdb.ServerVersionThreeFullStopSix),
					string(documentdb.ServerVersionFourFullStopZero),
					string(documentdb.ServerVersionFourFullStopTwo),
				}, false),
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
								string(documentdb.TypeContinuous),
								string(documentdb.TypePeriodic),
							}, false),
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
								string(documentdb.BackupStorageRedundancyGeo),
								string(documentdb.BackupStorageRedundancyLocal),
								string(documentdb.BackupStorageRedundancyZone),
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

			"connection_strings": {
				Type:      pluginsdk.TypeList,
				Computed:  true,
				Sensitive: true,
				Elem: &pluginsdk.Schema{
					Type:      pluginsdk.TypeString,
					Sensitive: true,
				},
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

			"tags": tags.Schema(),
		},
	}
}

func resourceCosmosDbAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] Preparing arguments for AzureRM Cosmos DB Account creation")

	id := parse.NewDatabaseAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_account", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	kind := d.Get("kind").(string)
	offerType := d.Get("offer_type").(string)

	var ipRangeFilter *[]documentdb.IPAddressOrRange
	if features.FourPointOhBeta() {
		ipRangeFilter = common.CosmosDBIpRangeFilterToIpRules(*utils.ExpandStringSlice(d.Get("ip_range_filter").(*pluginsdk.Set).List()))
	} else {
		ipRangeFilter = common.CosmosDBIpRangeFilterToIpRulesThreePointOh(d.Get("ip_range_filter").(string))
	}

	isVirtualNetworkFilterEnabled := d.Get("is_virtual_network_filter_enabled").(bool)
	enableFreeTier := d.Get("enable_free_tier").(bool)
	enableAutomaticFailover := d.Get("enable_automatic_failover").(bool)
	enableMultipleWriteLocations := d.Get("enable_multiple_write_locations").(bool)
	enableAnalyticalStorage := d.Get("analytical_storage_enabled").(bool)
	disableLocalAuthentication := d.Get("local_authentication_disabled").(bool)

	r, err := client.CheckNameExists(ctx, id.Name)
	if err != nil {
		// TODO: remove when https://github.com/Azure/azure-sdk-for-go/issues/9891 is fixed
		if !utils.ResponseWasStatusCode(r, http.StatusInternalServerError) {
			return fmt.Errorf("checking if CosmosDB Account %s: %+v", id, err)
		}
	} else {
		if !utils.ResponseWasNotFound(r) {
			return fmt.Errorf("CosmosDB Account %s already exists, please import the resource via terraform import", id.Name)
		}
	}
	geoLocations, err := expandAzureRmCosmosDBAccountGeoLocations(d)
	if err != nil {
		return fmt.Errorf("expanding %s geo locations: %+v", id, err)
	}

	publicNetworkAccess := documentdb.PublicNetworkAccessEnabled
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		publicNetworkAccess = documentdb.PublicNetworkAccessDisabled
	}

	networkByPass := documentdb.NetworkACLBypassNone
	if d.Get("network_acl_bypass_for_azure_services").(bool) {
		networkByPass = documentdb.NetworkACLBypassAzureServices
	}

	expandedIdentity, err := expandAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	capabilities := expandAzureRmCosmosDBAccountCapabilities(d)

	account := documentdb.DatabaseAccountCreateUpdateParameters{
		Location: pointer.To(location),
		Kind:     documentdb.DatabaseAccountKind(kind),
		Identity: expandedIdentity,
		DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:           pointer.To(offerType),
			IPRules:                            ipRangeFilter,
			IsVirtualNetworkFilterEnabled:      utils.Bool(isVirtualNetworkFilterEnabled),
			EnableFreeTier:                     utils.Bool(enableFreeTier),
			EnableAutomaticFailover:            utils.Bool(enableAutomaticFailover),
			ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                          &geoLocations,
			Capabilities:                       capabilities,
			VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:       utils.Bool(enableMultipleWriteLocations),
			PublicNetworkAccess:                publicNetworkAccess,
			EnableAnalyticalStorage:            utils.Bool(enableAnalyticalStorage),
			Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]interface{})),
			DisableKeyBasedMetadataWriteAccess: utils.Bool(!d.Get("access_key_metadata_writes_enabled").(bool)),
			NetworkACLBypass:                   networkByPass,
			NetworkACLBypassResourceIds:        utils.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]interface{})),
			DisableLocalAuth:                   utils.Bool(disableLocalAuthentication),
		},
		Tags: tags.Expand(t),
	}

	// These values may not have changed but they need to be in the update params...
	if v, ok := d.GetOk("default_identity_type"); ok {
		account.DatabaseAccountCreateUpdateProperties.DefaultIdentity = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("analytical_storage"); ok {
		account.DatabaseAccountCreateUpdateProperties.AnalyticalStorageConfiguration = expandCosmosDBAccountAnalyticalStorageConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("capacity"); ok {
		account.DatabaseAccountCreateUpdateProperties.Capacity = expandCosmosDBAccountCapacity(v.([]interface{}))
	}

	var createMode string
	if v, ok := d.GetOk("create_mode"); ok {
		createMode = v.(string)
		account.DatabaseAccountCreateUpdateProperties.CreateMode = documentdb.CreateMode(createMode)
	}

	if v, ok := d.GetOk("restore"); ok {
		account.DatabaseAccountCreateUpdateProperties.RestoreParameters = expandCosmosdbAccountRestoreParameters(v.([]interface{}))
	}

	if v, ok := d.GetOk("mongo_server_version"); ok {
		account.DatabaseAccountCreateUpdateProperties.APIProperties = &documentdb.APIProperties{
			ServerVersion: documentdb.ServerVersion(v.(string)),
		}
	}

	if v, ok := d.GetOk("backup"); ok {
		policy, err := expandCosmosdbAccountBackup(v.([]interface{}), false, createMode)
		if err != nil {
			return fmt.Errorf("expanding `backup`: %+v", err)
		}
		account.DatabaseAccountCreateUpdateProperties.BackupPolicy = policy
	} else if createMode != "" {
		return fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
	}

	if keyVaultKeyIDRaw, ok := d.GetOk("key_vault_key_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyIDRaw.(string))
		if err != nil {
			return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}
		account.DatabaseAccountCreateUpdateProperties.KeyVaultKeyURI = pointer.To(keyVaultKey.ID())
	}

	// additional validation on MaxStalenessPrefix as it varies depending on if the DB is multi region or not
	consistencyPolicy := account.DatabaseAccountCreateUpdateProperties.ConsistencyPolicy
	if len(geoLocations) > 1 && consistencyPolicy != nil && consistencyPolicy.DefaultConsistencyLevel == documentdb.DefaultConsistencyLevelBoundedStaleness {
		if msp := consistencyPolicy.MaxStalenessPrefix; msp != nil && pointer.From(msp) < 100000 {
			return fmt.Errorf("max_staleness_prefix (%d) must be greater then 100000 when more then one geo_location is used", *msp)
		}
		if mis := consistencyPolicy.MaxIntervalInSeconds; mis != nil && pointer.From(mis) < 300 {
			return fmt.Errorf("max_interval_in_seconds (%d) must be greater then 300 (5min) when more then one geo_location is used", *mis)
		}
	}

	err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id.ResourceGroup, id.Name, account, d)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbAccountRead(d, meta)
}

func resourceCosmosDbAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	// subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] Preparing arguments for AzureRM Cosmos DB Account update")

	id, err := parse.DatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	// get existing locations (if exists)
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	if existing.DatabaseAccountGetProperties == nil {
		return fmt.Errorf("retrieving %s: properties were nil", id)
	}

	configLocations, err := expandAzureRmCosmosDBAccountGeoLocations(d)
	if err != nil {
		return fmt.Errorf("expanding %s geo locations: %+v", id, err)
	}

	// Normalize Locations...
	cosmosLocations := make([]documentdb.Location, 0)
	cosmosLocationsMap := map[string]documentdb.Location{}

	if existing.DatabaseAccountGetProperties.Locations != nil {
		for _, l := range *existing.DatabaseAccountGetProperties.Locations {
			location := documentdb.Location{
				ID:               l.ID,
				LocationName:     l.LocationName,
				FailoverPriority: l.FailoverPriority,
				IsZoneRedundant:  l.IsZoneRedundant,
			}

			cosmosLocations = append(cosmosLocations, location)
			cosmosLocationsMap[azure.NormalizeLocation(*location.LocationName)] = location
		}
	}

	var capabilities *[]documentdb.Capability
	if existing.DatabaseAccountGetProperties.Capabilities != nil {
		capabilities = existing.DatabaseAccountGetProperties.Capabilities
		if d.HasChange("capabilities") {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'Capabilities'")

			newCapabilities := expandAzureRmCosmosDBAccountCapabilities(d)
			updateParameters := documentdb.DatabaseAccountUpdateParameters{
				DatabaseAccountUpdateProperties: &documentdb.DatabaseAccountUpdateProperties{
					Capabilities: newCapabilities,
				},
			}

			// Update Database 'capabilities'...
			future, err := client.Update(ctx, id.ResourceGroup, id.Name, updateParameters)
			if err != nil {
				return fmt.Errorf("updating CosmosDB Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to finish updating: %+v", id.Name, id.ResourceGroup, err)
			}

			capabilities = newCapabilities
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Updating 'Capabilities' [NO CHANGE]")
		}
	}

	// backup must be updated independently
	var backup documentdb.BasicBackupPolicy
	if existing.DatabaseAccountGetProperties.BackupPolicy != nil {
		backup = existing.DatabaseAccountGetProperties.BackupPolicy
		if d.HasChange("backup") {
			if v, ok := d.GetOk("backup"); ok {
				newBackup, err := expandCosmosdbAccountBackup(v.([]interface{}), d.HasChange("backup.0.type"), string(existing.DatabaseAccountGetProperties.CreateMode))
				if err != nil {
					return fmt.Errorf("expanding `backup`: %+v", err)
				}
				updateParameters := documentdb.DatabaseAccountUpdateParameters{
					DatabaseAccountUpdateProperties: &documentdb.DatabaseAccountUpdateProperties{
						BackupPolicy: newBackup,
					},
				}

				// Update Database 'backup'...
				future, err := client.Update(ctx, id.ResourceGroup, id.Name, updateParameters)
				if err != nil {
					return fmt.Errorf("updating CosmosDB Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
				}

				if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
					return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to finish updating: %+v", id.Name, id.ResourceGroup, err)
				}
				backup = newBackup
			} else if string(existing.CreateMode) != "" {
				return fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
			}
		}
	}

	updateRequired := false
	if props := existing.DatabaseAccountGetProperties; props != nil {
		location := azure.NormalizeLocation(pointer.From(existing.Location))
		offerType := d.Get("offer_type").(string)
		t := tags.Expand(d.Get("tags").(map[string]interface{}))
		kind := documentdb.DatabaseAccountKind(d.Get("kind").(string))
		isVirtualNetworkFilterEnabled := pointer.To(d.Get("is_virtual_network_filter_enabled").(bool))
		enableFreeTier := pointer.To(d.Get("enable_free_tier").(bool))
		enableAutomaticFailover := pointer.To(d.Get("enable_automatic_failover").(bool))
		enableAnalyticalStorage := pointer.To(d.Get("analytical_storage_enabled").(bool))
		disableLocalAuthentication := pointer.To(d.Get("local_authentication_disabled").(bool))

		networkByPass := documentdb.NetworkACLBypassNone
		if d.Get("network_acl_bypass_for_azure_services").(bool) {
			networkByPass = documentdb.NetworkACLBypassAzureServices
		}

		var ipRangeFilter *[]documentdb.IPAddressOrRange
		if features.FourPointOhBeta() {
			ipRangeFilter = common.CosmosDBIpRangeFilterToIpRules(*utils.ExpandStringSlice(d.Get("ip_range_filter").(*pluginsdk.Set).List()))
		} else {
			ipRangeFilter = common.CosmosDBIpRangeFilterToIpRulesThreePointOh(d.Get("ip_range_filter").(string))
		}

		publicNetworkAccess := documentdb.PublicNetworkAccessEnabled
		if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
			publicNetworkAccess = documentdb.PublicNetworkAccessDisabled
		}

		// NOTE: these fields are expanded directly into the
		// 'DatabaseAccountCreateUpdateParameters' below or
		// are included in the 'DatabaseAccountCreateUpdateParameters'
		// later, however we need to know if they changed or not...
		if d.HasChanges("consistency_policy", "virtual_network_rule", "cors_rule", "access_key_metadata_writes_enabled",
			"network_acl_bypass_for_azure_services", "network_acl_bypass_ids", "analytical_storage",
			"capacity", "create_mode", "restore", "key_vault_key_id", "mongo_server_version",
			"public_network_access_enabled", "ip_range_filter", "offer_type", "is_virtual_network_filter_enabled",
			"kind", "tags", "enable_free_tier", "enable_automatic_failover", "analytical_storage_enabled",
			"local_authentication_disabled") {
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

		account := documentdb.DatabaseAccountCreateUpdateParameters{
			Location: pointer.To(location),
			Kind:     kind,
			DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
				DatabaseAccountOfferType:           pointer.To(offerType),
				IPRules:                            ipRangeFilter,
				IsVirtualNetworkFilterEnabled:      isVirtualNetworkFilterEnabled,
				EnableFreeTier:                     enableFreeTier,
				EnableAutomaticFailover:            enableAutomaticFailover,
				Capabilities:                       capabilities,
				ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
				Locations:                          pointer.To(cosmosLocations),
				VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
				EnableMultipleWriteLocations:       props.EnableMultipleWriteLocations,
				PublicNetworkAccess:                publicNetworkAccess,
				EnableAnalyticalStorage:            enableAnalyticalStorage,
				Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]interface{})),
				DisableKeyBasedMetadataWriteAccess: pointer.To(!d.Get("access_key_metadata_writes_enabled").(bool)),
				NetworkACLBypass:                   networkByPass,
				NetworkACLBypassResourceIds:        utils.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]interface{})),
				DisableLocalAuth:                   disableLocalAuthentication,
				BackupPolicy:                       backup,
			},
			Tags: t,
		}

		accountProps := account.DatabaseAccountCreateUpdateProperties

		if keyVaultKeyIDRaw, ok := d.GetOk("key_vault_key_id"); ok {
			keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyIDRaw.(string))
			if err != nil {
				return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
			}
			accountProps.KeyVaultKeyURI = pointer.To(keyVaultKey.ID())
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
			accountProps.DefaultIdentity = pointer.To(v.(string))
		}

		// we need the following in the accountProps even if they have not changed...
		if v, ok := d.GetOk("analytical_storage"); ok {
			accountProps.AnalyticalStorageConfiguration = expandCosmosDBAccountAnalyticalStorageConfiguration(v.([]interface{}))
		}

		if v, ok := d.GetOk("capacity"); ok {
			accountProps.Capacity = expandCosmosDBAccountCapacity(v.([]interface{}))
		}

		var createMode string
		if v, ok := d.GetOk("create_mode"); ok {
			createMode = v.(string)
			accountProps.CreateMode = documentdb.CreateMode(createMode)
		}

		if v, ok := d.GetOk("restore"); ok {
			accountProps.RestoreParameters = expandCosmosdbAccountRestoreParameters(v.([]interface{}))
		}

		if v, ok := d.GetOk("mongo_server_version"); ok {
			accountProps.APIProperties = &documentdb.APIProperties{
				ServerVersion: documentdb.ServerVersion(v.(string)),
			}
		}

		// Only do this update if a value has changed above...
		if updateRequired {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'DatabaseAccountCreateUpdateParameters'")

			// Update the database...
			if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id.ResourceGroup, id.Name, account, d); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Update 'DatabaseAccountCreateUpdateParameters' [NO CHANGE]")
		}

		// Update the following properties independently after the initial CreateOrUpdate...
		if d.HasChange("enable_multiple_write_locations") {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'EnableMultipleWriteLocations'")

			enableMultipleWriteLocations := pointer.To(d.Get("enable_multiple_write_locations").(bool))
			if props.EnableMultipleWriteLocations != enableMultipleWriteLocations {
				accountProps.EnableMultipleWriteLocations = enableMultipleWriteLocations

				// Update the database...
				if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id.ResourceGroup, id.Name, account, d); err != nil {
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
			locationsUnchanged := make([]documentdb.Location, 0, len(cosmosLocationsMap))
			for _, value := range cosmosLocationsMap {
				locationsUnchanged = append(locationsUnchanged, value)
			}

			accountProps.Locations = pointer.To(locationsUnchanged)

			// Update the database...
			if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id.ResourceGroup, id.Name, account, d); err != nil {
				return fmt.Errorf("removing %q renamed `locations`: %+v", id, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Removing renamed 'Locations' [NO CHANGE]")
		}

		if d.HasChanges("geo_location") {
			log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'Locations'")
			// add any new/renamed locations
			accountProps.Locations = pointer.To(configLocations)

			// Update the database locations...
			err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id.ResourceGroup, id.Name, account, d)
			if err != nil {
				return fmt.Errorf("updating %q `locations`: %+v", id, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Updating 'Locations' [NO CHANGE]")
		}

		// Update Identity and Default Identity...
		identityChanged := false
		expandedIdentity, err := expandAccountIdentity(d.Get("identity").([]interface{}))
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
			identity := documentdb.DatabaseAccountUpdateParameters{
				Identity: pointer.To(documentdb.ManagedServiceIdentity{
					Type: documentdb.ResourceIdentityTypeNone,
				}),
			}

			// Update the database 'Identity' to 'None'...
			err = resourceCosmosDbAccountApiUpdate(client, ctx, id.ResourceGroup, id.Name, identity, d)
			if err != nil {
				return fmt.Errorf("updating 'identity' %q: %+v", id, err)
			}

			// If the Identity was removed from the configuration file it will be set as type None
			// so we can skip setting the Identity if it is going to be set to None...
			if expandedIdentity.Type != documentdb.ResourceIdentityTypeNone {
				log.Printf("[INFO] Updating AzureRM Cosmos DB Account: Updating 'Identity' to %q", expandedIdentity.Type)

				identity := documentdb.DatabaseAccountUpdateParameters{
					Identity: expandedIdentity,
				}

				// Update the database...
				err = resourceCosmosDbAccountApiUpdate(client, ctx, id.ResourceGroup, id.Name, identity, d)
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
			defaultIdentity := documentdb.DatabaseAccountUpdateParameters{
				DatabaseAccountUpdateProperties: &documentdb.DatabaseAccountUpdateProperties{
					DefaultIdentity: pointer.To(configDefaultIdentity),
				},
			}

			// Update the database...
			err = resourceCosmosDbAccountApiUpdate(client, ctx, id.ResourceGroup, id.Name, defaultIdentity, d)
			if err != nil {
				return fmt.Errorf("updating 'default_identity_type' %q: %+v", id, err)
			}
		} else {
			log.Printf("[INFO] [SKIP] AzureRM Cosmos DB Account: Updating 'DefaultIdentity' [NO CHANGE]")
		}
	}

	return resourceCosmosDbAccountRead(d, meta)
}

func resourceCosmosDbAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving CosmosDB Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(existing.Location))
	d.Set("kind", string(existing.Kind))

	identity, err := flattenAccountIdentity(existing.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := existing.DatabaseAccountGetProperties; props != nil {
		d.Set("offer_type", string(props.DatabaseAccountOfferType))

		if features.FourPointOhBeta() {
			d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilter(props.IPRules))
		} else {
			d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilterThreePointOh(props.IPRules))
		}

		d.Set("endpoint", props.DocumentEndpoint)
		d.Set("enable_free_tier", props.EnableFreeTier)
		d.Set("analytical_storage_enabled", props.EnableAnalyticalStorage)
		d.Set("public_network_access_enabled", props.PublicNetworkAccess == documentdb.PublicNetworkAccessEnabled)
		d.Set("default_identity_type", props.DefaultIdentity)
		d.Set("create_mode", props.CreateMode)

		if v := existing.IsVirtualNetworkFilterEnabled; v != nil {
			d.Set("is_virtual_network_filter_enabled", props.IsVirtualNetworkFilterEnabled)
		}

		if v := existing.EnableAutomaticFailover; v != nil {
			d.Set("enable_automatic_failover", props.EnableAutomaticFailover)
		}

		if v := existing.KeyVaultKeyURI; v != nil {
			d.Set("key_vault_key_id", props.KeyVaultKeyURI)
		}

		if v := existing.EnableMultipleWriteLocations; v != nil {
			d.Set("enable_multiple_write_locations", props.EnableMultipleWriteLocations)
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
			return fmt.Errorf("setting CosmosDB Account %q `consistency_policy` (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
		if apiProps := props.APIProperties; apiProps != nil {
			d.Set("mongo_server_version", apiProps.ServerVersion)
		}
		d.Set("network_acl_bypass_for_azure_services", props.NetworkACLBypass == documentdb.NetworkACLBypassAzureServices)
		d.Set("network_acl_bypass_ids", utils.FlattenStringSlice(props.NetworkACLBypassResourceIds))

		if v := existing.DisableLocalAuth; v != nil {
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
	if p := existing.ReadLocations; p != nil {
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
	if p := existing.WriteLocations; p != nil {
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
	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(keys.Response) {
			log.Printf("[DEBUG] Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List Write keys for CosmosDB Account %s: %s", id.Name, err)
	}
	d.Set("primary_key", keys.PrimaryMasterKey)
	d.Set("secondary_key", keys.SecondaryMasterKey)

	readonlyKeys, err := client.ListReadOnlyKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(keys.Response) {
			log.Printf("[DEBUG] Read Only Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List read-only keys for CosmosDB Account %s: %s", id.Name, err)
	}
	d.Set("primary_readonly_key", readonlyKeys.PrimaryReadonlyMasterKey)
	d.Set("secondary_readonly_key", readonlyKeys.SecondaryReadonlyMasterKey)

	connStringResp, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(keys.Response) {
			log.Printf("[DEBUG] Connection Strings were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List connection strings for CosmosDB Account %s: %s", id.Name, err)
	}

	var connStrings []string
	if connStringResp.ConnectionStrings != nil {
		connStrings = make([]string, len(*connStringResp.ConnectionStrings))
		for i, v := range *connStringResp.ConnectionStrings {
			connStrings[i] = *v.ConnectionString
			if propertyName, propertyExists := connStringPropertyMap[*v.Description]; propertyExists {
				d.Set(propertyName, v.ConnectionString) // lintignore:R001
			}
		}
	}

	d.Set("connection_strings", connStrings)

	return tags.FlattenAndSet(d, existing.Tags)
}

func resourceCosmosDbAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting CosmosDB Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	// the SDK now will return a `WasNotFound` response even when still deleting
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"NotFound"},
		MinTimeout: 30 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.Get(ctx, id.ResourceGroup, id.Name)
			if err2 != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return resp, "NotFound", nil
				}
				return nil, "", err2
			}

			return resp, "Deleting", nil
		},
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for CosmosDB Account %q (Resource Group %q) to be deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func resourceCosmosDbAccountApiUpdate(client *documentdb.DatabaseAccountsClient, ctx context.Context, resourceGroup string, name string, account documentdb.DatabaseAccountUpdateParameters, d *pluginsdk.ResourceData) error {
	future, err := client.Update(ctx, resourceGroup, name, account)
	if err != nil {
		return fmt.Errorf("updating CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to finish creating/updating: %+v", name, resourceGroup, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Creating", "Updating", "Deleting", "Initializing", "Dequeued", "Enqueued"},
		Target:                    []string{"Succeeded"},
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 2,
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.Get(ctx, resourceGroup, name)
			if err2 != nil || resp.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("reading CosmosDB Account %q after update (Resource Group %q): %+v", name, resourceGroup, err2)
			}
			status := "Succeeded"

			return resp, status, nil
		},
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to update: %+v", name, resourceGroup, err)
	}

	return nil
}

func resourceCosmosDbAccountApiCreateOrUpdate(client *documentdb.DatabaseAccountsClient, ctx context.Context, resourceGroup string, name string, account documentdb.DatabaseAccountCreateUpdateParameters, d *pluginsdk.ResourceData) error {
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, account)
	if err != nil {
		return fmt.Errorf("creating/updating CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to finish creating/updating: %+v", name, resourceGroup, err)
	}

	// if a replication location is added or removed it can take some time to provision
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating", "Updating", "Deleting", "Initializing", "Dequeued", "Enqueued"},
		Target:     []string{"Succeeded"},
		MinTimeout: 15 * time.Second,
		Delay:      30 * time.Second, // required because it takes some time before the 'creating' location shows up
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.Get(ctx, resourceGroup, name)
			if err2 != nil || resp.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("reading CosmosDB Account %q after create/update (Resource Group %q): %+v", name, resourceGroup, err2)
			}
			status := "Succeeded"
			if props := resp.DatabaseAccountGetProperties; props != nil {
				var locations []documentdb.Location

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

				for _, desiredLocation := range *account.Locations {
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

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to provision: %+v", name, resourceGroup, err)
	}

	return nil
}

func expandAzureRmCosmosDBAccountConsistencyPolicy(d *pluginsdk.ResourceData) *documentdb.ConsistencyPolicy {
	i := d.Get("consistency_policy").([]interface{})
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	consistencyLevel := input["consistency_level"].(string)
	policy := documentdb.ConsistencyPolicy{
		DefaultConsistencyLevel: documentdb.DefaultConsistencyLevel(consistencyLevel),
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
		policy.MaxIntervalInSeconds = utils.Int32(int32(maxInterval))
	}

	return &policy
}

func expandAzureRmCosmosDBAccountGeoLocations(d *pluginsdk.ResourceData) ([]documentdb.Location, error) {
	locations := make([]documentdb.Location, 0)
	for _, l := range d.Get("geo_location").(*pluginsdk.Set).List() {
		data := l.(map[string]interface{})

		location := documentdb.Location{
			LocationName:     pointer.To(azure.NormalizeLocation(data["location"].(string))),
			FailoverPriority: utils.Int32(int32(data["failover_priority"].(int))),
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

func expandAzureRmCosmosDBAccountCapabilities(d *pluginsdk.ResourceData) *[]documentdb.Capability {
	capabilities := d.Get("capabilities").(*pluginsdk.Set).List()
	s := make([]documentdb.Capability, 0)

	for _, c := range capabilities {
		m := c.(map[string]interface{})
		s = append(s, documentdb.Capability{Name: pointer.To(m["name"].(string))})
	}

	return &s
}

func expandAzureRmCosmosDBAccountVirtualNetworkRules(d *pluginsdk.ResourceData) *[]documentdb.VirtualNetworkRule {
	virtualNetworkRules := d.Get("virtual_network_rule").(*pluginsdk.Set).List()

	s := make([]documentdb.VirtualNetworkRule, len(virtualNetworkRules))
	for i, r := range virtualNetworkRules {
		m := r.(map[string]interface{})
		s[i] = documentdb.VirtualNetworkRule{
			ID:                               pointer.To(m["id"].(string)),
			IgnoreMissingVNetServiceEndpoint: pointer.FromBool(m["ignore_missing_vnet_service_endpoint"].(bool)),
		}
	}

	return &s
}

func flattenAzureRmCosmosDBAccountConsistencyPolicy(policy *documentdb.ConsistencyPolicy) []interface{} {
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

func flattenAzureRmCosmosDBAccountGeoLocations(account *documentdb.DatabaseAccountGetProperties) *pluginsdk.Set {
	locationSet := pluginsdk.Set{
		F: resourceAzureRMCosmosDBAccountGeoLocationHash,
	}
	if account == nil || account.FailoverPolicies == nil {
		return &locationSet
	}

	for _, l := range *account.FailoverPolicies {
		if l.ID == nil {
			continue
		}

		id := *l.ID
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

func findZoneRedundant(locations *[]documentdb.Location, id string) bool {
	if locations == nil {
		return false
	}

	for _, location := range *locations {
		if location.ID != nil && *location.ID == id {
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

func flattenAzureRmCosmosDBAccountCapabilities(capabilities *[]documentdb.Capability) *pluginsdk.Set {
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

func flattenAzureRmCosmosDBAccountVirtualNetworkRules(rules *[]documentdb.VirtualNetworkRule) *pluginsdk.Set {
	results := pluginsdk.Set{
		F: resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash,
	}

	if rules != nil {
		for _, r := range *rules {
			rule := map[string]interface{}{
				"id":                                   *r.ID,
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

func expandCosmosdbAccountBackup(input []interface{}, backupHasChange bool, createMode string) (documentdb.BasicBackupPolicy, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	attr := input[0].(map[string]interface{})

	switch attr["type"].(string) {
	case string(documentdb.TypeContinuous):
		if v := attr["interval_in_minutes"].(int); v != 0 && !backupHasChange {
			return nil, fmt.Errorf("`interval_in_minutes` can not be set when `type` in `backup` is `Continuous`")
		}
		if v := attr["retention_in_hours"].(int); v != 0 && !backupHasChange {
			return nil, fmt.Errorf("`retention_in_hours` can not be set when `type` in `backup` is `Continuous`")
		}
		if v := attr["storage_redundancy"].(string); v != "" && !backupHasChange {
			return nil, fmt.Errorf("`storage_redundancy` can not be set when `type` in `backup` is `Continuous`")
		}
		return documentdb.ContinuousModeBackupPolicy{
			Type: documentdb.TypeContinuous,
		}, nil

	case string(documentdb.TypePeriodic):
		if createMode != "" {
			return nil, fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
		}

		return documentdb.PeriodicModeBackupPolicy{
			Type: documentdb.TypePeriodic,
			PeriodicModeProperties: &documentdb.PeriodicModeProperties{
				BackupIntervalInMinutes:        utils.Int32(int32(attr["interval_in_minutes"].(int))),
				BackupRetentionIntervalInHours: utils.Int32(int32(attr["retention_in_hours"].(int))),
				BackupStorageRedundancy:        documentdb.BackupStorageRedundancy(attr["storage_redundancy"].(string)),
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown `type` in `backup`:%+v", attr["type"].(string))
	}
}

func flattenCosmosdbAccountBackup(input documentdb.BasicBackupPolicy) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	switch input.(type) {
	case documentdb.ContinuousModeBackupPolicy:
		return []interface{}{
			map[string]interface{}{
				"type": string(documentdb.TypeContinuous),
			},
		}, nil

	case documentdb.PeriodicModeBackupPolicy:
		policy, ok := input.AsPeriodicModeBackupPolicy()
		if !ok {
			return nil, fmt.Errorf("can not transit %+v into `backup` of `type` `Periodic`", input)
		}
		var interval, retention int
		if v := policy.PeriodicModeProperties.BackupIntervalInMinutes; v != nil {
			interval = int(*v)
		}
		if v := policy.PeriodicModeProperties.BackupRetentionIntervalInHours; v != nil {
			retention = int(*v)
		}
		var storageRedundancy documentdb.BackupStorageRedundancy
		if policy.PeriodicModeProperties.BackupStorageRedundancy != "" {
			storageRedundancy = policy.PeriodicModeProperties.BackupStorageRedundancy
		}
		return []interface{}{
			map[string]interface{}{
				"type":                string(documentdb.TypePeriodic),
				"interval_in_minutes": interval,
				"retention_in_hours":  retention,
				"storage_redundancy":  storageRedundancy,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown `type` in `backup`: %+v", input)
	}
}

func expandAccountIdentity(input []interface{}) (*documentdb.ManagedServiceIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := documentdb.ManagedServiceIdentity{
		Type: documentdb.ResourceIdentityType(string(expanded.Type)),
	}

	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*documentdb.ManagedServiceIdentityUserAssignedIdentitiesValue)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &documentdb.ManagedServiceIdentityUserAssignedIdentitiesValue{}
		}
	}

	return &out, nil
}

func flattenAccountIdentity(input *documentdb.ManagedServiceIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = pointer.To(identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		})

		if input.PrincipalID != nil {
			transform.PrincipalId = pointer.From(input.PrincipalID)
		}

		if input.TenantID != nil {
			transform.TenantId = pointer.From(input.TenantID)
		}

		if input.UserAssignedIdentities != nil {
			log.Printf("[DEBUG] input.UserAssignedIdentities ***NOT NULL***")

			for k, v := range input.UserAssignedIdentities {
				log.Printf("[DEBUG]  *** Parsing input.UserAssignedIdentities")
				details := identity.UserAssignedIdentityDetails{}

				if v.ClientID != nil {
					details.ClientId = v.ClientID
				}

				if v.PrincipalID != nil {
					details.PrincipalId = v.PrincipalID
				}

				transform.IdentityIds[k] = details

				log.Printf("[DEBUG]  *** Details: {ClientId: %q, PrincipalId: %q}", pointer.From(details.ClientId), pointer.From(details.PrincipalId))
				log.Printf("[DEBUG]  *** Current 'transform': %+v", transform)
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func expandCosmosDBAccountAnalyticalStorageConfiguration(input []interface{}) *documentdb.AnalyticalStorageConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &documentdb.AnalyticalStorageConfiguration{
		SchemaType: documentdb.AnalyticalStorageSchemaType(v["schema_type"].(string)),
	}
}

func expandCosmosDBAccountCapacity(input []interface{}) *documentdb.Capacity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &documentdb.Capacity{
		TotalThroughputLimit: utils.Int32(int32(v["total_throughput_limit"].(int))),
	}
}

func flattenCosmosDBAccountAnalyticalStorageConfiguration(input *documentdb.AnalyticalStorageConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var schemaType documentdb.AnalyticalStorageSchemaType
	if input.SchemaType != "" {
		schemaType = input.SchemaType
	}

	return []interface{}{
		map[string]interface{}{
			"schema_type": schemaType,
		},
	}
}

func flattenCosmosDBAccountCapacity(input *documentdb.Capacity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var totalThroughputLimit int32
	if input.TotalThroughputLimit != nil {
		totalThroughputLimit = *input.TotalThroughputLimit
	}

	return []interface{}{
		map[string]interface{}{
			"total_throughput_limit": totalThroughputLimit,
		},
	}
}

func expandCosmosdbAccountRestoreParameters(input []interface{}) *documentdb.RestoreParameters {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	restoreTimestampInUtc, _ := time.Parse(time.RFC3339, v["restore_timestamp_in_utc"].(string))

	return &documentdb.RestoreParameters{
		RestoreMode:           documentdb.RestoreModePointInTime,
		RestoreSource:         pointer.To(v["source_cosmosdb_account_id"].(string)),
		RestoreTimestampInUtc: &date.Time{Time: restoreTimestampInUtc},
		DatabasesToRestore:    expandCosmosdbAccountDatabasesToRestore(v["database"].(*pluginsdk.Set).List()),
	}
}

func expandCosmosdbAccountDatabasesToRestore(input []interface{}) *[]documentdb.DatabaseRestoreResource {
	results := make([]documentdb.DatabaseRestoreResource, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, documentdb.DatabaseRestoreResource{
			DatabaseName:    pointer.To(v["name"].(string)),
			CollectionNames: utils.ExpandStringSlice(v["collection_names"].(*pluginsdk.Set).List()),
		})
	}

	return &results
}

func flattenCosmosdbAccountRestoreParameters(input *documentdb.RestoreParameters) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var restoreSource string
	if input.RestoreSource != nil {
		restoreSource = *input.RestoreSource
	}

	var restoreTimestampInUtc string
	if input.RestoreTimestampInUtc != nil {
		restoreTimestampInUtc = input.RestoreTimestampInUtc.Format(time.RFC3339)
	}

	return []interface{}{
		map[string]interface{}{
			"database":                   flattenCosmosdbAccountDatabasesToRestore(input.DatabasesToRestore),
			"source_cosmosdb_account_id": restoreSource,
			"restore_timestamp_in_utc":   restoreTimestampInUtc,
		},
	}
}

func flattenCosmosdbAccountDatabasesToRestore(input *[]documentdb.DatabaseRestoreResource) []interface{} {
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

func checkCapabilitiesCanBeUpdated(kind string, oldCapabilities *[]documentdb.Capability, newCapabilities *[]documentdb.Capability) bool {
	// The feedback from service team : capabilities that can be added to an existing account
	canBeAddedCaps := []string{
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

func prepareCapabilities(capabilities interface{}) *[]documentdb.Capability {
	output := make([]documentdb.Capability, 0)
	for _, v := range capabilities.(*pluginsdk.Set).List() {
		m := v.(map[string]interface{})
		if c, ok := m["name"].(string); ok {
			cap := documentdb.Capability{
				Name: pointer.To(c),
			}
			output = append(output, cap)
		}
	}
	return &output
}
