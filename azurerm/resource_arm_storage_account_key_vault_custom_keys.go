package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageAccountKeyVaultCustomKeys() *schema.Resource {
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

			"account_encryption_source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(storage.MicrosoftStorage),
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

	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	enableBlobEncryption := d.Get("enable_blob_encryption").(bool)
	enableFileEncryption := d.Get("enable_file_encryption").(bool)
	keyVaultProperties := expandAzureRmStorageAccountKeyVaultProperties(d)

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
				KeySource:          storage.KeySource(storageAccountEncryptionSource),
				KeyVaultProperties: keyVaultProperties,
			},
			EnableHTTPSTrafficOnly: &enableHTTPSTrafficOnly,
			NetworkRuleSet:         networkRules,
		},
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

	d.Partial(true)

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

	if d.HasChange("account_encryption_source") {
		accountEncryptionSource := d.Get("account_encryption_source").(string)
		keyVaultProperties := expandAzureRmStorageAccountKeyVaultProperties(d)
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				Encryption: &storage.Encryption{
					KeySource:          storage.KeySource(accountEncryptionSource),
					KeyVaultProperties: keyVaultProperties,
				},
			},
		}

		_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
		if err != nil {
			return fmt.Errorf("Error updating Azure Storage Account account_encryption_source and key_vault_properties %q: %+v", storageAccountName, err)
		}

		d.SetPartial("account_encryption_source")
		d.SetPartial("key_vault_properties")
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

func resourceArmStorageAccountDelete(d *schema.ResourceData, meta interface{}) error {
	// Not sure what to do here...

	// ctx := meta.(*ArmClient).StopContext
	// client := meta.(*ArmClient).storageServiceClient

	// id, err := parseAzureResourceID(d.Id())
	// if err != nil {
	// 	return err
	// }
	// name := id.Path["storageAccounts"]
	// resourceGroup := id.ResourceGroup

	// read, err := client.GetProperties(ctx, resourceGroup, name)
	// if err != nil {
	// 	if utils.ResponseWasNotFound(read.Response) {
	// 		return nil
	// 	}

	// 	return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	// }

	// resp, err := client.Delete(ctx, resourceGroup, name)
	// if err != nil {
	// 	if !response.WasNotFound(resp.Response) {
	// 		return fmt.Errorf("Error issuing delete request for Storage Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	// 	}
	// }

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
