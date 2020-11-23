package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/files"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageShareFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageShareFileCreate,
		Read:   resourceArmStorageShareFileRead,
		Update: resourceArmStorageShareFileUpdate,
		Delete: resourceArmStorageShareFileDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: add validation
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
			"directory_name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "",
				ValidateFunc: validate.StorageShareDirectoryName,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "application/octet-stream",
			},

			"content_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
				// TODO check for 512 divisibility
				// ValidateFunc:
			},

			"content_encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"content_md5": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"content_disposition": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"metadata": MetaDataSchema(),
		},
	}
}

func resourceArmStorageShareFileCreate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	accountName := d.Get("storage_account_name").(string)
	shareName := d.Get("share_name").(string)
	fileName := d.Get("name").(string)
	directoryName := d.Get("directory_name").(string)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for File %q (Share %q): %s", accountName, fileName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Directories Client: %s", err)
	}

	existing, err := client.GetProperties(ctx, accountName, shareName, directoryName, fileName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing File %q (File Share %q / Storage Account %q / Resource Group %q): %s", fileName, shareName, accountName, account.ResourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		id := client.GetResourceID(accountName, shareName, directoryName, fileName)
		return tf.ImportAsExistsError("azurerm_storage_share_file", id)
	}

	input := files.CreateInput{
		MetaData:           metaData,
		ContentType:        utils.String(d.Get("content_type").(string)),
		ContentLength:      int64(d.Get("content_length").(int)),
		ContentEncoding:    utils.String(d.Get("content_encoding").(string)),
		ContentDisposition: utils.String(d.Get("content_disposition").(string)),
	}

	if v, ok := d.GetOk("content_md5"); ok {
		input.ContentMD5 = utils.String(v.(string))
	}

	if _, err := client.Create(ctx, accountName, shareName, directoryName, fileName, input); err != nil {
		return fmt.Errorf("Error creating File %q (File Share %q / Account %q): %+v", fileName, shareName, accountName, err)
	}

	// TODO Check if this is true
	/*
		// Storage Share Directories are eventually consistent
		log.Printf("[DEBUG] Waiting for File %q (File Share %q / Account %q) to become available", fileName, shareName, accountName)
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
		}*/

	resourceID := client.GetResourceID(accountName, shareName, directoryName, fileName)
	d.SetId(resourceID)

	return resourceArmStorageShareFileRead(d, meta)
}

func resourceArmStorageShareFileUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := files.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for File %q (Share %q): %s", id.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Directories Client: %s", err)
	}

	existing, err := client.GetProperties(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing File %q (File Share %q / Storage Account %q / Resource Group %q): %s", id.FileName, id.ShareName, id.AccountName, account.ResourceGroup, err)
		}
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := ExpandMetaData(metaDataRaw)

		account, err := storageClient.FindAccount(ctx, id.AccountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for File %q (Share %q): %s", id.AccountName, id.FileName, id.ShareName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
		}

		if _, err := client.SetMetaData(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName, metaData); err != nil {
			return fmt.Errorf("Error updating MetaData for File %q (File Share %q / Account %q): %+v", id.FileName, id.ShareName, id.AccountName, err)
		}
	}

	if d.HasChange("content_type") || d.HasChange("content_length") || d.HasChange("content_encoding") || d.HasChange("content_disposition") || d.HasChange("content_md5") {
		input := files.SetPropertiesInput{
			ContentType:        utils.String(d.Get("content_type").(string)),
			ContentLength:      utils.Int64(int64(d.Get("content_length").(int))),
			ContentEncoding:    utils.String(d.Get("content_encoding").(string)),
			ContentDisposition: utils.String(d.Get("content_disposition").(string)),
		}

		if v, ok := d.GetOk("content_md5"); ok {
			input.ContentMD5 = utils.String(v.(string))
		}

		if _, err := client.SetProperties(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName, input); err != nil {
			return fmt.Errorf("Error creating File %q (File Share %q / Account %q): %+v", id.FileName, id.ShareName, id.AccountName, err)
		}
	}

	return resourceArmStorageShareFileRead(d, meta)
}

func resourceArmStorageShareFileRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := files.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for File %q (Share %q): %s", id.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Share File %q (Share %s, Account %s) - assuming removed & removing from state", id.FileName, id.ShareName, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	props, err := client.GetProperties(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName)
	if err != nil {
		return fmt.Errorf("Error retrieving Storage Share %q (File Share %q / Account %q / Resource Group %q): %s", id.DirectoryName, id.ShareName, id.AccountName, account.ResourceGroup, err)
	}

	d.Set("name", id.FileName)
	d.Set("directory_name", id.DirectoryName)
	d.Set("share_name", id.ShareName)
	d.Set("storage_account_name", id.AccountName)

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %s", err)
	}
	d.Set("content_type", props.ContentType)
	d.Set("content_length", props.ContentLength)
	d.Set("content_encoding", props.ContentEncoding)
	d.Set("content_md5", props.ContentMD5)
	d.Set("content_disposition", props.ContentDisposition)

	return nil
}

func resourceArmStorageShareFileDelete(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := files.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for File %q (Share %q): %s", id.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share File Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	if _, err := client.Delete(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName); err != nil {
		return fmt.Errorf("Error deleting Storage Share File %q (File Share %q / Account %q / Resource Group %q): %s", id.FileName, id.ShareName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

/*
func storageShareDirectoryRefreshFunc(ctx context.Context, client *directories.Client, accountName, shareName, directoryName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, accountName, shareName, directoryName)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("Error retrieving Directory %q (File Share %q / Account %q): %s", directoryName, shareName, accountName, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
*/
