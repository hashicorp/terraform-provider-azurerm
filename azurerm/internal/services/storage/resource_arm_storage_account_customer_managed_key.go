package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	storageParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageAccountCustomerManagedKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageAccountCustomerManagedKeyCreateUpdate,
		Read:   resourceArmStorageAccountCustomerManagedKeyRead,
		Update: resourceArmStorageAccountCustomerManagedKeyCreateUpdate,
		Delete: resourceArmStorageAccountCustomerManagedKeyDelete,

		// TODO: this needs a custom ID validating importer
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.KeyVaultID,
			},

			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"key_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmStorageAccountCustomerManagedKeyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	vaultsClient := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountIDRaw := d.Get("storage_account_id").(string)
	storageAccountID, err := storageParse.AccountID(storageAccountIDRaw)
	if err != nil {
		return err
	}

	locks.ByName(storageAccountID.Name, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountID.Name, storageAccountResourceName)

	storageAccount, err := storageClient.GetProperties(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): %+v", storageAccountID.Name, storageAccountID.ResourceGroup, err)
	}
	if storageAccount.AccountProperties == nil {
		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): `properties` was nil", storageAccountID.Name, storageAccountID.ResourceGroup)
	}

	// since we're mutating the storage account here, we can use that as the ID
	resourceID := storageAccountIDRaw

	if d.IsNewResource() {
		// whilst this looks superflurious given encryption is enabled by default, due to the way
		// the Azure API works this technically can be nil
		if storageAccount.AccountProperties.Encryption != nil {
			if storageAccount.AccountProperties.Encryption.KeySource == storage.KeySourceMicrosoftKeyvault {
				return tf.ImportAsExistsError("azurerm_storage_account_customer_managed_key", resourceID)
			}
		}
	}

	keyVaultIDRaw := d.Get("key_vault_id").(string)
	keyVaultID, err := keyVaultParse.VaultID(keyVaultIDRaw)
	if err != nil {
		return err
	}

	keyVault, err := vaultsClient.Get(ctx, keyVaultID.ResourceGroup, keyVaultID.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", keyVaultID.Name, keyVaultID.ResourceGroup, err)
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
		return fmt.Errorf("Key Vault %q (Resource Group %q) must be configured for both Purge Protection and Soft Delete", keyVaultID.Name, keyVaultID.ResourceGroup)
	}

	keyVaultBaseURL, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultsClient, keyVaultIDRaw)
	if err != nil {
		return fmt.Errorf("Error looking up Key Vault URI from Key Vault %q (Resource Group %q): %+v", keyVaultID.Name, keyVaultID.ResourceGroup, err)
	}

	keyName := d.Get("key_name").(string)
	keyVersion := d.Get("key_version").(string)

	props := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			Encryption: &storage.Encryption{
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						Enabled: utils.Bool(true),
					},
					File: &storage.EncryptionService{
						Enabled: utils.Bool(true),
					},
				},
				KeySource: storage.KeySourceMicrosoftKeyvault,
				KeyVaultProperties: &storage.KeyVaultProperties{
					KeyName:     utils.String(keyName),
					KeyVersion:  utils.String(keyVersion),
					KeyVaultURI: utils.String(keyVaultBaseURL),
				},
			},
		},
	}

	if _, err = storageClient.Update(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, props); err != nil {
		return fmt.Errorf("Error updating Customer Managed Key for Storage Account %q (Resource Group %q): %+v", storageAccountID.Name, storageAccountID.ResourceGroup, err)
	}

	d.SetId(resourceID)
	return resourceArmStorageAccountCustomerManagedKeyRead(d, meta)
}

func resourceArmStorageAccountCustomerManagedKeyRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	vaultsClient := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountID, err := storageParse.AccountID(d.Id())
	if err != nil {
		return err
	}

	storageAccount, err := storageClient.GetProperties(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			log.Printf("[DEBUG] Storage Account %q could not be found in Resource Group %q - removing from state!", storageAccountID.Name, storageAccountID.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): %+v", storageAccountID.Name, storageAccountID.ResourceGroup, err)
	}
	if storageAccount.AccountProperties == nil {
		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): `properties` was nil", storageAccountID.Name, storageAccountID.ResourceGroup)
	}
	if storageAccount.AccountProperties.Encryption == nil || storageAccount.AccountProperties.Encryption.KeySource != storage.KeySourceMicrosoftKeyvault {
		log.Printf("[DEBUG] Customer Managed Key was not defined for Storage Account %q (Resource Group %q) - removing from state!", storageAccountID.Name, storageAccountID.ResourceGroup)
		d.SetId("")
		return nil
	}

	encryption := *storageAccount.AccountProperties.Encryption

	keyName := ""
	keyVaultURI := ""
	keyVersion := ""
	if props := encryption.KeyVaultProperties; props != nil {
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
		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): `properties.encryption.keyVaultProperties.keyVaultURI` was nil", storageAccountID.Name, storageAccountID.ResourceGroup)
	}

	keyVaultID, err := azure.GetKeyVaultIDFromBaseUrl(ctx, vaultsClient, keyVaultURI)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault ID from the Base URI %q: %+v", keyVaultURI, err)
	}

	// now we have the key vault uri we can look up the ID

	d.Set("storage_account_id", d.Id())
	d.Set("key_vault_id", keyVaultID)
	d.Set("key_name", keyName)
	d.Set("key_version", keyVersion)

	return nil
}

func resourceArmStorageAccountCustomerManagedKeyDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountID, err := storageParse.AccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(storageAccountID.Name, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountID.Name, storageAccountResourceName)

	// confirm it still exists prior to trying to update it, else we'll get an error
	storageAccount, err := storageClient.GetProperties(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return nil
		}

		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): %+v", storageAccountID.Name, storageAccountID.ResourceGroup, err)
	}

	// Since this isn't a real object, just modifying an existing object
	// "Delete" doesn't really make sense it should really be a "Revert to Default"
	// So instead of the Delete func actually deleting the Storage Account I am
	// making it reset the Storage Account to its default state
	props := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			Encryption: &storage.Encryption{
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						Enabled: utils.Bool(true),
					},
					File: &storage.EncryptionService{
						Enabled: utils.Bool(true),
					},
				},
				KeySource: storage.KeySourceMicrosoftStorage,
			},
		},
	}

	if _, err = storageClient.Update(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, props); err != nil {
		return fmt.Errorf("Error removing Customer Managed Key for Storage Account %q (Resource Group %q): %+v", storageAccountID.Name, storageAccountID.ResourceGroup, err)
	}

	return nil
}
