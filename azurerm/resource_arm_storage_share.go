package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

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

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	name := d.Get("name").(string)
	metaData := make(map[string]string) // TODO: support MetaData
	options := &storage.FileRequestOptions{}

	log.Printf("[INFO] Creating share %q in storage account %q", name, storageAccountName)
	reference := fileClient.GetShareReference(name)

	id := fmt.Sprintf("%s/%s/%s", name, resourceGroupName, storageAccountName)
	if requireResourcesToBeImported {
		exists, e := reference.Exists()
		if e != nil {
			return fmt.Errorf("Error checking if Share %q exists (Account %q / Resource Group %q): %s", name, storageAccountName, resourceGroupName, e)
		}

		if exists {
			return tf.ImportAsExistsError("azurerm_storage_share", id)
		}
	}

	err = reference.Create(options)
	if err != nil {
		return fmt.Errorf("Error creating Storage Share %q reference (storage account: %q) : %+v", name, storageAccountName, err)
	}

	log.Printf("[INFO] Setting share %q metadata in storage account %q", name, storageAccountName)
	reference.Metadata = metaData
	if err := reference.SetMetadata(options); err != nil {
		return fmt.Errorf("Error setting metadata on Storage Share %q: %+v", name, err)
	}

	log.Printf("[INFO] Setting share %q properties in storage account %q", name, storageAccountName)
	reference.Properties = storage.ShareProperties{
		Quota: d.Get("quota").(int),
	}
	if err := reference.SetProperties(options); err != nil {
		return fmt.Errorf("Error setting properties on Storage Share %q: %+v", name, err)
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

	if err := reference.FetchAttributes(nil); err != nil {
		return fmt.Errorf("Error fetching properties on Storage Share %q: %+v", name, err)
	}
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

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
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
	if err := reference.SetProperties(options); err != nil {
		return fmt.Errorf("Error setting properties on Storage Share %q: %+v", name, err)
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

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[INFO]Storage Account %q doesn't exist so the file won't exist", storageAccountName)
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
func validateArmStorageShareName(v interface{}, k string) (warnings []string, errors []error) {
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
	return warnings, errors
}
