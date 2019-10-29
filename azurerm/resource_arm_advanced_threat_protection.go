package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	azautorest "github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-getter/helper/url"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/queue/queues"
)

const blobStorageAccountDefaultAccessTier = "Hot"

func resourceArmAdvancedThreatProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAdvancedThreatProtectionCreate,
		Read:   resourceArmAdvancedThreatProtectionRead,
		Update: resourceArmAdvancedThreatProtectionUpdate,

		MigrateState:  resourceStorageAccountMigrateState,
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
				ValidateFunc: validateArmStorageAccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"location": azure.SchemaLocation(),
		},
	}
}

func resourceArmAdvancedThreatProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Storage.AccountsClient
	advancedThreatProtectionClient := meta.(*ArmClient).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	storageAccountName := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetProperties(ctx, resourceGroupName, storageAccountName, "")
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
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	enableBlobEncryption := d.Get("enable_blob_encryption").(bool)
	enableFileEncryption := d.Get("enable_file_encryption").(bool)
	enableHTTPSTrafficOnly := d.Get("enable_https_traffic_only").(bool)
	isHnsEnabled := d.Get("is_hns_enabled").(bool)

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	storageAccountEncryptionSource := d.Get("account_encryption_source").(string)

	parameters := storage.AccountCreateParameters{
		Location: &location,
		Sku: &storage.Sku{
			Name: storage.SkuName(storageType),
		},
		Tags: tags.Expand(t),
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
			NetworkRuleSet:         expandStorageAccountNetworkRules(d),
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

	// AccessTier is only valid for BlobStorage, StorageV2, and FileStorage accounts
	if accountKind == string(storage.BlobStorage) || accountKind == string(storage.StorageV2) || accountKind == string(storage.FileStorage) {
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

	// AccountTier must be Premium for FileStorage
	if accountKind == string(storage.FileStorage) {
		if string(parameters.Sku.Tier) == string(storage.StandardLRS) {
			return fmt.Errorf("A `account_tier` of `Standard` is not supported for FileStorage accounts.")
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

	// TODO: deprecate & split this out into it's own resource in 2.0
	// as this is not available in all regions, and presumably off by default
	// lets only try to set this value when true
	// TODO in 2.0 switch to guarding this with d.GetOkExists() ?
	if v := d.Get("enable_advanced_threat_protection").(bool); v {
		advancedThreatProtectionSetting := security.AdvancedThreatProtectionSetting{
			AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
				IsEnabled: utils.Bool(v),
			},
		}

		if _, err = advancedThreatProtectionClient.Create(ctx, d.Id(), advancedThreatProtectionSetting); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account enable_advanced_threat_protection %q: %+v", storageAccountName, err)
		}
	}

	if val, ok := d.GetOk("queue_properties"); ok {
		storageClient := meta.(*ArmClient).Storage
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

	return resourceArmAdvancedThreatProtectionRead(d, meta)
}

// resourceArmStorageAccountUpdate is unusual in the ARM API where most resources have a combined
// and idempotent operation for CreateOrUpdate. In particular updating all of the parameters
// available requires a call to Update per parameter...
func resourceArmAdvancedThreatProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Storage.AccountsClient
	advancedThreatProtectionClient := meta.(*ArmClient).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
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
		t := d.Get("tags").(map[string]interface{})

		opts := storage.AccountUpdateParameters{
			Tags: tags.Expand(t),
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

		if _, err := client.Update(ctx, resourceGroupName, storageAccountName, opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account Encryption %q: %+v", storageAccountName, err)
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

		d.SetPartial("enable_https_traffic_only")
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

		d.SetPartial("network_rules")
	}

	if d.HasChange("enable_advanced_threat_protection") {
		opts := security.AdvancedThreatProtectionSetting{
			AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
				IsEnabled: utils.Bool(d.Get("enable_advanced_threat_protection").(bool)),
			},
		}

		if _, err := advancedThreatProtectionClient.Create(ctx, d.Id(), opts); err != nil {
			return fmt.Errorf("Error updating Azure Storage Account enable_advanced_threat_protection %q: %+v", storageAccountName, err)
		}

		d.SetPartial("enable_advanced_threat_protection")
	}

	if d.HasChange("queue_properties") {
		storageClient := meta.(*ArmClient).Storage
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

		d.SetPartial("queue_properties")
	}

	d.Partial(false)
	return resourceArmAdvancedThreatProtectionRead(d, meta)
}

func resourceArmAdvancedThreatProtectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Storage.AccountsClient
	advancedThreatProtectionClient := meta.(*ArmClient).SecurityCenter.AdvancedThreatProtectionClient
	endpointSuffix := meta.(*ArmClient).environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
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

	keys, err := client.ListKeys(ctx, resGroup, name)
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

		if networkRules := props.NetworkRuleSet; networkRules != nil {
			if err := d.Set("network_rules", flattenStorageAccountNetworkRules(networkRules)); err != nil {
				return fmt.Errorf("Error setting `network_rules`: %+v", err)
			}
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

	// TODO in 2.0 switch to guarding this with d.GetOkExists()
	atp, err := advancedThreatProtectionClient.Get(ctx, d.Id())
	if err != nil {
		msg := err.Error()
		if !strings.Contains(msg, "The resource namespace 'Microsoft.Security' is invalid.") {
			if !strings.Contains(msg, "No registered resource provider found for location '") {
				if !strings.Contains(msg, "' and API version '2017-08-01-preview' for type ") {
					return fmt.Errorf("Error reading the advanced threat protection settings of AzureRM Storage Account %q: %+v", name, err)
				}
			}
		}
	} else {
		if atpp := atp.AdvancedThreatProtectionProperties; atpp != nil {
			d.Set("enable_advanced_threat_protection", atpp.IsEnabled)
		}
	}

	storageClient := meta.(*ArmClient).Storage
	account, err := storageClient.FindAccount(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q: %s", name, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", name)
	}

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

	return tags.FlattenAndSet(d, resp.Tags)
}
