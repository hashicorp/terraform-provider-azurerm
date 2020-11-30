package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/files"
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

			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ForceNew:     true,
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
		MetaData:           ExpandMetaData(d.Get("metadata").(map[string]interface{})),
		ContentType:        utils.String(d.Get("content_type").(string)),
		ContentEncoding:    utils.String(d.Get("content_encoding").(string)),
		ContentDisposition: utils.String(d.Get("content_disposition").(string)),
	}

	if v, ok := d.GetOk("content_md5"); ok {
		input.ContentMD5 = utils.String(v.(string))
	}

	var file *os.File
	if v, ok := d.GetOk("source"); ok {
		file, err = os.Open(v.(string))
		if err != nil {
			return fmt.Errorf("opening file : %s", err)
		}

		info, err := file.Stat()
		if err != nil {
			return fmt.Errorf("'stat'-ing File %q (File Share %q / Account %q): %+v", fileName, shareName, accountName, err)
		}

		input.ContentLength = info.Size()
	}

	if _, err := client.Create(ctx, accountName, shareName, directoryName, fileName, input); err != nil {
		return fmt.Errorf("creating File %q (File Share %q / Account %q): %+v", fileName, shareName, accountName, err)
	}

	if file != nil {
		if err := client.PutFile(ctx, accountName, shareName, directoryName, fileName, file, 4); err != nil {
			return fmt.Errorf("uploading File: %q (File Share %q / Account %q): %+v", fileName, shareName, accountName, err)
		}
	}

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

	if d.HasChange("content_type") || d.HasChange("content_encoding") || d.HasChange("content_disposition") || d.HasChange("content_md5") {
		input := files.SetPropertiesInput{
			ContentType:        utils.String(d.Get("content_type").(string)),
			ContentEncoding:    utils.String(d.Get("content_encoding").(string)),
			ContentDisposition: utils.String(d.Get("content_disposition").(string)),
			MetaData:           ExpandMetaData(d.Get("metadata").(map[string]interface{})),
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
