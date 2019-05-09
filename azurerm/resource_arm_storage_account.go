package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-02-01/storage"
	"github.com/hashicorp/go-getter/helper/url"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const blobStorageAccountDefaultAccessTier = "Hot"

func resourceArmStorageAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageAccountCreate,
		Read:   resourceArmStorageAccountRead,
		Update: resourceArmStorageAccountUpdate,
		Delete: resourceArmStorageAccountDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceStorageAccountMigrateState,
		SchemaVersion: 2,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageAccountName,
			},

			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"location": locationSchema(),

			"account_kind": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.Storage),
					string(storage.BlobStorage),
					string(storage.StorageV2),
				}, true),
				Default: string(storage.Storage),
			},

			"account_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Deprecated:       "This field has been split into `account_tier` and `account_replication_type`",
				ValidateFunc:     validateArmStorageAccountType,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"account_tier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"account_replication_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LRS",
					"ZRS",
					"GRS",
					"RAGRS",
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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

			"account_encryption_source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(storage.MicrosoftStorage),
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.MicrosoftKeyvault),
					string(storage.MicrosoftStorage),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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

			"enable_blob_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enable_file_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enable_https_traffic_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"is_hns_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"network_rules": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
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
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"virtual_network_subnet_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
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
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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

			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAzureRMStorageAccountTags,
			},
		},
	}
}

func validateAzureRMStorageAccountTags(v interface{}, _ string) (warnings []string, errors []error) {
	tagsMap := v.(map[string]interface{})

	if len(tagsMap) > 15 {
		errors = append(errors, fmt.Errorf("a maximum of 15 tags can be applied to each ARM resource"))
	}

	for k, v := range tagsMap {
		if len(k) > 128 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag key is 128 characters: %q is %d characters", k, len(k)))
		}

		value, err := tagValueToString(v)
		if err != nil {
			errors = append(errors, err)
		} else if len(value) > 256 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag value is 256 characters: the value for %q is %d characters", k, len(value)))
		}
	}

	return warnings, errors
}

func resourceArmStorageAccountCreate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).storageServiceClient

	storageAccountName := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported {
		existing, err := client.GetProperties(ctx, resourceGroupName, storageAccountName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Storage Account %q (Resource Group %q): %s", storageAccountName, resourceGroupName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_storage_account", *existing.ID)
		}
	}

	accountKind := d.Get("account_kind").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	enableBlobEncryption := d.Get("enable_blob_encryption").(bool)
	enableFileEncryption := d.Get("enable_file_encryption").(bool)
	enableHTTPSTrafficOnly := d.Get("enable_https_traffic_only").(bool)
	isHnsEnabled := d.Get("is_hns_enabled").(bool)

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	storageAccountEncryptionSource := d.Get("account_encryption_source").(string)

	networkRules := expandStorageAccountNetworkRules(d)

	parameters := storage.AccountCreateParameters{
		Location: &location,
		Sku: &storage.Sku{
			Name: storage.SkuName(storageType),
		},
		Tags: expandTags(tags),
		Kind: storage.Kind(accountKind),
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
			Encryption: &storage.Encryption{
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						Enabled: utils.Bool(enableBlobEncryption),
					},
					File: &storage.EncryptionService{
						Enabled: utils.Bool(enableFileEncryption),
					}},
				KeySource: storage.KeySource(storageAccountEncryptionSource),
			},
			EnableHTTPSTrafficOnly: &enableHTTPSTrafficOnly,
			NetworkRuleSet:         networkRules,
			IsHnsEnabled:           &isHnsEnabled,
		},
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

	// AccessTier is only valid for BlobStorage and StorageV2 accounts
	if accountKind == string(storage.BlobStorage) || accountKind == string(storage.StorageV2) {
		accessTier, ok := d.GetOk("access_tier")
		if !ok {
			// default to "Hot"
			accessTier = blobStorageAccountDefaultAccessTier
		}

		parameters.AccountPropertiesCreateParameters.AccessTier = storage.AccessTier(accessTier.(string))
	} else {
		if isHnsEnabled {
			return fmt.Errorf("`is_hns_enabled` can only be used with account kinds `StorageV2` and `BlobStorage`")
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

	account, err := client.GetProperties(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Storage Account %q: %+v", storageAccountName, err)
	}

	if account.ID == nil {
		return fmt.Errorf("Cannot read Storage Account %q (resource group %q) ID",
			storageAccountName, resourceGroupName)
	}
	log.Printf("[INFO] storage account %q ID: %q", storageAccountName, *account.ID)
	d.SetId(*account.ID)

	return resourceArmStorageAccountRead(d, meta)
}

// resourceArmStorageAccountUpdate is unusual in the ARM API where most resources have a combined
// and idempotent operation for CreateOrUpdate. In particular updating all of the parameters
// available requires a call to Update per parameter...
func resourceArmStorageAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).storageServiceClient
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	storageAccountName := id.Path["storageAccounts"]
	resourceGroupName := id.ResourceGroup

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	accountKind := d.Get("account_kind").(string)

	if accountKind == string(storage.BlobStorage) {
		if storageType == string(storage.StandardZRS) {
			return fmt.Errorf("A `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts.")
		}
	}

	d.Partial(true)

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

		d.SetPartial("account_replication_type")
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

		d.SetPartial("access_tier")
	}

	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})

		opts := storage.AccountUpdateParameters{
			Tags: expandTags(tags),
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account tags %q: %+v", storageAccountName, err)
		}

		d.SetPartial("tags")
	}

	if d.HasChange("enable_blob_encryption") || d.HasChange("enable_file_encryption") {
		encryptionSource := d.Get("account_encryption_source").(string)

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				Encryption: &storage.Encryption{
					Services:  &storage.EncryptionServices{},
					KeySource: storage.KeySource(encryptionSource),
				},
			},
		}

		if d.HasChange("enable_blob_encryption") {
			enableEncryption := d.Get("enable_blob_encryption").(bool)
			opts.Encryption.Services.Blob = &storage.EncryptionService{
				Enabled: utils.Bool(enableEncryption),
			}

			d.SetPartial("enable_blob_encryption")
		}

		if d.HasChange("enable_file_encryption") {
			enableEncryption := d.Get("enable_file_encryption").(bool)
			opts.Encryption.Services.File = &storage.EncryptionService{
				Enabled: utils.Bool(enableEncryption),
			}
			d.SetPartial("enable_file_encryption")
		}

		_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
		if err != nil {
			return fmt.Errorf("Error updating Azure Storage Account Encryption %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("custom_domain") {
		customDomain := expandStorageAccountCustomDomain(d)
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				CustomDomain: customDomain,
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

		d.SetPartial("enable_https_traffic_only")
	}

	if d.HasChange("identity") {
		storageAccountIdentity := expandAzureRmStorageAccountIdentity(d)

		opts := storage.AccountUpdateParameters{
			Identity: storageAccountIdentity,
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account identity %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("network_rules") {
		networkRules := expandStorageAccountNetworkRules(d)

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				NetworkRuleSet: networkRules,
			},
		}

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account network_rules %q: %+v", storageAccountName, err)
		}

		d.SetPartial("network_rules")
	}

	d.Partial(false)
	return resourceArmStorageAccountRead(d, meta)
}

func resourceArmStorageAccountRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).storageServiceClient
	endpointSuffix := meta.(*ArmClient).environment.StorageEndpointSuffix

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["storageAccounts"]
	resGroup := id.ResourceGroup

	resp, err := client.GetProperties(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading the state of AzureRM Storage Account %q: %+v", name, err)
	}

	keys, err := client.ListKeys(ctx, resGroup, name)
	if err != nil {
		return err
	}

	accessKeys := *keys.Keys
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("account_kind", resp.Kind)

	if sku := resp.Sku; sku != nil {
		d.Set("account_type", sku.Name)
		d.Set("account_tier", sku.Tier)
		d.Set("account_replication_type", strings.Split(fmt.Sprintf("%v", sku.Name), "_")[1])
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("access_tier", props.AccessTier)
		d.Set("enable_https_traffic_only", props.EnableHTTPSTrafficOnly)
		d.Set("is_hns_enabled", props.IsHnsEnabled)

		if customDomain := props.CustomDomain; customDomain != nil {
			if err := d.Set("custom_domain", flattenStorageAccountCustomDomain(customDomain)); err != nil {
				return fmt.Errorf("Error setting `custom_domain`: %+v", err)
			}
		}

		if encryption := props.Encryption; encryption != nil {
			if services := encryption.Services; services != nil {
				if blob := services.Blob; blob != nil {
					d.Set("enable_blob_encryption", blob.Enabled)
				}
				if file := services.File; file != nil {
					d.Set("enable_file_encryption", file.Enabled)
				}
			}
			d.Set("account_encryption_source", string(encryption.KeySource))
		}

		// Computed
		d.Set("primary_location", props.PrimaryLocation)
		d.Set("secondary_location", props.SecondaryLocation)

		if len(accessKeys) > 0 {
			pcs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *accessKeys[0].Value, endpointSuffix)
			d.Set("primary_connection_string", pcs)
		}

		if len(accessKeys) > 1 {
			scs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *accessKeys[1].Value, endpointSuffix)
			d.Set("secondary_connection_string", scs)
		}

		if err := flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d, props.PrimaryEndpoints); err != nil {
			return fmt.Errorf("error setting primary endpoints and hosts for blob, queue, table and file: %+v", err)
		}

		var primaryBlobConnectStr string
		if v := props.PrimaryEndpoints; v != nil {
			primaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, accessKeys[0].Value)
		}
		d.Set("primary_blob_connection_string", primaryBlobConnectStr)

		if err := flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d, props.SecondaryEndpoints); err != nil {
			return fmt.Errorf("error setting secondary endpoints and hosts for blob, queue, table: %+v", err)
		}

		var secondaryBlobConnectStr string
		if v := props.SecondaryEndpoints; v != nil {
			secondaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, accessKeys[1].Value)
		}
		d.Set("secondary_blob_connection_string", secondaryBlobConnectStr)

		networkRules := props.NetworkRuleSet
		if networkRules != nil {
			if err := d.Set("network_rules", flattenStorageAccountNetworkRules(networkRules)); err != nil {
				return fmt.Errorf("Error setting `network_rules`: %+v", err)
			}
		}
	}

	d.Set("primary_access_key", accessKeys[0].Value)
	d.Set("secondary_access_key", accessKeys[1].Value)

	identity := flattenAzureRmStorageAccountIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmStorageAccountDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).storageServiceClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["storageAccounts"]
	resourceGroup := id.ResourceGroup

	read, err := client.GetProperties(ctx, resourceGroup, name)
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

					id, err2 := parseAzureResourceID(*v.VirtualNetworkResourceID)
					if err2 != nil {
						return err2
					}

					networkName := id.Path["virtualNetworks"]
					virtualNetworkNames = append(virtualNetworkNames, networkName)
				}
			}
		}
	}

	azureRMLockMultipleByName(&virtualNetworkNames, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(&virtualNetworkNames, virtualNetworkResourceName)

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error issuing delete request for Storage Account %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

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
	networkRuleSet := &storage.NetworkRuleSet{}

	networkRuleSet.IPRules = expandStorageAccountIPRules(networkRule)
	networkRuleSet.VirtualNetworkRules = expandStorageAccountVirtualNetworks(networkRule)
	networkRuleSet.Bypass = expandStorageAccountBypass(networkRule)
	// Default Access is disabled when network rules are set.
	networkRuleSet.DefaultAction = storage.DefaultActionDeny

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

func flattenStorageAccountNetworkRules(input *storage.NetworkRuleSet) []interface{} {
	if len(*input.IPRules) == 0 && len(*input.VirtualNetworkRules) == 0 {
		return []interface{}{}
	}
	networkRules := make(map[string]interface{})

	networkRules["ip_rules"] = schema.NewSet(schema.HashString, flattenStorageAccountIPRules(input.IPRules))
	networkRules["virtual_network_subnet_ids"] = schema.NewSet(schema.HashString, flattenStorageAccountVirtualNetworks(input.VirtualNetworkRules))
	networkRules["bypass"] = schema.NewSet(schema.HashString, flattenStorageAccountBypass(input.Bypass))

	return []interface{}{networkRules}
}

func flattenStorageAccountIPRules(input *[]storage.IPRule) []interface{} {
	ipRules := make([]interface{}, len(*input))
	if input != nil {
		for i, ipRule := range *input {
			ipRules[i] = *ipRule.IPAddressOrRange
		}
	}

	return ipRules
}

func flattenStorageAccountVirtualNetworks(input *[]storage.VirtualNetworkRule) []interface{} {
	virtualNetworks := make([]interface{}, len(*input))

	if input != nil {
		for i, virtualNetwork := range *input {
			virtualNetworks[i] = *virtualNetwork.VirtualNetworkResourceID
		}
	}

	return virtualNetworks
}

func flattenStorageAccountBypass(input storage.Bypass) []interface{} {
	bypassValues := strings.Split(string(input), ", ")
	bypass := make([]interface{}, len(bypassValues))

	for i, value := range bypassValues {
		bypass[i] = value
	}

	return bypass
}

func validateArmStorageAccountName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile(`\A([a-z0-9]{3,24})\z`).MatchString(input) {
		errors = append(errors, fmt.Errorf("name can only consist of lowercase letters and numbers, and must be between 3 and 24 characters long"))
	}

	return warnings, errors
}

func validateArmStorageAccountType(v interface{}, _ string) (warnings []string, errors []error) {
	validAccountTypes := []string{"standard_lrs", "standard_zrs",
		"standard_grs", "standard_ragrs", "premium_lrs"}

	input := strings.ToLower(v.(string))

	for _, valid := range validAccountTypes {
		if valid == input {
			return
		}
	}

	errors = append(errors, fmt.Errorf("Invalid storage account type %q", input))
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

	d.Set(fmt.Sprintf("%s_%s_endpoint", ordinalString, typeString), endpoint)
	d.Set(fmt.Sprintf("%s_%s_host", ordinalString, typeString), host)
	return nil
}
