package storage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	azautorest "github.com/Azure/go-autorest/autorest"
	autorestAzure "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyvault "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	vnetParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	resource "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/queue/queues"
)

var (
	storageAccountResourceName = "azurerm_storage_account"
)

func resourceStorageAccount() *pluginsdk.Resource {
	upgraders := map[int]pluginsdk.StateUpgrade{
		0: migration.AccountV0ToV1{},
		1: migration.AccountV1ToV2{},
		2: migration.AccountV2ToV3{},
	}
	schemaVersion := 3

	return &pluginsdk.Resource{
		Create: resourceStorageAccountCreate,
		Read:   resourceStorageAccountRead,
		Update: resourceStorageAccountUpdate,
		Delete: resourceStorageAccountDelete,

		SchemaVersion:  schemaVersion,
		StateUpgraders: pluginsdk.StateUpgrades(upgraders),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StorageAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"account_kind": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.KindStorage),
					string(storage.KindBlobStorage),
					string(storage.KindBlockBlobStorage),
					string(storage.KindFileStorage),
					string(storage.KindStorageV2),
				}, false),
				Default: string(storage.KindStorageV2),
			},

			"account_tier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, false),
			},

			"account_replication_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LRS",
					"ZRS",
					"GRS",
					"RAGRS",
					"GZRS",
					"RAGZRS",
				}, false),
			},

			// Only valid for BlobStorage & StorageV2 accounts, defaults to "Hot" in create function
			"access_tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.AccessTierCool),
					string(storage.AccessTierHot),
				}, false),
			},

			"azure_files_authentication": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"directory_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(storage.DirectoryServiceOptionsAADDS),
								string(storage.DirectoryServiceOptionsAD),
							}, false),
						},

						"active_directory": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"storage_sid": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"domain_guid": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.IsUUID,
									},

									"domain_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"domain_sid": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"forest_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"netbios_domain_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
			"cross_tenant_replication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
			"custom_domain": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"use_subdomain": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"customer_managed_key": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						},

						"user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
					},
				},
			},

			"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_https_traffic_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"min_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(storage.MinimumTLSVersionTLS12),
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.MinimumTLSVersionTLS10),
					string(storage.MinimumTLSVersionTLS11),
					string(storage.MinimumTLSVersionTLS12),
				}, false),
			},

			"is_hns_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"nfsv3_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			// TODO: document this new field in 3.0
			"allow_nested_items_to_be_public": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"shared_access_key_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"network_rules": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"bypass": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(storage.BypassAzureServices),
									string(storage.BypassLogging),
									string(storage.BypassMetrics),
									string(storage.BypassNone),
								}, false),
							},
							Set: pluginsdk.HashString,
						},

						"ip_rules": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.StorageAccountIpRule,
							},
							Set: pluginsdk.HashString,
						},

						"virtual_network_subnet_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(storage.DefaultActionAllow),
								string(storage.DefaultActionDeny),
							}, false),
						},

						"private_link_access": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"endpoint_resource_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},

									"endpoint_tenant_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.IsUUID,
									},
								},
							},
						},
					},
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"blob_properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cors_rule": helpers.SchemaStorageAccountCorsRule(true),
						"delete_retention_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      7,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},

						"versioning_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"change_feed_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"default_service_version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.BlobPropertiesDefaultServiceVersion,
						},

						"last_access_time_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"container_delete_retention_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      7,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},
					},
				},
			},

			"queue_properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cors_rule": helpers.SchemaStorageAccountCorsRule(false),
						"logging": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"version": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"delete": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},
									"read": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},
									"write": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},
									"retention_policy_days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},
						"hour_metrics": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"version": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									// TODO 4.0: Remove this property and determine whether to enable based on existence of the out side block.
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},
									"include_apis": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},
									"retention_policy_days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},
						"minute_metrics": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"version": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									// TODO 4.0: Remove this property and determine whether to enable based on existence of the out side block.
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},
									"include_apis": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},
									"retention_policy_days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},
					},
				},
			},

			"routing": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"publish_internet_endpoints": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"publish_microsoft_endpoints": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"choice": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(storage.RoutingChoiceMicrosoftRouting),
								string(storage.RoutingChoiceInternetRouting),
							}, false),
							Default: string(storage.RoutingChoiceMicrosoftRouting),
						},
					},
				},
			},

			"share_properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cors_rule": helpers.SchemaStorageAccountCorsRule(true),

						"retention_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      7,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},

						"smb": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"versions": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"SMB2.1",
												"SMB3.0",
												"SMB3.1.1",
											}, false),
										},
									},

									"authentication_types": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"NTLMv2",
												"Kerberos",
											}, false),
										},
									},

									"kerberos_ticket_encryption_type": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"RC4-HMAC",
												"AES-256",
											}, false),
										},
									},

									"channel_encryption_type": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"AES-128-CCM",
												"AES-128-GCM",
												"AES-256-GCM",
											}, false),
										},
									},
								},
							},
						},
					},
				},
			},

			// lintignore:XS003
			"static_website": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"index_document": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"error_404_document": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"queue_encryption_key_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.KeyTypeService),
					string(storage.KeyTypeAccount),
				}, false),
				Default: string(storage.KeyTypeService),
			},

			"table_encryption_key_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.KeyTypeService),
					string(storage.KeyTypeAccount),
				}, false),
				Default: string(storage.KeyTypeService),
			},

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"large_file_share_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"primary_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_blob_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_blob_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_blob_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_blob_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_table_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_table_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_table_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_table_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_web_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_web_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_web_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_web_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_dfs_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_dfs_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_dfs_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_dfs_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_file_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_file_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_file_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_file_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
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

			"primary_blob_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_blob_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": {
				Type:         pluginsdk.TypeMap,
				Optional:     true,
				ValidateFunc: validate.StorageAccountTags,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				if d.HasChange("account_kind") {
					accountKind, changedKind := d.GetChange("account_kind")

					if accountKind != string(storage.KindStorage) && changedKind != string(storage.KindStorageV2) {
						log.Printf("[DEBUG] recreate storage account, could't be migrated from %s to %s", accountKind, changedKind)
						d.ForceNew("account_kind")
					} else {
						log.Printf("[DEBUG] storage account can be upgraded from %s to %s", accountKind, changedKind)
					}
				}

				if d.HasChange("large_file_share_enabled") {
					lfsEnabled, changedEnabled := d.GetChange("large_file_share_enabled")
					if lfsEnabled.(bool) && !changedEnabled.(bool) {
						return fmt.Errorf("`large_file_share_enabled` cannot be disabled once it's been enabled")
					}
				}
				return nil
			}),
			pluginsdk.ForceNewIfChange("account_replication_type", func(ctx context.Context, old, new, meta interface{}) bool {
				newAccRep := strings.ToUpper(new.(string))

				switch strings.ToUpper(old.(string)) {
				case "LRS", "GRS", "RAGRS":
					if newAccRep == "GZRS" || newAccRep == "RAGZRS" || newAccRep == "ZRS" {
						return true
					}
				case "ZRS", "GZRS", "RAGZRS":
					if newAccRep == "LRS" || newAccRep == "GRS" || newAccRep == "RAGRS" {
						return true
					}
				}
				return false
			}),
		),
	}
}

func resourceStorageAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	envName := meta.(*clients.Client).Account.Environment.Name
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	storageClient := meta.(*clients.Client).Storage
	keyVaultClient := meta.(*clients.Client).KeyVault
	resourceClient := meta.(*clients.Client).Resource
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewStorageAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.Name, storageAccountResourceName)
	defer locks.UnlockByName(id.Name, storageAccountResourceName)

	existing, err := client.GetProperties(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_storage_account", id.ID())
	}

	accountKind := d.Get("account_kind").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	enableHTTPSTrafficOnly := d.Get("enable_https_traffic_only").(bool)
	minimumTLSVersion := d.Get("min_tls_version").(string)
	isHnsEnabled := d.Get("is_hns_enabled").(bool)
	nfsV3Enabled := d.Get("nfsv3_enabled").(bool)
	allowBlobPublicAccess := d.Get("allow_nested_items_to_be_public").(bool)
	allowSharedKeyAccess := d.Get("shared_access_key_enabled").(bool)
	crossTenantReplication := d.Get("cross_tenant_replication_enabled").(bool)

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)

	parameters := storage.AccountCreateParameters{
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         &location,
		Sku: &storage.Sku{
			Name: storage.SkuName(storageType),
		},
		Tags: tags.Expand(t),
		Kind: storage.Kind(accountKind),
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
			EnableHTTPSTrafficOnly:      &enableHTTPSTrafficOnly,
			NetworkRuleSet:              expandStorageAccountNetworkRules(d, tenantId),
			IsHnsEnabled:                &isHnsEnabled,
			EnableNfsV3:                 &nfsV3Enabled,
			AllowSharedKeyAccess:        &allowSharedKeyAccess,
			AllowCrossTenantReplication: &crossTenantReplication,
		},
	}

	// For all Clouds except Public, China, and USGovernmentCloud, don't specify "allow_blob_public_access" and "min_tls_version" in request body.
	// https://github.com/hashicorp/terraform-provider-azurerm/issues/7812
	// https://github.com/hashicorp/terraform-provider-azurerm/issues/8083
	// USGovernmentCloud allow_blob_public_access and min_tls_version allowed as of issue 9128
	// https://github.com/hashicorp/terraform-provider-azurerm/issues/9128
	if envName != autorestAzure.PublicCloud.Name && envName != autorestAzure.USGovernmentCloud.Name && envName != autorestAzure.ChinaCloud.Name {
		if allowBlobPublicAccess || minimumTLSVersion != string(storage.MinimumTLSVersionTLS10) {
			return fmt.Errorf(`"allow_nested_items_to_be_public" and "min_tls_version" are not supported for a Storage Account located in %q`, envName)
		}
	} else {
		parameters.AccountPropertiesCreateParameters.AllowBlobPublicAccess = &allowBlobPublicAccess
		parameters.AccountPropertiesCreateParameters.MinimumTLSVersion = storage.MinimumTLSVersion(minimumTLSVersion)
	}

	storageAccountIdentity, err := expandAzureRmStorageAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return err
	}
	parameters.Identity = storageAccountIdentity

	if v, ok := d.GetOk("azure_files_authentication"); ok {
		expandAADFilesAuthentication, err := expandArmStorageAccountAzureFilesAuthentication(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("parsing `azure_files_authentication`: %v", err)
		}
		parameters.AzureFilesIdentityBasedAuthentication = expandAADFilesAuthentication
	}

	if _, ok := d.GetOk("custom_domain"); ok {
		parameters.CustomDomain = expandStorageAccountCustomDomain(d)
	}

	// BlobStorage does not support ZRS
	if accountKind == string(storage.KindBlobStorage) {
		if string(parameters.Sku.Name) == string(storage.SkuNameStandardZRS) {
			return fmt.Errorf("A `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts.")
		}
	}

	// AccessTier is only valid for BlobStorage, StorageV2, and FileStorage accounts
	if accountKind == string(storage.KindBlobStorage) || accountKind == string(storage.KindStorageV2) || accountKind == string(storage.KindFileStorage) {
		accessTier, ok := d.GetOk("access_tier")
		if !ok {
			// default to "Hot"
			accessTier = string(storage.AccessTierHot)
		}

		parameters.AccountPropertiesCreateParameters.AccessTier = storage.AccessTier(accessTier.(string))
	} else if isHnsEnabled && accountKind != string(storage.KindBlockBlobStorage) {
		return fmt.Errorf("`is_hns_enabled` can only be used with account kinds `StorageV2`, `BlobStorage` and `BlockBlobStorage`")
	}

	// NFSv3 is supported for standard general-purpose v2 storage accounts and for premium block blob storage accounts.
	// (https://docs.microsoft.com/en-us/azure/storage/blobs/network-file-system-protocol-support-how-to#step-5-create-and-configure-a-storage-account)
	if nfsV3Enabled &&
		!((accountTier == string(storage.SkuTierPremium) && accountKind == string(storage.KindBlockBlobStorage)) ||
			(accountTier == string(storage.SkuTierStandard) && accountKind == string(storage.KindStorageV2))) {
		return fmt.Errorf("`nfsv3_enabled` can only be used with account tier `Standard` and account kind `StorageV2`, or account tier `Premium` and account kind `BlockBlobStorage`")
	}
	if nfsV3Enabled && !isHnsEnabled {
		return fmt.Errorf("`nfsv3_enabled` can only be used when `is_hns_enabled` is `true`")
	}

	// AccountTier must be Premium for FileStorage
	if accountKind == string(storage.KindFileStorage) {
		if string(parameters.Sku.Tier) == string(storage.SkuNameStandardLRS) {
			return fmt.Errorf("A `account_tier` of `Standard` is not supported for FileStorage accounts.")
		}
	}

	// nolint staticcheck
	if v, ok := d.GetOkExists("large_file_share_enabled"); ok {
		parameters.LargeFileSharesState = storage.LargeFileSharesStateDisabled
		if v.(bool) {
			parameters.LargeFileSharesState = storage.LargeFileSharesStateEnabled
		}
	}

	if v, ok := d.GetOk("routing"); ok {
		parameters.RoutingPreference = expandArmStorageAccountRouting(v.([]interface{}))
	}

	// TODO 4.0
	// look into standardizing this across resources that support CMK and at the very least look at improving the UX
	// for encryption of blob, file, table and queue
	var encryption *storage.Encryption
	if v, ok := d.GetOk("customer_managed_key"); ok {
		if accountKind != string(storage.KindStorageV2) {
			return fmt.Errorf("customer managed key can only be used with account kind `StorageV2`")
		}
		if storageAccountIdentity.Type != storage.IdentityTypeUserAssigned {
			return fmt.Errorf("customer managed key can only be used with identity type `UserAssigned`")
		}
		encryption, err = expandStorageAccountCustomerManagedKey(ctx, keyVaultClient, resourceClient, v.([]interface{}))
		if err != nil {
			return err
		}
	}

	// By default (by leaving empty), the table and queue encryption key type is set to "Service". While users can change it to "Account" so that
	// they can further use CMK to encrypt table/queue data. Only the StorageV2 account kind supports the Account key type.
	// Also noted that the blob and file are always using the "Account" key type.
	// See: https://docs.microsoft.com/en-gb/azure/storage/common/account-encryption-key-create?tabs=portal
	queueEncryptionKeyType := d.Get("queue_encryption_key_type").(string)
	tableEncryptionKeyType := d.Get("table_encryption_key_type").(string)

	if accountKind != string(storage.KindStorageV2) {
		if queueEncryptionKeyType == string(storage.KeyTypeAccount) {
			return fmt.Errorf("`queue_encryption_key_type = \"Account\"` can only be used with account kind `StorageV2`")
		}
		if tableEncryptionKeyType == string(storage.KeyTypeAccount) {
			return fmt.Errorf("`table_encryption_key_type = \"Account\"` can only be used with account kind `StorageV2`")
		}
	}

	// if CMK is not supplied then only set storage encryption for queue and table, otherwise add it to the existing encryption block for CMK
	if encryption == nil {
		encryption = &storage.Encryption{
			KeySource: storage.KeySourceMicrosoftStorage,
			Services: &storage.EncryptionServices{
				Queue: &storage.EncryptionService{
					KeyType: storage.KeyType(queueEncryptionKeyType),
				},
				Table: &storage.EncryptionService{
					KeyType: storage.KeyType(tableEncryptionKeyType),
				},
			},
		}
	} else {
		encryption.Services.Queue = &storage.EncryptionService{
			KeyType: storage.KeyType(queueEncryptionKeyType),
		}
		encryption.Services.Table = &storage.EncryptionService{
			KeyType: storage.KeyType(queueEncryptionKeyType),
		}
	}

	infrastructureEncryption := d.Get("infrastructure_encryption_enabled").(bool)

	if infrastructureEncryption {
		if !((accountTier == string(storage.SkuTierPremium) && accountKind == string(storage.KindBlockBlobStorage)) ||
			(accountKind == string(storage.KindStorageV2))) {
			return fmt.Errorf("`infrastructure_encryption_enabled` can only be used with account kind `StorageV2`, or account tier `Premium` and account kind `BlockBlobStorage`")
		}
		encryption.RequireInfrastructureEncryption = &infrastructureEncryption
	}

	parameters.Encryption = encryption

	// Create
	future, err := client.Create(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating Azure Storage Account %q: %+v", id.Name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Azure Storage Account %q to be created: %+v", id.Name, err)
	}

	d.SetId(id.ID())

	// populate the cache
	account, err := client.GetProperties(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if err := storageClient.AddToCache(id.Name, account); err != nil {
		return fmt.Errorf("populating cache for %s: %+v", id, err)
	}

	if val, ok := d.GetOk("blob_properties"); ok {
		// FileStorage does not support blob settings
		if accountKind != string(storage.KindFileStorage) {
			blobClient := meta.(*clients.Client).Storage.BlobServicesClient

			blobProperties := expandBlobProperties(val.([]interface{}))

			// last_access_time_enabled and container_delete_retention_policy are not supported in USGov
			// Fix issue https://github.com/hashicorp/terraform-provider-azurerm/issues/11772
			if v := d.Get("blob_properties.0.last_access_time_enabled").(bool); v {
				blobProperties.LastAccessTimeTrackingPolicy = &storage.LastAccessTimeTrackingPolicy{
					Enable: utils.Bool(v),
				}
			}

			if v, ok := d.GetOk("blob_properties.0.container_delete_retention_policy"); ok {
				blobProperties.ContainerDeleteRetentionPolicy = expandBlobPropertiesDeleteRetentionPolicy(v.([]interface{}), false)
			}

			if _, err = blobClient.SetServiceProperties(ctx, id.ResourceGroup, id.Name, *blobProperties); err != nil {
				return fmt.Errorf("updating Azure Storage Account `blob_properties` %q: %+v", id.Name, err)
			}
		} else {
			return fmt.Errorf("`blob_properties` aren't supported for File Storage accounts.")
		}
	}

	if val, ok := d.GetOk("queue_properties"); ok {
		storageClient := meta.(*clients.Client).Storage
		account, err := storageClient.FindAccount(ctx, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Account %q: %s", id.Name, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", id.Name)
		}

		queueClient, err := storageClient.QueuesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("building Queues Client: %s", err)
		}

		queueProperties, err := expandQueueProperties(val.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `queue_properties` for Azure Storage Account %q: %+v", id.Name, err)
		}

		if err = queueClient.UpdateServiceProperties(ctx, account.ResourceGroup, id.Name, queueProperties); err != nil {
			return fmt.Errorf("updating Queue Properties for Storage Account %q: %+v", id.Name, err)
		}
	}

	if val, ok := d.GetOk("share_properties"); ok {
		// BlobStorage does not support file share settings
		// FileStorage Premium is supported
		if accountKind == string(storage.KindFileStorage) || accountKind != string(storage.KindBlobStorage) && accountKind != string(storage.KindBlockBlobStorage) && accountTier != string(storage.SkuTierPremium) {
			fileServiceClient := meta.(*clients.Client).Storage.FileServicesClient

			if _, err = fileServiceClient.SetServiceProperties(ctx, id.ResourceGroup, id.Name, expandShareProperties(val.([]interface{}))); err != nil {
				return fmt.Errorf("updating Azure Storage Account `share_properties` %q: %+v", id.Name, err)
			}
		} else {
			return fmt.Errorf("`share_properties` aren't supported for Blob Storage / Block Blob / StorageV2 Premium Storage accounts")
		}
	}

	if val, ok := d.GetOk("static_website"); ok {
		// static website only supported on StorageV2 and BlockBlobStorage
		if accountKind != string(storage.KindStorageV2) && accountKind != string(storage.KindBlockBlobStorage) {
			return fmt.Errorf("`static_website` is only supported for StorageV2 and BlockBlobStorage.")
		}
		storageClient := meta.(*clients.Client).Storage

		account, err := storageClient.FindAccount(ctx, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Account %q: %s", id.Name, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", id.Name)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps := expandStaticWebsiteProperties(val.([]interface{}))

		if _, err = accountsClient.SetServiceProperties(ctx, id.Name, staticWebsiteProps); err != nil {
			return fmt.Errorf("updating Azure Storage Account `static_website` %q: %+v", id.Name, err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}

func resourceStorageAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	envName := meta.(*clients.Client).Account.Environment.Name
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	keyVaultClient := meta.(*clients.Client).KeyVault
	resourceClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, storageAccountResourceName)
	defer locks.UnlockByName(id.Name, storageAccountResourceName)

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	accountKind := d.Get("account_kind").(string)

	if accountKind == string(storage.KindBlobStorage) || accountKind == string(storage.KindStorage) {
		if storageType == string(storage.SkuNameStandardZRS) {
			return fmt.Errorf("A `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts.")
		}
	}

	allowSharedKeyAccess := true
	if d.HasChange("shared_access_key_enabled") {
		allowSharedKeyAccess = d.Get("shared_access_key_enabled").(bool)

		// If AllowSharedKeyAccess is nil that breaks the Portal UI as reported in https://github.com/hashicorp/terraform-provider-azurerm/issues/11689
		// currently the Portal UI reports nil as false, and per the ARM API documentation nil is true. This manifests itself in the Portal UI
		// when a storage account is created by terraform that the AllowSharedKeyAccess is Disabled when it is actually Enabled, thus confusing out customers
		// to fix this, I have added this code to explicitly to set the value to true if is nil to workaround the Portal UI bug for our customers.
		// this is designed as a passive change, meaning the change will only take effect when the existing storage account is modified in some way if the
		// account already exists. since I have also switched up the default behaviour for net new storage accounts to always set this value as true, this issue
		// should automatically correct itself over time with these changes.
		// TODO: Remove code when Portal UI team fixes their code
	} else {
		existing, err := client.GetProperties(ctx, id.ResourceGroup, id.Name, "")
		if err == nil {
			if sharedKeyAccess := existing.AccountProperties.AllowSharedKeyAccess; sharedKeyAccess != nil {
				allowSharedKeyAccess = *sharedKeyAccess
			}
		} else {
			// Should never hit this, but added due to an abundance of caution
			return fmt.Errorf("retrieving Azure Storage Account %q AllowSharedKeyAccess: %+v", id.Name, err)
		}
	}
	// TODO: end remove changes when Portal UI team fixed their code

	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			AllowSharedKeyAccess: &allowSharedKeyAccess,
		},
	}

	if d.HasChange("cross_tenant_replication_enabled") {
		crossTenantReplication := d.Get("cross_tenant_replication_enabled").(bool)
		opts.AccountPropertiesUpdateParameters.AllowCrossTenantReplication = &crossTenantReplication
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
		return fmt.Errorf("updating Azure Storage Account AllowSharedKeyAccess %q: %+v", id.Name, err)
	}

	if d.HasChange("account_replication_type") {
		sku := storage.Sku{
			Name: storage.SkuName(storageType),
		}

		opts := storage.AccountUpdateParameters{
			Sku: &sku,
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account type %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("account_kind") {
		opts := storage.AccountUpdateParameters{
			Kind: storage.Kind(accountKind),
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account account_kind %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("access_tier") {
		accessTier := d.Get("access_tier").(string)

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				AccessTier: storage.AccessTier(accessTier),
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account access_tier %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})

		opts := storage.AccountUpdateParameters{
			Tags: tags.Expand(t),
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account tags %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("custom_domain") {
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				CustomDomain: expandStorageAccountCustomDomain(d),
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account Custom Domain %q: %+v", id.Name, err)
		}
	}

	// Updating `identity` should occur before updating `customer_managed_key`, as the latter depends on an identity.
	if d.HasChange("identity") {
		storageAccountIdentity, err := expandAzureRmStorageAccountIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return err
		}
		opts := storage.AccountUpdateParameters{
			Identity: storageAccountIdentity,
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account identity %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("customer_managed_key") {
		cmk := d.Get("customer_managed_key").([]interface{})
		encryption, err := expandStorageAccountCustomerManagedKey(ctx, keyVaultClient, resourceClient, cmk)
		if err != nil {
			return err
		}

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				Encryption: encryption,
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating %s Customer Managed Key: %+v", *id, err)
		}

	}

	if d.HasChange("enable_https_traffic_only") {
		enableHTTPSTrafficOnly := d.Get("enable_https_traffic_only").(bool)

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				EnableHTTPSTrafficOnly: &enableHTTPSTrafficOnly,
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account enable_https_traffic_only %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("min_tls_version") {
		minimumTLSVersion := d.Get("min_tls_version").(string)

		// For all Clouds except Public, China, and USGovernmentCloud, don't specify "min_tls_version" in request body.
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/8083
		// USGovernmentCloud "min_tls_version" allowed as of issue 9128
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/9128
		if envName != autorestAzure.PublicCloud.Name && envName != autorestAzure.USGovernmentCloud.Name && envName != autorestAzure.ChinaCloud.Name {
			if minimumTLSVersion != string(storage.MinimumTLSVersionTLS10) {
				return fmt.Errorf(`"min_tls_version" is not supported for a Storage Account located in %q`, envName)
			}
		} else {
			opts := storage.AccountUpdateParameters{
				AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
					MinimumTLSVersion: storage.MinimumTLSVersion(minimumTLSVersion),
				},
			}

			if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
				return fmt.Errorf("updating Azure Storage Account min_tls_version %q: %+v", id.Name, err)
			}
		}
	}

	if d.HasChange("allow_nested_items_to_be_public") {
		allowBlobPublicAccess := d.Get("allow_nested_items_to_be_public").(bool)

		// For all Clouds except Public, China, and USGovernmentCloud, don't specify "allow_blob_public_access" in request body.
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/7812
		// USGovernmentCloud "allow_blob_public_access" allowed as of issue 9128
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/9128
		if envName != autorestAzure.PublicCloud.Name && envName != autorestAzure.USGovernmentCloud.Name && envName != autorestAzure.ChinaCloud.Name {
			if allowBlobPublicAccess {
				return fmt.Errorf("allow_nested_items_to_be_public is not supported for a Storage Account located in %q", envName)
			}
		} else {
			opts := storage.AccountUpdateParameters{
				AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
					AllowBlobPublicAccess: &allowBlobPublicAccess,
				},
			}

			if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
				return fmt.Errorf("updating Azure Storage Account allow_blob_public_access %q: %+v", id.Name, err)
			}
		}
	}

	if d.HasChange("network_rules") {
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				NetworkRuleSet: expandStorageAccountNetworkRules(d, tenantId),
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account network_rules %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("large_file_share_enabled") {
		isEnabled := storage.LargeFileSharesStateDisabled
		if v := d.Get("large_file_share_enabled").(bool); v {
			isEnabled = storage.LargeFileSharesStateEnabled
		}
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				LargeFileSharesState: isEnabled,
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account network_rules %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("routing") {
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				RoutingPreference: expandArmStorageAccountRouting(d.Get("routing").([]interface{})),
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account routing %q: %+v", id.Name, err)
		}
	}

	// azure_files_authentication must be the last to be updated, cause it'll occupy the storage account for several minutes after receiving the response 200 OK. Issue: https://github.com/Azure/azure-rest-api-specs/issues/11272
	if d.HasChange("azure_files_authentication") {
		// due to service issue: https://github.com/Azure/azure-rest-api-specs/issues/12473, we need to update to None before changing its DirectoryServiceOptions
		old, new := d.GetChange("azure_files_authentication.0.directory_type")
		if old != new && new != string(storage.DirectoryServiceOptionsNone) {
			log.Print("[DEBUG] Disabling AzureFilesIdentityBasedAuthentication prior to changing DirectoryServiceOptions")
			dsNone := storage.AccountUpdateParameters{
				AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
					AzureFilesIdentityBasedAuthentication: &storage.AzureFilesIdentityBasedAuthentication{
						DirectoryServiceOptions: storage.DirectoryServiceOptionsNone,
					},
				},
			}
			if _, err := client.Update(ctx, id.ResourceGroup, id.Name, dsNone); err != nil {
				return fmt.Errorf("updating Azure Storage Account azure_files_authentication %q: %+v", id.Name, err)
			}
		}

		expandAADFilesAuthentication, err := expandArmStorageAccountAzureFilesAuthentication(d.Get("azure_files_authentication").([]interface{}))
		if err != nil {
			return err
		}
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				AzureFilesIdentityBasedAuthentication: expandAADFilesAuthentication,
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroup, id.Name, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account azure_files_authentication %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("blob_properties") {
		// FileStorage does not support blob settings
		if accountKind != string(storage.KindFileStorage) {
			blobClient := meta.(*clients.Client).Storage.BlobServicesClient
			blobProperties := expandBlobProperties(d.Get("blob_properties").([]interface{}))

			// last_access_time_enabled and container_delete_retention_policy are not supported in USGov
			// Fix issue https://github.com/hashicorp/terraform-provider-azurerm/issues/11772
			if d.HasChange("blob_properties.0.last_access_time_enabled") {
				lastAccessTimeTracking := false
				if v := d.Get("blob_properties.0.last_access_time_enabled").(bool); v {
					lastAccessTimeTracking = true
				}
				blobProperties.LastAccessTimeTrackingPolicy = &storage.LastAccessTimeTrackingPolicy{
					Enable: utils.Bool(lastAccessTimeTracking),
				}
			}

			if d.HasChange("blob_properties.0.container_delete_retention_policy") {
				blobProperties.ContainerDeleteRetentionPolicy = expandBlobPropertiesDeleteRetentionPolicy(d.Get("blob_properties.0.container_delete_retention_policy").([]interface{}), true)
			}

			if _, err = blobClient.SetServiceProperties(ctx, id.ResourceGroup, id.Name, *blobProperties); err != nil {
				return fmt.Errorf("updating Azure Storage Account `blob_properties` %q: %+v", id.Name, err)
			}
		} else {
			return fmt.Errorf("`blob_properties` aren't supported for File Storage accounts.")
		}
	}

	if d.HasChange("queue_properties") {
		storageClient := meta.(*clients.Client).Storage
		account, err := storageClient.FindAccount(ctx, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Account %q: %s", id.Name, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", id.Name)
		}

		queueClient, err := storageClient.QueuesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("building Queues Client: %s", err)
		}

		queueProperties, err := expandQueueProperties(d.Get("queue_properties").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `queue_properties` for Azure Storage Account %q: %+v", id.Name, err)
		}

		if err = queueClient.UpdateServiceProperties(ctx, account.ResourceGroup, id.Name, queueProperties); err != nil {
			return fmt.Errorf("updating Queue Properties for Storage Account %q: %+v", id.Name, err)
		}
	}

	if d.HasChange("share_properties") {
		// BlobStorage, BlockBlobStorage does not support file share settings
		// FileStorage Premium is supported
		if accountKind == string(storage.KindFileStorage) || accountKind != string(storage.KindBlobStorage) && accountKind != string(storage.KindBlockBlobStorage) && accountTier != string(storage.SkuTierPremium) {
			fileServiceClient := meta.(*clients.Client).Storage.FileServicesClient

			if _, err = fileServiceClient.SetServiceProperties(ctx, id.ResourceGroup, id.Name, expandShareProperties(d.Get("share_properties").([]interface{}))); err != nil {
				return fmt.Errorf("updating Azure Storage Account `file share_properties` %q: %+v", id.Name, err)
			}
		} else {
			return fmt.Errorf("`share_properties` aren't supported for Blob Storage /Block Blob /StorageV2 Premium Storage accounts")
		}
	}

	if d.HasChange("static_website") {
		// static website only supported on StorageV2 and BlockBlobStorage
		if accountKind != string(storage.KindStorageV2) && accountKind != string(storage.KindBlockBlobStorage) {
			return fmt.Errorf("`static_website` is only supported for StorageV2 and BlockBlobStorage.")
		}
		storageClient := meta.(*clients.Client).Storage

		account, err := storageClient.FindAccount(ctx, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Account %q: %s", id.Name, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", id.Name)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps := expandStaticWebsiteProperties(d.Get("static_website").([]interface{}))

		if _, err = accountsClient.SetServiceProperties(ctx, id.Name, staticWebsiteProps); err != nil {
			return fmt.Errorf("updating Azure Storage Account `static_website` %q: %+v", id.Name, err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}

func resourceStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading the state of AzureRM Storage Account %q: %+v", id.Name, err)
	}

	// handle the user not having permissions to list the keys
	d.Set("primary_connection_string", "")
	d.Set("secondary_connection_string", "")
	d.Set("primary_blob_connection_string", "")
	d.Set("secondary_blob_connection_string", "")
	d.Set("primary_access_key", "")
	d.Set("secondary_access_key", "")

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.Name, storage.ListKeyExpandKerb)
	if err != nil {
		// the API returns a 200 with an inner error of a 409..
		var hasWriteLock bool
		var doesntHavePermissions bool
		if e, ok := err.(azautorest.DetailedError); ok {
			if status, ok := e.StatusCode.(int); ok {
				hasWriteLock = status == http.StatusConflict
				doesntHavePermissions = status == http.StatusUnauthorized
			}
		}

		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("listing Keys for %s: %s", *id, err)
		}
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("edge_zone", flattenEdgeZone(resp.ExtendedLocation))
	d.Set("account_kind", resp.Kind)

	if sku := resp.Sku; sku != nil {
		d.Set("account_tier", sku.Tier)
		d.Set("account_replication_type", strings.Split(fmt.Sprintf("%v", sku.Name), "_")[1])
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("access_tier", props.AccessTier)
		if err := d.Set("azure_files_authentication", flattenArmStorageAccountAzureFilesAuthentication(props.AzureFilesIdentityBasedAuthentication)); err != nil {
			return fmt.Errorf("setting `azure_files_authentication`: %+v", err)
		}
		if err := d.Set("routing", flattenArmStorageAccountRouting(props.RoutingPreference)); err != nil {
			return fmt.Errorf("setting `routing`: %+v", err)
		}
		d.Set("enable_https_traffic_only", props.EnableHTTPSTrafficOnly)
		d.Set("is_hns_enabled", props.IsHnsEnabled)
		d.Set("nfsv3_enabled", props.EnableNfsV3)

		if crossTenantReplication := props.AllowCrossTenantReplication; crossTenantReplication != nil {
			d.Set("cross_tenant_replication_enabled", crossTenantReplication)
		}

		// There is a certain edge case that could result in the Azure API returning a null value for AllowBLobPublicAccess.
		// Since the field is a pointer, this gets marshalled to a nil value instead of a boolean.

		allowBlobPublicAccess := true
		if props.AllowBlobPublicAccess != nil {
			allowBlobPublicAccess = *props.AllowBlobPublicAccess
		}
		//lintignore:R001
		d.Set("allow_nested_items_to_be_public", allowBlobPublicAccess)

		// For all Clouds except Public, China, and USGovernmentCloud, "min_tls_version" is not returned from Azure so always persist the default values for "min_tls_version".
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/7812
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/8083
		// USGovernmentCloud "min_tls_version" allowed as of issue 9128
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/9128
		envName := meta.(*clients.Client).Account.Environment.Name
		if envName != autorestAzure.PublicCloud.Name && envName != autorestAzure.USGovernmentCloud.Name && envName != autorestAzure.ChinaCloud.Name {
			d.Set("min_tls_version", string(storage.MinimumTLSVersionTLS10))
		} else {
			// For storage account created using old API, the response of GET call will not return "min_tls_version", either.
			minTlsVersion := string(storage.MinimumTLSVersionTLS10)
			if props.MinimumTLSVersion != "" {
				minTlsVersion = string(props.MinimumTLSVersion)
			}
			d.Set("min_tls_version", minTlsVersion)
		}

		if customDomain := props.CustomDomain; customDomain != nil {
			if err := d.Set("custom_domain", flattenStorageAccountCustomDomain(customDomain)); err != nil {
				return fmt.Errorf("setting `custom_domain`: %+v", err)
			}
		}

		// Computed
		d.Set("primary_location", props.PrimaryLocation)
		d.Set("secondary_location", props.SecondaryLocation)

		if accessKeys := keys.Keys; accessKeys != nil {
			storageAccountKeys := *accessKeys
			if len(storageAccountKeys) > 0 {
				pcs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *storageAccountKeys[0].Value, endpointSuffix)
				d.Set("primary_connection_string", pcs)
			}

			if len(storageAccountKeys) > 1 {
				scs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *storageAccountKeys[1].Value, endpointSuffix)
				d.Set("secondary_connection_string", scs)
			}
		}

		if err := flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d, props.PrimaryEndpoints); err != nil {
			return fmt.Errorf("setting primary endpoints and hosts for blob, queue, table and file: %+v", err)
		}

		if accessKeys := keys.Keys; accessKeys != nil {
			storageAccountKeys := *accessKeys
			var primaryBlobConnectStr string
			if v := props.PrimaryEndpoints; v != nil {
				primaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, storageAccountKeys[0].Value)
			}
			d.Set("primary_blob_connection_string", primaryBlobConnectStr)
		}

		if err := flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d, props.SecondaryEndpoints); err != nil {
			return fmt.Errorf("setting secondary endpoints and hosts for blob, queue, table: %+v", err)
		}

		if accessKeys := keys.Keys; accessKeys != nil {
			storageAccountKeys := *accessKeys
			var secondaryBlobConnectStr string
			if v := props.SecondaryEndpoints; v != nil {
				secondaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, storageAccountKeys[1].Value)
			}
			d.Set("secondary_blob_connection_string", secondaryBlobConnectStr)
		}

		if err := d.Set("network_rules", flattenStorageAccountNetworkRules(props.NetworkRuleSet)); err != nil {
			return fmt.Errorf("setting `network_rules`: %+v", err)
		}

		if props.LargeFileSharesState != "" {
			d.Set("large_file_share_enabled", props.LargeFileSharesState == storage.LargeFileSharesStateEnabled)
		}

		allowSharedKeyAccess := true
		if props.AllowSharedKeyAccess != nil {
			allowSharedKeyAccess = *props.AllowSharedKeyAccess
		}
		d.Set("shared_access_key_enabled", allowSharedKeyAccess)

		// Setting the encryption key type to "Service" in PUT. The following GET will not return the queue/table in the service list of its response.
		// So defaults to setting the encryption key type to "Service" if it is absent in the GET response. Also, define the default value as "Service" in the schema.
		var (
			queueEncryptionKeyType = string(storage.KeyTypeService)
			tableEncryptionKeyType = string(storage.KeyTypeService)
		)
		if encryption := props.Encryption; encryption != nil && encryption.Services != nil {
			if encryption.Services.Queue != nil {
				queueEncryptionKeyType = string(encryption.Services.Queue.KeyType)
			}
			if encryption.Services.Table != nil {
				tableEncryptionKeyType = string(encryption.Services.Table.KeyType)
			}
		}
		d.Set("table_encryption_key_type", tableEncryptionKeyType)
		d.Set("queue_encryption_key_type", queueEncryptionKeyType)

		customerManagedKey, err := flattenStorageAccountCustomerManagedKey(id, props.Encryption)
		if err != nil {
			return err
		}

		if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
			return fmt.Errorf("setting `customer_managed_key`: %+v", err)
		}

		infrastructureEncryption := false
		if encryption := props.Encryption; encryption != nil && encryption.RequireInfrastructureEncryption != nil {
			infrastructureEncryption = *encryption.RequireInfrastructureEncryption
		}
		d.Set("infrastructure_encryption_enabled", infrastructureEncryption)
	}

	if accessKeys := keys.Keys; accessKeys != nil {
		storageAccountKeys := *accessKeys
		d.Set("primary_access_key", storageAccountKeys[0].Value)
		d.Set("secondary_access_key", storageAccountKeys[1].Value)
	}

	identity, err := flattenAzureRmStorageAccountIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	storageClient := meta.(*clients.Client).Storage
	account, err := storageClient.FindAccount(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Account %q: %s", id.Name, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.Name)
	}

	blobClient := storageClient.BlobServicesClient

	// FileStorage does not support blob settings
	if resp.Kind != storage.KindFileStorage {
		blobProps, err := blobClient.GetServiceProperties(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(blobProps.Response) {
				return fmt.Errorf("reading blob properties for AzureRM Storage Account %q: %+v", id.Name, err)
			}
		}

		if err := d.Set("blob_properties", flattenBlobProperties(blobProps)); err != nil {
			return fmt.Errorf("setting `blob_properties `for AzureRM Storage Account %q: %+v", id.Name, err)
		}
	}

	fileServiceClient := storageClient.FileServicesClient

	// FileStorage does not support blob kind, FileStorage Premium is supported
	if resp.Kind == storage.KindFileStorage || resp.Kind != storage.KindBlobStorage && resp.Kind != storage.KindBlockBlobStorage && resp.Sku != nil && resp.Sku.Tier != storage.SkuTierPremium {
		shareProps, err := fileServiceClient.GetServiceProperties(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(shareProps.Response) {
				return fmt.Errorf("reading share properties for AzureRM Storage Account %q: %+v", id.Name, err)
			}
		}

		if err := d.Set("share_properties", flattenShareProperties(shareProps)); err != nil {
			return fmt.Errorf("setting `share_properties `for AzureRM Storage Account %q: %+v", id.Name, err)
		}
	}

	// queue is only available for certain tier and kind (as specified below)
	if resp.Sku == nil {
		return fmt.Errorf("retrieving %s: `sku` was nil", *id)
	}

	if resp.Sku.Tier == storage.SkuTierStandard {
		if resp.Kind == storage.KindStorage || resp.Kind == storage.KindStorageV2 {
			queueClient, err := storageClient.QueuesClient(ctx, *account)
			if err != nil {
				return fmt.Errorf("building Queues Client: %s", err)
			}

			queueProps, err := queueClient.GetServiceProperties(ctx, account.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("reading queue properties for AzureRM Storage Account %q: %+v", id.Name, err)
			}

			if err := d.Set("queue_properties", flattenQueueProperties(queueProps)); err != nil {
				return fmt.Errorf("setting `queue_properties`: %+v", err)
			}
		}
	}

	var staticWebsite []interface{}

	// static website only supported on StorageV2 and BlockBlobStorage
	if resp.Kind == storage.KindStorageV2 || resp.Kind == storage.KindBlockBlobStorage {
		storageClient := meta.(*clients.Client).Storage

		account, err := storageClient.FindAccount(ctx, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Account %q: %s", id.Name, err)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps, err := accountsClient.GetServiceProperties(ctx, id.Name)
		if err != nil {
			if staticWebsiteProps.Response.Response != nil && !utils.ResponseWasNotFound(staticWebsiteProps.Response) {
				return fmt.Errorf("reading static website for AzureRM Storage Account %q: %+v", id.Name, err)
			}
		}

		staticWebsite = flattenStaticWebsiteProperties(staticWebsiteProps)
	}

	if err := d.Set("static_website", staticWebsite); err != nil {
		return fmt.Errorf("setting `static_website `for AzureRM Storage Account %q: %+v", id.Name, err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, storageAccountResourceName)
	defer locks.UnlockByName(id.Name, storageAccountResourceName)

	read, err := client.GetProperties(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// the networking api's only allow a single change to be made to a network layout at once, so let's lock to handle that
	virtualNetworkNames := make([]string, 0)
	if props := read.AccountProperties; props != nil {
		if rules := props.NetworkRuleSet; rules != nil {
			if vnr := rules.VirtualNetworkRules; vnr != nil {
				for _, v := range *vnr {
					if v.VirtualNetworkResourceID == nil {
						continue
					}

					id, err2 := vnetParse.SubnetID(*v.VirtualNetworkResourceID)
					if err2 != nil {
						return err2
					}

					networkName := id.VirtualNetworkName
					for _, virtualNetworkName := range virtualNetworkNames {
						if networkName == virtualNetworkName {
							continue
						}
					}
					virtualNetworkNames = append(virtualNetworkNames, networkName)
				}
			}
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("issuing delete request for %s: %+v", *id, err)
		}
	}

	// remove this from the cache
	storageClient.RemoveAccountFromCache(id.Name)

	return nil
}

func expandStorageAccountCustomDomain(d *pluginsdk.ResourceData) *storage.CustomDomain {
	domains := d.Get("custom_domain").([]interface{})
	if len(domains) == 0 {
		return &storage.CustomDomain{
			Name: utils.String(""),
		}
	}

	domain := domains[0].(map[string]interface{})
	name := domain["name"].(string)
	useSubDomain := domain["use_subdomain"].(bool)
	return &storage.CustomDomain{
		Name:             utils.String(name),
		UseSubDomainName: utils.Bool(useSubDomain),
	}
}

func flattenStorageAccountCustomDomain(input *storage.CustomDomain) []interface{} {
	domain := make(map[string]interface{})

	if v := input.Name; v != nil {
		domain["name"] = *v
	}

	// use_subdomain isn't returned
	return []interface{}{domain}
}

func expandStorageAccountCustomerManagedKey(ctx context.Context, keyVaultClient *keyvault.Client, resourceClient *resource.Client, input []interface{}) (*storage.Encryption, error) {
	if len(input) == 0 {
		return &storage.Encryption{}, nil
	}

	v := input[0].(map[string]interface{})

	keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(v["key_vault_key_id"].(string))
	if err != nil {
		return nil, err
	}

	keyVaultIdRaw, err := keyVaultClient.KeyVaultIDFromBaseUrl(ctx, resourceClient, keyId.KeyVaultBaseUrl)
	if err != nil {
		return nil, err
	}
	if keyVaultIdRaw == nil {
		return nil, fmt.Errorf("unexpected nil Key Vault ID retrieved at URL %s", keyId.KeyVaultBaseUrl)
	}
	keyVaultId, err := keyVaultParse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return nil, err
	}

	vaultsClient := keyVaultClient.VaultsClient
	if keyVaultId.SubscriptionId != vaultsClient.SubscriptionID {
		vaultsClient = keyVaultClient.KeyVaultClientForSubscription(keyVaultId.SubscriptionId)
	}

	keyVault, err := vaultsClient.Get(ctx, keyVaultId.ResourceGroup, keyVaultId.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *keyVaultId, err)
	}

	softDeleteEnabled := false
	purgeProtectionEnabled := false
	if props := keyVault.Properties; props != nil {
		if esd := props.EnableSoftDelete; esd != nil {
			softDeleteEnabled = *esd
		}
		if epp := props.EnablePurgeProtection; epp != nil {
			purgeProtectionEnabled = *epp
		}
	}
	if !softDeleteEnabled || !purgeProtectionEnabled {
		return nil, fmt.Errorf("%s must be configured for both Purge Protection and Soft Delete", *keyVaultId)
	}

	encryption := &storage.Encryption{
		Services: &storage.EncryptionServices{
			Blob: &storage.EncryptionService{
				Enabled: utils.Bool(true),
				KeyType: storage.KeyTypeAccount,
			},
			File: &storage.EncryptionService{
				Enabled: utils.Bool(true),
				KeyType: storage.KeyTypeAccount,
			},
		},
		EncryptionIdentity: &storage.EncryptionIdentity{
			EncryptionUserAssignedIdentity: utils.String(v["user_assigned_identity_id"].(string)),
		},
		KeySource: storage.KeySourceMicrosoftKeyvault,
		KeyVaultProperties: &storage.KeyVaultProperties{
			KeyName:     utils.String(keyId.Name),
			KeyVersion:  utils.String(keyId.Version),
			KeyVaultURI: utils.String(keyId.KeyVaultBaseUrl),
		},
	}

	return encryption, nil
}

func flattenStorageAccountCustomerManagedKey(storageAccountId *parse.StorageAccountId, input *storage.Encryption) ([]interface{}, error) {
	if input == nil || input.KeySource == storage.KeySourceMicrosoftStorage {
		return make([]interface{}, 0), nil
	}

	userAssignedIdentityId := ""
	keyName := ""
	keyVaultURI := ""
	keyVersion := ""

	if props := input.EncryptionIdentity; props != nil {
		if props.EncryptionUserAssignedIdentity != nil {
			userAssignedIdentityId = *props.EncryptionUserAssignedIdentity
		}
	}

	if props := input.KeyVaultProperties; props != nil {
		if props.KeyName != nil {
			keyName = *props.KeyName
		}
		if props.KeyVaultURI != nil {
			keyVaultURI = *props.KeyVaultURI
		}
		if props.KeyVersion != nil {
			keyVersion = *props.KeyVersion
		}
	}

	if keyVaultURI == "" {
		return nil, fmt.Errorf("retrieving %s: `properties.encryption.keyVaultProperties.keyVaultURI` was nil", *storageAccountId)
	}

	keyId, err := keyVaultParse.NewNestedItemID(keyVaultURI, "keys", keyName, keyVersion)
	if err != nil {
		return nil, err
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_key_id":          keyId.ID(),
			"user_assigned_identity_id": userAssignedIdentityId,
		},
	}, nil
}

func expandArmStorageAccountAzureFilesAuthentication(input []interface{}) (*storage.AzureFilesIdentityBasedAuthentication, error) {
	if len(input) == 0 {
		return &storage.AzureFilesIdentityBasedAuthentication{
			DirectoryServiceOptions: storage.DirectoryServiceOptionsNone,
		}, nil
	}

	v := input[0].(map[string]interface{})

	directoryOption := storage.DirectoryServiceOptions(v["directory_type"].(string))
	if _, ok := v["active_directory"]; directoryOption == storage.DirectoryServiceOptionsAD && !ok {
		return nil, fmt.Errorf("`active_directory` is required when `directory_type` is `AD`")
	}

	return &storage.AzureFilesIdentityBasedAuthentication{
		DirectoryServiceOptions:   directoryOption,
		ActiveDirectoryProperties: expandArmStorageAccountActiveDirectoryProperties(v["active_directory"].([]interface{})),
	}, nil
}

func expandArmStorageAccountActiveDirectoryProperties(input []interface{}) *storage.ActiveDirectoryProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &storage.ActiveDirectoryProperties{
		AzureStorageSid:   utils.String(v["storage_sid"].(string)),
		DomainGUID:        utils.String(v["domain_guid"].(string)),
		DomainName:        utils.String(v["domain_name"].(string)),
		DomainSid:         utils.String(v["domain_sid"].(string)),
		ForestName:        utils.String(v["forest_name"].(string)),
		NetBiosDomainName: utils.String(v["netbios_domain_name"].(string)),
	}
}

func expandArmStorageAccountRouting(input []interface{}) *storage.RoutingPreference {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &storage.RoutingPreference{
		RoutingChoice:             storage.RoutingChoice(v["choice"].(string)),
		PublishMicrosoftEndpoints: utils.Bool(v["publish_microsoft_endpoints"].(bool)),
		PublishInternetEndpoints:  utils.Bool(v["publish_internet_endpoints"].(bool)),
	}
}

func expandStorageAccountNetworkRules(d *pluginsdk.ResourceData, tenantId string) *storage.NetworkRuleSet {
	networkRules := d.Get("network_rules").([]interface{})
	if len(networkRules) == 0 {
		// Default access is enabled when no network rules are set.
		return &storage.NetworkRuleSet{DefaultAction: storage.DefaultActionAllow}
	}

	networkRule := networkRules[0].(map[string]interface{})
	networkRuleSet := &storage.NetworkRuleSet{
		IPRules:             expandStorageAccountIPRules(networkRule),
		VirtualNetworkRules: expandStorageAccountVirtualNetworks(networkRule),
		Bypass:              expandStorageAccountBypass(networkRule),
		ResourceAccessRules: expandStorageAccountPrivateLinkAccess(networkRule["private_link_access"].([]interface{}), tenantId),
	}

	if v := networkRule["default_action"]; v != nil {
		networkRuleSet.DefaultAction = storage.DefaultAction(v.(string))
	}

	return networkRuleSet
}

func expandStorageAccountIPRules(networkRule map[string]interface{}) *[]storage.IPRule {
	ipRulesInfo := networkRule["ip_rules"].(*pluginsdk.Set).List()
	ipRules := make([]storage.IPRule, len(ipRulesInfo))

	for i, ipRuleConfig := range ipRulesInfo {
		attrs := ipRuleConfig.(string)
		ipRule := storage.IPRule{
			IPAddressOrRange: utils.String(attrs),
			Action:           storage.ActionAllow,
		}
		ipRules[i] = ipRule
	}

	return &ipRules
}

func expandStorageAccountVirtualNetworks(networkRule map[string]interface{}) *[]storage.VirtualNetworkRule {
	virtualNetworkInfo := networkRule["virtual_network_subnet_ids"].(*pluginsdk.Set).List()
	virtualNetworks := make([]storage.VirtualNetworkRule, len(virtualNetworkInfo))

	for i, virtualNetworkConfig := range virtualNetworkInfo {
		attrs := virtualNetworkConfig.(string)
		virtualNetwork := storage.VirtualNetworkRule{
			VirtualNetworkResourceID: utils.String(attrs),
			Action:                   storage.ActionAllow,
		}
		virtualNetworks[i] = virtualNetwork
	}

	return &virtualNetworks
}

func expandStorageAccountBypass(networkRule map[string]interface{}) storage.Bypass {
	bypassInfo := networkRule["bypass"].(*pluginsdk.Set).List()

	var bypassValues []string
	for _, bypassConfig := range bypassInfo {
		bypassValues = append(bypassValues, bypassConfig.(string))
	}

	return storage.Bypass(strings.Join(bypassValues, ", "))
}

func expandStorageAccountPrivateLinkAccess(inputs []interface{}, tenantId string) *[]storage.ResourceAccessRule {
	privateLinkAccess := make([]storage.ResourceAccessRule, 0)
	if len(inputs) == 0 {
		return &privateLinkAccess
	}
	for _, input := range inputs {
		accessRule := input.(map[string]interface{})
		if v := accessRule["endpoint_tenant_id"].(string); v != "" {
			tenantId = v
		}
		privateLinkAccess = append(privateLinkAccess, storage.ResourceAccessRule{
			TenantID:   utils.String(tenantId),
			ResourceID: utils.String(accessRule["endpoint_resource_id"].(string)),
		})
	}

	return &privateLinkAccess
}

func expandBlobProperties(input []interface{}) *storage.BlobServiceProperties {
	props := storage.BlobServiceProperties{
		BlobServicePropertiesProperties: &storage.BlobServicePropertiesProperties{
			Cors: &storage.CorsRules{
				CorsRules: &[]storage.CorsRule{},
			},
			IsVersioningEnabled: utils.Bool(false),
			ChangeFeed: &storage.ChangeFeed{
				Enabled: utils.Bool(false),
			},
			DeleteRetentionPolicy: &storage.DeleteRetentionPolicy{
				Enabled: utils.Bool(false),
			},
		},
	}

	if len(input) == 0 || input[0] == nil {
		return &props
	}

	v := input[0].(map[string]interface{})

	deletePolicyRaw := v["delete_retention_policy"].([]interface{})
	props.BlobServicePropertiesProperties.DeleteRetentionPolicy = expandBlobPropertiesDeleteRetentionPolicy(deletePolicyRaw, true)
	corsRaw := v["cors_rule"].([]interface{})
	props.BlobServicePropertiesProperties.Cors = expandBlobPropertiesCors(corsRaw)

	props.IsVersioningEnabled = utils.Bool(v["versioning_enabled"].(bool))

	props.ChangeFeed = &storage.ChangeFeed{
		Enabled: utils.Bool(v["change_feed_enabled"].(bool)),
	}

	if version, ok := v["default_service_version"].(string); ok && version != "" {
		props.DefaultServiceVersion = utils.String(version)
	}

	return &props
}

func expandBlobPropertiesDeleteRetentionPolicy(input []interface{}, isupdate bool) *storage.DeleteRetentionPolicy {
	result := storage.DeleteRetentionPolicy{
		Enabled: utils.Bool(false),
	}
	if (len(input) == 0 || input[0] == nil) && !isupdate {
		return nil
	}

	if (len(input) == 0 || input[0] == nil) && isupdate {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &storage.DeleteRetentionPolicy{
		Enabled: utils.Bool(true),
		Days:    utils.Int32(int32(policy["days"].(int))),
	}
}

func expandBlobPropertiesCors(input []interface{}) *storage.CorsRules {
	blobCorsRules := storage.CorsRules{}

	if len(input) == 0 {
		return &blobCorsRules
	}

	corsRules := make([]storage.CorsRule, 0)
	for _, attr := range input {
		corsRuleAttr := attr.(map[string]interface{})
		corsRule := storage.CorsRule{}

		allowedOrigins := *utils.ExpandStringSlice(corsRuleAttr["allowed_origins"].([]interface{}))
		allowedHeaders := *utils.ExpandStringSlice(corsRuleAttr["allowed_headers"].([]interface{}))
		allowedMethods := *utils.ExpandStringSlice(corsRuleAttr["allowed_methods"].([]interface{}))
		exposedHeaders := *utils.ExpandStringSlice(corsRuleAttr["exposed_headers"].([]interface{}))
		maxAgeInSeconds := int32(corsRuleAttr["max_age_in_seconds"].(int))

		corsRule.AllowedOrigins = &allowedOrigins
		corsRule.AllowedHeaders = &allowedHeaders
		corsRule.AllowedMethods = &allowedMethods
		corsRule.ExposedHeaders = &exposedHeaders
		corsRule.MaxAgeInSeconds = &maxAgeInSeconds

		corsRules = append(corsRules, corsRule)
	}

	blobCorsRules.CorsRules = &corsRules

	return &blobCorsRules
}

func expandShareProperties(input []interface{}) storage.FileServiceProperties {
	props := storage.FileServiceProperties{
		FileServicePropertiesProperties: &storage.FileServicePropertiesProperties{
			Cors: &storage.CorsRules{
				CorsRules: &[]storage.CorsRule{},
			},
			ShareDeleteRetentionPolicy: &storage.DeleteRetentionPolicy{
				Enabled: utils.Bool(false),
			},
		},
	}

	if len(input) == 0 || input[0] == nil {
		return props
	}

	v := input[0].(map[string]interface{})

	props.FileServicePropertiesProperties.ShareDeleteRetentionPolicy = expandBlobPropertiesDeleteRetentionPolicy(v["retention_policy"].([]interface{}), false)

	props.FileServicePropertiesProperties.Cors = expandBlobPropertiesCors(v["cors_rule"].([]interface{}))

	props.ProtocolSettings = &storage.ProtocolSettings{
		Smb: expandSharePropertiesSMB(v["smb"].([]interface{})),
	}

	return props
}

func expandSharePropertiesSMB(input []interface{}) *storage.SmbSetting {
	if len(input) == 0 || input[0] == nil {
		return &storage.SmbSetting{
			Versions:                 utils.String(""),
			AuthenticationMethods:    utils.String(""),
			KerberosTicketEncryption: utils.String(""),
			ChannelEncryption:        utils.String(""),
		}
	}

	v := input[0].(map[string]interface{})

	return &storage.SmbSetting{
		Versions:                 utils.ExpandStringSliceWithDelimiter(v["versions"].(*pluginsdk.Set).List(), ";"),
		AuthenticationMethods:    utils.ExpandStringSliceWithDelimiter(v["authentication_types"].(*pluginsdk.Set).List(), ";"),
		KerberosTicketEncryption: utils.ExpandStringSliceWithDelimiter(v["kerberos_ticket_encryption_type"].(*pluginsdk.Set).List(), ";"),
		ChannelEncryption:        utils.ExpandStringSliceWithDelimiter(v["channel_encryption_type"].(*pluginsdk.Set).List(), ";"),
	}
}

func expandQueueProperties(input []interface{}) (queues.StorageServiceProperties, error) {
	var err error
	properties := queues.StorageServiceProperties{
		Cors: &queues.Cors{
			CorsRule: []queues.CorsRule{},
		},
		HourMetrics: &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
		MinuteMetrics: &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
		Logging: &queues.LoggingConfig{
			Version: "1.0",
			Delete:  false,
			Read:    false,
			Write:   false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
	}
	if len(input) == 0 {
		return properties, nil
	}

	attrs := input[0].(map[string]interface{})

	properties.Cors = expandQueuePropertiesCors(attrs["cors_rule"].([]interface{}))
	properties.Logging = expandQueuePropertiesLogging(attrs["logging"].([]interface{}))
	properties.MinuteMetrics, err = expandQueuePropertiesMetrics(attrs["minute_metrics"].([]interface{}))
	if err != nil {
		return properties, fmt.Errorf("expanding `minute_metrics`: %+v", err)
	}
	properties.HourMetrics, err = expandQueuePropertiesMetrics(attrs["hour_metrics"].([]interface{}))
	if err != nil {
		return properties, fmt.Errorf("expanding `hour_metrics`: %+v", err)
	}

	return properties, nil
}

func expandQueuePropertiesMetrics(input []interface{}) (*queues.MetricsConfig, error) {
	if len(input) == 0 {
		return &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		}, nil
	}

	metricsAttr := input[0].(map[string]interface{})

	metrics := &queues.MetricsConfig{
		Version: metricsAttr["version"].(string),
		Enabled: metricsAttr["enabled"].(bool),
	}

	if v, ok := metricsAttr["retention_policy_days"]; ok {
		if days := v.(int); days > 0 {
			metrics.RetentionPolicy = queues.RetentionPolicy{
				Days:    days,
				Enabled: true,
			}
		}
	}

	if v, ok := metricsAttr["include_apis"]; ok {
		includeAPIs := v.(bool)
		if metrics.Enabled {
			metrics.IncludeAPIs = &includeAPIs
		} else if includeAPIs {
			return nil, fmt.Errorf("`include_apis` may only be set when `enabled` is true")
		}
	}

	return metrics, nil
}

func expandQueuePropertiesLogging(input []interface{}) *queues.LoggingConfig {
	if len(input) == 0 {
		return &queues.LoggingConfig{
			Version: "1.0",
			Delete:  false,
			Read:    false,
			Write:   false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		}
	}

	loggingAttr := input[0].(map[string]interface{})
	logging := &queues.LoggingConfig{
		Version: loggingAttr["version"].(string),
		Delete:  loggingAttr["delete"].(bool),
		Read:    loggingAttr["read"].(bool),
		Write:   loggingAttr["write"].(bool),
	}

	if v, ok := loggingAttr["retention_policy_days"]; ok {
		if days := v.(int); days > 0 {
			logging.RetentionPolicy = queues.RetentionPolicy{
				Days:    days,
				Enabled: true,
			}
		}
	}

	return logging
}

func expandQueuePropertiesCors(input []interface{}) *queues.Cors {
	if len(input) == 0 {
		return &queues.Cors{}
	}

	corsRules := make([]queues.CorsRule, 0)
	for _, attr := range input {
		corsRuleAttr := attr.(map[string]interface{})
		corsRule := queues.CorsRule{}

		corsRule.AllowedOrigins = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_origins"].([]interface{})), ",")
		corsRule.ExposedHeaders = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["exposed_headers"].([]interface{})), ",")
		corsRule.AllowedHeaders = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_headers"].([]interface{})), ",")
		corsRule.AllowedMethods = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_methods"].([]interface{})), ",")
		corsRule.MaxAgeInSeconds = corsRuleAttr["max_age_in_seconds"].(int)

		corsRules = append(corsRules, corsRule)
	}

	cors := &queues.Cors{
		CorsRule: corsRules,
	}
	return cors
}

func expandStaticWebsiteProperties(input []interface{}) accounts.StorageServiceProperties {
	properties := accounts.StorageServiceProperties{
		StaticWebsite: &accounts.StaticWebsite{
			Enabled: false,
		},
	}
	if len(input) == 0 {
		return properties
	}

	properties.StaticWebsite.Enabled = true

	// @tombuildsstuff: this looks weird, doesn't it?
	// Since the presence of this block signifies the website's enabled however all fields within it are optional
	// TF Core returns a nil object when there's no keys defined within the block, rather than an empty map. As
	// such this hack allows us to have a Static Website block with only Enabled configured, without the optional
	// inner properties.
	if val := input[0]; val != nil {
		attr := val.(map[string]interface{})
		if v, ok := attr["index_document"]; ok {
			properties.StaticWebsite.IndexDocument = v.(string)
		}

		if v, ok := attr["error_404_document"]; ok {
			properties.StaticWebsite.ErrorDocument404Path = v.(string)
		}
	}

	return properties
}

func flattenArmStorageAccountAzureFilesAuthentication(input *storage.AzureFilesIdentityBasedAuthentication) []interface{} {
	if input == nil || input.DirectoryServiceOptions == storage.DirectoryServiceOptionsNone {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"directory_type":   input.DirectoryServiceOptions,
			"active_directory": flattenArmStorageAccountActiveDirectoryProperties(input.ActiveDirectoryProperties),
		},
	}
}

func flattenArmStorageAccountActiveDirectoryProperties(input *storage.ActiveDirectoryProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var azureStorageSid string
	if input.AzureStorageSid != nil {
		azureStorageSid = *input.AzureStorageSid
	}
	var domainGuid string
	if input.DomainGUID != nil {
		domainGuid = *input.DomainGUID
	}
	var domainName string
	if input.DomainName != nil {
		domainName = *input.DomainName
	}
	var domainSid string
	if input.DomainSid != nil {
		domainSid = *input.DomainSid
	}
	var forestName string
	if input.ForestName != nil {
		forestName = *input.ForestName
	}
	var netBiosDomainName string
	if input.NetBiosDomainName != nil {
		netBiosDomainName = *input.NetBiosDomainName
	}
	return []interface{}{
		map[string]interface{}{
			"storage_sid":         azureStorageSid,
			"domain_guid":         domainGuid,
			"domain_name":         domainName,
			"domain_sid":          domainSid,
			"forest_name":         forestName,
			"netbios_domain_name": netBiosDomainName,
		},
	}
}

func flattenArmStorageAccountRouting(input *storage.RoutingPreference) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var publishInternetEndpoints bool
	if input.PublishInternetEndpoints != nil {
		publishInternetEndpoints = *input.PublishInternetEndpoints
	}
	var publishMicrosoftEndpoints bool
	if input.PublishMicrosoftEndpoints != nil {
		publishMicrosoftEndpoints = *input.PublishMicrosoftEndpoints
	}

	return []interface{}{
		map[string]interface{}{
			"publish_internet_endpoints":  publishInternetEndpoints,
			"publish_microsoft_endpoints": publishMicrosoftEndpoints,
			"choice":                      input.RoutingChoice,
		},
	}
}

func flattenStorageAccountNetworkRules(input *storage.NetworkRuleSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	networkRules := make(map[string]interface{})

	networkRules["ip_rules"] = pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountIPRules(input.IPRules))
	networkRules["virtual_network_subnet_ids"] = pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountVirtualNetworks(input.VirtualNetworkRules))
	networkRules["bypass"] = pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountBypass(input.Bypass))
	networkRules["default_action"] = string(input.DefaultAction)
	networkRules["private_link_access"] = flattenStorageAccountPrivateLinkAccess(input.ResourceAccessRules)

	return []interface{}{networkRules}
}

func flattenStorageAccountIPRules(input *[]storage.IPRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ipRules := make([]interface{}, 0)
	for _, ipRule := range *input {
		if ipRule.IPAddressOrRange == nil {
			continue
		}

		ipRules = append(ipRules, *ipRule.IPAddressOrRange)
	}

	return ipRules
}

func flattenStorageAccountVirtualNetworks(input *[]storage.VirtualNetworkRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	virtualNetworks := make([]interface{}, 0)
	for _, virtualNetwork := range *input {
		if virtualNetwork.VirtualNetworkResourceID == nil {
			continue
		}

		virtualNetworks = append(virtualNetworks, *virtualNetwork.VirtualNetworkResourceID)
	}

	return virtualNetworks
}

func flattenStorageAccountPrivateLinkAccess(inputs *[]storage.ResourceAccessRule) []interface{} {
	if inputs == nil || len(*inputs) == 0 {
		return []interface{}{}
	}

	accessRules := make([]interface{}, 0)
	for _, input := range *inputs {
		var resourceId, tenantId string
		if input.ResourceID != nil {
			resourceId = *input.ResourceID
		}

		if input.TenantID != nil {
			tenantId = *input.TenantID
		}

		accessRules = append(accessRules, map[string]interface{}{
			"endpoint_resource_id": resourceId,
			"endpoint_tenant_id":   tenantId,
		})
	}

	return accessRules
}

func flattenBlobProperties(input storage.BlobServiceProperties) []interface{} {
	if input.BlobServicePropertiesProperties == nil {
		return []interface{}{}
	}

	flattenedCorsRules := make([]interface{}, 0)
	if corsRules := input.BlobServicePropertiesProperties.Cors; corsRules != nil {
		flattenedCorsRules = flattenBlobPropertiesCorsRule(corsRules)
	}

	flattenedDeletePolicy := make([]interface{}, 0)
	if deletePolicy := input.BlobServicePropertiesProperties.DeleteRetentionPolicy; deletePolicy != nil {
		flattenedDeletePolicy = flattenBlobPropertiesDeleteRetentionPolicy(deletePolicy)
	}

	flattenedContainerDeletePolicy := make([]interface{}, 0)
	if containerDeletePolicy := input.BlobServicePropertiesProperties.ContainerDeleteRetentionPolicy; containerDeletePolicy != nil {
		flattenedContainerDeletePolicy = flattenBlobPropertiesDeleteRetentionPolicy(containerDeletePolicy)
	}

	versioning, changeFeed := false, false
	if input.BlobServicePropertiesProperties.IsVersioningEnabled != nil {
		versioning = *input.BlobServicePropertiesProperties.IsVersioningEnabled
	}

	if v := input.BlobServicePropertiesProperties.ChangeFeed; v != nil && v.Enabled != nil {
		changeFeed = *v.Enabled
	}

	var defaultServiceVersion string
	if input.BlobServicePropertiesProperties.DefaultServiceVersion != nil {
		defaultServiceVersion = *input.BlobServicePropertiesProperties.DefaultServiceVersion
	}

	var LastAccessTimeTrackingPolicy bool
	if v := input.BlobServicePropertiesProperties.LastAccessTimeTrackingPolicy; v != nil && v.Enable != nil {
		LastAccessTimeTrackingPolicy = *v.Enable
	}

	return []interface{}{
		map[string]interface{}{
			"cors_rule":                         flattenedCorsRules,
			"delete_retention_policy":           flattenedDeletePolicy,
			"versioning_enabled":                versioning,
			"change_feed_enabled":               changeFeed,
			"default_service_version":           defaultServiceVersion,
			"last_access_time_enabled":          LastAccessTimeTrackingPolicy,
			"container_delete_retention_policy": flattenedContainerDeletePolicy,
		},
	}
}

func flattenBlobPropertiesCorsRule(input *storage.CorsRules) []interface{} {
	corsRules := make([]interface{}, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		allowedOrigins := make([]string, 0)
		if corsRule.AllowedOrigins != nil {
			allowedOrigins = *corsRule.AllowedOrigins
		}

		allowedMethods := make([]string, 0)
		if corsRule.AllowedMethods != nil {
			allowedMethods = *corsRule.AllowedMethods
		}

		allowedHeaders := make([]string, 0)
		if corsRule.AllowedHeaders != nil {
			allowedHeaders = *corsRule.AllowedHeaders
		}

		exposedHeaders := make([]string, 0)
		if corsRule.ExposedHeaders != nil {
			exposedHeaders = *corsRule.ExposedHeaders
		}

		maxAgeInSeconds := 0
		if corsRule.MaxAgeInSeconds != nil {
			maxAgeInSeconds = int(*corsRule.MaxAgeInSeconds)
		}

		corsRules = append(corsRules, map[string]interface{}{
			"allowed_headers":    allowedHeaders,
			"allowed_origins":    allowedOrigins,
			"allowed_methods":    allowedMethods,
			"exposed_headers":    exposedHeaders,
			"max_age_in_seconds": maxAgeInSeconds,
		})
	}

	return corsRules
}

func flattenBlobPropertiesDeleteRetentionPolicy(input *storage.DeleteRetentionPolicy) []interface{} {
	deleteRetentionPolicy := make([]interface{}, 0)

	if input == nil {
		return deleteRetentionPolicy
	}

	if enabled := input.Enabled; enabled != nil && *enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		deleteRetentionPolicy = append(deleteRetentionPolicy, map[string]interface{}{
			"days": days,
		})
	}

	return deleteRetentionPolicy
}

func flattenQueueProperties(input *queues.StorageServiceProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	queueProperties := make(map[string]interface{})

	if cors := input.Cors; cors != nil {
		if len(cors.CorsRule) > 0 {
			if cors.CorsRule[0].AllowedOrigins != "" {
				queueProperties["cors_rule"] = flattenQueuePropertiesCorsRule(input.Cors.CorsRule)
			}
		}
	}

	if logging := input.Logging; logging != nil {
		if logging.Version != "" {
			queueProperties["logging"] = flattenQueuePropertiesLogging(*logging)
		}
	}

	if hourMetrics := input.HourMetrics; hourMetrics != nil {
		if hourMetrics.Version != "" {
			queueProperties["hour_metrics"] = flattenQueuePropertiesMetrics(*hourMetrics)
		}
	}

	if minuteMetrics := input.MinuteMetrics; minuteMetrics != nil {
		if minuteMetrics.Version != "" {
			queueProperties["minute_metrics"] = flattenQueuePropertiesMetrics(*minuteMetrics)
		}
	}

	if len(queueProperties) == 0 {
		return []interface{}{}
	}
	return []interface{}{queueProperties}
}

func flattenQueuePropertiesMetrics(input queues.MetricsConfig) []interface{} {
	metrics := make(map[string]interface{})

	metrics["version"] = input.Version
	metrics["enabled"] = input.Enabled

	if input.IncludeAPIs != nil {
		metrics["include_apis"] = *input.IncludeAPIs
	}

	if input.RetentionPolicy.Enabled {
		metrics["retention_policy_days"] = input.RetentionPolicy.Days
	}

	return []interface{}{metrics}
}

func flattenQueuePropertiesCorsRule(input []queues.CorsRule) []interface{} {
	corsRules := make([]interface{}, 0)

	for _, corsRule := range input {
		attr := make(map[string]interface{})

		attr["allowed_origins"] = flattenCorsProperty(corsRule.AllowedOrigins)
		attr["allowed_methods"] = flattenCorsProperty(corsRule.AllowedMethods)
		attr["allowed_headers"] = flattenCorsProperty(corsRule.AllowedHeaders)
		attr["exposed_headers"] = flattenCorsProperty(corsRule.ExposedHeaders)
		attr["max_age_in_seconds"] = corsRule.MaxAgeInSeconds

		corsRules = append(corsRules, attr)
	}

	return corsRules
}

func flattenQueuePropertiesLogging(input queues.LoggingConfig) []interface{} {
	logging := make(map[string]interface{})

	logging["version"] = input.Version
	logging["delete"] = input.Delete
	logging["read"] = input.Read
	logging["write"] = input.Write

	if input.RetentionPolicy.Enabled {
		logging["retention_policy_days"] = input.RetentionPolicy.Days
	}

	return []interface{}{logging}
}

func flattenCorsProperty(input string) []interface{} {
	results := make([]interface{}, 0, len(input))

	origins := strings.Split(input, ",")
	for _, origin := range origins {
		results = append(results, origin)
	}

	return results
}

func flattenShareProperties(input storage.FileServiceProperties) []interface{} {
	if input.FileServicePropertiesProperties == nil {
		return []interface{}{}
	}

	flattenedCorsRules := make([]interface{}, 0)
	if corsRules := input.FileServicePropertiesProperties.Cors; corsRules != nil {
		flattenedCorsRules = flattenBlobPropertiesCorsRule(corsRules)
	}

	flattenedDeletePolicy := make([]interface{}, 0)
	if deletePolicy := input.FileServicePropertiesProperties.ShareDeleteRetentionPolicy; deletePolicy != nil {
		flattenedDeletePolicy = flattenBlobPropertiesDeleteRetentionPolicy(deletePolicy)
	}

	flattenedSMB := make([]interface{}, 0)
	if protocol := input.FileServicePropertiesProperties.ProtocolSettings; protocol != nil {
		flattenedSMB = flattenedSharePropertiesSMB(protocol.Smb)
	}

	if len(flattenedCorsRules) == 0 && len(flattenedDeletePolicy) == 0 && len(flattenedSMB) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"cors_rule":        flattenedCorsRules,
			"retention_policy": flattenedDeletePolicy,
			"smb":              flattenedSMB,
		},
	}
}

func flattenedSharePropertiesSMB(input *storage.SmbSetting) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	versions := []interface{}{}
	if input.Versions != nil {
		versions = utils.FlattenStringSliceWithDelimiter(input.Versions, ";")
	}

	authenticationMethods := []interface{}{}
	if input.AuthenticationMethods != nil {
		authenticationMethods = utils.FlattenStringSliceWithDelimiter(input.AuthenticationMethods, ";")
	}

	kerberosTicketEncryption := []interface{}{}
	if input.KerberosTicketEncryption != nil {
		kerberosTicketEncryption = utils.FlattenStringSliceWithDelimiter(input.KerberosTicketEncryption, ";")
	}

	channelEncryption := []interface{}{}
	if input.ChannelEncryption != nil {
		channelEncryption = utils.FlattenStringSliceWithDelimiter(input.ChannelEncryption, ";")
	}

	if len(versions) == 0 && len(authenticationMethods) == 0 && len(kerberosTicketEncryption) == 0 && len(channelEncryption) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"versions":                        versions,
			"authentication_types":            authenticationMethods,
			"kerberos_ticket_encryption_type": kerberosTicketEncryption,
			"channel_encryption_type":         channelEncryption,
		},
	}
}

func flattenStaticWebsiteProperties(input accounts.GetServicePropertiesResult) []interface{} {
	if storageServiceProps := input.StorageServiceProperties; storageServiceProps != nil {
		if staticWebsite := storageServiceProps.StaticWebsite; staticWebsite != nil {
			if !staticWebsite.Enabled {
				return []interface{}{}
			}

			return []interface{}{
				map[string]interface{}{
					"index_document":     staticWebsite.IndexDocument,
					"error_404_document": staticWebsite.ErrorDocument404Path,
				},
			}
		}
	}
	return []interface{}{}
}

func flattenStorageAccountBypass(input storage.Bypass) []interface{} {
	bypassValues := strings.Split(string(input), ", ")
	bypass := make([]interface{}, len(bypassValues))

	for i, value := range bypassValues {
		bypass[i] = value
	}

	return bypass
}

func expandAzureRmStorageAccountIdentity(input []interface{}) (*storage.Identity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := storage.Identity{
		Type: storage.IdentityType(string(expanded.Type)),
	}

	// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
	if expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.Type = storage.IdentityTypeSystemAssignedUserAssigned
	}

	// 'Failed to perform resource identity operation. Status: 'BadRequest'. Response:
	// {"error":{"code":"BadRequest",
	//  "message":"The request format was unexpected, a non-UserAssigned identity type should not contain: userAssignedIdentities"
	// }}
	// Upstream issue: https://github.com/Azure/azure-rest-api-specs/issues/17650
	if len(expanded.IdentityIds) > 0 {
		userAssignedIdentities := make(map[string]*storage.UserAssignedIdentity)
		for id := range expanded.IdentityIds {
			userAssignedIdentities[id] = &storage.UserAssignedIdentity{}
		}
		out.UserAssignedIdentities = userAssignedIdentities
	}

	return &out, nil
}

func flattenAzureRmStorageAccountIdentity(input *storage.Identity) (*[]interface{}, error) {
	var config *identity.SystemAndUserAssignedMap

	if input != nil {
		config = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: nil,
		}

		// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
		if input.Type == storage.IdentityTypeSystemAssignedUserAssigned {
			config.Type = identity.TypeSystemAssignedUserAssigned
		}

		if input.PrincipalID != nil {
			config.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			config.TenantId = *input.TenantID
		}
		identityIds := make(map[string]identity.UserAssignedIdentityDetails)
		for k, v := range input.UserAssignedIdentities {
			if v == nil {
				continue
			}

			details := identity.UserAssignedIdentityDetails{}

			if v.ClientID != nil {
				details.ClientId = utils.String(*v.ClientID)
			}
			if v.PrincipalID != nil {
				details.PrincipalId = utils.String(*v.PrincipalID)
			}

			identityIds[k] = details
		}

		config.IdentityIds = identityIds
	}

	return identity.FlattenSystemAndUserAssignedMap(config)
}

func getBlobConnectionString(blobEndpoint *string, acctName *string, acctKey *string) string {
	var endpoint string
	if blobEndpoint != nil {
		endpoint = *blobEndpoint
	}

	var name string
	if acctName != nil {
		name = *acctName
	}

	var key string
	if acctKey != nil {
		key = *acctKey
	}

	return fmt.Sprintf("DefaultEndpointsProtocol=https;BlobEndpoint=%s;AccountName=%s;AccountKey=%s", endpoint, name, key)
}

func flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d *pluginsdk.ResourceData, primary *storage.Endpoints) error {
	if primary == nil {
		return fmt.Errorf("primary endpoints should not be empty")
	}

	if err := setEndpointAndHost(d, "primary", primary.Blob, "blob"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.Dfs, "dfs"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.File, "file"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.Queue, "queue"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.Table, "table"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.Web, "web"); err != nil {
		return err
	}

	return nil
}

func flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d *pluginsdk.ResourceData, secondary *storage.Endpoints) error {
	if secondary == nil {
		return nil
	}

	if err := setEndpointAndHost(d, "secondary", secondary.Blob, "blob"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.Dfs, "dfs"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.File, "file"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.Queue, "queue"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.Table, "table"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.Web, "web"); err != nil {
		return err
	}
	return nil
}

func setEndpointAndHost(d *pluginsdk.ResourceData, ordinalString string, endpointType *string, typeString string) error {
	var endpoint, host string
	if v := endpointType; v != nil {
		endpoint = *v

		u, err := url.Parse(*v)
		if err != nil {
			return fmt.Errorf("invalid %s endpoint for parsing: %q", typeString, *v)
		}
		host = u.Host
	}

	// lintignore: R001
	d.Set(fmt.Sprintf("%s_%s_endpoint", ordinalString, typeString), endpoint)
	// lintignore: R001
	d.Set(fmt.Sprintf("%s_%s_host", ordinalString, typeString), host)
	return nil
}

func expandEdgeZone(input string) *storage.ExtendedLocation {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &storage.ExtendedLocation{
		Name: utils.String(normalized),
		Type: storage.ExtendedLocationTypesEdgeZone,
	}
}

func flattenEdgeZone(input *storage.ExtendedLocation) string {
	if input == nil || input.Type != storage.ExtendedLocationTypesEdgeZone || input.Name == nil {
		return ""
	}
	return edgezones.NormalizeNilable(input.Name)
}
