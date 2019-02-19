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

func resourceArmStorageAccountKeyVaultCustomKeys() *schema.Resource {
	return &schema.Resource{
		Read:   resourceArmStorageAccountVaultCustomKeysRead,
		Create: resourceArmStorageAccountVaultCustomKeysUpdate,
		Update: resourceArmStorageAccountVaultCustomKeysUpdate,
		Delete: resourceArmStorageAccountVaultCustomKeysDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceStorageAccountMigrateState,
		SchemaVersion: 2,

		Schema: map[string]*schema.Schema{
			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enable_blob_encryption": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"enable_file_encryption": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"account_encryption_source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(storage.MicrosoftKeyvault),
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.MicrosoftKeyvault),
					string(storage.MicrosoftStorage),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"key_vault_properties": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},
		},
	}
}

// resourceArmStorageAccountUpdate is unusual in the ARM API where most resources have a combined
// and idempotent operation for CreateOrUpdate. In particular updating all of the parameters
// available requires a call to Update per parameter...
func resourceArmStorageAccountVaultCustomKeysUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).storageServiceClient

	storageAccountId := d.Get("storage_account_id").(string)

	id, err := azure.ParseAzureResourceID(storageAccountId)
	if err != nil {
		return err
	}

	storageAccountName := id.Path["storageAccounts"]
	resourceGroupName := id.ResourceGroup
	encryptionSource := d.Get("account_encryption_source").(string)

	d.Partial(true)

	if d.HasChange("enable_blob_encryption") || d.HasChange("enable_file_encryption") {

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				Encryption: &storage.Encryption{
					Services:  &storage.EncryptionServices{},
					KeySource: storage.KeySource(storage.MicrosoftStorage),
				},
			},
		}

		if d.HasChange("enable_blob_encryption") {
			enableEncryption := d.Get("enable_blob_encryption").(bool)
			opts.Encryption.Services.Blob = &storage.EncryptionService{
				Enabled: utils.Bool(enableEncryption),
			}

			_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
			if err != nil {
				return fmt.Errorf("Error updating Azure Storage Account Encryption %q: %+v", storageAccountName, err)
			}

			d.SetPartial("enable_blob_encryption")
		}

		if d.HasChange("enable_file_encryption") {
			enableEncryption := d.Get("enable_file_encryption").(bool)
			opts.Encryption.Services.File = &storage.EncryptionService{
				Enabled: utils.Bool(enableEncryption),
			}

			_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
			if err != nil {
				return fmt.Errorf("Error updating Azure Storage Account Encryption %q: %+v", storageAccountName, err)
			}

			d.SetPartial("enable_file_encryption")
		}
	}

	// NOTE: If KeySource is KeyVault you also need to send key vault properties
	if encryptionSource == "Microsoft.Keyvault" {
		keyVaultProperties := expandAzureRmStorageAccountKeyVaultProperties(d)

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				Encryption: &storage.Encryption{
					KeySource:          storage.KeySource(storage.MicrosoftKeyvault),
					KeyVaultProperties: keyVaultProperties,
				},
			},
		}

		_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
		if err != nil {
			return fmt.Errorf("Error updating Azure Storage Account Encryption %q: %+v", storageAccountName, err)
		}

		d.SetPartial("account_encryption_source")
		d.SetPartial("key_vault_properties")
	}

	d.Partial(false)

	return resourceArmStorageAccountVaultCustomKeysRead(d, meta)
}

func resourceArmStorageAccountVaultCustomKeysRead(d *schema.ResourceData, meta interface{}) error {
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

	// keys, err := client.ListKeys(ctx, resGroup, name)
	// if err != nil {
	// 	return err
	// }

	//accessKeys := *keys.Keys
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
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

func resourceArmStorageAccountVaultCustomKeysDelete(d *schema.ResourceData, meta interface{}) error {
	// Thnik what to do here?
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
	keyVaultURI := v["key_vault_uri"].(string)

	return &storage.KeyVaultProperties{
		KeyName:     utils.String(keyName),
		KeyVersion:  utils.String(keyVersion),
		KeyVaultURI: utils.String(keyVaultURI),
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
