package storage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	azautorest "github.com/Azure/go-autorest/autorest"
	autorestAzure "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/go-getter/helper/url"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	msiparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	msiValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/queue/queues"
)

var storageAccountResourceName = "azurerm_storage_account"

func resourceStorageAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountCreate,
		Read:   resourceStorageAccountRead,
		Update: resourceStorageAccountUpdate,
		Delete: resourceStorageAccountDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AccountV0ToV1{},
			1: migration.AccountV1ToV2{},
		}),

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
					string(storage.Storage),
					string(storage.BlobStorage),
					string(storage.BlockBlobStorage),
					string(storage.FileStorage),
					string(storage.StorageV2),
				}, true),
				Default: string(storage.StorageV2),
			},

			"account_tier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"account_replication_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LRS",
					"ZRS",
					"GRS",
					"RAGRS",
					"GZRS",
					"RAGZRS",
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			// Only valid for BlobStorage & StorageV2 accounts, defaults to "Hot" in create function
			"access_tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.Cool),
					string(storage.Hot),
				}, true),
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

			"enable_https_traffic_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"min_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(storage.TLS10),
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.TLS10),
					string(storage.TLS11),
					string(storage.TLS12),
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

			"allow_blob_public_access": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
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
									string(storage.AzureServices),
									string(storage.Logging),
									string(storage.Metrics),
									string(storage.None),
								}, true),
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

			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(storage.IdentityTypeSystemAssigned),
								string(storage.IdentityTypeSystemAssignedUserAssigned),
								string(storage.IdentityTypeUserAssigned),
							}, true),
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"identity_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: msiValidate.UserAssignedIdentityID,
							},
						},
					},
				},
			},

			"blob_properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cors_rule": schemaStorageAccountCorsRule(true),
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
						"cors_rule": schemaStorageAccountCorsRule(false),
						"logging": {
							Type:     pluginsdk.TypeList,
							Optional: true,
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
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"version": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
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
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"version": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
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
								string(storage.MicrosoftRouting),
								string(storage.InternetRouting),
							}, false),
							Default: string(storage.MicrosoftRouting),
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
						"cors_rule": schemaStorageAccountCorsRule(true),

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

			//lintignore:XS003
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
		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			if d.HasChange("account_kind") {
				accountKind, changedKind := d.GetChange("account_kind")

				if accountKind != string(storage.Storage) && changedKind != string(storage.StorageV2) {
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
	}
}

func resourceStorageAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	envName := meta.(*clients.Client).Account.Environment.Name
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountName := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	locks.ByName(storageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountName, storageAccountResourceName)

	existing, err := client.GetProperties(ctx, resourceGroupName, storageAccountName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Storage Account %q (Resource Group %q): %s", storageAccountName, resourceGroupName, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_storage_account", *existing.ID)
	}

	accountKind := d.Get("account_kind").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	enableHTTPSTrafficOnly := d.Get("enable_https_traffic_only").(bool)
	minimumTLSVersion := d.Get("min_tls_version").(string)
	isHnsEnabled := d.Get("is_hns_enabled").(bool)
	nfsV3Enabled := d.Get("nfsv3_enabled").(bool)
	allowBlobPublicAccess := d.Get("allow_blob_public_access").(bool)

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	// this is the default behavior for the resource if the attribute is nil
	// we are making this change in Terraform https://github.com/terraform-providers/terraform-provider-azurerm/issues/11689
	// because the portal UI team has a bug in their code ignoring the ARM API documention which state that nil is true
	// TODO: Remove code when Portal UI team fixes their code
	allowSharedKeyAccess := true

	parameters := storage.AccountCreateParameters{
		Location: &location,
		Sku: &storage.Sku{
			Name: storage.SkuName(storageType),
		},
		Tags: tags.Expand(t),
		Kind: storage.Kind(accountKind),
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
			EnableHTTPSTrafficOnly: &enableHTTPSTrafficOnly,
			NetworkRuleSet:         expandStorageAccountNetworkRules(d, tenantId),
			IsHnsEnabled:           &isHnsEnabled,
			EnableNfsV3:            &nfsV3Enabled,
			// TODO: Remove AllowSharedKeyAcces assignment when Portal UI team fixes their code (e.g. nil is true)
			AllowSharedKeyAccess: &allowSharedKeyAccess,
		},
	}

	// For all Clouds except Public, China, and USGovernmentCloud, don't specify "allow_blob_public_access" and "min_tls_version" in request body.
	// https://github.com/terraform-providers/terraform-provider-azurerm/issues/7812
	// https://github.com/terraform-providers/terraform-provider-azurerm/issues/8083
	// USGovernmentCloud allow_blob_public_access and min_tls_version allowed as of issue 9128
	// https://github.com/terraform-providers/terraform-provider-azurerm/issues/9128
	if envName != autorestAzure.PublicCloud.Name && envName != autorestAzure.USGovernmentCloud.Name && envName != autorestAzure.ChinaCloud.Name {
		if allowBlobPublicAccess || minimumTLSVersion != string(storage.TLS10) {
			return fmt.Errorf(`"allow_blob_public_access" and "min_tls_version" are not supported for a Storage Account located in %q`, envName)
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
	if accountKind == string(storage.BlobStorage) {
		if string(parameters.Sku.Name) == string(storage.StandardZRS) {
			return fmt.Errorf("A `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts.")
		}
	}

	// AccessTier is only valid for BlobStorage, StorageV2, and FileStorage accounts
	if accountKind == string(storage.BlobStorage) || accountKind == string(storage.StorageV2) || accountKind == string(storage.FileStorage) {
		accessTier, ok := d.GetOk("access_tier")
		if !ok {
			// default to "Hot"
			accessTier = string(storage.Hot)
		}

		parameters.AccountPropertiesCreateParameters.AccessTier = storage.AccessTier(accessTier.(string))
	} else if isHnsEnabled && accountKind != string(storage.BlockBlobStorage) {
		return fmt.Errorf("`is_hns_enabled` can only be used with account kinds `StorageV2`, `BlobStorage` and `BlockBlobStorage`")
	}

	// NFSv3 is supported for standard general-purpose v2 storage accounts and for premium block blob storage accounts.
	// (https://docs.microsoft.com/en-us/azure/storage/blobs/network-file-system-protocol-support-how-to#step-5-create-and-configure-a-storage-account)
	if nfsV3Enabled &&
		!((accountTier == string(storage.Premium) && accountKind == string(storage.BlockBlobStorage)) ||
			(accountTier == string(storage.Standard) && accountKind == string(storage.StorageV2))) {
		return fmt.Errorf("`nfsv3_enabled` can only be used with account tier `Standard` and account kind `StorageV2`, or account tier `Premium` and account kind `BlockBlobStorage`")
	}
	if nfsV3Enabled && enableHTTPSTrafficOnly {
		return fmt.Errorf("`nfsv3_enabled` can only be used when `enable_https_traffic_only` is `false`")
	}
	if nfsV3Enabled && !isHnsEnabled {
		return fmt.Errorf("`nfsv3_enabled` can only be used when `is_hns_enabled` is `true`")
	}

	// AccountTier must be Premium for FileStorage
	if accountKind == string(storage.FileStorage) {
		if string(parameters.Sku.Tier) == string(storage.StandardLRS) {
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

	// Create
	future, err := client.Create(ctx, resourceGroupName, storageAccountName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Azure Storage Account %q: %+v", storageAccountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for Azure Storage Account %q to be created: %+v", storageAccountName, err)
	}

	account, err := client.GetProperties(ctx, resourceGroupName, storageAccountName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Storage Account %q: %+v", storageAccountName, err)
	}

	if account.ID == nil {
		return fmt.Errorf("Cannot read Storage Account %q (resource group %q) ID",
			storageAccountName, resourceGroupName)
	}
	log.Printf("[INFO] storage account %q ID: %q", storageAccountName, *account.ID)
	d.SetId(*account.ID)

	if val, ok := d.GetOk("blob_properties"); ok {
		// FileStorage does not support blob settings
		if accountKind != string(storage.FileStorage) {
			blobClient := meta.(*clients.Client).Storage.BlobServicesClient

			blobProperties := expandBlobProperties(val.([]interface{}))

			// last_access_time_enabled and container_delete_retention_policy are not supported in USGov
			// Fix issue https://github.com/terraform-providers/terraform-provider-azurerm/issues/11772
			if v := d.Get("blob_properties.0.last_access_time_enabled").(bool); v {
				blobProperties.LastAccessTimeTrackingPolicy = &storage.LastAccessTimeTrackingPolicy{
					Enable: utils.Bool(v),
				}
			}

			if v, ok := d.GetOk("blob_properties.0.container_delete_retention_policy"); ok {
				blobProperties.ContainerDeleteRetentionPolicy = expandBlobPropertiesDeleteRetentionPolicy(v.([]interface{}), false)
			}

			if _, err = blobClient.SetServiceProperties(ctx, resourceGroupName, storageAccountName, *blobProperties); err != nil {
				return fmt.Errorf("Error updating Azure Storage Account `blob_properties` %q: %+v", storageAccountName, err)
			}
		} else {
			return fmt.Errorf("`blob_properties` aren't supported for File Storage accounts.")
		}
	}

	if val, ok := d.GetOk("queue_properties"); ok {
		storageClient := meta.(*clients.Client).Storage
		account, err := storageClient.FindAccount(ctx, storageAccountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q: %s", storageAccountName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", storageAccountName)
		}

		queueClient, err := storageClient.QueuesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Queues Client: %s", err)
		}

		queueProperties, err := expandQueueProperties(val.([]interface{}))
		if err != nil {
			return fmt.Errorf("Error expanding `queue_properties` for Azure Storage Account %q: %+v", storageAccountName, err)
		}

		if err = queueClient.UpdateServiceProperties(ctx, account.ResourceGroup, storageAccountName, queueProperties); err != nil {
			return fmt.Errorf("updating Queue Properties for Storage Account %q: %+v", storageAccountName, err)
		}
	}

	if val, ok := d.GetOk("share_properties"); ok {
		// BlobStorage does not support file share settings
		// FileStorage Premium is supported
		if accountKind == string(storage.FileStorage) || accountKind != string(storage.BlobStorage) && accountKind != string(storage.BlockBlobStorage) && accountTier != string(storage.Premium) {
			fileServiceClient := meta.(*clients.Client).Storage.FileServicesClient

			if _, err = fileServiceClient.SetServiceProperties(ctx, resourceGroupName, storageAccountName, expandShareProperties(val.([]interface{}))); err != nil {
				return fmt.Errorf("updating Azure Storage Account `share_properties` %q: %+v", storageAccountName, err)
			}
		} else {
			return fmt.Errorf("`share_properties` aren't supported for Blob Storage / Block Blob / StorageV2 Premium Storage accounts")
		}
	}

	if val, ok := d.GetOk("static_website"); ok {
		// static website only supported on StorageV2 and BlockBlobStorage
		if accountKind != string(storage.StorageV2) && accountKind != string(storage.BlockBlobStorage) {
			return fmt.Errorf("`static_website` is only supported for StorageV2 and BlockBlobStorage.")
		}
		storageClient := meta.(*clients.Client).Storage

		account, err := storageClient.FindAccount(ctx, storageAccountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q: %s", storageAccountName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", storageAccountName)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps := expandStaticWebsiteProperties(val.([]interface{}))

		if _, err = accountsClient.SetServiceProperties(ctx, storageAccountName, staticWebsiteProps); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account `static_website` %q: %+v", storageAccountName, err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}

func resourceStorageAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	envName := meta.(*clients.Client).Account.Environment.Name
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	storageAccountName := id.Path["storageAccounts"]
	resourceGroupName := id.ResourceGroup

	locks.ByName(storageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountName, storageAccountResourceName)

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	accountKind := d.Get("account_kind").(string)

	if accountKind == string(storage.BlobStorage) {
		if storageType == string(storage.StandardZRS) {
			return fmt.Errorf("A `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts.")
		}
	}

	// AllowSharedKeyAccess can only be true due to issue: https://github.com/terraform-providers/terraform-provider-azurerm/issues/11460
	// if value is nil that brakes the Portal UI as reported in https://github.com/terraform-providers/terraform-provider-azurerm/issues/11689
	// currently the Portal UI reports nil as false, and per the ARM API documentation nil is true. This manafests itself in the Portal UI
	// when a storage account is created by terraform that the AllowSharedKeyAccess is Disabled when it is actually Enabled, thus confusing out customers
	// to fix this, I have added this code to explicitly to set the value to true if is nil to workaround the Portal UI bug for our customers.
	// this is designed as a passive change, meaning the change will only take effect when the existing storage account is modified in some way if the
	// account already exists. since I have also switched up the default behavor for net new storage accounts to always set this value as true, this issue
	// should automatically correct itself over time with these changes.
	// TODO: Remove code when Portal UI team fixes their code
	existing, err := client.GetProperties(ctx, resourceGroupName, storageAccountName, "")
	if err == nil {
		if sharedKeyAccess := existing.AccountProperties.AllowSharedKeyAccess; sharedKeyAccess == nil {
			allowSharedKeyAccess := true

			opts := storage.AccountUpdateParameters{
				AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
					AllowSharedKeyAccess: &allowSharedKeyAccess,
				},
			}

			if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
				return fmt.Errorf("Error updating Azure Storage Account AllowSharedKeyAccess %q: %+v", storageAccountName, err)
			}
		}
	} else {
		// Should never hit this, but added due to an abundance of caution
		return fmt.Errorf("Error retrieving Azure Storage Account %q AllowSharedKeyAccess: %+v", storageAccountName, err)
	}
	// TODO: end remove changes when Portal UI team fixed their code

	if d.HasChange("account_replication_type") {
		sku := storage.Sku{
			Name: storage.SkuName(storageType),
		}

		opts := storage.AccountUpdateParameters{
			Sku: &sku,
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account type %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("account_kind") {
		opts := storage.AccountUpdateParameters{
			Kind: storage.Kind(accountKind),
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account account_kind %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("access_tier") {
		accessTier := d.Get("access_tier").(string)

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				AccessTier: storage.AccessTier(accessTier),
			},
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account access_tier %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})

		opts := storage.AccountUpdateParameters{
			Tags: tags.Expand(t),
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account tags %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("custom_domain") {
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				CustomDomain: expandStorageAccountCustomDomain(d),
			},
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account Custom Domain %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("enable_https_traffic_only") {
		enableHTTPSTrafficOnly := d.Get("enable_https_traffic_only").(bool)

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				EnableHTTPSTrafficOnly: &enableHTTPSTrafficOnly,
			},
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account enable_https_traffic_only %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("min_tls_version") {
		minimumTLSVersion := d.Get("min_tls_version").(string)

		// For all Clouds except Public, China, and USGovernmentCloud, don't specify "min_tls_version" in request body.
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/8083
		// USGovernmentCloud "min_tls_version" allowed as of issue 9128
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/9128
		if envName != autorestAzure.PublicCloud.Name && envName != autorestAzure.USGovernmentCloud.Name && envName != autorestAzure.ChinaCloud.Name {
			if minimumTLSVersion != string(storage.TLS10) {
				return fmt.Errorf(`"min_tls_version" is not supported for a Storage Account located in %q`, envName)
			}
		} else {
			opts := storage.AccountUpdateParameters{
				AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
					MinimumTLSVersion: storage.MinimumTLSVersion(minimumTLSVersion),
				},
			}

			if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
				return fmt.Errorf("Error updating Azure Storage Account min_tls_version %q: %+v", storageAccountName, err)
			}
		}
	}

	if d.HasChange("allow_blob_public_access") {
		allowBlobPublicAccess := d.Get("allow_blob_public_access").(bool)

		// For all Clouds except Public, China, and USGovernmentCloud, don't specify "allow_blob_public_access" in request body.
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/7812
		// USGovernmentCloud "allow_blob_public_access" allowed as of issue 9128
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/9128
		if envName != autorestAzure.PublicCloud.Name && envName != autorestAzure.USGovernmentCloud.Name && envName != autorestAzure.ChinaCloud.Name {
			if allowBlobPublicAccess {
				return fmt.Errorf(`"allow_blob_public_access" is not supported for a Storage Account located in %q`, envName)
			}
		} else {
			opts := storage.AccountUpdateParameters{
				AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
					AllowBlobPublicAccess: &allowBlobPublicAccess,
				},
			}

			if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
				return fmt.Errorf("Error updating Azure Storage Account allow_blob_public_access %q: %+v", storageAccountName, err)
			}
		}
	}

	if d.HasChange("identity") {
		storageAccountIdentity, err := expandAzureRmStorageAccountIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return err
		}
		opts := storage.AccountUpdateParameters{
			Identity: storageAccountIdentity,
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account identity %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("network_rules") {
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				NetworkRuleSet: expandStorageAccountNetworkRules(d, tenantId),
			},
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account network_rules %q: %+v", storageAccountName, err)
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

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account network_rules %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("routing") {
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				RoutingPreference: expandArmStorageAccountRouting(d.Get("routing").([]interface{})),
			},
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account routing %q: %+v", storageAccountName, err)
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
			if _, err := client.Update(ctx, resourceGroupName, storageAccountName, dsNone); err != nil {
				return fmt.Errorf("updating Azure Storage Account azure_files_authentication %q: %+v", storageAccountName, err)
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

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("updating Azure Storage Account azure_files_authentication %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("blob_properties") {
		// FileStorage does not support blob settings
		if accountKind != string(storage.FileStorage) {
			blobClient := meta.(*clients.Client).Storage.BlobServicesClient
			blobProperties := expandBlobProperties(d.Get("blob_properties").([]interface{}))

			// last_access_time_enabled and container_delete_retention_policy are not supported in USGov
			// Fix issue https://github.com/terraform-providers/terraform-provider-azurerm/issues/11772
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

			if _, err = blobClient.SetServiceProperties(ctx, resourceGroupName, storageAccountName, *blobProperties); err != nil {
				return fmt.Errorf("Error updating Azure Storage Account `blob_properties` %q: %+v", storageAccountName, err)
			}
		} else {
			return fmt.Errorf("`blob_properties` aren't supported for File Storage accounts.")
		}
	}

	if d.HasChange("queue_properties") {
		storageClient := meta.(*clients.Client).Storage
		account, err := storageClient.FindAccount(ctx, storageAccountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q: %s", storageAccountName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", storageAccountName)
		}

		queueClient, err := storageClient.QueuesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Queues Client: %s", err)
		}

		queueProperties, err := expandQueueProperties(d.Get("queue_properties").([]interface{}))
		if err != nil {
			return fmt.Errorf("Error expanding `queue_properties` for Azure Storage Account %q: %+v", storageAccountName, err)
		}

		if err = queueClient.UpdateServiceProperties(ctx, account.ResourceGroup, storageAccountName, queueProperties); err != nil {
			return fmt.Errorf("updating Queue Properties for Storage Account %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("share_properties") {
		// BlobStorage, BlockBlobStorage does not support file share settings
		// FileStorage Premium is supported
		if accountKind == string(storage.FileStorage) || accountKind != string(storage.BlobStorage) && accountKind != string(storage.BlockBlobStorage) && accountTier != string(storage.Premium) {
			fileServiceClient := meta.(*clients.Client).Storage.FileServicesClient

			if _, err = fileServiceClient.SetServiceProperties(ctx, resourceGroupName, storageAccountName, expandShareProperties(d.Get("share_properties").([]interface{}))); err != nil {
				return fmt.Errorf("updating Azure Storage Account `file share_properties` %q: %+v", storageAccountName, err)
			}
		} else {
			return fmt.Errorf("`share_properties` aren't supported for Blob Storage /Block Blob /StorageV2 Premium Storage accounts")
		}
	}

	if d.HasChange("static_website") {
		// static website only supported on StorageV2 and BlockBlobStorage
		if accountKind != string(storage.StorageV2) && accountKind != string(storage.BlockBlobStorage) {
			return fmt.Errorf("`static_website` is only supported for StorageV2 and BlockBlobStorage.")
		}
		storageClient := meta.(*clients.Client).Storage

		account, err := storageClient.FindAccount(ctx, storageAccountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q: %s", storageAccountName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", storageAccountName)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps := expandStaticWebsiteProperties(d.Get("static_website").([]interface{}))

		if _, err = accountsClient.SetServiceProperties(ctx, storageAccountName, staticWebsiteProps); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account `static_website` %q: %+v", storageAccountName, err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}

func resourceStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["storageAccounts"]
	resGroup := id.ResourceGroup

	resp, err := client.GetProperties(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading the state of AzureRM Storage Account %q: %+v", name, err)
	}

	// handle the user not having permissions to list the keys
	d.Set("primary_connection_string", "")
	d.Set("secondary_connection_string", "")
	d.Set("primary_blob_connection_string", "")
	d.Set("secondary_blob_connection_string", "")
	d.Set("primary_access_key", "")
	d.Set("secondary_access_key", "")

	keys, err := client.ListKeys(ctx, resGroup, name, storage.Kerb)
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
			return fmt.Errorf("Error listing Keys for Storage Account %q (Resource Group %q): %s", name, resGroup, err)
		}
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
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
		d.Set("allow_blob_public_access", props.AllowBlobPublicAccess)
		// For all Clouds except Public, China, and USGovernmentCloud, "min_tls_version" is not returned from Azure so always persist the default values for "min_tls_version".
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/7812
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/8083
		// USGovernmentCloud "min_tls_version" allowed as of issue 9128
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/9128
		envName := meta.(*clients.Client).Account.Environment.Name
		if envName != autorestAzure.PublicCloud.Name && envName != autorestAzure.USGovernmentCloud.Name && envName != autorestAzure.ChinaCloud.Name {
			d.Set("min_tls_version", string(storage.TLS10))
		} else {
			// For storage account created using old API, the response of GET call will not return "min_tls_version", either.
			minTlsVersion := string(storage.TLS10)
			if props.MinimumTLSVersion != "" {
				minTlsVersion = string(props.MinimumTLSVersion)
			}
			d.Set("min_tls_version", minTlsVersion)
		}

		if customDomain := props.CustomDomain; customDomain != nil {
			if err := d.Set("custom_domain", flattenStorageAccountCustomDomain(customDomain)); err != nil {
				return fmt.Errorf("Error setting `custom_domain`: %+v", err)
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
			return fmt.Errorf("error setting primary endpoints and hosts for blob, queue, table and file: %+v", err)
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
			return fmt.Errorf("error setting secondary endpoints and hosts for blob, queue, table: %+v", err)
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
			return fmt.Errorf("Error setting `network_rules`: %+v", err)
		}

		if props.LargeFileSharesState != "" {
			d.Set("large_file_share_enabled", props.LargeFileSharesState == storage.LargeFileSharesStateEnabled)
		}
	}

	if accessKeys := keys.Keys; accessKeys != nil {
		storageAccountKeys := *accessKeys
		d.Set("primary_access_key", storageAccountKeys[0].Value)
		d.Set("secondary_access_key", storageAccountKeys[1].Value)
	}

	identity, err := flattenAzureRmStorageAccountIdentity(resp.Identity)
	if err != nil {
		return err
	}
	if err := d.Set("identity", identity); err != nil {
		return err
	}

	storageClient := meta.(*clients.Client).Storage
	account, err := storageClient.FindAccount(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q: %s", name, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", name)
	}

	blobClient := storageClient.BlobServicesClient

	// FileStorage does not support blob settings
	if resp.Kind != storage.FileStorage {
		blobProps, err := blobClient.GetServiceProperties(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(blobProps.Response) {
				return fmt.Errorf("Error reading blob properties for AzureRM Storage Account %q: %+v", name, err)
			}
		}

		if err := d.Set("blob_properties", flattenBlobProperties(blobProps)); err != nil {
			return fmt.Errorf("Error setting `blob_properties `for AzureRM Storage Account %q: %+v", name, err)
		}
	}

	fileServiceClient := storageClient.FileServicesClient

	// FileStorage does not support blob kind, FileStorage Premium is supported
	if resp.Kind == storage.FileStorage || resp.Kind != storage.BlobStorage && resp.Kind != storage.BlockBlobStorage && resp.Sku != nil && resp.Sku.Tier != storage.Premium {
		shareProps, err := fileServiceClient.GetServiceProperties(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(shareProps.Response) {
				return fmt.Errorf("reading share properties for AzureRM Storage Account %q: %+v", name, err)
			}
		}

		if err := d.Set("share_properties", flattenShareProperties(shareProps)); err != nil {
			return fmt.Errorf("setting `share_properties `for AzureRM Storage Account %q: %+v", name, err)
		}
	}

	// queue is only available for certain tier and kind (as specified below)
	if resp.Sku == nil {
		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): `sku` was nil", name, resGroup)
	}

	if resp.Sku.Tier == storage.Standard {
		if resp.Kind == storage.Storage || resp.Kind == storage.StorageV2 {
			queueClient, err := storageClient.QueuesClient(ctx, *account)
			if err != nil {
				return fmt.Errorf("Error building Queues Client: %s", err)
			}

			queueProps, err := queueClient.GetServiceProperties(ctx, account.ResourceGroup, name)
			if err != nil {
				return fmt.Errorf("Error reading queue properties for AzureRM Storage Account %q: %+v", name, err)
			}

			if err := d.Set("queue_properties", flattenQueueProperties(queueProps)); err != nil {
				return fmt.Errorf("setting `queue_properties`: %+v", err)
			}
		}
	}

	var staticWebsite []interface{}

	// static website only supported on StorageV2 and BlockBlobStorage
	if resp.Kind == storage.StorageV2 || resp.Kind == storage.BlockBlobStorage {
		storageClient := meta.(*clients.Client).Storage

		account, err := storageClient.FindAccount(ctx, name)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q: %s", name, err)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps, err := accountsClient.GetServiceProperties(ctx, name)
		if err != nil {
			if staticWebsiteProps.Response.Response != nil && !utils.ResponseWasNotFound(staticWebsiteProps.Response) {
				return fmt.Errorf("Error reading static website for AzureRM Storage Account %q: %+v", name, err)
			}
		}

		staticWebsite = flattenStaticWebsiteProperties(staticWebsiteProps)
	}

	if err := d.Set("static_website", staticWebsite); err != nil {
		return fmt.Errorf("Error setting `static_website `for AzureRM Storage Account %q: %+v", name, err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["storageAccounts"]
	resourceGroup := id.ResourceGroup

	locks.ByName(name, storageAccountResourceName)
	defer locks.UnlockByName(name, storageAccountResourceName)

	read, err := client.GetProperties(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return nil
		}

		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): %+v", name, resourceGroup, err)
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

					id, err2 := azure.ParseAzureResourceID(*v.VirtualNetworkResourceID)
					if err2 != nil {
						return err2
					}

					networkName := id.Path["virtualNetworks"]
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

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error issuing delete request for Storage Account %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	// remove this from the cache
	storageClient.RemoveAccountFromCache(name)

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
			Action:           storage.Allow,
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
			Action:                   storage.Allow,
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
			Enabled: false,
		},
		MinuteMetrics: &queues.MetricsConfig{
			Enabled: false,
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
		return properties, fmt.Errorf("Error expanding `minute_metrics`: %+v", err)
	}
	properties.HourMetrics, err = expandQueuePropertiesMetrics(attrs["hour_metrics"].([]interface{}))
	if err != nil {
		return properties, fmt.Errorf("Error expanding `hour_metrics`: %+v", err)
	}

	return properties, nil
}

func expandQueuePropertiesMetrics(input []interface{}) (*queues.MetricsConfig, error) {
	if len(input) == 0 {
		return &queues.MetricsConfig{}, nil
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
		return &queues.LoggingConfig{}
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

func expandAzureRmStorageAccountIdentity(vs []interface{}) (*storage.Identity, error) {
	if len(vs) == 0 {
		return &storage.Identity{
			Type: storage.IdentityTypeNone,
		}, nil
	}

	v := vs[0].(map[string]interface{})
	identity := storage.Identity{
		Type: storage.IdentityType(v["type"].(string)),
	}

	var identityIdSet []interface{}
	if identityIds, exists := v["identity_ids"]; exists {
		identityIdSet = identityIds.(*pluginsdk.Set).List()
	}

	// If type contains `UserAssigned`, `identity_ids` must be specified and have at least 1 element
	if identity.Type == storage.IdentityTypeUserAssigned || identity.Type == storage.IdentityTypeSystemAssignedUserAssigned {
		if len(identityIdSet) == 0 {
			return nil, fmt.Errorf("`identity_ids` must have at least 1 element when `type` includes `UserAssigned`")
		}

		userAssignedIdentities := make(map[string]*storage.UserAssignedIdentity)
		for _, id := range identityIdSet {
			userAssignedIdentities[id.(string)] = &storage.UserAssignedIdentity{}
		}

		identity.UserAssignedIdentities = userAssignedIdentities
	} else if len(identityIdSet) > 0 {
		// If type does _not_ contain `UserAssigned` (i.e. is set to `SystemAssigned` or defaulted to `None`), `identity_ids` is not allowed
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`; but `type` is currently %q", identity.Type)
	}

	return &identity, nil
}

func flattenAzureRmStorageAccountIdentity(identity *storage.Identity) ([]interface{}, error) {
	if identity == nil || identity.Type == storage.IdentityTypeNone {
		return make([]interface{}, 0), nil
	}

	var principalId, tenantId string
	if identity.PrincipalID != nil {
		principalId = *identity.PrincipalID
	}

	if identity.TenantID != nil {
		tenantId = *identity.TenantID
	}

	identityIds := make([]interface{}, 0)
	if identity.UserAssignedIdentities != nil {
		for key := range identity.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityID(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(identity.Type),
			"principal_id": principalId,
			"tenant_id":    tenantId,
			"identity_ids": pluginsdk.NewSet(pluginsdk.HashString, identityIds),
		},
	}, nil
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
