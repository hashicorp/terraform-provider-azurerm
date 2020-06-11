package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/datalakestore/paths"
)

func resourceArmStorageDataLakeGen2Path() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageDataLakeGen2PathCreate,
		Read:   resourceArmStorageDataLakeGen2PathRead,
		// Update: resourceArmStorageDataLakeGen2PathUpdate,
		Delete: resourceArmStorageDataLakeGen2PathDelete,

		// TODO
		// Importer: &schema.ResourceImporter{
		// 	State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		// 		storageClients := meta.(*clients.Client).Storage
		// 		ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, 5*time.Minute)
		// 		defer cancel()

		// 		id, err := filesystems.ParseResourceID(d.Id())
		// 		if err != nil {
		// 			return []*schema.ResourceData{d}, fmt.Errorf("Error parsing ID %q for import of Data Lake Gen2 File System: %v", d.Id(), err)
		// 		}

		// 		// we then need to look up the Storage Account ID
		// 		account, err := storageClients.FindAccount(ctx, id.AccountName)
		// 		if err != nil {
		// 			return []*schema.ResourceData{d}, fmt.Errorf("Error retrieving Account %q for Data Lake Gen2 File System %q: %s", id.AccountName, id.DirectoryName, err)
		// 		}
		// 		if account == nil {
		// 			return []*schema.ResourceData{d}, fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
		// 		}

		// 		d.Set("storage_account_id", account.ID)

		// 		return []*schema.ResourceData{d}, nil
		// 	},
		// },

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"storage_account_id": AccountIDSchema(),

			"filesystem_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageDataLakeGen2FileSystemName,
			},

			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // TODO handle rename
			},

			"resource": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"directory"}, false),
			},
		},
	}
}

func resourceArmStorageDataLakeGen2PathCreate(d *schema.ResourceData, meta interface{}) error {
	accountsClient := meta.(*clients.Client).Storage.AccountsClient
	client := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageID, err := parsers.ParseAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	// confirm the storage account exists, otherwise Data Plane API requests will fail
	storageAccount, err := accountsClient.GetProperties(ctx, storageID.ResourceGroup, storageID.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("Storage Account %q was not found in Resource Group %q!", storageID.Name, storageID.ResourceGroup)
		}

		return fmt.Errorf("Error checking for existence of Storage Account %q (Resource Group %q): %+v", storageID.Name, storageID.ResourceGroup, err)
	}

	fileSystemName := d.Get("filesystem_name").(string)
	path := d.Get("path").(string)

	resourceString := d.Get("resource").(string)
	var resource paths.PathResource
	switch resourceString {
	case "directory":
		resource = paths.PathResourceDirectory
	default:
		return fmt.Errorf("Unhandled resource type %q", resourceString)
	}

	id := client.GetResourceID(storageID.Name, fileSystemName, path)

	if features.ShouldResourcesBeImported() {
		resp, err := client.GetProperties(ctx, storageID.Name, fileSystemName, path)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for existence of existing Path %q in  File System %q (Account %q): %+v", path, fileSystemName, storageID.Name, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_storage_data_lake_gen2_path", id)
		}
	}

	log.Printf("[INFO] Creating Path %q in File System %q in Storage Account %q.", path, fileSystemName, storageID.Name)
	input := paths.CreateInput{
		Resource: resource,
	}

	if _, err := client.Create(ctx, storageID.Name, fileSystemName, path, input); err != nil {
		return fmt.Errorf("Error creating Path %q in File System %q in Storage Account %q: %s", path, fileSystemName, storageID.Name, err)
	}

	d.SetId(id)
	return resourceArmStorageDataLakeGen2PathRead(d, meta)
}

func resourceArmStorageDataLakeGen2PathUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Not implemented - update")
	// accountsClient := meta.(*clients.Client).Storage.AccountsClient
	// client := meta.(*clients.Client).Storage.FileSystemsClient
	// ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	// defer cancel()

	// id, err := filesystems.ParseResourceID(d.Id())
	// if err != nil {
	// 	return err
	// }

	// storageID, err := parsers.ParseAccountID(d.Get("storage_account_id").(string))
	// if err != nil {
	// 	return err
	// }

	// // confirm the storage account exists, otherwise Data Plane API requests will fail
	// storageAccount, err := accountsClient.GetProperties(ctx, storageID.ResourceGroup, storageID.Name, "")
	// if err != nil {
	// 	if utils.ResponseWasNotFound(storageAccount.Response) {
	// 		return fmt.Errorf("Storage Account %q was not found in Resource Group %q!", storageID.Name, storageID.ResourceGroup)
	// 	}

	// 	return fmt.Errorf("Error checking for existence of Storage Account %q (Resource Group %q): %+v", storageID.Name, storageID.ResourceGroup, err)
	// }

	// propertiesRaw := d.Get("properties").(map[string]interface{})
	// properties := ExpandMetaData(propertiesRaw)

	// log.Printf("[INFO] Updating Properties for File System %q in Storage Account %q.", id.DirectoryName, id.AccountName)
	// input := filesystems.SetPropertiesInput{
	// 	Properties: properties,
	// }
	// if _, err = client.SetProperties(ctx, id.AccountName, id.DirectoryName, input); err != nil {
	// 	return fmt.Errorf("Error updating Properties for File System %q in Storage Account %q: %s", id.DirectoryName, id.AccountName, err)
	// }

	// return resourceArmStorageDataLakeGen2FileSystemRead(d, meta)
}

func resourceArmStorageDataLakeGen2PathRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := paths.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, id.AccountName, id.FileSystemName, id.Path)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Path %q does not exist in File System %q in Storage Account %q - removing from state...", id.Path, id.FileSystemName, id.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Path %q in File System %q in Storage Account %q: %+v", id.Path, id.FileSystemName, id.AccountName, err)
	}

	d.Set("path", id.Path)

	d.Set("resource", resp.ResourceType)

	return nil
}

func resourceArmStorageDataLakeGen2PathDelete(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Not implemented - Delete")
	// client := meta.(*clients.Client).Storage.FileSystemsClient
	// ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	// defer cancel()

	// id, err := filesystems.ParseResourceID(d.Id())
	// if err != nil {
	// 	return err
	// }

	// resp, err := client.Delete(ctx, id.AccountName, id.DirectoryName)
	// if err != nil {
	// 	if !utils.ResponseWasNotFound(resp) {
	// 		return fmt.Errorf("Error deleting File System %q in Storage Account %q: %+v", id.DirectoryName, id.AccountName, err)
	// 	}
	// }

	// return nil
}
