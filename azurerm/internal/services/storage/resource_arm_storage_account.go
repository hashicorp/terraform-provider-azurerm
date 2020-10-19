package storage

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	azautorest "github.com/Azure/go-autorest/autorest"
	autorestAzure "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/go-getter/helper/url"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/queue/queues"
)

var storageAccountResourceName = "azurerm_storage_account"

func resourceArmStorageAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageAccountCreate,
		Read:   resourceArmStorageAccountRead,
		Update: resourceArmStorageAccountUpdate,
		Delete: resourceArmStorageAccountDelete,

		MigrateState:  ResourceStorageAccountMigrateState,
		SchemaVersion: 2,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateArmStorageAccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"account_kind": {
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"account_replication_type": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.Cool),
					string(storage.Hot),
				}, true),
			},

			"custom_domain": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"use_subdomain": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"enable_https_traffic_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"min_tls_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(storage.TLS10),
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.TLS10),
					string(storage.TLS11),
					string(storage.TLS12),
				}, false),
			},

			"is_hns_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"allow_blob_public_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"network_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bypass": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(storage.AzureServices),
									string(storage.Logging),
									string(storage.Metrics),
									string(storage.None),
								}, true),
							},
							Set: schema.HashString,
						},

						"ip_rules": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.StorageAccountIPRule,
							},
							Set: schema.HashString,
						},

						"virtual_network_subnet_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"default_action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(storage.DefaultActionAllow),
								string(storage.DefaultActionDeny),
							}, false),
						},
					},
				},
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, true),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"blob_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cors_rule": azure.SchemaStorageAccountCorsRule(true),
						"delete_retention_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:         schema.TypeInt,
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cors_rule": azure.SchemaStorageAccountCorsRule(false),
						"logging": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"delete": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"read": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"write": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"retention_policy_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},
						"hour_metrics": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"include_apis": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"retention_policy_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},
						"minute_metrics": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"include_apis": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"retention_policy_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},
					},
				},
			},

			"share_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cors_rule": azure.SchemaStorageAccountCorsRule(true),
						"delete_retention_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:         schema.TypeInt,
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

			"static_website": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"error_404_document": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"large_file_share_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"primary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_blob_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_blob_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_blob_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_blob_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_queue_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_queue_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_queue_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_queue_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_table_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_table_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_table_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_table_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_web_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_web_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_web_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_web_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_dfs_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_dfs_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_dfs_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_dfs_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_file_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_file_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_file_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_file_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_blob_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_blob_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: validateAzureRMStorageAccountTags,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
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
		},
	}
}

func validateAzureRMStorageAccountTags(v interface{}, _ string) (warnings []string, errors []error) {
	tagsMap := v.(map[string]interface{})

	if len(tagsMap) > 50 {
		errors = append(errors, fmt.Errorf("a maximum of 50 tags can be applied to storage account ARM resource"))
	}

	for k, v := range tagsMap {
		if len(k) > 128 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag key is 128 characters: %q is %d characters", k, len(k)))
		}

		value, err := tags.TagValueToString(v)
		if err != nil {
			errors = append(errors, err)
		} else if len(value) > 256 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag value is 256 characters: the value for %q is %d characters", k, len(value)))
		}
	}

	return warnings, errors
}

func resourceArmStorageAccountCreate(d *schema.ResourceData, meta interface{}) error {
	envName := meta.(*clients.Client).Account.Environment.Name
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
	allowBlobPublicAccess := d.Get("allow_blob_public_access").(bool)

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)

	parameters := storage.AccountCreateParameters{
		Location: &location,
		Sku: &storage.Sku{
			Name: storage.SkuName(storageType),
		},
		Tags: tags.Expand(t),
		Kind: storage.Kind(accountKind),
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
			EnableHTTPSTrafficOnly: &enableHTTPSTrafficOnly,
			NetworkRuleSet:         expandStorageAccountNetworkRules(d),
			IsHnsEnabled:           &isHnsEnabled,
		},
	}

	// For all Clouds except Public, don't specify "allow_blob_public_access" and "min_tls_version" in request body.
	// https://github.com/terraform-providers/terraform-provider-azurerm/issues/7812
	// https://github.com/terraform-providers/terraform-provider-azurerm/issues/8083
	if envName != autorestAzure.PublicCloud.Name {
		if allowBlobPublicAccess || minimumTLSVersion != string(storage.TLS10) {
			return fmt.Errorf(`"allow_blob_public_access" and "min_tls_version" are not supported for a Storage Account located in %q`, envName)
		}
	} else {
		parameters.AccountPropertiesCreateParameters.AllowBlobPublicAccess = &allowBlobPublicAccess
		parameters.AccountPropertiesCreateParameters.MinimumTLSVersion = storage.MinimumTLSVersion(minimumTLSVersion)
	}

	if _, ok := d.GetOk("identity"); ok {
		storageAccountIdentity := expandAzureRmStorageAccountIdentity(d)
		parameters.Identity = storageAccountIdentity
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
	} else if isHnsEnabled {
		return fmt.Errorf("`is_hns_enabled` can only be used with account kinds `StorageV2` and `BlobStorage`")
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

			if _, err = blobClient.SetServiceProperties(ctx, resourceGroupName, storageAccountName, blobProperties); err != nil {
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

		if _, err = queueClient.SetServiceProperties(ctx, storageAccountName, queueProperties); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account `queue_properties` %q: %+v", storageAccountName, err)
		}
	}

	if val, ok := d.GetOk("share_properties"); ok {
		// BlobStorage, Premium general purpose does not support file share settings
		if accountKind != string(storage.BlobStorage) && accountKind != string(storage.BlockBlobStorage) && accountTier != string(storage.Premium) {
			fileServiceClient := meta.(*clients.Client).Storage.FileServicesClient

			if _, err = fileServiceClient.SetServiceProperties(ctx, resourceGroupName, storageAccountName, expandShareProperties(val.([]interface{}))); err != nil {
				return fmt.Errorf("updating Azure Storage Account `file share_properties` %q: %+v", storageAccountName, err)
			}
		} else {
			return fmt.Errorf("`share_properties` aren't supported for Blob Storage / Block Blob Storage /Premium accounts")
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

	return resourceArmStorageAccountRead(d, meta)
}

func resourceArmStorageAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	envName := meta.(*clients.Client).Account.Environment.Name
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

		// For all Clouds except Public, don't specify "min_tls_version" in request body.
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/8083
		if envName != autorestAzure.PublicCloud.Name {
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

		// For all Clouds except Public, don't specify "allow_blob_public_access" in request body.
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/7812
		if envName != autorestAzure.PublicCloud.Name {
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
		opts := storage.AccountUpdateParameters{
			Identity: expandAzureRmStorageAccountIdentity(d),
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account identity %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("network_rules") {
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				NetworkRuleSet: expandStorageAccountNetworkRules(d),
			},
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account network_rules %q: %+v", storageAccountName, err)
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

	if d.HasChange("blob_properties") {
		// FileStorage does not support blob settings
		if accountKind != string(storage.FileStorage) {
			blobClient := meta.(*clients.Client).Storage.BlobServicesClient
			blobProperties := expandBlobProperties(d.Get("blob_properties").([]interface{}))

			if _, err = blobClient.SetServiceProperties(ctx, resourceGroupName, storageAccountName, blobProperties); err != nil {
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

		if _, err = queueClient.SetServiceProperties(ctx, storageAccountName, queueProperties); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account `queue_properties` %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("share_properties") {
		// BlobStorage, BlockBlobStorage, Premium does not support file share settings
		if accountKind != string(storage.BlobStorage) && accountKind != string(storage.BlockBlobStorage) && accountTier != string(storage.Premium) {
			fileServiceClient := meta.(*clients.Client).Storage.FileServicesClient

			if _, err = fileServiceClient.SetServiceProperties(ctx, resourceGroupName, storageAccountName, expandShareProperties(d.Get("share_properties").([]interface{}))); err != nil {
				return fmt.Errorf("updating Azure Storage Account `file share_properties` %q: %+v", storageAccountName, err)
			}
		} else {
			return fmt.Errorf("`share_properties` aren't supported for Blob Storage /Block Blob Storage / Premium accounts")
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

	return resourceArmStorageAccountRead(d, meta)
}

func resourceArmStorageAccountRead(d *schema.ResourceData, meta interface{}) error {
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
		d.Set("enable_https_traffic_only", props.EnableHTTPSTrafficOnly)
		d.Set("is_hns_enabled", props.IsHnsEnabled)
		d.Set("allow_blob_public_access", props.AllowBlobPublicAccess)
		// For all Clouds except Public, "min_tls_version" is not returned from Azure so always persist the default values for "min_tls_version".
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/7812
		// https://github.com/terraform-providers/terraform-provider-azurerm/issues/8083
		if meta.(*clients.Client).Account.Environment.Name != autorestAzure.PublicCloud.Name {
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

	identity := flattenAzureRmStorageAccountIdentity(resp.Identity)
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
	fileServiceClient := storageClient.FileServicesClient

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

	// FileStorage does not support blob settings
	if resp.Kind != storage.BlobStorage && resp.Kind != storage.BlockBlobStorage && resp.Sku.Tier != storage.Premium {
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

			queueProps, err := queueClient.GetServiceProperties(ctx, name)
			if err != nil {
				if queueProps.Response.Response != nil && !utils.ResponseWasNotFound(queueProps.Response) {
					return fmt.Errorf("Error reading queue properties for AzureRM Storage Account %q: %+v", name, err)
				}
			}

			if err := d.Set("queue_properties", flattenQueueProperties(queueProps)); err != nil {
				return fmt.Errorf("Error setting `queue_properties `for AzureRM Storage Account %q: %+v", name, err)
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

func resourceArmStorageAccountDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandStorageAccountCustomDomain(d *schema.ResourceData) *storage.CustomDomain {
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

func expandStorageAccountNetworkRules(d *schema.ResourceData) *storage.NetworkRuleSet {
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
	}

	if v := networkRule["default_action"]; v != nil {
		networkRuleSet.DefaultAction = storage.DefaultAction(v.(string))
	}

	return networkRuleSet
}

func expandStorageAccountIPRules(networkRule map[string]interface{}) *[]storage.IPRule {
	ipRulesInfo := networkRule["ip_rules"].(*schema.Set).List()
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
	virtualNetworkInfo := networkRule["virtual_network_subnet_ids"].(*schema.Set).List()
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
	bypassInfo := networkRule["bypass"].(*schema.Set).List()

	var bypassValues []string
	for _, bypassConfig := range bypassInfo {
		bypassValues = append(bypassValues, bypassConfig.(string))
	}

	return storage.Bypass(strings.Join(bypassValues, ", "))
}

func expandBlobProperties(input []interface{}) storage.BlobServiceProperties {
	props := storage.BlobServiceProperties{
		BlobServicePropertiesProperties: &storage.BlobServicePropertiesProperties{
			Cors: &storage.CorsRules{
				CorsRules: &[]storage.CorsRule{},
			},
			DeleteRetentionPolicy: &storage.DeleteRetentionPolicy{
				Enabled: utils.Bool(false),
			},
		},
	}

	if len(input) == 0 || input[0] == nil {
		return props
	}

	v := input[0].(map[string]interface{})

	deletePolicyRaw := v["delete_retention_policy"].([]interface{})
	props.BlobServicePropertiesProperties.DeleteRetentionPolicy = expandBlobPropertiesDeleteRetentionPolicy(deletePolicyRaw)

	corsRaw := v["cors_rule"].([]interface{})
	props.BlobServicePropertiesProperties.Cors = expandBlobPropertiesCors(corsRaw)

	return props
}

func expandBlobPropertiesDeleteRetentionPolicy(input []interface{}) *storage.DeleteRetentionPolicy {
	deleteRetentionPolicy := storage.DeleteRetentionPolicy{
		Enabled: utils.Bool(false),
	}

	if len(input) == 0 {
		return &deleteRetentionPolicy
	}

	policy := input[0].(map[string]interface{})
	days := policy["days"].(int)
	deleteRetentionPolicy.Enabled = utils.Bool(true)
	deleteRetentionPolicy.Days = utils.Int32(int32(days))

	return &deleteRetentionPolicy
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

	deletePolicyRaw := v["delete_retention_policy"].([]interface{})
	props.FileServicePropertiesProperties.ShareDeleteRetentionPolicy = expandBlobPropertiesDeleteRetentionPolicy(deletePolicyRaw)

	corsRaw := v["cors_rule"].([]interface{})
	props.FileServicePropertiesProperties.Cors = expandBlobPropertiesCors(corsRaw)

	return props
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

func flattenStorageAccountNetworkRules(input *storage.NetworkRuleSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	networkRules := make(map[string]interface{})

	networkRules["ip_rules"] = schema.NewSet(schema.HashString, flattenStorageAccountIPRules(input.IPRules))
	networkRules["virtual_network_subnet_ids"] = schema.NewSet(schema.HashString, flattenStorageAccountVirtualNetworks(input.VirtualNetworkRules))
	networkRules["bypass"] = schema.NewSet(schema.HashString, flattenStorageAccountBypass(input.Bypass))
	networkRules["default_action"] = string(input.DefaultAction)

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

	if len(flattenedCorsRules) == 0 && len(flattenedDeletePolicy) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"cors_rule":               flattenedCorsRules,
			"delete_retention_policy": flattenedDeletePolicy,
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

func flattenQueueProperties(input queues.StorageServicePropertiesResponse) []interface{} {
	if input.Response.Response == nil {
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

	if len(flattenedCorsRules) == 0 && len(flattenedDeletePolicy) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"cors_rule":               flattenedCorsRules,
			"delete_retention_policy": flattenedDeletePolicy,
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

func ValidateArmStorageAccountName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile(`\A([a-z0-9]{3,24})\z`).MatchString(input) {
		errors = append(errors, fmt.Errorf("name (%q) can only consist of lowercase letters and numbers, and must be between 3 and 24 characters long", input))
	}

	return warnings, errors
}

func expandAzureRmStorageAccountIdentity(d *schema.ResourceData) *storage.Identity {
	identities := d.Get("identity").([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := identity["type"].(string)
	return &storage.Identity{
		Type: &identityType,
	}
}

func flattenAzureRmStorageAccountIdentity(identity *storage.Identity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	if identity.Type != nil {
		result["type"] = *identity.Type
	}
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}
	if identity.TenantID != nil {
		result["tenant_id"] = *identity.TenantID
	}

	return []interface{}{result}
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

func flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d *schema.ResourceData, primary *storage.Endpoints) error {
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

func flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d *schema.ResourceData, secondary *storage.Endpoints) error {
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

func setEndpointAndHost(d *schema.ResourceData, ordinalString string, endpointType *string, typeString string) error {
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
