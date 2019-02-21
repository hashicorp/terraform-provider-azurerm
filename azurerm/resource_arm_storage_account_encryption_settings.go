package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageAccountEncryptionSettings() *schema.Resource {
	return &schema.Resource{
		Read:   resourceArmStorageAccountEncryptionSettingsRead,
		Create: resourceArmStorageAccountEncryptionSettingsCreateUpdate,
		Update: resourceArmStorageAccountEncryptionSettingsCreateUpdate,
		Delete: resourceArmStorageAccountEncryptionSettingsDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 2,

		Schema: map[string]*schema.Schema{
			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"key_vault_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
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

			"account_encryption_source": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.MicrosoftKeyvault),
					string(storage.MicrosoftStorage),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"key_vault_properties": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func resourceArmStorageAccountEncryptionSettingsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	vaultClient := meta.(*ArmClient).keyVaultClient
	client := meta.(*ArmClient).storageServiceClient

	storageAccountId := d.Get("storage_account_id").(string)

	id, err := azure.ParseAzureResourceID(storageAccountId)
	if err != nil {
		return err
	}

	storageAccountName := id.Path["storageAccounts"]
	resourceGroupName := id.ResourceGroup
	enableBlobEncryption := d.Get("enable_blob_encryption").(bool)
	enableFileEncryption := d.Get("enable_file_encryption").(bool)
	encryptionSource := d.Get("account_encryption_source").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	pKeyVaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)

	if err != nil {
		return fmt.Errorf("Error looking up Key Vault URI from id %q: %+v", keyVaultId, err)
	}

	keyVaultProperties := expandAzureRmStorageAccountKeyVaultProperties(d)
	keyVaultProperties.KeyVaultURI = utils.String(pKeyVaultBaseUrl)

	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			Encryption: &storage.Encryption{
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						Enabled: utils.Bool(enableBlobEncryption),
					},
					File: &storage.EncryptionService{
						Enabled: utils.Bool(enableFileEncryption),
					}},
				KeySource:          storage.KeySource(encryptionSource),
				KeyVaultProperties: keyVaultProperties,
			},
		},
	}

	_, err = client.Update(ctx, resourceGroupName, storageAccountName, opts)
	if err != nil {
		return fmt.Errorf("Error updating Azure Storage Account Encryption %q: %+v", storageAccountName, err)
	}

	resourceId := fmt.Sprintf("%s/encryptionSettings", storageAccountId)
	d.SetId(resourceId)

	return resourceArmStorageAccountEncryptionSettingsRead(d, meta)
}

func resourceArmStorageAccountEncryptionSettingsRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).storageServiceClient

	storageAccountId := d.Get("storage_account_id").(string)

	id, err := parseAzureResourceID(storageAccountId)
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

	if props := resp.AccountProperties; props != nil {
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

			if keyVaultProperties := encryption.KeyVaultProperties; keyVaultProperties != nil {
				if err := d.Set("key_vault_properties", flattenAzureRmStorageAccountKeyVaultProperties(keyVaultProperties)); err != nil {
					return fmt.Errorf("Error flattening `key_vault_properties`: %+v", err)
				}
			}
		}
	}

	return nil
}

func resourceArmStorageAccountEncryptionSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).storageServiceClient

	storageAccountId := d.Get("storage_account_id").(string)

	id, err := azure.ParseAzureResourceID(storageAccountId)
	if err != nil {
		return err
	}

	storageAccountName := id.Path["storageAccounts"]
	resourceGroupName := id.ResourceGroup

	// Since this isn't a real object, just modifying an existing object
	// "Delete" doesn't really make sense it should really be a "Revert to Default"
	// So instead of the Delete func actually deleting the Storage Account I am
	// making it reset the Storage Account to it's default state
	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			Encryption: &storage.Encryption{
				KeySource: storage.KeySource(storage.MicrosoftStorage),
			},
		},
	}

	_, err = client.Update(ctx, resourceGroupName, storageAccountName, opts)
	if err != nil {
		return fmt.Errorf("Error updating Azure Storage Account Encryption %q: %+v", storageAccountName, err)
	}

	return nil
}

func expandAzureRmStorageAccountKeyVaultProperties(d *schema.ResourceData) *storage.KeyVaultProperties {
	vs := d.Get("key_vault_properties").([]interface{})
	if vs == nil || len(vs) == 0 {
		return &storage.KeyVaultProperties{}
	}

	v := vs[0].(map[string]interface{})
	keyName := v["key_name"].(string)
	keyVersion := v["key_version"].(string)

	return &storage.KeyVaultProperties{
		KeyName:    utils.String(keyName),
		KeyVersion: utils.String(keyVersion),
	}
}

func flattenAzureRmStorageAccountKeyVaultProperties(keyVaultProperties *storage.KeyVaultProperties) []interface{} {
	if keyVaultProperties == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	if keyVaultProperties.KeyName != nil {
		result["key_name"] = *keyVaultProperties.KeyName
	}
	if keyVaultProperties.KeyVersion != nil {
		result["key_version"] = *keyVaultProperties.KeyVersion
	}
	if keyVaultProperties.KeyVaultURI != nil {
		result["key_vault_uri"] = *keyVaultProperties.KeyVaultURI
	}

	return []interface{}{result}
}
