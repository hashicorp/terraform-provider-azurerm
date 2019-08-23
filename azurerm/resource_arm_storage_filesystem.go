package azurerm

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	azauto "github.com/Azure/go-autorest/autorest/azure"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/dfs"
)

var (
	timeout = int32(5000)
	maxrst  = int32(1000)
)

func resourceArmStorageFilesystem() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmStorageFilesystemCreate,
		Read:          resourceArmStorageFilesystemRead,
		Delete:        resourceArmStorageFilesystemDelete,
		SchemaVersion: 1,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageFilesystemName,
			},
			"resource_group_name": azure.SchemaResourceGroupNameDeprecated(),
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

//Following the naming convention as laid out in the docs
func validateArmStorageFilesystemName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^\$root$|^[0-9a-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q: %q",
			k, value))
	}
	if len(value) < 3 || len(value) > 63 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 3 and 63 characters: %q", k, value))
	}
	if regexp.MustCompile(`^-`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot begin with a hyphen: %q", k, value))
	}
	return warnings, errors
}

func resourceArmStorageFilesystemCreate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	storageClient := armClient.storage
	ctx := armClient.StopContext

	filesystemName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Filesystem %q (Account %s): %s", filesystemName, accountName, err)
	}
	if resourceGroup == nil {
		return fmt.Errorf("Unable to locate Resource Group for Storage Filesystem %q (Account %s)", filesystemName, accountName)
	}

	client, err := storageClient.FileSystemsClient(ctx, *resourceGroup, accountName)
	if err != nil {
		return fmt.Errorf("Error building Filesystems Client: %s", err)
	}

	reqId, _ := uuid.NewUUID()
	resp, err := client.GetFilesystemProperties(ctx, accountName, filesystemName, reqId.String(), &timeout, "")
	id := fmt.Sprintf("https://%s.dfs.%s/%s", accountName, armClient.environment.StorageEndpointSuffix, filesystemName)
	if features.ShouldResourcesBeImported() {
		exists, err := ((*resp.Response).StatusCode != http.StatusNotFound), err
		if err != nil {
			return fmt.Errorf("Error checking for existence of existing Filesystem %q (Account %q / Resource Group %q): %+v", filesystemName, accountName, *resourceGroup, err)
		}

		if exists {
			return tf.ImportAsExistsError("azurerm_storage_filesystem", id)
		}
	}

	log.Printf("[INFO] Creating filesystem %q in storage account %q.", filesystemName, accountName)
	err = resource.Retry(120*time.Second, checkFilesystemIsCreated(ctx, client, accountName, filesystemName))
	if err != nil {
		return fmt.Errorf("Error creating filesystem %q in storage account %q: %s", filesystemName, accountName, err)
	}

	d.SetId(id)
	return resourceArmStorageFilesystemRead(d, meta)
}

// resourceAzureStorageFilesystemRead does all the necessary API calls to
// read the status of the storage filesystem off Azure.
func resourceArmStorageFilesystemRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	storageClient := armClient.storage
	ctx := armClient.StopContext

	id, err := parseStorageFilesystemID(d.Id(), armClient.environment)
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.storageAccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage FileSystem %q (Account %s): %s", id.filesystemName, id.storageAccountName, err)
	}
	if resourceGroup == nil {
		log.Printf("[DEBUG] Unable to locate Resource Group for Storage FileSystem %q (Account %s) - assuming removed & removing from state", id.filesystemName, id.storageAccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileSystemsClient(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}

	var filesystem *dfs.Filesystem
	var continuation string
	continuation = ""

	for {
		reqId, _ := uuid.NewUUID()
		resp, err := client.ListFilesystem(ctx, id.storageAccountName, id.filesystemName, continuation, &maxrst, reqId.String(), &timeout, "")
		if err != nil {
			return fmt.Errorf("Failed to retrieve storage resp in account %q: %s", id.filesystemName, err)
		}
		for _, c := range *resp.Filesystems {
			if *(c.Name) == id.filesystemName {
				filesystem = &c
				break
			}
		}

		if resp.Response.Header.Get("x-ms-continuation") == "" {
			break
		}

		continuation = resp.Response.Header.Get("x-ms-continuation")
	}

	if filesystem == nil {
		log.Printf("[INFO] Storage filesystem %q does not exist in account %q, removing from state...", id.filesystemName, id.storageAccountName)
		d.SetId("")
		return nil
	}

	d.Set("name", id.filesystemName)
	d.Set("storage_account_name", id.storageAccountName)
	d.Set("resource_group_name", resourceGroup)

	output := make(map[string]interface{})

	output["last_modified"] = filesystem.LastModified
	output["etag"] = filesystem.ETag

	if err := d.Set("properties", output); err != nil {
		return fmt.Errorf("Error setting `properties`: %+v", err)
	}

	return nil
}

// resourceAzureStorageFilesystemDelete does all the necessary API calls to
// delete a storage filesystem off Azure.
func resourceArmStorageFilesystemDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	storageClient := armClient.storage
	ctx := armClient.StopContext

	id, err := parseStorageFilesystemID(d.Id(), armClient.environment)
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.storageAccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage FileSystem %q (Account %s): %s", id.filesystemName, id.storageAccountName, err)
	}
	if resourceGroup == nil {
		return fmt.Errorf("Unable to locate Resource Group for Storage FileSystem %q (Account %s)", id.filesystemName, id.storageAccountName)
	}

	client, err := storageClient.FileSystemsClient(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return fmt.Errorf("Error building FileSystems Client: %s", err)
	}

	log.Printf("[INFO] Deleting storage filesystem %q in account %q", id.filesystemName, id.storageAccountName)

	reqId, _ := uuid.NewUUID()
	if _, err := client.DeleteFilesystem(ctx, id.storageAccountName, id.filesystemName, "", "", reqId.String(), &timeout, ""); err != nil {
		return fmt.Errorf("Error deleting storage filesystem %q from storage account %q: %s", id.filesystemName, id.storageAccountName, err)
	}

	return nil
}

func checkFilesystemIsCreated(ctx context.Context, filesystemClient *dfs.Client, accountName string, name string) func() *resource.RetryError {
	return func() *resource.RetryError {
		if _, err := createIfNotExists(ctx, filesystemClient, accountName, name); err != nil {
			return resource.RetryableError(err)
		}

		return nil
	}
}

// CreateIfNotExists creates a storage datalake gen2 filesystem if it does not exist. Returns
// true if filesystem is newly created or false if filesystem already exists.
func createIfNotExists(ctx context.Context, filesystemClient *dfs.Client, accountName string, name string) (bool, error) {
	reqId, _ := uuid.NewUUID()
	resp, err := filesystemClient.CreateFilesystem(ctx, accountName, name, "", reqId.String(), &timeout, "")
	if resp.Response.Response != nil {
		defer drainRespBody(resp.Response.Response)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusConflict {
			return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated, nil
		}
	}
	return false, err
}

func drainRespBody(resp *http.Response) {
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

type storageFilesystemId struct {
	id                 string
	storageAccountName string
	filesystemName     string
}

func parseStorageFilesystemID(input string, environment azauto.Environment) (*storageFilesystemId, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("Error parsing %q as URI: %+v", input, err)
	}

	// remove the leading `/`
	segments := strings.Split(strings.TrimPrefix(uri.Path, "/"), "/")
	if len(segments) < 1 {
		return nil, fmt.Errorf("Expected number of segments in the path to be < 1 but got %d", len(segments))
	}

	storageAccountName := strings.Replace(uri.Host, fmt.Sprintf(".dfs.%s", environment.StorageEndpointSuffix), "", 1)
	filesystemName := segments[0]

	id := storageFilesystemId{
		id:                 input,
		storageAccountName: storageAccountName,
		filesystemName:     filesystemName,
	}
	return &id, nil
}
