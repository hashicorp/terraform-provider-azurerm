package azurerm

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataLakeStoreFile() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmDataLakeStoreFileCreate,
		Read:          resourceArmDataLakeStoreFileRead,
		Delete:        resourceArmDataLakeStoreFileDelete,
		MigrateState:  resourceDataLakeStoreFileMigrateState,
		SchemaVersion: 1,
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
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"remote_file_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDataLakeStoreRemoteFilePath(),
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
	client := meta.(*clients.Client).Datalake.StoreFilesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	chunkSize := 4 * 1024 * 1024

	log.Printf("[INFO] preparing arguments for Date Lake Store File creation.")

	accountName := d.Get("account_name").(string)
	remoteFilePath := d.Get("remote_file_path").(string)
	localFilePath := d.Get("local_file_path").(string)

	// example.azuredatalakestore.net/test/example.txt
	id := fmt.Sprintf("%s.%s%s", accountName, client.AdlsFileSystemDNSSuffix, remoteFilePath)

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetFileStatus(ctx, accountName, remoteFilePath, utils.Bool(true))
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Lake Store File %q (Account %q): %s", remoteFilePath, accountName, err)
			}
		}

		if existing.FileStatus != nil && existing.FileStatus.ModificationTime != nil {
			return tf.ImportAsExistsError("azurerm_data_lake_store_file", id)
		}
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("error opening file %q: %+v", localFilePath, err)
	}
	defer func(c io.Closer) {
		if err := c.Close(); err != nil {
			log.Printf("[DEBUG] Error closing Data Lake Store File %q: %+v", localFilePath, err)
		}
	}(file)

	if _, err = client.Create(ctx, accountName, remoteFilePath, nil, nil, filesystem.DATA, nil, nil); err != nil {
		return fmt.Errorf("Error issuing create request for Data Lake Store File %q : %+v", remoteFilePath, err)
	}

	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		flag := filesystem.DATA
		if n < chunkSize {
			// last chunk
			flag = filesystem.CLOSE
		}
		chunk := ioutil.NopCloser(bytes.NewReader(buffer[:n]))

		if _, err = client.Append(ctx, accountName, remoteFilePath, chunk, nil, flag, nil, nil); err != nil {
			return fmt.Errorf("Error transferring chunk for Data Lake Store File %q : %+v", remoteFilePath, err)
		}
	}

	d.SetId(id)
	return resourceArmDataLakeStoreFileRead(d, meta)
}

func resourceArmDataLakeStoreFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreFilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseDataLakeStoreFileId(d.Id(), client.AdlsFileSystemDNSSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetFileStatus(ctx, id.storageAccountName, id.filePath, utils.Bool(true))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Data Lake Store File %q was not found (Account %q)", id.filePath, id.storageAccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Data Lake Store File %q (Account %q): %+v", id.filePath, id.storageAccountName, err)
	}

	d.Set("account_name", id.storageAccountName)
	d.Set("remote_file_path", id.filePath)

	return nil
}

func resourceArmDataLakeStoreFileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreFilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseDataLakeStoreFileId(d.Id(), client.AdlsFileSystemDNSSuffix)
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.storageAccountName, id.filePath, utils.Bool(false))
	if err != nil {
		if !response.WasNotFound(resp.Response.Response) {
			return fmt.Errorf("Error issuing delete request for Data Lake Store File %q (Account %q): %+v", id.filePath, id.storageAccountName, err)
		}
	}

	return nil
}

type dataLakeStoreFileId struct {
	storageAccountName string
	filePath           string
}

func parseDataLakeStoreFileId(input string, suffix string) (*dataLakeStoreFileId, error) {
	// Example: tomdevdls1.azuredatalakestore.net/test/example.txt
	// we add a scheme to the start of this so it parses correctly
	uri, err := url.Parse(fmt.Sprintf("https://%s", input))
	if err != nil {
		return nil, fmt.Errorf("Error parsing %q as URI: %+v", input, err)
	}

	// TODO: switch to pulling this from the Environment when it's available there
	// BUG: https://github.com/Azure/go-autorest/issues/312
	replacement := fmt.Sprintf(".%s", suffix)
	accountName := strings.Replace(uri.Host, replacement, "", -1)

	file := dataLakeStoreFileId{
		storageAccountName: accountName,
		filePath:           uri.Path,
	}
	return &file, nil
}

func validateDataLakeStoreRemoteFilePath() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (warnings []string, errors []error) {
		val := v.(string)

		if !strings.HasPrefix(val, "/") {
			errors = append(errors, fmt.Errorf("%q must start with `/`", k))
		}

		return warnings, errors
	}
}
