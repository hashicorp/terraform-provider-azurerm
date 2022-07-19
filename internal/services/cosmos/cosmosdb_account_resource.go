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

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb"
	"github.com/Azure/go-autorest/autorest/date"
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			"default_identity_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "FirstPartyIdentity",
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
				ForceNew: true,
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
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"EnableAggregationPipeline",
								"EnableCassandra",
								"EnableGremlin",
								"EnableTable",
								"EnableServerless",
								"EnableMongo",
								"MongoDBv3.4",
								"mongoEnableDocLevelTTL",
								"DisableRateLimitingResponses",
								"AllowSelfServeUpgradeToMongo36",
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

			"identity": commonschema.SystemAssignedIdentityOptional(),

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
		// todo remove when https://github.com/Azure/azure-sdk-for-go/issues/9891 is fixed
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

	account := documentdb.DatabaseAccountCreateUpdateParameters{
		Location: utils.String(location),
		Kind:     documentdb.DatabaseAccountKind(kind),
		Identity: expandedIdentity,
		DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:           utils.String(offerType),
			IPRules:                            ipRangeFilter,
			IsVirtualNetworkFilterEnabled:      utils.Bool(isVirtualNetworkFilterEnabled),
			EnableFreeTier:                     utils.Bool(enableFreeTier),
			EnableAutomaticFailover:            utils.Bool(enableAutomaticFailover),
			ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                          &geoLocations,
			Capabilities:                       expandAzureRmCosmosDBAccountCapabilities(d),
			VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:       utils.Bool(enableMultipleWriteLocations),
			PublicNetworkAccess:                publicNetworkAccess,
			EnableAnalyticalStorage:            utils.Bool(enableAnalyticalStorage),
			Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]interface{})),
			DisableKeyBasedMetadataWriteAccess: utils.Bool(!d.Get("access_key_metadata_writes_enabled").(bool)),
			NetworkACLBypass:                   networkByPass,
			NetworkACLBypassResourceIds:        utils.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]interface{})),
			DisableLocalAuth:                   utils.Bool(disableLocalAuthentication),
			DefaultIdentity:                    utils.String(d.Get("default_identity_type").(string)),
		},
		Tags: tags.Expand(t),
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
		account.DatabaseAccountCreateUpdateProperties.KeyVaultKeyURI = utils.String(keyVaultKey.ID())
	}

	// additional validation on MaxStalenessPrefix as it varies depending on if the DB is multi region or not
	consistencyPolicy := account.DatabaseAccountCreateUpdateProperties.ConsistencyPolicy
	if len(geoLocations) > 1 && consistencyPolicy != nil && consistencyPolicy.DefaultConsistencyLevel == documentdb.DefaultConsistencyLevelBoundedStaleness {
		if msp := consistencyPolicy.MaxStalenessPrefix; msp != nil && *msp < 100000 {
			return fmt.Errorf("max_staleness_prefix (%d) must be greater then 100000 when more then one geo_location is used", *msp)
		}
		if mis := consistencyPolicy.MaxIntervalInSeconds; mis != nil && *mis < 300 {
			return fmt.Errorf("max_interval_in_seconds (%d) must be greater then 300 (5min) when more then one geo_location is used", *mis)
		}
	}

	err = resourceCosmosDbAccountApiUpsert(client, ctx, id.ResourceGroup, id.Name, account, d)
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

	id := parse.NewDatabaseAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := d.Get("location").(string)
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

	newLocations, err := expandAzureRmCosmosDBAccountGeoLocations(d)
	if err != nil {
		return fmt.Errorf("expanding %s geo locations: %+v", id, err)
	}

	// get existing locations (if exists)
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	oldLocations := make([]documentdb.Location, 0)
	oldLocationsMap := map[string]documentdb.Location{}
	for _, l := range *resp.Locations {
		location := documentdb.Location{
			ID:               l.ID,
			LocationName:     l.LocationName,
			FailoverPriority: l.FailoverPriority,
			IsZoneRedundant:  l.IsZoneRedundant,
		}

		oldLocations = append(oldLocations, location)
		oldLocationsMap[azure.NormalizeLocation(*location.LocationName)] = location
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

	// cannot update properties and add/remove replication locations or updating enabling of multiple
	// write locations at the same time. so first just update any changed properties
	account := documentdb.DatabaseAccountCreateUpdateParameters{
		Location: utils.String(location),
		Kind:     documentdb.DatabaseAccountKind(kind),
		Identity: expandedIdentity,
		DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:           utils.String(offerType),
			IPRules:                            ipRangeFilter,
			IsVirtualNetworkFilterEnabled:      utils.Bool(isVirtualNetworkFilterEnabled),
			EnableFreeTier:                     utils.Bool(enableFreeTier),
			EnableAutomaticFailover:            utils.Bool(enableAutomaticFailover),
			Capabilities:                       expandAzureRmCosmosDBAccountCapabilities(d),
			ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                          &oldLocations,
			VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:       resp.EnableMultipleWriteLocations,
			PublicNetworkAccess:                publicNetworkAccess,
			EnableAnalyticalStorage:            utils.Bool(enableAnalyticalStorage),
			Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]interface{})),
			DisableKeyBasedMetadataWriteAccess: utils.Bool(!d.Get("access_key_metadata_writes_enabled").(bool)),
			NetworkACLBypass:                   networkByPass,
			NetworkACLBypassResourceIds:        utils.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]interface{})),
			DisableLocalAuth:                   utils.Bool(disableLocalAuthentication),
			DefaultIdentity:                    utils.String(d.Get("default_identity_type").(string)),
		},
		Tags: tags.Expand(t),
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

	if keyVaultKeyIDRaw, ok := d.GetOk("key_vault_key_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyIDRaw.(string))
		if err != nil {
			return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}
		account.DatabaseAccountCreateUpdateProperties.KeyVaultKeyURI = utils.String(keyVaultKey.ID())
	}

	if v, ok := d.GetOk("mongo_server_version"); ok {
		account.DatabaseAccountCreateUpdateProperties.APIProperties = &documentdb.APIProperties{
			ServerVersion: documentdb.ServerVersion(v.(string)),
		}
	}

	if v, ok := d.GetOk("backup"); ok {
		policy, err := expandCosmosdbAccountBackup(v.([]interface{}), d.HasChange("backup.0.type"), createMode)
		if err != nil {
			return fmt.Errorf("expanding `backup`: %+v", err)
		}
		account.DatabaseAccountCreateUpdateProperties.BackupPolicy = policy
	} else if createMode != "" {
		return fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
	}

	if err = resourceCosmosDbAccountApiUpsert(client, ctx, id.ResourceGroup, id.Name, account, d); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// Update the property independently after the initial upsert as no other properties may change at the same time.
	account.DatabaseAccountCreateUpdateProperties.EnableMultipleWriteLocations = utils.Bool(enableMultipleWriteLocations)
	if *resp.EnableMultipleWriteLocations != enableMultipleWriteLocations {
		if err = resourceCosmosDbAccountApiUpsert(client, ctx, id.ResourceGroup, id.Name, account, d); err != nil {
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
		locationsUnchanged := make([]documentdb.Location, 0, len(oldLocationsMap))
		for _, value := range oldLocationsMap {
			locationsUnchanged = append(locationsUnchanged, value)
		}

		account.DatabaseAccountCreateUpdateProperties.Locations = &locationsUnchanged
		if err = resourceCosmosDbAccountApiUpsert(client, ctx, id.ResourceGroup, id.Name, account, d); err != nil {
			return fmt.Errorf("removing %s renamed locations: %+v", id, err)
		}
	}

	// add any new/renamed locations
	account.DatabaseAccountCreateUpdateProperties.Locations = &newLocations
	err = resourceCosmosDbAccountApiUpsert(client, ctx, id.ResourceGroup, id.Name, account, d)
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

	id, err := parse.DatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving CosmosDB Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))

	d.Set("kind", string(resp.Kind))

	if v := resp.Identity; v != nil {
		if err := d.Set("identity", flattenAccountIdentity(v)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
	}

	if props := resp.DatabaseAccountGetProperties; props != nil {
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
		defaultIdentity := props.DefaultIdentity
		if defaultIdentity == nil || *defaultIdentity == "" {
			defaultIdentity = utils.String("FirstPartyIdentity")
		}
		d.Set("default_identity_type", defaultIdentity)
		d.Set("create_mode", props.CreateMode)

		if v := resp.IsVirtualNetworkFilterEnabled; v != nil {
			d.Set("is_virtual_network_filter_enabled", props.IsVirtualNetworkFilterEnabled)
		}

		if v := resp.EnableAutomaticFailover; v != nil {
			d.Set("enable_automatic_failover", props.EnableAutomaticFailover)
		}

		if v := resp.KeyVaultKeyURI; v != nil {
			d.Set("key_vault_key_id", props.KeyVaultKeyURI)
		}

		if v := resp.EnableMultipleWriteLocations; v != nil {
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
		if v := resp.DisableLocalAuth; v != nil {
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
	if p := resp.ReadLocations; p != nil {
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
	if p := resp.WriteLocations; p != nil {
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
		}
	}
	d.Set("connection_strings", connStrings)

	return tags.FlattenAndSet(d, resp.Tags)
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

func resourceCosmosDbAccountApiUpsert(client *documentdb.DatabaseAccountsClient, ctx context.Context, resourceGroup string, name string, account documentdb.DatabaseAccountCreateUpdateParameters, d *pluginsdk.ResourceData) error {
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, account)
	if err != nil {
		return fmt.Errorf("creating/updating CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the CosmosDB Account %q (Resource Group %q) to finish creating/updating: %+v", name, resourceGroup, err)
	}

	// if a replication location is added or removed it can take some time to provision
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating", "Updating", "Deleting", "Initializing"},
		Target:     []string{"Succeeded"},
		MinTimeout: 30 * time.Second,
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
		policy.MaxStalenessPrefix = utils.Int64(int64(stalenessPrefix))
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
			LocationName:     utils.String(azure.NormalizeLocation(data["location"].(string))),
			FailoverPriority: utils.Int32(int32(data["failover_priority"].(int))),
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
			return nil, fmt.Errorf("Each `geo_location` needs to have a unique failover_prioroty. Multiple instances of '%d' found", priority)
		}

		if _, ok := byName[name]; ok {
			return nil, fmt.Errorf("Each `geo_location` needs to be in unique location. Multiple instances of '%s' found", name)
		}

		if priority > locationsCount-1 {
			return nil, fmt.Errorf("The maximum value for a failover priority = (total number of regions - 1). '%d' was found", priority)
		}

		byPriorities[priority] = location
		byName[name] = location
	}

	// and must have one of 0 priority
	if _, ok := byPriorities[0]; !ok {
		return nil, fmt.Errorf("There needs to be a `geo_location` with a failover_priority of 0")
	}

	return locations, nil
}

func expandAzureRmCosmosDBAccountCapabilities(d *pluginsdk.ResourceData) *[]documentdb.Capability {
	capabilities := d.Get("capabilities").(*pluginsdk.Set).List()
	s := make([]documentdb.Capability, 0)

	for _, c := range capabilities {
		m := c.(map[string]interface{})
		s = append(s, documentdb.Capability{Name: utils.String(m["name"].(string))})
	}

	return &s
}

func expandAzureRmCosmosDBAccountVirtualNetworkRules(d *pluginsdk.ResourceData) *[]documentdb.VirtualNetworkRule {
	virtualNetworkRules := d.Get("virtual_network_rule").(*pluginsdk.Set).List()

	s := make([]documentdb.VirtualNetworkRule, len(virtualNetworkRules))
	for i, r := range virtualNetworkRules {
		m := r.(map[string]interface{})
		s[i] = documentdb.VirtualNetworkRule{
			ID:                               utils.String(m["id"].(string)),
			IgnoreMissingVNetServiceEndpoint: utils.Bool(m["ignore_missing_vnet_service_endpoint"].(bool)),
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
	if account == nil || *account.FailoverPolicies == nil {
		return &locationSet
	}

	for _, l := range *account.FailoverPolicies {
		id := *l.ID
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
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &documentdb.ManagedServiceIdentity{
		Type: documentdb.ResourceIdentityType(string(expanded.Type)),
	}, nil
}

func flattenAccountIdentity(input *documentdb.ManagedServiceIdentity) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAssigned(transform)
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
		RestoreSource:         utils.String(v["source_cosmosdb_account_id"].(string)),
		RestoreTimestampInUtc: &date.Time{Time: restoreTimestampInUtc},
		DatabasesToRestore:    expandCosmosdbAccountDatabasesToRestore(v["database"].(*pluginsdk.Set).List()),
	}
}

func expandCosmosdbAccountDatabasesToRestore(input []interface{}) *[]documentdb.DatabaseRestoreResource {
	results := make([]documentdb.DatabaseRestoreResource, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, documentdb.DatabaseRestoreResource{
			DatabaseName:    utils.String(v["name"].(string)),
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
