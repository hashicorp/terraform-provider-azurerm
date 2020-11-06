package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/directories"
)

func resourceArmStorageShareDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageShareDirectoryCreate,
		Read:   resourceArmStorageShareDirectoryRead,
		Update: resourceArmStorageShareDirectoryUpdate,
		Delete: resourceArmStorageShareDirectoryDelete,
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageShareDirectoryName,
			},
			"share_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"metadata": MetaDataSchema(),
		},
	}
}

func resourceArmStorageShareDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	accountName := d.Get("storage_account_name").(string)
	shareName := d.Get("share_name").(string)
	directoryName := d.Get("name").(string)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Directory %q (Share %q): %s", accountName, directoryName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.FileShareDirectoriesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Directories Client: %s", err)
	}

	existing, err := client.Get(ctx, accountName, shareName, directoryName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Directory %q (File Share %q / Storage Account %q / Resource Group %q): %s", directoryName, shareName, accountName, account.ResourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		id := client.GetResourceID(accountName, shareName, directoryName)
		return tf.ImportAsExistsError("azurerm_storage_share_directory", id)
	}

	if _, err := client.Create(ctx, accountName, shareName, directoryName, metaData); err != nil {
		return fmt.Errorf("Error creating Directory %q (File Share %q / Account %q): %+v", directoryName, shareName, accountName, err)
	}

	// Storage Share Directories are eventually consistent
	log.Printf("[DEBUG] Waiting for Directory %q (File Share %q / Account %q) to become available", directoryName, shareName, accountName)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   storageShareDirectoryRefreshFunc(ctx, client, accountName, shareName, directoryName),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   d.Timeout(schema.TimeoutCreate),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Directory %q (File Share %q / Account %q) to become available: %s", directoryName, shareName, accountName, err)
	}

	resourceID := client.GetResourceID(accountName, shareName, directoryName)
	d.SetId(resourceID)

	return resourceArmStorageShareDirectoryRead(d, meta)
}

func resourceArmStorageShareDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := directories.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Directory %q (Share %q): %s", id.AccountName, id.DirectoryName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileShareDirectoriesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client: %s", err)
	}

	if _, err := client.SetMetaData(ctx, id.AccountName, id.ShareName, id.DirectoryName, metaData); err != nil {
		return fmt.Errorf("Error updating MetaData for Directory %q (File Share %q / Account %q): %+v", id.DirectoryName, id.ShareName, id.AccountName, err)
	}

	return resourceArmStorageShareDirectoryRead(d, meta)
}

func resourceArmStorageShareDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := directories.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Directory %q (Share %q): %s", id.AccountName, id.DirectoryName, id.ShareName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Share Directory %q (Share %s, Account %s) - assuming removed & removing from state", id.DirectoryName, id.ShareName, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileShareDirectoriesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	props, err := client.Get(ctx, id.AccountName, id.ShareName, id.DirectoryName)
	if err != nil {
		return fmt.Errorf("Error retrieving Storage Share %q (File Share %q / Account %q / Resource Group %q): %s", id.DirectoryName, id.ShareName, id.AccountName, account.ResourceGroup, err)
	}

	d.Set("name", id.DirectoryName)
	d.Set("share_name", id.ShareName)
	d.Set("storage_account_name", id.AccountName)

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %s", err)
	}

	return nil
}

func resourceArmStorageShareDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := directories.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Directory %q (Share %q): %s", id.AccountName, id.DirectoryName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileShareDirectoriesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	if _, err := client.Delete(ctx, id.AccountName, id.ShareName, id.DirectoryName); err != nil {
		return fmt.Errorf("Error deleting Storage Share %q (File Share %q / Account %q / Resource Group %q): %s", id.DirectoryName, id.ShareName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

func storageShareDirectoryRefreshFunc(ctx context.Context, client *directories.Client, accountName, shareName, directoryName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, accountName, shareName, directoryName)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("Error retrieving Directory %q (File Share %q / Account %q): %s", directoryName, shareName, accountName, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
