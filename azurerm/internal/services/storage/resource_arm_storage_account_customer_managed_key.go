package storage

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageAccountCustomerManagedKey() *schema.Resource {
	return &schema.Resource{
		Read:          resourceArmStorageAccountCustomerManagedKeyRead,
		Create:        resourceArmStorageAccountCustomerManagedKeyCreateUpdate,
		Update:        resourceArmStorageAccountCustomerManagedKeyCreateUpdate,
		Delete:        resourceArmStorageAccountCustomerManagedKeyDelete,
		SchemaVersion: 2,
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
				ValidateFunc: azure.ValidateResourceID,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"key_vault_access_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"key_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"key_vault_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmStorageAccountCustomerManagedKeyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyVaultId := d.Get("key_vault_id").(string)
	keyName := d.Get("key_name").(string)
	keyVersion := d.Get("key_version").(string)
	storageAccountId := d.Get("storage_account_id").(string)
	encryptionServices := true

	storageAccountName, storageAccountResourceGroupName, err := azure.KeyVaultGetResourceNameResource(storageAccountId, "storageAccounts")
	if err != nil {
		return err
	}

	keyVaultAccountName, keyVaultResourceGroupName, err := azure.KeyVaultGetResourceNameResource(keyVaultId, "vaults")
	if err != nil {
		return err
	}

	// First check to see if the key vault is configured correctly or not
	if !azure.KeyVaultIsSoftDeleteAndPurgeProtected(ctx, vaultClient, keyVaultId) {
		return fmt.Errorf("Key Vault %q (Resource Group %q) is not configured correctly, please make sure that both 'soft_delete_enabled' and 'purge_protection_enabled' arguments are set to 'true'", keyVaultAccountName, keyVaultResourceGroupName)
	}

	pKeyVaultBaseURL, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Key Vault URI from Key Vault %q (Resource Group %q): %+v", keyVaultAccountName, keyVaultResourceGroupName, err)
	}

	props := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			Encryption: &storage.Encryption{
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						Enabled: utils.Bool(encryptionServices),
					},
					File: &storage.EncryptionService{
						Enabled: utils.Bool(encryptionServices),
					},
				},
				KeySource: storage.MicrosoftKeyvault,
				KeyVaultProperties: &storage.KeyVaultProperties{
					KeyName:     utils.String(keyName),
					KeyVersion:  utils.String(keyVersion),
					KeyVaultURI: utils.String(pKeyVaultBaseURL),
				},
			},
		},
	}

	_, err = storageClient.Update(ctx, storageAccountResourceGroupName.(string), storageAccountName.(string), props)
	if err != nil {
		return fmt.Errorf("Error updating Azure Storage Account %q (Resource Group %q) Customer Managed Key : %+v", storageAccountName, storageAccountResourceGroupName, err)
	}

	resourceId := fmt.Sprintf("%s/customerManagedKey", storageAccountId)
	d.SetId(resourceId)

	return resourceArmStorageAccountCustomerManagedKeyRead(d, meta)
}

func resourceArmStorageAccountCustomerManagedKeyRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountId := d.Get("storage_account_id").(string)
	keyVaultId := d.Get("key_vault_id").(string)
	keyVaultPolicy := d.Get("key_vault_access_policy_id").(string)

	d.Set("storage_account_id", storageAccountId)
	d.Set("key_vault_id", keyVaultId)
	d.Set("key_vault_access_policy_id", keyVaultPolicy)

	name, resGroup, err := azure.KeyVaultGetResourceNameResource(storageAccountId, "storageAccounts")
	if err != nil {
		return err
	}

	resp, err := storageClient.GetProperties(ctx, resGroup.(string), name.(string), "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading the state of AzureRM Storage Account %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if props := resp.AccountProperties; props != nil {
		if encryption := props.Encryption; encryption != nil {
			if keyVaultProperties := encryption.KeyVaultProperties; keyVaultProperties != nil {
				d.Set("key_name", *keyVaultProperties.KeyName)
				d.Set("key_version", *keyVaultProperties.KeyVersion)
				d.Set("key_vault_uri", *keyVaultProperties.KeyVaultURI)
			}
		}
	}

	return nil
}

func resourceArmStorageAccountCustomerManagedKeyDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountId := d.Get("storage_account_id").(string)

	storageAccountName, resourceGroupName, err := azure.KeyVaultGetResourceNameResource(storageAccountId, "storageAccounts")
	if err != nil {
		return err
	}

	// Since this isn't a real object, just modifying an existing object
	// "Delete" doesn't really make sense it should really be a "Revert to Default"
	// So instead of the Delete func actually deleting the Storage Account I am
	// making it reset the Storage Account to it's default state
	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			Encryption: &storage.Encryption{
				KeySource: storage.MicrosoftStorage,
			},
		},
	}

	_, err = storageClient.Update(ctx, resourceGroupName.(string), storageAccountName.(string), opts)
	if err != nil {
		return fmt.Errorf("Error deleting AzureRM Storage Account %q (Resource Group %q) %q: %+v", storageAccountName.(string), resourceGroupName.(string), err)
	}

	return nil
}

func resourceArmStorageAccountCustomerManagedKeyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()

	d.Set("storage_account_id", id)
	d.Set("key_vault_id", "")
	d.Set("key_vault_access_policy_id", "")

	name, resGroup, err := azure.KeyVaultGetResourceNameResource(id, "storageAccounts")
	if err != nil {
		return nil, err
	}

	resp, err := storageClient.GetProperties(ctx, resGroup.(string), name.(string), "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil, nil
		}
		return nil, fmt.Errorf("Error importing the state of AzureRM Storage Account %q (Resource Group %q): %+v", name.(string), resGroup.(string), err)
	}

	if props := resp.AccountProperties; props != nil {
		if encryption := props.Encryption; encryption != nil {
			if keyVaultProperties := encryption.KeyVaultProperties; keyVaultProperties != nil {
				d.Set("key_name", keyVaultProperties.KeyName)
				d.Set("key_version", keyVaultProperties.KeyVersion)
				d.Set("key_vault_uri", keyVaultProperties.KeyVaultURI)
			}
		}
	}

	resourceId := fmt.Sprintf("%s/customerManagedKey", id)
	d.SetId(resourceId)

	results := make([]*schema.ResourceData, 1)

	results[0] = d
	return results, nil
}
