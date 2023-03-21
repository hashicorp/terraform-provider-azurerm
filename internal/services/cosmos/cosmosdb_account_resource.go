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

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2021-10-15/cosmosdb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
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
*/
var capabilitiesToKindMap = map[string]interface{}{
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongo)):                    []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongo16MBDocumentSupport)): []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)):     []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)):     []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs)): []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableCassandra)):                []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableGremlin)):                  []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableTable)):                    []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableServerless)):               []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableAggregationPipeline)):      []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesMongoDBv34)):                     []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesMongoEnableDocLevelTTL)):         []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesDisableRateLimitingResponses)):   []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36)): []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
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
	return &pluginsdk.Resource{
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
			_, err := cosmosdb.ParseDatabaseAccountID(id)
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

			"default_identity_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.Any(
					validation.StringMatch(regexp.MustCompile(`^UserAssignedIdentity(.)+$`), "It may start with `UserAssignedIdentity`"),
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
								string(cosmosdb.DefaultConsistencyLevelBoundedStaleness),
								string(cosmosdb.DefaultConsistencyLevelConsistentPrefix),
								string(cosmosdb.DefaultConsistencyLevelEventual),
								string(cosmosdb.DefaultConsistencyLevelSession),
								string(cosmosdb.DefaultConsistencyLevelStrong),
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
					string(cosmosdb.ServerVersionThreePointTwo),
					string(cosmosdb.ServerVersionThreePointSix),
					string(cosmosdb.ServerVersionFourPointZero),
					string(cosmosdb.ServerVersionFourPointTwo),
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
								string(cosmosdb.BackupPolicyTypeContinuous),
								string(cosmosdb.BackupPolicyTypePeriodic),
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
	log.Printf("[INFO] preparing arguments for AzureRM Cosmos DB Account creation.")

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

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	kind := d.Get("kind").(string)
	offerType := d.Get("offer_type").(string)
	var ipRangeFilter *[]cosmosdb.IPAddressOrRange
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

	dbAccountNameId := cosmosdb.NewDatabaseAccountNameID(id.DatabaseAccountName)
	r, err := client.DatabaseAccountsCheckNameExists(ctx, dbAccountNameId)
	if err != nil {
		// todo remove when https://github.com/Azure/azure-sdk-for-go/issues/9891 is fixed
		if !response.WasStatusCode(r.HttpResponse, http.StatusInternalServerError) && !response.WasNotFound(r.HttpResponse) {
			return fmt.Errorf("checking if CosmosDB Account %s: %+v", id, err)
		}
	} else {
		if !response.WasNotFound(r.HttpResponse) {
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

	expandedIdentity, err := expandAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	capabilities := expandAzureRmCosmosDBAccountCapabilities(d)

	account := cosmosdb.DatabaseAccountCreateUpdateParameters{
		Location: utils.String(location),
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
			VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:       utils.Bool(enableMultipleWriteLocations),
			PublicNetworkAccess:                &publicNetworkAccess,
			EnableAnalyticalStorage:            utils.Bool(enableAnalyticalStorage),
			Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]interface{})),
			DisableKeyBasedMetadataWriteAccess: utils.Bool(!d.Get("access_key_metadata_writes_enabled").(bool)),
			NetworkAclBypass:                   &networkByPass,
			NetworkAclBypassResourceIds:        utils.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]interface{})),
			DisableLocalAuth:                   utils.Bool(disableLocalAuthentication),
		},
		Tags: expandTags(t),
	}

	if v, ok := d.GetOk("default_identity_type"); ok {
		account.Properties.DefaultIdentity = utils.String(v.(string))
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

	if keyVaultKeyIDRaw, ok := d.GetOk("key_vault_key_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyIDRaw.(string))
		if err != nil {
			return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}
		account.Properties.KeyVaultKeyUri = utils.String(keyVaultKey.ID())
	}

	// additional validation on MaxStalenessPrefix as it varies depending on if the DB is multi region or not
	consistencyPolicy := account.Properties.ConsistencyPolicy
	if len(geoLocations) > 1 && consistencyPolicy != nil && consistencyPolicy.DefaultConsistencyLevel == cosmosdb.DefaultConsistencyLevelBoundedStaleness {
		if msp := consistencyPolicy.MaxStalenessPrefix; msp != nil && *msp < 100000 {
			return fmt.Errorf("max_staleness_prefix (%d) must be greater then 100000 when more then one geo_location is used", *msp)
		}
		if mis := consistencyPolicy.MaxIntervalInSeconds; mis != nil && *mis < 300 {
			return fmt.Errorf("max_interval_in_seconds (%d) must be greater then 300 (5min) when more then one geo_location is used", *mis)
		}
	}

	err = resourceCosmosDbAccountApiUpsert(client, ctx, id, account, d)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbAccountRead(d, meta)
}

func resourceCosmosDbAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Cosmos DB Account update.")

	id := cosmosdb.NewDatabaseAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := d.Get("location").(string)
	t := d.Get("tags").(map[string]interface{})

	kind := d.Get("kind").(string)
	offerType := d.Get("offer_type").(string)
	var ipRangeFilter *[]cosmosdb.IPAddressOrRange
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

	newLocations, err := expandAzureRmCosmosDBAccountGeoLocations(d)
	if err != nil {
		return fmt.Errorf("expanding %s geo locations: %+v", id, err)
	}

	// get existing locations (if exists)
	resp, err := client.DatabaseAccountsGet(ctx, id)
	if err != nil {
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	capabilities := make([]cosmosdb.Capability, 0)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.Capabilities != nil {
				for _, v := range *props.Capabilities {
					c := cosmosdb.Capability{
						Name: v.Name,
					}
					capabilities = append(capabilities, c)
				}
			}
		}
	}

	if d.HasChange("capabilities") {

		newCapabilities := expandAzureRmCosmosDBAccountCapabilities(d)

		updateParameters := cosmosdb.DatabaseAccountUpdateParameters{
			Properties: &cosmosdb.DatabaseAccountUpdateProperties{
				Capabilities: newCapabilities,
			},
		}

		err = client.DatabaseAccountsUpdateThenPoll(ctx, id, updateParameters)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		resp, err := client.DatabaseAccountsGet(ctx, id)
		if err != nil {
			return fmt.Errorf("re-retrieving %s after update: %+v", id.ID(), err)
		}

		capabilities = make([]cosmosdb.Capability, 0)
		oldLocations := make([]cosmosdb.Location, 0)
		oldLocationsMap := map[string]cosmosdb.Location{}
		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				if props.Capabilities != nil {
					for _, v := range *props.Capabilities {
						c := cosmosdb.Capability{
							Name: v.Name,
						}
						capabilities = append(capabilities, c)
					}
				}

				if props.Locations != nil {
					for _, l := range *props.Locations {
						location := cosmosdb.Location{
							Id:               l.Id,
							LocationName:     l.LocationName,
							FailoverPriority: l.FailoverPriority,
							IsZoneRedundant:  l.IsZoneRedundant,
						}
						oldLocations = append(oldLocations, location)
						oldLocationsMap[azure.NormalizeLocation(*location.LocationName)] = location
					}
				}
			}
		}
	}

	oldLocations := make([]cosmosdb.Location, 0)
	oldLocationsMap := map[string]cosmosdb.Location{}
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.Locations != nil {
				for _, l := range *props.Locations {
					location := cosmosdb.Location{
						Id:               l.Id,
						LocationName:     l.LocationName,
						FailoverPriority: l.FailoverPriority,
						IsZoneRedundant:  l.IsZoneRedundant,
					}
					oldLocations = append(oldLocations, location)
					oldLocationsMap[azure.NormalizeLocation(*location.LocationName)] = location
				}
			}
		}
	}

	publicNetworkAccess := cosmosdb.PublicNetworkAccessEnabled
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		publicNetworkAccess = cosmosdb.PublicNetworkAccessDisabled
	}

	networkByPass := cosmosdb.NetworkAclBypassNone
	if d.Get("network_acl_bypass_for_azure_services").(bool) {
		networkByPass = cosmosdb.NetworkAclBypassAzureServices
	}

	expandedIdentity, err := expandAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	// cannot update properties and add/remove replication locations or updating enabling of multiple
	// write locations at the same time. so first just update any changed properties
	account := cosmosdb.DatabaseAccountCreateUpdateParameters{
		Location: utils.String(location),
		Kind:     pointer.To(cosmosdb.DatabaseAccountKind(kind)),
		Identity: expandedIdentity,
		Properties: cosmosdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:           cosmosdb.DatabaseAccountOfferType(offerType),
			IPRules:                            ipRangeFilter,
			IsVirtualNetworkFilterEnabled:      utils.Bool(isVirtualNetworkFilterEnabled),
			EnableFreeTier:                     utils.Bool(enableFreeTier),
			EnableAutomaticFailover:            utils.Bool(enableAutomaticFailover),
			Capabilities:                       &capabilities,
			ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                          oldLocations,
			VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			PublicNetworkAccess:                &publicNetworkAccess,
			EnableAnalyticalStorage:            utils.Bool(enableAnalyticalStorage),
			Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]interface{})),
			DisableKeyBasedMetadataWriteAccess: utils.Bool(!d.Get("access_key_metadata_writes_enabled").(bool)),
			NetworkAclBypass:                   &networkByPass,
			NetworkAclBypassResourceIds:        utils.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]interface{})),
			DisableLocalAuth:                   utils.Bool(disableLocalAuthentication),
		},
		Tags: expandTags(t),
	}

	if resp.Model != nil && resp.Model.Properties != nil {
		account.Properties.EnableMultipleWriteLocations = resp.Model.Properties.EnableMultipleWriteLocations
	}

	// d.GetOk cannot identify whether user sets the property that is added Optional and Computed when the property isn't set in TF config file. Because d.GetOk always gets the property value from the last apply when the property isn't set in TF config file. So it has to identify it using `d.GetRawConfig()`
	if v := d.GetRawConfig().AsValueMap()["default_identity_type"]; !v.IsNull() {
		account.Properties.DefaultIdentity = utils.String(v.AsString())
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

	if keyVaultKeyIDRaw, ok := d.GetOk("key_vault_key_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyIDRaw.(string))
		if err != nil {
			return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}
		account.Properties.KeyVaultKeyUri = utils.String(keyVaultKey.ID())
	}

	if v, ok := d.GetOk("mongo_server_version"); ok {
		account.Properties.ApiProperties = &cosmosdb.ApiProperties{
			ServerVersion: pointer.To(cosmosdb.ServerVersion(v.(string))),
		}
	}

	if v, ok := d.GetOk("backup"); ok {
		policy, err := expandCosmosdbAccountBackup(v.([]interface{}), d.HasChange("backup.0.type"), createMode)
		if err != nil {
			return fmt.Errorf("expanding `backup`: %+v", err)
		}
		account.Properties.BackupPolicy = policy
	} else if createMode != "" {
		return fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
	}

	if err = resourceCosmosDbAccountApiUpsert(client, ctx, id, account, d); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// Update the property independently after the initial upsert as no other properties may change at the same time.
	account.Properties.EnableMultipleWriteLocations = utils.Bool(enableMultipleWriteLocations)
	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.EnableMultipleWriteLocations != nil && *resp.Model.Properties.EnableMultipleWriteLocations != enableMultipleWriteLocations {
		if err = resourceCosmosDbAccountApiUpsert(client, ctx, id, account, d); err != nil {
			return fmt.Errorf("updating %s EnableMultipleWriteLocations: %+v", id, err)
		}
	}

	// determine if any locations have been renamed/priority reordered and remove them
	removedOne := false
	for _, l := range newLocations {
		if ol, ok := oldLocationsMap[*l.LocationName]; ok {
			if *l.FailoverPriority != *ol.FailoverPriority {
				if *l.FailoverPriority == 0 {
					return fmt.Errorf("cannot change the failover priority of %s location %s to %d", id, *l.LocationName, *l.FailoverPriority)
				}
				delete(oldLocationsMap, *l.LocationName)
				removedOne = true
				continue
			}
		}
	}

	if removedOne {
		locationsUnchanged := make([]cosmosdb.Location, 0, len(oldLocationsMap))
		for _, value := range oldLocationsMap {
			locationsUnchanged = append(locationsUnchanged, value)
		}

		account.Properties.Locations = locationsUnchanged
		if err = resourceCosmosDbAccountApiUpsert(client, ctx, id, account, d); err != nil {
			return fmt.Errorf("removing %s renamed locations: %+v", id, err)
		}
	}

	// add any new/renamed locations
	account.Properties.Locations = newLocations
	err = resourceCosmosDbAccountApiUpsert(client, ctx, id, account, d)
	if err != nil {
		return fmt.Errorf("updating %s locations: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbAccountRead(d, meta)
}

func resourceCosmosDbAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseDatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DatabaseAccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DatabaseAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {

		cosmosdbLocation := ""
		if model.Location != nil {
			cosmosdbLocation = location.NormalizeNilable(model.Location)
		}
		d.Set("location", cosmosdbLocation)

		kind := ""
		if model.Kind != nil {
			kind = string(*model.Kind)
		}
		d.Set("kind", kind)

		if model.Identity != nil {
			identity, err := flattenAccountIdentity(model.Identity)
			if err != nil {
				return err
			}
			d.Set("identity", identity)
		}

		if props := model.Properties; props != nil {
			d.Set("offer_type", props.DatabaseAccountOfferType)
			if features.FourPointOhBeta() {
				d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilter(props.IPRules))
			} else {
				d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilterThreePointOh(props.IPRules))
			}
			d.Set("endpoint", props.DocumentEndpoint)

			d.Set("enable_free_tier", props.EnableFreeTier)
			d.Set("analytical_storage_enabled", props.EnableAnalyticalStorage)
			publicNetworkAccess := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = *props.PublicNetworkAccess == cosmosdb.PublicNetworkAccessEnabled
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)
			d.Set("default_identity_type", props.DefaultIdentity)
			d.Set("create_mode", props.CreateMode)

			if v := props.IsVirtualNetworkFilterEnabled; v != nil {
				d.Set("is_virtual_network_filter_enabled", props.IsVirtualNetworkFilterEnabled)
			}

			if v := props.EnableAutomaticFailover; v != nil {
				d.Set("enable_automatic_failover", props.EnableAutomaticFailover)
			}

			if v := props.KeyVaultKeyUri; v != nil {
				d.Set("key_vault_key_id", props.KeyVaultKeyUri)
			}

			if v := props.EnableMultipleWriteLocations; v != nil {
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
				d.Set("mongo_server_version", apiProps.ServerVersion)
			}

			networkAclBypassForAzureServices := false
			if props.NetworkAclBypass != nil {
				networkAclBypassForAzureServices = *props.NetworkAclBypass == cosmosdb.NetworkAclBypassAzureServices
			}
			d.Set("network_acl_bypass_for_azure_services", networkAclBypassForAzureServices)
			d.Set("network_acl_bypass_ids", utils.FlattenStringSlice(props.NetworkAclBypassResourceIds))
			if v := props.DisableLocalAuth; v != nil {
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

			readEndpoints := make([]string, 0)
			if p := props.ReadLocations; p != nil {
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
			if p := props.WriteLocations; p != nil {
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

		if keysModel := keys.Model; keysModel != nil {
			primaryKey := ""
			if keysModel.PrimaryMasterKey != nil {
				primaryKey = *keysModel.PrimaryMasterKey
			}
			d.Set("primary_key", primaryKey)

			secondaryKey := ""
			if keysModel.SecondaryMasterKey != nil {
				secondaryKey = *keysModel.SecondaryMasterKey
			}
			d.Set("secondary_key", secondaryKey)
		}

		readonlyKeys, err := client.DatabaseAccountsListReadOnlyKeys(ctx, *id)
		if err != nil {
			if response.WasNotFound(keys.HttpResponse) {
				log.Printf("[DEBUG] Read Only Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.DatabaseAccountName, id.ResourceGroupName)
				d.SetId("")
				return nil
			}

			return fmt.Errorf("[ERROR] Unable to List read-only keys for CosmosDB Account %s: %s", id.DatabaseAccountName, err)
		}

		if readOnlyKeysModel := readonlyKeys.Model; readOnlyKeysModel != nil {

			readOnlyPrimaryKey := ""
			if readOnlyKeysModel.PrimaryReadonlyMasterKey != nil {
				readOnlyPrimaryKey = *readOnlyKeysModel.PrimaryReadonlyMasterKey
			}
			d.Set("primary_readonly_key", readOnlyPrimaryKey)

			readOnlySecondaryKey := ""
			if readOnlyKeysModel.SecondaryReadonlyMasterKey != nil {
				readOnlySecondaryKey = *readOnlyKeysModel.SecondaryReadonlyMasterKey
			}
			d.Set("secondary_readonly_key", readOnlySecondaryKey)
		}

		connStringResp, err := client.DatabaseAccountsListConnectionStrings(ctx, *id)
		if err != nil {
			if response.WasNotFound(keys.HttpResponse) {
				log.Printf("[DEBUG] Connection Strings were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.DatabaseAccountName, id.ResourceGroupName)
				d.SetId("")
				return nil
			}

			return fmt.Errorf("[ERROR] Unable to List connection strings for CosmosDB Account %s: %s", id.DatabaseAccountName, err)
		}

		var connStrings []string
		if connStringModel := connStringResp.Model; connStringModel != nil {
			if connStringModel.ConnectionStrings != nil {
				connStrings = make([]string, len(*connStringModel.ConnectionStrings))
				for i, v := range *connStringModel.ConnectionStrings {
					connStrings[i] = *v.ConnectionString
					if propertyName, propertyExists := connStringPropertyMap[*v.Description]; propertyExists {
						d.Set(propertyName, *v.ConnectionString)
					}
				}
			}
		}
		d.Set("connection_strings", connStrings)

		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceCosmosDbAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseDatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	err = client.DatabaseAccountsDeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
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
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}

func resourceCosmosDbAccountApiUpsert(client *cosmosdb.CosmosDBClient, ctx context.Context, id cosmosdb.DatabaseAccountId, account cosmosdb.DatabaseAccountCreateUpdateParameters, d *pluginsdk.ResourceData) error {
	err := client.DatabaseAccountsCreateOrUpdateThenPoll(ctx, id, account)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	// if a replication location is added or removed it can take some time to provision
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating", "Updating", "Deleting", "Initializing"},
		Target:     []string{"Succeeded"},
		MinTimeout: 30 * time.Second,
		Delay:      30 * time.Second, // required because it takes some time before the 'creating' location shows up
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.DatabaseAccountsGet(ctx, id)
			if err2 != nil || resp.HttpResponse.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("reading %s: %+v", id, err2)
			}
			status := "Succeeded"
			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
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
		return fmt.Errorf("waiting for the %s to provision: %+v", id, err)
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
		policy.MaxStalenessPrefix = utils.Int64(int64(stalenessPrefix))
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
			LocationName:     utils.String(azure.NormalizeLocation(data["location"].(string))),
			FailoverPriority: utils.Int64(int64(data["failover_priority"].(int))),
			IsZoneRedundant:  utils.Bool(data["zone_redundant"].(bool)),
		}

		locations = append(locations, location)
	}

	// TODO maybe this should be in a CustomizeDiff
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
			return nil, fmt.Errorf("each `geo_location` needs to be in unique location. Multiple instances of '%s' found", name)
		}

		if priority > locationsCount-1 {
			return nil, fmt.Errorf("the maximum value for a failover priority = (total number of regions - 1). '%d' was found", priority)
		}

		byPriorities[priority] = location
		byName[name] = location
	}

	// and must have one of 0 priority
	if _, ok := byPriorities[0]; !ok {
		return nil, fmt.Errorf("there needs to be a `geo_location` with a failover_priority of 0")
	}

	return locations, nil
}

func expandAzureRmCosmosDBAccountCapabilities(d *pluginsdk.ResourceData) *[]cosmosdb.Capability {
	capabilities := d.Get("capabilities").(*pluginsdk.Set).List()
	s := make([]cosmosdb.Capability, 0)

	for _, c := range capabilities {
		m := c.(map[string]interface{})
		s = append(s, cosmosdb.Capability{Name: utils.String(m["name"].(string))})
	}

	return &s
}

func expandAzureRmCosmosDBAccountVirtualNetworkRules(d *pluginsdk.ResourceData) *[]cosmosdb.VirtualNetworkRule {
	virtualNetworkRules := d.Get("virtual_network_rule").(*pluginsdk.Set).List()

	s := make([]cosmosdb.VirtualNetworkRule, len(virtualNetworkRules))
	for i, r := range virtualNetworkRules {
		m := r.(map[string]interface{})
		s[i] = cosmosdb.VirtualNetworkRule{
			Id:                               utils.String(m["id"].(string)),
			IgnoreMissingVNetServiceEndpoint: utils.Bool(m["ignore_missing_vnet_service_endpoint"].(bool)),
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
	if account == nil || *account.FailoverPolicies == nil {
		return &locationSet
	}

	for _, l := range *account.FailoverPolicies {
		id := *l.Id
		lb := map[string]interface{}{
			"id":                id,
			"location":          azure.NormalizeLocation(*l.LocationName),
			"failover_priority": int(*l.FailoverPriority),
			// there is not zone redundancy information in the FailoverPolicies currently, we have to search it by `id` in the Locations property.
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

func isServerlessCapacityMode(accResp cosmosdb.DatabaseAccountsGetOperationResponse) bool {
	if model := accResp.Model; model != nil {
		if props := model.Properties; props != nil && props.Capabilities != nil {
			for _, v := range *props.Capabilities {
				if v.Name != nil && *v.Name == "EnableServerless" {
					return true
				}
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
		location := azure.NormalizeLocation(m["location"].(string))
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
			return nil, fmt.Errorf("`interval_in_minutes` can not be set when `type` in `backup` is `Continuous`")
		}
		if v := attr["retention_in_hours"].(int); v != 0 && !backupHasChange {
			return nil, fmt.Errorf("`retention_in_hours` can not be set when `type` in `backup` is `Continuous`")
		}
		if v := attr["storage_redundancy"].(string); v != "" && !backupHasChange {
			return nil, fmt.Errorf("`storage_redundancy` can not be set when `type` in `backup` is `Continuous`")
		}
		return cosmosdb.ContinuousModeBackupPolicy{
			MigrationState: &cosmosdb.BackupPolicyMigrationState{
				TargetType: pointer.To(cosmosdb.BackupPolicyTypeContinuous),
			},
		}, nil

	case string(cosmosdb.BackupPolicyTypePeriodic):
		if createMode != "" {
			return nil, fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
		}

		return cosmosdb.PeriodicModeBackupPolicy{
			MigrationState: &cosmosdb.BackupPolicyMigrationState{
				TargetType: pointer.To(cosmosdb.BackupPolicyTypePeriodic),
			},
			PeriodicModeProperties: &cosmosdb.PeriodicModeProperties{
				BackupIntervalInMinutes:        utils.Int64(int64(attr["interval_in_minutes"].(int))),
				BackupRetentionIntervalInHours: utils.Int64(int64(attr["retention_in_hours"].(int))),
				BackupStorageRedundancy:        pointer.To(cosmosdb.BackupStorageRedundancy(attr["storage_redundancy"].(string))),
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown `type` in `backup`:%+v", attr["type"].(string))
	}
}

func flattenCosmosdbAccountBackup(input cosmosdb.BackupPolicy) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	switch input.(type) {
	case cosmosdb.ContinuousModeBackupPolicy:
		return []interface{}{
			map[string]interface{}{
				"type": string(cosmosdb.BackupPolicyTypeContinuous),
			},
		}, nil

	case cosmosdb.PeriodicModeBackupPolicy:
		policy, ok := input.(cosmosdb.PeriodicModeBackupPolicy)
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
		var storageRedundancy cosmosdb.BackupStorageRedundancy
		if policy.PeriodicModeProperties != nil && policy.PeriodicModeProperties.BackupStorageRedundancy != nil && *policy.PeriodicModeProperties.BackupStorageRedundancy != "" {
			storageRedundancy = *policy.PeriodicModeProperties.BackupStorageRedundancy
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

func expandAccountIdentity(input []interface{}) (*identity.LegacySystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := identity.LegacySystemAndUserAssignedMap{
		Type: expanded.Type,
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.IdentityIds = make(map[string]identity.UserAssignedIdentityDetails)
		for k := range expanded.IdentityIds {
			out.IdentityIds[k] = identity.UserAssignedIdentityDetails{}
		}
	}
	return &out, nil
}

func flattenAccountIdentity(input *identity.LegacySystemAndUserAssignedMap) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        input.Type,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
			PrincipalId: input.PrincipalId,
			TenantId:    input.TenantId,
		}

		for k, v := range input.IdentityIds {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientId,
				PrincipalId: v.PrincipalId,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
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
	if input.SchemaType != nil && *input.SchemaType != "" {
		schemaType = *input.SchemaType
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

	restoreParameters := &cosmosdb.RestoreParameters{
		RestoreMode:        pointer.To(cosmosdb.RestoreModePointInTime),
		RestoreSource:      utils.String(v["source_cosmosdb_account_id"].(string)),
		DatabasesToRestore: expandCosmosdbAccountDatabasesToRestore(v["database"].(*pluginsdk.Set).List()),
	}

	restoreTimestampInUtc, _ := time.Parse(time.RFC3339, v["restore_timestamp_in_utc"].(string))
	restoreParameters.SetRestoreTimestampInUtcAsTime(restoreTimestampInUtc)

	return restoreParameters
}

func expandCosmosdbAccountDatabasesToRestore(input []interface{}) *[]cosmosdb.DatabaseRestoreResource {
	results := make([]cosmosdb.DatabaseRestoreResource, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, cosmosdb.DatabaseRestoreResource{
			DatabaseName:    utils.String(v["name"].(string)),
			CollectionNames: utils.ExpandStringSlice(v["collection_names"].(*pluginsdk.Set).List()),
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
		restoreTimeAsTime, _ := input.GetRestoreTimestampInUtcAsTime()
		restoreTimestampInUtc = restoreTimeAsTime.String()
	}

	return []interface{}{
		map[string]interface{}{
			"database":                   flattenCosmosdbAccountDatabasesToRestore(input.DatabasesToRestore),
			"source_cosmosdb_account_id": restoreSource,
			"restore_timestamp_in_utc":   restoreTimestampInUtc,
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

func checkCapabilitiesCanBeUpdated(kind string, oldCapabilities *[]cosmosdb.Capability, newCapabilities *[]cosmosdb.Capability) bool {
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
			cap := cosmosdb.Capability{
				Name: utils.String(c),
			}
			output = append(output, cap)
		}
	}
	return &output
}
