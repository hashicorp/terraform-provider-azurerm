package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var storageAccountBlobSettingsResourceName = "azurerm_storage_account_blob_settings"

func resourceArmStorageAccountBlobSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageAccountBlobSettingsCreateUpdate,
		Read:   resourceArmStorageAccountBlobSettingsRead,
		Update: resourceArmStorageAccountBlobSettingsCreateUpdate,
		Delete: resourceArmStorageAccountBlobSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enable_soft_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"soft_delete_retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(1, 365),
			},
		},
	}
}

func resourceArmStorageAccountBlobSettingsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storageBlobServicesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Storage Account Blob Settings creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	enableSoftDelete := d.Get("enable_soft_delete").(bool)
	softDeleteRetentionDays := int32(d.Get("soft_delete_retention_days").(int))

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.GetServiceProperties(ctx, resourceGroupName, storageAccountName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("error checking for presence of existing Storage Account Blob Settings (Storage Account %q / Resource Group %q): %s",
					storageAccountName, resourceGroupName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError(storageAccountBlobSettingsResourceName, *existing.ID)
		}
	}

	properties := storage.BlobServicePropertiesProperties{
		DeleteRetentionPolicy: &storage.DeleteRetentionPolicy{
			Enabled: &enableSoftDelete,
			Days:    &softDeleteRetentionDays,
		},
	}

	propertiesWrapper := storage.BlobServiceProperties{
		BlobServicePropertiesProperties: &properties,
	}

	_, err := client.SetServiceProperties(ctx, resourceGroupName, storageAccountName, propertiesWrapper)

	if err != nil {
		return fmt.Errorf("error while updataing Storage Account Blob Settings (Storage Account %q / Resource Group %q): %s",
			storageAccountName, resourceGroupName, err)
	}

	read, err := client.GetServiceProperties(ctx, resourceGroupName, storageAccountName)

	if err != nil || read.ID == nil {
		return fmt.Errorf("can't read ID of Storage Account Blob Settings (Storage Account %q / Resource Group %q): %s",
			storageAccountName, resourceGroupName, err)
	}

	d.SetId(*read.ID)

	return resourceArmStorageAccountBlobSettingsRead(d, meta)
}

func resourceArmStorageAccountBlobSettingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storageBlobServicesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] Reading Azure ARM Storage Account Blob Settings.")

	id, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	resourceGroupName := id.ResourceGroup
	storageAccountName := id.Path["storageAccounts"]

	resp, err := client.GetServiceProperties(ctx, resourceGroupName, storageAccountName)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error making Read request on Azure Storage Account %q: %+v", storageAccountName, err)
	}

	d.Set("resource_group_name", resourceGroupName)
	d.Set("storage_account_name", storageAccountName)

	if props := resp.BlobServicePropertiesProperties; props != nil {
		if policy := props.DeleteRetentionPolicy; policy != nil {
			if policy.Enabled != nil {
				d.Set("enable_soft_delete", policy.Enabled)
			}
			if policy.Days != nil {
				d.Set("soft_delete_retention_days", policy.Days)
			}
		}
	}

	return nil
}

func resourceArmStorageAccountBlobSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storageBlobServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroupName := id.ResourceGroup
	storageAccountName := id.Path["storageAccounts"]
	enableSoftDelete := false
	propertiesWrapper := storage.BlobServiceProperties{
		BlobServicePropertiesProperties: &storage.BlobServicePropertiesProperties{
			DeleteRetentionPolicy: &storage.DeleteRetentionPolicy{
				Enabled: &enableSoftDelete,
			},
		},
	}

	_, err = client.SetServiceProperties(ctx, resourceGroupName, storageAccountName, propertiesWrapper)

	if err != nil {
		return fmt.Errorf("error deleting Storage Account Blob Settings (Storage Account %q / Resource Group %q): %s",
			storageAccountName, resourceGroupName, err)
	}

	return nil
}
