package azurerm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataLakeStoreFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataLakeStoreFileCreate,
		Read:   resourceArmDataLakeStoreFileRead,
		Delete: resourceArmDataLakeStoreFileDelete,

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"remote_file_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateFilePath(),
			},

			"local_file_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmDataLakeStoreFileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreFilesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Date Lake Store File creation.")

	accountName := d.Get("account_name").(string)
	remoteFilePath := d.Get("remote_file_path").(string)
	localFilePath := d.Get("local_file_path").(string)

	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("error opening file %q: %+v", localFilePath, err)
	}
	defer file.Close()

	// Read the file contents
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = client.Create(ctx, accountName, remoteFilePath, ioutil.NopCloser(bytes.NewReader(fileContents)), utils.Bool(false), filesystem.CLOSE, nil, nil)
	if err != nil {
		return fmt.Errorf("Error issuing create request for Data Lake Store File %q : %+v", remoteFilePath, err)
	}

	d.SetId(remoteFilePath)

	return resourceArmDataLakeStoreFileRead(d, meta)
}

func resourceArmDataLakeStoreFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreFilesClient
	ctx := meta.(*ArmClient).StopContext

	accountName := d.Get("account_name").(string)
	remoteFilePath := d.Id()

	resp, err := client.GetFileStatus(ctx, accountName, remoteFilePath, utils.Bool(true))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Data Lake Store File %q was not found (Account %q)", remoteFilePath, accountName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Data Lake Store File %q (Account %q): %+v", remoteFilePath, accountName, err)
	}

	return nil
}

func resourceArmDataLakeStoreFileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreFilesClient
	ctx := meta.(*ArmClient).StopContext

	accountName := d.Get("account_name").(string)
	remoteFilePath := d.Id()
	resp, err := client.Delete(ctx, accountName, remoteFilePath, utils.Bool(false))
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request for Data Lake Store File %q (Account %q): %+v", remoteFilePath, accountName, err)
	}

	return nil
}
