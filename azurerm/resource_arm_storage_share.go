package azurerm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmStorageShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageShareCreate,
		Read:   resourceArmStorageShareRead,
		Update: resourceArmStorageShareUpdate,
		Delete: resourceArmStorageShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		MigrateState:  resourceStorageShareMigrateState,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageShareName,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5120,
				ValidateFunc: validation.IntBetween(1, 5120),
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceArmStorageShareCreate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(waitCtx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	name := d.Get("name").(string)
	metaData := make(map[string]string) // TODO: support MetaData
	options := &storage.FileRequestOptions{}

	reference := fileClient.GetShareReference(name)
	exists, err := reference.Exists()
	if err != nil {
		return fmt.Errorf("Error checking if the Share %q already exists within Storage Account %q: %+v", name, storageAccountName, err)
	}

	id := fmt.Sprintf("%s/%s/%s", name, resourceGroupName, storageAccountName)
	if exists {
		return tf.ImportAsExistsError("azurerm_storage_share", id)
	}

	log.Printf("[INFO] Creating share %q in storage account %q", name, storageAccountName)
	err = reference.Create(options)

	log.Printf("[INFO] Setting share %q metadata in storage account %q", name, storageAccountName)
	reference.Metadata = metaData
	reference.SetMetadata(options)

	log.Printf("[INFO] Setting share %q properties in storage account %q", name, storageAccountName)
	reference.Properties = storage.ShareProperties{
		Quota: d.Get("quota").(int),
	}
	err = reference.SetProperties(options)
	if err != nil {
		return fmt.Errorf("Error setting properties on Share %q within Storage Account %q: %+v", name, storageAccountName, err)
	}

	d.SetId(id)
	return resourceArmStorageShareRead(d, meta)
}

func resourceArmStorageShareRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		return fmt.Errorf("ID was not in the expected format - expected `{name}/{resourceGroup}/{storageAccountName}` got %q", id)
	}
	name := id[0]
	resourceGroupName := id[1]
	storageAccountName := id[2]

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing file %q from state", storageAccountName, d.Id())
		d.SetId("")
		return nil
	}

	reference := fileClient.GetShareReference(name)
	exists, err := reference.Exists()
	if err != nil {
		return fmt.Errorf("Error testing existence of share %q: %s", name, err)
	}

	if !exists {
		log.Printf("[INFO] Share %q no longer exists, removing from state...", name)
		d.SetId("")
		return nil
	}

	url := reference.URL()
	if url == "" {
		log.Printf("[INFO] URL for %q is empty", name)
	}
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	d.Set("storage_account_name", storageAccountName)
	d.Set("url", url)

	reference.FetchAttributes(nil)
	d.Set("quota", reference.Properties.Quota)

	return nil
}

func resourceArmStorageShareUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		return fmt.Errorf("ID was not in the expected format - expected `{name}/{resourceGroup}/{storageAccountName}` got %q", id)
	}
	name := id[0]
	resourceGroupName := id[1]
	storageAccountName := id[2]

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(waitCtx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	options := &storage.FileRequestOptions{}

	reference := fileClient.GetShareReference(name)

	log.Printf("[INFO] Setting share %q properties in storage account %q", name, storageAccountName)
	reference.Properties = storage.ShareProperties{
		Quota: d.Get("quota").(int),
	}
	err = reference.SetProperties(options)
	if err != nil {
		return fmt.Errorf("Error setting properties on Share %q within Storage Account %q: %+v", name, storageAccountName, err)
	}

	return resourceArmStorageShareRead(d, meta)
}

func resourceArmStorageShareDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		return fmt.Errorf("ID was not in the expected format - expected `{name}/{resourceGroup}/{storageAccountName}` got %q", id)
	}
	name := id[0]
	resourceGroupName := id[1]
	storageAccountName := id[2]

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(waitCtx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[INFO] Storage Account %q doesn't exist so the file won't exist", storageAccountName)
		return nil
	}

	reference := fileClient.GetShareReference(name)
	options := &storage.FileRequestOptions{}

	if _, err = reference.DeleteIfExists(options); err != nil {
		return fmt.Errorf("Error deleting storage file %q: %s", name, err)
	}

	d.SetId("")
	return nil
}

//Following the naming convention as laid out in the docs https://msdn.microsoft.com/library/azure/dn167011.aspx
func validateArmStorageShareName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9a-z-]+$`).MatchString(value) {
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
	if regexp.MustCompile(`[-]{2,}`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q does not allow consecutive hyphens: %q", k, value))
	}
	return
}
